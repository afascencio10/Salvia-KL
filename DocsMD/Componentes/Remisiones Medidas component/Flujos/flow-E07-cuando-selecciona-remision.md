━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando el usuario selecciona una fila o todas
   Tipo: User Interaction
   Código: E-07
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Disparado por: checkbox de una fila, o el checkbox "Seleccionar todas" del encabezado
Solo vista supervisor. El resto de las vistas no pinta checkboxes ni la barra.


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

PASO 1 — Una fila

  Si el id no está en selectedIds, se agrega.
  Si ya está, se quita.
  No llama al API.

PASO 2 — Todas

  todosSeleccionados = hay filas visibles y cada una está en selectedIds.

  SI todosSeleccionados:
    → selectedIds = []
    → se vacía la selección entera, también ids que hubieran quedado de un filtro anterior

  SI no:
    → selectedIds = ids de las filas visibles en este momento
    → un id marcado que el filtro ocultó deja de estar seleccionado

PASO 3 — Barra

  SI selectedIds.length > 0:
    → "↺ Reasignar N seleccionada" o "seleccionadas"
    → el clic abre E-08

  La selección sobrevive a E-02 … E-06, salvo el caso de "todas" descrito arriba.
  E-10 la vacía después de guardar.

  → FIN EJECUCIÓN ✓
