━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando confirma cierre por no competencia
   Tipo: User Interaction / Backend write
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Disparado por: clic en "Confirmar cierre", step 'postContacto',
               con pendingClose === 'competencia'

Precondición: motivoTexto.trim() !== ''  (el botón está :disabled si no)

INPUT: {
  icode:       reporte abierto en el modal
  motivoTexto: texto escrito en el textarea  → observación obligatoria
}


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND — gestion-reporte-modal.js
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━


PASO 1 — Validar localmente (defensa adicional a :disabled)

  SI motivoTexto.trim() === "":
    → TERMINAR ejecución (no debería ocurrir, el botón ya lo impide)


PASO 2 — Llamar a cerrar

  saving = true
  POST /api/v1/victim-contacts/{icode}/cerrar
  body: { close_reason: 'co', close_observation: motivoTexto.trim() }


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  BACKEND — POST /api/v1/victim-contacts/:icode/cerrar
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━


PASO 3 — Validar precondiciones

  Igual a M-02 PASO 2 (existe, no cerrado, canManage).

  SI NO — validación adicional específica de este motivo:
  SI close_observation vacío, ausente o solo espacios:
    → HTTP 400 "la observación es obligatoria para este motivo"
    → TERMINAR ejecución
  SI close_observation (sin espacios laterales) supera 500 caracteres:
    → HTTP 400 "la observación no puede superar 500 caracteres"   // varchar(500); sin esta validación daba 500
    → TERMINAR ejecución


PASO 4 — Resolver enum y actualizar

  reasonRow = enum categoría 'victim_contact_close_reason', code 'co'

  UPDATE victim_contact SET
    victim_contact_close_reason      = reasonRow.id,
    victim_contact_close_observation = :close_observation,
    victim_contact_closed_date       = NOW()
  WHERE victim_contact_i_code = :icode


PASO 5 — Responder

  { ok: true, closeReason: 'co', category: 'Inválido' }

  → FIN EJECUCIÓN BACKEND ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND — transición de vista
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━


PASO 6 — Mostrar confirmación

  saving = false
  pendingClose = ''; motivoTexto = ''
  closeReasonLabel = 'Inválido por Cierre por no competencia — ' + motivoTexto original
  categoria = 'Inválido'
  step = 'cerrado'

  → FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — resueltos tras la implementación
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

| Variable / decisión                                                              | Resolución |
|------------------------------------------------------------------------------------|------------|
| Longitud máxima de la observación | 500 caracteres: el textarea tiene `maxlength="500"` y el backend responde 400 si se supera (casos `M08-C3`, `M09-C3`) |
