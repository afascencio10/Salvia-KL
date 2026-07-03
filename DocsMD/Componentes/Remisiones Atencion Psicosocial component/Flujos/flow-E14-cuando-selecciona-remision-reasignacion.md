━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando el usuario selecciona o deselecciona una remisión para reasignación
   Tipo: User Interaction
   Código: E-14
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Precondición: prop :reasignacion === true

Disparado por:
  (A) checkbox de fila
  (B) checkbox "Todos" del encabezado


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

PASO 1 — Validar modo reasignación

  SI reasignacion !== true → TERMINAR


PASO 2 — Regla de selección

  SI el usuario intenta MARCAR una remisión:
    SI status === 'cerrado':
      → Checkbox deshabilitado / revertir si aplica
      → Alerta: "No puedes reasignar remisiones en estado cerrado"
      → TERMINAR
    SI NO → permitir (cualquier combinación de estados)


PASO 3 — Actualizar selectedRemisiones (solo página actual)

  Checkbox fila:
    SI checked → agregar (sin duplicar por id)
    SI unchecked → quitar

  Checkbox "Todos":
    SI unchecked → quitar todas las reasignables de la página actual
    SI checked → agregar todas las reasignables de la página (status <> cerrado)

PASO 4 — Mostrar u ocultar botón "Reasignar"

  SI selectedRemisiones.length > 0:
    → Mostrar TableToolbar con ReassignBtn

  SI selectedRemisiones.length === 0:
    → Ocultar TableToolbar

  // No emite evento al padre — E-15 lo hace al presionar el botón

  → FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  LIMPIEZA DE SELECCIÓN
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

| Evento | Acción |
|---|---|
| E-05 Limpiar filtros | selectedRemisiones = [] |
| E-06 Cambio página | selectedRemisiones = [] |
| Cualquier filtro E-03…E-13 | selectedRemisiones = [] |
| fetchRemisiones() | selectedRemisiones = [] |
| reload() | selectedRemisiones = [] |
