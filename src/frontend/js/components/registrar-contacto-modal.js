/**
 * registrar-contacto-modal
 * Modal para registrar un contacto/sesión en una remisión psicosocial.
 * Se abre vía: this.$refs.registrarContactoModal.open(psicosocialId)
 * Emite @saved cuando se guarda exitosamente.
 */

(function injectRegistrarContactoStyles() {
    if (document.getElementById('rcm-styles')) return;
    var style = document.createElement('style');
    style.id = 'rcm-styles';
    style.textContent = `
        .rcm-overlay { position:fixed;inset:0;background:rgba(0,0,0,.45);display:flex;align-items:center;justify-content:center;z-index:10000;padding:16px }
        .rcm-modal { background:#fff;border-radius:16px;width:100%;max-width:500px;max-height:90vh;overflow-y:auto;box-shadow:0 20px 40px rgba(0,0,0,.18);padding:32px }
        .rcm-title { font-size:1.2rem;font-weight:800;color:#111827;margin:0 0 6px }
        .rcm-subtitle { font-size:.85rem;color:#6b7280;margin:0 0 24px }
        .rcm-row { display:flex;gap:14px;margin-bottom:20px }
        .rcm-field { flex:1 }
        .rcm-label { font-size:.78rem;font-weight:600;color:#374151;display:block;margin-bottom:6px }
        .rcm-input { width:100%;padding:10px 14px;border:1px solid #d1d5db;border-radius:10px;font-size:.85rem;outline:none }
        .rcm-input:focus { border-color:#5106A7 }
        .rcm-select { width:100%;padding:10px 14px;border:1px solid #d1d5db;border-radius:10px;font-size:.85rem;outline:none;appearance:auto }
        .rcm-textarea { width:100%;padding:10px 14px;border:1px solid #d1d5db;border-radius:10px;font-size:.85rem;outline:none;resize:vertical;min-height:80px }
        .rcm-toggle-group { display:flex;gap:10px;margin-bottom:20px;flex-wrap:wrap }
        .rcm-toggle-btn { padding:10px 18px;border-radius:10px;border:1.5px solid #e5e7eb;background:#fff;font-size:.82rem;font-weight:600;cursor:pointer;display:flex;align-items:center;gap:8px;transition:all .15s;color:#374151 }
        .rcm-toggle-btn:hover { border-color:#c4b5fd }
        .rcm-toggle-btn.active-contacto { border-color:#5106A7;background:#ede9fe;color:#5106A7 }
        .rcm-toggle-btn.active-agendar { border-color:#5106A7;background:#f3e8ff;color:#5106A7 }
        .rcm-toggle-btn.active-realizar { border-color:#16a34a;background:#dcfce7;color:#15803d }
        .rcm-footer { display:flex;gap:12px;justify-content:flex-end;margin-top:28px }
        .rcm-btn-cancel { padding:10px 24px;border-radius:10px;border:1px solid #e5e7eb;background:#fff;color:#374151;font-size:.85rem;font-weight:600;cursor:pointer }
        .rcm-btn-save { padding:10px 24px;border-radius:10px;border:none;background:#c4b5fd;color:#1f2937;font-size:.85rem;font-weight:700;cursor:pointer }
        .rcm-btn-save:hover { background:#a78bfa }
    `;
    document.head.appendChild(style);
})();

app.component('registrar-contacto-modal', {
    delimiters: ['${', '}'],
    emits: ['saved'],

    data() {
        return {
            visible: false,
            psicosocialId: null,
            guardando: false,

            form: {
                fecha: '',
                hora: '',
                tipo: null, // 'contacto' | 'agendar' | 'realizar'
                nota: '',
                sessionType: '',
                sessionFecha: '',
                sessionHora: '',
            },
        };
    },

    computed: {
        puedeGuardar() {
            if (!this.form.fecha || !this.form.hora || !this.form.tipo) return false;
            if (this.form.tipo === 'contacto' && !this.form.nota) return false;
            if (this.form.tipo === 'agendar' && (!this.form.sessionType || !this.form.sessionFecha || !this.form.sessionHora)) return false;
            if (this.form.tipo === 'realizar' && !this.form.sessionType) return false;
            return true;
        },
    },

    methods: {
        open(psicosocialId) {
            this.psicosocialId = psicosocialId;
            this.visible = true;
            this.guardando = false;
            var now = new Date();
            this.form = {
                fecha: now.toISOString().substring(0, 10),
                hora: now.toTimeString().substring(0, 5),
                tipo: null,
                nota: '',
                sessionType: '',
                sessionFecha: now.toISOString().substring(0, 10),
                sessionHora: now.toTimeString().substring(0, 5),
            };
        },

        cerrar() {
            this.visible = false;
            this.psicosocialId = null;
        },

        selectTipo(tipo) {
            this.form.tipo = tipo;
        },

        async guardar() {
            if (!this.puedeGuardar || this.guardando) return;
            this.guardando = true;

            var payload = {
                psicosocialId: this.psicosocialId,
                contactDate: this.form.fecha,
                contactTime: this.form.hora,
                type: this.form.tipo,
            };

            if (this.form.tipo === 'contacto') {
                payload.summary = this.form.nota;
            } else if (this.form.tipo === 'agendar') {
                payload.sessionType = this.form.sessionType;
                payload.scheduledDate = this.form.sessionFecha;
                payload.scheduledTime = this.form.sessionHora;
            } else if (this.form.tipo === 'realizar') {
                payload.sessionType = this.form.sessionType;
            }

            try {
                var res = await fetch('/api/v1/psychosocial-support/' + this.psicosocialId + '/contacts', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify(payload),
                });
                if (!res.ok) throw new Error('Error ' + res.status);
                this.$emit('saved');
                this.cerrar();
            } catch (e) {
                console.error('[registrar-contacto-modal] Error:', e);
            } finally {
                this.guardando = false;
            }
        },
    },

    template: `
<div v-if="visible" class="rcm-overlay" @click.self="cerrar">
    <div class="rcm-modal">
        <p class="rcm-title">Registrar contacto</p>
        <p class="rcm-subtitle">Registra lo que ocurrió en este contacto con la paciente.</p>

        <!-- Fecha y hora -->
        <div class="rcm-row">
            <div class="rcm-field">
                <label class="rcm-label">Fecha del contacto</label>
                <input type="date" class="rcm-input" v-model="form.fecha" />
            </div>
            <div class="rcm-field">
                <label class="rcm-label">Hora</label>
                <input type="time" class="rcm-input" v-model="form.hora" />
            </div>
        </div>

        <!-- Toggle: ¿Qué ocurrió? -->
        <p class="rcm-label" style="margin-bottom:10px">¿Qué ocurrió en este contacto?</p>
        <div class="rcm-toggle-group">
            <button type="button" :class="['rcm-toggle-btn', form.tipo === 'contacto' ? 'active-contacto' : '']" @click="selectTipo('contacto')">
                💬 Solo contacto
            </button>
            <button type="button" :class="['rcm-toggle-btn', form.tipo === 'agendar' ? 'active-agendar' : '']" @click="selectTipo('agendar')">
                📅 Agendé sesión
            </button>
            <button type="button" :class="['rcm-toggle-btn', form.tipo === 'realizar' ? 'active-realizar' : '']" @click="selectTipo('realizar')">
                ✅ Realicé sesión
            </button>
        </div>

        <!-- Solo contacto → Nota -->
        <div v-if="form.tipo === 'contacto'" style="margin-bottom:16px">
            <label class="rcm-label">Nota</label>
            <textarea class="rcm-textarea" v-model="form.nota" placeholder="Ej: Llamada breve para confirmar cita, paciente reporta estar bien..." maxlength="255"></textarea>
        </div>

        <!-- Agendé sesión → Tipo + Fecha/Hora sesión -->
        <template v-if="form.tipo === 'agendar'">
            <div style="margin-bottom:16px">
                <label class="rcm-label">Tipo de sesión</label>
                <select class="rcm-select" v-model="form.sessionType">
                    <option value="" disabled>Seleccione...</option>
                    <option value="PRIMER_CONTACTO">Primer contacto</option>
                    <option value="PRIMERA_ATENCION">Primera atención</option>
                    <option value="SEGUIMIENTO">Seguimiento</option>
                    <option value="CIERRE">Cierre</option>
                </select>
            </div>
            <div class="rcm-row">
                <div class="rcm-field">
                    <label class="rcm-label">Fecha de la sesión</label>
                    <input type="date" class="rcm-input" v-model="form.sessionFecha" />
                </div>
                <div class="rcm-field">
                    <label class="rcm-label">Hora de la sesión</label>
                    <input type="time" class="rcm-input" v-model="form.sessionHora" />
                </div>
            </div>
        </template>

        <!-- Realicé sesión → Solo tipo -->
        <div v-if="form.tipo === 'realizar'" style="margin-bottom:16px">
            <label class="rcm-label">Tipo de sesión</label>
            <select class="rcm-select" v-model="form.sessionType">
                <option value="" disabled>Seleccione...</option>
                <option value="PRIMER_CONTACTO">Primer contacto</option>
                <option value="PRIMERA_ATENCION">Primera atención</option>
                <option value="SEGUIMIENTO">Seguimiento</option>
                <option value="CIERRE">Cierre</option>
            </select>
        </div>

        <!-- Footer -->
        <div class="rcm-footer">
            <button class="rcm-btn-cancel" @click="cerrar">Cancelar</button>
            <button class="rcm-btn-save" :disabled="!puedeGuardar || guardando" :style="{ opacity: puedeGuardar && !guardando ? '1' : '.5' }" @click="guardar">
                \${ guardando ? 'Guardando...' : 'Guardar' }
            </button>
        </div>
    </div>
</div>
    `,
});
