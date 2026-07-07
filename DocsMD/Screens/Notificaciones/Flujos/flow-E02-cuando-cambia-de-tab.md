# flow-E02 — Cuando cambia de tab

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando cambia de tab
   Tipo: User Interaction
   Funciones: switchTab() · loadOficios()
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  tab: 'todos' | 'gestionar'
}

PASO 1 — Actualizar estado local
  currentTab  = tab
  currentPage = 0

PASO 2 — Recargar datos desde el backend

GET /api/v1/entity-letters?{params}

Parámetro adicional según tab:
  tab === 'gestionar' → manageableOnly=true
  tab === 'todos'     → (sin manageableOnly)

El filtro de tab se aplica en BD, no en memoria.

PASO 3 — Renderizar tabla con la primera página de resultados
```
