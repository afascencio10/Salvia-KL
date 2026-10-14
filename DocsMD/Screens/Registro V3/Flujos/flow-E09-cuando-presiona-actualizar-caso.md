━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando presiona "Actualizar Caso"
   Tipo: User Interaction
   Función: actualizarCaso(registro, caso) → components/cabezote-busqueda.js
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  registro: perfil del Registro encontrado   → RegistroResumen
  caso:     { icode, estado: 'Iniciado' }     → CasoResumen sobre el que se hizo click (solo visible si estado === 'Iniciado')
  agentId:  id del Agente Integral             → sesión
}

PASO 1 — Registrar trazabilidad del nuevo episodio

[POST] /api/v1/victim-cases/{caso.icode}/nuevo-episodio
// payload:
{
  "agent_id": agentId,
  "fecha":    now()
}
→ Se agrega una entrada a la bitácora: "{agentId} registró un nuevo episodio sobre este Caso el {fecha}"

PASO 2 — Cargar el Caso completo

[GET] /api/v1/victim-cases/{caso.icode}
→ caso = { icode, estado: 'Iniciado', ...todas las respuestas ya diligenciadas }
→ registro = perfil calculado de CI+DP+CA

PASO 3 — Abrir el Formulario
  moduloActual = 'HV'   // el nuevo episodio se registra sobre Hechos — la tríada víctima-agresor no cambia
  stage = 'formulario'

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión                                                                  | Paso afectado |
|------------------------------------------------------------------------------------------|---------------|
| Si "Actualizar Caso" registra un episodio nuevo de Hechos, ¿se sobrescriben las respuestas anteriores de HV o se necesita un historial de episodios dentro del mismo Caso? El Guion Funcional no lo detalla | PASO 3 |
