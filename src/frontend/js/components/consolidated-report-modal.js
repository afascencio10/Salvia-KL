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
            elapsedSeconds: 0,
            jobId: null,
            pollingInterval: null,
            jobReady: false
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
            this.jobReady = false;
            this.visible = true;

            // Verificar si hay un job pendiente en localStorage
            var savedJobId = localStorage.getItem('report_job_id');
            if (savedJobId) {
                this.jobId = savedJobId;
                this.loading = true;
                this.startPolling();
            }
        },
        close: function() {
            this.visible = false;
            this.stopProgressTimer();
            this.stopPolling();
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
            this.jobReady = false;
            this.startProgressTimer();

            // Iniciar generación asíncrona
            fetch('/api/v1/reportes/iniciar', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ start_date: this.startDate, end_date: this.endDate })
            })
            .then(async response => {
                if (!response.ok) {
                    const err = await response.json();
                    throw new Error(err.error || err.message || 'Error al iniciar reporte');
                }
                return response.json();
            })
            .then(data => {
                this.jobId = data.jobId;
                localStorage.setItem('report_job_id', data.jobId);
                this.startPolling();
            })
            .catch(error => {
                this.stopProgressTimer();
                this.loading = false;
                this.errorMessage = error.message || 'Error al conectar con el servidor.';
            });
        },

        startPolling: function() {
            this.stopPolling();
            this.pollingInterval = setInterval(() => {
                this.checkJobStatus();
            }, 5000);
            // Verificar inmediatamente
            this.checkJobStatus();
        },

        stopPolling: function() {
            if (this.pollingInterval) {
                clearInterval(this.pollingInterval);
                this.pollingInterval = null;
            }
        },

        checkJobStatus: function() {
            if (!this.jobId) return;
            fetch('/api/v1/reportes/estado/' + this.jobId)
                .then(res => res.json())
                .then(data => {
                    if (data.status === 'listo') {
                        this.stopPolling();
                        this.stopProgressTimer();
                        this.loading = false;
                        this.jobReady = true;
                        this.progressMessage = '¡Tu reporte está listo!';
                    } else if (data.status === 'error') {
                        this.stopPolling();
                        this.stopProgressTimer();
                        this.loading = false;
                        this.jobReady = false;
                        this.errorMessage = data.error || 'Error al generar el reporte';
                        localStorage.removeItem('report_job_id');
                        this.jobId = null;
                    } else {
                        // Sigue procesando
                        if (data.progress) this.progressMessage = data.progress;
                    }
                })
                .catch(() => {
                    // Error de red — no detenemos el polling, reintentará
                });
        },

        downloadReport: function() {
            if (!this.jobId) return;
            var self = this;
            fetch('/api/v1/reportes/descargar/' + this.jobId)
                .then(res => {
                    if (!res.ok) throw new Error('Error al descargar');
                    return res.blob();
                })
                .then(blob => {
                    var url = window.URL.createObjectURL(blob);
                    var a = document.createElement('a');
                    a.href = url;
                    a.download = 'reporte-seguimientos_' + self.startDate + '_' + self.endDate + '.xlsx';
                    document.body.appendChild(a);
                    a.click();
                    setTimeout(function() {
                        document.body.removeChild(a);
                        window.URL.revokeObjectURL(url);
                    }, 100);
                    localStorage.removeItem('report_job_id');
                    self.jobId = null;
                    self.jobReady = false;
                    self.visible = false;
                })
                .catch(err => {
                    self.errorMessage = err.message || 'Error al descargar el reporte';
                });
        }
    }
});
