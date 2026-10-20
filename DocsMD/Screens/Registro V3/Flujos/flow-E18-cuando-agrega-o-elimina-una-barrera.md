━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando agrega o elimina una Barrera
   Tipo: User Interaction
   Función: agregarBarrera() / eliminarBarrera(index) → components/dinamic-form.js (repeater genérico)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

> Modela el repeater real "Identificación de Barreras" de Seguimiento
> (`repeater_group` `5fd3ecdc-2e5f-4b31-97ef-8a994580586a`, `min_repetitions = 1`
> en Seguimiento) — ver `registro-v3-data.md` § D9.

INPUT: {
  index: posición del ítem a eliminar, solo aplica para eliminarBarrera
}

PASO 1 — Agregar una barrera

SI se presiona "+ Siguiente barrera":
  → Crea un nuevo ítem vacío en `respuestas.barreras`
  → Renderiza sus campos (Sector de la barrera, bloque de barreras del sector elegido, entidad asignada, categorías transversales, fecha, funcionario/dependencia, descripción, gestión)
  → CONTINÚA FLUJO GENERAL

PASO 2 — Eliminar una barrera

SI se presiona "✕ Eliminar" sobre un ítem:
  → Mostrar confirmación "¿Eliminar esta barrera del registro?"
  SI confirma:
    → Quita el ítem en `index`
    → Recalcula la barra de progreso del módulo GS
    → CONTINÚA FLUJO GENERAL
  SI NO confirma:
    → TERMINAR ejecución

PASO 3 — Recalcular visibilidad del bloque de barreras según sector

Cuando cambia "Sector de la barrera" de un ítem:
  SEGÚN valor:
    CASO 'salud':       → Muestra "Barreras identificadas en Salud" + "Otra barrera en Salud ¿cuál?" (condicional a incluir "Otra"); oculta los bloques de Justicia y Protección
    CASO 'justicia':    → Muestra el bloque de Justicia; oculta Salud y Protección
    CASO 'proteccion':  → Muestra el bloque de Protección; oculta Salud y Justicia
    CASO 'otras_instituciones': → Oculta los 3 bloques sectoriales; queda disponible el campo "Entidad a la que se asignará la barrera" y las 4 categorías transversales
  → Las 4 categorías transversales (institucionales/talento humano, económicas, territoriales, enfoque diferencial) están siempre visibles sin importar el sector

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión                                                                     | Paso afectado |
|-----------------------------------------------------------------------------------------------|---------------|
| G-V3-20 — catálogo exacto de opciones de cada multiselect (barreras por sector, categorías transversales, gestión de la barrera) — no extraído de la tabla `option` real en este análisis | PASO 1, PASO 3 |
| G-V3-21 — tabla real donde se persisten las respuestas del repeater | PASO 1 |
| Si Registro V3 debe heredar `min_repetitions = 1` de Seguimiento (obliga a registrar al menos 1 barrera una vez que Q17 = Sí) o dejarlo en 0 (puede responder Sí y no alcanzar a registrar ninguna en este contacto) | PASO 1 |
