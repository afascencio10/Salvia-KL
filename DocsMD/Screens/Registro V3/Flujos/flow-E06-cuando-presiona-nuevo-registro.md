━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando presiona "Nuevo Registro"
   Tipo: User Interaction
   Estado: 🟡 Implementado con cambios — sincronizado con el código 2026-10-31
   Funciones: crearNuevoRegistro() → registro-v3-app.js · CreateDraft() → victim_case_v3_controller.go / service / repository
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  buscar.tipoDocumento, buscar.documento, buscar.telefono   → mismos criterios de la búsqueda (pueden estar vacíos)
  agente:  UserICode de la sesión                           → backend
}

(El botón "+ Nuevo registro" está **siempre visible y habilitado**: no exige haber buscado ni confirmar que se trata de otra persona.)

PASO 1 — Crear el Caso en Borrador
  [POST] /api/v1/victim-cases-v3/draft
  // payload (lo que manda el frontend; `municipio_residencia` y `canal` existen en la API pero no se envían):
  {
    "tipo_documento": buscar.tipoDocumento,
    "documento":      buscar.documento,
    "telefono":       buscar.telefono
  }

  Backend (CreateDraft), todo en una transacción:
    1. general_user_profile  con nombres/apellidos "Pendiente", documento y municipio dummy (cc / 0000000000 / 00000000) salvo los recibidos
    2. general_user  con login y contraseña aleatorios (hash bcrypt), estado 'e', idioma 'sp'
    3. rel_role_general_user  con el rol 'us'
    4. victim_case  en estado 'bo' (Borrador) con nombres "Pendiente", el documento recibido o el dummy, `agent_id` = agente de la sesión
    5. victim_case_form2  con TODAS las columnas NOT NULL rellenas con un valor dummy según su tipo (bool/enum → el id de "No", texto → "", número → 0, fecha → 1900-01-01, hora → 00:00:00)
    6. SI llegó teléfono: victim_case_form2_victim_phone = teléfono y se marca CI_4 como respondida
    7. victim_case_new_user_credentials = {login, pass} en texto plano (jsonb)
  Respuesta: { "icode": <victim_case_i_code> }

PASO 2 — Abrir el Formulario
  SI la respuesta trae icode: abrirCaso(icode)   → ver flow-E07 (carga el Caso y abre en CI)
  SI NO: no pasa nada (no hay mensaje de error)

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión | Paso afectado |
|---|---|
| Planeado y no implementado: checkbox "Verifiqué las coincidencias: se trata de otra persona" que habilitaba el botón cuando la búsqueda tenía resultados | Todo el evento |
| El canal de origen no se captura: la API acepta `canal` pero el frontend no lo envía | PASO 1 |
| Sin job de expiración para Borradores abandonados (heredado de V2-real, G-V3-17); cada "Nuevo registro" crea además un `general_user` con credenciales en texto plano | PASO 1 |
| Si el POST falla no se muestra ningún error al agente | PASO 2 |
