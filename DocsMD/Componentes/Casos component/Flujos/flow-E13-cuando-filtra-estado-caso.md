━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando filtra por estado del caso
   Tipo: User Interaction
   Código: E-13
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Disparado por: cambio de valor en el DropdownFilter **"Estado del caso"**
               (`filter.type === 'dropdown'`, `filter.key === 'estado_caso'`)
               Handler: `setDropdownFilter('estado_caso', value)`

> **Definición de negocio:** Filtra casos cuyo `salvia.victim_case.victim_case_status`
> coincide exactamente con el código seleccionado. Las etiquetas visibles en el
> dropdown son legibles; el `value` enviado al backend es el **código** de 1–2
> caracteres (misma convención que `get_case_detail_sv.html` → `labelEstado`).

INPUT: {
  filterValue:   código seleccionado en el `<select>`
                 opciones UI (quemadas en el padre):
                   'ra' | 'is' | 'cd' | 'ex' | 'r' | 'fc'
                 valor vacío "" → quitar filtro de activeDropdowns['estado_caso']
}


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  MAPEO CÓDIGO → ETIQUETA (referencia get_case_detail_sv)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

| `filter_value` (código) | Etiqueta en UI     | Campo BD                    |
|-------------------------|--------------------|-----------------------------|
| `ra`                    | Activo             | victim_case_status          |
| `is`                    | Con novedad        | victim_case_status          |
| `cd`                    | Cerrado            | victim_case_status          |
| `ex`                    | Vencido            | victim_case_status          |
| `r`                     | Por aprobar        | victim_case_status          |
| `fc`                    | Recontacto         | victim_case_status          |


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND — casos-component.js
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━


PASO 1 — Verificar si el filtro ya está activo

  SI activeDropdowns['estado_caso'] === filterValue:
    → No hacer nada
    → TERMINAR ejecución

  SI filterValue === "":
    → delete activeDropdowns['estado_caso']
    → CONTINÚA PASO 2

  SI NO:
    → activeDropdowns['estado_caso'] = filterValue
    → CONTINÚA PASO 2

  // Combinable con chip casos_nuevos, riesgo, seguimientos, agentId, search, sort


PASO 2 — Resetear búsqueda y paginación

  searchText   = ""
  currentPage  = 1
  loading      = true
  loadError    = null
  cases        = []


PASO 3 — Consultar backend

  GET /api/v1/cases/list
    ?filter_estado_caso={filterValue}    // ra | is | cd | ex | r | fc
    &filter_key=persona_asignada         // si prop agentId (Mis casos)
    &filter_value={agentId}
    &chip_filter=casos_nuevos            // si chip activo (E-09)
    &filter_riesgo=...                   // otros dropdowns activos
    &filter_seguimientos_ejecutados=...
    &search=...
    &sort={sortBy}
    &order={sortOrder}
    &page=1
    &page_size={pageSize}

  SI respuesta no ok:
    → loadError = data.error || 'Error al aplicar el filtro'
    → TERMINAR ejecución

  SI respuesta ok:
    → cases, totalCases, loading = false
    → CONTINÚA PASO 4


PASO 4 — Actualizar vista

  filteredCases = cases
  → Columna `case_status` muestra caseStatusLabel(case.status) por fila
  → EmptyState si no hay resultados

  → FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  BACKEND — GET /api/v1/cases/list
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

**Fuente de dato:** `salvia.victim_case.victim_case_status`  
**Ya presente en SELECT base (E-01 PASO 5)** — no requiere JOIN adicional.

PASO 5 — Agregar cláusula WHERE

  CASO filter_estado_caso presente
       (o dropdown_filter_key legacy = 'estado_caso'):

    WHERE vc.victim_case_status = {código}

  // Comparación exacta al código de 1–2 caracteres
  // Se combina con AND al resto de filtros aditivos


PASO 6 — Respuesta

  Idéntica a E-01 PASO 9. Cada ítem incluye `status` con el código;
  la columna `case_status` en UI muestra la etiqueta traducida.


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  CONFIGURACIÓN EN EL PADRE (pantalla consumidora)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Registrar en `availableFilters`:

```js
{
  key: 'estado_caso',
  type: 'dropdown',
  label: 'Estado del caso',
  options: [
    { value: 'ra', label: 'Activo' },
    { value: 'is', label: 'Con novedad' },
    { value: 'cd', label: 'Cerrado' },
    { value: 'ex', label: 'Vencido' },
    { value: 'r',  label: 'Por aprobar' },
    { value: 'fc', label: 'Recontacto' },
  ]
}
```

Registrar en `columns` (si la pantalla debe mostrar la columna):

```js
{ key: 'case_status', label: 'Estado del caso' }
```

Helper en `casos-component.js` (implementación):

```js
caseStatusLabel: function(code) {
  var map = {
    ra: 'Activo', is: 'Con novedad', cd: 'Cerrado',
    ex: 'Vencido', r: 'Por aprobar', fc: 'Recontacto'
  };
  return map[code] || '—';
}
```


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

| Variable / decisión | Propuesta |
|---------------------|-----------|
| ¿Badge CSS por estado (badgeEstado en detalle SV)? | Solo texto en v1 |
| Código desconocido en columna | Mostrar "—" |
| filter_value inválido en backend | Ignorar cláusula o 400 |
| ¿Estados adicionales en BD no listados en labelEstado? | Pendiente catálogo completo |


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
📋 ARCHIVOS A TOCAR EN IMPLEMENTACIÓN (referencia)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

| Capa | Archivo |
|------|---------|
| Docs | Este flujo, E-01, E-02, `casos-component-events.md`, `casos-component-interface.md` |
| Frontend | `casos-component.js`, `case_component.html`, `get_my_cases.html` |
| Backend | `cases_list_repository.go`, `cases_list_controller.go`, `cases_list_service.go` |
