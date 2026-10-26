━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando cancela el modal de reasignación
   Tipo: User Interaction
   Código: E-09
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Disparado por: botón "Cancelar" del modal


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

PASO 1 — Cerrar sin guardar

  reasignarAbierto  = false
  reasignarAgenteId = ""

  selectedIds no se toca.
  Si se abrió desde "Asignar agente", la selección sigue siendo esa única fila
  y la barra de reasignar sigue visible.
  Si se abrió desde la barra, la selección múltiple sigue igual.

  No hay POST. No se recarga el listado.

  → FIN EJECUCIÓN ✓
