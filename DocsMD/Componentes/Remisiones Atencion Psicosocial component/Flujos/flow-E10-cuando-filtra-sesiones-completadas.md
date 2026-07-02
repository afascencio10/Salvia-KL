━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando filtra por sesiones completadas
   Tipo: User Interaction
   Código: E-10
   Prerequisito: M-02 (tabla team_contact)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Disparado por: cambio en DropdownFilter `sesiones_completadas`


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Opciones fijas del dropdown: `{ value:'0', label:'0' }` … `{ value:'6', label:'6' }`
(Alineado con los 6 puntos de la barra de sesiones en E-01)

PASO 1 — activeFilters['sesiones_completadas'] = value (o delete si "")

PASO 2 — currentPage = 1; selectedRemisiones = []; fetchRemisiones()

  → FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  BACKEND
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

El conteo de sesiones **no** proviene de una columna en `psychosocial_support`.
Se calcula desde `salvia.team_contact`:

```sql
AND (
  SELECT COUNT(*)::int
  FROM salvia.team_contact tc
  WHERE tc.psicosocial_id = ps.id
    AND tc.is_psico_session = true
    AND tc.is_completed = true
    AND tc.deleted_at IS NULL
) = {filter_sesiones_completadas}::int
```

Coincidencia **exacta** con el número de sesiones psicosociales completadas.

> **Cambio post-reunión:** reemplaza filtro por `psychosocial_support.session_count`.
