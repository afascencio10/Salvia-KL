app.component('consolidated-report-modal', {
    template: window.__consolidatedReportModalTpl || '#tpl-consolidated-report-modal',
    delimiters: ['${', '}'],
    props: {
        currentRole: {
            type: String,
            required: true
        }
    },
    data: function() {
        return {
            visible: false,
            startDate: '',
            endDate: '',
            loading: false,
            errorMessage: '',
            progressMessage: 'Consultando base de datos...',
            progressInterval: null,
            elapsedSeconds: 0
        };
    },
    methods: {
        open: function() {
            const today = new Date();
            
            // Fecha final: hoy (YYYY-MM-DD local)
            const yearEnd = today.getFullYear();
            const monthEnd = String(today.getMonth() + 1).padStart(2, '0');
            const dayEnd = String(today.getDate()).padStart(2, '0');
            this.endDate = `${yearEnd}-${monthEnd}-${dayEnd}`;
            
            // Fecha inicial: hoy menos 1 mes
            const oneMonthAgo = new Date();
            oneMonthAgo.setMonth(today.getMonth() - 1);
            
            const yearStart = oneMonthAgo.getFullYear();
            const monthStart = String(oneMonthAgo.getMonth() + 1).padStart(2, '0');
            const dayStart = String(oneMonthAgo.getDate()).padStart(2, '0');
            this.startDate = `${yearStart}-${monthStart}-${dayStart}`;

            this.errorMessage = '';
            this.loading = false;
            this.visible = true;
        },
        close: function() {
            if (this.loading) {
                return; // Bloquear cierre durante carga
            }
            this.visible = false;
            this.stopProgressTimer();
        },
        onBackdropClick: function() {
            this.close();
        },
        startProgressTimer: function() {
            this.elapsedSeconds = 0;
            this.progressMessage = 'Consultando base de datos...';
            this.progressInterval = setInterval(() => {
                this.elapsedSeconds++;
                if (this.elapsedSeconds <= 3) {
                    this.progressMessage = 'Consultando base de datos...';
                } else if (this.elapsedSeconds <= 7) {
                    this.progressMessage = 'Escribiendo registros en las hojas de cálculo...';
                } else if (this.elapsedSeconds <= 12) {
                    this.progressMessage = 'Comprimiendo archivo para descarga...';
                } else {
                    this.progressMessage = 'Iniciando descarga en el navegador...';
                }
            }, 1000);
        },
        stopProgressTimer: function() {
            if (this.progressInterval) {
                clearInterval(this.progressInterval);
                this.progressInterval = null;
            }
        },
        submitReport: function() {
            this.errorMessage = '';
            
            // Validaciones locales
            if (!this.startDate || !this.endDate) {
                this.errorMessage = 'Ambas fechas son obligatorias.';
                return;
            }

            const start = new Date(this.startDate);
            const end = new Date(this.endDate);

            if (start > end) {
                this.errorMessage = 'La fecha inicial no puede ser posterior a la fecha final.';
                return;
            }

            const diffTime = Math.abs(end - start);
            const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));
            if (diffDays > 366) {
                this.errorMessage = 'El rango no puede superar 1 año (366 días).';
                return;
            }

            this.loading = true;
            this.startProgressTimer();

            fetch('/api/v1/reportes/seguimientos-consolidado', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    start_date: this.startDate,
                    end_date: this.endDate
                })
            })
            .then(async response => {
                if (!response.ok) {
                    // Si falla, parseamos el error JSON
                    const errorText = await response.text();
                    let errMsg = 'Error interno al generar el reporte.';
                    try {
                        const errJson = JSON.parse(errorText);
                        if (errJson && errJson.message) {
                            errMsg = errJson.message;
                        }
                    } catch (e) {
                        // Usar valor por defecto
                    }
                    throw new Error(errMsg);
                }
                return response.blob();
            })
            .then(blob => {
                this.stopProgressTimer();
                this.loading = false;
                
                // Forzar descarga del archivo
                const url = window.URL.createObjectURL(blob);
                const a = document.createElement('a');
                a.href = url;
                a.download = `reporte-seguimientos_${this.startDate}_${this.endDate}.xlsx`;
                document.body.appendChild(a);
                a.click();
                
                // Limpieza
                setTimeout(() => {
                    document.body.removeChild(a);
                    window.URL.revokeObjectURL(url);
                }, 100);

                this.visible = false; // Cerrar el modal al completarse con éxito
            })
            .catch(error => {
                this.stopProgressTimer();
                this.loading = false;
                this.errorMessage = error.message || 'Error al conectar con el servidor.';
            });
        }
    }
});
