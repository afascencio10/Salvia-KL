━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando confirma cierre por decisión de la persona
   Tipo: User Interaction / Backend write
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Disparado por: clic en "Confirmar cierre", step 'postContacto',
               con pendingClose === 'decision'

Precondición: motivoTexto.trim() !== ''  (el botón está :disabled si no)

INPUT: {
  icode:       reporte abierto en el modal
  motivoTexto: texto escrito en el textarea  → "manifestación de la persona"
}


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND / BACKEND
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Idéntico a M-08 en todos sus pasos (mismo endpoint, misma validación de
observación obligatoria), con una sola diferencia:

  body: { close_reason: 'dp', close_observation: motivoTexto.trim() }

  → reasonRow = enum categoría 'victim_contact_close_reason', code 'dp'
  → closeReasonLabel final = 'Inválido por Cierre por decisión de la persona — ' + motivoTexto

El texto del label del textarea también cambia en la UI ("Registre la
manifestación de la persona:" en vez de "Registre el motivo de no
competencia:"), pero es un detalle puramente visual sin efecto en el flujo.


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Mismos que M-08 (ver flow-M08).
