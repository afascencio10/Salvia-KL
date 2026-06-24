# `Detalle del Caso` — Interfaz de la Pantalla

Vista completa de un caso VBG: datos de la víctima, hechos, seguimientos programados, barreras, derivaciones y timeline histórico.

## Archivos relevantes

| Archivo | Rol |
|---|---|
| `src/frontend/html/salvia/case_detail/get_case_detail_sv.html` | Template principal — lógica Vue y estructura completa |
| `src/frontend/js/components/case-timeline.js` | Componente timeline del caso (tab Timeline) |
| `src/frontend/js/components/case-info.js` | Info completa del caso — usado en modal y modo inline |
| `src/frontend/js/components/follow-up-contact-modal.js` | Modal de contacto para iniciar un seguimiento |

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
│   └── LoadingText "Cargando información del caso..."
│
├── [v-if error && !cargando]
│   └── ErrorPanel
│       ├── [v-if errorTipo === 404] Icon "🔍"
│       ├── [v-else] Icon "⚠️"
│       ├── [v-if errorTipo === 404] Title "Caso no encontrado"
│       ├── [v-else] Title "No se pudo cargar el caso"
│       ├── Description amigable
│       ├── Link "← Volver al listado de casos"
│       ├── [v-if errorTipo !== 404] Btn "🔄 Reintentar"  → reintentar()
│       └── <details> detalle técnico colapsable
│
└── [v-if caso && !cargando]
    ├── Header  (.cd-header)
    │   ├── HeaderLeft  (.cd-header-left)
    │   │   ├── Avatar  (.cd-avatar)  inicial del nombre de la víctima
    │   │   ├── Info
    │   │   │   ├── Name  (.cd-case-name)  caso.victimCaseVictimName
    │   │   │   └── Sub  (.cd-case-sub)  edad · municipio, departamento
    │   │   └── Badges
    │   │       ├── BadgeRiesgo  :class=badgeRiesgo(nivelRiesgo)
    │   │       └── BadgeEstado  :class=badgeEstado(caso.victimCaseStatus)
    │   └── HeaderRight  (.cd-header-right)
    │       ├── BtnInfoCaso "📋 Info caso"  → mostrarInfoCaso = true
    │       ├── [v-if puedeReasignar] BtnReasignar "Reasignar"  → abrirModalReasignar()
    │       └── [v-if puedeReasignar || userRole === 'ro'] BtnNuevoSeg
    │               :disabled si followUpsV2.length >= 8
    │               Label "Máximo alcanzado (8)" | "+ Nuevo seguimiento"
    │               → abrirModalNuevoSeg()
    │
    ├── [v-if proximoSeguimiento] BannerProximoSeg
    │   └── fecha, hora, agente  ·  BtnVerSeg "Ver →"  → tabActiva = 'seguimientos'
    │
    └── Panel  (.cd-panel)
        ├── Tabs  (.cd-tabs)
        │   └── TabBtn × 5  [v-for tabs]
        │       ids: info | derivaciones | barreras | seguimientos | timeline
        │
        └── TabContent  (.cd-tab-content)
            │
            ├── [tabActiva === 'info']
            │   ├── SeccionResumen  (.ci-section)  [collapsible: infoSeccion.resumen — abierta por defecto]
            │   │   ├── MetricasRow: nivelRiesgo · estadoCaso · denunciasAnteriores · ultimoAgente
            │   │   ├── Grid datos víctima (nombre, doc, edad, hijos, municipio, dpto, agresor, fecha reporte)
            │   │   └── [v-if contexto existe] Grid contexto del caso
            │   │       (tipoViolencia, subtipo, ámbito, nacionalidad, género, territorio, plan, ajuste)
            │   │
            │   ├── [v-if tipoViolencia.length || subtipoViolencia.length || ambitoViolencia.length]
            │   │   SeccionHechos  (.ci-section)  [collapsible: infoSeccion.hechos — cerrada por defecto]
            │   │   ├── Grid hechos de violencia
            │   │   └── [v-if hechosReportados.length] HechosTimeline × N  (hechos de seguimientos)
            │   │
            │   └── SeccionInfoCompleta  (.ci-section)  [collapsible: infoSeccion.completa — cerrada por defecto]
            │       └── [v-if infoSeccion.completa]
            │           CaseInfo  // src/frontend/js/components/case-info.js
            │           :case-id="caseICode"  mode="inline"
            │
            ├── [tabActiva === 'derivaciones']
            │   ├── DerivSeccion "Medidas de Emergencia"  (fondo rojo claro)
            │   │   └── Item × N  [v-for emergencyMeasures]  type · issuingAuthority · fecha · status · notes
            │   ├── DerivSeccion "Psicosocial"  (fondo azul claro)
            │   │   └── Item × N  [v-for psychosocialSupports]  type · provider · sessionCount · fecha · status · notes
            │   └── DerivSeccion "Estabilización Socioeconómica"  (fondo verde claro)
            │       └── Item × N  [v-for economicStabilizations]  type · institution · benefit · fecha · status · notes
            │
            ├── [tabActiva === 'barreras']
            │   ├── Stats: total · abiertas (status OPEN) · resueltas (!OPEN)
            │   ├── [v-if barriers.length === 0] EmptyState "No se han identificado barreras"
            │   └── BarreraItem × N  [v-for barriers]
            │       indicator (⚠ abierta | ✓ resuelta) + sector · descripción · fecha · badgeStatus
            │
            ├── [tabActiva === 'seguimientos']
            │   ├── Stats: total · ejecutados (REALIZADO o is_completed) · pendientes
            │   ├── [v-if followUpsV2.length === 0] EmptyState
            │   └── SeguimientoItem × N  [v-for followUpsV2]  @click=toggleSegExpand
            │       ├── Indicator (✓ | número) + badge status + [v-if attempts > 0] badge intentos N/9
            │       ├── Meta: fecha · hora · agente · último intento
            │       ├── [v-if seg.summary] Resumen
            │       ├── [si REALIZADO] Link "Ver detalle →"  /salvia/seguimiento/:id
            │       ├── [si PENDIENTE] Actions
            │       │   ├── [agente propio && rol op/ro] BtnIniciar "▶ Iniciar"  → iniciarSeguimiento(seg)
            │       │   ├── [otro agente && rol op/ro] BtnIniciar :disabled
            │       │   ├── [v-if puedeReasignar || ro-dueño] BtnEditar "✏️ Editar"  → editarSeguimiento(seg)
            │       │   └── [v-if puedeReasignar || agente propio] BtnReasignarSeg "📋 Reasignar"  → reasignarSeguimiento(seg, idx)
            │       └── [v-if segExpandido === seg.id && tiene intentos]
            │           HistorialIntentos × N  (número · contestó/no contestó · razón · fechaHora)
            │
            └── [tabActiva === 'timeline']
                └── CaseTimeline  // src/frontend/js/components/case-timeline.js
                    :case-id="caseICode"

─── MODALES ─────────────────────────────────────────────────────────────────────────────────

FollowUpContactModal  // src/frontend/js/components/follow-up-contact-modal.js
    ref="contactModal"  @completed="onFollowUpContactCompleted"

CaseInfo (modo modal)  // src/frontend/js/components/case-info.js
    :visible="mostrarInfoCaso"  @close="mostrarInfoCaso = false"

[v-if toast.show] AlertModal  (.cd-alert-modal)
    [type === 'success'] icon ✓  /  [type === 'error'] icon ✕
    BtnAceptar  → toast.show = false

[v-if modalNuevoSeg] ModalNuevoSeg  (.ns-modal)
    ├── <input type="date"> Fecha*  :min=fechaHoy
    ├── <input type="time"> Hora
    ├── [userRole === 'sv'] <select> Agente*  [v-for agentesRO]
    ├── [userRole !== 'sv'] <input> Agente  :disabled  (fijo al usuario actual)
    ├── <textarea> Notas / Observaciones
    └── Footer: BtnCancelar · BtnGuardar :disabled si !nuevoSegValido || guardandoSeg

[v-if modalReasignar] ModalReasignarCaso  (.ns-modal)
    ├── <select> Operador*  [v-for operadores]
    └── Footer: BtnCancelar · BtnAsignar :disabled si !reasignarOperador || reasignando

[v-if modalReasignarSeg] ModalReasignarSeg  (.ns-modal)
    ├── Aviso de continuidad
    ├── <select> Agente*  [v-for agentesParaReasignarSeg]
    └── Footer: BtnCancelar · BtnConfirmar :disabled si !reasignarSegAgente || reasignandoSeg

[v-if modalEditarSeg] ModalEditarSeg  (.ns-modal)
    ├── <input type="date"> Fecha*  :min=fechaHoy
    ├── <input type="time"> Hora
    └── Footer: BtnCancelar · BtnGuardar :disabled si !editarSegValido || guardandoEditSeg
```

---

## Niveles de riesgo

| nivelRiesgo | Label | Badge |
|---|---|---|
| 0 / null | Sin evaluar | `badge-cerrado` |
| 1 | Bajo | `badge-bajo` |
| 2 | Medio | `badge-medio` |
| 3 o 4 | Alto | `badge-alto` |

## Estados del caso

| code | Label | Badge |
|---|---|---|
| `ra` | Activo | `badge-activo` |
| `is` | Con novedad | `badge-seguimiento` |
| `cd` | Cerrado | `badge-cerrado` |
| `ex` | Vencido | `badge-vencido` |

## Controles por rol

| Control | Roles visibles |
|---|---|
| Btn "Reasignar" (caso) | `sv` |
| Btn "+ Nuevo seguimiento" | `sv`, `ro` |
| Btn "▶ Iniciar" (habilitado) | `op`, `ro` — solo si `seg.agent_id === userICode` |
| Btn "✏️ Editar" seguimiento | `sv`, o `ro` dueño del caso |
| Btn "📋 Reasignar" seguimiento | `sv`, o agente propio del seguimiento |
