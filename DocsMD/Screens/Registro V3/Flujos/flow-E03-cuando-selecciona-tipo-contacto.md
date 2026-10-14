━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando selecciona el tipo de contacto
   Tipo: User Interaction
   Función: seleccionarTipoContacto(tipo) → components/habeas-data-gate.js
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  tipo: 'victima_caso_nuevo' | 'victima_seguimiento' | 'tercero_reporte'   → <select> TipoContactoBlock
}

PASO 1 — Enrutar según el tipo de contacto

SEGÚN tipo:
  CASO 'victima_caso_nuevo':
    → stage = 'cabezote'
    → Vue re-evalúa stage → se muestra CabezoteBusqueda
    → CONTINÚA FLUJO GENERAL

  CASO 'victima_seguimiento':
    → Redirige a la pantalla existente de Formulario de Seguimiento (dinamic-form, form_id `2d0aeb46-...`)
    → router.push('/salvia/seguimiento/nuevo')  // o la ruta real equivalente — confirmar
    → TERMINAR ejecución en esta pantalla

  CASO 'tercero_reporte':
    → Redirige al Formulario de Reporte (Tercero) — fuera de alcance de esta ronda de planeación
    → router.push('/salvia/reportes/nuevo')  // ruta por confirmar, el Formulario de Reporte no está planeado todavía
    → TERMINAR ejecución en esta pantalla

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión                                                      | Paso afectado |
|------------------------------------------------------------------------------|---------------|
| Ruta real de destino para "Víctima - Seguimiento"                            | CASO 'victima_seguimiento' |
| El Formulario de Reporte (Tercero) no está construido ni planeado todavía — confirmar si ya existe algo o si esta ronda debe dejar solo un placeholder | CASO 'tercero_reporte' |
| G-V3-11 — El `<select>` de Tipo de contacto hoy solo tiene las 3 opciones del Excel. El cliente pide verificar si hacen falta opciones adicionales para enrutar a Formulario Familiares (feminicidio) y Formulario Hombres (obs. #11) — no implementado | PASO 1 |
