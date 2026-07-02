━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
📋 ÍNDICE — Filtros dropdown (E-02)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

El evento **E-02** agrupa los filtros de tipo `dropdown`. Cada uno tiene flujo
independiente con cláusula SQL propia en el backend.

| Código | Filtro UI | `filter.key` | Flujo |
|--------|-----------|--------------|-------|
| **E-09** | Estado remisión | `estado_remision` | [flow-E09](./flow-E09-cuando-filtra-estado-remision.md) |
| **E-10** | Sesiones completadas | `sesiones_completadas` | [flow-E10](./flow-E10-cuando-filtra-sesiones-completadas.md) |
| **E-11** | Dupla asignada | `dupla_asignada` | [flow-E11](./flow-E11-cuando-filtra-dupla-asignada.md) |
| **E-12** | Nivel de riesgo | `nivel_riesgo` | [flow-E12](./flow-E12-cuando-filtra-nivel-riesgo.md) |
| **E-13** | Equipo remitente | `equipo_remitente` | [flow-E13](./flow-E13-cuando-filtra-equipo-remitente.md) |

---

## Comportamiento común (E-09 … E-13)

1. Resetea `currentPage = 1` y limpia `selectedRemisiones` (→ E-14).
2. Actualiza `activeFilters[filterKey] = value` o elimina la key si value === `""`.
3. Llama a `fetchRemisiones()` → `GET /api/v1/psychosocial-support/list` con **todos** los filtros activos (AND).
4. Mantiene `sortBy = 'created_at'` y `sortOrder = 'desc'` vigentes.
5. No altera los inputs de búsqueda (`numero_identidad`, `telefono`) ni el autocomplete ya seleccionado.
6. No elimina `defaultFilter` (`scopeFilter`) — `agent_id` / `dupla_id` del prop persisten.
7. Tras cada cambio: `fetchRemisiones()` y, si `mostrarCards === true`, `fetchStats()` con los mismos filtros.

---

## Filtros con flujo propio (fuera de E-02)

| Filtro | Tipo | Evento |
|---|---|---|
| Número de identidad | `search` | E-03 |
| Teléfono | `search` | E-04 |
| Profesional asignada | `autocomplete` | E-07 + E-08 |
| Limpiar filtros | botón | E-05 |

---

## Handler en `remisiones-psicosocial-component.js`

```js
setDropdownFilter(filterKey, value) {
  // Delega según filterKey → E-09 … E-13
  // Todos comparten fetchRemisiones() al final
}
```
