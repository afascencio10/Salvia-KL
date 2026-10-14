━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando presiona "Limpiar" en el Cabezote
   Tipo: User Interaction
   Función: limpiarBusqueda() → components/cabezote-busqueda.js
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {} (no recibe datos — es un reset)

PASO 1 — Resetear criterios y resultados
  busqueda = { documento: '', telefono: '' }
  resultados = { reportes: [], registros: [], borradores: [] }
  buscado = false

PASO 2 — Renderizar
  Vue re-evalúa buscado → se oculta la sección de resultados, se muestran los 3 inputs vacíos

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Ninguno — evento trivial.
