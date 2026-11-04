━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando carga la pantalla con un registro de origen
   Tipo: Lifecycle                                   (numeración E19 provisional: puede chocar con la de otros WP al integrar)
   Estado: ✅ Implementado y probado — WP-1 T7 (2026-11-02) (ver QA/E19-carga-con-reporte-qa-cases.md)
   Funciones: mounted() → abrirDesdeReporte() → abrirCaso(icode)  (registro-v3-app.js)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: { ?caso=<icode> en la URL    → generada por el «Registro de caso» del modal de gestión de reportes (redirect) o por «Ir al borrador del caso» }

PASO 1 — mounted() (existente): cargarEnumCatalog(), cargarGeoDepartamentos() y el respaldo beforeunload
         🆕 al final: this.abrirDesdeReporte()                 // única línea nueva en código existente

PASO 2 — 🆕 abrirDesdeReporte() (método nuevo)
  icode = new URLSearchParams(location.search).get('caso')
  SI vacío o no tiene forma de UUID: no hace nada — el flujo normal arranca en la etapa 'habeas' → TERMINAR ejecución
  SI NO → CONTINÚA FLUJO GENERAL

PASO 3 — Abrir el registro creado desde el reporte
  → abrirCaso(icode)                                     // existente, sin cambios:
      GET /api/v1/victim-cases-v3/:icode
      answers = solo las claves en answeredKeys (CI_2, CI_4, HV_1 y, si aplica, CI_8)
      etapa = 'formulario' ; moduloActivo = 'CI'
  Se omiten habeas, tipo de contacto y Cabezote: el consentimiento ya se tomó en el reporte V2.

PASO 4 — Recargar la página con el mismo ?caso reabre el mismo registro (el parámetro no se limpia).

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión | Paso afectado |
|---|---|
| Aviso visual «Registro creado desde el reporte {radicado}» (propuesta: no; el timeline del caso ya muestra «Borrador de caso iniciado») | PASO 3 |
| `abrirCaso` no maneja errores de red: si el GET falla la persona queda en la etapa habeas | PASO 3 |
| **Hallazgo para el equipo de V3 (mitigado en la conversión):** `fanOutPerfilV3` solo omite la réplica del Perfil cuando el documento está VACÍO, pero `CreateDraft` guarda el temporal `cc/0000000000`: con dos borradores sin documento, guardar CI/DP/CA en uno pisaba las respuestas del otro (probado: el teléfono). Los borradores creados desde un reporte quedan con el documento vacío, así que no les pasa. Los borradores normales sin documento (p. ej. creados con búsqueda solo por teléfono) siguen expuestos; corrección sugerida (1 línea, archivo de V3 que WP-1 no toca): en `fanOutPerfilV3`, `if doc.DocType == "" \|\| doc.DocNumber == "" \|\| doc.DocNumber == dummyDocNumberV3 { return nil }` | PASO 3 |
