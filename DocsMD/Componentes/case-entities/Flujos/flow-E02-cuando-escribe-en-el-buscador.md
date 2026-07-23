━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando escribe en el buscador
   Tipo: User Interaction
   Función: v-model="busqueda"  (computed: entidadesFiltradas)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  texto:  valor tecleado por el usuario   → v-model del <input> buscador
}

PASO 1 — Actualizar variable reactiva `busqueda` con el texto tecleado

PASO 2 — Vue re-evalúa el computed `entidadesFiltradas`
  entidadesFiltradas = entidades.filter(e =>
    (!sectorSeleccionado || e.sector === sectorSeleccionado) AND
    (!busqueda.trim() || e.entityBranchName.toLowerCase().includes(busqueda.trim().toLowerCase()))
  )

  → No se hace ninguna llamada a un sistema externo — opera enteramente sobre `entidades`,
    el arreglo ya cargado en el evento E01.

PASO 3 — Vue re-renderiza el Grid con `entidadesFiltradas`
  SI entidadesFiltradas.length === 0:
    → Muestra EmptyState "No hay entidades que coincidan con los filtros"

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión                                    | Paso afectado |
|----------------------------------------------------------|---------------|
| Ninguno — evento trivial, sin dependencias externas.      | —             |
