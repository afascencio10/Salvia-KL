━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando presiona "Continuar" en Autorización
   Tipo: User Interaction
   Función: goTo('perfil') — botón del paso 'auth'
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {}  (solo requiere que form.autorizacion.code === 'y', ya validado por :disabled)

PASO 1 — Navegar al paso siguiente
  this.step = 'perfil'   // goTo('perfil')

No hay validación adicional, ni llamada a backend, ni transformación de datos.
Los valores de form ya diligenciados se conservan.
