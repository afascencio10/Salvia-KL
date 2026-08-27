/**
 * registrar-contacto-modal
 * Modal para registrar un contacto/sesión en una remisión psicosocial.
 * Ofrece dos caminos:
 *   1. Contacto inmediato → crea con fecha actual y redirige al formulario psicosocial
 *   2. Programar contacto → pide fecha/hora, valida disponibilidad (2h) y crea agendado
 *
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
        .rcm-input { width:100%;padding:10px 14px;border:1px solid #d1d5db;border-radius:10px;font-size:.85rem;outline:none;box-sizing:border-box }
        .rcm-input:focus { border-color:#5106A7 }
        .rcm-footer { display:flex;gap:12px;justify-content:flex-end;margin-top:28px }
        .rcm-btn-cancel { padding:10px 24px;border-radius:10px;border:1px solid #e5e7eb;background:#fff;color:#374151;font-size:.85rem;font-weight:600;cursor:pointer }
        .rcm-btn-save { padding:10px 24px;border-radius:10px;border:none;background:#5106A7;color:#fff;font-size:.85rem;font-weight:700;cursor:pointer }
        .rcm-btn-save:hover { background:#3b0764 }
        .rcm-btn-save:disabled { opacity:.5;cursor:not-allowed }
        .rcm-option-card { border:2px solid #e5e7eb;border-radius:12px;padding:16px 20px;cursor:pointer;transition:all .15s;margin-bottom:12px }
        .rcm-option-card:hover { border-color:#c4b5fd }
        .rcm-option-card.active { border-color:#5106A7;background:#f5f3ff }
        .rcm-option-title { font-size:.9rem;font-weight:700;color:#111827;margin:0 0 4px }
        .rcm-option-desc { font-size:.78rem;color:#6b7280;margin:0 }
        .rcm-alert { background:#fef2f2;border:1px solid #fecaca;border-radius:8px;padding:10px 14px;font-size:.8rem;color:#dc2626;margin-bottom:16px }
        .rcm-success { background:#f0fdf4;border:1px solid #bbf7d0;border-radius:8px;padding:10px 14px;font-size:.8rem;color:#15803d;margin-bottom:16px }
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
            modo: null, // 'inmediato' | 'programado'
            // Campos para programado
            fecha: '',
            hora: '',
            // Fecha mínima (hoy local, impide seleccionar pasado)
            fechaMinima: (function() { var d = new Date(); return d.getFullYear() + '-' + String(d.getMonth()+1).padStart(2,'0') + '-' + String(d.getDate()).padStart(2,'0'); })(),
            // Validación de disponibilidad
            disponibilidadOk: null, // null = no verificado, true/false
            disponibilidadMsg: '',
            verificando: false,
            errorMsg: '',
        };
    },

    computed: {
        puedeGuardar() {
            if (!this.modo) return false;
            if (this.modo === 'inmediato') return true;
            if (this.modo === 'programado') {
                if (!this.fecha || !this.hora) return false;
                if (this.fecha < this.fechaMinima) return false;
                return this.disponibilidadOk === true;
            }
            return false;
        },
    },

    methods: {
        open(psicosocialId) {
            this.psicosocialId = psicosocialId;
            this.visible = true;
            this.guardando = false;
            this.modo = null;
            this.fecha = '';
            this.hora = '';
            this.disponibilidadOk = null;
            this.disponibilidadMsg = '';
            this.verificando = false;
            this.errorMsg = '';
        },

        cerrar() {
            this.visible = false;
            this.psicosocialId = null;
        },

        seleccionarModo(modo) {
            this.modo = modo;
            this.disponibilidadOk = null;
            this.disponibilidadMsg = '';
            this.errorMsg = '';
        },

        async verificarDisponibilidad() {
            if (!this.fecha || !this.hora) return;
            this.verificando = true;
            this.disponibilidadOk = null;
            this.disponibilidadMsg = '';
            try {
                var params = new URLSearchParams({
                    date: this.fecha,
                    time: this.hora,
                    mode: 'individual'
                });
                var res = await fetch('/api/v1/psychosocial-support/' + this.psicosocialId + '/availability?' + params.toString());
                var data = await res.json();
                this.disponibilidadOk = data.available;
                this.disponibilidadMsg = data.message || (data.available ? 'Horario disponible' : 'Conflicto de horario');
            } catch (e) {
                this.disponibilidadOk = false;
                this.disponibilidadMsg = 'Error al verificar disponibilidad';
            } finally {
                this.verificando = false;
            }
        },

        async guardar() {
            if (!this.puedeGuardar || this.guardando) return;
            this.guardando = true;
            this.errorMsg = '';

            var payload = { psicosocialId: this.psicosocialId };

            if (this.modo === 'inmediato') {
                // Contacto inmediato: fecha/hora actual, tipo "realizar"
                var now = new Date();
                payload.contactDate = now.toISOString().substring(0, 10);
                payload.contactTime = now.toTimeString().substring(0, 5);
                payload.type = 'realizar';
                payload.sessionType = 'PRIMER_CONTACTO';
            } else {
                // Programado: agendar con fecha/hora seleccionada
                payload.contactDate = this.fecha;
                payload.contactTime = this.hora;
                payload.type = 'agendar';
                payload.sessionType = 'PRIMER_CONTACTO';
                payload.scheduledDate = this.fecha;
                payload.scheduledTime = this.hora;
            }

            try {
                var res = await fetch('/api/v1/psychosocial-support/' + this.psicosocialId + '/contacts', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify(payload),
                });
                if (!res.ok) {
                    var err = await res.json();
                    throw new Error(err.error || 'Error al guardar');
                }
                var data = await res.json();

                if (this.modo === 'inmediato') {
                    // Solo redirigir al formulario si es ps/ts (profesional que atiende)
                    var role = this.$root ? (this.$root.currentRole || '') : '';
                    if ((role === 'ps' || role === 'ts') && data && data.id) {
                        window.location.href = '/salvia/psicosocial/registrar/' + this.psicosocialId + '?contact_id=' + data.id;
                        return;
                    }
                    // Si es sv u otro rol, solo confirmar y cerrar
                    this.$emit('saved');
                    this.cerrar();
                } else {
                    this.$emit('saved');
                    this.cerrar();
                }
            } catch (e) {
                this.errorMsg = e.message || 'Error al registrar contacto';
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
        <p class="rcm-subtitle">Selecciona cómo deseas registrar este contacto con la paciente.</p>

        <!-- Opciones -->
        <div :class="['rcm-option-card', modo === 'inmediato' ? 'active' : '']" @click="seleccionarModo('inmediato')">
            <p class="rcm-option-title">⚡ Contacto inmediato</p>
            <p class="rcm-option-desc">Registrar un contacto ahora y comenzar el formulario psicosocial de inmediato.</p>
        </div>
        <div :class="['rcm-option-card', modo === 'programado' ? 'active' : '']" @click="seleccionarModo('programado')">
            <p class="rcm-option-title">📅 Programar contacto</p>
            <p class="rcm-option-desc">Agendar una sesión para una fecha y hora específica (se valida disponibilidad).</p>
        </div>

        <!-- Campos para programado -->
        <template v-if="modo === 'programado'">
            <div class="rcm-row" style="margin-top:16px">
                <div class="rcm-field">
                    <label class="rcm-label">Fecha</label>
                    <input type="date" class="rcm-input" v-model="fecha" :min="fechaMinima" @change="disponibilidadOk = null" />
                </div>
                <div class="rcm-field">
                    <label class="rcm-label">Hora</label>
                    <input type="time" class="rcm-input" v-model="hora" @change="disponibilidadOk = null" />
                </div>
            </div>
            <button type="button" @click="verificarDisponibilidad" :disabled="!fecha || !hora || verificando"
                style="width:100%;padding:10px;border-radius:8px;border:1px solid #5106A7;background:#fff;color:#5106A7;font-size:.82rem;font-weight:600;cursor:pointer;margin-bottom:12px"
                :style="{ opacity: (!fecha || !hora || verificando) ? '.5' : '1' }">
                \${ verificando ? 'Verificando...' : '🔍 Verificar disponibilidad' }
            </button>
            <div v-if="disponibilidadOk === true" class="rcm-success">\${ disponibilidadMsg }</div>
            <div v-if="disponibilidadOk === false" class="rcm-alert">\${ disponibilidadMsg }</div>
        </template>

        <!-- Error -->
        <div v-if="errorMsg" class="rcm-alert">\${ errorMsg }</div>

        <!-- Footer -->
        <div class="rcm-footer">
            <button class="rcm-btn-cancel" @click="cerrar">Cancelar</button>
            <button class="rcm-btn-save" :disabled="!puedeGuardar || guardando" @click="guardar">
                \${ guardando ? 'Guardando...' : (modo === 'inmediato' ? 'Iniciar atención' : 'Agendar sesión') }
            </button>
        </div>
    </div>
</div>
    `,
});
