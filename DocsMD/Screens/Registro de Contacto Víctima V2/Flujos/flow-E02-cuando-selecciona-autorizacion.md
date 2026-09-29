━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando selecciona la autorización de datos personales
   Tipo: User Interaction
   Función: v-model="form.autorizacion" (paso 'auth')
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  autorizacion:  valor seleccionado ('Si' | 'No')   → <select> del usuario
}

PASO 1 — Actualizar form.autorizacion con el valor seleccionado

SI form.autorizacion === 'No':
  → Vue re-evalúa el binding v-if del aviso de corte
  → Se muestra el bloque "No es posible continuar con el reporte sin la
    autorización de tratamiento de datos. Hasta aquí llega el formulario."
  → El botón "Continuar" queda deshabilitado (:disabled="form.autorizacion !== 'Si'")
  → CONTINÚA FLUJO GENERAL (el usuario permanece en el paso 'auth')

SI form.autorizacion === 'Si':
  → El aviso de corte se oculta si estaba visible
  → El botón "Continuar" se habilita
  → CONTINÚA FLUJO GENERAL

No hay llamada a backend en este evento.
