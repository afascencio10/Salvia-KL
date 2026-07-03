/**
 * reasignar-remisiones-modal
 * Modal para reasignar remisiones psicosociales seleccionadas.
 */
(function() {
    var vueApp = (typeof app !== 'undefined') ? app : (typeof home !== 'undefined' ? home : null);
    if (!vueApp) {
        console.error('[reasignar-remisiones-modal] No se encontró app/home.');
        return;
    }

    var ENDPOINT_PROFESSIONALS = '/api/v1/psychosocial-support/profesionales-reasignacion';
    var ENDPOINT_DUPLAS = '/api/v1/duplas/reasignacion';
    var ENDPOINT_REASSIGN = '/api/v1/psychosocial-support/reasignar-bulk';
    var STATUS_CERRADO = 'cerrado';

    vueApp.component('reasignar-remisiones-modal', {
        delimiters: ['${', '}'],

        emits: ['closed', 'reassigned'],

        data: function() {
            return {
                visible: false,
                remisiones: [],
                asignarEnDupla: false,
                selectedProfessionalId: '',
                selectedDuplaId: '',
                professionalGroups: [],
                duplaOptions: [],
                loadingOptions: false,
                optionsError: null,
                saving: false,
                saveError: null,
                confirmVisible: false,
                _loadSeq: 0,
            };
        },

        computed: {
            canConfirm: function() {
                if (this.loadingOptions || this.saving) {
                    return false;
                }
                if (this.optionsError) {
                    return false;
                }
                return this.asignarEnDupla
                    ? !!this.selectedDuplaId
                    : !!this.selectedProfessionalId;
            },

            selectionSubtitle: function() {
                var n = this.remisiones.length;
                if (n === 1) {
                    return '1 remisión seleccionada.';
                }
                return n + ' remisiones seleccionadas.';
            },

            selectedAssigneeLabel: function() {
                if (this.asignarEnDupla) {
                    var dupla = this.duplaOptions.find(function(d) {
                        return d.id === this.selectedDuplaId;
                    }.bind(this));
                    return dupla ? dupla.label : '';
                }
                var self = this;
                var found = '';
                this.professionalGroups.forEach(function(group) {
                    (group.professionals || []).forEach(function(prof) {
                        if (prof.icode === self.selectedProfessionalId) {
                            found = prof.fullName;
                        }
                    });
                });
                return found;
            },
        },

        methods: {

            // RRM-01 — Abrir modal
            open: function(remisiones) {
                if (!remisiones || !Array.isArray(remisiones) || remisiones.length === 0) {
                    return;
                }

                var valid = remisiones.filter(function(r) {
                    return r && r.status !== STATUS_CERRADO;
                });
                if (!valid.length) {
                    return;
                }

                this.visible = true;
                this.remisiones = valid.slice();
                this.asignarEnDupla = false;
                this.selectedProfessionalId = '';
                this.selectedDuplaId = '';
                this.professionalGroups = [];
                this.duplaOptions = [];
                this.optionsError = null;
                this.saveError = null;
                this.saving = false;
                this.confirmVisible = false;

                this.loadOptions();
                this.focusFirstControl();
            },

            focusFirstControl: function() {
                var self = this;
                this.$nextTick(function() {
                    if (!self.visible || !self.$el) {
                        return;
                    }
                    var toggle = self.$el.querySelector('.rrm-toggle');
                    if (toggle) {
                        toggle.focus();
                    }
                });
            },

            // RRM-02 — Alternar switch dupla
            onToggleDupla: function() {
                if (this.saving || !this.visible) {
                    return;
                }
                this.asignarEnDupla = !this.asignarEnDupla;
                this.selectedProfessionalId = '';
                this.selectedDuplaId = '';
                this.optionsError = null;
                this.loadOptions();
            },

            // RRM-03 — Cargar opciones del select
            loadOptions: function() {
                var self = this;
                if (!self.visible) {
                    return;
                }

                self._loadSeq = (self._loadSeq || 0) + 1;
                var seq = self._loadSeq;

                self.loadingOptions = true;
                self.optionsError = null;

                var url = self.asignarEnDupla ? ENDPOINT_DUPLAS : ENDPOINT_PROFESSIONALS;

                fetch(url)
                    .then(function(res) {
                        return res.json().then(function(data) {
                            return { ok: res.ok, data: data };
                        });
                    })
                    .then(function(result) {
                        if (seq !== self._loadSeq) {
                            return;
                        }

                        if (!result.ok) {
                            self.optionsError = (result.data && result.data.error)
                                || 'Error al cargar las opciones';
                            self.professionalGroups = [];
                            self.duplaOptions = [];
                            return;
                        }

                        if (self.asignarEnDupla) {
                            var duplas = (result.data && result.data.duplas) || [];
                            self.duplaOptions = Array.isArray(duplas) ? duplas : [];
                            self.professionalGroups = [];
                            if (self.duplaOptions.length === 0) {
                                self.optionsError = 'No hay duplas disponibles';
                            }
                        } else {
                            var groups = (result.data && result.data.groups) || [];
                            self.professionalGroups = Array.isArray(groups) ? groups : [];
                            self.duplaOptions = [];
                            var total = self.professionalGroups.reduce(function(acc, g) {
                                return acc + ((g.professionals && g.professionals.length) || 0);
                            }, 0);
                            if (total === 0) {
                                self.optionsError = 'No hay profesionales disponibles';
                            }
                        }
                    })
                    .catch(function() {
                        if (seq !== self._loadSeq) {
                            return;
                        }
                        self.optionsError = 'Error al cargar las opciones';
                        self.professionalGroups = [];
                        self.duplaOptions = [];
                    })
                    .finally(function() {
                        if (seq !== self._loadSeq) {
                            return;
                        }
                        self.loadingOptions = false;
                    });
            },

            // RRM-04 — Cerrar modal
            onBackdropClick: function() {
                if (this.saving || this.confirmVisible) {
                    return;
                }
                this.close();
            },

            close: function() {
                if (this.saving || this.confirmVisible) {
                    return;
                }

                this._loadSeq = (this._loadSeq || 0) + 1;
                this.visible = false;
                this.remisiones = [];
                this.asignarEnDupla = false;
                this.selectedProfessionalId = '';
                this.selectedDuplaId = '';
                this.professionalGroups = [];
                this.duplaOptions = [];
                this.loadingOptions = false;
                this.optionsError = null;
                this.saveError = null;
                this.saving = false;
                this.confirmVisible = false;

                this.$emit('closed', {});
            },

            cancelConfirm: function() {
                this.confirmVisible = false;
            },

            // RRM-05 — Confirmar reasignación
            confirmReassign: function() {
                if (!this.canConfirm) {
                    return;
                }
                this.saveError = null;
                this.confirmVisible = true;
            },

            acceptConfirm: function() {
                this.confirmVisible = false;
                this.runReassignSave();
            },

            runReassignSave: function() {
                var self = this;
                self.saveError = null;
                self.saving = true;

                var assignMode = self.asignarEnDupla ? 'dupla' : 'professional';
                var remisionIds = self.remisiones.map(function(r) { return r.id; }).filter(Boolean);

                fetch(ENDPOINT_REASSIGN, {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({
                        remision_ids: remisionIds,
                        assign_mode: assignMode,
                        professional_id: assignMode === 'professional' ? self.selectedProfessionalId : null,
                        dupla_id: assignMode === 'dupla' ? self.selectedDuplaId : null,
                    }),
                })
                    .then(function(res) {
                        return res.json().then(function(data) {
                            return { ok: res.ok, data: data };
                        });
                    })
                    .then(function(result) {
                        if (!result.ok) {
                            self.saveError = (result.data && result.data.error)
                                || 'Error al reasignar las remisiones';
                            self.saving = false;
                            return;
                        }

                        var payload = {
                            remisiones: self.remisiones.slice(),
                            assignMode: assignMode,
                            updated: result.data.remisiones_updated || 0,
                            teamContactsUpdated: result.data.team_contacts_updated || 0,
                        };

                        self.saving = false;
                        self.$emit('reassigned', payload);
                        self.forceClose();
                    })
                    .catch(function() {
                        self.saveError = 'Error al reasignar las remisiones';
                        self.saving = false;
                    });
            },

            forceClose: function() {
                this._loadSeq = (this._loadSeq || 0) + 1;
                this.visible = false;
                this.remisiones = [];
                this.asignarEnDupla = false;
                this.selectedProfessionalId = '';
                this.selectedDuplaId = '';
                this.professionalGroups = [];
                this.duplaOptions = [];
                this.loadingOptions = false;
                this.optionsError = null;
                this.saveError = null;
                this.saving = false;
                this.confirmVisible = false;
            },
        },

        template: (typeof window !== 'undefined' && window.__reasignarRemisionesModalTpl)
            ? window.__reasignarRemisionesModalTpl
            : '#tpl-reasignar-remisiones-modal',
    });
})();
