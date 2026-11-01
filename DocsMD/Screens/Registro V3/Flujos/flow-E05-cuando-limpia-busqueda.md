> ⛔ **No implementado (verificado en el código el 2026-10-31).** El Cabezote no tiene botón "Limpiar"; para empezar otra búsqueda se vuelve a escribir en los inputs. Se conserva el diseño original por si se retoma.

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
