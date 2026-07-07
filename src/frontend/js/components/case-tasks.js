/**
 * case-tasks
 * Componente reutilizable: gestión de tareas de un caso.
 * Muestra tareas pendientes y completadas, permite completarlas con modales según tipo.
 *
 * Props:
 *   - caseId   (String, required): icode del caso
 *   - userId   (String, optional): icode del usuario logueado (para permisos)
 *   - userRole (String, optional): rol del usuario ('sv', 'ro', 'op')
 * Events:
 *   - @task-completed: emitido cuando se completa una tarea
 */

(function injectCaseTasksStyles() {
    if (document.getElementById('case-tasks-styles')) return;
    var style = document.createElement('style');
    style.id = 'case-tasks-styles';
    style.textContent = `
        .ct-container { padding: 16px 0 }
        .ct-summary { display:flex;gap:16px;margin-bottom:20px }
        .ct-card { flex:1;border:1px solid #e5e7eb;border-radius:10px;padding:16px;text-align:center }
        .ct-card-label { font-size:.72rem;color:#6b7280;text-transform:uppercase;font-weight:600;margin:0 0 4px }
        .ct-card-value { font-size:1.4rem;font-weight:700;margin:0 }
        .ct-card-pending .ct-card-value { color:#f97316 }
        .ct-card-done .ct-card-value { color:#22c55e }
        .ct-list { display:flex;flex-direction:column;gap:10px }
        .ct-list::before { content:none;display:none }
        .ct-task { border:none!important;border-radius:10px;padding:14px 16px;display:flex;align-items:center;justify-content:space-between;transition:background .12s;background:#f9fafb;box-shadow:none }
        .ct-task:hover { background:#f9fafb }
        .ct-task-info { flex:1 }
        .ct-task-category { font-size:.68rem;color:#9ca3af;text-transform:uppercase;font-weight:600 }
        .ct-task-desc { font-size:.84rem;color:#374151;margin:4px 0 0;font-weight:500 }
        .ct-task-meta { font-size:.72rem;color:#9ca3af;margin:4px 0 0 }
        .ct-badge { display:inline-block;padding:2px 8px;border-radius:10px;font-size:.68rem;font-weight:600 }
        .ct-badge-todo { background:#fff7ed;color:#f97316;border:1px solid #fed7aa }
        .ct-badge-done { background:#f0fdf4;color:#22c55e;border:1px solid #bbf7d0 }
        .ct-btn-complete { background:#5106A7;color:#fff;border:none;padding:6px 14px;border-radius:6px;font-size:.78rem;font-weight:600;cursor:pointer }
        .ct-btn-complete:hover { background:#3d0480 }
        .ct-btn-complete:disabled { opacity:.5;cursor:not-allowed }
        .ct-empty { text-align:center;padding:40px 0;color:#9ca3af;font-size:.85rem }
        .ct-loading { text-align:center;padding:40px 0;color:#9ca3af }
        .ct-modal-backdrop { position:fixed;inset:0;background:rgba(0,0,0,.5);display:flex;align-items:center;justify-content:center;z-index:10000 }
        .ct-modal { background:#fff;border-radius:12px;width:95%;max-width:550px;padding:24px;box-shadow:0 20px 50px rgba(0,0,0,.2) }
        .ct-modal-title { font-size:1rem;font-weight:700;margin:0 0 16px;color:#111827 }
        .ct-modal-field { margin-bottom:14px }
        .ct-modal-label { font-size:.75rem;font-weight:600;color:#6b7280;margin-bottom:4px;display:block }
        .ct-modal-input { width:100%;border:1px solid #d1d5db;border-radius:8px;padding:10px 12px;font-size:.84rem;resize:vertical;min-height:80px }
        .ct-modal-footer { display:flex;justify-content:flex-end;gap:10px;margin-top:16px }
        .ct-modal-btn-cancel { background:#f3f4f6;color:#374151;border:none;padding:8px 16px;border-radius:6px;font-size:.82rem;cursor:pointer }
        .ct-modal-btn-save { background:#5106A7;color:#fff;border:none;padding:8px 16px;border-radius:6px;font-size:.82rem;font-weight:600;cursor:pointer }
        .ct-modal-btn-save:disabled { opacity:.5;cursor:not-allowed }
        .ct-tabs { display:flex;gap:0;border-bottom:2px solid #e5e7eb;margin-bottom:16px }
        .ct-tab { padding:8px 16px;font-size:.8rem;font-weight:600;color:#6b7280;cursor:pointer;border-bottom:2px solid transparent;margin-bottom:-2px }
        .ct-tab.active { color:#5106A7;border-bottom-color:#5106A7 }
    `;
    document.head.appendChild(style);
})();

app.component('case-tasks', {
    delimiters: ['${', '}'],
    props: {
        caseId:   { type: String, required: true },
        userId:   { type: String, default: '' },
        userRole: { type: String, default: '' },
        barrierId: { type: String, default: '' },
    },
    emits: ['task-completed'],
    data: function() {
        return {
            loading: false,
            tasks: [],
            activeTab: 'pending',
            reasignarTask: null,
            reasignarAgente: '',
            agentesEquipo: [],
            reasignando: false,
        };
    },
    computed: {
        pendingTasks: function() {
            return this.tasks.filter(function(t) { return t.status === 'ToDo'; });
        },
        completedTasks: function() {
            return this.tasks.filter(function(t) { return t.status === 'Done'; });
        },
        visibleTasks: function() {
            return this.activeTab === 'pending' ? this.pendingTasks : this.completedTasks;
        },
    },
    mounted: function() {
        this.cargar();
    },
    methods: {
        async cargar() {
            this.loading = true;
            try {
                var url = '/api/v1/case-tasks?caseId=' + this.caseId;
                if (this.barrierId) url += '&barrierId=' + this.barrierId;
                var res = await fetch(url);
                if (!res.ok) throw new Error('Error ' + res.status);
                this.tasks = await res.json();
                if (!Array.isArray(this.tasks)) this.tasks = [];
            } catch (e) {
                console.error('Error cargando tareas:', e);
                this.tasks = [];
            } finally {
                this.loading = false;
            }
        },
        abrirModal: function(task) {
            if (this.$refs.taskModal) {
                this.$refs.taskModal.open(task.id);
            }
        },
        onTaskCompleted: function() {
            this.cargar();
            this.$emit('task-completed');
        },
        puedeCompletar: function(task) {
            if (this.userRole === 'ps' || this.userRole === 'ts') return false;
            if (this.userRole === 'sv') return true;
            return task.assignedUserId === this.userId;
        },
        async abrirReasignar(task) {
            this.reasignarTask = task;
            this.reasignarAgente = '';
            this.reasignando = false;
            // Cargar agentes del equipo
            try {
                var res = await fetch('/api/v1/equipo-operadores');
                if (res.ok) {
                    this.agentesEquipo = await res.json();
                    if (!Array.isArray(this.agentesEquipo)) this.agentesEquipo = [];
                }
            } catch(e) { this.agentesEquipo = []; }
        },
        async confirmarReasignar() {
            if (!this.reasignarAgente || !this.reasignarTask) return;
            this.reasignando = true;
            try {
                var res = await fetch('/api/v1/case-tasks/' + this.reasignarTask.id + '/reassign', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ assignedUserId: this.reasignarAgente })
                });
                if (!res.ok) throw new Error('Error');
                this.reasignarTask = null;
                this.cargar();
            } catch(e) {
                console.error('Error reasignando:', e);
            } finally {
                this.reasignando = false;
            }
        },
        formatFecha: function(fecha) {
            if (!fecha) return '';
            return fecha.substring(0, 10);
        },
        labelTipo: function(task) {
            var type = task.type || '';
            var labels = {
                'escribir_oficio_cierre': 'Escribir oficio cierre de caso',
                'escribir_oficio_barrera': 'Escribir oficio barrera',
                'corregir_oficio_barrera': 'Corregir oficio barrera',
                'llamar_entidad_barrera': 'Llamar a entidad barrera',
                'gestion_enlace_territorial': 'Gestión barrera enlace territorial',
            };
            return labels[type] || type || task.category || 'Tarea';
        },
    },
    template: `
<div class="ct-container">
    <div v-if="loading" class="ct-loading"><i class="fa fa-spinner fa-spin"></i> Cargando tareas...</div>
    <template v-else>
        <!-- Tabs: Pendientes / Completadas -->
        <div v-if="tasks.length > 0" style="display:flex;gap:4px;margin-bottom:14px;border-bottom:1px solid #e5e7eb;padding-bottom:0">
            <button @click="activeTab = 'pending'" style="padding:8px 16px;font-size:.78rem;font-weight:600;border:none;cursor:pointer;border-bottom:2px solid transparent;background:none" :style="activeTab === 'pending' ? 'color:#5106A7;border-bottom-color:#5106A7' : 'color:#6b7280'">
                Pendientes (\${ pendingTasks.length })
            </button>
            <button @click="activeTab = 'completed'" style="padding:8px 16px;font-size:.78rem;font-weight:600;border:none;cursor:pointer;border-bottom:2px solid transparent;background:none" :style="activeTab === 'completed' ? 'color:#5106A7;border-bottom-color:#5106A7' : 'color:#6b7280'">
                Completadas (\${ completedTasks.length })
            </button>
        </div>

        <!-- Lista de tareas pendientes -->
        <template v-if="activeTab === 'pending'">
            <div v-if="pendingTasks.length === 0" class="ct-empty">
                No hay tareas pendientes
            </div>
            <div v-else class="ct-list">
                <div v-for="task in pendingTasks" :key="task.id" class="ct-task" style="flex-direction:row;align-items:center;gap:12px;padding:14px 20px">
                    <div style="width:34px;height:34px;border-radius:8px;display:flex;align-items:center;justify-content:center;flex-shrink:0" :style="task.type && task.type.indexOf('oficio') !== -1 ? 'background:#f3f4f6' : task.type && (task.type.indexOf('llamar') !== -1 || task.type.indexOf('llamada') !== -1) ? 'background:#fce7f3' : task.type && task.type.indexOf('comite') !== -1 ? 'background:#fef9c3' : 'background:#ede9fe'">
                        <i :class="task.type && task.type.indexOf('oficio') !== -1 ? 'fa fa-file-alt' : task.type && (task.type.indexOf('llamar') !== -1 || task.type.indexOf('llamada') !== -1) ? 'fa fa-phone' : task.type && task.type.indexOf('comite') !== -1 ? 'fa fa-clipboard-list' : 'fa fa-tasks'" :style="task.type && task.type.indexOf('oficio') !== -1 ? 'color:#1f2937' : task.type && (task.type.indexOf('llamar') !== -1 || task.type.indexOf('llamada') !== -1) ? 'color:#db2777' : task.type && task.type.indexOf('comite') !== -1 ? 'color:#92400e' : 'color:#7c3aed'" style="font-size:.82rem"></i>
                    </div>
                    <div style="flex:1;min-width:0">
                        <div style="display:flex;align-items:center;gap:8px;margin-bottom:2px">
                            <span style="font-size:.7rem;font-weight:800;color:#1f2937;text-transform:uppercase;letter-spacing:.04em">\${ task.category }</span>
                            <span class="ct-badge ct-badge-todo">Pendiente</span>
                        </div>
                        <div style="font-size:.84rem;font-weight:700;color:#1f2937">\${ task.description || labelTipo(task) }</div>
                        <div style="font-size:.73rem;color:#9ca3af;margin-top:2px">Asignado a: \${ task.assignedUserName || 'Sin asignar' }</div>
                    </div>
                    <div style="display:flex;gap:6px;flex-shrink:0">
                        <button v-if="puedeCompletar(task)" class="ct-btn-complete" @click="abrirModal(task)">Gestionar</button>
                        <button v-if="false && userRole === 'sv'" class="ct-btn-complete" style="background:#3b82f6" @click="abrirReasignar(task)">Reasignar</button>
                    </div>
                </div>
            </div>
        </template>

        <!-- Lista de tareas completadas -->
        <template v-if="activeTab === 'completed'">
            <div v-if="completedTasks.length === 0" class="ct-empty">
                No hay tareas completadas aún
            </div>
            <div v-else class="ct-list">
                <div v-for="task in completedTasks" :key="task.id" class="ct-task" style="flex-direction:row;align-items:center;gap:12px;padding:14px 20px;background:#f9fafb">
                    <div style="width:34px;height:34px;border-radius:8px;display:flex;align-items:center;justify-content:center;flex-shrink:0;background:#dcfce7">
                        <i class="fa fa-check" style="color:#16a34a;font-size:.82rem"></i>
                    </div>
                    <div style="flex:1;min-width:0">
                        <div style="display:flex;align-items:center;gap:8px;margin-bottom:2px">
                            <span style="font-size:.7rem;font-weight:800;color:#1f2937;text-transform:uppercase;letter-spacing:.04em">\${ task.category }</span>
                            <span class="ct-badge ct-badge-done">Completada</span>
                        </div>
                        <div style="font-size:.84rem;font-weight:700;color:#1f2937">\${ task.description || labelTipo(task) }</div>
                        <div style="font-size:.73rem;color:#9ca3af;margin-top:2px">Completada por: \${ task.assignedUserName || '—' }</div>
                    </div>
                    <div style="display:flex;gap:6px;flex-shrink:0">
                        <button class="ct-btn-complete" style="background:#16a34a" @click="$refs.taskHistory.open(task.id)">Ver detalle</button>
                    </div>
                </div>
            </div>
        </template>
    </template>

    <!-- Modal de tareas se maneja desde el padre via @open-task-modal -->
    <!-- Modal ver detalle tarea completada (autosuficiente) -->
    <case-task-history ref="taskHistory"></case-task-history>
    <!-- Modal gestionar tarea pendiente (autosuficiente) -->
    <case-task-modal ref="taskModal" :current-user-id="userId" @completed="onTaskCompleted"></case-task-modal>
    <!-- Modal reasignar tarea -->
    <div v-if="reasignarTask" class="ct-modal-backdrop" @click.self="reasignarTask = null">
        <div class="ct-modal">
            <p class="ct-modal-title">Reasignar tarea</p>
            <div class="ct-modal-field">
                <label class="ct-modal-label">Seleccionar agente</label>
                <select v-model="reasignarAgente" class="ct-modal-input" style="min-height:auto;resize:none;padding:8px 12px">
                    <option value="" disabled>Seleccione un agente...</option>
                    <option v-for="ag in agentesEquipo" :key="ag.icode" :value="ag.icode">\${ ag.fullName || ag.icode }</option>
                </select>
            </div>
            <div class="ct-modal-footer">
                <button class="ct-modal-btn-cancel" @click="reasignarTask = null">Cancelar</button>
                <button class="ct-modal-btn-save" :disabled="!reasignarAgente || reasignando" @click="confirmarReasignar">
                    \${ reasignando ? 'Reasignando...' : 'Reasignar' }
                </button>
            </div>
        </div>
    </div>
</div>
    `
});
