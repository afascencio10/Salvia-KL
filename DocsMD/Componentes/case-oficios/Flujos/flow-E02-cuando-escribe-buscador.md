━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando escribe en el buscador
   Tipo: User Interaction
   Función: reactividad de Vue sobre `buscador` (v-model, sin método explícito)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  texto:  contenido del <input> del buscador   → tecleado por el usuario
}


PASO 1 — Actualizar el estado del buscador

  buscador = texto   // v-model, sin debounce (ver GAP)


PASO 2 — Vue recalcula oficiosFiltrados (computed)

  oficiosFiltrados = oficios.filter(o => cumpleBuscador(o) && cumpleTema(o) && cumpleEstado(o))

  cumpleBuscador(o):
    SI buscador.trim() === '':
      → true  (sin filtro activo)
    SI NO:
      q = buscador.toLowerCase()
      asunto = o.subject || o.asuntoRadicado || o.asuntoRespuesta || ''
      → (o.entidad || '').toLowerCase().includes(q)
        OR (o.urlKofax || '').toLowerCase().includes(q)
        OR asunto.toLowerCase().includes(q)


PASO 3 — Vue re-renderiza el grid

  SI oficiosFiltrados.length === 0:
    → Muestra EmptyState "No hay oficios que coincidan con los filtros."
  SI NO:
    → Renderiza las OficioCard resultantes

→ FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión                                                    | Paso afectado |
|--------------------------------------------------------------------------|---------------|
| ¿Necesita debounce? Con filtrado 100% en memoria sobre una lista ya      | PASO 1        |
| cargada (sin llamadas al API por cada tecla), probablemente no —        |               |
| queda a criterio de UX si el volumen de oficios por caso crece mucho.   |               |
