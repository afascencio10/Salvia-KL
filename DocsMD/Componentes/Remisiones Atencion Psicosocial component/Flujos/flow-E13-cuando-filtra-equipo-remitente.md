━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando filtra por equipo remitente
   Tipo: User Interaction
   Código: E-13
   Prerequisito: M-02 (campo submitted_by_team)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Disparado por: cambio en DropdownFilter `equipo_remitente`

> Filtra por la columna `submitted_by_team` en `psychosocial_support`.
> El nombre del remitente sigue resolviéndose vía JOIN en `submitted_by`.
> El componente **solo lee** ambos campos; no implementa cómo se persisten al crear.


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Opciones cargadas en E-01:
  `GET /api/v1/psychosocial-support/equipos-remitentes`
  → equipos distintos de `ps.submitted_by_team` WHERE NOT NULL

PASO 1 — activeFilters['equipo_remitente'] = value (o delete si "")

PASO 2 — fetchRemisiones()

  → FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  BACKEND — filtro listado
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

```sql
AND ps.submitted_by_team = {filter_equipo_remitente}
```

Remisiones con `submitted_by_team` NULL quedan excluidas al filtrar por equipo concreto.


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  BACKEND — endpoint equipos-remitentes (E-01)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

```sql
SELECT DISTINCT ps.submitted_by_team AS team
FROM salvia.psychosocial_support ps
WHERE ps.deleted_at IS NULL
  AND ps.submitted_by_team IS NOT NULL
  AND ps.submitted_by_team <> ''
ORDER BY team
```

> **Cambio post-reunión:** ya no requiere JOIN a `security.general_user`.
> El equipo se lee directamente de `submitted_by_team`.
