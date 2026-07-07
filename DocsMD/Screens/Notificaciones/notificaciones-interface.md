# Notificaciones — Interfaz

Pantalla que gestiona el flujo completo de oficios enviados a entidades para atender barreras institucionales de casos VBG. El operador (`op`) y revisor operativo (`ro`) proyectan y corrigen oficios; el agente de notificaciones (`an`) los revisa, radica y registra respuestas. Cada oficio sigue un flujo de estados: `por_proyectar → para_revisar → aprobacion_juridica → para_radicar → radicado → respondido`, con una rama de corrección que devuelve al estado `para_revisar`.

La pantalla carga oficios paginados desde el backend (5 por página), aplica filtros y tabs directamente en la base de datos, y permite gestionarlos mediante modales específicos por estado.

> **v2 planificado (rol `an`):** ver [notificaciones-plan-v2-an.md](./notificaciones-plan-v2-an.md)

---

## Archivos relevantes

| Archivo | Descripción |
|---|---|
| `src/frontend/html/salvia/notifications/notificaciones.html` | Template principal: encabezado, tabs, filtros, tabla y paginador |
| `src/frontend/html/salvia/notifications/modal_proyectar.html` | Modal para estado `por_proyectar` |
| `src/frontend/html/salvia/notifications/modal_revisar.html` | Modal para estado `para_revisar` |
| `src/frontend/html/salvia/notifications/modal_corregir.html` | Modal para estado `en_correccion` |
| `src/frontend/html/salvia/notifications/modal_aprobar.html` | Modal activo para radicación (`aprobacion_juridica`; acción `radicar`) |
| `src/frontend/html/salvia/notifications/modal_radicar.html` | **Obsoleto — no se usa** (legacy; casi idéntico a `modal_aprobar`) |
| `src/frontend/html/salvia/notifications/modal_registrar_respuesta.html` | Modal para estado `radicado` |
| `src/frontend/js/components/notifications.js` | Lógica Vue: carga paginada, filtros server-side, modales y acciones |
| `src/salvia/controller/entity_letter_controller.go` | `GET /api/v1/entity-letters` — listado paginado con filtros |
| `src/salvia/controller/entity_branch_api_controller.go` | `GET /api/v1/entity-branches` — devuelve `{id, icode, name}` por `town_code` |
| `src/internal/repository/entity_letter_repository.go` | Consultas SQL paginadas con JOINs y filtros ILIKE |
| `DocsMD/Screens/Notificaciones/notificaciones-api.md` | Documentación del endpoint de listado paginado |
| `DocsMD/Screens/Notificaciones/notificaciones-plan-v2-an.md` | Plan v2: visibilidad global `an`, 3 campos auditoría, tabs, filtros e historial |

---

## Árbol de interfaz

```
Pantalla: Notificaciones
│
├── Encabezado
│   ├── Título: "Notificaciones Salvia"
│   └── Subtítulo: "Gestión de oficios asignados a tu rol"
│
├── Tabs [rol op / ro — v1 actual]
│   ├── Tab "Todos mis oficios"
│   └── Tab "Oficios por gestionar"
│       └── Badge con pendingCount del API
│
├── Tabs [rol an — v2 planificado]
│   ├── Tab "Todos"
│   ├── Tab "Mis Oficios"
│   └── Tab "Oficios por gestionar"
│       └── Badge: total gestionables del sistema
│
├── [v-if: isLoading] Estado de carga
│   ├── Spinner
│   └── Texto "Cargando oficios..."
│
├── [v-if: loadError && !isLoading] Estado de error
│   ├── Ícono de advertencia
│   └── Mensaje de error
│
├── [v-if: !isLoading && !loadError] Filtros
│   ├── Select "Estado del oficio"
│   │   └── Opciones: Todos | Por proyectar | Para revisar | En corrección |
│   │       Aprobación jurídica | Para radicar | Radicado | Respondido
│   ├── Input texto "Número de identidad"
│   ├── Input texto "Entidad"
│   ├── Input texto "Número radicado"
│   └── [v2, solo rol an] 3 dropdowns (mismo catálogo de agentes an activos):
│       ├── Select "Revisado por"
│       ├── Select "Radicado por"
│       └── Select "Respuesta registrada por"
│
├── [v-if: !isLoading && !loadError] Tabla de oficios
│   │
│   ├── Encabezado de tabla [v2]: Caso | Barrera | Información del Oficio | Historial Agentes
│   ├── Encabezado de tabla [v1 actual]: Caso | Barrera | Información del Oficio
│   │
│   ├── [v-for: paginatedOficios] Fila por oficio
│   │   │
│   │   ├── Columna CASO
│   │   │   ├── Nombre completo de la víctima
│   │   │   ├── Código del caso · Número de documento
│   │   │   ├── Tag de prioridad del caso
│   │   │   └── Link "Ver caso →" → /salvia/casos/:caseId/detalle
│   │   │
│   │   ├── Columna BARRERA
│   │   │   ├── Tags de sectores de la barrera
│   │   │   ├── Organización de la barrera
│   │   │   ├── Descripción de la barrera
│   │   │   └── Link "Ver barrera →"
│   │   │
│   │   ├── Columna INFORMACIÓN DEL OFICIO
│   │   │   ├── Tag de estado del oficio
│   │   │   ├── [v-if: canManage] Botón "Gestionar"
│   │   │   ├── Tag de prioridad del oficio (normal | alta)
│   │   │   ├── [v-if: entidad] Nombre de entidad · Municipio
│   │   │   ├── [v-if: officialDependency] Funcionario
│   │   │   ├── [v-if: subject] Asunto
│   │   │   ├── [v-if: urlKofax] Link al documento en Kofax
│   │   │   ├── Número radicado · Fecha de creación del oficio
│   │   │   └── [v-if: correoEntidad] Correo de la entidad
│   │   │
│   │   └── Columna HISTORIAL AGENTES [v2 planificado]
│   │       ├── [v-if: notificationUserIdReview]
│   │       │   Revisado por: {nombre agente}
│   │       ├── [v-if: notificationUserIdRadicado]
│   │       │   Radicado por: {nombre agente}
│   │       ├── [v-if: notificationUserIdResponse]
│   │       │   Respuesta por: {nombre agente}
│   │       └── [v-if: ninguno] "—"
│   │
│   └── [v-else] Estado vacío
│       ├── Ícono de bandeja
│       └── Texto "No se encontraron oficios con los filtros seleccionados."
│
└── [v-if: numPages > 1] Paginador
    ├── Botón "Anterior" (deshabilitado en página 0)
    ├── Botones numéricos con "..." para rangos largos
    └── Botón "Siguiente" (deshabilitado en última página)

─────────────────────────────────────────────────
Modales (renderizados siempre, visibles según activeModal)
─────────────────────────────────────────────────

├── Modal "Proyectar oficio" [estado: por_proyectar | rol: op, ro]
│   └── Footer: [Cancelar] [Registrar]
│       → Sin cambio en campos notificationUserId*
│
├── Modal "Revisar oficio" [estado: para_revisar | rol: an]
│   └── Footer: [Cancelar] [Por corregir] [Marcar como revisado]
│       → por_corregir  → notificationUserIdReview = userId (reemplaza)
│       → revisar       → notificationUserIdReview = userId (reemplaza)
│       → review_by     → se mantiene lógica actual en revisar
│
├── Modal "Corregir oficio" [estado: en_correccion | rol: op, ro]
│   └── Sin cambio en campos notificationUserId*
│
├── Modal "Radicar oficio — Aprobación jurídica" [estado: aprobacion_juridica | rol: an]
│   │   Modal único de radicación en uso (reemplaza modal_radicar.html)
│   └── Footer: [Cancelar] [Por corregir] [Marcar oficio como radicado]
│       → radicar → notificationUserIdRadicado = userId (reemplaza)
│       → radicado_by → se mantiene lógica actual
│
├── modal_radicar.html [OBSOLETO — no referenciar en v2]
│   └── Archivo legacy; flujo real usa modal_aprobar.html
│
└── Modal "Registrar respuesta" [estado: radicado | rol: an]
    └── Footer: [Cancelar] [Registrar respuesta]
        → registrar_respuesta → notificationUserIdResponse = userId (reemplaza)
```

---

## Columna "Historial Agentes" — especificación UI

| Aspecto | Detalle |
|---|---|
| Visible para | Todos los roles (`op`, `ro`, `an`) |
| Fuente de datos | `notificationUserIdReview`, `notificationUserIdRadicado`, `notificationUserIdResponse` del API |
| Resolución de nombre | Catálogo de agentes `an` cargado en frontend, o campos `*Name` enriquecidos desde backend |
| Líneas mostradas | Solo las acciones que tengan agente registrado |
| Campos null (históricos) | No mostrar línea; no error; columna `"—"` si todo vacío |
| Estilo sugerido | Texto compacto, labels en gris, nombre del agente en negrita |
