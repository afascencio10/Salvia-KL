━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando el usuario presiona un botón de acción
   Tipo: User Interaction
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Disparado por: clic en un ActionBtn dentro de la columna de acciones
               de cualquier fila de la tabla (<tr> × N)

INPUT: {
  buttonId:   identificador del botón presionado   → btn.id del prop :buttons
  case:       objeto completo del caso de esa fila → elemento de filteredCases
}


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND — casos-component.js
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━


PASO 1 — Emitir evento hacia el componente padre

  this.$emit('action-clicked', {
    buttonId: buttonId,   // id del botón presionado, definido en prop :buttons
    case:     case        // objeto completo del caso:
                          //   { id, i_code, names, lastNames, docNumber,
                          //     creationDate, status, agentNames, agentLastNames,
                          //     team, riskStatus, nextFollowUpDate }
  })

  // El componente padre es quien interpreta buttonId y decide qué acción ejecutar.
  // casos-component no navega, no abre modales ni hace llamadas al backend.

  → FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  RESPONSABILIDAD DEL PADRE
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

El componente padre escucha @action-clicked y ejecuta la lógica según buttonId.

Ejemplos de uso típico:

  CASO buttonId === 'ver_detalle':
    → Navegar a /salvia/casos/{case.i_code}/detalle

  CASO buttonId === 'hacer_seguimiento':
    → Navegar a /salvia/hacer-seguimiento/{followUpId}
      (el padre resuelve el followUpId a partir de case.i_code)

  CASO buttonId === 'asignar':
    → Abrir modal de asignación con case como contexto


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

| Variable / decisión                                                              | Paso afectado |
|----------------------------------------------------------------------------------|---------------|
| IDs exactos de los botones que usará cada pantalla padre (ver_detalle, etc.)     | Responsabilidad padre |
| ¿El botón puede tener estado (disabled) según condición del caso?                | Interfaz      |
| ¿Los botones tienen iconos además del label?                                     | Interfaz      |
