# `Registro V3` — Interfaz de la Pantalla

Pantalla de ingreso y registro de hechos de VBG/VPP: gate de habeas data →
triage de tipo de contacto → Cabezote de búsqueda de duplicidad → Formulario
de Ingreso Registro multi-módulo (`<dinamic-form>` re-seedeado con las
secciones de `registro-v3-secciones.md`). Reemplaza en el tiempo a
`set_victim_case_v2.html`, pero coexiste con él mientras haya Casos abiertos
con el formulario V2 (ver GAP G-V3-03).

## Archivos relevantes

| Archivo | Rol |
|---|---|
| `src/frontend/html/salvia/victim_case/set_victim_case_v3.html` *(NUEVO)* | Template anfitrión — monta Habeas data, Cabezote y `<dinamic-form>` |
| `src/frontend/js/components/habeas-data-gate.js` *(NUEVO)* | Bloque de autorización + tipo de contacto |
| `src/frontend/js/components/cabezote-busqueda.js` *(NUEVO)* | Buscador de duplicidad + 5 acciones |
| `src/frontend/js/components/dinamic-form.js` *(EXISTENTE)* | Motor de formulario DB-driven — re-seedeado, no modificado en su lógica |
| `src/frontend/js/components/field-renderer.js` *(EXISTENTE, si aplica)* | Render por tipo de pregunta dentro de `dinamic-form` |
| `src/frontend/css/victim-case-v3.css` *(NUEVO)* | Estilos de Habeas data, Cabezote, sidebar de módulos y badges |
| `src/cmd/seed/seed_registro_v3.sql` *(NUEVO)* | Seed de `form`/`form_section`/`question`/`option`/`visibility_condition` para V3 |

---

## Árbol de interfaz

```
RegistroV3Screen  (set_victim_case_v3.html)
│
├── [stage === 'habeasData']
│   └── HabeasDataGate  // components/habeas-data-gate.js
│       ├── LetreroAgente  "Debida diligencia" — "El Sistema SALVIA del Ministerio de Igualdad y Equidad le informa que,
│       │   como parte de su deber legal de protección de derechos y atención integral, después de esta atención la
│       │   entidad podrá realizar las acciones necesarias para dar trámite y seguimiento a la información
│       │   suministrada, y activar las rutas de atención y remitir la información a las entidades competentes
│       │   cuando se identifiquen situaciones que puedan requerir protección." (texto verbatim del mockup)
│       ├── LetreroAgente  "Respeto OSIGD" — "Garantice el reconocimiento y respeto de la identidad de género de la
│       │   persona usuaria mediante el uso prioritario de su nombre identitario y de los pronombres acordes con su
│       │   autoidentificación, registrando el nombre jurídico únicamente cuando sea necesario para fines
│       │   administrativos." (texto verbatim del mockup)
│       ├── BannerLegal  "Autorización para el tratamiento de datos..." (texto largo, Ley 1581/Res. 772)
│       ├── BtnSí  "Sí" → responderHabeasData(true)
│       └── BtnNo  "No" → responderHabeasData(false)
│
├── [stage === 'formularioCerrado']
│   └── MensajeCierre  "El formulario no puede continuar sin autorización."
│
├── [stage === 'tipoContacto']
│   └── TipoContactoBlock
│       └── <select>  (.sv-select)
│           ├── <option> "Víctima - Caso Nuevo"
│           ├── <option> "Víctima - Seguimiento"
│           └── <option> "Tercero - Reporte"
│           → seleccionarTipoContacto()
│
├── [stage === 'cabezote']
│   └── CabezoteBusqueda  // components/cabezote-busqueda.js
│       ├── Header  (.cb-header)  "Registros SALVIA · Formulario de Ingreso Registro"
│       ├── BannerRecordatorio  "Siempre verifica..."  (.cb-alert)
│       ├── InputsBusqueda  (.cb-inputs)
│       │   ├── <input>  placeholder="Número de documento"
│       │   └── <input>  placeholder="Número de teléfono"
│       │       (**sin campo de nombre** — corregido por diagrama del cliente: la búsqueda que determina
│       │       coincidencia es solo por documento y teléfono, nunca por nombre; ver registro-v3-data.md D7)
│       ├── BtnBuscar  "Buscar" → buscarEnCabezote()
│       ├── BtnLimpiar  "Limpiar" → limpiarBusqueda()
│       │
│       └── [v-if buscado]
│           ├── [v-if resultados.reportes.length]
│           │   └── ReporteResumen × N  [v-for resultados.reportes]
│           │       └── Badge "Pendiente de validación"  (solo lectura)
│           │
│           ├── [v-if resultados.registros.length]
│           │   └── RegistroResumen × N  [v-for resultados.registros]
│           │       ├── PerfilResumen  (.cb-perfil)  nombre · documento · teléfono
│           │       ├── BadgeEstadoRegistro  [Incompleto|Completo]  (calculado, no columna — ver registro-v3-data.md D1)
│           │       └── CasoResumen × N  [v-for registro.casos]
│           │           ├── Texto  "Caso {id}"  ·  BadgeEstadoCaso [Borrador|Iniciado|Terminado]
│           │           ├── Texto  resumen (tipo violencia · ámbito · relación agresor)
│           │           └── [v-if caso.estado === 'Iniciado']
│           │               BtnActualizarCaso  "Actualizar Caso" → actualizarCaso(registro, caso)
│           │
│           ├── [v-if resultados.registros.length]
│           │   └── BtnNuevoCasoEnRegistro  "Nuevo Caso en Registro existente" → nuevoCasoEnRegistro(registro)
│           │
│           ├── [v-if resultados.borradores.length]
│           │   └── BorradorResumen × N  [v-for resultados.borradores]
│           │       ├── Texto  teléfono · fecha de creación · canal
│           │       └── BtnProcederFormulario  "Proceder al formulario" → procederAlFormulario(borrador)
│           │
│           ├── [v-if resultados.registros.length]
│           │   └── <checkbox>  "Verifiqué las coincidencias: se trata de otra persona"  v-model="busqueda.otraPersona"
│           └── BtnNuevoRegistro  "+ Nuevo Registro" → crearNuevoRegistro()
│               :disabled="resultados.registros.length > 0 && !busqueda.otraPersona"
│
└── [stage === 'formulario']
    └── FormularioIngresoRegistro
        ├── Header  (.fr-header)
        │   ├── NombreIdentitarioPrevalente  "{tratamiento} {nombreIdentitario} ({NOMBRES_JURIDICOS})"
        │   ├── BadgeEstadoCaso  [Borrador|Iniciado|Terminado]
        │   ├── BadgeEstadoRegistro  [Incompleto|Completo]
        │   ├── BtnAlertaInminente  "⚠ Alerta inminente" → activarAlertaInminente()
        │   │   (proveniente del prototipo React, no de la fuente Drive — ver flow-E13, placement/criterios provisionales)
        │   └── [v-if sesion.alertaDup || tamizaje.nivel en ['Alto','Extremo'] || tamizaje.ideacion || sexo_gestacion]
        │       AlertasHeader
        │       ├── [v-if sesion.alertaDup] AlertaDuplicidad  "El documento coincide con el Registro {idRegistro}..." → BtnResolver
        │       ├── [v-if tamizaje.nivel en Alto/Extremo] AlertaRiesgo  "Alerta por riesgo {nivel}: {ACCION_NIVEL[nivel]}..."
        │       ├── [v-if tamizaje.ideacion] AlertaIdeacionSuicida  "Verifique si fue ideación o amenaza/intento..."
        │       └── [v-if DP.embarazo === 'Sí'] MarcaVulnerabilidad  "Marca de vulnerabilidad: víctima en estado de gestación"
        │           (agregado del mockup — no estaba en ningún archivo previo, complementa el motor de alertas de flow-E11)
        │
        ├── SidebarModulos  (.fr-sidebar)
        │   ├── <checkbox>  "Mostrar solo obligatorias de 1er contacto"  v-model="ses.soloC1"
        │   │   (filtra el `<dinamic-form>` a solo preguntas [1er contacto] — agregado del mockup, no contradice nada)
        │   ├── GrupoBloque  "Datos y Registro de la víctima"
        │   │   └── ModuloItem × 3  [v-for: CI, DP, CA]
        │   │       ├── BarraProgreso  :pct="progresoModulo(m)"
        │   │       └── :class="activo si moduloActual === m"  → cambiarModulo(m)
        │   │
        │   ├── GrupoBloque  "Datos del Caso de violencia"
        │   │   └── ModuloItem × 4  [v-for: HV, AG, TZ, GS]
        │   │       ├── BarraProgreso  :pct="progresoModulo(m)"
        │   │       ├── 🔒  [v-if !triadaCompleta && m en ['TZ','GS']]
        │   │       │   :disabled, tooltip "Completa Víctima + Hechos + Agresor para desbloquear"
        │   │       └── :class="activo"  → cambiarModulo(m)  [v-if habilitado]
        │   │
        │   └── GrupoBloque  "Registro de gestión"  (tomado del mockup original — omitido en la primera
        │       │                                     versión de este plan, agregado tras re-leer el PDF
        │       │                                     de observaciones del cliente, obs. #17)
        │       ├── BtnRutas  "↗ Rutas" → abrirRutas()
        │       ├── BtnBarreras  "↗ Barreras" → abrirBarreras()
        │       └── BtnNotificaciones  "↗ Notificaciones" → abrirNotificaciones()
        │           (placeholder — conexión real con GS sin precisar, ver GAP G-V3-12 y evento E16)
        │
        └── <dinamic-form>  // components/dinamic-form.js (EXISTENTE)
            :form-id="REGISTRO_V3_FORM_ID"  :section="moduloActual"
            │
            ├── [moduloActual === 'CI'] SectionCI
            │   └── FieldRenderer × N  [v-for preguntasDeSeccion('CI')]
            │       ├── [type === 'text']  <input type="text">
            │       ├── [type === 'dropdown']  <select>
            │       ├── [type === 'boolean']  BtnSí / BtnNo
            │       └── [type === 'multiple']  <checkbox> × N  [v-for opciones]
            │
            ├── [moduloActual === 'DP'] SectionDP  (misma estructura FieldRenderer × N)
            │   └── [v-if pregunta.tipo === 'escala']
            │       EscalaDificultad × 9  [v-for actividadesEscala]
            │       └── <select>  opciones: Sin dificultad | Alguna | Mucha | No puede
            │
            ├── [moduloActual === 'CA'] SectionCA  (FieldRenderer × N)
            │
            ├── [moduloActual === 'HV'] SectionHV  (FieldRenderer × N)
            │   └── [v-if pregunta.id === 'subtipo_violencia']
            │       <checkbox> × N  [v-for subtiposDisponibles(tipoSeleccionado)]
            │
            ├── [moduloActual === 'AG'] SectionAG  (FieldRenderer × N)
            │
            ├── [moduloActual === 'TZ'] SectionTamizaje
            │   ├── PreguntaGeneral × 7  [v-for preguntasGenerales]
            │   ├── [v-if esOsigd]
            │   │   └── PreguntaFuerteVPP  "¿...motivados por prejuicio...?"
            │   ├── [v-if wasPartner]
            │   │   └── PreguntaRamaPareja × 23  [v-for tamizajePareja]
            │   │       └── [v-if pregunta.condicionalOsigd && !esOsigd]  (oculta)
            │   ├── [v-else]
            │   │   └── PreguntaRamaNoPareja × 24  [v-for tamizajeNoPareja]
            │   │       ├── [v-if pregunta.condicionalOsigd && !esOsigd]  (oculta)
            │   │       └── [v-if pregunta.condicionalAmbitoLaboral && !esAmbitoLaboral]  (oculta)
            │   └── BannerResultadoTamizaje  (.tz-banner)
            │       "{tipoRiesgo}: nivel {nivelRiesgo}"  :class="colorPorNivel(nivelRiesgo)"
            │
            └── [moduloActual === 'GS'] SectionGS
                ├── ParrafoSintesis  (.gs-parrafo-sintesis)
                │   └── [v-if seccionA || seccionB || seccionC || seccionD]
                │       SeccionTexto × N  [v-for seccionesConDatos]
                │       (ver registro-v3-parrafo-sintesis.md — agregado tras obs. #18 del cliente)
                │
                ├── SubseccionA  "Denuncia virtual y gestión general"
                │   └── FieldRenderer × N  [v-for preguntasDeSeccion('GS', 'A')]
                │
                ├── SubseccionB  "Remisión a equipos SALVIA"  // reestructurado sobre Seguimiento real, ver registro-v3-data.md D9
                │   ├── FieldRenderer  "¿Requiere remisión?" · "¿Cuáles equipos?"
                │   ├── [v-if equipos.includes('medidas_emergencia') && equipos.includes('atencion_psico')]
                │   │   BannerAdvertencia  "No se puede remitir a estos dos equipos simultáneamente"
                │   ├── [v-if equipos.includes('medidas_emergencia')]
                │   │   <checkbox> × 6  [v-for opcionesMedidasEmergencia]
                │   ├── [v-if equipos.includes('atencion_psico')]
                │   │   ├── BannerInfo  "Criterio obligatorio + mínimo 3 puntos"
                │   │   └── <checkbox> × 7  [v-for criteriosPsicosocial]  (1 obligatorio marcado aparte + 6 con puntaje)
                │   ├── [v-if equipos.includes('atencion_hombres')]
                │   │   <checkbox> × 1  criterioHombresObligatorio
                │   ├── [v-if equipos.includes('discapacidad')]
                │   │   <checkbox> × 2  [v-for serviciosDiscapacidad]
                │   ├── [v-if equipos.includes('estabilizacion')]
                │   │   <checkbox> × 3  [v-for criteriosEstabilizacion]  (mínimo 1)
                │   └── <textarea>  "Elementos que evidencian la remisión"
                │
                ├── SubseccionC  "Registrar Entidades"  // repeater, modela Identificación de Entidades real
                │   ├── FieldRenderer  "¿Ha puesto en conocimiento...?" · "¿Salvia debe activar rutas...?" · "¿Está de acuerdo con la activación?"
                │   ├── [v-if haPuestoEnConocimiento]
                │   │   RepeaterEntidadesAcudidas
                │   │   ├── EntidadItem × N  [v-for entidadesAcudidas]
                │   │   │   └── FieldRenderer × N  (sector, depto, ciudad, municipio, entidad, radicado, fecha, acciones completadas/pendientes)
                │   │   └── BtnAgregarEntidad  "+ Siguiente entidad" → agregarEntidad('acudidas')
                │   └── [v-if salviaDebeActivar]
                │       RepeaterEntidadesActivacion
                │       ├── EntidadItem × N  [v-for entidadesActivacion]
                │       │   └── FieldRenderer × N  (igual que arriba + canal de activación + "otra información")
                │       └── BtnAgregarEntidad  "+ Siguiente entidad" → agregarEntidad('activacion')
                │
                ├── SubseccionD  "Registrar Barreras"  // repeater, modela Identificación de Barreras real
                │   ├── FieldRenderer  "¿Identificó barreras institucionales durante esta atención?"  (GAP G-V3-19, pregunta nueva)
                │   └── [v-if identificoBarreras]
                │       RepeaterBarreras
                │       ├── BarreraItem × N  [v-for barreras]
                │       │   └── FieldRenderer × N  (sector, barreras por sector + "otra", entidad asignada, categorías transversales, fecha, funcionario, descripción, gestión)
                │       └── BtnAgregarBarrera  "+ Siguiente barrera" → agregarBarrera()
                │
                ├── SubseccionE  "Alerta Inminente"
                │   └── BtnActivarAlertaInminente  "Activar Alerta Inminente" → activarAlertaInminente(origen: 'GS')
                │       (mismo flujo que el botón del header, ver flow-E13 — aquí con campo de justificación)
                │
                └── BtnFinalizar  "Finalizar registro y enrutar a Seguimiento"  → finalizarRegistro()
                    :disabled si !obligatoriasCompletas
```

### Pantalla de Cierre (tras `finalizarRegistro()`, antes de redirigir a Seguimiento)

> Agregada al revisar `Mockup_Reportes_Registro (1).html` — la primera versión
> de este plan hacía que `flow-E15` redirigiera directo a Seguimiento sin una
> pantalla de confirmación intermedia. El mockup muestra una pantalla de
> cierre con resumen antes de navegar — adición razonable, no contradice nada
> ya definido. Ver `flow-E15` actualizado.

```
CierreRegistro  (reemplaza a FormularioIngresoRegistro cuando ses.fase === 'cierre')
│
├── IconoExito  "✓" + "Formulario de Ingreso Registro terminado"
│
├── StatCards × 4
│   ├── StatCard  "Registro de víctima"  →  {idRegistro} · {estado}
│   ├── StatCard  "Caso"  →  {id} · {estado} · "en Seguimiento"
│   ├── StatCard  "Resultado del Tamizaje"  →  {tipoRiesgo} · {nivelRiesgo} · {puntaje} pts
│   └── StatCard  "Origen"  →  {canal} · {idReporte || '—'}
│
├── [v-if pendientes2doContacto > 0]
│   NotaPendientes  "Quedan N variables obligatorias de máx. 2º contacto marcadas como 'Pendiente'..."
│
├── TablaActuaciones  (Acción | Conecta con | Responsable | Fecha | Estado)  [v-for actuacionesGeneradas]
│
├── CardBitacora  "Bitácora del Caso · primera entrada"  →  {texto} · {agente} · {fecha}
│
├── CardRemisiones  "Remisiones a procesos especializados"  [v-for remisiones]
│
└── BotonesNavegacion
    ├── BtnVerExpediente  "Ver Expediente de víctima"  (placeholder, Expediente diferido — D2)
    ├── BtnIrCabezote  "Ir al Cabezote de búsqueda"  → stage = 'cabezote'
    └── BtnMisReportes  "Mis Reportes"  → fuera de alcance
```

---

## Estados visuales de los badges

| Badge | Valores | Dónde se muestra |
|---|---|---|
| `BadgeEstadoRegistro` | Incompleto · Completo | Cabezote (por Registro encontrado) y Header del Formulario (calculado, no almacenado) |
| `BadgeEstadoCaso` | Borrador · Iniciado · Terminado · **Cerrado** | Cabezote (por Caso) y Header del Formulario — "Cerrado" solo se alcanza desde el Módulo de Seguimiento, fuera de esta pantalla, pero el badge debe reconocer el valor si se reabre un Caso ya cerrado (confirmado como 4 estados por obs. #1 del cliente, el Guion Funcional de Drive solo documentaba 3) |
| `BannerResultadoTamizaje` | 4 niveles × 2 tipos = 8 combinaciones de color/texto | Dentro del módulo TZ, al final |

## Reglas de bloqueo del sidebar

| Módulo | Se bloquea cuando... |
|---|---|
| TZ, GS | `!triadaCompleta` — CI+DP+CA (Víctima) y HV+AG (Hecho+Agresor) no tienen aún todas las obligatorias [1er contacto] diligenciadas (CA no aporta: sus 4 preguntas son [máx. 2º contacto]). Umbral único, consistente con la transición Borrador→Iniciado de `flow-E12` — ver `registro-v3-data.md`. |
