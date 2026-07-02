# `barrier-follow-up-timeline` — Interfaz del Componente

Componente de solo lectura que recibe el `barrierId` de una barrera y renderiza en orden cronológico los registros de seguimiento (`barrier_follow_up`) como un listado de cards.

## Archivos relevantes

| Archivo | Rol |
|---|---|
| `src/frontend/js/components/barrier-follow-up-timeline.js` | Componente principal — template, estilos y lógica Vue |

---

## Árbol de interfaz

```
barrier-follow-up-timeline  (.bft-container)
│
├── [v-if cargando]
│   └── LoadingMsg  (.bft-loading)  "Cargando seguimientos..."
│
├── [v-if error && !cargando]
│   └── ErrorMsg  (.bft-error)  error
│
└── [v-if !cargando && !error]
    │
    ├── [v-if seguimientos.length === 0]
    │   └── EmptyState  (.bft-empty)
    │       "No hay seguimientos registrados para esta barrera."
    │
    └── [v-else]
        CardList  (.bft-list)
        └── FollowUpCard × N  (.bft-card)  [v-for seguimientos]  — orden: más antiguo primero
            │
            ├── CardHeader  (.bft-card-header)
            │   ├── Fecha  (.bft-fecha)  seguimiento.createdAt  (dd mmm. yyyy)
            │   └── Autor  (.bft-autor)  seguimiento.actorName
            │
            ├── TagRow  (.bft-tag-row)
            │   SEGÚN seguimiento.closesBarrier:
            │     true  → Tag  "Cierra barrera"   (.bft-tag .bft-tag--green)
            │     false, persists=true  → Tag  "Persiste"   (.bft-tag .bft-tag--orange)
            │     false, persists=false → Tag  "No persiste" (.bft-tag .bft-tag--gray)
            │
            ├── FieldRespuesta  (.bft-field)
            │   ├── Label  "Respuesta institucional"
            │   └── Value  (.bft-value)  seguimiento.institutionalResponseLabel
            │
            ├── [v-if seguimiento.managementActions]
            │   FieldGestion  (.bft-field)
            │   ├── Label  "Gestión realizada"
            │   └── ChipList  (.bft-chips)
            │       Chip × N  [v-for resolveGestionCSV(seguimiento.managementActions)]
            │
            ├── [v-if seguimiento.actions]
            │   FieldActuaciones  (.bft-field)
            │   ├── Label  "Actuaciones"
            │   └── Text  (.bft-text)  seguimiento.actions
            │
            └── [v-if seguimiento.closesBarrier && seguimiento.closureReason]
                FieldMotivoCierre  (.bft-field)
                ├── Label  "Motivo del cierre"
                └── Value  (.bft-value)  seguimiento.closureReasonLabel
```

---

## Mapas de etiquetas

### `institutionalResponse` → etiqueta

| Valor | Etiqueta |
|---|---|
| `respuesta_oficial` | Respondió de manera oficial a Salvia |
| `contacto_victima` | Se contactó con la víctima |
| `sin_respuesta` | No se obtuvo respuesta |

### `managementActions` (CSV) → chips

| Valor | Etiqueta |
|---|---|
| `activacion_ruta_interinstitucional` | Activación de ruta interinstitucional |
| `alerta_barreras` | Alerta por barreras |
| `articulacion_institucional` | Articulación institucional |
| `escalamiento_organismo_control` | Escalamiento a organismo de control |
| `gestion_llamada` | Gestión administrativa - Llamada |
| `orientacion_llamada` | Orientación y enrutamiento - Llamada |

### `closureReason` → etiqueta

| Valor | Etiqueta |
|---|---|
| `resuelta` | Resuelta de manera efectiva |
| `superada_parcialmente` | Superada parcialmente |
| `instalada_entidad` | Instalada en entidad competente |
| `no_gestionable` | No gestionable desde la competencia institucional |
| `no_voluntad` | Expresa no voluntad de accionar institucional |
| `clasificacion_incorrecta` | Clasificada incorrectamente |
