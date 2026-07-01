━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando filtra por estado de remisión
   Tipo: User Interaction
   Código: E-09
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Disparado por: cambio en DropdownFilter `estado_remision`
Handler: `setDropdownFilter('estado_remision', value)`


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

PASO 1 — Actualizar activeFilters

  SI value === ""  → delete activeFilters['estado_remision']
  SI NO            → activeFilters['estado_remision'] = value

PASO 2 — Reset paginación y selección → fetchRemisiones()

  → FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  BACKEND
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Estados en BD — español snake_case (patrón `entity_letter.go`):

| `filter_estado_remision` | Label dropdown | Condición de negocio |
|---|---|---|
| `abierto` | Abiertos | Al crear la remisión |
| `en_gestion` | En gestión | Primer contacto logrado |
| `en_devolucion` | En devolución | No cumplió criterios |
| `cerrado` | Cerrados | Por cualquier motivo |

```sql
AND ps.status = {filter_estado_remision}
```

Constantes Go: `PsychosocialSupportStatusAbierto`, `EnGestion`, `EnDevolucion`, `Cerrado`.

Default al crear remisión: `abierto` (ver M-01).
