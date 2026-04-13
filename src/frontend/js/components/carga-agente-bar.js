/**
 * carga-agente-bar
 * Barra horizontal con pills de carga por agente (mockup: CARGA POR AGENTE).
 * Props:
 *   - agentes (Array): [{agent_id, total}]
 *   - fecha   (String): fecha seleccionada
 */
app.component('carga-agente-bar', {
    delimiters: ['${', '}'],
    props: {
        agentes: { type: Array, default: () => [] },
        fecha:   { type: String, default: '' },
    },
    data() {
        return {
            fechaSel: this.fecha || new Date().toISOString().split('T')[0],
        };
    },
    watch: {
        fecha(v) { this.fechaSel = v; }
    },
    methods: {
        pillColor(total) {
            var colors = ['#8b5cf6', '#f59e0b', '#ef4444', '#22c55e', '#3b82f6', '#ec4899', '#14b8a6', '#f97316'];
            var idx = Math.abs(total) % colors.length;
            return colors[idx];
        },
        cambiarFecha() {
            this.$emit('cambiar-fecha', this.fechaSel);
        }
    },
    template: `
    <div class="d-flex align-items-center mb-3 p-2 rounded" style="background: #f9fafb; border: 1px solid #e5e7eb; gap: 10px; flex-wrap: wrap;">
        <span class="text-muted small font-weight-bold" style="text-transform: uppercase; letter-spacing: 0.05em;">
            Carga por agente:
        </span>
        <input type="date" class="form-control form-control-sm" style="width: 140px;"
               v-model="fechaSel" @change="cambiarFecha">

        <div class="d-flex" style="gap: 6px; flex-wrap: wrap;">
            <span v-if="agentes.length === 0" class="text-muted small">Sin datos</span>
            <span v-for="a in agentes" :key="a.agent_id"
                  class="d-inline-flex align-items-center px-2 py-1 rounded-pill text-white"
                  :style="{ background: pillColor(a.total), fontSize: '12px', gap: '4px', fontWeight: '600' }">
                \${ a.agent_id }
                <span class="d-inline-flex align-items-center justify-content-center rounded-circle bg-white"
                      style="width: 20px; height: 20px; font-size: 11px; font-weight: 700; color: #374151;">
                    \${ a.total }
                </span>
                <span style="font-weight: 400; opacity: 0.9;">pend.</span>
            </span>
        </div>
    </div>
    `
});
