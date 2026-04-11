/**
 * seguimiento-card
 * Componente reutilizable que muestra el resumen de un seguimiento.
 * Props:
 *   - titulo    (String)  : título del seguimiento
 *   - estado    (String)  : 'pendiente' | 'en_progreso' | 'completado'
 *   - fecha     (String)  : fecha programada
 *   - descripcion (String): descripción breve
 */
app.component('seguimiento-card', {
    delimiters: ['${', '}'],
    props: {
        titulo:      { type: String, required: true },
        estado:      { type: String, default: 'pendiente' },
        fecha:       { type: String, default: '' },
        descripcion: { type: String, default: '' },
    },
    computed: {
        estadoClass() {
            const map = {
                'pendiente':    'badge bg-warning text-dark',
                'en_progreso':  'badge bg-primary',
                'completado':   'badge bg-success',
            };
            return map[this.estado] || 'badge bg-secondary';
        },
        estadoLabel() {
            const map = {
                'pendiente':   'Pendiente',
                'en_progreso': 'En progreso',
                'completado':  'Completado',
            };
            return map[this.estado] || this.estado;
        }
    },
    template: `
        <div class="card shadow-sm mb-3">
            <div class="card-header d-flex justify-content-between align-items-center">
                <strong>\${titulo}</strong>
                <span :class="estadoClass">\${estadoLabel}</span>
            </div>
            <div class="card-body">
                <p class="card-text text-muted mb-1" v-if="descripcion">\${descripcion}</p>
                <small class="text-muted" v-if="fecha">
                    <i class="fas fa-calendar-alt mr-1"></i> \${fecha}
                </small>
            </div>
            <div class="card-footer text-right">
                <button class="btn btn-sm btn-outline-primary" @click="$emit('seleccionar', titulo)">
                    Ver detalle <i class="fas fa-arrow-right ml-1"></i>
                </button>
            </div>
        </div>
    `
});
