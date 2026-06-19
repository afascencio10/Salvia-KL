# Listado de Reportes — Interfaz de la pantalla

Pantalla legacy adaptada desde `get_victim_cases_sv.html` para listar **solo reportes** (`victim_contact` sin caso). Visible para roles **`sv`** y **`ro`**.

**Regla de implementación:** no eliminar líneas existentes — **comentar** con bloque `<!-- REPORTES-ONLY: ... -->` o `// REPORTES-ONLY: ...` en JavaScript, indicando qué se ocultó y por qué.

---

## Archivos HTML

| Archivo | Rol |
|---|---|
| `src/frontend/html/salvia/victim_case/get_victim_cases_sv.html` | Supervisor (`sv`) — **base de referencia** |
| `src/frontend/html/salvia/victim_case/get_victim_cases_ro.html` | Revisor operativo (`ro`) — cambios espejo |

Entry point backend: `GET /salvia/casos` → `VictimCaseFacade.VictimCaseGET` → template según rol.

---

## Árbol de interfaz (estado objetivo)

```
get_victim_cases_sv / get_victim_cases_ro  (.salvia-main)
│
├── WelcomeSection
│   ├── [sv]  h2 "Casos Salvia"
│   └── [ro]  h2 + botones header (Crear Nuevo caso, Mis Casos)
│               // RECOMENDACIÓN: comentar botones de casos en header de ro
│               // si la pantalla pasa a ser exclusiva de reportes
│
├── [COMENTAR] Sección "Encontrar caso"
│   // Buscador por tipo/número de documento → consulta victimCases
│   // No aplica en pantalla solo-reportes
│
├── hr
│
├── NavTabs  (.nav)
│   ├── Tab "Recontacto"          → goToTab('recontact_v')  → filter fcv
│   ├── Tab "Inválidos"           → goToTab('recontact_i')  → filter fci
│   ├── [COMENTAR] Enrutados      → routing / ra
│   ├── [COMENTAR] Enrutados por aprobar → r
│   ├── [COMENTAR] Vencidos      → ex
│   ├── [COMENTAR] Completados    → cd
│   ├── [COMENTAR] Con novedades  → is
│   ├── [sv] Botón notificaciones → getAlerts()
│   ├── [sv] Botón mapa reportes  → getReports()   // opcional comentar
│   └── [sv] Botón reasignar      → assignOperator() // opcional comentar
│
├── FilterBarReportes  (.search-form)  ← NUEVO
│   ├── Input  "Nombres de la víctima"     v-model=filterNames
│   ├── Input  "Apellidos de la víctima"   v-model=filterLastNames
│   ├── Input  "Teléfono de contacto"      v-model=filterPhone
│   ├── Btn "Buscar"    → applyReportFilters()   // E-02
│   └── Btn "Limpiar"   → clearReportFilters()   // E-02
│
└── TabContent  (#myTabContent)
    │
    ├── [ACTIVO] Bloque reportes  (v-if currentFilter == 'fcv' || 'fci')
    │   ├── EmptyState  si victimContacts.length == 0
    │   ├── Table  (.table)
    │   │   thead: Fecha | Nombres y Apellidos | Riesgo feminicidio | Estado | Acciones
    │   │   tbody: v-for victimContact in victimContacts
    │   │       Estado: tag "Recontacto" (v) | "Recontacto Inválido" (i)
    │   │       Acciones: menuToolsContactTable
    │   │         → Ver:    /salvia/primer_contacto/{icode}/
    │   │         → Crear caso: /salvia/casos/{icode}/nuevo
    │   └── PaginationBar
    │
    └── [COMENTAR] Bloque casos  (v-if currentFilter != 'fcv' && != 'fci')
        // Tabla victimCases completa + paginación + acciones de caso
```

---

## Cambios por sección — `get_victim_cases_sv.html`

### 1. Comentar sección "Encontrar caso" (líneas ~7–32)

Envolver en comentario HTML:

```html
<!-- REPORTES-ONLY: buscador por documento de casos — deshabilitado
<div class="mt-3 mb-4"> ... </div>
-->
```

También comentar método `searchCase()` en el `<script>` o dejarlo comentado con nota.

### 2. Comentar tabs de casos en `<ul class="nav">` (líneas ~41–55)

Mantener activos solo:

```html
<li class="nav-item" v-on:click="goToTab('recontact_v')">...</li>
<li class="nav-item" v-on:click="goToTab('recontact_i')">...</li>
```

Comentar: `routing`, `routedToApprove`, `expired`, `completed`, `issues`.

### 3. Agregar FilterBarReportes (después de `<hr>`, antes de tabs o después de tabs)

```html
<div class="mt-3 mb-4">
  <h3>Filtrar reportes</h3>
  <div class="search-form">
    <!-- 3 inputs + botones Buscar / Limpiar -->
  </div>
</div>
```

Estado Vue nuevo en `data()`:

```javascript
filterNames: "",
filterLastNames: "",
filterPhone: "",
```

### 4. Comentar bloque tabla de casos (líneas ~132–212)

```html
<!-- REPORTES-ONLY: tabla de casos — deshabilitada
<div v-if="currentFilter != 'fcv' && currentFilter != 'fci'"> ... </div>
-->
```

### 5. Comentar overlay reasignación (opcional, líneas ~216–237)

Solo aplica a casos. Comentar `#assignOperatorOverlay` y métodos `submit('openAssign'|'assign'|'closeAssign')` si ya no hay tabla de casos.

### 6. Simplificar `goToTab()` en script

Comentar cases del `switch` excepto `recontact_v` y `recontact_i`:

```javascript
switch (tabName) {
  case "recontact_v": filter = "fcv"; break;
  case "recontact_i": filter = "fci"; break;
  // REPORTES-ONLY: tabs de casos comentados
  // case "routing": filter = "ra"; break;
  // ...
}
```

### 7. Adaptar `goToPage()` y agregar métodos de filtro

- Eliminar rama que llama `salviaVictimCaseGETFormPath` (casos) — **comentar**, no borrar.
- En rama de contactos, append query params de filtro:

```javascript
buildContactListUrl() {
  var url = "{{.salviaVictimContactGETFormPath}}/f/" + this.currentFilter + "/p/" + this.currentPage;
  var params = [];
  if (this.filterNames.trim())     params.push("names="     + encodeURIComponent(this.filterNames.trim()));
  if (this.filterLastNames.trim()) params.push("lastNames=" + encodeURIComponent(this.filterLastNames.trim()));
  if (this.filterPhone.trim())     params.push("phone="     + encodeURIComponent(this.filterPhone.trim()));
  if (params.length) url += "?" + params.join("&");
  return url;
}
```

---

## Cambios espejo — `get_victim_cases_ro.html`

Aplicar los **mismos bloques** numerados arriba en las secciones equivalentes:

| Sección sv | Equivalente ro (aprox.) |
|---|---|
| Líneas 7–32 buscador caso | Líneas 17–47 |
| Nav tabs 41–55 | Nav tabs 56–73 |
| Tabla casos 132–212 | Tabla casos 140–212 |
| Script goToTab/goToPage | Script 387–419 |

**Diferencias a preservar en `ro`:**

- Header con botones "Crear Nuevo caso" y "Mis Casos" — comentar si la pantalla queda solo-reportes, o dejar con nota TODO.
- No tiene botones `getReports()` ni `assignOperator()` del sv.

---

## Botón "Crear caso" — todos los roles

### Tabla de acciones (`menuToolsContactTable`)

Generada en backend desde `Menu.go` → `menu_tool_get_victim_contacts`.

| Rol | Estado actual | Cambio |
|---|---|---|
| `op` | Ver + Crear caso | Sin cambio |
| `ro` | Ver + Crear caso | Sin cambio |
| `et` | Ver + Crear caso | Sin cambio |
| `sv` | Solo Ver | **Agregar** entrada `menu_tool_new_victim_case` |

Agregar en `Menu.go` bloque `sv`:

```go
{
    "label":       Locale["sp"]["menu_tool_new_victim_case"],
    "path":        "/" + Locale["sp"]["salvia"] + "/" + Locale["sp"]["VictimCase"],
    "action":      Locale["sp"]["new"],
    "customClass": "fas fa-plus-circle",
},
```

URL resultante en tabla: `/salvia/casos/{victimContact.icode}/nuevo`

### Permisos

`set_victim_case` ya está habilitado para `sv` y `ro` en `PermissionsByRole`. No requiere cambio de permisos.

### Detalle del reporte (opcional fase 2)

Agregar botón "Crear caso" en `get_victim_contact_sv_v2.html` (hoy solo tiene Regresar), alineado con `get_victim_contact.html`.

---

## Columnas de la tabla (sin cambio inicial)

| Columna | Fuente |
|---|---|
| Fecha Creación | `victimContact.creationDate` |
| Nombres y Apellidos | `lastNames`, `names` |
| Riesgo de feminicidio | *(vacío hoy)* |
| Estado | `status` → tag Recontacto / Inválido |
| Acciones | `menuToolsContactTable` |

> Fase posterior: enriquecer columnas con datos de `victim_contact_form2` (tipo reporte, teléfono, etc.) alineado con `reportes-component`.

---

## Convención de comentarios

```html
<!-- REPORTES-ONLY [2026-06]: <descripción breve>
     Motivo: pantalla adaptada a solo reportes (sv/ro)
     Original: ...
-->
... código comentado ...
<!-- /REPORTES-ONLY -->
```

```javascript
// REPORTES-ONLY [2026-06]: rama de paginación de casos
// if (this.currentFilter != "fcv" && ...) { ... }
```
