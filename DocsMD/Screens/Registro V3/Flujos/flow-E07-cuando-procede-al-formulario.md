━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando presiona "Proceder al formulario"
   Tipo: User Interaction
   Función: procederAlFormulario(borrador) → components/cabezote-busqueda.js
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  borrador: { icode, telefono, fechaCreacion, canal }   → objeto BorradorResumen sobre el que se hizo click
  agentId:  id del Agente Integral                       → sesión
}

PASO 1 — Registrar trazabilidad de retome

[POST] /api/v1/victim-cases/{borrador.icode}/retome
// payload:
{
  "agent_id": agentId,   // quién retoma
  "fecha":    now()      // cuándo
}
→ Se agrega una entrada a la bitácora/log del Caso: "{agentId} retomó el Borrador el {fecha}"

PASO 2 — Cargar el Caso completo

[GET] /api/v1/victim-cases/{borrador.icode}
→ caso = { icode, estado: 'Borrador', ...respuestas ya diligenciadas }
→ registro = perfil calculado de CI+DP+CA ya diligenciado (estado 'Incompleto' por definición, un Borrador nunca es Completo)

PASO 3 — Determinar el módulo de entrada
  moduloActual = módulo con menor progreso dentro de los ya habilitados (típicamente CI si nada se ha diligenciado, o el primer módulo incompleto)

PASO 4 — Abrir el Formulario
  stage = 'formulario'

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión                                  | Paso afectado |
|-----------------------------------------------------------|---------------|
| Regla exacta para elegir "el módulo de entrada" cuando hay progreso parcial en varios módulos no contiguos | PASO 3 |
