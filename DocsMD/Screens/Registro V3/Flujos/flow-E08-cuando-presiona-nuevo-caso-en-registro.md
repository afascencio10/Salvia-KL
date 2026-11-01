> ⛔ **No implementado (verificado en el código el 2026-10-31).** El Cabezote no tiene el botón "Nuevo Caso en Registro existente" ni existe el endpoint `nuevo-caso-mismo-registro`. Está diferido junto con la separación Registro/Caso (decisión D1); hoy cada "+ Nuevo registro" crea un Caso independiente (con sincronización del Perfil al guardar, ver flow-E11). Se conserva el diseño original por si se retoma.

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando presiona "Nuevo Caso en Registro existente"
   Tipo: User Interaction
   Función: nuevoCasoEnRegistro(registroEncontrado) → components/cabezote-busqueda.js
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

> Etiqueta del botón (detalle de `Mockup_Reportes_Registro (1).html`, no
> cambia la lógica, solo el texto mostrado): si el Registro tiene algún Caso
> todavía activo (Borrador/Iniciado/Terminado, ninguno Cerrado) con tríada
> distinta → el botón dice **"Nuevo Caso en Registro existente"** (este
> flujo). Si **todos** los Casos del Registro ya están Cerrados y llega un
> nuevo hecho → el mismo botón cambia su texto a **"Nuevo hecho tras cierre
> de Caso"**, misma lógica de backend, sin distinción funcional.

INPUT: {
  registroEncontrado: { perfil, estadoRegistro, casos }   → RegistroResumen sobre el que se hizo click
  agentId: id del Agente Integral                          → sesión
}

PASO 1 — Identificar el Caso más reciente del mismo documento

  casoMasReciente = registroEncontrado.casos ordenados por fecha descendente [0]
  // Es la fuente del Perfil de víctima a copiar — ver registro-v3-data.md D1

PASO 2 — Crear el nuevo Caso copiando el Perfil

[POST] /api/v1/victim-cases/nuevo-caso-mismo-registro
// payload:
{
  "source_icode":  casoMasReciente.icode,   // de dónde se copia CI+DP+CA
  "agent_id":      agentId,
  "victim_case_status": "bo"               // Borrador
}

  Lógica del endpoint (backend):
    Consultar victim_case_form2 WHERE icode = source_icode
    Copiar las columnas de CI+DP+CA (no las de HV/AG/TZ/GS — esas son del Caso, no del Perfil)
    Crear fila nueva de victim_case + victim_case_form2 con esas columnas copiadas
    Generar nuevo icode

→ respuesta.icode → nuevo identificador del Caso

PASO 3 — Abrir el Formulario en el módulo HV
  registro = { ...registroEncontrado.perfil, estado: registroEncontrado.estadoRegistro }
  caso = { icode: respuesta.icode, estado: 'Borrador' }
  moduloActual = 'HV'   // CI/DP/CA ya vienen completos desde la copia
  stage = 'formulario'

→ Ver flujo: Cuando actualiza una respuesta del formulario (si el agente decide corregir algún campo de CI/DP/CA copiado, dispara el sub-flujo de sincronización fan-out)

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión                                                        | Paso afectado |
|---------------------------------------------------------------------------------|---------------|
| Si `estadoRegistro` era "Incompleto" en el Caso origen, el nuevo Caso nace con un Perfil también incompleto — confirmar que esto es aceptable (el agente debería completar el Perfil en cualquiera de los Casos, y el fan-out lo propaga) | PASO 2 |
