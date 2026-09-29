━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando selecciona el ajuste razonable
   Tipo: User Interaction
   Función: v-model="form.ajusteRazonable" (paso 'perfil')
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  ajusteRazonable:  valor seleccionado (10 opciones — Guía-intérprete, Traducción,
                     Apoyos, Formatos accesibles, Física, Dispositivos, Confianza,
                     Transporte, No informa, No requiere)  → <select> del usuario
}

PASO 1 — Actualizar form.ajusteRazonable con el valor seleccionado

PASO 2 — Habilitar/deshabilitar "Continuar" del paso 'perfil'
  según :disabled="!form.tipoReporte || !form.ajusteRazonable" (ver también E-04)

No hay llamada a backend ni efecto sobre otros campos del formulario.
