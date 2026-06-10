━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando el usuario presiona "Reasignar Casos"
   Tipo: User Interaction
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Disparado por: clic en ReassignBtn ("Reasignar Casos")
               en TableToolbar (.cc-table-toolbar), encima de la tabla

Precondición: prop :reasignacion === true
               selectedCases.length > 0

INPUT: {
  cases:   casos seleccionados   → selectedCases (Array<Object>)
}


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND — casos-component.js
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━


PASO 1 — Validar precondiciones

  SI reasignacion !== true || selectedCases.length === 0:
    → El botón no debería estar visible; no hacer nada
    → TERMINAR ejecución


PASO 2 — Emitir evento hacia el componente padre

  this.$emit('reasignar-casos', {
    cases: selectedCases   // array de objetos caso completos, mismo shape que E-05
  })

  // El componente padre es quien interpreta el evento y decide qué acción ejecutar.
  // casos-component no abre modales, no reasigna ni hace llamadas al backend.

  → FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  RESPONSABILIDAD DEL PADRE
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

El componente padre escucha @reasignar-casos y ejecuta la lógica de reasignación.

Ejemplo de uso típico:

  @reasignar-casos="onReasignarCasos"

  onReasignarCasos({ cases }):
    → this.$refs.reasignarModal.open(cases)
    → Ver planeación del modal:
        reasignar-casos-modal-interface.md
        reasignar-casos-modal-events.md (M-01 … M-05)
    → Tras reasignación exitosa (M-05, pendiente), el padre recarga la tabla


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

| Variable / decisión                                                              | Paso afectado |
|----------------------------------------------------------------------------------|---------------|
| ¿El padre limpia selectedCases tras reasignar exitosa o lo hace el hijo?         | Responsabilidad padre / ref |
| Planeación del guardado (M-05)                                                   | [flow-M05](./flow-M05-cuando-confirma-reasignacion.md) |
| Al cerrar el modal exitosamente, el padre debe recargar la tabla                 | Responsabilidad padre + M-05 |
