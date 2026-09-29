/**
 * case-tasks
 * Componente reutilizable: gestión de tareas de un caso.
 * Muestra tareas pendientes y completadas, permite completarlas con modales según tipo.
 *
 * Props:
 *   - caseId   (String, required): icode del caso
 *   - userId   (String, optional): icode del usuario logueado (para permisos)
 *   - userRole (String, optional): rol del usuario ('sv', 'ro', 'op', 'ps', 'ts', 'et')
 *   - allowPsicoManage (Boolean, optional): si true, ps/ts pueden gestionar
 *     tareas asignadas a ellos (p. ej. detalle de remisión psicosocial).
 *     Default false: en detalle de caso ps/ts son solo lectura.
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
        psychosocialId: { type: String, default: '' },
        puedeGestionPropia: { type: Boolean, default: false },
        allowPsicoManage: { type: Boolean, default: false },
        entityBranchId: { default: null },
        entityName: { type: String, default: '' },
        barrierDescription: { type: String, default: '' },
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
        scopedTasks: function() {
            var self = this;
            // Rol Entidad ('et'): filtro original del equipo — solo sus tareas asignadas.
            if (this.userRole === 'et') {
                return this.tasks.filter(function(t) {
                    return String(t.assignedUserId || '') === String(self.userId || '');
                });
            }
            // Equipo Técnico ('te'): ve las tareas asignadas a sí mismo y, además, las
            // tareas "comite_caso" del pool (sin dueño), disponibles para que cualquier
            // integrante del equipo técnico las tome y gestione. El backend ya aplica
            // este mismo criterio; aquí se replica para mantener la vista consistente.
            if (this.userRole === 'te') {
                return this.tasks.filter(function(t) {
                    if (String(t.assignedUserId || '') === String(self.userId || '')) return true;
                    if (t.type === 'comite_caso' && !t.assignedUserId) return true;
                    return false;
                });
            }
            return this.tasks;
        },
        pendingTasks: function() {
            return this.scopedTasks.filter(function(t) { return t.status === 'ToDo'; });
        },
        completedTasks: function() {
            return this.scopedTasks.filter(function(t) { return t.status === 'Done'; });
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
                if (this.psychosocialId) url += '&psychosocialId=' + this.psychosocialId;
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
        abrirModalGestionPropia: function() {
            if (!this.$refs.taskModal) return;
            this.$refs.taskModal.openArticulacion({
                barrierId: this.barrierId,
                caseId: this.caseId,
                entityBranchId: this.entityBranchId,
                entityName: this.entityName || '',
                barrierDescription: this.barrierDescription || '',
            });
        },
        // esComiteDelPool: tarea "comite_caso" sin dueño que cualquier integrante del
        // Equipo Técnico ('te') puede tomar y gestionar.
        esComiteDelPool: function(task) {
            return task && task.type === 'comite_caso' && !task.assignedUserId;
        },
        puedeCompletar: function(task) {
            if ((this.userRole === 'ps' || this.userRole === 'ts') && !this.allowPsicoManage) return false;
            if (this.userRole === 'sv') return true;
            // El Equipo Técnico puede gestionar las comite_caso del pool.
            if (this.userRole === 'te' && this.esComiteDelPool(task)) return true;
            return String(task.assignedUserId || '') === String(this.userId || '');
        },
        esAsignadoOSv: function(task) {
            // Para tareas de validar_remision: solo el asignado puede gestionar (el modal requiere ser el profesional)
            if (task.type === 'validar_remision') {
                return String(task.assignedUserId || '') === String(this.userId || '');
            }
            // Para las demás tareas: sv siempre puede, o ser el asignado
            if (this.userRole === 'sv') return true;
            // El Equipo Técnico puede gestionar las comite_caso del pool aunque no le estén asignadas.
            if (this.userRole === 'te' && this.esComiteDelPool(task)) return true;
            return String(task.assignedUserId || '') === String(this.userId || '');
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
                'gestion_propia': 'Gestión propia',
                'articulacion_territorial': 'Articulación territorial',
                'compromiso_con_enlace': 'Compromiso con Enlace',
                'solicitud_entidad': 'Hacer solicitudes a la entidad',
                'respuesta_entidad': 'Responder solicitud de la entidad',
            };
            return labels[type] || type || task.category || 'Tarea';
        },
        tituloTarea: function(task) {
            if (task.type === 'solicitud_entidad') {
                return 'Hacer solicitudes a ' + (task.entityName || 'la entidad');
            }
            return task.description || this.labelTipo(task);
        },
        // entidadTarea resuelve el nombre de la entidad relacionada a la tarea.
        // El backend expone la entidad por dos rutas: 'entityName' (cuando la tarea
        // apunta a una entidad concreta del caso vía entity_case) y 'barrierEntityName'
        // (la sede/entidad ligada a la barrera). Se prioriza la primera.
        entidadTarea: function(task) {
            return (task && (task.entityName || task.barrierEntityName)) || '';
        },
        // solicitudTarea es el texto de la solicitud que se le hace a la entidad.
        // Corresponde a la descripción de la tarea del enlace territorial.
        solicitudTarea: function(task) {
            return (task && task.description) || '';
        },
        // esTareaEnlace indica si la tarjeta corresponde a una gestión hacia una
        // entidad (barrera/solicitud), para mostrar el bloque "Entidad / Solicitud".
        // Se muestra cuando hay una entidad resuelta o la tarea nace de una barrera.
        esTareaEnlace: function(task) {
            if (!task) return false;
            return !!this.entidadTarea(task) || !!task.barrierId;
        },
    },
    template: `
<div class="ct-container">
    <div v-if="loading" class="ct-loading"><i class="fa fa-spinner fa-spin"></i> Cargando tareas...</div>
    <template v-else>
        <!-- Tabs: Pendientes / Completadas + acción Registrar gestión propia -->
        <div v-if="tasks.length > 0 || puedeGestionPropia" style="display:flex;align-items:center;justify-content:space-between;gap:8px;margin-bottom:14px;border-bottom:1px solid #e5e7eb;padding-bottom:0">
            <div v-if="tasks.length > 0" style="display:flex;gap:4px">
                <button @click="activeTab = 'pending'" style="padding:8px 16px;font-size:.78rem;font-weight:600;border:none;cursor:pointer;border-bottom:2px solid transparent;background:none" :style="activeTab === 'pending' ? 'color:#5106A7;border-bottom-color:#5106A7' : 'color:#6b7280'">
                    Pendientes (\${ pendingTasks.length })
                </button>
                <button @click="activeTab = 'completed'" style="padding:8px 16px;font-size:.78rem;font-weight:600;border:none;cursor:pointer;border-bottom:2px solid transparent;background:none" :style="activeTab === 'completed' ? 'color:#5106A7;border-bottom-color:#5106A7' : 'color:#6b7280'">
                    Completadas (\${ completedTasks.length })
                </button>
            </div>
            <div v-else></div>
            <button v-if="puedeGestionPropia" @click="abrirModalGestionPropia" class="ct-btn-complete" style="flex-shrink:0;margin-bottom:8px">
                + Registrar gestión propia
            </button>
        </div>

        <!-- Lista de tareas pendientes -->
        <template v-if="activeTab === 'pending'">
            <div v-if="pendingTasks.length === 0" class="ct-empty">
                No hay tareas pendientes
            </div>
            <div v-else class="ct-list">
                <div v-for="task in pendingTasks" :key="task.id" class="ct-task" style="flex-direction:row;align-items:center;gap:12px;padding:14px 20px">
                    <div v-if="task.type === 'solicitud_entidad'" style="width:34px;height:34px;border-radius:8px;display:flex;align-items:center;justify-content:center;flex-shrink:0;background:#dbeafe">
                        <i class="fa fa-folder" style="color:#2563eb;font-size:.82rem"></i>
                    </div>
                    <div v-else-if="task.type === 'respuesta_entidad'" style="width:34px;height:34px;border-radius:8px;display:flex;align-items:center;justify-content:center;flex-shrink:0;background:#ecfdf5">
                        <i class="fa fa-reply" style="color:#059669;font-size:.82rem"></i>
                    </div>
                    <div v-else style="width:34px;height:34px;border-radius:8px;display:flex;align-items:center;justify-content:center;flex-shrink:0" :style="task.type && task.type.indexOf('oficio') !== -1 ? 'background:#f3f4f6' : task.type && (task.type.indexOf('llamar') !== -1 || task.type.indexOf('llamada') !== -1) ? 'background:#fce7f3' : task.type && task.type.indexOf('comite') !== -1 ? 'background:#fef9c3' : 'background:#ede9fe'">
                        <i :class="task.type && task.type.indexOf('oficio') !== -1 ? 'fa fa-file-alt' : task.type && (task.type.indexOf('llamar') !== -1 || task.type.indexOf('llamada') !== -1) ? 'fa fa-phone' : task.type && task.type.indexOf('comite') !== -1 ? 'fa fa-clipboard-list' : 'fa fa-tasks'" :style="task.type && task.type.indexOf('oficio') !== -1 ? 'color:#1f2937' : task.type && (task.type.indexOf('llamar') !== -1 || task.type.indexOf('llamada') !== -1) ? 'color:#db2777' : task.type && task.type.indexOf('comite') !== -1 ? 'color:#92400e' : 'color:#7c3aed'" style="font-size:.82rem"></i>
                    </div>
                    <div style="flex:1;min-width:0">
                        <div style="display:flex;align-items:center;gap:8px;margin-bottom:2px">
                            <span style="font-size:.7rem;font-weight:800;color:#1f2937;text-transform:uppercase;letter-spacing:.04em">\${ task.category }</span>
                            <span class="ct-badge ct-badge-todo">Pendiente</span>
                        </div>
                        <div style="font-size:.84rem;font-weight:700;color:#1f2937">\${ tituloTarea(task) }</div>
                        <!-- Entidad y solicitud: guían al enlace territorial sobre cómo gestionar -->
                        <div v-if="esTareaEnlace(task)" style="margin-top:4px;display:flex;flex-direction:column;gap:2px">
                            <div v-if="entidadTarea(task)" style="font-size:.74rem;color:#374151"><i class="fa fa-building" style="color:#2563eb;margin-right:5px"></i><strong>Entidad:</strong> \${ entidadTarea(task) }</div>
                            <div v-if="solicitudTarea(task)" style="font-size:.74rem;color:#374151"><i class="fa fa-file-lines" style="color:#7c3aed;margin-right:5px"></i><strong>Solicitud:</strong> \${ solicitudTarea(task) }</div>
                        </div>
                        <div style="font-size:.73rem;color:#9ca3af;margin-top:2px">Asignado a: \${ task.assignedUserName || 'Sin asignar' }</div>
                    </div>
                    <div style="display:flex;gap:6px;flex-shrink:0">
                        <button v-if="puedeCompletar(task) && esAsignadoOSv(task)" class="ct-btn-complete" @click="abrirModal(task)">Gestionar</button>
                        <button v-else-if="puedeCompletar(task) && !esAsignadoOSv(task)" class="ct-btn-complete" disabled style="opacity:.5;cursor:not-allowed" title="Solo el profesional asignado puede gestionar esta tarea">Gestionar</button>
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
                        <div style="font-size:.84rem;font-weight:700;color:#1f2937">\${ tituloTarea(task) }</div>
                        <!-- Entidad y solicitud: guían al enlace territorial sobre cómo gestionar -->
                        <div v-if="esTareaEnlace(task)" style="margin-top:4px;display:flex;flex-direction:column;gap:2px">
                            <div v-if="entidadTarea(task)" style="font-size:.74rem;color:#374151"><i class="fa fa-building" style="color:#2563eb;margin-right:5px"></i><strong>Entidad:</strong> \${ entidadTarea(task) }</div>
                            <div v-if="solicitudTarea(task)" style="font-size:.74rem;color:#374151"><i class="fa fa-file-lines" style="color:#7c3aed;margin-right:5px"></i><strong>Solicitud:</strong> \${ solicitudTarea(task) }</div>
                        </div>
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
