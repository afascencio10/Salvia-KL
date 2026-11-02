━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando cambia de módulo en el sidebar
   Tipo: User Interaction
   Estado: 🟡 Implementado con cambios — sincronizado con el código 2026-10-31
   Funciones: irAModulo(id) · moduloBloqueado(id) → registro-v3-app.js
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

> El "sidebar" planeado es hoy una fila de pestañas (`.nav-tabs`) con los 7 módulos, sin barras de progreso.

INPUT: {
  id:  'CI' | 'DP' | 'CA' | 'HV' | 'AG' | 'TZ' | 'GS'   → clic en la pestaña
}

PASO 1 — Verificar si el módulo está bloqueado
  moduloBloqueado(id) = (id === 'TZ' || id === 'GS') && !triadaCompleta

SI moduloBloqueado(id):
  → mensajeModuloBloqueado = "Completa primero los módulos de Contacto, Datos personales, Contacto de apoyo, Hechos y Agresor para habilitar Tamizaje|Gestión Salvia."
  → No cambia de módulo (la pestaña bloqueada se ve atenuada, con 🔒 y title "Completa primero CI/DP/CA/HV/AG")
  → TERMINAR ejecución

SI NO:
  → mensajeModuloBloqueado = ''
  → CONTINÚA FLUJO GENERAL

PASO 2 — Cambiar de módulo
  moduloActivo = id; mensajeGuardado = ''
  Vue muestra el bloque del módulo: ModuloGenerico (CI/DP/CA/HV/AG), Tamizaje (TZ) o GestionSalvia (GS)
  Las respuestas no guardadas del módulo anterior se conservan en `answers` (no hay guardado automático al cambiar)

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━CAMBIOS WP-3 (2026-11-02):
  PASO 1 (REG-13): moduloBloqueado(id) = (id === 'TZ') && !triadaCompleta. Gestión Salvia ya no se bloquea. El aviso queda solo para el Tamizaje:
    "Completa primero los módulos de Contacto, Datos personales, Contacto de apoyo, Hechos y Agresor para habilitar Tamizaje."
  PASO 1 (REG-35): el cliente nunca baja triadaCompleta de true a false y el servidor devuelve la tríada real al guardar cualquier módulo:
    guardar GS o TZ ya no vuelve a poner los candados (se cierra el GAP G-V3-24).
  PASO 2 (REG-01): los ítems del menú vertical (uno por módulo de modulos[]) llaman a irAModulo; también lo llaman "Siguiente", "Anterior" y el avance
    automático tras guardar. moduloSiguiente/moduloAnterior(id) recorren modulos[] y devuelven el primer módulo NO bloqueado (o null).
    Se limpian mensajeGuardado, mensajeGuardadoEsError, mensajeAvance y errorGuardado.
  Reactivo (REG-02): el menú muestra el porcentaje de cada módulo y el encabezado el global; se recalculan solos al cambiar las respuestas.

━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión | Paso afectado |
|---|---|
| Planeado y no implementado: barra de progreso por módulo y regreso automático a la tríada si `triadaCompleta` pasa a false (las keys respondidas se acumulan en `victim_case_v3_answered_keys` y nunca se quitan, así que la tríada no se "des-completa") | PASO 1 |
