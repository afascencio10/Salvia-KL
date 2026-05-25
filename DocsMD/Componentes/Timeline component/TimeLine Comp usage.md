# Componente `case-timeline`

Componente Vue reutilizable que renderiza el timeline de eventos de un caso. Sigue el mismo patrón de `dinamic-form`: es autónomo, inyecta sus propios estilos y hace su propio fetch al backend.

---

## Archivos relevantes

| Archivo | Descripción |
|---|---|
| `src/frontend/js/components/case-timeline.js` | Definición del componente (estilos + lógica + template) |
| `src/internal/repository/case_timeline_event_repository.go` | Repositorio — método `GetByCaseID` con JOIN al nombre del actor |
| `src/internal/models/case_timeline_event.go` | Modelo `CaseTimelineEvent`, tabla `salvia.case_timeline_event` |
| `src/salvia/controller/case_detail_controller.go` | Endpoint `GET /api/v1/casos/:id/timeline-events` |

---

## Props

| Prop | Tipo | Requerido | Descripción |
|---|---|---|---|
| `case-id` | `String` | ✅ Sí | `icode` del caso (`VictimCase.VictimCaseICode`) |
| `barrier-id` | `String` | ❌ No | UUID de una barrera. Si se pasa, el componente solo muestra los eventos vinculados a esa barrera |

---

## Cómo incluirlo en una ruta nueva

### 1. Registrar el componente antes del `mount`

El componente llama a `app.component(...)` cuando se carga, por lo que el script debe incluirse **después** de que la instancia Vue esté creada pero **antes** del `mount`.

El patrón es siempre el mismo — al final del bloque `<script>` principal, justo antes de cerrar:

```html
<!-- Al final del <script> principal, antes de hacer mount -->
var app = home;   <!-- alias necesario si tu app se llama "home" -->
</script>

<script src="/static/js/components/case-timeline.js"></script>

<script>
home.mount('#app');
</script>
```

> **Nota:** Si tu app ya se llama `app` (ej. `var app = Vue.createApp({...})`), no necesitas el alias — solo agrega el `<script src>` antes del `mount`.

---

### 2. Usar el componente en el template

**Caso básico** — muestra todos los eventos del caso:

```html
<case-timeline :case-id="caseICode"></case-timeline>
```

**Filtrado por barrera** — muestra solo los eventos vinculados a una barrera específica:

```html
<case-timeline :case-id="caseICode" :barrier-id="barreraActual.id"></case-timeline>
```

**Con valor estático desde Go template** (sin binding Vue):

```html
<case-timeline case-id="{{.caseICode}}"></case-timeline>
```

---

### 3. Exponer `caseICode` en el data de Vue

El componente necesita recibir el `icode` del caso como prop. La forma más limpia es tenerlo en el `data()` de la app:

```js
var CASE_ICODE = "{{.caseICode}}";  // variable Go template → JS

var home = Vue.createApp({
  data() {
    return {
      caseICode: CASE_ICODE,
      // ... resto del data
    };
  }
});
```

---

## Qué hace el componente internamente

1. **Al montar** llama a `GET /api/v1/casos/:caseId/timeline-events?barrierId=xxx`
2. El backend consulta `salvia.case_timeline_event` con un LEFT JOIN a `security.general_user` + `security.general_user_profile` usando `event_user_id = general_user_i_code` para resolver el nombre del actor
3. Los eventos llegan ordenados por fecha descendente
4. El componente renderiza filtros pill por `category` en la parte superior
5. Cada evento muestra: icono Font Awesome (del campo `icon`), badge con el `type`, fecha, descripción y nombre del actor

---

## Categorías de filtro disponibles

| Valor en DB (`category`) | Label en UI |
|---|---|
| *(vacío)* | Todos |
| `General` | General |
| `Barreras` | Barreras |
| `Seguimientos` | Seguimientos |
| `Oficios` | Oficios |
| `Medidas` | Medidas |
| `Psicosocial` | Psicosocial |
| `Estabilización` | Estabilización |

---

## Campo `icon` — valores esperados

El campo `icon` del modelo acepta dos formatos, el componente los normaliza:

| Formato guardado | Ejemplo | Resultado renderizado |
|---|---|---|
| Clase FA completa | `"fa fa-calendar"` | Se usa tal cual |
| Solo nombre del ícono | `"calendar-check"` | El componente agrega `"fa fa-"` → `"fa fa-calendar-check"` |

Íconos usados actualmente en el código:

| Evento | `icon` |
|---|---|
| Seguimiento ejecutado | `"calendar-check"` |
| Seguimiento editado | `"fa fa-calendar"` |
| Seguimiento pospuesto | `"fa fa-calendar"` |
| Intento fallido | `"fa fa-calendar"` |

> **Convención a futuro:** guardar solo el nombre del ícono sin prefijo (`"calendar-check"`, `"user"`, etc.) y dejar que el componente construya la clase.

---

## Ejemplo completo — nueva ruta

```html
<!-- mi_nueva_ruta.html -->

<div id="app" v-cloak>
  <!-- ... contenido de la ruta ... -->

  <!-- Tab o sección de timeline -->
  <div v-if="seccionActiva === 'timeline'">
    <case-timeline :case-id="caseICode"></case-timeline>
  </div>
</div>

<script>
var CASE_ICODE = "{{.caseICode}}";

var home = Vue.createApp({
  delimiters: ['${', '}'],
  data() {
    return {
      caseICode: CASE_ICODE,
      seccionActiva: 'timeline',
    };
  },
  // ...
  var app = home;
});
</script>

<script src="/static/js/components/case-timeline.js"></script>

<script>
home.mount('#app');
</script>
```
