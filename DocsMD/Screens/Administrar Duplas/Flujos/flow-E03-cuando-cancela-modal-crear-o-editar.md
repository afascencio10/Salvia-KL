# flow-E03 — Cuando cancela modal crear o editar

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando cancela modal crear o editar
   Tipo: User Interaction
   Función: closeFormModal()
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  modal.visible:  true
  modal.kind:     'form'
}

Disparado por:
  • Botón "Cancelar"
  • ✕ del encabezado (si existe)
  • Clic en backdrop / Escape (si se habilita)


PASO 1 — Verificar que se puede cerrar

SI isSaving === true:
  → No cerrar (evitar perder request en vuelo)
  → TERMINAR

SI isSaving === false:
  → CONTINÚA FLUJO GENERAL


PASO 2 — Cerrar y limpiar estado del formulario

  modal.visible = false
  modal.kind = null
  modal.mode = null
  editingDuplaId = null
  form = { name: '', psychologistId: '', socialWorkerId: '' }
  saveError = null
  warningPs = null
  warningTs = null

  → Vue oculta el modal
  → La lista de profesionales / duplas NO se recarga (no hubo cambios)

→ FIN EJECUCIÓN ✓
```

---

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

| Variable / decisión | Paso afectado |
|---|---|
| ¿Confirmar “¿Descartar cambios?” si el form está dirty? | PASO 1 |
| ¿Cierre con Escape / backdrop habilitado? | Disparadores |
