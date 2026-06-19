━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando el usuario presiona un botón de acción
   Tipo: User Interaction
   Código: E-05
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Disparado por: clic en un ActionBtn dentro de la columna de acciones
               de cualquier fila de la tabla (<tr> × N)
               Solo existe si el padre pasó :buttons con al menos un elemento.

INPUT: {
  buttonId:   identificador del botón presionado   → btn.id del prop :buttons
  report:     objeto completo del reporte de esa fila → elemento de filteredReports
}


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND — reportes-component.js
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━


PASO 1 — Emitir evento hacia el componente padre

  this.$emit('action-clicked', {
    buttonId: buttonId,
    report:   report
  })

  // El componente NO interpreta buttonId.
  // NO navega, NO abre modales, NO llama al backend.
  // Toda acción posterior es responsabilidad exclusiva del padre.

  → FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  CONFIGURACIÓN DEL PADRE
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

El padre define los botones que necesite vía prop `:buttons`:

```html
<reportes-component
  :columns="reportColumns"
  :buttons="[
    { id: 'ver_detalle',  label: 'Ver detalle' },
    { id: 'crear_caso',   label: 'Crear caso' },
    { id: 'invalidar',    label: 'Invalidar' }
  ]"
  @action-clicked="onReportAction"
/>
```

El padre escucha `@action-clicked` y ejecuta la lógica según `buttonId`:

  CASO buttonId === 'ver_detalle':
    → Navegar o abrir modal con el reporte completo

  CASO buttonId === 'crear_caso':
    → Iniciar flujo de registro de caso con report.icode

  CASO buttonId === 'invalidar':
    → Llamar a InvalidateVictimContactByICode (status → 'i')


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

| Variable / decisión                                                              | Paso afectado |
|----------------------------------------------------------------------------------|---------------|
| ¿Qué botones concretos definirá la pantalla padre?                               | Configuración del padre |
| ¿Existe ya una pantalla de detalle de reporte?                                   | Responsabilidad del padre |
