━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando cierra el modal
   Tipo: User Interaction
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Disparado por: (A) clic en "✕" del header
               (B) clic en el backdrop
               (C) clic en "Cerrar", step 'cerrado'

Precondición: saving !== true  (no se puede cerrar a mitad de un guardado)

INPUT: {} — no recibe datos


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND — gestion-reporte-modal.js
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━


PASO 1 — Validar precondición

  SI saving === true:
    → TERMINAR ejecución (no hace nada mientras hay un guardado en curso)


PASO 2 — Limpiar estado interno y ocultar

  visible = false
  step = null
  report = null
  intentos = []
  pendingClose = ''; motivoTexto = ''; atencionEspecializada = 'No'


PASO 3 — Avisar al padre

  this.$emit('updated')

  // Siempre se emite, sin distinguir si hubo un cambio real (se cerró el
  // reporte, se registró un intento) o el usuario solo miró y cerró sin
  // tocar nada. Recargar casos-component de más es barato; no recargar
  // cuando sí hubo cambio dejaría la tabla desactualizada — se prefiere
  // la opción simple y siempre-correcta sobre rastrear un flag "dirty".

  → FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Ninguno.
