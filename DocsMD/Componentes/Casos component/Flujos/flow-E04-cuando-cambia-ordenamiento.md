━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando el usuario cambia el ordenamiento
   Tipo: User Interaction
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Disparado por: cambio en el <select> SortSelector de FilterBar → SearchAndSort

INPUT: {
  newSortBy:   criterio de orden seleccionado   → valor del <select>
               valores: 'registration_date' | 'next_follow_up'
}


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND — casos-component.js  (ordenamiento del lado del cliente)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━


PASO 1 — Determinar dirección del ordenamiento

  SI newSortBy === sortBy (mismo criterio que el activo):
    → sortOrder = sortOrder === 'asc' ? 'desc' : 'asc'   // invertir dirección

  SI newSortBy !== sortBy (criterio diferente):
    → sortBy    = newSortBy
    → sortOrder = 'desc'    // reiniciar a descendente por defecto


PASO 2 — Reordenar filteredCases del lado del cliente

  SEGÚN sortBy:

    CASO 'registration_date':
      → Ordenar filteredCases por case.creationDate  {sortOrder}
        SI sortOrder === 'asc':  más antiguo primero
        SI sortOrder === 'desc': más reciente primero

    CASO 'next_follow_up':
      → Ordenar filteredCases por case.nextFollowUpDate  {sortOrder}
        SI case.nextFollowUpDate === null:
          → Siempre al final (NULLS LAST), independiente del sortOrder


PASO 3 — Actualizar la vista

  → Renderizar tabla con el nuevo orden de filteredCases

  → FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

| Variable / decisión                                                              | Paso afectado |
|----------------------------------------------------------------------------------|---------------|
| ¿El ordenamiento inicial por defecto es 'desc' de fecha de registro?             | PASO 1        |
| ¿Se muestra algún indicador visual (flecha ↑↓) en el selector activo?           | Interfaz      |
| ¿El orden 'asc' de next_follow_up pone primero los más urgentes (más próximos)? | PASO 2        |
