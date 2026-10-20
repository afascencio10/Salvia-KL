━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando presiona "Reasignar Reportes"
   Tipo: User Interaction
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Disparado por: clic en ReassignBtn ("Reasignar Reportes" — el texto sale del
               prop item-label-plural) dentro de casos-component

Precondición: :reasignacion === true (solo rol sv) Y selectedCases.length > 0

INPUT: {
  reports:  reportes seleccionados   → selectedCases (Array<ReportListItem>)
}


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND — casos-component.js (sin cambios — ya documentado como E-15)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━


PASO 1 — Emitir hacia el padre

  this.$emit('reasignar-casos', { cases: selectedCases })

  // Mismo evento genérico ya usado por casos. El payload se llama 'cases'
  // por consistencia con la firma existente, aunque contenga reportes.


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND — list_reportes.html (padre)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━


PASO 2 — onReasignarReportes({ cases })

  this.$refs.reasignarModal.open(cases)   // el modal ya está montado con entity-kind="reportes" (prop), solo para sv

  // reasignar-casos-modal extendido (ver related-tables.md GAP-3): no exige
  // equipo común entre los reportes seleccionados — a diferencia de casos,
  // un reporte no tiene equipo hasta que se le asigna un agente. Hoy eso lo
  // garantiza el modo contingencia (REASSIGN_CROSS_TEAM_CONTINGENCY = true).

  → FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — resueltos tras la implementación
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

| Variable / decisión                                                              | Resolución |
|------------------------------------------------------------------------------------|------------|
| Texto del botón: ¿"Reasignar Casos" (genérico) o parametrizable? | Parametrizado con el prop `item-label-plural` de casos-component: "Reasignar Reportes" (caso `LR-E05-C2`). Con el default `'casos'` las pantallas de casos no cambian |
| Regla de equipo común al seleccionar (E-14 de casos-component) — para reportes, ¿se desactiva o se compara `victim_contact_team`? | **No se adaptó la regla: la desactiva el modo contingencia** (`REASSIGN_CROSS_TEAM_CONTINGENCY = true`, `casos-component.js:43`). Con el flag en `false`, la regla lee `caseTeam` —el reporte trae `team`— y no se podría marcar ninguna fila. Riesgo anotado en `casos-component-interface.md` (casos `LR-E05-C3`) |
