# Notificaciones — Interfaz

Pantalla que gestiona el flujo completo de oficios enviados a entidades para atender barreras institucionales de casos VBG. El operador (`op`) proyecta y corrige oficios; el agente de notificaciones (`an`) los revisa, aprueba, radica y registra respuestas. Cada oficio sigue un flujo de estados: `por_proyectar → para_revisar → aprobacion_juridica → para_radicar → radicado → respondido`, con una rama de corrección que devuelve al estado `para_revisar`.

La pantalla carga todos los oficios asignados al usuario según su rol, los muestra en una tabla paginada (5 por página) y permite gestionarlos mediante modales específicos por estado. Los filtros y tabs operan de forma reactiva sin recargar la página.

---

## Archivos relevantes

| Archivo | Descripción |
|---|---|
| `src/frontend/html/salvia/notifications/notificaciones.html` | Template principal: encabezado, tabs, filtros, tabla y paginador |
| `src/frontend/html/salvia/notifications/modal_proyectar.html` | Modal para estado `por_proyectar` |
| `src/frontend/html/salvia/notifications/modal_revisar.html` | Modal para estado `para_revisar` |
| `src/frontend/html/salvia/notifications/modal_corregir.html` | Modal para estado `en_correccion` |
| `src/frontend/html/salvia/notifications/modal_aprobar.html` | Modal para estado `aprobacion_juridica` |
| `src/frontend/html/salvia/notifications/modal_radicar.html` | Modal para estado `para_radicar` |
| `src/frontend/html/salvia/notifications/modal_registrar_respuesta.html` | Modal para estado `radicado` |
| `src/frontend/js/components/notifications.js` | Lógica Vue: carga, filtros, paginación, modales y acciones |

---

## Árbol de interfaz

```
Pantalla: Notificaciones
│
├── Encabezado
│   ├── Título: "Notificaciones Salvia"
│   └── Subtítulo: "Gestión de oficios asignados a tu rol"
│
├── Tabs
│   ├── Tab "Todos mis oficios"
│   └── Tab "Oficios por gestionar"
│       └── Badge con conteo de oficios donde canManage = true
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
│   └── Input texto "Número radicado"
│
├── [v-if: !isLoading && !loadError] Tabla de oficios
│   │
│   ├── Encabezado de tabla: Caso | Barrera | Información del Oficio
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
│   │   │   └── Link "Ver barrera →" (pendiente de implementar)
│   │   │
│   │   └── Columna INFORMACIÓN DEL OFICIO
│   │       ├── Tag de estado del oficio
│   │       ├── [v-if: canManage] Botón "Gestionar"
│   │       ├── Tag de prioridad del oficio (normal | alta)
│   │       ├── [v-if: urlKofax] Link al documento en Kofax
│   │       ├── Número radicado · Fecha de creación del oficio
│   │       └── [v-if: correoEntidad] Correo de la entidad
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

├── Modal "Proyectar oficio" [estado: por_proyectar | rol: op]
│   ├── Título del oficio (solo lectura)
│   ├── Hint de instrucciones
│   ├── Select "Nivel" (municipal | departamental | nacional) *requerido
│   ├── Input "Entidad" *requerido
│   ├── Input "Ruta del oficio en el Kofax" *requerido
│   ├── Select "Prioridad del oficio" (normal | alta) *requerido
│   └── Footer: [Cancelar] [Registrar]
│
├── Modal "Revisar oficio" [estado: para_revisar | rol: an]
│   ├── Título del oficio (solo lectura)
│   ├── Link al Kofax (solo lectura)
│   ├── Entidad y Nivel (solo lectura)
│   ├── Select "Razón de corrección" (requerido si elige "Por corregir")
│   └── Footer: [Cancelar] [Por corregir] [Marcar como revisado]
│
├── Modal "Corregir oficio" [estado: en_correccion | rol: op]
│   ├── Título del oficio (solo lectura)
│   ├── Hint de instrucciones
│   ├── [v-if: reasonCorrection] Razón de corrección del revisor (solo lectura, fondo rojo)
│   ├── Link al Kofax (solo lectura)
│   └── Footer: [Cancelar] [Marcar oficio como corregido]
│
├── Modal "Radicar oficio — Aprobación jurídica" [estado: aprobacion_juridica | rol: an]
│   ├── Título del oficio (solo lectura)
│   ├── Link al Kofax (solo lectura)
│   ├── Entidad y Nivel (solo lectura)
│   ├── Input "Asunto" *requerido
│   ├── Input "Correo entidad" *requerido
│   ├── Input "Número radicado" *requerido
│   └── Footer: [Cancelar] [Marcar oficio como radicado]
│
├── Modal "Radicar oficio" [estado: para_radicar | rol: an]
│   ├── Título del oficio (solo lectura)
│   ├── Input "Asunto" *requerido
│   ├── Input "Correo entidad" *requerido
│   ├── Input "Número radicado" *requerido
│   ├── Input "Ruta del oficio en el Kofax" *requerido
│   └── Footer: [Cancelar] [Marcar oficio como radicado]
│
└── Modal "Registrar respuesta" [estado: radicado | rol: an]
    ├── Título del oficio (solo lectura)
    ├── Hint de instrucciones
    ├── Input date "Fecha de respuesta" *requerido
    ├── Input email "Correo del remitente" *requerido
    ├── Input "Asunto de la respuesta" *requerido
    ├── Input "Respuesta recibida por" *requerido
    └── Footer: [Cancelar] [Registrar respuesta]
```
