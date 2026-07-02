━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando carga el componente
   Tipo: Lifecycle
   Función: mounted()
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  barrierId:   UUID de la barrera    → prop recibido del padre
}

PASO 1 — Activar estado de carga
  → cargando = true
  → error = null
  → seguimientos = []

PASO 2 — Consultar seguimientos de la barrera
  GET /api/v1/barriers/:barrierId/follow-ups

  → resultado: array de BarrierFollowUp ordenado por createdAt ASC

SI error de red o respuesta no-200:
  → error = mensaje de error recibido
  → cargando = false
  → TERMINAR ejecución

SI NO:
  → CONTINÚA FLUJO GENERAL

PASO 3 — Asignar resultados
  → seguimientos = respuesta del API
  → cargando = false

PASO 4 — Vue re-evalúa el template
  SI seguimientos.length === 0:
    → Renderiza EmptyState "No hay seguimientos registrados para esta barrera."

  SI NO:
    → Renderiza Timeline con un FollowUpItem por cada elemento
    → Cada item muestra: fecha, autor, tag persiste/cerrada, respuesta institucional,
      chips de gestión (si los hay), actuaciones (si las hay), motivo de cierre (si aplica)

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión                              | Paso afectado |
|--------------------------------------------------|---------------|
| El backend debe resolver actorName desde created_by_id — definir si lo hace en el endpoint o si el modelo BarrierFollowUp expone un campo calculado | PASO 2 |
| Endpoint `GET /api/v1/barriers/:barrierId/follow-ups` aún no existe — debe crearse | PASO 2 |
