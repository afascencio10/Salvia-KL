━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando carga el módulo «Ruta de atención integral»
   Tipo: Lifecycle
   Estado: ✅ Implementado 2026-11-04 (WP-4 · REG-17)
   Funciones: irAModulo('RA') · watch de moduloActivo → registro-v3-app.js
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

> No hay petición propia: el Caso ya se cargó con `abrirCaso()` (`GET /api/v1/victim-cases-v3/{icode}`), que trae `answers` (con `GS_14` y `GS_17` si ya se respondieron), `entities`, `barriers`, `entityMeta`, `barrierMeta` y `entitiesReadOnly`.

INPUT: { moduloActivo: 'RA' }   → clic en la pestaña «Ruta de atención integral» (no tiene candado: `moduloBloqueado('RA')` es `false`)

PASO 1 — `irAModulo('RA')`: limpia los mensajes y fija `moduloActivo = 'RA'`
PASO 2 — `watch` de `moduloActivo`: si el formulario de entidad no es de tipo `acudida`, se reinicia (`entidadVacia('acudida')`), junto con el buscador, los errores y el formulario de barrera abierto
PASO 3 — Se pinta el bloque:
  aviso de borrador (`!entitiesReadOnly`) o de solo lectura
  `GS_14`, y si es «Sí» (o si `GS_17` es «Sí» y no hay entidades acudidas): las tarjetas de `entidadesAcudidas` y, si `!entitiesReadOnly`, el formulario de entidad acudida
  `GS_17` y su nota

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión | Paso afectado |
|---|---|
| `guardarYSalir` y `beforeunload` (de WP-3) no guardan desde este módulo: lo contestado en RA se pierde si se sale sin pulsar «Guardar sección» («Para integración») | PASO 1 |
