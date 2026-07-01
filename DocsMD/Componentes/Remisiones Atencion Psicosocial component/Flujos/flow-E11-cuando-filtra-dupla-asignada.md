━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando filtra por dupla asignada
   Tipo: User Interaction
   Código: E-11
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Disparado por: cambio en DropdownFilter `dupla_asignada`

Opciones cargadas en E-01 desde `GET /api/v1/duplas`:
  label = `dupla.name`  ·  value = `dupla.id`


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

PASO 1 — activeFilters['dupla_id'] = value (o delete si "")

PASO 2 — currentPage = 1; selectedRemisiones = []; fetchRemisiones()

  GET /api/v1/psychosocial-support/list
    &filter_dupla_id={value}

  → FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  BACKEND
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

```sql
AND ps.dupla_id = {filter_dupla_id}
```

Opción adicional en UI (opcional): `{ value: '__sin_dupla__', label: 'Sin dupla' }`

```sql
AND ps.dupla_id IS NULL
```

Prerequisito: tabla `salvia.dupla` creada en M-01.
