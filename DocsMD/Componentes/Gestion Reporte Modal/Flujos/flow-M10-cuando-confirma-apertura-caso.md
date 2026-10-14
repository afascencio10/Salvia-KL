━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando confirma apertura de caso
   Tipo: User Interaction / Backend write
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

⏸️ DIFERIDO — no implementado. Desde 2026-10-13 el backend RECHAZA close_reason='ap'
con HTTP 400 (deferredCloseReasonCodes en el servicio de gestión), para que tampoco se
llegue por API a "Cerrado - válido" sin abrir un caso. Al construir este evento hay que
quitar 'ap' de ese mapa. El diseño de abajo se conserva tal cual.

Disparado por: clic en "Confirmar apertura de caso", step 'postContacto',
               con pendingClose === 'apertura'

INPUT: {
  icode:                 reporte abierto en el modal
  atencionEspecializada: 'No' | 'Si'   // solo relevante si report.tipo en (Hombre, Familiar) — GAP-6
}


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND — gestion-reporte-modal.js
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━


PASO 1 — Llamar a cerrar (categoría Válido)

  saving = true
  POST /api/v1/victim-contacts/{icode}/cerrar
  body: { close_reason: 'ap', atencion_especializada: atencionEspecializada === 'Si' }


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  BACKEND — POST /api/v1/victim-contacts/:icode/cerrar
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━


PASO 2 — Validar precondiciones

  Igual a M-02 PASO 2 (existe, no cerrado, canManage).


PASO 3 — Resolver enum y actualizar

  reasonRow = enum categoría 'victim_contact_close_reason', code 'ap'

  UPDATE victim_contact SET
    victim_contact_close_reason = reasonRow.id,
    victim_contact_closed_date  = NOW()
  WHERE victim_contact_i_code = :icode

  // atencion_especializada (GAP-6): no tiene columna propia todavía — no se
  // persiste hasta que se confirme que el requisito es real. Si se
  // confirma, se agrega como columna nullable adicional en victim_contact
  // o como parte del payload que recibe el formulario de caso.


PASO 4 — Responder

  { ok: true, closeReason: 'ap', category: 'Válido', reportType: vc.report_type }

  → FIN EJECUCIÓN BACKEND ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND — redirección según tipo (GAP-7)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━


PASO 5 — Decidir destino

  saving = false

  SI report.tipo IN ('Victima', 'Tercero'):
    → window.location = "/salvia/casos/" + icode + "/nuevo"
    // Ruta YA EXISTENTE (VictimCasePOST_GET), ya pre-llena el formulario de
    // caso a partir del VictimContactDTO. No se construye nada nuevo aquí.

  SI NO (tipo IN ('Familiar', 'Hombre')):
    // El mockup señala literalmente que estos formularios "no hacen parte
    // de este mockup" (línea 662) — y la lectura de código confirma que
    // tampoco existen hoy en la aplicación real.
    → mostrar mensaje: "Este reporte quedó cerrado como válido. La apertura
      de caso para reportes de tipo {tipo} aún no tiene una pantalla
      destino — contacte al equipo de producto."
    → step = 'cerrado'   // el reporte ya quedó cerrado en BD (PASO 3),
                          // independientemente de que la navegación no
                          // tenga a dónde ir todavía

  → FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

| Variable / decisión                                                              | Paso afectado |
|------------------------------------------------------------------------------------|---------------|
| GAP-6 — ¿"Remitir a Atención Especializada" es un requisito real? ¿Dónde se guarda si sí? | PASO 1, PASO 3 |
| GAP-7 — ¿Qué debe pasar hoy si se intenta abrir caso desde un reporte Familiar/Hombre, dado que el formulario destino no existe? ¿Se bloquea el botón "Apertura de caso" directamente para esos tipos en vez de dejar cerrar el reporte sin destino? | PASO 5 — alternativa a evaluar con el usuario: ocultar la opción 'apertura' completa para tipo Familiar/Hombre hasta que exista destino |
| ¿Cerrar el reporte como "Válido" sin poder navegar a ningún lado (rama Familiar/Hombre) es aceptable, o se prefiere NO cerrar el reporte hasta que exista la pantalla destino? | PASO 3 vs. PASO 5 — depende de la respuesta a GAP-7 |
