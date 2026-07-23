━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando selecciona un sector en el filtro
   Tipo: User Interaction
   Función: v-model="sectorSeleccionado"  (computed: entidadesFiltradas)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  sector:  valor elegido en el <select>   → "" (Todos) | "he" | "js" | "pm"
}

PASO 1 — Actualizar variable reactiva `sectorSeleccionado`

PASO 2 — Vue re-evalúa el computed `entidadesFiltradas` (misma fórmula que en E02,
  combinando sector + búsqueda de texto) — sobre los datos ya cargados en E01.

SI sectorSeleccionado === "" (opción "Todos"):
  → No aplica filtro de sector
  → CONTINÚA FLUJO GENERAL

SI NO (sector específico elegido):
  → Filtra entidades cuyo `.sector` coincida
  → CONTINÚA FLUJO GENERAL

PASO 3 — Vue re-renderiza el Grid con el resultado

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión                                    | Paso afectado |
|----------------------------------------------------------|---------------|
| Ninguno — evento trivial, sin dependencias externas.      | —             |
