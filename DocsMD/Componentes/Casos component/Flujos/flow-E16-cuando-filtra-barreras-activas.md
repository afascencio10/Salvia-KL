━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando filtra por barreras activas
   Tipo: User Interaction
   Código: E-16
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Disparado por: cambio de valor en el DropdownFilter **"Barreras activas"**
               (`filter.type === 'dropdown'`, `filter.key === 'barreras_activas'`)
               Handler: `setDropdownFilter('barreras_activas', value)`

> **Definición de negocio:** Filtra casos que tienen **al menos una** barrera con
> `status = 'OPEN'` en `salvia.barrier_v2` (modelo `BarrierV2` en
> `src/internal/models/barrier_v2.go`). Un caso puede tener varias barreras abiertas
> en distintos sectores. El dropdown expone una sola opción operativa: `OPEN`
> (etiqueta UI **"Abierta"**). Valor vacío `""` quita el filtro.

INPUT: {
  filterValue:   valor seleccionado en el `<select>`
                 opciones UI (quemadas en el padre):
                   'OPEN'   → etiqueta "Abierta"
                 valor vacío "" → quitar filtro de activeDropdowns['barreras_activas']
}


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  MAPEO VALOR → ETIQUETA
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

| `filter_value` | Etiqueta en UI | Campo BD              | Modelo Go                    |
|----------------|----------------|-----------------------|------------------------------|
| `OPEN`         | Abierta        | barrier_v2.status     | `BarrierV2StatusOpen`        |


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND — casos-component.js
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━


PASO 1 — Verificar si el filtro ya está activo

  SI activeDropdowns['barreras_activas'] === filterValue:
    → No hacer nada
    → TERMINAR ejecución

  SI filterValue === "":
    → delete activeDropdowns['barreras_activas']
    → CONTINÚA PASO 2

  SI NO:
    → activeDropdowns['barreras_activas'] = filterValue   // 'OPEN'
    → CONTINÚA PASO 2

  // Combinable con chip casos_nuevos, riesgo, equipo, seguimientos, estado_caso,
  // persona_asignada, search, sort


PASO 2 — Resetear búsqueda y paginación

  searchText   = ""
  currentPage  = 1
  loading      = true
  loadError    = null
  cases        = []


PASO 3 — Consultar backend

  GET /api/v1/cases/list
    ?filter_barreras_activas=OPEN
    &filter_key=persona_asignada         // si prop agentId (Mis casos)
    &filter_value={agentId}
    &chip_filter=casos_nuevos            // si chip activo (E-09)
    &filter_riesgo=...                   // otros dropdowns activos
    &filter_seguimientos_ejecutados=...
    &filter_estado_caso=...
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
  → Columna `barriers` (si visible): formatOpenBarriers(case.openBarriers) por fila
  → EmptyState si no hay resultados

  → FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  BACKEND — GET /api/v1/cases/list
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

**Fuente de dato:** `salvia.barrier_v2` — modelo `BarrierV2`  
**Relación:** `barrier_v2.case_id = victim_case.victim_case_i_code`  
**Subquery de barreras OPEN en SELECT base:** ya definida en E-01 PASO 5 (`open_barriers_json`).

PASO 5 — Agregar cláusula WHERE

  CASO filter_barreras_activas presente y valor = 'OPEN':

    WHERE EXISTS (
      SELECT 1
      FROM   salvia.barrier_v2 b
      WHERE  b.case_id   = vc.victim_case_i_code
      AND    b.status    = 'OPEN'
      AND    b.deleted_at IS NULL
    )

  // Solo aceptar 'OPEN' como valor válido; otros valores → ignorar cláusula o 400
  // Se combina con AND al resto de filtros aditivos


PASO 6 — Respuesta

  Idéntica a E-01 PASO 9. Cada ítem incluye `openBarriers` con todas las barreras
  OPEN del caso (puede ser array con varios elementos). La columna `barriers` en UI
  agrega sectores únicos: "Abierta → Sector: Salud, Protección".


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  CONFIGURACIÓN EN EL PADRE (pantalla consumidora)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Registrar en `availableFilters`:

```js
{
  key: 'barreras_activas',
  type: 'dropdown',
  label: 'Barreras activas',
  options: [
    { value: 'OPEN', label: 'Abierta' },
  ]
}
```

Registrar en `columns` (si la pantalla debe mostrar la columna):

```js
{ key: 'barriers', label: 'Barreras' }
```

Helpers en `casos-component.js` (implementación):

```js
barrierSectorLabel: function(code) {
  var map = {
    salud: 'Salud',
    justicia: 'Justicia',
    proteccion: 'Protección',
    otras_instituciones: 'Otras instituciones',
    barrera_transversal: 'Barrera Transversal'
  };
  return map[code] || code;
},

formatOpenBarriers: function(openBarriers) {
  if (!openBarriers || openBarriers.length === 0) return '';
  var sectors = [];
  openBarriers.forEach(function(b) {
    var label = this.barrierSectorLabel(b.sector);
    if (label && sectors.indexOf(label) === -1) sectors.push(label);
  }, this);
  if (sectors.length === 0) return 'Abierta';
  return 'Abierta → Sector: ' + sectors.join(', ');
}
```


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

| Variable / decisión | Propuesta |
|---------------------|-----------|
| ¿`case_id` en barrier_v2 usa `victim_case_i_code` o UUID `victim_case_id`? | Confirmar con datos reales |
| filter_value distinto de OPEN | Ignorar cláusula o responder 400 |
| ¿Mostrar columna barreras solo en pantallas con filtro activo? | No — columna independiente del filtro |
| Estados Articulada / MANAGED en columna | No mostrar (solo OPEN en E-01) |


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
📋 ARCHIVOS A TOCAR EN IMPLEMENTACIÓN (referencia)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

| Capa | Archivo |
|------|---------|
| Docs | Este flujo, E-01, E-02, `casos-component-events.md`, `casos-component-interface.md` |
| Frontend | `casos-component.js`, pantallas consumidoras |
| Backend | `cases_list_repository.go`, `cases_list_controller.go`, `cases_list_service.go` |
| Modelo | `src/internal/models/barrier_v2.go` |
