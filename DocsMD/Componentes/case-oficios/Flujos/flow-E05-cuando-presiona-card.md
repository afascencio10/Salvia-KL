━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando presiona una card de oficio
   Tipo: User Interaction
   Función: abrirDetalle(oficio)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  oficio:  objeto EntityLetter ya presente en memoria (oficios[])   → de la card presionada
}


PASO 1 — Mostrar el modal con el oficio seleccionado

  oficioSeleccionado = oficio
  // Sin fetch adicional — el objeto completo ya está cargado desde E01


PASO 2 — Vue re-evalúa reactivamente

  → El modal (.co-modal-overlay) aparece con v-if="oficioSeleccionado"
  → Se renderizan las secciones según los campos presentes en el oficio:
      Oficio (siempre) · Radicación y Respuesta (si aplica) ·
      Corrección (si reasonCorrection) · Barrera relacionada (si barrierId) ·
      Caso relacionado (si caseId — en este componente siempre está presente,
      ya que la lista se cargó filtrada por caseId)

→ FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión                                                    | Paso afectado |
|--------------------------------------------------------------------------|---------------|
| Como este componente vive dentro de Detalle del Caso, la sección "Caso   | PASO 2        |
| relacionado" del modal es redundante (el usuario ya está en ese caso).  |               |
| ¿Se omite esa sección aquí, a diferencia de oficios-list.js donde sí     |               |
| aporta (porque ahí se puede ver oficios de distintos casos)?            |               |
