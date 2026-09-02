(function () {
    'use strict';

    var cfg = window.__CasoEntidadDetalleConfig || {};
    var entityCaseId = cfg.entityCaseId || '';

    var home = Vue.createApp({
        delimiters: ['${', '}'],
        data: function () {
            return {
                detail: null,
                loading: true,
                error: null,
                errorTitle: '',
                errorMessage: '',
                errorCode: 0,
                activeTab: 'info',
                toast: { visible: false, message: '', type: 'info' }
            };
        },
        mounted: function () {
            document.title = cfg.windowTitle || 'Detalle Caso Entidad';
            var appEl = document.getElementById('app');
            if (appEl) appEl.style.display = 'block';
            this.fetchDetail();
        },
        methods: {
            fetchDetail: function () {
                var self = this;
                self.loading = true;
                self.error = null;

                if (!entityCaseId) {
                    self.loading = false;
                    self.error = 'ID no proporcionado';
                    self.errorTitle = 'Identificador no válido';
                    self.errorMessage = 'No se proporcionó un identificador de caso-entidad válido.';
                    self.errorCode = 400;
                    return;
                }

                fetch('/api/v1/entity-cases/' + encodeURIComponent(entityCaseId) + '/detail', {
                    method: 'GET',
                    credentials: 'same-origin',
                    headers: { 'Accept': 'application/json' }
                })
                .then(function (res) {
                    if (!res.ok) {
                        self.errorCode = res.status;
                        if (res.status === 404) {
                            self.errorTitle = 'Relación no encontrada';
                            self.errorMessage = 'El caso-entidad solicitado no existe o fue eliminado.';
                        } else if (res.status === 401) {
                            self.errorTitle = 'No autorizado';
                            self.errorMessage = 'Su sesión expiró o no tiene permisos para ver este detalle.';
                        } else {
                            self.errorTitle = 'Error del servidor';
                            self.errorMessage = 'Ocurrió un problema al obtener la información. Intente nuevamente.';
                        }
                        throw new Error('HTTP ' + res.status);
                    }
                    return res.json();
                })
                .then(function (data) {
                    self.detail = data;
                    self.loading = false;
                })
                .catch(function (err) {
                    self.error = err.message || 'Error desconocido';
                    self.loading = false;
                });
            },

            statusBadgeClass: function (label) {
                if (!label) return '';
                var l = label.toLowerCase();
                if (l.indexOf('activo') !== -1 || l.indexOf('aprobado') !== -1) return 'ced-badge-status--activo';
                if (l.indexOf('cerrado') !== -1) return 'ced-badge-status--cerrado';
                if (l.indexOf('vencido') !== -1) return 'ced-badge-status--vencido';
                if (l.indexOf('novedad') !== -1) return 'ced-badge-status--novedad';
                return '';
            },

            showToast: function (message, type) {
                var self = this;
                self.toast = { visible: true, message: message, type: type || 'info' };
                setTimeout(function () { self.toast.visible = false; }, 4000);
            }
        }
    });

    home.mount('#app');
})();
