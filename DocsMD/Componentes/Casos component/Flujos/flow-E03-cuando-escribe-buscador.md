━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando el usuario escribe en el buscador
   Tipo: User Interaction
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Disparado por: input en el SearchFilter (type='search') dentro de FilterGroup
               Se activa en cada cambio del campo — con debounce de 400ms

INPUT: {
  searchText:   texto escrito por el usuario   → v-model del <input type="text">
}

> **Nota:** La búsqueda llama al backend porque el componente usa paginación.
> No es posible buscar entre registros que no están en la página actual.
> El searchText actúa en paralelo al filtro activo: ambos se envían juntos.


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND — casos-component.js
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━


PASO 1 — Actualizar estado y resetear paginación

  searchText  = searchText.trim()
  currentPage = 1        // siempre volver a la primera página al buscar
  loading     = true
  cases       = []


PASO 2 — Verificar si el campo quedó vacío

  SI searchText === "":
    → Llamar al backend igual (PASO 3) para restaurar la lista completa
      del filtro activo (sin parámetro de búsqueda)
    → CONTINÚA PASO 3

  SI searchText !== "":
    → CONTINÚA PASO 3


PASO 3 — Consultar backend combinando filtro activo + texto de búsqueda

  GET /api/v1/cases/list
    ?filter_key={activeFilter.key}
    &filter_value={activeFilter.value}
    &search={searchText}               // vacío si el usuario borró el campo
    &sort={sortBy}
    &order={sortOrder}
    &page=1
    &page_size={pageSize}

  SI respuesta no ok (status != 2xx):
    → loading   = false
    → loadError = data.error || 'Error al buscar casos'
    → TERMINAR ejecución

  SI respuesta ok:
    → cases      = data.cases
    → totalCases = data.total
    → loading    = false
    → CONTINÚA PASO 4


PASO 4 — Actualizar vista

  filteredCases = cases

  SI filteredCases.length === 0:
    → Mostrar EmptyState "No se encontraron casos"

  SI filteredCases.length > 0:
    → Renderizar tabla con filteredCases
    → Renderizar PaginationBar actualizado

  → FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  BACKEND — GET /api/v1/cases/list  (parámetro search)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

El endpoint ya existente recibe el parámetro adicional `search`.

SI search != "":
  → Agregar cláusula WHERE a la query base del PASO 5 de E-01:

    AND (
      vc.victim_case_i_code ILIKE '%' || search || '%'
      OR
      [campo_telefono] ILIKE '%' || search || '%'   // ⚠️ GAP: campo exacto del teléfono
    )

SI search == "" o no viene:
  → No agregar cláusula adicional (comportamiento normal del filtro activo)

El resto del flujo (ordenamiento, paginación, construcción de response)
es idéntico al PASO 7, 8 y 9 de E-01.


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

| Variable / decisión                                                                  | Paso afectado |
|--------------------------------------------------------------------------------------|---------------|
| Campo exacto del teléfono de la víctima y en qué tabla vive                          | PASO 3 (backend) |
| ¿El debounce de 400ms es correcto o se prefiere otro valor?                          | PASO 1        |
| ¿La búsqueda usa ILIKE (parcial) o solo coincidencia exacta por i_code?              | PASO 3 (backend) |
| ¿Se puede combinar búsqueda Y filtro activo al mismo tiempo (actualmente sí)?        | PASO 3        |
