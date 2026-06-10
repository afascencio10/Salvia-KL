━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando el usuario selecciona o deselecciona un caso para reasignación
   Tipo: User Interaction
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Disparado por:
  (A) cambio en el checkbox de una fila (.cc-case-checkbox)
  (B) cambio en el checkbox "Todos" del encabezado (.cc-select-all-checkbox)

Precondición: prop :reasignacion === true

INPUT (fila): {
  case:       objeto completo del caso de esa fila   → elemento de filteredCases (página actual)
  checked:    nuevo estado del checkbox               → true si se marcó, false si se desmarcó
}

INPUT (todos): {
  checked:    nuevo estado del checkbox del encabezado
  pageCases:  casos visibles en la página actual     → filteredCases
}


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND — casos-component.js
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━


PASO 1 — Validar que el modo reasignación está activo

  SI reasignacion !== true:
    → No renderizar columna de checkbox ni ejecutar este flujo
    → TERMINAR ejecución


PASO 2 — Regla de mismo equipo (caseTeam)

  Campo de comparación: case.caseTeam (victim_case_team)

  SI el usuario intenta MARCAR un caso:
    → SI selectedCases está vacío → permitir (cualquier equipo con caseTeam definido)
    → SI selectedCases NO está vacío:
        SI case.caseTeam !== selectedCases[0].caseTeam:
          → Revertir el checkbox a desmarcado
          → Mostrar alerta flotante sobre la tabla (.cc-table-alert-overlay):
            título "Selección no permitida"
            mensaje "Solo puedes seleccionar casos de un mismo equipo"
            auto-cierre a los 4.5s o botón cerrar
          → TERMINAR ejecución
        → SI coincide → CONTINÚA PASO 3

  SI case.caseTeam está vacío → no permitir selección


PASO 3 — Actualizar selectedCases (solo página actual)

  Alcance: la selección solo aplica a casos de filteredCases (página actual).
  Al cambiar de página o recargar la tabla, selectedCases se limpia.

  Checkbox de fila — toggleCaseSelection(case, event):
    SI checked === true  → agregar case (sin duplicar por i_code)
    SI checked === false → quitar case de selectedCases

  Checkbox "Todos" — toggleSelectAllCurrentPage(event):
    SI checked === false:
      → Quitar de selectedCases todos los casos de la página actual

    SI checked === true:
      → SI la página tiene casos de más de un equipo Y selectedCases está vacío:
          → Revertir checkbox, mostrar toast de mismo equipo, TERMINAR
      → Determinar targetTeam:
          SI hay selección previa → selectedCases[0].caseTeam
          SI NO → caseTeam del primer caso de la página
      → Agregar todos los casos de la página cuyo caseTeam === targetTeam
        que aún no estén en selectedCases


PASO 4 — Mostrar u ocultar botón "Reasignar Casos"

  SI selectedCases.length > 0:
    → Mostrar TableToolbar (.cc-table-toolbar) con contador y ReassignBtn

  SI selectedCases.length === 0:
    → Ocultar TableToolbar

  // No se emite evento hacia el padre en este flujo.
  // E-15 (modal de reasignación en el padre) queda pendiente de planeación.

  → FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  NOTAS DE COMPORTAMIENTO
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

| Situación | Comportamiento |
|---|---|
| Cambio de página (E-06) | selectedCases se limpia al recargar (fetchCases) |
| Cambio de filtro u ordenamiento (E-02, E-03, E-04) | selectedCases se limpia al recargar |
| reasignacion pasa de true a false | selectedCases se limpia; se oculta columna y botón |
| Página con equipos mixtos + "Todos" sin selección previa | Toast de error; no selecciona ninguno |
| Página con selección parcial del equipo A | "Todos" marca solo los casos del equipo A en la página |
| Checkbox encabezado | Estado indeterminate si hay selección parcial en la página |
