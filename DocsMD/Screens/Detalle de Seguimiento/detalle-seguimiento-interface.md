# `get_follow_up_detail` — Interfaz de la Pantalla

Vista de solo lectura (o edición limitada) de un seguimiento completado. Muestra dos tabs: **Resumen** con la información de la víctima, datos del seguimiento, barreras identificadas y remisiones; y **Formulario** con el formulario dinámico asociado al seguimiento.

Ruta: `/salvia/seguimiento/:id`
Template: `src/frontend/html/salvia/follow_up_detail/get_follow_up_detail.html`
Facade: `FollowUpGET`

---

## Archivos relevantes

| Archivo | Rol |
|---|---|
| `src/frontend/html/salvia/follow_up_detail/get_follow_up_detail.html` | Template principal — estructura, lógica Vue y carga de datos |
| `src/frontend/js/components/dinamic-form.js` | Componente de formulario dinámico embebido en el tab Formulario |
| `src/frontend/css/styles.css` | Estilos de layout, tabs, cards y estados visuales (`ms-*`, `vc-*`) |
| `src/salvia/facades/FollowUpFacade.go` | Facade que renderiza el template e inyecta `followUpId` |

---

## Árbol de interfaz

```
follow-up-detail-app  (.ms-page-wrapper)
│
├── Breadcrumb  (.ms-breadcrumb)
│   ├── BtnVolver  (.ms-breadcrumb-btn)  → goBack()
│   └── Texto  "{ victimCase.name } / Seguimiento #{ followUp.id }"
│
└── TabsWrapper  (.ms-tabs-wrapper)
    │
    ├── TabsHeader  (.ms-tabs-header)
    │   └── TabBtn × N  [v-for TABS]  (.ms-tab-button)  @click → activeTab = tab.id
    │       Estados: active | inactive
    │       Tabs disponibles: "Resumen" (fas fa-chart-bar) | "Formulario" (fas fa-clipboard-list)
    │
    └── TabContent  (.ms-tab-content)
        │
        ├── [v-if activeTab === 'resumen']  ── TAB RESUMEN ──────────────────
        │   │
        │   ├── [v-if victimInfo]  SeccionVictima  (.ms-detail-section)
        │   │   └── CardVictima  (.vc-card)
        │   │       ├── Header  (.vc-header)
        │   │       │   ├── Avatar  (.vc-avatar)  <i fas fa-user>
        │   │       │   ├── Identity  (.vc-identity)
        │   │       │   │   ├── Nombre  (.vc-name)   victimInfo.Names + victimInfo.LastNames
        │   │       │   │   └── Ubicacion  (.vc-location)  victimInfo.TownName
        │   │       │   └── Meta  (.vc-meta)
        │   │       │       ├── Agente  followUp.agent
        │   │       │       └── FechaRealizacion  followUp.completedAt
        │   │       ├── Divider  (.vc-divider)
        │   │       └── InfoGrid  (.ms-info-grid)
        │   │           ├── Telefono          victimInfo.Phone
        │   │           ├── IdentidadGenero   victimInfo.GenderIdentity
        │   │           ├── OrientacionSexual victimInfo.SexualOrientation
        │   │           ├── TelefonoContacto  victimInfo.ContactPhone
        │   │           └── Edad              victimInfo.Age + " años"
        │   │
        │   ├── SeccionInfoSeguimiento  (.ms-detail-section)
        │   │   └── InfoGrid  (.ms-info-grid)
        │   │       ├── FechaRealizacion  followUp.completedAt
        │   │       └── Agente            followUp.agent
        │   │
        │   └── SeccionPlanAtencion  (.ms-detail-section)
        │       └── CarePlan  (.ms-care-plan)
        │           │
        │           ├── [v-if activeBarriers.length > 0]  ListaBarreras  (.ms-barriers-list)
        │           │   ├── BarrerasHeader  (.ms-barriers-header)
        │           │   │   ├── Icono + Titulo  "Barreras identificadas"
        │           │   │   ├── TooltipInfo  (title: descripcion de barrera)
        │           │   │   └── Contador  (.ms-barriers-count)  activeBarriers.length
        │           │   └── BarreraItem × N  [v-for activeBarriers]  (.ms-barrier-item)
        │           │       └── Info  (.ms-barrier-info)
        │           │           ├── Sector       (.ms-barrier-sector)  barrera.sector
        │           │           └── Descripcion  (.ms-barrier-description)  barrera.description
        │           │
        │           ├── [v-else]  EmptyStateBarreras  (.ms-empty-state-small)
        │           │   └── Texto  "No se identificaron barreras institucionales en este seguimiento."
        │           │
        │           ├── [v-if emergencyMeasures.length > 0 || psychosocialSupports.length > 0 || economicStabilizations.length > 0]
        │           │   Remisiones  (.ms-referrals)
        │           │   ├── RemisionesHeader  (.ms-referrals-header)  "Remisiones a equipos Salvia"
        │           │   └── RemisionesList  (.ms-referrals-list)
        │           │       ├── MedidaEmergencia × N  [v-for emergencyMeasures]  (.ms-referral-item)
        │           │       │   ├── Badge  (.ms-badge--emergency)  "Medida: { em.type }"
        │           │       │   ├── TooltipInfo  "Autoridad: { em.issuingAuthority }"
        │           │       │   └── Notas  em.notes
        │           │       │
        │           │       ├── ApoyoPsicosocial × N  [v-for psychosocialSupports]  (.ms-referral-item)
        │           │       │   ├── Badge  (.ms-badge--psychosocial)  "Psicosocial: { ps.type }"
        │           │       │   ├── TooltipInfo  (descripcion fija de evaluacion psicosocial)
        │           │       │   └── Notas  ps.notes
        │           │       │
        │           │       └── EstabilizacionEconomica × N  [v-for economicStabilizations]  (.ms-referral-item)
        │           │           ├── Badge  (.ms-badge--economic)  "Económica: { es.type }"
        │           │           ├── TooltipInfo  "Institución: { es.institution }"
        │           │           └── Beneficio  es.benefit
        │           │
        │           └── [v-else]  EmptyStateRemisiones  (.ms-empty-state-small)
        │               └── Texto  "No se registraron remisiones a equipos Salvia en este seguimiento."
        │
        └── [v-if activeTab === 'formulario']  ── TAB FORMULARIO ──────────
            │
            ├── [v-if formId && submissionId]
            │   dinamic-form  // src/frontend/js/components/dinamic-form.js
            │       :form-id="formId"
            │       :submission-id="submissionId"
            │       :can-edit="canEdit"
            │       :form-state="formState"
            │       @answers-updated → onAnswersUpdated()
            │
            └── [v-else]  EmptyStateFormulario  (.ms-empty-state-small)
                └── Texto  "Este seguimiento no tiene un formulario asociado."
```

---

## Estados de `canEdit`

| Condición | canEdit |
|---|---|
| Caso cerrado (`caseInfo.status === 'cd'`) | `false` permanente |
| Seguimiento REALIZADO y han pasado > 5 días desde `completed_at` | `false` |
| Cualquier otro estado | `true` |

---

## `formState` — estado externo inyectado al formulario dinámico

| Campo | Valores posibles | Cómo se calcula |
|---|---|---|
| `psysocialRemisionState` | `''` \| `'Remisión: SI cumple'` \| `'Remisión: NO cumple'` | `onAnswersUpdated()`: evalúa la pregunta `71c42c4a-f640-47ad-b2c1-5d4c18480449`; requiere `criterio_obligatorio` + puntaje ≥ 3 con los demás criterios |
