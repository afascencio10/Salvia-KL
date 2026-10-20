/**
 * gestion-reporte-modal
 * Modal de gestión de un reporte individual — validación previa, esquema
 * 3x3 de recontacto, y cierre. Abierto desde Listado Reportes / Mis
 * Reportes al presionar "Ver Reporte" (action-clicked buttonId==='ver').
 *
 * Diseño tomado literal de Mockup_Reportes_Registro.html (PASO 6a-6e).
 * Ver DocsMD/Componentes/Gestion Reporte Modal/.
 *
 * M-03 ("Reporte de prueba") — removido (2026-10-12): los reportes de
 * prueba se generan y limpian en el ambiente de pruebas, no necesitan
 * manejo en producción.
 * M-10 ("Apertura de caso") — diferido, no implementado en esta fase
 * (coincide con la nota de FigJam). El paso 'postContacto' solo ofrece
 * "no competencia" y "decisión de la persona".
 *
 * Método público:
 *   open(report)  — abre el modal (M-01)
 *
 * Emite:
 *   updated  — al cerrarse el modal, siempre (M-11)
 */
(function() {
    var vueApp = (typeof app !== 'undefined') ? app : (typeof home !== 'undefined' ? home : null);
    if (!vueApp) {
        console.error('[gestion-reporte-modal] No se encontró app/home.');
        return;
    }

    // El esquema 3x3 cuenta "días" en America/Bogota (backend). El log del modal debe agrupar y mostrar las
    // horas con esa MISMA zona, no con la del navegador: si no, un usuario fuera de Colombia vería badges X/3
    // que no coinciden con "Días completos". Mismo criterio que formatFechaHora en Detalle del Caso.
    var BOGOTA_TZ = 'America/Bogota';

    vueApp.component('gestion-reporte-modal', {
        delimiters: ['${', '}'],

        emits: ['updated'],

        data: function() {
            return {
                visible: false,
                icode: '',
                loading: false,
                loadError: null,
                step: null, // 'validacionPrevia' | 'esquema3x3' | 'postContacto' | 'cerrado'

                rowPreview: null, // fila de la tabla (ReportListItem): permite pintar el header antes de que llegue el detalle
                report: null, // detalle crudo del backend (GET /api/v1/victim-contacts/:icode)
                intentos: [],
                intentosHoy: 0,
                diasCompletos: 0,
                canManage: false,

                closeReasonCode: '',
                closeReasonLabel: '',
                closeCategory: '',
                closeObservationShown: '',

                pendingClose: '', // '' | 'competencia' | 'decision'
                motivoTexto: '',

                saving: false,
                actionError: null,
            };
        },

        computed: {
            nombreMostrado: function() {
                if (this.report) {
                    return this.report.nombreMostrado;
                }
                var row = this.rowPreview;
                return row ? (((row.names || '') + ' ' + (row.lastNames || '')).trim()) : '';
            },
            reportTypeLabel: function() {
                if (this.report) {
                    return this.report.reportType || '';
                }
                return this.rowPreview ? (this.rowPreview.reportType || '') : '';
            },
            // Historial de intentos agrupado por día (más reciente primero) —
            // agregado por Felipe (2026-10-12): ver día a día cómo va el 3x3.
            intentosPorDia: function() {
                var groups = {};
                var order = [];
                this.intentos.forEach(function(it) {
                    var d = new Date(it.date);
                    var dayKey = isNaN(d.getTime()) ? 'desconocido' : d.toLocaleDateString('es-CO', { timeZone: BOGOTA_TZ });
                    if (!groups[dayKey]) {
                        groups[dayKey] = { fecha: dayKey, fallidos: 0, items: [] };
                        order.push(dayKey);
                    }
                    groups[dayKey].items.push(it);
                    if (it.result === 'fallido') {
                        groups[dayKey].fallidos++;
                    }
                });
                return order.reverse().map(function(k) { return groups[k]; });
            },
        },

        methods: {

            // ------------------------------------------------------------
            // M-01 — Abrir modal
            // ------------------------------------------------------------

            open: function(reportRow) {
                if (!reportRow) {
                    return;
                }
                var icode = reportRow.i_code || reportRow.icode || String(reportRow.id || '');
                if (!icode) {
                    return;
                }

                this.visible = true;
                this.icode = icode;
                this.loading = true;
                this.loadError = null;
                this.step = null;
                this.rowPreview = reportRow;
                this.report = null;
                this.intentos = [];
                this.intentosHoy = 0;
                this.diasCompletos = 0;
                this.canManage = false;
                this.pendingClose = '';
                this.motivoTexto = '';
                this.actionError = null;
                this.saving = false;

                this.fetchDetail();
            },

            fetchDetail: function() {
                var self = this;
                fetch('/api/v1/victim-contacts/' + encodeURIComponent(self.icode))
                    .then(function(res) {
                        return res.json().then(function(data) {
                            return { ok: res.ok, data: data };
                        });
                    })
                    .then(function(result) {
                        if (!result.ok) {
                            self.loadError = (result.data && result.data.error) || 'Error al cargar el reporte';
                            self.loading = false;
                            return;
                        }
                        self.applyDetail(result.data);
                        self.loading = false;
                    })
                    .catch(function() {
                        self.loadError = 'Error al cargar el reporte';
                        self.loading = false;
                    });
            },

            applyDetail: function(data) {
                this.report = data;
                this.intentos = data.intentos || [];
                this.intentosHoy = data.intentosHoy || 0;
                this.diasCompletos = data.diasCompletos || 0;
                this.canManage = !!data.canManage;
                this.closeReasonCode = data.closeReasonCode || '';
                this.closeObservationShown = data.closeObservation || '';
                // Mismo texto que se arma al cerrar (closeReasonLabelFor): motivo + " — " + observación (co/dp).
                // Sin esto, al reabrir un reporte cerrado el motivo salía sin la observación.
                this.closeReasonLabel = (data.closeReasonLabel || '') +
                    (this.closeObservationShown ? ' — ' + this.closeObservationShown : '');
                this.closeCategory = data.closeCategory || '';

                var huboExitoso = this.intentos.some(function(i) { return i.result === 'exitoso'; });

                if (this.closeReasonCode) {
                    this.step = 'cerrado';
                } else if (huboExitoso) {
                    // Corrección 2026-10-12 (ver flow-M01): sin esto, reabrir un
                    // reporte con contacto exitoso sin cierre terminado volvía
                    // a esquema3x3 en vez de postContacto.
                    this.step = 'postContacto';
                } else if (this.intentos.length > 0) {
                    this.step = 'esquema3x3';
                } else {
                    this.step = 'validacionPrevia';
                }
            },

            // ------------------------------------------------------------
            // M-02 — "No información mínima" (M-03 "Reporte de prueba" removido)
            // ------------------------------------------------------------

            cerrarPrevio: function(code) {
                this.runClose(code, '');
            },

            // ------------------------------------------------------------
            // M-04 — "Procede recontacto" (transición local, sin backend)
            // ------------------------------------------------------------

            irAEsquema3x3: function() {
                if (!this.canManage || this.step !== 'validacionPrevia') {
                    return;
                }
                this.step = 'esquema3x3';
            },

            // ------------------------------------------------------------
            // M-05 / M-06 — Registrar intento (fallido / exitoso)
            // ------------------------------------------------------------

            registrarIntento: function(resultado) {
                var self = this;
                if (!self.canManage || self.step !== 'esquema3x3' || self.saving) {
                    return;
                }
                self.saving = true;
                self.actionError = null;

                fetch('/api/v1/victim-contacts/' + encodeURIComponent(self.icode) + '/intento', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ result: resultado }),
                })
                    .then(function(res) {
                        return res.json().then(function(data) {
                            return { ok: res.ok, data: data };
                        });
                    })
                    .then(function(result) {
                        if (!result.ok) {
                            self.actionError = (result.data && result.data.error) || 'Error al registrar el intento';
                            self.saving = false;
                            return;
                        }

                        self.intentos = self.intentos.concat([{
                            date: new Date().toISOString(),
                            result: resultado,
                            agentName: result.data.agentName || '',
                        }]);
                        self.intentosHoy = result.data.intentosHoy || 0;
                        self.diasCompletos = result.data.diasCompletos || 0;
                        self.saving = false;

                        if (resultado === 'exitoso') {
                            // Mismo retardo visual que el mockup (setTimeout 400ms) —
                            // da tiempo a ver el intento en el log antes de cambiar de vista.
                            setTimeout(function() {
                                self.step = 'postContacto';
                                self.pendingClose = '';
                            }, 400);
                        }
                    })
                    .catch(function() {
                        self.actionError = 'Error al registrar el intento';
                        self.saving = false;
                    });
            },

            // ------------------------------------------------------------
            // M-07 — "Cerrar por no contacto" (solo si diasCompletos >= 3)
            // ------------------------------------------------------------

            cerrarPorNoContacto: function() {
                if (this.diasCompletos < 3) {
                    return;
                }
                this.runClose('nc', '');
            },

            // ------------------------------------------------------------
            // M-08 / M-09 — Cierre post-contacto (no competencia / decisión)
            // ------------------------------------------------------------

            seleccionarPendingClose: function(value) {
                this.pendingClose = value;
                this.motivoTexto = '';
            },

            cancelarPendingClose: function() {
                this.pendingClose = '';
                this.motivoTexto = '';
            },

            confirmarCierrePostContacto: function() {
                if (!this.motivoTexto || !this.motivoTexto.trim()) {
                    return;
                }
                var code = this.pendingClose === 'competencia' ? 'co' : 'dp';
                this.runClose(code, this.motivoTexto.trim());
            },

            // ------------------------------------------------------------
            // Cierre genérico (M-02, M-07, M-08, M-09)
            // ------------------------------------------------------------

            runClose: function(code, observation) {
                var self = this;
                if (!self.canManage || self.saving) {
                    return;
                }
                self.saving = true;
                self.actionError = null;

                var body = { close_reason: code };
                if (observation) {
                    body.close_observation = observation;
                }

                fetch('/api/v1/victim-contacts/' + encodeURIComponent(self.icode) + '/cerrar', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify(body),
                })
                    .then(function(res) {
                        return res.json().then(function(data) {
                            return { ok: res.ok, data: data };
                        });
                    })
                    .then(function(result) {
                        self.saving = false;
                        if (!result.ok) {
                            self.actionError = (result.data && result.data.error) || 'Error al cerrar el reporte';
                            return;
                        }
                        self.closeReasonCode = code;
                        self.closeCategory = result.data.category || (code === 'ap' ? 'Válido' : 'Inválido');
                        self.closeReasonLabel = self.closeReasonLabelFor(code, observation);
                        self.pendingClose = '';
                        self.motivoTexto = '';
                        self.step = 'cerrado';
                    })
                    .catch(function() {
                        self.saving = false;
                        self.actionError = 'Error al cerrar el reporte';
                    });
            },

            closeReasonLabelFor: function(code, observation) {
                var map = {
                    ni: 'Inválido por No información Mínima',
                    nc: 'Inválido por Cierre por no contacto',
                    co: 'Inválido por Cierre por no competencia',
                    dp: 'Inválido por Cierre por decisión de la persona',
                    ap: 'Válido - Apertura de Caso Salvia',
                };
                var label = map[code] || code;
                if (observation) {
                    label += ' — ' + observation;
                }
                return label;
            },

            // ------------------------------------------------------------
            // M-11 — Cerrar modal (siempre emite 'updated')
            // ------------------------------------------------------------

            close: function() {
                if (this.saving) {
                    return;
                }
                this.visible = false;
                this.$emit('updated');
            },

            onBackdropClick: function() {
                this.close();
            },

            // ------------------------------------------------------------
            // Helpers de presentación
            // ------------------------------------------------------------

            formatDateTime: function(dateStr) {
                if (!dateStr) return '—';
                var d = new Date(dateStr);
                if (isNaN(d.getTime())) return '—';
                return d.toLocaleDateString('es-CO', { timeZone: BOGOTA_TZ }) + ' ' +
                    d.toLocaleTimeString('es-CO', { hour: '2-digit', minute: '2-digit', timeZone: BOGOTA_TZ });
            },

            resultLabel: function(result) {
                return result === 'exitoso' ? 'Exitoso' : 'Fallido';
            },
        },

        template: (typeof window !== 'undefined' && window.__gestionReporteModalTpl)
            ? window.__gestionReporteModalTpl
            : '#tpl-gestion-reporte-modal',
    });
})();
