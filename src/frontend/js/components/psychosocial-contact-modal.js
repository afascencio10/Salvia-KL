/**
 * psychosocial-contact-modal
 * Flujo de modales del esquema 3x3 de Atención Psicosocial. Autónomo (no usa
 * follow_up_v2). Invocado por `ref` desde <psychosocial-contact-card> vía start(process).
 *
 *  ¿Contestó? ──No──> registra intento fallido (con nota) y cierra (o abre Acciones al 3º del día)
 *             ──Sí──> registra contacto exitoso ──> Consentimiento
 *                        ──No acepta──> (formulario de cierre — futuro)
 *                        ──Sí acepta──> Sesión (inmediata / agendada)
 *  Acciones (3 fallidos/día): 3a próximo intento · 3b añadir intento (<50) · 3c cierre (9 en 3 días)
 */
(function injectPsychosocialContactModalStyles() {
    if (document.getElementById('psc-modal-styles')) return;
    var style = document.createElement('style');
    style.id = 'psc-modal-styles';
    style.textContent = `
        .psc-modal-overlay { position: fixed; inset: 0; background: rgba(17,24,39,.55); display: flex; align-items: center; justify-content: center; z-index: 1050; padding: 1rem; }
        .psc-modal-content { background: #fff; border-radius: 1rem; padding: 2rem; max-width: 28rem; width: 100%; text-align: center; box-shadow: 0 20px 25px -5px rgba(0,0,0,.1); }
        .psc-icon { width: 4rem; height: 4rem; border-radius: 9999px; display: flex; align-items: center; justify-content: center; margin: 0 auto 1.25rem; font-size: 1.5rem; }
        .psc-title { font-size: 1.25rem; font-weight: 700; color: #111827; margin-bottom: .25rem; }
        .psc-sub { font-size: .875rem; color: #6b7280; margin-bottom: .25rem; }
        .psc-name { font-size: 1rem; font-weight: 700; color: #111827; margin-bottom: 1.25rem; }
        .psc-banner { background: #fee2e2; border: 1px solid #ef4444; border-radius: .5rem; padding: .6rem; margin-bottom: 1.25rem; color: #b91c1c; font-size: .8125rem; font-weight: 600; }
        .psc-label { display: block; text-align: left; font-size: .8125rem; font-weight: 700; color: #374151; margin-bottom: .375rem; }
        .psc-input { width: 100%; padding: .7rem .9rem; border-radius: .6rem; border: 1px solid #e5e7eb; font-size: .875rem; box-sizing: border-box; font-family: inherit; margin-bottom: 1.1rem; }
        .psc-btn { padding: .8rem 1.1rem; border-radius: .7rem; font-weight: 600; font-size: .875rem; cursor: pointer; border: none; }
        .psc-btn:disabled { opacity: .55; cursor: not-allowed; }
        .psc-btn-primary { background: #8b5cf6; color: #fff; }
        .psc-btn-ghost { background: #fff; border: 1px solid #e5e7eb; color: #4b5563; }
        .psc-btn-danger { background: #dc2626; color: #fff; }
        .psc-row { display: flex; gap: .75rem; }
        .psc-row > .psc-btn { flex: 1; }
        .psc-cancel { margin-top: 1.25rem; background: none; border: none; font-size: .75rem; color: #9ca3af; font-weight: 600; cursor: pointer; }
        .psc-stack { display: flex; flex-direction: column; gap: .6rem; }
    `;
    document.head.appendChild(style);
})();

app.component('psychosocial-contact-modal', {
    delimiters: ['${', '}'],
    props: {
        currentUser:   { type: String, default: '' },
        currentUserId: { type: String, default: '' },
        userTeam:      { type: String, default: '' }
    },
    data: function() {
        return {
            proc: null,
            modalContesto: false,
            modalConsentimiento: false,
            modalSesion: false,
            modalAcciones: false,
            modalNextAttempt: false,
            modalCierre: false,
            closureLoading: false,
            closureFormId: null,
            closureSubmissionId: null,
            attemptDateTime: '',
            note: '',
            scheduledDate: '',
            scheduledTime: '',
            nextAttemptDateTime: '',
            lastSuccessAttemptId: null,
            loading: false
        };
    },
    computed: {
        // El modal de acciones se abre por dos disparadores: umbral diario o
        // elegibilidad de cierre. El texto se adapta usando los conteos reales.
        accionesPorDia: function() {
            return this.proc && (this.proc.dailyFailedCount || 0) >= 3;
        },
        accionesTitle: function() {
            return this.accionesPorDia ? 'Límite de intentos del día' : 'Cierre por imposibilidad disponible';
        },
        accionesSub: function() {
            if (!this.proc) return '';
            if (this.accionesPorDia) {
                var n = this.proc.dailyFailedCount || 0;
                return 'Se registraron ' + n + ' intento' + (n === 1 ? '' : 's') + ' fallido' + (n === 1 ? '' : 's') + ' hoy. ¿Qué deseas hacer?';
            }
            var total = this.proc.totalCount || 0;
            var dias = this.proc.distinctDaysCount || 0;
            return 'Este proceso acumula ' + total + ' intento' + (total === 1 ? '' : 's') +
                ' en ' + dias + ' día' + (dias === 1 ? '' : 's') + ' distinto' + (dias === 1 ? '' : 's') +
                ', suficiente para cerrar por imposibilidad de contacto. ¿Qué deseas hacer?';
        }
    },
    methods: {
        // ── API pública ────────────────────────────────────────────────
        start: function(process) {
            this.proc = process || {};
            this.note = '';
            this.initAttemptDateTime();
            // Abre directo el modal de acciones si (a) hay 3 fallidos del día (RN-04)
            // o (b) el proceso ya es elegible para cierre por imposibilidad (RN-06):
            // en ese caso la opción de cierre debe estar siempre a la vista.
            if ((this.proc.dailyFailedCount || 0) >= 3 || this.proc.closureEligible) {
                this.modalAcciones = true;
            } else {
                this.modalContesto = true;
            }
        },

        // ── Helpers ────────────────────────────────────────────────────
        initAttemptDateTime: function() {
            var tz = (new Date()).getTimezoneOffset() * 60000;
            this.attemptDateTime = (new Date(Date.now() - tz)).toISOString().slice(0, 16);
        },
        toISO: function(localValue) {
            if (!localValue) return null;
            var p = localValue.split(/[-T:]/);
            if (p.length < 5) return null;
            return new Date(+p[0], +p[1] - 1, +p[2], +p[3], +p[4]).toISOString();
        },
        closeAll: function() {
            this.modalContesto = this.modalConsentimiento = this.modalSesion = false;
            this.modalAcciones = this.modalNextAttempt = this.modalCierre = false;
            this.proc = null; this.note = ''; this.scheduledDate = ''; this.scheduledTime = '';
            this.nextAttemptDateTime = ''; this.lastSuccessAttemptId = null; this.loading = false;
            this.closureFormId = null; this.closureSubmissionId = null; this.closureLoading = false;
        },
        flashSuccess: function() {
            var ov = document.getElementById('successOverlay');
            if (ov) { ov.style.display = 'flex'; setTimeout(function(){ ov.style.display = 'none'; }, 1800); }
        },
        showError: function(msg) {
            if (window.Swal) { Swal.fire({ icon: 'error', title: 'Error', text: msg }); }
            else { alert(msg); }
        },
        futureNote: function(msg) {
            if (window.Swal) { Swal.fire({ icon: 'info', title: 'Pendiente', text: msg }); }
            else { alert(msg); }
        },

        // ── ¿Contestó? ─────────────────────────────────────────────────
        handleYesAnswer: function() {
            var self = this;
            if (this.loading) return;
            this.loading = true;
            this.postAttempt(true).then(function(data) {
                self.loading = false;
                if (!data) return;
                self.lastSuccessAttemptId = data.attempt ? data.attempt.id : null;
                self.$emit('completed', { type: 'contacted', psicosocialId: self.proc.id });
                self.modalContesto = false;
                self.modalConsentimiento = true;
            });
        },
        handleNoAnswer: function() {
            var self = this;
            if (this.loading) return;
            this.loading = true;
            this.postAttempt(false).then(function(data) {
                self.loading = false;
                if (!data) return;
                self.$emit('completed', { type: 'attempt', psicosocialId: self.proc.id, counters: data.counters });
                var counters = data.counters || {};
                if (counters.dailyThresholdReached) {
                    // Refrescar el proc con los contadores nuevos y abrir Acciones
                    self.proc.dailyFailedCount = counters.dailyFailedCount;
                    self.proc.totalCount = counters.totalCount;
                    self.proc.distinctDaysCount = counters.distinctDaysCount;
                    self.proc.closureEligible = counters.closureEligible;
                    self.modalContesto = false;
                    self.modalAcciones = true;
                } else {
                    self.closeAll();
                    self.flashSuccess();
                }
            });
        },
        postAttempt: function(wasAnswered) {
            var self = this;
            var body = { was_answered: wasAnswered, note: this.note || null, attempt_at: this.toISO(this.attemptDateTime) };
            return fetch('/api/v1/psychosocial/' + encodeURIComponent(this.proc.id) + '/contact-attempts', {
                method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(body)
            }).then(function(res) {
                return res.json().then(function(json) {
                    if (!res.ok) { self.showError(json.error || 'No se pudo registrar el intento'); return null; }
                    return json;
                });
            }).catch(function(e) { self.showError(e.message); return null; });
        },

        // ── Consentimiento ─────────────────────────────────────────────
        handleConsent: function(accepted) {
            var self = this;
            if (this.loading || !this.lastSuccessAttemptId) { if(!this.lastSuccessAttemptId) this.showError('No hay intento exitoso asociado'); return; }
            this.loading = true;
            fetch('/api/v1/contact-attempts/' + encodeURIComponent(this.lastSuccessAttemptId) + '/consent', {
                method: 'PATCH', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify({ consent_given: accepted })
            }).then(function(res) {
                return res.json().then(function(json) {
                    self.loading = false;
                    if (!res.ok) { self.showError(json.error || 'Error al registrar consentimiento'); return; }
                    self.modalConsentimiento = false;
                    if (accepted) {
                        self.modalSesion = true;
                    } else {
                        self.openClosureForm('no_consentimiento');
                    }
                });
            }).catch(function(e) { self.loading = false; self.showError(e.message); });
        },

        // ── Sesión ─────────────────────────────────────────────────────
        scheduleSession: function(immediate) {
            var self = this;
            if (this.loading) return;
            if (!immediate && (!this.scheduledDate || !this.scheduledTime)) { this.showError('Selecciona fecha y hora'); return; }
            this.loading = true;
            var body = immediate
                ? { immediate: true }
                : { immediate: false, scheduled_date: this.scheduledDate, scheduled_time: this.scheduledTime };
            fetch('/api/v1/psychosocial/' + encodeURIComponent(this.proc.id) + '/sessions', {
                method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(body)
            }).then(function(res) {
                return res.json().then(function(json) {
                    self.loading = false;
                    if (!res.ok) { self.showError(json.error || 'Error al agendar la sesión'); return; }
                    self.$emit('completed', { type: 'scheduled', psicosocialId: self.proc.id, sessionId: json.session ? json.session.id : null });
                    if (immediate && json.redirectUrl) {
                        self.futureNote('Sesión inmediata creada. El formulario de atención está pendiente de implementar (' + json.redirectUrl + ').');
                    }
                    self.closeAll();
                    self.flashSuccess();
                });
            }).catch(function(e) { self.loading = false; self.showError(e.message); });
        },

        // ── Acciones (3 fallidos/día) ──────────────────────────────────
        openNextAttempt: function() {
            var tz = (new Date()).getTimezoneOffset() * 60000;
            this.nextAttemptDateTime = (new Date(Date.now() - tz + 86400000)).toISOString().slice(0, 16);
            this.modalAcciones = false;
            this.modalNextAttempt = true;
        },
        saveNextAttempt: function() {
            var self = this;
            var iso = this.toISO(this.nextAttemptDateTime);
            if (!iso) { this.showError('Selecciona fecha y hora'); return; }
            this.loading = true;
            fetch('/api/v1/psychosocial/' + encodeURIComponent(this.proc.id) + '/next-attempt', {
                method: 'PUT', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify({ next_contact_attempt_at: iso })
            }).then(function(res) {
                return res.json().then(function(json) {
                    self.loading = false;
                    if (!res.ok) { self.showError(json.error || 'Error al fijar el próximo intento'); return; }
                    self.$emit('completed', { type: 'next-attempt', psicosocialId: self.proc.id, nextContactAttemptAt: json.nextContactAttemptAt });
                    self.closeAll();
                    self.flashSuccess();
                });
            }).catch(function(e) { self.loading = false; self.showError(e.message); });
        },
        addAttempt: function() {
            if ((this.proc.totalCount || 0) >= 50) return;
            this.modalAcciones = false;
            this.initAttemptDateTime();
            this.note = '';
            this.modalContesto = true;
        },
        closeByImpossibility: function() {
            this.modalAcciones = false;
            this.openClosureForm('imposibilidad_contacto_3x3');
        },

        // ── Cierre (formulario dinámico) ───────────────────────────────
        openClosureForm: function(reason) {
            var self = this;
            if (this.closureLoading) return;
            // Abrir el modal de inmediato con loader; el form se carga en segundo plano.
            this.closureLoading = true;
            this.closureFormId = null;
            this.closureSubmissionId = null;
            this.modalConsentimiento = false;
            this.modalAcciones = false;
            this.modalCierre = true;
            fetch('/api/v1/psychosocial/' + encodeURIComponent(this.proc.id) + '/init-closure-form?reason=' + encodeURIComponent(reason), {
                method: 'POST', headers: { 'Content-Type': 'application/json' }
            }).then(function(res) {
                return res.json().then(function(json) {
                    self.closureLoading = false;
                    if (!res.ok) { self.modalCierre = false; self.showError(json.error || 'No se pudo iniciar el formulario de cierre'); return; }
                    self.closureFormId = json.form_id;
                    self.closureSubmissionId = json.submission_id;
                });
            }).catch(function(e) { self.closureLoading = false; self.modalCierre = false; self.showError(e.message); });
        },
        onClosureCompleted: function() {
            var self = this;
            var psId = this.proc ? this.proc.id : null;
            fetch('/api/v1/psychosocial/' + encodeURIComponent(psId) + '/close', {
                method: 'POST', headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ submission_id: this.closureSubmissionId })
            }).then(function(res) {
                return res.json().then(function(json) {
                    if (!res.ok) { self.showError(json.error || 'No se pudo cerrar el proceso'); return; }
                    self.$emit('completed', { type: 'closed', psicosocialId: psId, status: json.status, motivo: json.motivo });
                    self.closeAll();
                    self.flashSuccess();
                });
            }).catch(function(e) { self.showError(e.message); });
        }
    },
    template: `
    <div>
      <!-- ¿Contestó? -->
      <div v-if="modalContesto" class="psc-modal-overlay" @click.self="closeAll">
        <div class="psc-modal-content">
          <div class="psc-icon" style="background:#f5f3ff;color:#8b5cf6;"><i class="fas fa-phone-alt"></i></div>
          <h3 class="psc-title">¿Contestó la llamada?</h3>
          <p class="psc-sub">Estamos contactando a</p>
          <p class="psc-name">\${ proc && proc.caseName ? proc.caseName : 'la ciudadana' }</p>
          <div v-if="proc && proc.totalCount > 0" class="psc-banner">
            <i class="fas fa-exclamation-triangle"></i> \${ proc.totalCount } intento\${ proc.totalCount > 1 ? 's' : '' } previo\${ proc.totalCount > 1 ? 's' : '' }
          </div>
          <label class="psc-label"><i class="far fa-calendar-alt"></i> Fecha y hora del intento</label>
          <input type="datetime-local" class="psc-input" v-model="attemptDateTime" aria-label="Fecha y hora del intento">
          <label class="psc-label">Nota (opcional)</label>
          <textarea class="psc-input" rows="2" v-model="note" placeholder="Ej: No contestó, celular apagado / Sí contestó, pide que lo llamen otro día" aria-label="Nota"></textarea>
          <div class="psc-row">
            <button class="psc-btn psc-btn-ghost" :disabled="loading" @click="handleNoAnswer">No contestó</button>
            <button class="psc-btn psc-btn-primary" :disabled="loading" @click="handleYesAnswer">Sí contestó <i class="fas fa-check"></i></button>
          </div>
          <button class="psc-cancel" :disabled="loading" @click="closeAll">CANCELAR</button>
        </div>
      </div>

      <!-- Consentimiento -->
      <div v-if="modalConsentimiento" class="psc-modal-overlay" @click.self="closeAll">
        <div class="psc-modal-content">
          <div class="psc-icon" style="background:#eff6ff;color:#2563eb;"><i class="fas fa-file-signature"></i></div>
          <h3 class="psc-title">Consentimiento Informado</h3>
          <p class="psc-sub" style="margin-bottom:1.5rem;">¿La ciudadana acepta continuar con la atención psicosocial?</p>
          <div class="psc-row">
            <button class="psc-btn psc-btn-ghost" :disabled="loading" @click="handleConsent(false)">No acepta</button>
            <button class="psc-btn psc-btn-primary" :disabled="loading" @click="handleConsent(true)">Sí acepta</button>
          </div>
          <button class="psc-cancel" :disabled="loading" @click="closeAll">CANCELAR</button>
        </div>
      </div>

      <!-- Sesión -->
      <div v-if="modalSesion" class="psc-modal-overlay" @click.self="closeAll">
        <div class="psc-modal-content">
          <div class="psc-icon" style="background:#f0fdf4;color:#16a34a;"><i class="far fa-calendar-check"></i></div>
          <h3 class="psc-title">¿Cuándo se realiza la sesión?</h3>
          <p class="psc-sub" style="margin-bottom:1.25rem;">Puedes iniciarla de inmediato o agendarla</p>
          <div class="psc-stack">
            <button class="psc-btn psc-btn-primary" :disabled="loading" @click="scheduleSession(true)">De inmediato</button>
            <div style="border-top:1px solid #e5e7eb;margin:.4rem 0;"></div>
            <label class="psc-label">Agendar — Fecha</label>
            <input type="date" class="psc-input" v-model="scheduledDate" aria-label="Fecha de la sesión">
            <label class="psc-label">Agendar — Hora</label>
            <input type="time" class="psc-input" v-model="scheduledTime" aria-label="Hora de la sesión">
            <button class="psc-btn psc-btn-ghost" :disabled="loading" @click="scheduleSession(false)">Agendar sesión</button>
          </div>
          <button class="psc-cancel" :disabled="loading" @click="closeAll">CANCELAR</button>
        </div>
      </div>

      <!-- Acciones -->
      <div v-if="modalAcciones" class="psc-modal-overlay" @click.self="closeAll">
        <div class="psc-modal-content">
          <div class="psc-icon" :style="accionesPorDia ? 'background:#fee2e2;color:#dc2626;' : 'background:#ede9fe;color:#7c3aed;'"><i class="fas fa-exclamation-circle"></i></div>
          <h3 class="psc-title">\${ accionesTitle }</h3>
          <p class="psc-sub" style="margin-bottom:1.25rem;">\${ accionesSub }</p>
          <div class="psc-stack">
            <button class="psc-btn psc-btn-primary" :disabled="loading" @click="openNextAttempt"><i class="far fa-calendar-alt"></i> Fijar próximo intento</button>
            <button class="psc-btn psc-btn-ghost" :disabled="loading || (proc && proc.totalCount >= 50)" @click="addAttempt"><i class="fas fa-plus-circle"></i> Añadir intento \${ proc && proc.totalCount >= 50 ? '(tope 50)' : '' }</button>
            <button v-if="proc && proc.closureEligible" class="psc-btn psc-btn-danger" :disabled="loading" @click="closeByImpossibility"><i class="fas fa-times-circle"></i> Cerrar por imposibilidad</button>
          </div>
          <button class="psc-cancel" :disabled="loading" @click="closeAll">VOLVER</button>
        </div>
      </div>

      <!-- Próximo intento (3a) -->
      <div v-if="modalNextAttempt" class="psc-modal-overlay" @click.self="closeAll">
        <div class="psc-modal-content">
          <div class="psc-icon" style="background:#f5f3ff;color:#8b5cf6;"><i class="far fa-calendar-plus"></i></div>
          <h3 class="psc-title">Próximo intento de contacto</h3>
          <p class="psc-sub" style="margin-bottom:1.25rem;">Elige la fecha y hora exactas</p>
          <input type="datetime-local" class="psc-input" v-model="nextAttemptDateTime" aria-label="Próximo intento">
          <div class="psc-row">
            <button class="psc-btn psc-btn-ghost" :disabled="loading" @click="closeAll">Cancelar</button>
            <button class="psc-btn psc-btn-primary" :disabled="loading" @click="saveNextAttempt">Guardar</button>
          </div>
        </div>
      </div>

      <!-- Formulario de Cierre de proceso psicosocial (dinamic-form) -->
      <div v-if="modalCierre" class="psc-modal-overlay" @click.self="closeAll">
        <div class="psc-modal-content" style="max-width:900px; width:95%; max-height:90vh; overflow:auto; text-align:left;">
          <div style="display:flex; justify-content:space-between; align-items:center; margin-bottom:1rem; border-bottom:1px solid #e5e7eb; padding-bottom:.75rem;">
            <h3 class="psc-title" style="margin:0;">Formulario de Cierre de proceso psicosocial</h3>
            <button @click="closeAll" style="background:none;border:none;font-size:1.25rem;color:#9ca3af;cursor:pointer;"><i class="fas fa-times"></i></button>
          </div>
          <!-- Loader idéntico al del dinamic-form para que el cambio sea invisible -->
          <div v-if="!closureFormId || !closureSubmissionId" style="display:flex;align-items:center;justify-content:center;width:100%;padding:48px 0;gap:12px;color:#6b7280;font-size:14px;">
            <svg style="width:20px;height:20px;animation:df-spin 0.8s linear infinite;flex-shrink:0" viewBox="0 0 24 24" fill="none">
              <circle cx="12" cy="12" r="10" stroke="#e5e7eb" stroke-width="3"/>
              <path d="M12 2a10 10 0 0 1 10 10" stroke="#7c3aed" stroke-width="3" stroke-linecap="round"/>
            </svg>
            Cargando formulario...
          </div>
          <dinamic-form
            v-else
            :form-id="closureFormId"
            :submission-id="closureSubmissionId"
            @form-completed="onClosureCompleted">
          </dinamic-form>
        </div>
      </div>
    </div>
    `
});
