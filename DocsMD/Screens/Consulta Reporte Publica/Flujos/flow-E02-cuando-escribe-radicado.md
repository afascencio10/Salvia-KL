━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando escribe el radicado
   Tipo: User Interaction
   Estado: ✅ Implementado y probado — WP-1 T4 (2026-11-02)
   Función: onRadicadoInput()
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: { valor: texto del campo radicado → lo que la persona escribe o pega }

PASO 1 — Normalizar
  limpio = valor.toUpperCase() sin todo carácter que no sea A-Z o 0-9
  SI limpio tiene más de 8 símbolos → recortar a 8
  SI limpio tiene más de 4 símbolos → insertar "-" tras el 4.º   → "K7QM-4R2X"
  radicado = resultado
PASO 2 — SI mensajeError ≠ '' → mensajeError = ''
PASO 3 — Validar el formato (reactivo)
  radicadoValido = /^[A-HJ-NP-Z2-9]{4}-[A-HJ-NP-Z2-9]{4}$/.test(radicado)
  → :disabled de «Consultar» = !radicadoValido || consultando || bloqueado

No hay llamada a backend. (I, O, 0 y 1 se aceptan al escribir pero el botón sigue deshabilitado: no existen en el alfabeto.)

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Ninguno.
