/**
 * admin-duplas.js
 * Lógica Vue para la pantalla "Administrar Duplas" (administrar_duplas.html).
 *
 * Requiere que el HTML haya expuesto previamente la configuración del servidor
 * en el objeto global window.__AdminDuplasConfig con las propiedades:
 *   - windowTitle {string|any}
 *   - currentUser {string}
 *   - currentRole {string}
 *   - menu        {object}
 */

(function () {
    'use strict';

    var cfg = window.__AdminDuplasConfig || {};

    var home = Vue.createApp({
        delimiters: ['${', '}'],

        data: function () {
            return {
                windowTitle: cfg.windowTitle || 'Administrar Duplas',
                currentUser: cfg.currentUser || '',
                currentRole: cfg.currentRole || '',
                menu: cfg.menu || {},

                isLoading: false,
                loadError: null,

                psychologists: [],
                socialWorkers: [],
                duplas: [],

                modal: {
                    visible: false,
                    kind: null,   // 'form' | 'delete' | 'deleteBlocked'
                    mode: null,   // 'create' | 'edit'
                },
                editingDuplaId: null,
                form: {
                    name: '',
                    psychologistId: '',
                    socialWorkerId: '',
                },
                warningPs: null,
                warningTs: null,
                saveError: null,
                isSaving: false,
                nameFieldError: false,
                toast: { visible: false, message: '', type: 'success' },
                _toastTimer: null,

                deleteTarget: null,
                deleteError: null,
                isDeleting: false,
                blockedMessage: null,

                availablePsychologists: [],
                availableSocialWorkers: [],
            };
        },

        computed: {
            canSave: function () {
                return !this.isSaving
                    && this.form.name.trim().length > 0
                    && this.form.psychologistId !== ''
                    && this.form.socialWorkerId !== '';
            },
        },

        mounted: function () {
            document.title = typeof this.windowTitle === 'string'
                ? this.windowTitle
                : String(this.windowTitle || 'Administrar Duplas');
            document.getElementById('app').style.display = 'block';
            this.reloadScreen();
        },

        methods: {
            initialOf: function (name) {
                if (!name || typeof name !== 'string') return '?';
                var t = name.trim();
                return t ? t.charAt(0).toUpperCase() : '?';
            },

            // E01 — Cuando carga la pantalla
            reloadScreen: function () {
                if (this.currentRole !== 'sv') {
                    this.loadError = 'Rol no autorizado para administrar duplas.';
                    this.isLoading = false;
                    return;
                }

                this.isLoading = true;
                this.loadError = null;

                var self = this;
                Promise.all([
                    fetch('/api/v1/duplas/profesionales').then(function (res) {
                        return res.json().then(function (data) {
                            return { ok: res.ok, status: res.status, data: data };
                        });
                    }),
                    fetch('/api/v1/duplas/activas').then(function (res) {
                        return res.json().then(function (data) {
                            return { ok: res.ok, status: res.status, data: data };
                        });
                    }),
                ]).then(function (results) {
                    var proResult = results[0];
                    var duplaResult = results[1];

                    if (proResult.status === 401 || duplaResult.status === 401) {
                        window.location.href = '/static/landing.html';
                        return;
                    }

                    if (!proResult.ok || !duplaResult.ok) {
                        self.loadError = 'No se pudo cargar la información de duplas. Intenta de nuevo.';
                        self.isLoading = false;
                        return;
                    }

                    var psychologists = (proResult.data.psychologists || []).map(function (p) {
                        return {
                            icode: p.icode,
                            fullName: p.fullName,
                            enDupla: false,
                            duplaId: null,
                            duplaName: null,
                        };
                    });
                    var socialWorkers = (proResult.data.socialWorkers || []).map(function (p) {
                        return {
                            icode: p.icode,
                            fullName: p.fullName,
                            enDupla: false,
                            duplaId: null,
                            duplaName: null,
                        };
                    });
                    var duplas = (duplaResult.data.items || []).map(function (d) {
                        return {
                            id: d.id,
                            name: d.name,
                            psychologistId: d.psychologistId,
                            psychologistName: d.psychologistName,
                            socialWorkerId: d.socialWorkerId,
                            socialWorkerName: d.socialWorkerName,
                        };
                    });

                    var assignedPsychByUserId = {};
                    var assignedTsByUserId = {};
                    duplas.forEach(function (d) {
                        if (d.psychologistId) {
                            assignedPsychByUserId[d.psychologistId] = { duplaId: d.id, duplaName: d.name };
                        }
                        if (d.socialWorkerId) {
                            if (!assignedTsByUserId[d.socialWorkerId]) {
                                assignedTsByUserId[d.socialWorkerId] = [];
                            }
                            assignedTsByUserId[d.socialWorkerId].push({
                                duplaId: d.id,
                                duplaName: d.name,
                            });
                        }
                    });

                    psychologists.forEach(function (p) {
                        var info = assignedPsychByUserId[p.icode];
                        if (info) {
                            p.enDupla = true;
                            p.duplaId = info.duplaId;
                            p.duplaName = info.duplaName;
                        } else {
                            p.enDupla = false;
                            p.duplaId = null;
                            p.duplaName = null;
                        }
                    });

                    socialWorkers.forEach(function (p) {
                        var list = assignedTsByUserId[p.icode];
                        if (list && list.length) {
                            p.enDupla = true;
                            p.duplaId = list[0].duplaId;
                            p.duplaName = list.map(function (x) { return x.duplaName; }).join(', ');
                        } else {
                            p.enDupla = false;
                            p.duplaId = null;
                            p.duplaName = null;
                        }
                    });

                    self.psychologists = psychologists;
                    self.socialWorkers = socialWorkers;
                    self.duplas = duplas;
                    self.isLoading = false;
                }).catch(function () {
                    self.loadError = 'No se pudo cargar la información de duplas. Intenta de nuevo.';
                    self.isLoading = false;
                });
            },

            // E02
            openDuplaModal: function (mode, dupla) {
                if (this.isLoading || this.isSaving) return;

                this.saveError = null;
                this.nameFieldError = false;
                this.warningPs = null;
                this.warningTs = null;

                if (mode === 'edit' && dupla) {
                    this.modal.mode = 'edit';
                    this.editingDuplaId = dupla.id;
                    this.form = {
                        name: dupla.name || '',
                        psychologistId: dupla.psychologistId || '',
                        socialWorkerId: dupla.socialWorkerId || '',
                    };
                } else {
                    this.modal.mode = 'create';
                    this.editingDuplaId = null;
                    this.form = { name: '', psychologistId: '', socialWorkerId: '' };
                }

                this.computeAvailableProfessionals();
                this.modal.kind = 'form';
                this.modal.visible = true;
            },

            computeAvailableProfessionals: function () {
                // Solo la psicóloga es exclusiva de una dupla.
                // La trabajadora social puede pertenecer a varias.
                var occupiedPsych = {};
                var editingId = this.editingDuplaId;
                (this.duplas || []).forEach(function (d) {
                    if (editingId && d.id === editingId) return;
                    if (d.psychologistId) occupiedPsych[d.psychologistId] = true;
                });

                this.availablePsychologists = (this.psychologists || []).filter(function (p) {
                    return !occupiedPsych[p.icode];
                });
                this.availableSocialWorkers = (this.socialWorkers || []).slice();

                this.warningPs = this.availablePsychologists.length === 0
                    ? 'Todas las psicólogas están asignadas a una dupla.'
                    : null;
                this.warningTs = this.availableSocialWorkers.length === 0
                    ? 'No hay trabajadoras sociales activas.'
                    : null;
            },

            // E03
            closeFormModal: function () {
                if (this.isSaving) return;
                this.modal.visible = false;
                this.modal.kind = null;
                this.modal.mode = null;
                this.editingDuplaId = null;
                this.form = { name: '', psychologistId: '', socialWorkerId: '' };
                this.saveError = null;
                this.nameFieldError = false;
                this.warningPs = null;
                this.warningTs = null;
                this.availablePsychologists = [];
                this.availableSocialWorkers = [];
            },

            showToast: function (message, type) {
                var self = this;
                if (this._toastTimer) {
                    clearTimeout(this._toastTimer);
                }
                this.toast = {
                    visible: true,
                    message: message || '',
                    type: type || 'success',
                };
                this._toastTimer = setTimeout(function () {
                    self.toast.visible = false;
                }, 3200);
            },

            // E04 — Cuando guarda modal
            saveDuplaModal: function () {
                if (this.isSaving) return;
                if (!this.modal.visible || this.modal.kind !== 'form') return;

                var name = (this.form.name || '').trim();
                var psychologistId = this.form.psychologistId || '';
                var socialWorkerId = this.form.socialWorkerId || '';

                if (!name || name.length > 36 || !psychologistId || !socialWorkerId || psychologistId === socialWorkerId) {
                    this.saveError = 'Completa nombre, psicóloga y trabajadora social.';
                    return;
                }

                this.isSaving = true;
                this.saveError = null;
                this.nameFieldError = false;

                var isEdit = this.modal.mode === 'edit';
                var url = isEdit
                    ? '/api/v1/duplas/' + encodeURIComponent(this.editingDuplaId)
                    : '/api/v1/duplas';
                var method = isEdit ? 'PUT' : 'POST';
                var self = this;

                fetch(url, {
                    method: method,
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({
                        name: name,
                        psychologistId: psychologistId,
                        socialWorkerId: socialWorkerId,
                    }),
                }).then(function (res) {
                    return res.json().then(function (data) {
                        return { ok: res.ok, status: res.status, data: data };
                    }).catch(function () {
                        return { ok: res.ok, status: res.status, data: {} };
                    });
                }).then(function (result) {
                    if (result.status === 401) {
                        window.location.href = '/static/landing.html';
                        return;
                    }

                    if (result.status === 200 || result.status === 201) {
                        self.isSaving = false;
                        self.closeFormModal();
                        self.showToast(
                            isEdit ? 'La dupla se editó exitosamente' : 'La dupla se creó exitosamente',
                            'success'
                        );
                        self.reloadScreen();
                        return;
                    }

                    self.isSaving = false;
                    var apiError = (result.data && result.data.error) ? String(result.data.error) : '';

                    if (result.status === 409 && apiError.indexOf('nombre de la dupla en uso') !== -1) {
                        self.saveError = 'nombre de la dupla en uso';
                        self.nameFieldError = true;
                        return;
                    }

                    if (result.status === 409 || result.status === 400) {
                        self.saveError = apiError || 'No se pudo guardar la dupla. Intenta de nuevo.';
                        return;
                    }

                    self.saveError = 'No se pudo guardar la dupla. Intenta de nuevo.';
                }).catch(function () {
                    self.isSaving = false;
                    self.saveError = 'No se pudo guardar la dupla. Intenta de nuevo.';
                });
            },

            // E05
            openDeleteConfirm: function (dupla) {
                if (this.isLoading || this.isDeleting || !dupla) return;
                this.deleteTarget = { id: dupla.id, name: dupla.name };
                this.deleteError = null;
                this.modal.kind = 'delete';
                this.modal.visible = true;
            },

            // E06
            closeDeleteConfirm: function () {
                if (this.isDeleting) return;
                this.modal.visible = false;
                this.modal.kind = null;
                this.deleteTarget = null;
                this.deleteError = null;
            },

            // E07 — Cuando confirma eliminar dupla
            confirmDeleteDupla: function () {
                if (!this.deleteTarget || this.isDeleting) return;

                this.isDeleting = true;
                this.deleteError = null;

                var self = this;
                var duplaId = this.deleteTarget.id;

                fetch('/api/v1/duplas/' + encodeURIComponent(duplaId), {
                    method: 'DELETE',
                }).then(function (res) {
                    return res.json().then(function (data) {
                        return { ok: res.ok, status: res.status, data: data };
                    }).catch(function () {
                        return { ok: res.ok, status: res.status, data: {} };
                    });
                }).then(function (result) {
                    if (result.status === 401) {
                        window.location.href = '/static/landing.html';
                        return;
                    }

                    if (result.status === 200) {
                        self.isDeleting = false;
                        self.closeDeleteConfirm();
                        self.showToast('Dupla eliminada. Los profesionales quedaron disponibles.', 'success');
                        self.reloadScreen();
                        return;
                    }

                    if (result.status === 409) {
                        self.isDeleting = false;
                        self.modal.visible = false;
                        self.modal.kind = null;
                        self.deleteTarget = null;
                        self.deleteError = null;
                        self.blockedMessage = (result.data && result.data.message)
                            || 'Esta dupla no se puede eliminar porque se está utilizando en una sesión.';
                        self.modal.kind = 'deleteBlocked';
                        self.modal.visible = true;
                        return;
                    }

                    if (result.status === 404) {
                        self.isDeleting = false;
                        self.closeDeleteConfirm();
                        self.showToast('La dupla ya no existe o fue eliminada.', 'error');
                        self.reloadScreen();
                        return;
                    }

                    self.isDeleting = false;
                    self.deleteError = 'No se pudo eliminar la dupla. Intenta de nuevo.';
                }).catch(function () {
                    self.isDeleting = false;
                    self.deleteError = 'No se pudo eliminar la dupla. Intenta de nuevo.';
                });
            },

            // E08
            closeDeleteBlockedModal: function () {
                this.modal.visible = false;
                this.modal.kind = null;
                this.blockedMessage = null;
                this.deleteTarget = null;
            },

            logout: function () {
                logout('/seguridad/logout');
            },
        },
    });

    var app = home;
    window.app = app;
    window.home = home;
    home.mount('#app');
})();
