/**
 * seguimientos-filtros
 * Tabs + filtros inline para Seguimientos del Área.
 * Emite: @filtrar(filtros), @cambiar-tab(tab)
 */
app.component('seguimientos-filtros', {
    delimiters: ['${', '}'],
    props: {
        agentes: { type: Array, default: () => [] },
        estados: { type: Array, default: () => [] },
        total:   { type: Number, default: 0 },
    },
    data() {
        return {
            tabActiva: 'pendientes',
            fechaInicio: '',
            fechaFin: '',
            estadoSel: '',
            agenteSel: '',
            busqueda: '',
        };
    },
    computed: {
        tabs() {
            return [
                { key: 'pendientes',  label: 'Vista rápida' },
                { key: 'esta_semana', label: 'Esta semana' },
                { key: 'por_fecha',   label: 'Por fecha' },
                { key: 'por_buscar',  label: 'Por buscar' },
                { key: 'historico',   label: 'Histórico' },
            ];
        },
        mostrarFiltrosFecha() {
            return this.tabActiva === 'por_fecha' || this.tabActiva === 'historico';
        },
        mostrarBusqueda() {
            return this.tabActiva === 'por_buscar';
        }
    },
    methods: {
        seleccionarTab(key) {
            this.tabActiva = key;
            this.$emit('cambiar-tab', key);
            this.emitirFiltros();
        },
        emitirFiltros() {
            var tab = this.tabActiva;
            if (tab === 'esta_semana') {
                var hoy = new Date();
                var lunes = new Date(hoy);
                lunes.setDate(hoy.getDate() - hoy.getDay() + 1);
                var domingo = new Date(lunes);
                domingo.setDate(lunes.getDate() + 6);
                this.fechaInicio = lunes.toISOString().split('T')[0];
                this.fechaFin = domingo.toISOString().split('T')[0];
                tab = '';
            }
            if (tab === 'por_buscar') tab = '';
            if (tab === 'por_fecha') tab = '';
            if (tab === 'historico') tab = 'todos';

            this.$emit('filtrar', {
                tab: tab,
                fecha_inicio: this.fechaInicio,
                fecha_fin: this.fechaFin,
                estado: this.estadoSel,
                agente_id: this.agenteSel,
                busqueda: this.busqueda,
            });
        },
    },
    template: `
    <div class="mb-3">
        <!-- Tabs -->
        <div class="d-flex align-items-center mb-2" style="gap: 4px; flex-wrap: wrap;">
            <button v-for="t in tabs" :key="t.key"
                    class="btn btn-sm"
                    :class="tabActiva === t.key ? 'btn-primary' : 'btn-outline-secondary'"
                    @click="seleccionarTab(t.key)"
                    style="border-radius: 20px; font-size: 13px; padding: 4px 14px;">
                \${ t.label }
            </button>
        </div>

        <!-- Filtros inline -->
        <div class="d-flex align-items-center" style="gap: 8px; flex-wrap: wrap;">
            <span class="text-muted small">Filtros:</span>

            <!-- Fechas (solo en tabs que lo requieren) -->
            <template v-if="mostrarFiltrosFecha">
                <input type="date" class="form-control form-control-sm" style="width: 140px;"
                       v-model="fechaInicio" @change="emitirFiltros">
                <input type="date" class="form-control form-control-sm" style="width: 140px;"
                       v-model="fechaFin" @change="emitirFiltros">
            </template>

            <!-- Búsqueda (solo en tab buscar) -->
            <template v-if="mostrarBusqueda">
                <input type="text" class="form-control form-control-sm" style="width: 200px;"
                       placeholder="dd/mm/yyyy" v-model="busqueda" @input="emitirFiltros">
            </template>

            <!-- Estado -->
            <select class="form-control form-control-sm" style="width: 150px;"
                    v-model="estadoSel" @change="emitirFiltros">
                <option value="">Todos los estados</option>
                <option v-for="e in estados" :key="e" :value="e">\${ e }</option>
            </select>

            <!-- Agente -->
            <select class="form-control form-control-sm" style="width: 160px;"
                    v-model="agenteSel" @change="emitirFiltros">
                <option value="">Todos los agentes</option>
                <option v-for="a in agentes" :key="a.agent_id" :value="a.agent_id">\${ a.agent_id }</option>
            </select>

            <!-- Contador -->
            <span class="ml-auto text-muted small">\${ total } resultados</span>
        </div>
    </div>
    `
});
