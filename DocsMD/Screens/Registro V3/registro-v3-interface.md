# `/salvia/registro-v3` — Interfaz de la Pantalla

> 🔄 **Actualizado el 2026-11-02 (WP-2: REG-07/08/09/10/11/12/34, ver [`changelogNov2026.md`](changelogNov2026.md))** — hora del hecho como selector (HV_3), pregunta «Tipo de discapacidad» (DP_51), zona de residencia con opciones, Tamizaje en tablas y por etapas, y avisos visibles cuando un guardado o una carga fallan.
>
> ✅ **Sincronizado con el código el 2026-10-31** (versión implementada del 2026-10-14/15). La planeación original —con `<dinamic-form>`, sidebar con barras de progreso, componentes `habeas-data-gate`/`cabezote-busqueda` y el botón "Alerta inminente" del header— **no se construyó así**: la decisión A1 (ver [`registro-v3-data.md`](registro-v3-data.md) § D10) descartó `dinamic-form` y todo quedó en una app Vue a medida. Lo planeado y no construido está en la sección "Planeado y no implementado" al final.

> 🔄 **Actualizado el 2026-11-02 (WP-3: REG-01/02/03/04/05/06/13/35/37).** La pantalla se llama «Registro de víctima»; el menú de módulos es vertical (se arma con `modulos[]`) con porcentaje de avance, hay barra «Anterior / Siguiente» y guardar avanza al siguiente módulo; solo el Tamizaje se bloquea; el encabezado muestra nombre, documento y estado legibles; los faltantes salen en rojo; «Finalizar» guarda GS primero. Detalle en [`changelogNov2026.md`](changelogNov2026.md).

> 🔄 **Actualizado el 2026-11-04 (WP-4: REG-14/15/16/17/18/19/20/38\*, REG-04 c; ver [`changelogNov2026.md`](changelogNov2026.md)).** «Gestión Salvia» pasa a **«Plan de acción Salvia»** con el orden a–i definido por el cliente y sin los títulos A–D (REG-15); nace el módulo **«Ruta de atención integral»** (8.º, entre Tamizaje y Plan de acción; REG-17, carpeta [`../Ruta de atención integral/`](../Ruta%20de%20atenci%C3%B3n%20integral/)); `GS_14` y `GS_17` se contestan allí. La barrera es el formulario de Seguimiento v4 (REG-18), el canal de activación es de selección múltiple como en Seguimiento (REG-20), agregar entidad funciona como en Gestión institucional (REG-19), `GS_20` «elementos psicosociales» es nueva (REG-16), `GS_2` se retira (REG-14) y la «Gestión» de la barrera muestra etiqueta (REG-04 c). Entidades y barreras se guardan en un **borrador JSON** mientras el Caso es Borrador y quedan de **solo lectura** una vez materializadas (REG-38\*, ver [`Flujos/flow-subflujo-materializar-borrador.md`](Flujos/flow-subflujo-materializar-borrador.md)).

Pantalla de ingreso y registro de hechos de VBG/VPP para los roles `op` y `sv` (permiso `manage_registro_v3`). Recorre seis etapas dentro de una sola página: autorización de datos personales → tipo de contacto → Cabezote de búsqueda de duplicidad → Formulario de siete módulos (CI, DP, CA, HV, AG, TZ, GS) → Cierre. No hay entrada en el menú lateral: se llega por la URL. Guarda por módulo en `/api/v1/victim-cases-v3/*` y escribe directo sobre `victim_case`/`victim_case_form2`, las mismas tablas que lee Seguimiento.

## Archivos relevantes

| Archivo | Rol |
|---|---|
| `src/frontend/html/salvia/registro_v3/registro_v3.html` | Template único: las 6 etapas con `v-if="etapa===…"` |
| `src/frontend/js/components/registro-v3-app.js` | App Vue 3 `RegistroV3App` (Options API, delimitadores `${ }`): estado, validaciones, llamadas a la API, párrafo síntesis |
| `src/frontend/js/components/registro-v3-campos.js` | Catálogo de preguntas por módulo (tipo, obligatoriedad, condicionales), máscaras de documento, textos propios de enums, preguntas de Tamizaje, equipos de remisión y criterios, sectores y canales |
| `src/salvia/facades/RegistroV3Facade.go` | `RegistroV3GET`: valida `manage_registro_v3` y renderiza la página (`MainRouter.go:230`) |
| `src/salvia/service/victim_case_v3_draft.go` · `victim_case_v3_tasks.go` *(WP-4)* | Borrador JSON de entidades y barreras, vaciado al iniciar el Caso, y tareas, oficios y remisiones que se crean entonces |
| `scripts/gen_registro_v3_barrera_catalog.py` *(WP-4)* | Genera `BARRERA_V4_V3`, `ENTIDAD_V4_V3` y `CANAL_ACTIVACION_V3` desde el seed de Seguimiento v4 |
| `src/salvia/controller/victim_case_v3_controller.go` · `service/victim_case_v3_service.go` · `internal/repository/victim_case_v3_repository.go` · `victim_case_v3_fields.go` | API `/api/v1/victim-cases-v3/*` (ver [`registro-v3-index.md`](registro-v3-index.md#api-que-usa-la-pantalla)) |
| `src/salvia/controller/entity_branch_api_controller.go` · `location_controller.go` · `form_controller.go` *(compartidos)* | Cascada de ubicación, sedes de entidades, obligaciones y competencias que usa GS § C/D |
| `src/frontend/templates/layouts/standard_scripts.html` · `overlay.html` *(compartidos)* | Scripts y popups globales |

---

## Árbol de interfaz

```
RegistroV3Screen  (registro_v3.html)  <main.salvia-main> .container
│
├── Titulo  "Registro de víctima"  (.welcome-section)
│
├── [etapa === 'habeas']
│   HabeasData
│   ├── Titulo  "Autorización de datos personales"  (.tituloform)
│   ├── Texto × 3  (Ley 1581/2012 y Resolución 772/2024 · derechos del titular y quejas ante la SIC · activación de rutas y remisión a entidades competentes)
│   ├── BtnAcepto  "Acepto y continuar"  → aceptarHabeasData()   (etapa = 'tipoContacto')
│   └── BtnNoAcepto  "No acepto"  → rechazarHabeasData()          (etapa = 'formularioCerrado')
│
├── [etapa === 'formularioCerrado']
│   Aviso  (.alert-warning)  "La persona no autorizó el tratamiento de sus datos personales. No se crea ningún registro ni se guarda ninguna respuesta."
│
├── [etapa === 'tipoContacto']
│   TipoContacto  "Tipo de contacto"
│   ├── Lista (.list-group)  → seleccionarTipoContacto(tipo)
│   │   ├── "Víctima - Caso Nuevo"            (victima_caso_nuevo)  → etapa = 'cabezote'
│   │   ├── "Víctima - Seguimiento"           (victima_seguimiento)  ┐
│   │   ├── "Tercero - Reporte"               (tercero_reporte)      │ → tipoContactoNoDisponible = true
│   │   ├── "Familiares (feminicidio)"        (familiares)           │
│   │   └── "Hombres"                         (hombres)              ┘
│   └── [v-if tipoContactoNoDisponible]
│       Aviso  (.alert-info)  "Este tipo de contacto se gestiona en otra pantalla del sistema, fuera de Registro V3. Cierra este formulario y continúa desde el menú principal."
│
├── [etapa === 'cabezote']
│   Cabezote
│   ├── Titulo  "Buscar registro existente"  +  Nota  "Búsqueda únicamente por documento y/o teléfono (nunca por nombre)."
│   ├── <input>  buscar.documento  "Documento"
│   ├── <input>  buscar.telefono   "Teléfono"
│   ├── BtnBuscar  "Buscar"  :disabled="buscando"  → ejecutarBusqueda()
│   ├── [v-if errorBusqueda]  Aviso (.alert-danger)  ${errorBusqueda}      ← falla al buscar, al crear el borrador o al abrir un registro ("No se pudo abrir el registro…", flow-E21)
│   ├── [v-if resultadosBusqueda.length]
│   │   Tabla  (.table)  columnas: Nombre · Documento · Teléfono · Estado · Registro · (acción)
│   │   └── Fila × N  [v-for resultadosBusqueda]
│   │       ├── ${names} ${lastNames}
│   │       ├── ${docType}-${docNumber}
│   │       ├── ${phone}
│   │       ├── ${status}                                   ← código crudo del Caso: bo / in / te …
│   │       ├── Badge  [Completo (.bg-success) | Incompleto (.bg-secondary)]  ← registroCompleto
│   │       └── BtnAbrir  "Abrir"  → abrirCaso(icode)
│   ├── [v-else-if buscando]  "Buscando..."
│   └── BtnNuevoRegistro  "+ Nuevo registro"  → crearNuevoRegistro()    ← siempre visible, sin condición
│
├── [etapa === 'formulario']
│   Formulario
│   ├── [v-if resumenCaso]
│   │   Encabezado  "${nombreEncabezado()} — ${documentoEncabezado()}"  +  Badge (.bg-info) ${etiquetaEstado(status)}  +  Barra de avance global (.progress) "Avance N %"   ← sin nombre: "Persona sin identificar"; documento provisional (0000000000): "Documento sin registrar"; tras cada guardado se vuelve a leer el resumen del Caso (WP-3 · REG-03/04/02)
│   ├── [v-if triadaCompleta]
│   │   Aviso  (.alert-success)  "La información mínima de contacto, víctima y hecho está completa — el Caso pasó a Iniciado."
│   ├── BtnFinalizarLlamada  "Finalizar llamada (guardar y salir)"  :disabled="guardando"  → guardarYSalir()
│   ├── Menú de módulos  (ul.nav.nav-pills.flex-column.v3-menu, columna izquierda)  × N  [v-for modulos]   ← un solo módulo visible a la vez (moduloActivo); se construye con modulos[] (sin módulos escritos a mano)  (WP-3 · REG-01)
│   │   CI "Contacto e identificación" · DP "Datos personales" · CA "Contacto de apoyo" · HV "Hechos de violencia"
│   │   AG "Características del agresor" · TZ "Tamizaje de riesgo" · RA "Ruta de atención integral" · GS "Plan de acción Salvia"   ← 8 módulos (WP-4: RA nuevo, GS renombrado)
│   │   ├── Ítem: etiqueta + [TZ con !triadaCompleta] " 🔒" y texto atenuado + porcentaje de avance "NN %" ("—" si el módulo no tiene obligatorias de 1er contacto: CA y TZ)  → irAModulo(id) (no navega si está bloqueado)
│   │   └── Solo TZ se bloquea (REG-13): la Ruta de atención integral y el Plan de acción se pueden abrir desde el Borrador (GS_1, GS_4, GS_5, GS_14, GS_16, GS_17 y GS_18 siguen siendo obligatorias al finalizar)
│   ├── [v-if mensajeModuloBloqueado]  Aviso (.alert-warning)  "Completa primero los módulos de Contacto, Datos personales, Contacto de apoyo, Hechos y Agresor para habilitar Tamizaje."
│   ├── BarraNavegacion  (.v3-nav)  BtnAnterior "← Anterior" → irAAnterior() · "Módulo i de N" · BtnSiguiente "Siguiente →" → irASiguiente() · [v-if mensajeAvance] "✓ Sección guardada. Pasaste de «X» a «Y»."   ← saltan los módulos bloqueados y NO guardan (E22/E23, WP-3 · REG-01)
│   ├── [v-if errorGuardado]  Aviso (.alert-danger #registro-v3-error-guardado)  ${errorGuardado}   ← un guardado falló (dato inválido con la pregunta a corregir, error del servidor, sin conexión, sesión vencida); se hace scroll hasta él (flow-E21)
│   │
│   ├── [moduloActivo ∈ CI, DP, CA, HV, AG]
│   │   ModuloGenerico
│   │   ├── Campo × N  [v-for camposVisibles(moduloActivo)]        ← visibilidad por condicionales del catálogo
│   │   │   ├── <label>  ${campo.label}
│   │   │   ├── [tipo 'text']      <input>  (.is-invalid si soloLetras con solo números, o CI_6 fuera de máscara)
│   │   │   │   ├── [soloLetras inválido]  "No puede ser solo números."
│   │   │   │   ├── [CI_6]  "Formato esperado: …"  y  "No cumple el formato esperado."
│   │   │   │   └── [CI_6 y avisoDuplicidad]  Aviso (.alert-warning)  "Ya existe un Registro con este documento: …"
│   │   │   ├── [tipo 'textarea']  <textarea>
│   │   │   ├── [tipo 'date']      <input type="date">
│   │   │   ├── [tipo 'time']      <input type="time" step="60">   ← HV_3 «¿A qué hora ocurrió el hecho?» (REG-09)
│   │   │   ├── [tipo 'dropdown']  <select>  (opciones fijas, p. ej. departamentos, o enum del catálogo; DP_30 zona de residencia = Rural | Urbano; DP_51 «Tipo de discapacidad» = 7 opciones, solo si DP_2 = Sí)
│   │   │   ├── [tipo 'boolean']   Radio Sí / No
│   │   │   └── [tipo 'multiple']  Checkbox × N  (con tope `max`, p. ej. HV_11 = 2)
│   │   ├── BtnGuardar  "Guardar sección"  :disabled="guardando"  → guardarModulo(moduloActivo)
│   │   └── Mensaje  (.ms-2  :class="{'v3-msg-error': mensajeGuardadoEsError}")  ${mensajeGuardado}      ← "Sección guardada.", las validaciones del cliente o, si el guardado falló, "No se guardó." (el detalle sale en el aviso rojo de arriba)   ← los faltantes y errores salen en rojo (REG-05); si se guardó, se pasa al siguiente módulo no bloqueado (REG-01)
│   │
│   ├── [moduloActivo === 'TZ']
│   │   Tamizaje
│   │   ├── Nota  "Primero se responden las 5 preguntas iniciales (marcadas «inicial»). Si alguna se responde "Sí", se habilita el resto del Tamizaje; si todas son "No", ya se puede calcular y guardar."
│   │   ├── TablaGenerales  (<table> Pregunta | Sí | No · colgroup 70/15/15 · <caption> oculto)   (REG-12)
│   │   │   └── Fila × 7  [v-for TAMIZAJE_GENERALES_V3]  Pregunta + badge "inicial" en TZ_1, 2, 3, 4, 7  ·  2 radios (mismo name, aria-label "Sí — …" / "No — …")  → tzToggle()
│   │   │       └── TZ_5 y TZ_6 solo visibles con tzRevelarResto
│   │   ├── [v-show tzRevelarResto]
│   │   │   ├── Nota  "Rama detectada automáticamente según AG — Características del agresor: Pareja/expareja | No pareja"
│   │   │   └── TablaResto  (<table> Pregunta | Sí | No · colgroup 70/15/15)
│   │   │       ├── Fila  TZ_8 (VPP, prejuicio)       ← se muestra siempre que se revele el resto (no depende de esOsigd)
│   │   │       └── Fila × 24  [v-for rama pareja (TZP_1..24) o no pareja (TZNP_1..24)]   ← las 24, sin ocultar condicionales
│   │   ├── [v-if tzInicialesRespondidas && !tzRevelarResto]  Nota  "Las cinco preguntas iniciales están en «No»: no se pregunta el resto del Tamizaje. Ya puedes calcular y guardar el resultado."   (REG-34)
│   │   ├── BtnCalcular  "Calcular y guardar Tamizaje"  [v-show tzRevelarResto || tzInicialesRespondidas]  :disabled="tzGuardando"  → guardarTamizaje()
│   │   ├── [v-if tzResultado]
│   │   │   Resultado  (.alert-danger si riesgo alto/extremo, si no .alert-info)
│   │   │   "Tipo de riesgo: {feminicidio|homicidio_abyecto} — Nivel: {bajo|moderado|alto|extremo} — Puntaje: N"  [+ nota de empate]
│   │   │   └── [v-if riesgoEsAlto()]  "⚠️ Nivel de riesgo alto — considere activar la Alerta Inminente en Gestión Salvia."
│   │   └── [v-if hayIdeacionSuicida()]  Aviso (.alert-danger)  "⚠️ Se marcó ideación suicida. Debe abordarse con remisión psicosocial (Gestión Salvia § B, Atención Psicosocial)."
│   │
│   ├── [moduloActivo === 'RA']   (WP-4 · REG-17)
│   │   RutaAtencionIntegral   ← informativo: a qué rutas o entidades acudió la persona y si encontró barreras; replica «Identificación de entidades» de Seguimiento v4 sin el componente dinámico (textos en `ENTIDAD_V4_V3`, generado del seed)
│   │   ├── Texto guía  +  Aviso de borrador [v-if !entitiesReadOnly] (.alert-warning) "Borrador: las entidades, las barreras y la remisión… se crean cuando el Caso pase a Iniciado…" | [v-else] (.alert-info) "El Caso ya dejó de ser borrador: … solo lectura"   ← fragmento `rv3_aviso_borrador`
│   │   ├── Radio  GS_14  "¿Ha puesto en conocimiento de alguna entidad o institución los hechos de violencia reportados? Sí | No"
│   │   ├── [v-if GS_14 === 'true' || (GS_17 === 'true' && !entidadesAcudidas.length)]
│   │   │   ├── TarjetaEntidad × N  [v-for entidadesAcudidas]   (fragmento `rv3_tarjeta_entidad`, igual que en el Plan de acción)
│   │   │   └── [v-if !entitiesReadOnly]  FormularioEntidadAcudida  (.card)  "Entidad a la que ya acudió"
│   │   │       ├── BuscadorEntidad  (fragmento `rv3_entidad_buscador`, REG-19): Departamento\* → Ciudad\* → Municipio\* → Sector\* (selects en cascada, cada uno deshabilitado hasta elegir el anterior) → «Buscar entidad»\* (sugerencias; la lista de sedes se pide UNA vez al elegir municipio y sector, el filtro es local, sin tildes, con debounce de 250 ms) → chips «Acciones institucionales» Completadas / Pendientes (GET /api/v1/entity-branches/:id/obligations; sin «otra»; una acción no puede estar en las dos)
│   │   │       ├── [v-if EntityBranchID]  «Número de radicado (NUNC, medida de protección). Escriba tipo y número»\*  ·  «Fecha en la que reportó a la entidad»\* (date)  ·  «Otra información asociada a la Ruta»\* (→ objetivo)
│   │   │       ├── [v-if PendingActions.length]  (D17)  «¿La persona está de acuerdo con que SALVIA active la Ruta de Atención y las gestiones institucionales correspondientes?»\* Sí | No → [Sí] «¿Salvia debe activar rutas en el marco de la debida diligencia?»\* → [ambas Sí] «Seleccione los canales utilizados para la activación o gestión de la Ruta de Atención»\* (5 casillas; «Ninguna» excluyente) + «¿Qué gestiones se realizaron desde Salvia para la activación de ruta?»\*
│   │   │       ├── [v-if entidadFaltantes().length]  "Falta: …"   ·   [v-if errorEntidad]  texto de error
│   │   │       └── BtnAgregarEntidad  "+ Siguiente entidad"  :disabled="!EntityBranchID || entidadFaltantes().length"  → agregarEntidad()
│   │   ├── Radio  GS_17  "¿Se identificaron barreras? Sí | No"   ·   [v-if GS_17 === 'true']  Nota "Registra cada barrera con el botón «+ Agregar barrera a esta entidad» de la tarjeta de la entidad donde ocurrió." (+ " Primero agrega la entidad." si no hay entidades)
│   │   └── BtnGuardar  "Guardar sección"  → guardarRutaAtencion()   ← valida GS_14 y GS_17 y envía GS con enviarSeccion('GS'); no avanza solo (la barra «Siguiente» sí funciona)  +  Mensaje ${mensajeGuardado}
│   │
│   └── [moduloActivo === 'GS']
│       PlanDeAccionSalvia   (WP-4 · REG-15: antes «Gestión Salvia»; orden a–i del cliente, sin los títulos A–D)
│       ├── a) [v-if parrafoSintesis()]  Aviso (.alert-secondary)  "Párrafo síntesis: …"       ← recalculado en cada render, no se guarda
│       ├── Aviso de borrador / solo lectura  (fragmento `rv3_aviso_borrador`)
│       ├── b) GS_1 (boolean)       ← fragmento `rv3_gs_campos` con `camposGs('GS_1')`
│       ├── c) Radio  GS_15  "¿Está de acuerdo con que SALVIA realice la activación de la Ruta de Atención? Sí | No"   ← siempre visible (antes solo con GS_14 o GS_16)
│       ├── d) Radio  GS_16  "¿Salvia debe activar rutas en el marco de la debida diligencia? Sí | No"
│       │   ├── [v-if entidadesActivacion.length]  TarjetaEntidad × N  [v-for entidadesActivacion]   ← las activaciones guardadas se ven siempre (la acudida solo aparece en la Ruta)
│       │   │   (fragmento `rv3_tarjeta_entidad`)  Badge sector + nombre de la sede + ubicación · «Canales de activación: …» (de `entityMeta`; los datos anteriores «Canal de activación: …» en `solicitud` se siguen viendo) · «Acciones institucionales:» chips «✓ Completada» / «Pendiente» (ids en `entityMeta`, etiquetas del catálogo de la sede)
│       │   │   ├── [v-if !entitiesReadOnly]  BtnEliminarEntidad  "Eliminar entidad"  → eliminarEntidad(relId)  (window.confirm; en el borrador se eliminan también sus barreras)
│       │   │   ├── [v-if barrerasDeEntidad(e).length]  Tabla  Descripción · Gestión (`etiquetaGestionBarrera`: etiqueta, no código; «—» si no hay) · «Eliminar» [v-if !entitiesReadOnly] → eliminarBarrera(id)
│       │   │   ├── [v-if GS_17 === 'true' && !entitiesReadOnly && no hay formulario abierto]  BtnAgregarBarrera  "+ Agregar barrera a esta entidad"  → abrirBarreraParaEntidad(e)   ← también en la tarjeta de las acudidas de la Ruta (D8, mientras REG-29 no se defina)
│       │   │   └── [v-if barreraParaEntidad === e.entityBranchId && !entitiesReadOnly]  FormularioBarrera  (ver abajo)
│       │   └── [v-if GS_16 === 'true' && !entitiesReadOnly]  FormularioActivacion  (.card)  «Activación de ruta»  (tipo fijo `activacion`)
│       │       ├── BuscadorEntidad  (fragmento `rv3_entidad_buscador`; solo acciones pendientes)
│       │       ├── [v-if EntityBranchID]  «Solicitud a la entidad» (textarea)  ·  «¿Cuenta con un N° de radicado, NUNC, MP, etc.?»  ·  «Objetivo con la entidad»\* (textarea, obligatorio)
│       │       ├── [v-if EntityBranchID && GS_15 === 'true' && GS_16 === 'true']  «Seleccione los canales utilizados…»\* (5 casillas: Notificación SALVIA · Articulación institucional telefónica o por mensajería · Denuncia fácil · Mi Salvia · Ninguna, excluyente)  +  «¿Qué gestiones se realizaron desde Salvia para la activación de ruta?»\*   |  [v-else]  nota «Los canales de activación se piden cuando…»
│       │       └── [v-if falta algo] "Falta: …" · [v-if errorEntidad] error · BtnAgregarEntidad "+ Siguiente entidad" → agregarEntidad()   (envía `Agreement`/`DueDiligence` como copias de GS_15 y GS_16)
│       ├── e) GS_3 (textarea)      ← `camposGs('GS_3')`
│       ├── f) Remisión a otros equipos de Salvia   (todos los campos `:disabled="entitiesReadOnly"`; GS_5 solo si ya está respondida, para no dejar sin salida a Finalizar)
│       │   ├── Radio  GS_5  "Requiere remisión: Sí | No"
│       │   └── [v-if GS_5 === 'true']
│       │       ├── Equipo × 6  [v-for GS_6]  Checkbox  Medidas de Emergencia · Atención Psicosocial · Estabilización · Atención Hombres · Salvia Dignidad · Discapacidad
│       │       │   ├── [equipo marcado con criterios]  Checkbox × N de criterios (GS_9..GS_12)  +  "✓ Cumple criterios de remisión" | "✗ Aún no cumple los criterios mínimos"
│       │       │   └── [equipo 'me' marcado]  Checkbox × 6 de medidas de emergencia (GS_7)
│       │       └── g) <textarea>  GS_13  "Justificación de la remisión"
│       ├── h) GS_20 (multiple)  "Seleccione los elementos psicosociales integrados en este acompañamiento"  (14 opciones del cliente + «Ninguno», excluyente en ambos sentidos; deseable)   ·   i) GS_4 (textarea)      ← `camposGs('GS_20,GS_4')`
│       ├── § E  "Alerta Inminente"  (sin cambios)
│       │   ├── Radio  GS_18  "¿Activar alerta inminente? Sí | No"
│       │   └── [v-if GS_18 === 'true']  <textarea>  GS_19  "Justificación de la alerta"
│       ├── BtnGuardar  "Guardar sección"  → guardarModulo('GS')   +  Mensaje ${mensajeGuardado}
│       ├── [v-if errorFinalizar]  Aviso (.alert-danger)  ${errorFinalizar}
│       └── BtnFinalizar  "Finalizar registro y enrutar a Seguimiento"  → finalizarRegistro()   ← primero guarda GS (mismas validaciones de «Guardar sección»), luego valida en el servidor y, si falta algo, lista los campos (REG-37); un 401 se avisa sin perder lo escrito (REG-21)
│
│   FormularioBarrera  (WP-4 · REG-18; fragmento `rv3_form_barrera`, dentro de la TarjetaEntidad)   ← se pinta recorriendo `BARRERA_V4_V3` (generado del seed de Seguimiento v4); sin competencias
│       ├── «Sector en el que se presentó la barrera»\*  [Salud | Justicia | Protección | Transversales]
│       ├── [sector Salud | Justicia | Protección]  texto informativo del sector · «Barreras en …» (casillas) · «¿A qué institución acudió?» (select) · «Otra barrera en sector …» (texto)
│       ├── [sector Transversales]  4 listas de casillas: institucionales y de talento humano · económicas y socioeconómicas · territoriales y geográficas · por ausencia de enfoque diferencial
│       ├── «Fecha en la que se presentó la barrera (aproximada)»\*  ·  «Funcionario/a o dependencia donde se presentó la barrera»\*  ·  «Departamento»\* → «Municipio»\* (carga los municipios del departamento)  ·  «Vereda o zona»\* [Cabecera | Vereda o zona rural | No precisa]  ·  «Descripción de la barrera (tiempo, modo y lugar)»\*
│       ├── «Gestión de la barrera»\*  [Activación de ruta interinstitucional (Generar oficio) | Alerta por barreras]  → [con «Alerta»] «¿Cuál es su perspectiva profesional frente a la barrera?»\* + aviso «Alerta por barreras notificaría al equipo técnico…»
│       ├── «Descripción de la gestión realizada con relación a las barreras»\*
│       ├── [v-if errorBarrera] error  ·  BtnCancelar → cerrarBarreraParaEntidad()  ·  BtnAgregarBarrera "+ Agregar barrera"  :disabled="barreraFaltantes().length"  → agregarBarrera()   (cambiar el sector o el departamento limpia lo que dejó de aplicar)
│
└── [etapa === 'cierre']
    Cierre
    ├── Aviso  (.alert-success)  "Registro finalizado."
    ├── [v-if mensajeCierre]  Aviso (.alert-info)  "Se solicitó la Denuncia fácil ante la FGN." · "Se confirma remisión a: …"
    ├── Tarjeta × 4  (.card)  Estado ("Terminado") · Nivel de riesgo (tzResultado.riskLevel | "—"; rojo si alto/extremo) · Entidades (n) · Barreras (n)
    ├── "Párrafo síntesis"  ${parrafoSintesis()}
    ├── "Tabla de actuaciones"  Tipo · Detalle   (una fila por entidad "nombre (sector)" + " — Canales: …" si tiene canales de activación, y por barrera "descripción")  |  "Sin entidades ni barreras registradas."
    ├── "Bitácora"  ${GS_4}  |  "Sin acompañamiento registrado."
    └── BtnVolver  "Volver al Cabezote"  → volverACabezote()
```

---

## Entrada desde un reporte (WP-1 T7)

`GET /salvia/registro-v3?caso=<icode>` abre **directo el formulario** (etapa `formulario`, módulo CI) del borrador creado desde un reporte: se omiten Habeas data, tipo de contacto y Cabezote porque el consentimiento ya se tomó en el reporte V2. Solo se repueblan las claves en `answeredKeys` (CI_2 nombres y apellidos, CI_4 teléfono, HV_1 hechos y, si el reporte tenía ajuste razonable, CI_8). Tipo y número de documento quedan **vacíos** (los borradores normales llevan el temporal `cc` / `0000000000`) hasta que el agente complete CI_5 y CI_6 — así `fanOutPerfilV3` no replica el Perfil entre borradores sin documento; CI_8 depende de CI_7 (no se marca CI_7). Sin `?caso=` (o con un valor que no tenga forma de UUID) el flujo normal arranca en la autorización. Ver [`Flujos/flow-E19-cuando-carga-con-registro-de-origen.md`](Flujos/flow-E19-cuando-carga-con-registro-de-origen.md).

## Estados y reglas

| Elemento | Valores / regla |
|---|---|
| `etapa` | `habeas` → `tipoContacto` → `cabezote` → `formulario` → `cierre`; ramas `formularioCerrado` (no aceptó) y `tipoContactoNoDisponible` (aviso dentro de `tipoContacto`) |
| Estado del Caso (`status`) | `bo` Borrador (al crear) → `in` Iniciado (cuando `triadaCompleta`) → `te` Terminado (al finalizar). El Cabezote y el encabezado muestran el **nombre** (Borrador, Iniciado, Terminado, Cerrado…; 11 códigos, "Desconocido" para uno nuevo), nunca el código crudo (REG-04). `cd` (Cerrado) lo pone Seguimiento |
| `triadaCompleta` | Las 32 obligatorias [1er contacto] de CI, DP, HV y AG (`triadaFieldsV3`: CI_1-7 y CI_9, DP_1, 2, 5, 8, 9, 11, 12, 14, 17, 27-29, HV_1, 2, 4, 6, 8, 9, 11, 13, 14, AG_1-3) ya respondidas alguna vez. CA no aporta. Lo calcula el backend (`SaveSection` y `GetCase`) y desbloquea TZ y GS |
| Badge del Cabezote | "Completo" si el Caso más reciente con ese documento tiene respondidas las obligatorias de CI + DP (+ DP_10, DP_15) + CA_1-4; si no, "Incompleto". Se calcula, no se guarda |
| Borrador de entidades y barreras (WP-4 · REG-38\*) | Con el Caso en Borrador (`bo`) las entidades, las barreras y los datos de Seguimiento v4 sin columna se guardan en `victim_case.victim_case_v3_draft_data` (ids `d~<icode>~<uuid>`); al pasar a Iniciado (`in`) se materializan en `entity_case`, `entity_case_obligation` y `barrier_v2` y se crean las tareas, los oficios y las remisiones. `GetCase` devuelve `entityMeta`, `barrierMeta` y `entitiesReadOnly` |
| Solo lectura (WP-4 · D20) | Con `entitiesReadOnly` (el Caso ya no es Borrador) no hay formulario de entidad, ni «Eliminar», ni «+ Agregar barrera», y los campos de remisión (GS_6–GS_13) quedan deshabilitados; GS_5 solo si ya estaba respondida; GS_14–GS_17 siguen respondibles (Finalizar las exige). El backend responde 409 a las altas y bajas |
| Bloqueo del menú | Solo TZ, con `!triadaCompleta` (REG-13). `triadaCompleta` no baja de true a false en pantalla y el servidor la devuelve real para cualquier módulo (REG-35) |
| Obligatoriedad por pregunta | `C1` = [1er contacto] (bloquea "Guardar sección" si está visible y vacía), `C2` = [máx. 2º contacto] (se exige solo al finalizar, y solo las incondicionales), `D` = Deseable |

## Planeado y no implementado (diferencias con la planeación original)

| Planeado | Qué hay en el código |
|---|---|
| Componentes `habeas-data-gate`, `cabezote-busqueda`, `dinamic-form`, `field-renderer` y CSS `victim-case-v3.css` | Todo en `registro_v3.html` + `registro-v3-app.js` + `registro-v3-campos.js`; sin CSS propio |
| Letreros "Debida diligencia" y "Respeto OSIGD" con texto verbatim del mockup | No existen; la etapa Habeas solo trae 3 párrafos legales |
| Selector de Tipo de contacto como `<select>` con 3 opciones que redirige | 5 botones; solo "Víctima - Caso Nuevo" continúa, las otras 4 muestran un aviso (no hay URL resuelta) |
| Cabezote con banner, botón "Limpiar", campo de tipo de documento, y resultados agrupados en Reportes pendientes / Registros con Casos / Borradores | Dos inputs, un botón "Buscar" y una tabla plana; **no** consulta reportes (`victim_contact`); sin "Limpiar"; el filtro por tipo de documento existe en la API pero no tiene input |
| "Proceder al formulario", "Nuevo Caso en Registro existente" (E-08), "Actualizar Caso" (E-09) y checkbox "Verifiqué las coincidencias" | No existen: el único acceso a un Caso existente es el botón "Abrir"; "+ Nuevo registro" siempre está habilitado |
| Encabezado del Formulario con nombre identitario prevalente, badge de Registro, alertas (duplicidad, riesgo, ideación, gestación) y botón "⚠ Alerta inminente" | Una línea con tratamiento + nombres + documento + badge de estado crudo; el aviso de duplicidad sale bajo CI_6 y los de riesgo/ideación dentro de TZ; la Alerta Inminente son los campos GS_18/GS_19 (decisión C7) |
| Sidebar con 3 grupos, barras de progreso por módulo, checkbox "solo obligatorias de 1er contacto" y grupo "Registro de gestión" (Rutas/Barreras/Notificaciones) | Menú vertical de módulos con porcentaje de avance por módulo y global (WP-3 · REG-01/02); sin grupos, sin filtro y sin "Registro de gestión" (E-16 fuera de alcance) |
| Módulo DP con 9 preguntas de escala de dificultad, y marca de vulnerabilidad por gestación | DP_3 es un checklist de actividades; DP_26 (gestación) se muestra siempre (`__gestacionAplica` devuelve `true`) y no genera marca |
| Tamizaje: pregunta VPP solo si `esOsigd`; preguntas OSIGD y de ámbito laboral ocultas según el caso (`TAMIZAJE_*_OSIGD_V3`, `*_LABORAL_V3`) | TZ_8 se muestra siempre tras revelar el resto y las 24 preguntas de la rama aparecen todas: las constantes de filtro existen pero ningún render las usa. `esOsigd` solo decide el tipo de riesgo |
| Alertas de riesgo enviadas a coordinación (`NotificationService`) | Solo avisos en pantalla; no hay notificación ni log |
| Pantalla de Cierre con 4 `StatCard` de Registro/Caso/Tamizaje/Origen, tabla de actuaciones con responsable/fecha/estado, bitácora, remisiones y 3 botones | 4 tarjetas distintas (Estado, Nivel de riesgo, Entidades, Barreras), tabla simple (Tipo/Detalle), bitácora = GS_4 y un solo botón "Volver al Cabezote" (la ruta de Seguimiento exige un `follow_up_v2` que aún no existe) |
| Repeaters de Entidades y Barreras independientes (`dinamic-form`) | Barreras anidadas dentro de cada Entidad (decisión de 2026-10-15); Entidades en `entity_case`/`entity_case_obligation` y Barreras en `barrier_v2` |
