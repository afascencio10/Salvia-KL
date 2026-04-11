/**
 * seguimientos-tabla
 * Tabla de seguimientos con checkboxes, badges de riesgo/estado, y acción Ejecutar.
 * Emite: @reagendar(id), @paginar(page), @seleccionar(item)
 */
app.component('seguimientos-tabla', {
    delimiters: ['${', '}'],
    props: {
        items:   { type: Array,   default: () => [] },
        total:   { type: Number,  default: 0 },
        page:    { type: Number,  default: 0 },
        limit:   { type: Number,  default: 20 },
        loading: { type: Boolean, default: false },
    },
    data() {
        return {
            seleccionados: [],
            todosSel: false,
        };
    },
    computed: {
        totalPages() { return Math.ceil(this.total / this.limit) || 1; },
        paginaActual() { return this.page + 1; },
    },
    watch: {
        items() { this.seleccionados = []; this.todosSel = false; }
    },
    methods: {
        toggleTodos() {
            if (this.todosSel) {
                this.seleccionados = this.items.map(function(i) { return i.id; });
            } else {
                this.seleccionados = [];
            }
        },
        riesgoBadge(risk) {
            var map = {
                'EXTREMO':  'sa-badge sa-badge-extremo',
                'ALTO':     'sa-badge sa-badge-alto',
                'MODERADO': 'sa-badge sa-badge-moderado',
                'BAJO':     'sa-badge sa-badge-bajo',
            };
            return map[risk] || 'sa-badge sa-badge-default';
        },
        estadoBadge(status) {
            var map = {
                'PENDIENTE':    'sa-badge sa-badge-pendiente',
                'REALIZADO':    'sa-badge sa-badge-realizado',
                'VENCIDO':      'sa-badge sa-badge-vencido',
                'REPROGRAMADO': 'sa-badge sa-badge-reprogramado',
            };
            return map[status] || 'sa-badge sa-badge-default';
        },
        formatFecha(d) {
            if (!d) return '-';
            return new Date(d).toLocaleDateString('es-CO', { year: 'numeric', month: '2-digit', day: '2-digit' });
        },
        formatHora(d) {
            if (!d) return '-';
            return new Date(d).toLocaleTimeString('es-CO', { hour: '2-digit', minute: '2-digit' });
        },
        irPagina(p) {
            if (p >= 0 && p < this.totalPages) this.$emit('paginar', p);
        },
    },
    template: `
    <div>
        <!-- Loading -->
        <div v-if="loading" class="text-center py-4">
            <div class="spinner-border text-primary" role="status"><span class="sr-only">Cargando...</span></div>
        </div>

        <!-- Tabla -->
        <div v-else class="table-responsive">
            <table class="table table-hover mb-0" style="font-size: 14px;">
                <thead style="background: #f9fafb;">
                    <tr>
                        <th style="width: 30px;">
                            <input type="checkbox" v-model="todosSel" @change="toggleTodos">
                        </th>
                        <th>Caso <i class="fas fa-sort text-muted" style="font-size:10px;"></i></th>
                        <th>Riesgo <i class="fas fa-sort text-muted" style="font-size:10px;"></i></th>
                        <th>Agente <i class="fas fa-sort text-muted" style="font-size:10px;"></i></th>
                        <th>Fecha ↑</th>
                        <th>Hora <i class="fas fa-sort text-muted" style="font-size:10px;"></i></th>
                        <th>Estado <i class="fas fa-sort text-muted" style="font-size:10px;"></i></th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-if="items.length === 0">
                        <td colspan="7" class="text-center text-muted py-4">
                            <i class="fas fa-inbox fa-2x mb-2 d-block"></i>
                            No se encontraron seguimientos
                        </td>
                    </tr>
                    <tr v-for="item in items" :key="item.id"
                        style="cursor: pointer;"
                        @click="$emit('seleccionar', item)">
                        <td @click.stop>
                            <input type="checkbox" :value="item.id" v-model="seleccionados">
                        </td>
                        <td>
                            <div style="font-weight: 600;">\${ item.case_id }</div>
                            <div class="text-muted" style="font-size: 12px;">SAL-\${ item.sequence_number }</div>
                        </td>
                        <td><span :class="riesgoBadge(item.risk_status)">\${ item.risk_status }</span></td>
                        <td>\${ item.agent_id }</td>
                        <td>\${ formatFecha(item.scheduled_date) }</td>
                        <td>\${ formatHora(item.scheduled_date) }</td>
                        <td>
                            <span :class="estadoBadge(item.status)">\${ item.status }</span>
                            <a v-if="item.status === 'PENDIENTE'" href="#"
                               class="ml-2 text-primary small"
                               @click.stop.prevent="$emit('reagendar', item.id)"
                               title="Re agendar y priorizar">
                                Ejecutar
                            </a>
                        </td>
                    </tr>
                </tbody>
            </table>
        </div>

        <!-- Paginación -->
        <div v-if="total > 0" class="d-flex justify-content-between align-items-center px-3 py-2" style="border-top: 1px solid #e5e7eb;">
            <small class="text-muted">
                Página \${ paginaActual } de \${ totalPages } — \${ total } registros
            </small>
            <nav>
                <ul class="pagination pagination-sm mb-0">
                    <li class="page-item" :class="{ disabled: page === 0 }">
                        <a class="page-link" href="#" @click.prevent="irPagina(page - 1)">Anterior</a>
                    </li>
                    <li class="page-item" :class="{ disabled: paginaActual >= totalPages }">
                        <a class="page-link" href="#" @click.prevent="irPagina(page + 1)">Siguiente</a>
                    </li>
                </ul>
            </nav>
        </div>
    </div>
    `
});
