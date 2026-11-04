━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando responde «¿Se identificaron barreras?»
   Tipo: User Interaction
   Estado: ✅ Implementado 2026-11-04 (WP-4 · REG-17)
   Funciones: setBool('GS_17', valor) → registro-v3-app.js
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: { GS_17: true | false }   → radios «Sí / No»

PASO 1 — `answers['GS_17'] = 'true' | 'false'` (se guarda con «Guardar sección», RA-E08)
PASO 2 — Reactividad:
  SI es «Sí»: cada tarjeta (de esta lista y del Plan de acción) ofrece «+ Agregar barrera a esta entidad» (si el Caso sigue en Borrador) y se muestra la nota «Registra cada barrera con el botón… de la tarjeta de la entidad donde ocurrió.»; si no hay ninguna entidad, se agrega «Primero agrega la entidad.»; si no hay entidades acudidas, el formulario de entidad se muestra aunque `GS_14` no sea «Sí»
  SI es «No»: se ocultan los botones; las barreras ya registradas siguen en las tablas

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión | Paso afectado |
|---|---|
| Sin mínimo de barreras: «Sí» no obliga a registrar ninguna | PASO 2 |
