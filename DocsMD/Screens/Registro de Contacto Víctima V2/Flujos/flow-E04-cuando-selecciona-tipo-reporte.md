━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando selecciona el tipo de reporte
   Tipo: User Interaction
   Función: v-model="form.tipoReporte" (paso 'perfil')
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  tipoReporte:  valor seleccionado ('Victima' | 'Tercero' | 'Familiar' | 'Hombre')  → <select> del usuario
}

PASO 1 — Actualizar form.tipoReporte con el valor seleccionado

PASO 2 — Vue recalcula qué bloques del paso 'formulario' serán visibles
  (efecto reactivo del mismo evento — no se dispara al llegar al paso
  'formulario', sino que queda determinado desde este momento):

  SEGÚN tipoReporte:
    CASO 'Victima':
      → Al llegar a 'formulario' se ocultará "Datos de quien reporta"
      → Se mostrará "Datos de la persona afectada" con teléfono y disposición
        de llamada de la víctima
      → CONTINÚA

    CASO 'Tercero':
      → Se mostrará "Datos de quien reporta" con selector de anonimato (E-08),
        rol de cuidado y "¿la víctima tiene conocimiento?"
      → Se mostrará "Datos de la persona afectada" con teléfono y disposición
        de llamada de la víctima
      → CONTINÚA

    CASO 'Familiar':
      → Se mostrará "Datos de quien reporta" con parentesco, teléfono y
        disposición de llamada DEL REPORTANTE (no de la víctima)
      → Se mostrará "Datos de la persona afectada" SIN teléfono ni disposición
        (nota informativa: la víctima falleció)
      → CONTINÚA

    CASO 'Hombre':
      → Se mostrará "Datos de quien reporta" con teléfono y disposición de
        llamada DEL REPORTANTE
      → Se OCULTA por completo el bloque "Datos de la persona afectada"
      → CONTINÚA

  DEFAULT ('' — sin seleccionar):
    → Ningún bloque del paso 'formulario' es alcanzable aún
    → El botón "Continuar" del paso 'perfil' permanece deshabilitado

PASO 3 — Habilitar/deshabilitar "Continuar" del paso 'perfil'
  según :disabled="!form.tipoReporte || !form.ajusteRazonable" (ver también E-05)

No hay llamada a backend en este evento.
