━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando busca en el Cabezote
   Tipo: User Interaction
   Estado: 🟡 Implementado con cambios — sincronizado con el código 2026-10-31
   Funciones: ejecutarBusqueda() → registro-v3-app.js · Search() → victim_case_v3_controller.go
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

> La búsqueda de duplicidad es **solo por número de documento y teléfono — nunca por nombre** (D7 en
> `registro-v3-data.md`). El nombre se muestra en los resultados pero no participa en el filtro.

INPUT: {
  documento:  texto, opcional   → <input> "Documento"        (buscar.documento)
  telefono:   texto, opcional   → <input> "Teléfono"         (buscar.telefono)
  tipoDocumento: texto, opcional → estado buscar.tipoDocumento; la pantalla no tiene input para él, así que siempre va vacío
}

Validar formulario:
  • Al menos uno de documento / teléfono debe tener valor
  SI los dos están vacíos:
    → errorBusqueda = "Digite al menos documento o teléfono."
    → TERMINAR ejecución (sin llamada al backend)

PASO 1 — Buscar
  buscando = true
  [GET] /api/v1/victim-cases-v3/buscar
  // query string:
  {
    "tipo_documento": buscar.tipoDocumento,   // solo si no está vacío
    "documento":      buscar.documento,       // solo si no está vacío
    "telefono":       buscar.telefono         // solo si no está vacío
  }

  Backend (Search, VictimCaseV3Controller):
    SI documento y teléfono vienen vacíos: responde 400 "digite al menos un criterio (documento o teléfono)"
    repo.FindByDocumentOrPhone:
      SI hay documento: WHERE victim_case_victim_doc_number = documento [AND victim_case_victim_doc_type = tipo]   // el teléfono se IGNORA
      SI NO:            WHERE victim_case_form2_victim_phone = teléfono
    PARA CADA fila: registroCompleto = ¿el Caso MÁS RECIENTE con ese tipo+número de documento tiene respondidas
                    las obligatorias de CI + DP (+ DP_10, DP_15) + CA_1-4?   (se calcula una vez por documento)
    Respuesta: { "registros": [{ icode, names, lastNames, docType, docNumber, phone, status, registroCompleto }] }

PASO 2 — Mostrar resultados
  resultadosBusqueda = data.registros || []
  buscando = false
  SI hay resultados: tabla con Nombre · Documento (tipo-número) · Teléfono · Estado (código crudo bo/in/te/…) · badge Completo/Incompleto · botón "Abrir"
  SI NO hay resultados: no se muestra ninguna tabla ni mensaje de "sin coincidencias"; "+ Nuevo registro" sigue visible
  SI la petición falla: errorBusqueda = "Error al buscar."
  → CONTINÚA FLUJO GENERAL

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión | Paso afectado |
|---|---|
| Planeado y no implementado: buscar **Reportes pendientes de validación** (`/api/v1/reportes/buscar` nunca existió; los reportes ciudadanos están en `victim_contact` y no se consultan) y agrupar los resultados en Reportes / Registros con sus Casos / Borradores | PASO 1, PASO 2 |
| Si se escribe documento y teléfono a la vez, el teléfono no filtra: la coincidencia es solo por documento | PASO 1 (backend) |
| Sin mensaje "sin coincidencias" ni indicación de qué Caso es Borrador: el estado se muestra con el código crudo | PASO 2 |
| G-V3-15 — alternativa si el agente no tiene documento ni teléfono | Todo el evento |
