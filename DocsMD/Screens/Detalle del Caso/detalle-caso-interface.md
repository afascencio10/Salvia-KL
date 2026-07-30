# `Detalle del Caso` — Interfaz de la Pantalla

Vista completa de un caso VBG: datos de la víctima, hechos, derivaciones, gestión institucional (entidades + barreras + oficios), seguimientos, tareas y timeline histórico.

Ruta: `/salvia/casos/:id/detalle`
Template: `src/frontend/html/salvia/case_detail/get_case_detail_sv.html`
Facade: `CaseDetailGET` (`CaseDetailFacade.go`)

## Archivos relevantes

| Archivo | Rol |
|---|---|
| `src/frontend/html/salvia/case_detail/get_case_detail_sv.html` | Template principal — lógica Vue y estructura completa (6 tabs) |
| `src/frontend/css/case_detail.css` | Estilos de layout, tabs, badges y estados visuales (`cd-*`, `ci-*`, `ns-*`) — cargado globalmente en `templates/layouts/header.html` |
| `src/frontend/js/components/case-timeline.js` | Timeline del caso (tab Timeline) |
| `src/frontend/js/components/case-info.js` | Info completa del caso — usado en modal y en modo inline (tab Info General) |
| `src/frontend/js/components/follow-up-contact-modal.js` | Modal de contacto para iniciar un seguimiento |
| `src/frontend/js/components/case-tasks.js` | Listado de tareas del caso — usado en el tab Tareas (todas) y, filtrado por barrera, dentro de cada `BarreraItem` en el tab Gestión institucional |
| `src/frontend/js/components/case-task-modal.js` | Modal para completar una `CaseTask` — montado a nivel de pantalla, ver nota en sección Modales |
| `src/frontend/js/components/case-task-history.js` | Modal de solo lectura para ver el detalle de una `CaseTask` completada — montado a nivel de pantalla, ver nota en sección Modales |
| `src/frontend/js/components/case-oficios.js` | Oficios del caso — embebido al final del tab Gestión institucional |
| `src/frontend/js/components/case-entities.js` | Entidades relacionadas con el caso — embebido al inicio del tab Gestión institucional, antes de la lista de barreras |
| `src/salvia/facades/CaseDetailFacade.go` | Facade que renderiza el template e inyecta `caseICode`, `userRole`, `userTeam`, `userICode`, `currentUser` |

---

## Árbol de interfaz

```
DetalleCaso  (.container)
│
├── Breadcrumb  (.cd-breadcrumb)
│   ├── Link "← Casos reportados"
│   └── [v-if caso] Span  caso.victimCaseVictimName
│
├── [v-if cargando]
│   └── LoadingText  "Cargando información del caso..."
│
├── [v-if error && !cargando]
│   └── ErrorPanel
│       ├── [v-if errorTipo === 404] Icon  "🔍"
│       ├── [v-else] Icon  "⚠️"
│       ├── [v-if errorTipo === 404] Title  "Caso no encontrado"
│       ├── [v-else] Title  "No se pudo cargar el caso"
│       ├── Description  (texto amigable según errorTipo)
│       ├── Link "← Volver al listado de casos"
│       ├── [v-if errorTipo !== 404] BtnReintentar "🔄 Reintentar"  → reintentar()
│       └── <details>  "Detalle técnico"  (colapsable, muestra `error`)
│
└── [v-if caso && !cargando]
    ├── Header  (.cd-header)
    │   ├── HeaderLeft  (.cd-header-left)
    │   │   ├── Avatar  (.cd-avatar)  primera letra de victimCaseVictimName (o "?")
    │   │   ├── Info
    │   │   │   ├── Name  (.cd-case-name)  caso.victimCaseVictimName
    │   │   │   └── Sub  (.cd-case-sub)  edad · municipio, departamento
    │   │   └── Badges  (.cd-badges)
    │   │       ├── BadgeRiesgo  :class=badgeRiesgo(nivelRiesgo)
    │   │       └── BadgeEstado  :class=badgeEstado(caso.victimCaseStatus)
    │   └── HeaderRight  (.cd-header-right)
    │       ├── BtnInfoCaso  "📋 Info caso"  → mostrarInfoCaso = true
    │       ├── [v-if puedeReasignar]  BtnReasignar  "Reasignar"  → abrirModalReasignar()
    │       ├── [v-if userRole === 'op' && (!caso.agentId || caso.agentId !== userICode)]
    │       │   BtnAsignarme  "👤 Asignarme este caso" | "Asignando..."  :disabled=autoAsignando  → autoAsignarCaso()
    │       └── [v-if puedeReasignar || userRole === 'ro' || (userRole === 'op' && caso.agentId === userICode)]
    │           BtnNuevoSeg
    │               :disabled si followUpsV2.length >= 8
    │               Label "Máximo alcanzado (8)" | "+ Nuevo seguimiento"
    │               → abrirModalNuevoSeg()
    │
    ├── [v-if proximoSeguimiento]  BannerProximoSeg
    │   └── "📅 Próximo seguimiento programado"  fecha · [hora si scheduled_time != '00:00'] · [👤 agente si agent_id]
    │       BtnVerSeg "Ver →"  → tabActiva = 'seguimientos'
    │
    ├── [v-if caseTasks.filter(status === 'ToDo').length > 0]  BannerTareasPendientesGlobal
    │   └── "Tienes N tarea(s) pendiente(s) con este caso"
    │       BtnVerTareas "Ver tareas"  → tabActiva = 'tareas'
    │
    └── Panel  (.cd-panel)
        ├── Tabs  (.cd-tabs)
        │   └── TabBtn × 6  [v-for tabs]  :class="active"
        │       ids: info | derivaciones | barreras | seguimientos | tareas | timeline
        │       labels: "Info General" | "Derivaciones" | "🏛️ Gestión institucional" | "Seguimientos" | "Tareas" | "Timeline"
        │       (id interno "barreras" sin cambiar — solo el label visible pasó de "Barreras" a "Gestión institucional")
        │
        └── TabContent  (.cd-tab-content)
            │
            ├── [tabActiva === 'info']  ── TAB INFO GENERAL ──────────────────
            │   │
            │   ├── SeccionResumen  (.ci-section)  [collapsible: infoSeccion.resumen — abierta por defecto]
            │   │   ├── MetricasRow (4): NivelRiesgo (badge) · EstadoCaso (badge) · DenunciasAnteriores (número) · AgenteResponsable (ultimoAgente)
            │   │   ├── Grid "Datos de la víctima"
            │   │   │   nombre completo · nombre identitario (nombreIdentitario) · tipo doc · número doc ·
            │   │   │   edad (edadCalculada o form1.age) · tiene hijos (form1.childrenNumber) ·
            │   │   │   municipio · departamento · tipo de agresor (tipoAgresorResumen) · fecha de reporte
            │   │   └── [v-if hay algún dato de contexto]  Grid "Contexto del caso"
            │   │       tipoViolencia · subtipoViolencia (lista si > 1) · ambitoViolencia ·
            │   │       nacionalidadResumen · generoResumen · territorioOcurrencia ·
            │   │       planAtencion (lista si > 1) · [v-if ajusteRazonable.length] ajusteRazonable
            │   │
            │   ├── [v-if tipoViolencia.length || subtipoViolencia.length || ambitoViolencia.length]
            │   │   SeccionHechos  (.ci-section)  [collapsible: infoSeccion.hechos — cerrada por defecto]
            │   │   ├── Grid "Hechos de violencia"
            │   │   │   tipoViolencia · subtipoViolencia · ambitoViolencia · territorioOcurrencia (lugar) · planAtencion
            │   │   └── [v-if hechosReportados.length]  SeccionHechosReportados
            │   │       "Hechos reportados en seguimientos"  (borde superior)
            │   │       └── HechoCard × N  [v-for hechosReportados]  (.bg-fef2f2)
            │   │           "Hecho #{idx+1}"  ·  fecha (h.date)  ·  h.description (o "Sin descripción")
            │   │
            │   └── SeccionInfoCompleta  (.ci-section)  [collapsible: infoSeccion.completa — cerrada por defecto]
            │       └── [v-if infoSeccion.completa]
            │           CaseInfo  // src/frontend/js/components/case-info.js
            │           :case-id="caseICode"  mode="inline"
            │
            ├── [tabActiva === 'derivaciones']  ── TAB DERIVACIONES ─────────────
            │   │
            │   ├── DerivSeccion "🔵 Medidas de Emergencia"  (fondo rojo claro)
            │   │   ├── Contador  emergencyMeasures.length + " registro(s)"
            │   │   ├── [v-if length === 0]  "Sin registros"
            │   │   └── Item × N  [v-for emergencyMeasures]
            │   │       type · [issuingAuthority] · fecha (createdAt) · [notes] · BadgeStatus (item.status || "ACTIVE", estilo fijo "proceso")
            │   │
            │   ├── DerivSeccion "🩷 Psicosocial"  (fondo azul claro)
            │   │   ├── Contador  psychosocialSupports.length + " registro(s)"
            │   │   ├── [v-if length === 0]  "Sin registros"
            │   │   └── Item × N  [v-for psychosocialSupports]
            │   │       type · [provider] · fecha · "Sesiones: N" (sessionCount) · [notes] · BadgeStatus (estilo fijo "completado")
            │   │
            │   └── DerivSeccion "📦 Estabilizacion Socioeconomica"  (fondo verde claro)
            │       ├── Contador  economicStabilizations.length + " registro(s)"
            │       ├── [v-if length === 0]  "Sin registros"
            │       └── Item × N  [v-for economicStabilizations]
            │           type · [institution] · [benefit] · fecha · [notes] · BadgeStatus (estilo fijo "articulando")
            │
            ├── [tabActiva === 'barreras']  ── TAB GESTIÓN INSTITUCIONAL (antes "Barreras") ──
            │   │
            │   ├── SeccionEntidades  (borde inferior, margen 24px)
            │   │   "🏛️ Entidades"
            │   │   CaseEntities  // src/frontend/js/components/case-entities.js
            │   │   :case-id="caseICode"  :user-id="userICode"  :user-role="userRole"
            │   │   (ver DocsMD/Componentes/case-entities/ para el árbol interno completo —
            │   │    filtros, cards, modal de agregar; lectura para cualquier rol con acceso a
            │   │    esta pantalla, agregar solo para sv/op/ro)
            │   │
            │   ├── [v-if tareasPendientesTotal().length > 0]  BannerTareasPendientesLocal
            │   │   └── "⚠️ Tienes N tarea(s) pendiente(s) con este caso"
            │   │       BtnVerTareas "Ver tareas"  → tabActiva = 'tareas'
            │   │       (misma condición y acción que BannerTareasPendientesGlobal, repetido dentro del tab)
            │   │
            │   ├── Stats (4): Total barreras · Abiertas (status === 'OPEN') · En gestión (barrerasEnGestionCount: status === 'En Gestion') · Resueltas (barreraEsResuelta: Articulada|MANAGED|CLOSED|FINISH|CLOSE|RESOLVED)
            │   │
            │   ├── [v-if barriers.length === 0]  EmptyState  "No se han identificado barreras en este caso"
            │   │
            │   ├── BarreraItem × N  [v-for barriers]  :class=barreraRowClass(status)
            │   │   ├── Indicador  (.cd-seg-indicator)  :class=barreraIndicatorClass(status)
            │   │   │   icono: ⚠ (OPEN) | ◷ (En Gestion) | ✓ (cualquier otro)
            │   │   ├── Sector  (b.sector || "Sin sector")  +  BadgeStatus  labelBarreraStatus(status)
            │   │   ├── [v-if b.specificBarriers]  "Barreras identificadas"
            │   │   │   └── Chip × N  [v-for resolveBarrierCSV(specificBarriers)]  (catálogo ~40 valores por sector salud/justicia/protección)
            │   │   ├── [v-if b.otherBarrierDesc]  "Otra barrera"  → texto libre
            │   │   ├── [v-if b.specificInstitutions]  "Instituciones"
            │   │   │   └── Chip × N  [v-for resolveInstitutionCSV(specificInstitutions)]  (catálogo de instituciones)
            │   │   ├── [v-if b.institutionName]  "Institución específica"  → texto
            │   │   ├── Fecha  📅 formatFecha(b.createdAt)
            │   │   ├── BtnToggleDetalle  "▼ Ver descripción y tareas" | "▲ Ocultar detalle"
            │   │   │   → barrExpandido = barrExpandido === b.id ? null : b.id
            │   │   ├── Link  "📄 Ver detalle barrera →"  → /salvia/barreras/:id
            │   │   └── [v-if barrExpandido === b.id]
            │   │       ├── [v-if b.description]  "Descripción"  → texto
            │   │       └── CaseTasks  // src/frontend/js/components/case-tasks.js
            │   │           :case-id="caseICode"  :barrier-id="b.id"  :user-id="userICode"  :user-role="userRole"
            │   │
            │   └── SeccionOficios  (borde superior, margen 24px)
            │       "📄 Oficios del caso"
            │       CaseOficios  // src/frontend/js/components/case-oficios.js
            │       :case-id="caseICode"
            │
            ├── [tabActiva === 'seguimientos']  ── TAB SEGUIMIENTOS ────────────
            │   │
            │   ├── Stats (3): Total (followUpsV2.length) · Ejecutados (status REALIZADO || is_completed) · Pendientes (status PENDIENTE && !is_completed)
            │   ├── [v-if followUpsV2.length === 0]  EmptyState  "No hay seguimientos registrados para este caso"
            │   └── SeguimientoItem × N  [v-for followUpsV2]  (fila completa clickeable)  @click → toggleSegExpand(seg.id)
            │       ├── Indicador  ✓ (si completado) | número de orden (idx+1)
            │       ├── "Seguimiento #{idx+1}"  +  BadgeStatus  seg.status
            │       ├── [v-if seg.attempts > 0]  BadgeIntentos  "{attempts}/9 intentos"
            │       │   color: rojo si ≥9 · naranja si ≥3 · gris si menos
            │       ├── Meta: 📅 fecha · [🕐 hora si != '00:00'] · [👤 agente si agent_id] · [Último intento si last_attempt_at]
            │       ├── [v-if seg.summary]  Resumen  (texto)
            │       ├── [si completado]  Link "Ver detalle →"  /salvia/seguimiento/:id  (@click.stop)
            │       ├── [si NO completado]  Actions  (@click.stop)
            │       │   ├── [v-if (rol ro/op) && seg.agent_id === userICode]  BtnIniciar "▶ Iniciar"  → iniciarSeguimiento(seg)
            │       │   ├── [v-else-if rol ro/op]  BtnIniciar  :disabled  title="Este seguimiento está asignado a otro agente"
            │       │   ├── [v-if puedeReasignar || (rol ro/op && caso.agentId === userICode)]
            │       │   │   BtnEditar "✏️ Editar"  → editarSeguimiento(seg)
            │       │   └── [v-if puedeReasignar || (rol ro/op && seg.agent_id === userICode)]
            │       │       BtnReasignarSeg "📋 Reasignar"  → reasignarSeguimiento(seg, idx)
            │       └── [v-if segExpandido === seg.id && seg.follow_up_attempts.length > 0]  (@click.stop)
            │           HistorialIntentos  "Historial de intentos (N)"
            │           └── IntentoItem × N  [v-for follow_up_attempts]
            │               número (ai+1) · "Contestó" (verde) | "No contestó" (rojo) según was_answered · [razón] · fechaHora
            │
            ├── [tabActiva === 'tareas']  ── TAB TAREAS ─────────────────────
            │   └── CaseTasks  // src/frontend/js/components/case-tasks.js
            │       :case-id="caseICode"  :user-id="userICode"  :user-role="userRole"
            │       (sin filtro de barrera — muestra todas las tareas del caso)
            │
            └── [tabActiva === 'timeline']  ── TAB TIMELINE ──────────────────
                └── CaseTimeline  // src/frontend/js/components/case-timeline.js
                    :case-id="caseICode"

─── MODALES ─────────────────────────────────────────────────────────────────────────────────

FollowUpContactModal  // src/frontend/js/components/follow-up-contact-modal.js
    ref="contactModal"  :current-user :current-user-id="userICode" :user-team="userTeam"
    @completed="onFollowUpContactCompleted"

CaseTaskModal  // src/frontend/js/components/case-task-modal.js
    ref="taskModal"  :current-user-id="userICode"  @completed="onCaseTaskCompleted"
    ⚠️ Montado a nivel de pantalla, pero su único disparador (`$refs.taskModal.open(testTaskId)`) vivía en un
    bloque de UI de desarrollo dentro del tab Derivaciones que hoy está **comentado** ("pruebas finalizadas").
    Actualmente no hay ningún control activo que abra este modal — `case-tasks.js` usa su **propia** instancia
    interna de `case-task-modal`, independiente de esta.

CaseTaskHistory  // src/frontend/js/components/case-task-history.js
    ref="taskHistory"   (modal de solo lectura)
    ⚠️ Misma situación que CaseTaskModal: su único trigger (`$refs.taskHistory.open(testHistoryTaskId)`) está en
    el bloque de UI de desarrollo comentado. No hay control activo que lo abra hoy.

CaseInfo (modo modal)  // src/frontend/js/components/case-info.js
    :case-id="caseICode"  :visible="mostrarInfoCaso"  @close="mostrarInfoCaso = false"

[v-if toast.show]  AlertModal  (.cd-alert-modal)
    [type === 'success']  icon ✓  "¡Operación exitosa!"   /   [type === 'error']  icon ✕  "Error"
    toast.message
    BtnAceptar  → toast.show = false

[v-if modalNuevoSeg]  ModalNuevoSeg  (.ns-modal)
    ├── [v-if nuevoSegError]  AlertaError  (mensaje + botón cerrar inline)
    ├── <input type="date">  Fecha*  :min=fechaHoy
    ├── <input type="time">  Hora
    ├── [v-if userRole === 'sv']  <select>  Agente*  [v-for agentesRO]
    ├── [v-else]  <input>  Agente  :disabled  (fijo al usuario actual)
    ├── <textarea>  Notas / Observaciones
    ├── [v-if !nuevoSegValido && segTocado]  Aviso  "Complete los campos obligatorios (fecha y agente) para guardar."
    └── Footer: BtnCancelar · BtnGuardar  :disabled si !nuevoSegValido || guardandoSeg

[v-if modalReasignar]  ModalReasignarCaso  (.ns-modal)
    ├── <select>  Operador*  [v-for operadores]
    └── Footer: BtnCancelar · BtnAsignar  :disabled si !reasignarOperador || reasignando

[v-if modalReasignarSeg]  ModalReasignarSeg  (.ns-modal)
    ├── Aviso de continuidad  "Se recomienda que el seguimiento continue con el agente que inicio la gestion..."
    ├── <select>  Agente*  [v-for agentesParaReasignarSeg]
    └── Footer: BtnCancelar · BtnConfirmar  :disabled si !reasignarSegAgente || reasignandoSeg

[v-if modalEditarSeg]  ModalEditarSeg  (.ns-modal)
    ├── <input type="date">  Fecha*  :min=fechaHoy
    ├── <input type="time">  Hora
    └── Footer: BtnCancelar · BtnGuardar  :disabled si !editarSegValido || guardandoEditSeg
```

---

## Niveles de riesgo

| nivelRiesgo | Label | Badge |
|---|---|---|
| 0 / null | Sin evaluar | `badge-cerrado` |
| 1 | Bajo | `badge-bajo` |
| 2 | Medio | `badge-medio` |
| 3 | Alto | `badge-alto` |
| 4 | Extremo | `badge-alto` (misma clase visual que Alto, label distinto) |

> `nivelRiesgo` prioriza `form2.riskLevel` si existe y es > 0; si no, cae a `form1.femicideRisk` (`'1'` → 3/Alto, `'0'` → 1/Bajo).

## Estados del caso

| code | Label | Badge |
|---|---|---|
| `ra` | Activo | `badge-activo` |
| `is` | Con novedad | `badge-seguimiento` |
| `cd` | Cerrado | `badge-cerrado` |
| `ex` | Vencido | `badge-vencido` |
| `r` | Por aprobar | `badge-cerrado` (sin estilo propio, cae al default) |
| `fc` | Recontacto | `badge-cerrado` (sin estilo propio, cae al default) |

## Estados de barrera (`barrier_v2.status`)

| status | Label (`labelBarreraStatus`) | Indicador | Grupo visual |
|---|---|---|---|
| `OPEN` | Abierta | ⚠ | pending |
| `En Gestion` | En gestión | ◷ | gestion |
| `Articulada` | Articulada | ✓ | done — cuenta como "Resuelta" |
| `MANAGED` | Gestionada | ✓ | done — cuenta como "Resuelta" |
| `CLOSED` | Cerrada | ✓ | done — cuenta como "Resuelta" |
| `FINISH` | Cerrada | ✓ | done — cuenta como "Resuelta" |
| `CLOSE` | Cerrada | ✓ | done — cuenta como "Resuelta" |
| `RESOLVED` | Resuelta | ✓ | done — cuenta como "Resuelta" |

## Controles por rol

| Control | Condición |
|---|---|
| Btn "Reasignar" (caso) | `sv` |
| Btn "👤 Asignarme este caso" | `op`, solo si el caso no tiene agente o el agente no es el usuario actual |
| Btn "+ Nuevo seguimiento" | `sv`, `ro`, o `op` dueño del caso (`caso.agentId === userICode`) |
| Btn "▶ Iniciar" (habilitado) | `op`, `ro` — solo si `seg.agent_id === userICode` |
| Btn "✏️ Editar" seguimiento | `sv`, o (`ro`/`op`) dueño del caso (`caso.agentId === userICode`) |
| Btn "📋 Reasignar" seguimiento | `sv`, o (`ro`/`op`) agente propio del seguimiento (`seg.agent_id === userICode`) |

## Resolución de agentes para reasignar un seguimiento (`agentesParaReasignarSeg`)

| Rol actual | Lista mostrada |
|---|---|
| `sv` | `agentesRO` (ya filtrado por team del caso) |
| `op` | Solo él mismo — array de un elemento `[{ icode: userICode, fullName: currentUser }]` |
| `ro` (u otro) | `todosAgentes` filtrado por `team === caso.victimCaseTeam` (o `userTeam` si el caso no tiene team) |
