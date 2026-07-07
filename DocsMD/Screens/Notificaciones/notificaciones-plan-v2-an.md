# Notificaciones — Plan v2 (Agente de Notificaciones)

> **Estado:** Implementado (jul 2026)  
> **Origen:** Cambios acordados en reunión (jul 2026)  
> **Última actualización:** modal único de radicación + compatibilidad datos históricos (null-safe)

---

## Resumen del cambio de negocio

| Antes (v1) | Después (v2) |
|---|---|
| Cada agente `an` solo ve oficios donde `notification_user_id = su id` | Todos los agentes `an` ven **todos** los `entity_letter` |
| Solo puede gestionar oficios asignados a él | Cualquier agente `an` puede gestionar **cualquier** oficio en estado gestionable |
| Tabs `an`: "Todos mis oficios" / "Oficios por gestionar" | Tabs `an`: **"Todos"** / **"Mis Oficios"** / **"Oficios por gestionar"** |
| Sin historial visible de agentes por acción | Nueva columna **"Historial Agentes"** en la tabla |
| 4 filtros de texto/select | + **3 dropdowns** de agente (Revisado por, Radicado por, Respuesta registrada por) |

Los roles `op` y `ro` **no cambian** en alcance ni tabs (siguen filtrando por `agent_id`).

---

## Campos nuevos en `EntityLetter` (3 campos)

### Convención de nombres

| Capa | Patrón | Ejemplo |
|---|---|---|
| Go struct | `NotificationUserID{Accion}` | `NotificationUserIDReview` |
| Columna SQL | `notification_user_id_{accion}` | `notification_user_id_review` |
| JSON API | `notificationUserId{Accion}` | `notificationUserIdReview` |

Todos son `*string`, `varchar(36)`, nullable — almacenan el `userICode` del agente `an` que ejecutó la acción.

> **`NotificationUserIDJuridica` no se crea** — omitido por decisión de negocio.

### Tabla de campos

| Campo Go | Columna SQL | JSON | Acción(es) modal | Modal | Campo se actualiza |
|---|---|---|---|---|---|
| `NotificationUserIDReview` | `notification_user_id_review` | `notificationUserIdReview` | `revisar` **o** `por_corregir` | `modal_revisar.html` | Siempre **reemplaza** el valor anterior |
| `NotificationUserIDRadicado` | `notification_user_id_radicado` | `notificationUserIdRadicado` | `radicar` | `modal_aprobar.html` ("Marcar oficio como radicado") | Asigna / reemplaza |
| `NotificationUserIDResponse` | `notification_user_id_response` | `notificationUserIdResponse` | `registrar_respuesta` | `modal_registrar_respuesta.html` ("Registrar respuesta") | Asigna / reemplaza |

> **`modal_radicar.html` no se usa** en el flujo actual. Es casi idéntico a `modal_aprobar.html` pero queda **sin referencia activa**. La radicación se gestiona únicamente desde `modal_aprobar.html`. En implementación v2: `modalByStatus.para_radicar` debe apuntar a `'aprobar'` (o eliminar el template `modal_radicar` del árbol de modales).

### Modales activos vs. obsoletos

| Template | Estado | ¿En uso? |
|---|---|---|
| `modal_revisar.html` | `para_revisar` | Sí |
| `modal_aprobar.html` | `aprobacion_juridica` (y `para_radicar` si aplica) | **Sí — único modal de radicación** |
| `modal_radicar.html` | `para_radicar` (referencia legacy en JS) | **No — obsoleto, no implementar en v2** |
| `modal_registrar_respuesta.html` | `radicado` | Sí |

### Reglas de persistencia (E07)

```
acción revisar       (modal_revisar) → notification_user_id_review   = userId  [reemplaza si existe]
acción por_corregir  (modal_revisar) → notification_user_id_review   = userId  [reemplaza si existe]
acción radicar       (modal_aprobar)  → notification_user_id_radicado = userId  [reemplaza si existe]
acción registrar_respuesta           → notification_user_id_response  = userId  [reemplaza si existe]
```

El frontend **no envía** estos campos; el backend los asigna desde `userId` del payload.

> **Nota:** `por_corregir` desde `modal_aprobar.html` **no** actualiza `notification_user_id_review` — solo las acciones del modal revisar.

### Fragmento propuesto para el modelo Go

```go
// ── Auditoría: agente de notificaciones por acción ───────────────────────────
NotificationUserIDReview   *string `gorm:"type:varchar(36);column:notification_user_id_review"   json:"notificationUserIdReview,omitempty"`
NotificationUserIDRadicado *string `gorm:"type:varchar(36);column:notification_user_id_radicado" json:"notificationUserIdRadicado,omitempty"`
NotificationUserIDResponse *string `gorm:"type:varchar(36);column:notification_user_id_response" json:"notificationUserIdResponse,omitempty"`
```

### Migración SQL propuesta

```sql
ALTER TABLE salvia.entity_letter
    ADD COLUMN IF NOT EXISTS notification_user_id_review   VARCHAR(36),
    ADD COLUMN IF NOT EXISTS notification_user_id_radicado VARCHAR(36),
    ADD COLUMN IF NOT EXISTS notification_user_id_response VARCHAR(36);

CREATE INDEX IF NOT EXISTS idx_el_notif_review   ON salvia.entity_letter (notification_user_id_review);
CREATE INDEX IF NOT EXISTS idx_el_notif_radicado ON salvia.entity_letter (notification_user_id_radicado);
CREATE INDEX IF NOT EXISTS idx_el_notif_response ON salvia.entity_letter (notification_user_id_response);
```

### Relación con campos existentes

| Campo existente | Decisión |
|---|---|
| `notification_user_id` | **Se mantiene** en el modelo. Ya no es relevante para visibilidad, tabs ni filtros v2. |
| `review_by` | **Lógica actual intacta** — se sigue llenando en `revisar` con el mismo `userId`. Independiente del nuevo campo. |
| `radicado_by` | **Lógica actual intacta** — se sigue llenando en `radicar` con el mismo `userId`. Independiente del nuevo campo. |
| `register_by` | Sin cambio — acción `proyectar` (`op`) |
| `response_review_by` | Sin cambio — texto libre del modal registrar respuesta |

Los campos nuevos y los legacy pueden almacenar el **mismo ID** sin conflicto; cumplen propósitos distintos (auditoría de notificaciones vs. auditoría legacy).

### Datos históricos — sin backfill

- Los registros existentes **no se migran** a los 3 campos nuevos.
- Las columnas nuevas se crean **nullable**; oficios antiguos tendrán `NULL` en los 3 campos.
- **No debe generarse error** en backend ni frontend cuando los campos son `null`:
  - API: omitir en JSON con `omitempty` o devolver `null` sin fallar.
  - Columna Historial Agentes: ocultar líneas vacías; si los 3 son `null` → mostrar `"—"`.
  - Tab "Mis Oficios": oficios antiguos **no aparecen** ahí hasta que el agente ejecute una acción v2.
  - Filtros dropdown: sin valor seleccionado no filtran; no error si el oficio no tiene agente en ese campo.
  - Resolución de nombre: si el `userICode` no está en catálogo, mostrar el ID o `"—"` sin romper la fila.

---

## Catálogo de agentes para dropdowns

- **Fuente:** usuarios **activos** con rol **`an`**.
- **Uso:** los **3 dropdowns** comparten el **mismo listado**.
- **Valor del option:** `userICode`.
- **Label del option:** nombre completo del agente (ej. `names + lastNames`).
- **Opción por defecto:** "Todos" (valor vacío — sin filtro).

### Endpoint propuesto

```
GET /api/v1/users?role=an&active=true
```

> Si ya existe un endpoint equivalente en el proyecto, reutilizarlo. Cargar una sola vez al montar la pantalla (rol `an`).

---

## Tabs por rol (UI)

### Rol `op` / `ro` — sin cambio

| Tab | Query API |
|---|---|
| Todos mis oficios | `agentId={userId}` |
| Oficios por gestionar | `agentId={userId}&manageableOnly=true` |

### Rol `an` — nuevo

| Tab | Valor `currentTab` | Query API |
|---|---|---|
| **Todos** | `todos` | Sin filtro de usuario — todos los oficios paginados |
| **Mis Oficios** | `mis_oficios` | `mineOnly=true&notificationAgentId={userId}` |
| **Oficios por gestionar** | `gestionar` | `manageableOnly=true` — todos los gestionables del sistema |

#### Tab "Mis Oficios" — criterio en BD

Oficio incluido si el usuario logueado aparece en **cualquiera** de los 3 campos nuevos:

```sql
(
  el.notification_user_id_review   = :userId OR
  el.notification_user_id_radicado = :userId OR
  el.notification_user_id_response = :userId
)
```

#### Tab "Oficios por gestionar"

Sin filtro por agente. Estados gestionables de `an`:

`para_revisar`, `aprobacion_juridica`, `para_radicar`, `radicado`

`canManage` para rol `an` depende **solo del estado**, no de asignación.

#### Badge `pendingCount`

Para `an`: contar **todos** los oficios gestionables del sistema.

---

## Filtros (dropdowns) — solo rol `an`

| Label UI | Query param API | Filtra columna |
|---|---|---|
| Revisado por | `notificationUserIdReview` | `notification_user_id_review` |
| Radicado por | `notificationUserIdRadicado` | `notification_user_id_radicado` |
| Respuesta registrada por | `notificationUserIdResponse` | `notification_user_id_response` |

Los 3 dropdowns muestran el **mismo catálogo** de agentes `an` activos.

Los filtros existentes (estado, identidad, entidad, radicado) se mantienen y se combinan con AND.

---

## Nueva columna: "Historial Agentes"

### Ubicación

Cuarta columna en la tabla, visible para **todos los roles** que acceden a la pantalla (`op`, `ro`, `an`).

```
Encabezado: Caso | Barrera | Información del Oficio | Historial Agentes
```

### Contenido por fila

Mostrar solo las líneas cuyo campo tenga valor. Resolver `userICode` → nombre legible usando el catálogo de agentes `an` (o nombres enriquecidos desde API si el backend los incluye en el JOIN).

```
Revisado por:       {nombre agente}     ← notificationUserIdReview
Radicado por:       {nombre agente}     ← notificationUserIdRadicado
Respuesta por:      {nombre agente}     ← notificationUserIdResponse
```

Si un campo es `null`, no se muestra esa línea. Si los tres son `null`, mostrar `"—"` **sin error ni advertencia** (comportamiento esperado en oficios históricos).

### Datos en API

Incluir en `EntityLetterWithRelations`:

```json
{
  "notificationUserIdReview": "USR001",
  "notificationUserIdRadicado": "USR002",
  "notificationUserIdResponse": null,
  "notificationUserReviewName": "María García",
  "notificationUserRadicadoName": "Pedro López",
  "notificationUserResponseName": null
}
```

> Los campos `*Name` son opcionales en v2. Alternativa: resolver nombres en frontend con el catálogo ya cargado para los dropdowns.

---

## Cambios en API (`GET /api/v1/entity-letters`)

### Parámetros nuevos

| Parámetro | Tipo | Rol | Descripción |
|---|---|---|---|
| `mineOnly` | bool | `an` | Tab "Mis Oficios" |
| `notificationAgentId` | string | `an` | ID del agente logueado (pareja de `mineOnly`) |
| `notificationUserIdReview` | string | `an` | Filtro dropdown |
| `notificationUserIdRadicado` | string | `an` | Filtro dropdown |
| `notificationUserIdResponse` | string | `an` | Filtro dropdown |

### Parámetros que dejan de usarse para `an`

| Parámetro | Motivo |
|---|---|
| `notificationUserId` como filtro principal de listado | Reemplazado por tabs `todos` / `mineOnly` |

---

## Cambios en interfaz (HTML + JS)

| Elemento | Cambio |
|---|---|
| Tabs (`an`) | 3 tabs: Todos, Mis Oficios, Oficios por gestionar |
| Filtros (`an`) | + 3 dropdowns con catálogo compartido de agentes `an` |
| Tabla header | + columna "Historial Agentes" |
| Tabla fila | + bloque con las 3 líneas de historial |
| `notificaciones.html` | Ajustar grid CSS de 3 → 4 columnas |
| `notifications.js` | `buildOficiosUrl`, tabs `an`, carga catálogo agentes, render historial |
| `notificaciones.css` | Estilos columna historial |

---

## Archivos a modificar (implementación futura)

| Capa | Archivo | Cambio |
|---|---|---|
| Modelo | `entity_letter.go`, `entity_letter_view.go` | 3 campos nuevos (+ nombres opcionales en DTO) |
| BD | migración SQL | ALTER TABLE + 3 índices |
| Repository | `entity_letter_repository.go` | Filtros `mineOnly`, 3 dropdowns, listado global `an` |
| Service | `entity_letter_service.go` | `PerformAction` persiste 3 campos; mantiene `review_by`/`radicado_by` |
| Controller | `entity_letter_controller.go` | Nuevos query params |
| Users API | endpoint o reutilizar existente | Listar agentes `an` activos |
| Frontend HTML | `notificaciones.html` | 3 tabs `an`, 3 dropdowns, columna historial |
| Frontend JS | `notifications.js` | Tabs, filtros, historial null-safe, catálogo agentes; `modalByStatus.para_radicar` → `'aprobar'` |
| Frontend CSS | `notificaciones.css` | Grid 4 columnas |
| Docs | flujos E01–E07, API, interfaz | Alinear con v2 |

---

## Decisiones cerradas

| Tema | Decisión |
|---|---|
| `modal_radicar.html` | No se usa; radicación solo vía `modal_aprobar.html` |
| Backfill datos antiguos | **No.** Campos nuevos quedan `NULL`; UI/backend null-safe |
| Columna historial | Visible para todos los roles; sin error si campos vacíos |

---

## Orden de implementación sugerido

1. Migración BD + modelo Go + DTO `EntityLetterWithRelations`
2. Endpoint catálogo agentes `an` activos
3. `PerformAction` — persistir 3 campos nuevos (sin tocar lógica `review_by`/`radicado_by`)
4. Repository — filtros tabs, 3 dropdowns, listado global `an`
5. Controller — query params
6. Frontend — tabs `an`, dropdowns, columna Historial Agentes, CSS grid
7. Actualizar documentación de flujos E01–E07 y API
8. Pruebas manuales por rol
