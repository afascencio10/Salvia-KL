━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando abre el modal
   Tipo: Lifecycle
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Disparado por: el padre llama open(report) tras recibir action-clicked
               con buttonId==='ver' (Listado Reportes E-04 / Mis Reportes E-02)

INPUT: {
  report:  objeto reporte de la fila (ReportListItem)  → parámetro de open()
}


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND — gestion-reporte-modal.js
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━


PASO 1 — Abrir y mostrar loading

  visible = true
  loading = true
  icode = report.icode


PASO 2 — Pedir el detalle completo al backend

  GET /api/v1/victim-contacts/{icode}

  // La fila de la tabla (ReportListItem) no trae: datos del reportante,
  // ajuste razonable, descripción de hechos, ni la lista de intentos.


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  BACKEND — GET /api/v1/victim-contacts/:icode (nuevo)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━


PASO 3 — Cargar victim_contact + su origen (GAP-1)

  vc = SELECT * FROM victim_contact WHERE victim_contact_i_code = :icode
  SI vc no existe:
    → HTTP 404
    → TERMINAR ejecución

  SI NO, una sola consulta arma el detalle leyendo el origen que exista (GAP-1):
    reportType        = enum de victim_contact_form2.report_type (V1 o V2) → etiqueta corta (Víctima/Tercero/Familiar/Hombre)
    reporterNames/LastNames/Email y parentesco = victim_contact_form2 + victim_contact_form2_v2 (solo reportes V2)
    telefono          = el primero que no sea 0/NULL de: form2.victim_col_phone → form2.reporter_phone → form1.phone
                        // en Familiar/Hombre la víctima no tiene teléfono (se guarda 0): el de contacto es el del reportante
    ajusteRazonable   = ajustes de victim_contact_form2 (relación de enums) → texto vía Locale
    descripcionHechos = form2.facts_description → form1.facts_description
  // Un reporte "solo-form1" (lo registró un agente; 2 951 de 3 503) no trae tipo, reportante ni ajustes.


PASO 4 — Cargar intentos y calcular contadores

  intentos = SELECT * FROM victim_contact_attempt
             WHERE victim_contact_attempt_victim_contact = :icode
               AND victim_contact_attempt_deleted_at IS NULL
             ORDER BY victim_contact_attempt_date ASC

  intentosHoy  = COUNT(intentos WHERE result='fallido' AND fecha::date = hoy en America/Bogota)
  // La columna es timestamp SIN zona y guarda UTC: el día Bogotá sale de
  // (fecha AT TIME ZONE 'UTC' AT TIME ZONE 'America/Bogota')::date  (a secas, AT TIME ZONE 'America/Bogota' desplaza el día)
  // Cada intento trae agentName = nombre del perfil de quien lo registró.

  // "Día completo" = día con >= 3 intentos fallidos (no cualquier día con
  // intentos) — corrección 2026-10-12, ver flow-M07 PASO 2 para el
  // detalle y la fuente en el mockup.
  fallidosPorDia = agrupar intentos WHERE result='fallido' POR fecha::date, contar cada grupo
  diasCompletos  = COUNT(días EN fallidosPorDia CON count >= 3)


PASO 5 — Determinar canManage (GAP-2, confirmado 2026-10-11)

  canManage = (sesión.role === 'op' OR sesión.role === 'sv') AND closeReasonCode es vacío
  // Decisión del usuario: op y sv pueden gestionar cualquier reporte visible
  // para ellos, no solo el agente con victim_contact_agent_id propio. En la
  // práctica es siempre true aquí (solo op/sv llegan a abrir este modal) —
  // se mantiene como defensa en profundidad, no como restricción de UI.
  // En un reporte cerrado (o legado invalidado) canManage es false: no hay nada que gestionar.


PASO 6 — Responder

  // Forma PLANA (ReportDetail en internal/models/report_list_item.go), no anidada en "report":
  {
    icode, nombreMostrado, reportType,
    reporterNames, reporterLastNames, parentesco, reporterEmail, telefono, ajusteRazonable, descripcionHechos,
    intentos: [ { date (ISO, UTC), result ('fallido'|'exitoso'), agentName } ],
    intentosHoy, diasCompletos,
    closeReasonCode,     // 'ni'|'nc'|'co'|'dp'|'ap', o 'legado' si el reporte lo invalidó el flujo anterior; '' si está abierto
    closeReasonLabel,    // "Inválido por No información Mínima", …; en un legado, su victim_contact_status_description
    closeCategory,       // 'Válido' (solo ap) | 'Inválido' | ''
    closeObservation,    // observación de co/dp
    canManage
  }

  → FIN EJECUCIÓN BACKEND ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND — determinar step inicial
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━


PASO 7 — Elegir vista

  loading = false
  huboExitoso = response.intentos.some(i => i.result === 'exitoso')

  SI response.closeReasonCode != '':     // incluye 'legado'
    step = 'cerrado'
    // motivo mostrado = closeReasonLabel + " — " + closeObservation (si hay observación), igual que justo tras cerrar
  SI NO SI huboExitoso:
    step = 'postContacto'
    // Corrección (2026-10-11): antes esta rama no existía — un intento
    // exitoso sin cierre terminado (modal cerrado a mitad de M-08/M-09)
    // reabría en 'esquema3x3' en vez de 'postContacto'. Esa es la razón de
    // fondo por la que M-06 persiste el intento exitoso en vez de ser una
    // transición puramente local como M-04: sin esta fila en BD, reabrir el
    // modal no tendría forma de saber que ya hubo contacto, y el agente
    // podría terminar cerrando por "no contacto" un reporte donde sí
    // contactó a la persona.
  SI NO SI response.intentos.length > 0:
    step = 'esquema3x3'
  SI NO:
    step = 'validacionPrevia'

  → FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — resueltos tras la implementación
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

| Variable / decisión                                                              | Resolución |
|------------------------------------------------------------------------------------|------------|
| Origen exacto del reporte (V1 / V2 / agente) y cómo leer su form correspondiente   | El detalle lee `form2`, `form2_v2` y `form1` según existan (ver PASO 3). Confirmar con negocio si hay un 4.º origen |
| ¿Qué pasa si un intento ya se registró "hoy" según el reloj del servidor pero el agente cruzó medianoche durante la sesión abierta? | Se recalcula siempre en el backend al abrir; el cliente no cachea contadores entre aperturas (caso `M01-C24`) |
