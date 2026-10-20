━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando presiona "Nuevo Registro"
   Tipo: User Interaction
   Función: crearNuevoRegistro() → components/cabezote-busqueda.js
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  busqueda: { documento, telefono }   → mismos criterios usados en la búsqueda (E04, sin nombre — ver D7)
  agentId:  id del Agente Integral     → sesión
}

> Agregado al revisar `Mockup_Reportes_Registro (1).html` (fuente de menor
> prioridad, solo aporta un detalle de UX que no contradice nada definido):
> si la búsqueda del Cabezote (E04) tuvo coincidencias, el botón "+ Nuevo
> Registro" permanece deshabilitado hasta que el agente marque la casilla
> "Verifiqué las coincidencias: se trata de otra persona" — evita crear un
> Registro duplicado por descuido.

Validar formulario:
  SI resultados.registros.length > 0 Y !busqueda.otraPersona:
    → BtnNuevoRegistro permanece deshabilitado
    → TERMINAR ejecución

PASO 1 — Crear el Caso en Borrador prediligenciado

> Reutiliza el mecanismo ya implementado y probado de V2-real (`CreateDraft`
> + `UpsertForm2FromAnswers`, ver `registro-v3-data.md` § D8) en vez de
> inventar uno nuevo: la fila de `victim_case` se crea de inmediato con
> valores dummy en los 5 campos núcleo que aún no se conocen (nombres,
> apellidos, tipo/número de documento, municipio de residencia — "Pendiente"/
> "cc"/"0000000000"/"00000000"), y la fila de `victim_case_form2` se crea
> en el mismo momento con "No" como dummy en todas sus columnas NOT NULL
> (incluidas las ~55 nuevas de Tamizaje V3, VPP, tipo/nivel de riesgo y
> remisión). Cada valor dummy se sobrescribe con la respuesta real en cuanto
> el agente la diligencia, sin importar en qué módulo ni en qué orden.

[POST] /api/v1/victim-cases/draft
// payload:
{
  "victim_case_phone_hint":    busqueda.telefono,
  "victim_case_document_hint": busqueda.documento,
  "victim_case_status":        "bo",                  // Borrador — código ya existe en VICTIM_CASE_STATUS
  "victim_case_agent_id":      agentId,
  "victim_case_channel":       "web" // o el canal real de origen, por confirmar — ver GAPS
}
→ Internamente el backend aplica `CreateDraft` (victim_case con dummies) +
  `UpsertForm2FromAnswers` (victim_case_form2 con dummies) sobre las
  respuestas vacías que trae el formulario recién abierto — no hace falta
  esperar a ningún `saveSection` para que la fila exista
→ respuesta.icode → nuevo identificador del Caso

PASO 2 — Abrir el Formulario
  registro = { estado: 'Incompleto' }  // calculado, no viene de BD — ver registro-v3-data.md D1
  caso = { icode: respuesta.icode, estado: 'Borrador' }
  moduloActual = 'CI'
  stage = 'formulario'

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión                                              | Paso afectado |
|----------------------------------------------------------------------|---------------|
| Catálogo real de "canal de origen" a registrar automáticamente (Línea 155, web, app móvil, PQRS) y cómo se captura en esta pantalla si el agente la abre manualmente desde Mi Salvia | PASO 1 |
| ~~Endpoint `/api/v1/victim-cases/draft` es propuesto~~ — **resuelto**: reutiliza `CreateDraft`/`UpsertForm2FromAnswers` de V2-real, ver D8. Falta solo adaptar `victim_case_form_fields.go` con las columnas nuevas de V3 | PASO 1 |
| G-V3-16 — si V3 también necesita crear un `general_user` con credenciales para la víctima como hace V2-real, o si eso es específico de V2 (ver D8) | PASO 1 |
