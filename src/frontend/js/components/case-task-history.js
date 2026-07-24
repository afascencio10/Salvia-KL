/**
 * case-task-history
 * Modal de solo lectura para ver el detalle de una CaseTask ya completada.
 * El padre lo abre vía: this.$refs.taskHistory.open(taskId)
 * No emite eventos ni recibe props — se controla exclusivamente con open(taskId).
 */

/* ─── Estilos ────────────────────────────────────────────────────────────── */
(function injectCaseTaskHistoryStyles() {
    if (document.getElementById('case-task-history-styles')) return;
    const style = document.createElement('style');
    style.id = 'case-task-history-styles';
    style.textContent = `
        .cth-overlay {
            position: fixed;
            inset: 0;
            background: rgba(0,0,0,.45);
            display: flex;
            align-items: center;
            justify-content: center;
            z-index: 1050;
            padding: 16px;
        }

        .cth-modal {
            position: relative;
            background: #fff;
            border-radius: 16px;
            width: 100%;
            max-width: 540px;
            max-height: 90vh;
            display: flex;
            flex-direction: column;
            box-shadow: 0 20px 40px rgba(0,0,0,.18);
            overflow: hidden;
        }

        .cth-header {
            display: flex;
            align-items: center;
            justify-content: space-between;
            padding: 20px 24px 16px;
            border-bottom: 1px solid #e5e7eb;
            flex-shrink: 0;
        }
        .cth-header-left {
            display: flex;
            align-items: center;
            gap: 12px;
        }
        .cth-header-icon {
            width: 40px;
            height: 40px;
            border-radius: 10px;
            background: #f0fdf4;
            color: #16a34a;
            display: flex;
            align-items: center;
            justify-content: center;
            font-size: 18px;
            flex-shrink: 0;
        }
        .cth-header-title {
            font-size: 16px;
            font-weight: 700;
            color: #111827;
            margin: 0;
        }
        .cth-header-sub {
            font-size: 12px;
            color: #9ca3af;
            margin: 2px 0 0;
        }
        .cth-close-btn {
            background: none;
            border: none;
            font-size: 20px;
            color: #9ca3af;
            cursor: pointer;
            padding: 4px;
            line-height: 1;
            transition: color .15s;
        }
        .cth-close-btn:hover { color: #4b5563; }

        .cth-body {
            flex: 1;
            overflow-y: auto;
            padding: 20px 24px;
        }

        .cth-meta-row {
            display: flex;
            flex-direction: column;
            gap: 4px;
            margin-bottom: 16px;
            padding-bottom: 16px;
            border-bottom: 1px dashed #e5e7eb;
        }
        .cth-completado-por {
            font-size: 13px;
            font-weight: 600;
            color: #374151;
        }
        .cth-fecha {
            font-size: 12px;
            color: #9ca3af;
        }

        .cth-descripcion {
            font-size: 13px;
            color: #374151;
            background: #f9fafb;
            border-radius: 8px;
            padding: 10px 12px;
            margin-bottom: 16px;
            line-height: 1.5;
        }

        .cth-section-title {
            font-size: 11px;
            font-weight: 700;
            color: #9ca3af;
            text-transform: uppercase;
            letter-spacing: .05em;
            margin: 16px 0 10px;
        }
        .cth-section-title:first-child { margin-top: 0; }

        .cth-field {
            margin-bottom: 12px;
        }
        .cth-field-label {
            display: block;
            font-size: 11px;
            font-weight: 600;
            color: #9ca3af;
            text-transform: uppercase;
            letter-spacing: .03em;
            margin-bottom: 3px;
        }
        .cth-field-value {
            font-size: 13px;
            color: #111827;
            line-height: 1.4;
        }

        .cth-chips {
            display: flex;
            flex-wrap: wrap;
            gap: 8px;
            margin-bottom: 14px;
        }
        .cth-chip {
            display: inline-block;
            padding: 5px 12px;
            border-radius: 999px;
            background: #eff6ff;
            color: #1d4ed8;
            font-size: 12px;
            font-weight: 600;
        }

        .cth-nota-corregido {
            font-size: 13px;
            color: #374151;
            background: #eff6ff;
            border: 1px solid #bfdbfe;
            border-radius: 8px;
            padding: 12px 14px;
            line-height: 1.5;
        }

        /* Error banner */
        .cth-error-banner {
            padding: 10px 14px;
            background: #fef2f2;
            border: 1px solid #fecaca;
            border-radius: 8px;
            color: #b91c1c;
            font-size: 13px;
        }

        /* Loading */
        .cth-loading {
            display: flex;
            flex-direction: column;
            align-items: center;
            padding: 40px 0;
            gap: 12px;
            color: #9ca3af;
            font-size: 13px;
        }
        @keyframes cth-spin { to { transform: rotate(360deg); } }
        .cth-spinner {
            width: 28px; height: 28px;
            animation: cth-spin .7s linear infinite;
        }

        /* Footer */
        .cth-footer {
            display: flex;
            padding: 16px 24px;
            border-top: 1px solid #e5e7eb;
            flex-shrink: 0;
        }
        .cth-btn-close {
            flex: 1;
            padding: 10px 16px;
            border-radius: 8px;
            font-size: 13px;
            font-weight: 600;
            cursor: pointer;
            transition: background .15s;
            border: 1px solid #e5e7eb;
            background: #fff;
            color: #374151;
            outline: none;
        }
        .cth-btn-close:hover { background: #f9fafb; }
    `;
    document.head.appendChild(style);
})();

/* ─── Componente Vue ──────────────────────────────────────────────────────── */
app.component('case-task-history', {
    delimiters: ['${', '}'],

    data() {
        return {
            visible:  false,
            tarea:    null,
            cargando: false,
            error:    null,
        };
    },

    computed: {
        titulo() {
            if (!this.tarea) return 'Detalle de tarea';
            return {
                gestion_llamada:  'Gestión de Llamada',
                proyectar_oficio: 'Proyectar Oficio',
                comite_caso:      'Decisiones del Comité',
                'Corregir oficio': 'Corregir Oficio',
                gestion_propia:   'Gestión Propia',
            }[this.tarea.type] || 'Detalle de tarea';
        },

        iconoTipo() {
            if (!this.tarea) return 'fa-tasks';
            return {
                gestion_llamada:  'fa-phone-alt',
                proyectar_oficio: 'fa-file-signature',
                comite_caso:      'fa-users',
                'Corregir oficio': 'fa-pencil-alt',
                gestion_propia:   'fa-link',
            }[this.tarea.type] || 'fa-tasks';
        },

        ubicacion() {
            const f = (this.tarea && this.tarea.formData) || {};
            return [f.departamentoNombre, f.ciudadNombre, f.municipioNombre]
                .filter(Boolean)
                .join(' / ');
        },
    },

    methods: {
        /* ── Public API ──────────────────────────────────────────────────── */

        async open(taskId) {
            this.visible  = true;
            this.tarea    = null;
            this.error    = null;
            this.cargando = true;

            try {
                const res = await fetch('/api/v1/case-tasks/' + taskId);
                if (res.status === 401) { window.location.href = '/static/landing.html'; return; }
                if (!res.ok) throw new Error('Error ' + res.status + ' al cargar la tarea.');
                this.tarea = await res.json();
            } catch (e) {
                this.error = 'No se pudo cargar la tarea. ' + e.message;
            } finally {
                this.cargando = false;
            }
        },

        cerrar() {
            this.visible  = false;
            this.tarea    = null;
            this.error    = null;
        },

        /* ── Utilidades ──────────────────────────────────────────────────── */

        formatearFecha(fecha) {
            if (!fecha) return '';
            const d = new Date(fecha);
            if (isNaN(d.getTime())) return fecha;
            return d.toLocaleDateString('es-CO', { day: '2-digit', month: 'long', year: 'numeric' })
                + ' · ' + d.toLocaleTimeString('es-CO', { hour: '2-digit', minute: '2-digit' });
        },
    },

    template: `
<div>
<div v-if="visible" class="cth-overlay" @click.self="cerrar">
<div class="cth-modal">

    <!-- Header -->
    <div class="cth-header">
        <div class="cth-header-left">
            <div class="cth-header-icon">
                <i :class="'fas ' + iconoTipo"></i>
            </div>
            <div>
                <p class="cth-header-title">\${ titulo }</p>
                <p v-if="tarea" class="cth-header-sub">ID \${ tarea.id.substring(0,8) }…</p>
            </div>
        </div>
        <button class="cth-close-btn" @click="cerrar">✕</button>
    </div>

    <!-- Body -->
    <div class="cth-body">

        <!-- Cargando -->
        <div v-if="cargando" class="cth-loading">
            <svg class="cth-spinner" viewBox="0 0 24 24" fill="none">
                <circle cx="12" cy="12" r="10" stroke="#e5e7eb" stroke-width="3"/>
                <path d="M12 2a10 10 0 0 1 10 10" stroke="#16a34a" stroke-width="3" stroke-linecap="round"/>
            </svg>
            Cargando tarea...
        </div>

        <!-- Error -->
        <div v-else-if="error" class="cth-error-banner">
            <i class="fas fa-exclamation-circle"></i> \${ error }
        </div>

        <!-- Detalle -->
        <template v-else-if="tarea">

            <div class="cth-meta-row">
                <p class="cth-completado-por">Completado por: \${ tarea.assignedUserName || 'Sin asignar' }</p>
                <p class="cth-fecha">Completado el \${ formatearFecha(tarea.completedAt) }</p>
            </div>

            <p v-if="tarea.description" class="cth-descripcion">\${ tarea.description }</p>

            <!-- Detalle formData: gestion_llamada -->
            <template v-if="tarea.type === 'gestion_llamada'">
                <p class="cth-section-title">Detalle de la gestión</p>
                <div class="cth-field">
                    <label class="cth-field-label">Ubicación</label>
                    <div class="cth-field-value">\${ ubicacion || '—' }</div>
                </div>
                <div class="cth-field">
                    <label class="cth-field-label">Entidad</label>
                    <div class="cth-field-value">\${ tarea.formData.entidadNombre || '—' }</div>
                </div>
                <div class="cth-field">
                    <label class="cth-field-label">Funcionario</label>
                    <div class="cth-field-value">\${ tarea.formData.funcionario || '—' }</div>
                </div>
                <div class="cth-field" v-if="tarea.formData.descripcion">
                    <label class="cth-field-label">Notas</label>
                    <div class="cth-field-value">\${ tarea.formData.descripcion }</div>
                </div>
                <div class="cth-field">
                    <label class="cth-field-label">¿Generó oficio?</label>
                    <div class="cth-field-value">\${ tarea.formData.generaOficio ? 'Sí' : 'No' }</div>
                </div>
                <template v-if="tarea.formData.generaOficio">
                    <div class="cth-field">
                        <label class="cth-field-label">Asunto</label>
                        <div class="cth-field-value">\${ tarea.formData.asunto || '—' }</div>
                    </div>
                    <div class="cth-field">
                        <label class="cth-field-label">Ruta Kofax</label>
                        <div class="cth-field-value">\${ tarea.formData.rutaKofax || '—' }</div>
                    </div>
                </template>
            </template>

            <!-- Detalle formData: proyectar_oficio -->
            <template v-else-if="tarea.type === 'proyectar_oficio'">
                <p class="cth-section-title">Detalle del oficio proyectado</p>
                <div class="cth-field">
                    <label class="cth-field-label">Ubicación</label>
                    <div class="cth-field-value">\${ ubicacion || '—' }</div>
                </div>
                <div class="cth-field">
                    <label class="cth-field-label">Entidad</label>
                    <div class="cth-field-value">\${ tarea.formData.entidadNombre || '—' }</div>
                </div>
                <div class="cth-field">
                    <label class="cth-field-label">Funcionario</label>
                    <div class="cth-field-value">\${ tarea.formData.funcionario || '—' }</div>
                </div>
                <div class="cth-field">
                    <label class="cth-field-label">Asunto</label>
                    <div class="cth-field-value">\${ tarea.formData.asunto || '—' }</div>
                </div>
                <div class="cth-field">
                    <label class="cth-field-label">Ruta Kofax</label>
                    <div class="cth-field-value">\${ tarea.formData.rutaKofax || '—' }</div>
                </div>
            </template>

            <!-- Detalle formData: comite_caso -->
            <template v-else-if="tarea.type === 'comite_caso'">
                <p class="cth-section-title">Decisiones tomadas</p>
                <div class="cth-chips">
                    <span class="cth-chip" v-for="d in (tarea.formData.decisionesTexto || [])" :key="d">\${ d }</span>
                </div>
                <div class="cth-field" v-if="tarea.formData.observacionesOficio">
                    <label class="cth-field-label">Obs. oficio</label>
                    <div class="cth-field-value">\${ tarea.formData.observacionesOficio }</div>
                </div>
                <div class="cth-field" v-if="tarea.formData.observacionesRecomendaciones">
                    <label class="cth-field-label">Obs. recomendaciones</label>
                    <div class="cth-field-value">\${ tarea.formData.observacionesRecomendaciones }</div>
                </div>
                <template v-if="tarea.formData.nivelMecanismo">
                    <div class="cth-field">
                        <label class="cth-field-label">Nivel del mecanismo</label>
                        <div class="cth-field-value">\${ tarea.formData.nivelMecanismoTexto || tarea.formData.nivelMecanismo }</div>
                    </div>
                    <div class="cth-field" v-if="tarea.formData.observacionesMecanismo">
                        <label class="cth-field-label">Obs. mecanismo</label>
                        <div class="cth-field-value">\${ tarea.formData.observacionesMecanismo }</div>
                    </div>
                </template>
            </template>

            <!-- Detalle: Corregir oficio -->
            <template v-else-if="tarea.type === 'Corregir oficio'">
                <p class="cth-nota-corregido">El oficio fue corregido por el agente que lo proyectó</p>
            </template>

            <!-- Detalle formData: gestion_propia -->
            <template v-else-if="tarea.type === 'gestion_propia'">
                <p class="cth-section-title">Detalle de la gestión</p>
                <div class="cth-field">
                    <label class="cth-field-label">Tipo de gestión</label>
                    <div class="cth-field-value">\${ (tarea.formData && tarea.formData.subtipo) || '—' }</div>
                </div>
            </template>

        </template>

    </div>

    <!-- Footer -->
    <div class="cth-footer">
        <button class="cth-btn-close" @click="cerrar">Cerrar</button>
    </div>

</div>
</div>
</div>
    `,
});
