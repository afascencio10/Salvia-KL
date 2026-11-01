━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando presiona "Atrás" en Perfil
   Tipo: User Interaction
   Función: goTo('auth') — botón del paso 'perfil'
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {}

PASO 1 — Navegar al paso anterior
  this.step = 'auth'   // goTo('auth')

No hay validación, ni llamada a backend, ni limpieza de form. tipoReporte y
ajusteRazonable quedan en memoria por si el usuario vuelve a avanzar.
