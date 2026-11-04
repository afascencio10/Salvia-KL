━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando responde «¿Ha puesto en conocimiento de alguna entidad…?»
   Tipo: User Interaction
   Estado: ✅ Implementado 2026-11-04 (WP-4 · REG-17)
   Funciones: setBool('GS_14', valor) → registro-v3-app.js
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: { GS_14: true | false }   → radios «Sí / No»

PASO 1 — `answers['GS_14'] = 'true' | 'false'` (estado local; se guarda con «Guardar sección», RA-E08)
PASO 2 — Reactividad:
  SI es «Sí»: se muestran las tarjetas de las entidades acudidas y el formulario de entidad (si el Caso sigue en Borrador)
  SI es «No»: se ocultan el formulario y la lista **sin borrar nada** (las tarjetas ya registradas se vuelven a ver al responder «Sí»; si `GS_17` es «Sí» y no hay entidades acudidas, el formulario se muestra igual)

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión | Paso afectado |
|---|---|
| Responder «No» con entidades ya agregadas no las elimina | PASO 2 |
