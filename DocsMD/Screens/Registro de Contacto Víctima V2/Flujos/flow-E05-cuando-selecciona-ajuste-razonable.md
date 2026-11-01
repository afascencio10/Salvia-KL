━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando selecciona el ajuste razonable
   Tipo: User Interaction
   Función: v-model="form.ajusteRazonable" (paso 'perfil')
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  ajusteRazonable:  objeto enum {icode, code, name}; 12 opciones (sc, gi, id, ap, in, ad, dt, pr, tr, ni, nr, ns —
                     textos en la tabla de etiquetas de registro-contacto-victima-v2-interface.md)  → <select> del usuario
}

PASO 1 — Actualizar form.ajusteRazonable con el valor seleccionado

PASO 2 — Habilitar/deshabilitar "Continuar" del paso 'perfil'
  según :disabled="!form.tipoReporte || !form.ajusteRazonable" (ver también E-04)

No hay llamada a backend ni efecto sobre otros campos del formulario.
