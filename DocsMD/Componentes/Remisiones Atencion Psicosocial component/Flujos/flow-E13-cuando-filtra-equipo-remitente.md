━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando filtra por equipo remitente
   Tipo: User Interaction
   Código: E-13
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Disparado por: cambio en DropdownFilter `equipo_remitente`

> Filtra por el equipo del usuario guardado en `submitted_by`.
> El componente **solo lee** ese id y resuelve el team vía JOIN — no implementa
> cómo se persiste `submitted_by` al crear la remisión.


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Opciones cargadas en E-01:
  `GET /api/v1/psychosocial-support/equipos-remitentes`
  → equipos distintos de `submitter_gu.general_user_team` donde `submitted_by IS NOT NULL`

PASO 1 — activeFilters['equipo_remitente'] = value (o delete si "")

PASO 2 — fetchRemisiones()

  → FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  BACKEND — filtro listado
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

```sql
AND submitter_gu.general_user_team = {filter_equipo_remitente}
```

Remisiones con `submitted_by` NULL quedan excluidas al filtrar por equipo concreto.


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  BACKEND — endpoint equipos-remitentes (E-01)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

```sql
SELECT DISTINCT submitter_gu.general_user_team AS team
FROM salvia.psychosocial_support ps
JOIN security.general_user submitter_gu
     ON submitter_gu.general_user_i_code = ps.submitted_by
WHERE ps.deleted_at IS NULL
  AND ps.submitted_by IS NOT NULL
  AND submitter_gu.general_user_team <> ''
ORDER BY team
```

**Nota:** `submitted_by` se lee tal cual está en BD. El listado y el filtro de
equipo remitente resuelven nombre y team vía JOIN. Cómo se guarda el id al
crear la remisión está fuera del alcance de este componente.
