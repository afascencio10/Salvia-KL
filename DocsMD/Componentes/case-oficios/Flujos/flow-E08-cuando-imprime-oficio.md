━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando imprime el oficio
   Tipo: User Interaction
   Función: printOficio()
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  (ninguno — imprime la ventana actual, con el modal de detalle abierto)
}


PASO 1 — Invocar el diálogo de impresión del navegador

  window.print()
  // Imprime toda la ventana visible tal como está en el DOM en ese momento


→ FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión                                                    | Paso afectado |
|--------------------------------------------------------------------------|---------------|
| window.print() imprime toda la página (incluyendo el resto de           | PASO 1        |
| Detalle del Caso detrás del overlay), no solo el modal. ¿Hace falta      |               |
| un estilo @media print que oculte todo excepto .co-modal, como          |               |
| probablemente ya tenga oficios-list.js (a verificar en su CSS)?         |               |
