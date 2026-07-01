━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando filtra por sesiones completadas
   Tipo: User Interaction
   Código: E-10
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Disparado por: cambio en DropdownFilter `sesiones_completadas`


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Opciones fijas del dropdown: `{ value:'0', label:'0' }` … `{ value:'6', label:'6' }`
(Alineado con `MAX_SESSIONS = 6` quemado en E-01)

PASO 1 — activeFilters['sesiones_completadas'] = value (o delete si "")

PASO 2 — currentPage = 1; selectedRemisiones = []; fetchRemisiones()

  → FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  BACKEND
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

```sql
AND ps.session_count = {filter_sesiones_completadas}::int
```

Coincidencia **exacta** con `psychosocial_support.session_count`.
