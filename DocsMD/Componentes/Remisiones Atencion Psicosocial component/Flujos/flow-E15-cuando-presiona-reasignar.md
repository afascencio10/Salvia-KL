━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando el usuario presiona "Reasignar"
   Tipo: User Interaction
   Código: E-15
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Precondición:
  :reasignacion === true
  selectedRemisiones.length >= 1

Disparado por: clic en ReassignBtn "Reasignar" (.rps-reassign-btn)


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

PASO 1 — Validar selección

  SI selectedRemisiones.length === 0:
    → No renderizar botón (guard defensivo)
    → TERMINAR


PASO 2 — Emitir evento hacia el padre

  $emit('reasignar-remisiones', {
    remisiones: selectedRemisiones   // objetos completos de la respuesta del listado
  })

  // El componente NO abre modal ni llama al backend.
  // El padre abre reasignar-remisiones-modal → RRM-01
  // Ver: ../reasignar-remisiones-modal-interface.md


PASO 3 — Estado post-emisión

  // No limpiar selectedRemisiones aquí — el padre puede recargar la tabla tras éxito.

  → FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  PAYLOAD EMITIDO
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

```js
{
  remisiones: [
    {
      id: "uuid",
      caseICode: "SAL-001",
      status: "abierto",
      duplaId: "...",
      agentId: "...",
      // ... resto de campos del listado E-01
    },
    ...
  ]
}
```

**Responsabilidad del padre:**

```javascript
onReasignarRemisiones(payload) {
  this.$refs.reasignarModal.open(payload.remisiones);  // RRM-01
}
```

Documentación modal: [reasignar-remisiones-modal-interface.md](../reasignar-remisiones-modal-interface.md)
