# Planeación — Pantalla de Notificaciones

> ⚠️ **MD temporal de validación.** No implementar hasta que sea aprobado.
> Archivo: `src/frontend/html/salvia/notifications/notificaciones.html`
> JS: `src/frontend/js/components/notifications.js`

---

## 1. Descripción de la interfaz actual

La pantalla gestiona el flujo completo de **oficios enviados a entidades** (`entity_letter`) para atender barreras institucionales de casos VBG.

**Roles que acceden:**
- `op` (Operador/Agente Seguimiento) — proyecta y corrige oficios
- `an` (Agente de Notificaciones) — revisa, aprueba, radica y registra respuestas

**Flujo de estados del oficio:**
```
por_proyectar → para_revisar → aprobacion_juridica → para_radicar → radicado → respondido
                     ↓
               en_correccion → [vuelve a para_revisar]
```

---

## 2. Árbol de interfaz

```
Pantalla: Notificaciones
│
├── Encabezado
│   ├── Título: "Notificaciones Salvia"
│   └── Subtítulo: "Gestión de oficios asignados a tu rol"
│
├── Tabs
│   ├── Tab "Todos mis oficios"
│   └── Tab "Oficios por gestionar" [badge con conteo de pendientes]
│
├── [Si isLoading] Estado de carga — spinner + texto
│
├── [Si loadError] Estado de error — ícono + mensaje
│
├── [Si cargado sin error] Filtros
│   ├── Select: Estado del oficio (por_proyectar | para_revisar | en_correccion | aprobacion_juridica | para_radicar | radicado | respondido)
│   ├── Input texto: Número de identidad
│   ├── Input texto: Entidad
│   └── Input texto: Número radicado
│
├── [Si cargado sin error] Tabla de oficios (paginada, 5 por página)
│   └── Fila por oficio
│       ├── Columna CASO
│       │   ├── Nombre completo de la víctima
│       │   ├── Código del caso + número documento
│       │   ├── Tag de prioridad del caso
│       │   └── Link "Ver caso →" (va a /salvia/casos/:id/detalle)
│       │
│       ├── Columna BARRERA
│       │   ├── Tags de sectores de la barrera
│       │   ├── Organización de la barrera
│       │   ├── Descripción de la barrera
│       │   └── Link "Ver barrera →" (sin implementar aún — href="#")
│       │
│       └── Columna INFORMACIÓN DEL OFICIO
│           ├── Tag de estado del oficio
│           ├── [Si canManage] Botón "Gestionar"
│           ├── Tag de prioridad del oficio (normal | alta)
│           ├── [Si urlKofax] Link al documento en Kofax
│           ├── Número radicado + Fecha de creación
│           └── [Si correoEntidad] Correo de la entidad
│
├── [Si sin resultados] Estado vacío — ícono + texto
│
└── [Si numPages > 1] Paginador
    ├── Botón "Anterior"
    ├── Botones numerados (con "..." para rangos largos)
    └── Botón "Siguiente"

Modales (uno activo a la vez según estado del oficio seleccionado):
│
├── Modal "Proyectar oficio" [estado: por_proyectar | rol: op]
│   ├── Nivel (select: municipal | departamental | nacional) *requerido
│   ├── Entidad (texto) *requerido
│   ├── Ruta del oficio en el Kofax (texto) *requerido
│   └── Prioridad del oficio (select: normal | alta) *requerido
│
├── Modal "Revisar oficio" [estado: para_revisar | rol: an]
│   ├── Link al Kofax (solo lectura)
│   ├── Entidad y Nivel (solo lectura)
│   ├── Razón de corrección (select — requerido si elige "Por corregir")
│   └── Acciones: [Cancelar] [Por corregir] [Marcar como revisado]
│
├── Modal "Corregir oficio" [estado: en_correccion | rol: op]
│   ├── Razón de corrección registrada por el revisor (solo lectura, fondo rojo)
│   ├── Link al Kofax (solo lectura)
│   └── Acción: [Marcar oficio como corregido]
│
├── Modal "Radicar oficio (Aprobación jurídica)" [estado: aprobacion_juridica | rol: an]
│   ├── Link al Kofax (solo lectura)
│   ├── Entidad y Nivel (solo lectura)
│   ├── Asunto *requerido
│   ├── Correo entidad *requerido
│   ├── Número radicado *requerido
│   └── Acción: [Marcar oficio como radicado]
│
├── Modal "Radicar oficio" [estado: para_radicar | rol: an]
│   ├── Asunto *requerido
│   ├── Correo entidad *requerido
│   ├── Número radicado *requerido
│   ├── Ruta del oficio en el Kofax *requerido
│   └── Acción: [Marcar oficio como radicado]
│
└── Modal "Registrar respuesta" [estado: radicado | rol: an]
    ├── Fecha de respuesta *requerido
    ├── Correo del remitente *requerido
    ├── Asunto de la respuesta *requerido
    ├── Respuesta recibida por *requerido
    └── Acción: [Registrar respuesta]
```

---

## 3. Inventario de eventos

### E01 — Cuando carga la pantalla
**Tipo:** Lifecycle
**Descripción:** Se ejecuta al montar el componente Vue. Verifica el rol del usuario, construye la URL del API con el filtro correspondiente y carga los oficios.
**Requerido:** Sí

### E02 — Cuando cambia de tab
**Tipo:** User Interaction
**Descripción:** El usuario presiona "Todos mis oficios" o "Oficios por gestionar". Cambia `currentTab` y resetea la página a 0. Vue filtra reactivamente `filteredOficios`.
**Requerido:** Sí

### E03 — Cuando aplica filtros
**Tipo:** User Interaction
**Descripción:** El usuario escribe o selecciona en cualquiera de los 4 filtros (estado, identidad, entidad, radicado). Vue filtra reactivamente y resetea la página a 0.
**Requerido:** Sí

### E04 — Cuando cambia de página
**Tipo:** User Interaction
**Descripción:** El usuario presiona Anterior, Siguiente o un número de página. Cambia `currentPage` y hace scroll al top.
**Requerido:** Sí

### E05 — Cuando presiona "Gestionar"
**Tipo:** User Interaction
**Descripción:** El usuario abre el modal correspondiente al estado actual del oficio. Pre-carga el formulario con los datos existentes del oficio seleccionado.
**Requerido:** Sí

### E06 — Cuando cancela el modal
**Tipo:** User Interaction
**Descripción:** El usuario presiona "Cancelar", la X del modal o el overlay. Cierra el modal y limpia el estado del formulario.
**Requerido:** Sí

### E07 — Cuando confirma acción en modal
**Tipo:** User Interaction
**Descripción:** El usuario presiona el botón de acción principal del modal (proyectar, revisar, corregir, radicar, registrar respuesta). Valida los campos requeridos, envía el payload al endpoint `PATCH /api/v1/entity-letters/:id/action` y actualiza el estado del oficio en la lista sin recargar la pantalla.
**Requerido:** Sí

---

## 4. Flujos lógicos

---

### Flujo E01 — Cuando carga la pantalla

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando carga la pantalla
   Tipo: Lifecycle
   Función: mounted() → loadOficios()
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  currentRole:   rol del usuario    → window.NotifConfig.currentRole (inyectado por Go)
  currentUserId: ID del usuario     → window.NotifConfig.currentUserId
}

PASO 1 — Verificar que el rol es válido (op | an)

SI rol NO está en ['op', 'an']:
  → Asignar loadError = 'Rol no autorizado para acceder a esta pantalla.'
  → TERMINAR ejecución

SI rol es válido:
  → CONTINÚA FLUJO GENERAL

PASO 2 — Activar estado de carga (isLoading = true, loadError = null)

PASO 3 — Construir URL del API según rol

SI currentRole === 'op':
  → url = '/api/v1/entity-letters?limit=100&page=0&agentId={currentUserId}'

SI currentRole === 'an':
  → url = '/api/v1/entity-letters?limit=100&page=0&notificationUserId={currentUserId}'

PASO 4 — Llamar al API

GET {url}

→ resultado: lista de entity_letter o error

PASO 5 — Manejar respuesta

SI status === 200:
  → Extraer array de items (response o response.items)
  → Ordenar por createdAt DESC
  → Mapear cada item a objeto oficio con mapApiToOficio()
    (incluye calcular canManage según estado y rol)
  → Asignar a this.oficios
  → isLoading = false

SI status === 401:
  → Redirigir a /static/landing.html

SI otro error:
  → loadError = 'No se pudo cargar la lista de oficios. Intenta de nuevo.'
  → isLoading = false

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión                              | Paso afectado |
|--------------------------------------------------|---------------|
| ¿Qué campos exactos devuelve el endpoint GET?    | PASO 4        |
| ¿El endpoint pagina? ¿limit=100 es suficiente?   | PASO 3        |
```

---

### Flujo E05 — Cuando presiona "Gestionar"

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando presiona "Gestionar"
   Tipo: User Interaction
   Función: openGestionarModal(oficio)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  oficio:  objeto del oficio seleccionado  → fila clickeada en la tabla
}

PASO 1 — Guardar el oficio seleccionado en selectedOficio

PASO 2 — Determinar qué modal abrir según el estado del oficio

SEGÚN oficio.status:
  CASO 'por_proyectar':       → activeModal = 'proyectar'
  CASO 'para_revisar':        → activeModal = 'revisar'
  CASO 'en_correccion':       → activeModal = 'corregir'
  CASO 'aprobacion_juridica': → activeModal = 'aprobar'
  CASO 'para_radicar':        → activeModal = 'radicar'
  CASO 'radicado':            → activeModal = 'registrar_respuesta'
  DEFAULT:                    → activeModal = null (modal no se muestra)

PASO 3 — Pre-cargar el formulario con datos existentes del oficio
  modalForm.nivel                = oficio.nivel || ''
  modalForm.kofaxPath            = oficio.urlKofax || ''
  modalForm.prioridad            = oficio.letterPriority || 'normal'
  modalForm.entidad              = oficio.entidad || oficio.barrierOrg || ''
  modalForm.asunto               = oficio.asuntoRadicado || ''
  modalForm.correoEntidad        = oficio.correoEntidad || ''
  modalForm.numeroRadicado       = oficio.numeroRadicado || ''
  modalForm.fechaRespuesta       = ''
  modalForm.correoRemitente      = oficio.correo || ''
  modalForm.asuntoRespuesta      = ''
  modalForm.respuestaRecibidaPor = ''
  modalForm.reasonCorrection     = ''

PASO 4 — Limpiar errores y mostrar el modal
  saveError = null
  showModal = true
  → Vue muestra el modal correspondiente a activeModal
```

---

### Flujo E07 — Cuando confirma acción en modal

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando confirma acción en modal
   Tipo: User Interaction
   Función: submitModal(action)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  action:         string con la acción a ejecutar   → botón presionado en el modal
  selectedOficio: oficio actualmente seleccionado   → state del componente
  modalForm:      campos del formulario del modal   → inputs del usuario
  currentUserId:  ID del usuario en sesión          → window.NotifConfig
}

PASO 1 — Construir payload base
  payload = { action: action, userId: currentUserId }

┌─────────────────────────────────────────────────────────┐
│  SUB-FLUJO: Validación y payload por tipo de acción     │
└─────────────────────────────────────────────────────────┘

  SI action === 'proyectar':
    Validar:
      • nivel:      required
      • entidad:    required
      • kofaxPath:  required
    SI falla validación:
      → alert() con el campo faltante
      → TERMINAR ejecución
    SI ok:
      → payload.nivel    = modalForm.nivel
      → payload.entidad  = modalForm.entidad
      → payload.urlKofax = modalForm.kofaxPath
      → payload.priority = modalForm.prioridad || 'normal'

  SI action === 'por_corregir':
    Validar:
      • reasonCorrection: required
    SI falla:
      → alert()
      → TERMINAR ejecución
    SI ok:
      → payload.reasonCorrection = modalForm.reasonCorrection

  SI action === 'radicar':
    Validar:
      • asunto:         required
      • correoEntidad:  required
      • numeroRadicado: required
    SI falla:
      → alert()
      → TERMINAR ejecución
    SI ok:
      → payload.asuntoRadicado = modalForm.asunto
      → payload.correoEntidad  = modalForm.correoEntidad
      → payload.numeroRadicado = modalForm.numeroRadicado

  SI action === 'registrar_respuesta':
    Validar:
      • fechaRespuesta:       required
      • correoRemitente:      required
      • asuntoRespuesta:      required
      • respuestaRecibidaPor: required
    SI falla:
      → alert()
      → TERMINAR ejecución
    SI ok:
      → payload.responseDate     = modalForm.fechaRespuesta
      → payload.correoRemitente  = modalForm.correoRemitente
      → payload.asuntoRespuesta  = modalForm.asuntoRespuesta
      → payload.responseReviewBy = modalForm.respuestaRecibidaPor

  SI action === 'revisar' o 'corregir':
    → No requiere campos adicionales
    → payload solo contiene { action, userId }

  → FIN SUB-FLUJO → CONTINÚA FLUJO GENERAL

PASO 2 — Activar estado de guardado (isSaving = true, saveError = null)

PASO 3 — Enviar al API

PATCH /api/v1/entity-letters/{selectedOficio.id}/action

payload: {
  action:  acción ejecutada           → determinada en Sub-flujo
  userId:  currentUserId              → sesión
  ...campos adicionales según acción
}

PASO 4 — Manejar respuesta

SI status === 200:
  → Ocultar overlay de éxito después de 1.2s
  → Encontrar el oficio en this.oficios por ID
  → Actualizar en la lista:
      status    = response.state
      canManage = canManageForRole(response.state)
  → Cerrar modal (closeModal)
  → isSaving = false

SI status === 401:
  → Redirigir a /static/landing.html

SI status === 422:
  → saveError = response.error || 'Transición de estado no permitida.'
  → isSaving = false

SI otro error:
  → saveError = response.error || 'No se pudo guardar el oficio. Intenta de nuevo.'
  → isSaving = false

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión                                          | Paso afectado |
|--------------------------------------------------------------|---------------|
| ¿El modal 'aprobar' llama a action='radicar'? (mismo modal)  | PASO 1        |
| ¿Qué campos devuelve el PATCH en la respuesta?               | PASO 4        |
| Las validaciones usan alert() — ¿se cambiará a inline?       | Sub-flujo     |
```

---

## 5. Observaciones y mejoras sugeridas

| # | Observación | Impacto |
|---|---|---|
| 1 | El link "Ver barrera →" tiene `href="#"` — no lleva a ningún lado | Medio |
| 2 | Las validaciones del modal usan `alert()` — no es consistente con el resto del sistema | Bajo |
| 3 | `mapApiToOficio` asigna `priority: 'Alto'` hardcodeado en lugar de usar el valor real del API | Medio |
| 4 | El modal `aprobar` y el modal `radicar` ejecutan la misma `action: 'radicar'` — pueden unificarse o aclararse | Bajo |
| 5 | No hay flujo documentado para E02, E03, E04 — son triviales (solo cambian estado Vue reactivo) | Sin impacto |
