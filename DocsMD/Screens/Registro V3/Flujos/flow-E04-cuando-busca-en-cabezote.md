━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando busca en el Cabezote
   Tipo: User Interaction
   Función: buscarEnCabezote() → components/cabezote-busqueda.js
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

> **Corregido por diagrama del cliente** (imagen de flujo recibida, ver
> `registro-v3-data.md` D7): la búsqueda de duplicidad se hace **solo por
> número de documento y número de teléfono — nunca por nombre y apellido**.
> Esto reemplaza el diseño original (3 campos, incluido nombre) que traía
> el Guion Funcional/Excel de Drive. El nombre sigue **mostrándose** en los
> resultados (para que el agente confirme visualmente que es la persona
> correcta) pero ya no es un criterio de búsqueda ni de coincidencia.

INPUT: {
  documento: texto, opcional     → <input> número de documento
  telefono:  texto, opcional     → <input> número de teléfono
}

Validar formulario:
  • Al menos uno de los 2 campos debe tener valor
  SI los 2 están vacíos:
    → Mostrar mensaje "Digite al menos un criterio" inline
    → TERMINAR ejecución

PASO 1 — Buscar Reportes pendientes de validación

[GET] /api/v1/reportes/buscar
// filter:
{
  "documento": documento,   // input del usuario
  "telefono":  telefono
}
→ resultados.reportes = respuesta.data.reportes

PASO 2 — Buscar Registros (víctimas) y sus Casos asociados

[GET] /api/v1/victim-cases/buscar
// filter:
{
  "documento": documento,
  "telefono":  telefono
}
→ Agrupa la respuesta por documento de identidad (mismo documento = mismo "Registro", ver registro-v3-data.md D1)
→ El `nombre` de cada resultado se incluye en la respuesta solo para mostrarlo en PerfilResumen — no participa en el filtro
→ Para cada grupo, calcula BadgeEstadoRegistro (Incompleto/Completo) evaluando las obligatorias de CI+DP+CA del `victim_case` más reciente del grupo
→ resultados.registros = lista de { perfil, estadoRegistro, casos: [...] }

PASO 3 — Buscar Borradores asociados

Dentro de la misma respuesta de PASO 2, filtrar los `victim_case` con estado = Borrador
→ resultados.borradores = lista de { telefono, fechaCreacion, canal, icode }

PASO 4 — Mostrar resultados
  buscado = true
  Vue re-evalúa buscado y resultados → se renderiza la sección de resultados (ReporteResumen × N, RegistroResumen × N con sus CasoResumen anidados, BorradorResumen × N, y el BtnNuevoRegistro al final)

SI resultados.reportes.length === 0 Y resultados.registros.length === 0 Y resultados.borradores.length === 0:
  → Mostrar "Sin coincidencias por documento o teléfono"
  → El BtnNuevoRegistro sigue visible (siempre lo está tras buscar)
  → CONTINÚA FLUJO GENERAL

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión                                                | Paso afectado |
|------------------------------------------------------------------------|---------------|
| Endpoint `/api/v1/reportes/buscar` no existe todavía — el módulo de Reportes (intake de terceros) no está construido en esta ronda | PASO 1 |
| El Guion Funcional de Drive describía "búsqueda aproximada (tildes, orden de nombres)" — esa redacción asumía búsqueda por nombre, que el diagrama del cliente elimina. Para documento/teléfono probablemente basta con normalizar espacios/guiones, no coincidencia difusa — confirmar | PASO 1, PASO 2 |
| Si debe mantenerse como opción secundaria/"avanzada" un campo de nombre oculto por defecto, para los casos raros donde el agente no tiene ni documento ni teléfono a mano (ej. la víctima solo recuerda su nombre) — el diagrama no contempla ese caso límite | Todo el evento |
