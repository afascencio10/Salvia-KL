# `/salvia/registro-v3` — Interfaz de la Pantalla

> ✅ **Sincronizado con el código el 2026-10-31** (versión implementada del 2026-10-14/15). La planeación original —con `<dinamic-form>`, sidebar con barras de progreso, componentes `habeas-data-gate`/`cabezote-busqueda` y el botón "Alerta inminente" del header— **no se construyó así**: la decisión A1 (ver [`registro-v3-data.md`](registro-v3-data.md) § D10) descartó `dinamic-form` y todo quedó en una app Vue a medida. Lo planeado y no construido está en la sección "Planeado y no implementado" al final.

Pantalla de ingreso y registro de hechos de VBG/VPP para los roles `op` y `sv` (permiso `manage_registro_v3`). Recorre seis etapas dentro de una sola página: autorización de datos personales → tipo de contacto → Cabezote de búsqueda de duplicidad → Formulario de siete módulos (CI, DP, CA, HV, AG, TZ, GS) → Cierre. No hay entrada en el menú lateral: se llega por la URL. Guarda por módulo en `/api/v1/victim-cases-v3/*` y escribe directo sobre `victim_case`/`victim_case_form2`, las mismas tablas que lee Seguimiento.

## Archivos relevantes

| Archivo | Rol |
|---|---|
| `src/frontend/html/salvia/registro_v3/registro_v3.html` | Template único: las 6 etapas con `v-if="etapa===…"` |
| `src/frontend/js/components/registro-v3-app.js` | App Vue 3 `RegistroV3App` (Options API, delimitadores `${ }`): estado, validaciones, llamadas a la API, párrafo síntesis |
| `src/frontend/js/components/registro-v3-campos.js` | Catálogo de preguntas por módulo (tipo, obligatoriedad, condicionales), máscaras de documento, textos propios de enums, preguntas de Tamizaje, equipos de remisión y criterios, sectores y canales |
| `src/salvia/facades/RegistroV3Facade.go` | `RegistroV3GET`: valida `manage_registro_v3` y renderiza la página (`MainRouter.go:230`) |
| `src/salvia/controller/victim_case_v3_controller.go` · `service/victim_case_v3_service.go` · `internal/repository/victim_case_v3_repository.go` · `victim_case_v3_fields.go` | API `/api/v1/victim-cases-v3/*` (ver [`registro-v3-index.md`](registro-v3-index.md#api-que-usa-la-pantalla)) |
| `src/salvia/controller/entity_branch_api_controller.go` · `location_controller.go` · `form_controller.go` *(compartidos)* | Cascada de ubicación, sedes de entidades, obligaciones y competencias que usa GS § C/D |
| `src/frontend/templates/layouts/standard_scripts.html` · `overlay.html` *(compartidos)* | Scripts y popups globales |

---

## Árbol de interfaz

```
RegistroV3Screen  (registro_v3.html)  <main.salvia-main> .container
│
├── Titulo  "Registro V3"  (.welcome-section)
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
│   ├── [v-if errorBusqueda]  Aviso (.alert-danger)  ${errorBusqueda}
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
│   │   Encabezado  "${tratamiento()} ${names} ${lastNames} — ${docType}-${docNumber}"  +  Badge (.bg-info) ${status}
│   ├── [v-if triadaCompleta]
│   │   Aviso  (.alert-success)  "La información mínima de contacto, víctima y hecho está completa — el Caso pasó a Iniciado."
│   ├── BtnFinalizarLlamada  "Finalizar llamada (guardar y salir)"  :disabled="guardando"  → guardarYSalir()
│   ├── Tabs  (.nav.nav-tabs)  × 7  [v-for modulos]       ← un solo módulo visible a la vez (moduloActivo)
│   │   CI "Contacto e identificación" · DP "Datos personales" · CA "Contacto de apoyo" · HV "Hechos de violencia"
│   │   AG "Características del agresor" · TZ "Tamizaje de riesgo" · GS "Gestión Salvia"
│   │   └── [TZ, GS con !triadaCompleta]  texto atenuado + " 🔒" + title "Completa primero CI/DP/CA/HV/AG"  → irAModulo(id) no navega
│   ├── [v-if mensajeModuloBloqueado]  Aviso (.alert-warning)  "Completa primero los módulos de Contacto, Datos personales, Contacto de apoyo, Hechos y Agresor para habilitar Tamizaje|Gestión Salvia."
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
│   │   │   ├── [tipo 'dropdown']  <select>  (opciones fijas, p. ej. departamentos, o enum del catálogo)
│   │   │   ├── [tipo 'boolean']   Radio Sí / No
│   │   │   └── [tipo 'multiple']  Checkbox × N  (con tope `max`, p. ej. HV_11 = 2)
│   │   ├── BtnGuardar  "Guardar sección"  :disabled="guardando"  → guardarModulo(moduloActivo)
│   │   └── Mensaje  ${mensajeGuardado}
│   │
│   ├── [moduloActivo === 'TZ']
│   │   Tamizaje
│   │   ├── Nota  "Se evalúan primero 7 preguntas generales. Si alguna de las marcadas como obligatorias se responde "Sí", se habilita el resto del Tamizaje."
│   │   ├── PreguntaGeneral × 7  [v-for TAMIZAJE_GENERALES_V3]  Radio Sí/No  (badge "inicial" en TZ_1, 2, 3, 4, 7)  → tzToggle()
│   │   │   └── TZ_5 y TZ_6 solo visibles con tzRevelarResto
│   │   ├── [v-show tzRevelarResto]
│   │   │   ├── PreguntaVPP  TZ_8  Radio Sí/No            ← se muestra siempre que se revele el resto (no depende de esOsigd)
│   │   │   ├── Nota  "Rama detectada automáticamente según AG — Características del agresor: Pareja/expareja | No pareja"
│   │   │   ├── Pregunta × 24  [v-for rama pareja (TZP_1..24) o no pareja (TZNP_1..24)]  Radio Sí/No   ← las 24, sin ocultar condicionales
│   │   │   └── BtnCalcular  "Calcular y guardar Tamizaje"  :disabled="tzGuardando"  → guardarTamizaje()
│   │   ├── [v-if tzResultado]
│   │   │   Resultado  (.alert-danger si riesgo alto/extremo, si no .alert-info)
│   │   │   "Tipo de riesgo: {feminicidio|homicidio_abyecto} — Nivel: {bajo|moderado|alto|extremo} — Puntaje: N"  [+ nota de empate]
│   │   │   └── [v-if riesgoEsAlto()]  "⚠️ Nivel de riesgo alto — considere activar la Alerta Inminente en Gestión Salvia."
│   │   └── [v-if hayIdeacionSuicida()]  Aviso (.alert-danger)  "⚠️ Se marcó ideación suicida. Debe abordarse con remisión psicosocial (Gestión Salvia § B, Atención Psicosocial)."
│   │
│   └── [moduloActivo === 'GS']
│       GestionSalvia
│       ├── [v-if parrafoSintesis()]  Aviso (.alert-secondary)  "Párrafo síntesis: …"       ← recalculado en cada render, no se guarda
│       ├── § A  "Denuncia y gestión general"      Campo × N [v-for camposVisibles('GS')]: GS_1 (boolean) · GS_2 (multiple) · GS_3 (textarea) · GS_4 (textarea)
│       ├── § B  "Remisión a otros equipos de Salvia"
│       │   ├── Radio  GS_5  "Requiere remisión: Sí | No"
│       │   └── [v-if GS_5 === 'true']
│       │       ├── Equipo × 6  [v-for GS_6]  Checkbox  Medidas de Emergencia · Atención Psicosocial · Estabilización · Atención Hombres · Salvia Dignidad · Discapacidad
│       │       │   ├── [equipo marcado con criterios]  Checkbox × N de criterios (GS_9..GS_12)  +  "✓ Cumple criterios de remisión" | "✗ Aún no cumple los criterios mínimos"
│       │       │   └── [equipo 'me' marcado]  Checkbox × 6 de medidas de emergencia (GS_7)
│       │       └── <textarea>  GS_13  "Justificación de la remisión"
│       ├── § C  "Registrar entidades (acudidas / activación de ruta)"
│       │   ├── Radio  GS_14  "¿Ya se reportó a alguna entidad? Sí | No"
│       │   ├── Radio  GS_16  "¿Salvia debe activar rutas en el marco de la debida diligencia? Sí | No"
│       │   ├── [v-if GS_14 === 'true' || GS_16 === 'true']  Radio  GS_15  "¿Está de acuerdo con que SALVIA realice la activación de la Ruta de Atención? Sí | No"
│       │   ├── [v-if entidades.length]
│       │   │   TarjetaEntidad × N  [v-for entidades]  (.card)
│       │   │   ├── Badge sector  +  nombre de la sede  +  ubicación (departamento · ciudad · municipio)
│       │   │   ├── BtnEliminarEntidad  "Eliminar entidad"  → eliminarEntidad(relId)   (window.confirm "…También se eliminarán sus barreras asociadas." — el backend NO las elimina, ver GAP G-V3-22)
│       │   │   ├── [v-if barrerasDeEntidad(e).length]  Tabla  Descripción · Gestión · "Eliminar" → eliminarBarrera(id)  (window.confirm)
│       │   │   ├── [v-if GS_17 === 'true' && no hay formulario abierto]  BtnAgregarBarrera  "+ Agregar barrera a esta entidad"  → abrirBarreraParaEntidad(e)
│       │   │   └── [v-if barreraParaEntidad === e.entityBranchId]  FormularioBarrera  (ver § D)
│       │   └── BuscadorEntidad  (.card)
│       │       ├── <select>  Tipo  [Entidad acudida (acudida) | Activación de ruta (activacion)]
│       │       ├── <select>  Departamento → Ciudad → Municipio  (en cascada; GET /api/v1/locations/…)
│       │       ├── <select>  Sector  [Salud | Justicia | Protección | Otras instituciones]
│       │       ├── [v-if entidadesDisponibles.length]  <select>  Entidad  (GET /api/v1/entity-branches?town_code&sector)
│       │       ├── [v-if EntityBranchID]  <input> "N° radicado"  ·  [Kind = activacion] <select> "Canal de activación" [Notificación | Llamada | Mi Salvia]  ·  <input> "Objetivo con la entidad" (acudida) | "Información adicional" (activación)
│       │       ├── [v-if EntityBranchID]  Checklist "Acciones institucionales": Completadas (solo acudida) y Pendientes  (GET /api/v1/entity-branches/:id/obligations)
│       │       ├── [v-if errorEntidad]  texto de error
│       │       └── BtnAgregarEntidad  "+ Siguiente entidad"  :disabled="!EntityBranchID"  → agregarEntidad()
│       ├── § D  "Barreras institucionales"
│       │   ├── Radio  GS_17  "¿Se identificaron barreras? Sí | No"
│       │   ├── [v-if GS_17 === 'true' && !entidades.length]  Nota  "Primero agrega la entidad relacionada en la sección C -- las barreras se registran dentro de cada entidad."
│       │   └── FormularioBarrera (dentro de la TarjetaEntidad)
│       │       ├── [sector con catálogo propio]  Checkbox × N  "Barreras específicas del sector"
│       │       ├── [sector otras_instituciones]  <input>  "¿Cuál otra barrera?"
│       │       ├── Categoría transversal × 4  Checkbox × N  (institucionales y de talento humano · económicas · territoriales · por ausencia de enfoque diferencial)
│       │       ├── Checklist  "Competencias institucionales de la entidad"  (GET /api/v1/competences?sector=)
│       │       ├── <input type="date">  "Fecha en la que se presentó la barrera"  ·  <input>  "Funcionario(a) o dependencia donde se presentó"  ·  <textarea>  "Descripción de la barrera"
│       │       ├── Checklist  "Gestión realizada"  (6 opciones)
│       │       ├── [v-if errorBarrera]  texto de error
│       │       ├── BtnCancelar  "Cancelar"  → cerrarBarreraParaEntidad()
│       │       └── BtnAgregarBarrera  "+ Agregar barrera"  :disabled="!Description || !Date || !Official"  → agregarBarrera()
│       ├── § E  "Alerta Inminente"
│       │   ├── Radio  GS_18  "¿Activar alerta inminente? Sí | No"
│       │   └── [v-if GS_18 === 'true']  <textarea>  GS_19  "Justificación de la alerta"
│       ├── BtnGuardar  "Guardar sección"  → guardarModulo('GS')   +  Mensaje ${mensajeGuardado}
│       ├── [v-if errorFinalizar]  Aviso (.alert-danger)  ${errorFinalizar}
│       └── BtnFinalizar  "Finalizar registro y enrutar a Seguimiento"  → finalizarRegistro()   ← siempre habilitado; el backend decide
│
└── [etapa === 'cierre']
    Cierre
    ├── Aviso  (.alert-success)  "Registro finalizado."
    ├── [v-if mensajeCierre]  Aviso (.alert-info)  "Se solicitó la Denuncia fácil ante la FGN." · "Se confirma remisión a: …"
    ├── Tarjeta × 4  (.card)  Estado ("Terminado") · Nivel de riesgo (tzResultado.riskLevel | "—"; rojo si alto/extremo) · Entidades (n) · Barreras (n)
    ├── "Párrafo síntesis"  ${parrafoSintesis()}
    ├── "Tabla de actuaciones"  Tipo · Detalle   (una fila por entidad "nombre (sector)" y por barrera "descripción")  |  "Sin entidades ni barreras registradas."
    ├── "Bitácora"  ${GS_4}  |  "Sin acompañamiento registrado."
    └── BtnVolver  "Volver al Cabezote"  → volverACabezote()
```

---

## Estados y reglas

| Elemento | Valores / regla |
|---|---|
| `etapa` | `habeas` → `tipoContacto` → `cabezote` → `formulario` → `cierre`; ramas `formularioCerrado` (no aceptó) y `tipoContactoNoDisponible` (aviso dentro de `tipoContacto`) |
| Estado del Caso (`status`) | `bo` Borrador (al crear) → `in` Iniciado (cuando `triadaCompleta`) → `te` Terminado (al finalizar). El Cabezote y el encabezado muestran el **código crudo**, no una etiqueta. `cd` (Cerrado) lo pone Seguimiento |
| `triadaCompleta` | Las 32 obligatorias [1er contacto] de CI, DP, HV y AG (`triadaFieldsV3`: CI_1-7 y CI_9, DP_1, 2, 5, 8, 9, 11, 12, 14, 17, 27-29, HV_1, 2, 4, 6, 8, 9, 11, 13, 14, AG_1-3) ya respondidas alguna vez. CA no aporta. Lo calcula el backend (`SaveSection` y `GetCase`) y desbloquea TZ y GS |
| Badge del Cabezote | "Completo" si el Caso más reciente con ese documento tiene respondidas las obligatorias de CI + DP (+ DP_10, DP_15) + CA_1-4; si no, "Incompleto". Se calcula, no se guarda |
| Bloqueo del sidebar | Solo TZ y GS, con `!triadaCompleta` |
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
| Sidebar con 3 grupos, barras de progreso por módulo, checkbox "solo obligatorias de 1er contacto" y grupo "Registro de gestión" (Rutas/Barreras/Notificaciones) | Pestañas horizontales sin progreso ni filtro; sin "Registro de gestión" (E-16 fuera de alcance) |
| Módulo DP con 9 preguntas de escala de dificultad, y marca de vulnerabilidad por gestación | DP_3 es un checklist de actividades; DP_26 (gestación) se muestra siempre (`__gestacionAplica` devuelve `true`) y no genera marca |
| Tamizaje: pregunta VPP solo si `esOsigd`; preguntas OSIGD y de ámbito laboral ocultas según el caso (`TAMIZAJE_*_OSIGD_V3`, `*_LABORAL_V3`) | TZ_8 se muestra siempre tras revelar el resto y las 24 preguntas de la rama aparecen todas: las constantes de filtro existen pero ningún render las usa. `esOsigd` solo decide el tipo de riesgo |
| Alertas de riesgo enviadas a coordinación (`NotificationService`) | Solo avisos en pantalla; no hay notificación ni log |
| Pantalla de Cierre con 4 `StatCard` de Registro/Caso/Tamizaje/Origen, tabla de actuaciones con responsable/fecha/estado, bitácora, remisiones y 3 botones | 4 tarjetas distintas (Estado, Nivel de riesgo, Entidades, Barreras), tabla simple (Tipo/Detalle), bitácora = GS_4 y un solo botón "Volver al Cabezote" (la ruta de Seguimiento exige un `follow_up_v2` que aún no existe) |
| Repeaters de Entidades y Barreras independientes (`dinamic-form`) | Barreras anidadas dentro de cada Entidad (decisión de 2026-10-15); Entidades en `entity_case`/`entity_case_obligation` y Barreras en `barrier_v2` |
