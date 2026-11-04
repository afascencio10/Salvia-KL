━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando presiona «Anterior» entre módulos
   Tipo: User Interaction
   Estado: ✅ Implementado — 2026-11-02 (WP-3 · REG-01)
   Función: irAAnterior() · moduloAnterior(id) → registro-v3-app.js
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  moduloActivo:   módulo visible                         → estado local
  modulos[]:      lista ordenada de módulos              → estado local (el menú y las flechas se arman con ella; no hay módulos escritos a mano)
  triadaCompleta: define qué módulo está bloqueado       → estado local
}

PASO 1 — Calcular el destino
  destino = moduloAnterior(moduloActivo): primer módulo anterior (hacia atrás) en modulos[] con moduloBloqueado(id) == false; null si no hay.
  SI destino == null:
    → el botón está deshabilitado (:disabled="!moduloAnterior(moduloActivo)"), sin acción
    → TERMINAR ejecución
  SI NO:
    → CONTINÚA FLUJO GENERAL

PASO 2 — Cambiar de módulo
  irAModulo(destino)  → Ver flujo: E10 Cuando cambia de módulo

PASO 3 — Reflejo en pantalla (reactividad de Vue; parte de este evento)
  El menú resalta el nuevo módulo y la barra muestra "Módulo i de N". NO se guarda nada: las respuestas sin guardar siguen en `answers`.

Ejemplos (orden CI, DP, CA, HV, AG, TZ, GS):
  • Tríada incompleta (TZ bloqueado): el anterior de GS es AG.   • Tríada completa: el anterior de GS es TZ.   • Desde CI: null.

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión | Paso afectado |
|---|---|
| No advierte si el módulo actual tiene cambios sin guardar (igual que las pestañas anteriores) | PASO 2 |
| El número E23 es de WP-3; si otra entrega usa el mismo número al integrar, se renumera | — |
