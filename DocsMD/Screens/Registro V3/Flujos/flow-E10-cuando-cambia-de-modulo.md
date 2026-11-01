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

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión | Paso afectado |
|---|---|
| `triadaCompleta` solo se actualiza con la respuesta de un guardado o al abrir el Caso; el cliente no la recalcula al editar. Además, el backend responde `triadaCompleta:false` al guardar TZ/GS (retorna antes de evaluar) y `enviarSeccion` la asigna sin distinguirlo: **tras "Guardar sección" en GS (o "Finalizar llamada" estando en GS) las pestañas TZ/GS vuelven a mostrar 🔒 y el aviso verde de Iniciado desaparece**, aunque el Caso siga Iniciado en BD (se corrige al reabrirlo). Por lectura de código; no verificado en el navegador | PASO 1 |
| Planeado y no implementado: barra de progreso por módulo y regreso automático a la tríada si `triadaCompleta` pasa a false (las keys respondidas se acumulan en `victim_case_v3_answered_keys` y nunca se quitan, así que la tríada no se "des-completa") | PASO 1 |
