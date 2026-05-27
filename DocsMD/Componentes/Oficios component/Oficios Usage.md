# Componente `OficiosList` — Guía de uso

Componente Vue 3 de **solo lectura** que lista los oficios (`EntityLetter`) filtrados por caso, barrera o usuario. Muestra una tabla con el estado del oficio, la entidad y el documento, y abre un modal con el detalle completo incluyendo las secciones de Barrera y Caso cuando están disponibles.

---

## Archivos del componente

| Archivo | Ruta | Descripción |
|---|---|---|
| Template HTML | `frontend/html/salvia/notifications/oficios_list.html` | Markup Vue + punto de montaje |
| Lógica JS | `frontend/js/components/oficios-list.js` | App Vue, llamadas a la API, mapeo de datos |
| Estilos CSS | `frontend/css/oficios-list.css` | Estilos con prefijo `ol-` |

---

## Props del componente

Todos los parámetros son **opcionales**. Si no se pasa ninguno, carga la lista paginada general.

| Prop | Tipo | Descripción |
|---|---|---|
| `caseId` | `string \| null` | Filtra oficios del caso con ese UUID |
| `barrierId` | `string \| null` | Filtra oficios de la barrera con ese UUID |
| `userId` | `string \| null` | Filtra oficios donde el usuario es agente **o** notificador (dos llamadas paralelas + merge) |

### Comportamiento de la API según el prop activo

| Prop activo | Endpoint | Respuesta |
|---|---|---|
| `caseId` | `GET /api/v1/entity-letters?caseId={id}` | `[]EntityLetter` — sin datos de relaciones |
| `barrierId` | `GET /api/v1/entity-letters?barrierId={id}` | `[]EntityLetter` — sin datos de relaciones |
| `userId` | `GET /api/v1/entity-letters?agentId={id}` **+** `GET /api/v1/entity-letters?notificationUserId={id}` | `[]EntityLetterWithRelations` — incluye datos de caso y barrera |
| ninguno | `GET /api/v1/entity-letters?limit=100&page=0` | `[]EntityLetter` paginado general |

> **Nota sobre relaciones:** cuando se filtra por `userId`, el modal muestra las cards completas de **Caso** (nombre de víctima, código, documento) y **Barrera** (sectores, organización, descripción). Con `caseId` o `barrierId`, el modal solo muestra los IDs con sus respectivos enlaces de navegación.

---

## Prioridad de parámetros

El componente resuelve los parámetros en este orden:

```
1. data-* attributes del div  →  uso server-side (Go template)
2. window.OficiosListConfig   →  uso desde JS del padre (Angular-like)
3. null                       →  sin filtro, lista general
```

---

## Métodos de integración

Existen tres formas de pasar parámetros al componente según el contexto de la pantalla.

---

### Método 1 — Go template (recomendado para pantallas de detalle)

Ideal cuando la pantalla ya tiene el ID disponible como variable Go (ej: detalle de caso, detalle de barrera).

#### Paso 1 — Registrar el partial en el facade

Agregar `oficios_list.html` al slice de `extraTemplates` del facade de la pantalla:

```go
// src/salvia/facades/MiPantallaFacade.go

var miPantallaTemplates = []string{
    // ... otros partials ...
    "frontend/html/salvia/notifications/oficios_list.html",
}
```

#### Paso 2 — Pasar las variables en `RenderTemplate`

```go
common_facades.RenderTemplate(
    c,
    "mi-entidad",
    "salvia",
    "mi-modulo/",
    salvia_config.HTML_Templates,
    "mi-pantalla",
    extraTemplates,
    utils.DEFAULT_VIEW,
    utils.DEFAULT_PANIC_TEMPLATE,
    map[string]interface{}{
        "windowTitle": "Mi Pantalla",
        // ... otros campos ...

        // Props para OficiosList — pasar solo el que aplique:
        "OficiosListCaseId":    caseId,    // filtra por caso
        // "OficiosListBarrierId": barrierId, // filtra por barrera
        // "OficiosListUserId":    userId,    // filtra por agente/notificador
    },
    utils.GetFullHtmlFuncMap(),
)
```

#### Paso 3 — Incluir el componente en el HTML de la pantalla

```html
{{/* mi_pantalla.html */}}

<!-- ... contenido de la pantalla ... -->

{{ template "notifications/oficios_list.html" . }}
```

El template Go renderiza automáticamente los `data-*` correctos:

```html
<!-- Resultado renderizado -->
<div id="oficios-list-root" data-case-id="019d020a-..."></div>
```

---

### Método 2 — `window.OficiosListConfig` desde JS (Angular-like)

Ideal cuando la pantalla ya tiene su propio archivo JS (`miPantalla.js`) y el ID está disponible en el cliente. El script de la pantalla publica la configuración **antes** de que `oficios-list.js` cargue.

El orden de ejecución de scripts garantiza que `window.OficiosListConfig` esté disponible cuando el componente se monta.

#### En el JS de la pantalla (`miPantalla.js`)

```javascript
// Al final del IIFE de la pantalla, después de montar la app principal:
home.mount('#app');

// Pasar parámetros al componente hijo OficiosList
window.OficiosListConfig = {
    caseId: '019d020a-16cf-7c25-8744-5c1fd6f43a72',
    // barrierId: '34e367b4-ea59-4936-b080-5e5491a60380',
    // userId: cfg.currentUserId,
};
```

#### En el HTML de la pantalla

```html
{{/* mi_pantalla.html */}}

<!-- Primero el script de la pantalla (publica OficiosListConfig) -->
<script src="/static/js/components/miPantalla.js"></script>

<!-- Luego el componente (lee OficiosListConfig al montarse) -->
{{ template "notifications/oficios_list.html" . }}
```

> El componente también acepta el ID de la variable de sesión del usuario:
> ```javascript
> window.OficiosListConfig = { userId: cfg.currentUserId };
> ```

---

### Método 3 — Sin parámetros (lista general)

Si no se pasa ningún prop por ninguno de los dos métodos, el componente hace un `GET /api/v1/entity-letters?limit=100&page=0` y muestra todos los oficios paginados.

```html
{{/* Solo incluir el template, sin variables ni config JS */}}
{{ template "notifications/oficios_list.html" . }}
```

---

## Ejemplo completo — Pantalla de detalle de caso

Este ejemplo muestra cómo integrar el componente en una hipotética pantalla de detalle de caso (`case_detail.html`) que ya tiene el `caseId` disponible.

### `CaseDetailFacade.go`

```go
package salvia_facades

import (
    common_facades "bitsflow/common/facades"
    "bitsflow/common/utils"
    salvia_config "bitsflow/salvia/config"
    "net/http"

    "github.com/gin-contrib/sessions"
    "github.com/gin-gonic/gin"
)

var caseDetailTemplates = []string{
    // Partials propios de la pantalla
    "frontend/html/salvia/case_detail/case_detail.html",
    // Componente reutilizable de oficios
    "frontend/html/salvia/notifications/oficios_list.html",
}

func CaseDetailGET(c *gin.Context) {
    session := sessions.Default(c)
    sessionIDVal := session.Get("userData")
    if sessionIDVal == nil {
        c.Redirect(http.StatusTemporaryRedirect, "/static/landing.html")
        return
    }

    sessionID := sessionIDVal.(string)
    s, err := utils.GetCommonSession(sessionID)
    if err != nil {
        c.Redirect(http.StatusTemporaryRedirect, "/static/landing.html")
        return
    }

    caseId := c.Param("caseId")

    extraTemplates := append(utils.GetFullHtmlTemplates(), caseDetailTemplates...)

    common_facades.RenderTemplate(
        c,
        "case-detail",
        "salvia",
        "case_detail/",
        salvia_config.HTML_Templates,
        "case_detail",
        extraTemplates,
        utils.DEFAULT_VIEW,
        utils.DEFAULT_PANIC_TEMPLATE,
        map[string]interface{}{
            "windowTitle":        "Detalle del Caso",
            "currentUser":        s.Names + " " + s.LastNames,
            "currentRole":        s.CurrentRole,
            "currentUserId":      s.UserICode,
            "caseId":             caseId,
            // Pasar el caseId al componente de oficios
            "OficiosListCaseId":  caseId,
        },
        utils.GetFullHtmlFuncMap(),
    )
}
```

### `case_detail.html`

```html
{{/* Pantalla de detalle de caso */}}

<div class="container">
    <h2>Detalle del Caso</h2>

    <!-- ... información del caso ... -->

    <!-- Componente de oficios filtrado por este caso -->
    {{ template "notifications/oficios_list.html" . }}
</div>

{{ template "standard_scripts.html" . }}
{{ template "overlay.html" . }}

<script src="/static/js/components/caseDetail.js"></script>
```

---

## Ejemplo completo — Pantalla de detalle de barrera

```go
// En el facade de barrera:
map[string]interface{}{
    "windowTitle":           "Detalle de Barrera",
    "barrierId":             barrierId,
    "OficiosListBarrierId":  barrierId,  // filtra el componente por esta barrera
}
```

```html
{{/* barrier_detail.html */}}

<!-- ... información de la barrera ... -->

{{ template "notifications/oficios_list.html" . }}
```

---

## Ejemplo completo — Pantalla de seguimiento por usuario

Útil en pantallas donde se quieren mostrar todos los oficios donde el usuario participa (como agente de seguimiento o de notificaciones).

```go
// En el facade:
map[string]interface{}{
    "windowTitle":       "Mis Oficios",
    "OficiosListUserId": s.UserICode,  // userId del usuario en sesión
}
```

```html
{{/* mis_oficios.html */}}
{{ template "notifications/oficios_list.html" . }}
```

> Esta variante usa `EntityLetterWithRelations` como respuesta, por lo que el modal mostrará las cards completas de **Caso** y **Barrera** con nombre de víctima, sectores y descripción.

---

## Checklist de integración

Al agregar el componente a una pantalla nueva verificar:

- [ ] `"frontend/html/salvia/notifications/oficios_list.html"` está en el slice de `extraTemplates` del facade
- [ ] Se pasa al menos una de las tres variables: `OficiosListCaseId`, `OficiosListBarrierId` o `OficiosListUserId` (opcional, si no se pasa carga todo)
- [ ] El HTML de la pantalla incluye `{{ template "notifications/oficios_list.html" . }}`
- [ ] Si se usa el Método 2 (JS), `window.OficiosListConfig` se asigna **antes** de que `oficios-list.js` cargue

---

## Estructura interna del modal de detalle

El modal que abre el botón "Ver" organiza la información en secciones condicionales:

| Sección | Siempre visible | Condición de visibilidad |
|---|---|---|
| **Oficio** | Sí | — |
| **Radicación y Respuesta** | No | Si existe `numero_radicado`, `asunto_radicado` o `response_date` |
| **Corrección requerida** | No | Si existe `reason_correction` |
| **Barrera relacionada** | No | Si existe `barrier_id` |
| **Caso relacionado** | No | Si existe `case_id` |
