━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
📋 ÍNDICE — Filtros chip y dropdown (E-02)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

El evento **E-02** original se dividió en tres flujos independientes,
uno por cada control de filtro. Cada uno tiene su propia lógica de negocio
y cláusula SQL en el backend.

| Código | Filtro UI            | Tipo       | Flujo |
|--------|----------------------|------------|-------|
| **E-09** | Casos nuevos         | `chip`     | [flow-E09-cuando-filtra-casos-nuevos.md](./flow-E09-cuando-filtra-casos-nuevos.md) |
| **E-10** | Nivel de riesgo      | `dropdown` | [flow-E10-cuando-filtra-nivel-riesgo.md](./flow-E10-cuando-filtra-nivel-riesgo.md) |
| **E-11** | Por equipo           | `dropdown` | [flow-E11-cuando-filtra-equipo.md](./flow-E11-cuando-filtra-equipo.md) |
| **E-12** | Seguimientos ejecutados | `dropdown` | [flow-E12-cuando-filtra-seguimientos-ejecutados.md](./flow-E12-cuando-filtra-seguimientos-ejecutados.md) |
| **E-13** | Estado del caso        | `dropdown` | [flow-E13-cuando-filtra-estado-caso.md](./flow-E13-cuando-filtra-estado-caso.md) |

---

## Comportamiento común (todos los filtros E-09 / E-10 / E-11 / E-12 / E-13)

1. Resetean `searchText`, `currentPage = 1` y limpian `cases` antes de consultar.
2. Llaman a `GET /api/v1/cases/list` con `filter_key` y `filter_value` (si aplica).
3. Mantienen `sortBy` y `sortOrder` vigentes.
4. Si el usuario elige **"Todos / Todas"** (dropdown) o **desactiva el chip**, vuelven a `defaultFilter`.

---

## Filtros con flujo propio (fuera de E-02)

| Filtro              | Tipo           | Evento |
|---------------------|----------------|--------|
| Persona asignada    | `autocomplete` | E-07 + E-08 |
| Buscar documento/tel. | `search`     | E-03 |

---

## Handlers en `casos-component.js`

| Handler                 | Delega a |
|-------------------------|----------|
| `toggleChipFilter(key)` | E-09 cuando `key === 'casos_nuevos'` |
| `setDropdownFilter(key, value)` | E-10 cuando `key === 'riesgo'` |
| `setDropdownFilter(key, value)` | E-11 cuando `key === 'equipo'` |
| `setDropdownFilter(key, value)` | E-12 cuando `key === 'seguimientos_ejecutados'` |
| `setDropdownFilter(key, value)` | E-13 cuando `key === 'estado_caso'` |
