━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando filtra por número de seguimientos ejecutados
   Tipo: User Interaction
   Código: E-12
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Disparado por: cambio de valor en el DropdownFilter **"Seguimientos ejecutados"**
               (`filter.type === 'dropdown'`, `filter.key === 'seguimientos_ejecutados'`)
               Handler: `setDropdownFilter('seguimientos_ejecutados', value)`

> **Definición de negocio:** Filtra casos cuya cantidad de seguimientos con
> `status = 'REALIZADO'` en `salvia.follow_up_v2` coincide exactamente con el
> número seleccionado. Es el mismo conteo que muestra la columna
> `completed_follow_ups` (E-01).

INPUT: {
  filterValue:   valor seleccionado en el `<select>` (string numérica)
                 opciones UI: '0' | '1' | '2' | ... | '10'  (rango propuesto)
                 valor vacío "" → resetear al defaultFilter
}


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND — casos-component.js
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━


PASO 1 — Verificar si el filtro ya está activo

  SI activeDropdown.key === 'seguimientos_ejecutados'
     Y activeDropdown.value === filterValue:
    → No hacer nada
    → TERMINAR ejecución

  SI filterValue === "":
    → activeDropdown = null (solo si key === 'seguimientos_ejecutados')
    → CONTINÚA PASO 2

  SI NO:
    → activeDropdown = { key: 'seguimientos_ejecutados', value: filterValue }
    → CONTINÚA PASO 2

  // activeDropdown es independiente de agentId y activeChipKey


PASO 2 — Resetear búsqueda y paginación

  searchText   = ""
  currentPage  = 1
  loading      = true
  loadError    = null
  cases        = []


PASO 3 — Consultar backend

  GET /api/v1/cases/list
    ?dropdown_filter_key=seguimientos_ejecutados
    &dropdown_filter_value={filterValue}   // "0" | "1" | "2" | ...
    &filter_key=persona_asignada           // si prop agentId (Mis casos)
    &filter_value={agentId}
    &chip_filter=casos_nuevos              // si chip activo (E-09)
    &search=...                            // si búsqueda activa (E-03)
    &sort={sortBy}
    &order={sortOrder}
    &page=1
    &page_size={pageSize}

  // dropdown_filter_* es aditivo: combina con agentId, chip, search y paginación

  SI respuesta no ok:
    → loadError = data.error || 'Error al aplicar el filtro'
    → TERMINAR ejecución

  SI respuesta ok:
    → cases, totalCases, loading = false
    → CONTINÚA PASO 4


PASO 4 — Actualizar vista

  filteredCases = cases
  → Renderizar tabla (columna completed_follow_ups refleja el conteo filtrado)
  → EmptyState si no hay resultados

  → FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  BACKEND — GET /api/v1/cases/list  (dropdown_filter_key = seguimientos_ejecutados)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

**Fuente de dato:** `salvia.follow_up_v2` — modelo `FollowUpV2`  
**Campo de relación:** `follow_up_v2.case_id` = `victim_case.victim_case_i_code`  
**Estado contado:** `FollowUpStatusRealizado` = `"REALIZADO"`  
**Subconsulta (misma que E-01 para la columna):**

```sql
(
  SELECT COUNT(*)::int
  FROM   salvia.follow_up_v2 fu
  WHERE  fu.case_id = vc.victim_case_i_code
  AND    fu.status  = 'REALIZADO'
  AND    fu.deleted_at IS NULL
) AS completed_follow_ups_count
```

PASO 5 — Agregar cláusula WHERE

  CASO dropdown_filter_key = 'seguimientos_ejecutados'
       (o filter_key legacy = 'seguimientos_ejecutados'):

    WHERE (
      SELECT COUNT(*)::int
      FROM   salvia.follow_up_v2 fu
      WHERE  fu.case_id = vc.victim_case_i_code
      AND    fu.status  = 'REALIZADO'
      AND    fu.deleted_at IS NULL
    ) = {nivel_entero}

  // El backend parsea filter_value a entero (0, 1, 2, ...)
  // Comparación de igualdad exacta, no "al menos N"


PASO 6 — Respuesta

  Idéntica a E-01 PASO 9. Cada ítem incluye `completedFollowUpsCount`
  acorde al filtro aplicado.


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  CONFIGURACIÓN EN EL PADRE (pantalla consumidora)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Registrar en `availableFilters`:

```js
{
  key: 'seguimientos_ejecutados',
  type: 'dropdown',
  label: 'Seguimientos ejecutados',
  options: [
    { value: '0', label: '0' },
    { value: '1', label: '1' },
    { value: '2', label: '2' },
    // ... hasta el máximo acordado (propuesta: 0–10)
  ]
}
```

Registrar en `columns` (si la pantalla debe mostrar la columna):

```js
{ key: 'completed_follow_ups', label: 'Seguimientos ejecutados' }
```


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

| Variable / decisión | Propuesta |
|---------------------|-----------|
| Rango del select (0–10 fijo vs. dinámico desde BD) | 0–10 fijo en v1 |
| ¿Opción "5 o más" para casos con muchos seguimientos? | Pendiente negocio |
| ¿Combinar con chip casos_nuevos y riesgo/equipo? | Sí, aditivo (mismo patrón E-10) |
| Validación filter_value no numérico | 400 o ignorar filtro |


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
📋 ARCHIVOS A TOCAR EN IMPLEMENTACIÓN (referencia)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

| Capa | Archivo |
|------|---------|
| Docs | Este flujo, E-01, E-09, `casos-component-events.md`, `casos-component-interface.md`, `flow-E02` índice |
| Frontend | `casos-component.js`, `case_component.html`, pantallas padre (`get_my_cases.html`, etc.) |
| Backend | `cases_list_repository.go`, `case_list_item.go`, `cases_list_controller.go` |
| Modelo | `follow_up_v2.go` (constante `FollowUpStatusRealizado` ya definida) |
