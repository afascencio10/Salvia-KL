━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando selecciona el tipo de contacto
   Tipo: User Interaction
   Estado: 🟡 Implementado con cambios — sincronizado con el código 2026-10-31
   Función: seleccionarTipoContacto(tipo) → registro-v3-app.js
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  tipo:  'victima_caso_nuevo' | 'victima_seguimiento' | 'tercero_reporte' | 'familiares' | 'hombres'   → botón de la lista "Tipo de contacto"
}

SEGÚN tipo:
  CASO 'victima_caso_nuevo':
    → etapa = 'cabezote'
    → Vue muestra el Cabezote de búsqueda
    → CONTINÚA FLUJO GENERAL

  CASO 'victima_seguimiento' | 'tercero_reporte' | 'familiares' | 'hombres':
    → tipoContactoNoDisponible = true
    → Vue muestra "Este tipo de contacto se gestiona en otra pantalla del sistema, fuera de Registro V3. Cierra este formulario y continúa desde el menú principal."
    → No hay redirección ni enlace: ninguna fuente resuelve la ruta destino
    → TERMINAR ejecución (el agente puede elegir otra opción)

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión | Paso afectado |
|---|---|
| Rutas reales de destino de "Víctima - Seguimiento" (el Formulario de Seguimiento exige un `follow_up_v2` ya existente: no hay una ruta que acepte "abrir desde cero"), "Tercero - Reporte" y de los formularios de Familiares (feminicidio) y Hombres (G-V3-11) | CASO 'victima_seguimiento' · 'tercero_reporte' · 'familiares' · 'hombres' |
