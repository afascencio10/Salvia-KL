# `gestion-reporte-modal` — Interfaz del Componente

Modal reutilizable que concentra **toda** la lógica nueva de "Gestión reportes" (FigJam, canal `l25jrzyf`): ver el detalle de un reporte, decidir si procede el recontacto, registrar intentos bajo el esquema 3x3 y cerrar el reporte (o derivar a apertura de caso). Lo abren `Listado Reportes` y `Mis Reportes` por igual desde el botón de fila "Ver Reporte" (`action-clicked`, `buttonId==='ver'`).

A diferencia de `casos-component`/`reasignar-casos-modal` (que ya existen y se extienden), este componente es **enteramente nuevo**. Su diseño visual y su máquina de estados no se inventaron en esta planeación: están tomados **literal** del prototipo ya validado `Mockup_Reportes_Registro.html` (sección "PASO 6a–6e" + "PASO 7", líneas 440–665) — mismas etiquetas de botón, mismos textos de ayuda, misma secuencia de pasos. Solo se adapta de "páginas completas del mockup" a "pasos dentro de un modal", y se le agregan las llamadas HTTP reales que el mockup simula en memoria.

## Archivos relevantes

| Archivo | Rol |
|---|---|
| `src/frontend/js/components/gestion-reporte-modal.js` | Componente Vue — template y lógica (nuevo) |
| `src/frontend/css/gestion-reporte-modal.css` | Estilos (nuevo) |
| `src/salvia/controller/victim_contact_detail_controller.go` | `GET /api/v1/victim-contacts/:icode` (nuevo) |
| `src/salvia/controller/victim_contact_attempt_controller.go` | `POST /api/v1/victim-contacts/:icode/intento` (nuevo) |
| `src/salvia/controller/victim_contact_close_controller.go` | `POST /api/v1/victim-contacts/:icode/cerrar` (nuevo) |

## Relación con otros componentes

```
list_reportes.html / get_my_reportes.html (padres)
│
├── casos-component
│   └── emite action-clicked { buttonId:'ver', case }     → Listado/Mis Reportes E-04/E-02
│
└── gestion-reporte-modal
    ├── ref="gestionModal"
    ├── open(report)   ← padre llama tras el action-clicked
    └── emite 'updated' al cerrarse → el padre recarga casos-component
```

---

## Árbol de interfaz

```
gestion-reporte-modal
│
├── [v-if !visible]  → no renderiza nada
│
└── [v-if visible]
    ModalBackdrop  (.grm-backdrop)
    │  @click.self → close()   // solo si step no implica guardado en curso
    │
    └── ModalDialog  (.grm-dialog)  role="dialog" aria-modal="true"
        │
        ├── ModalHeader  (.grm-header)
        │   ├── ICode  (.grm-icode)        report.icode
        │   ├── Title  (.grm-title)        nombreMostrado || "Reporte"
        │   ├── TypeBadge  (.grm-type-badge)  [v-if reportType]  reportType  (uppercase por CSS)
        │   │   // Mientras llega el detalle, Title y TypeBadge salen de la FILA de la tabla (names + lastNames
        │   │   // y reportType); al llegar usan nombreMostrado/reportType del backend. Un reporte sin nombre
        │   │   // (tipo Hombre / registrado por un agente) deja el título en "Reporte" sobre el radicado.
        │   └── CloseBtn  (.grm-close)  "✕"  → close()
        │
        └── ModalBody  (.grm-body)
            │  [v-if loading]         Spinner  (.grm-loading)  "Cargando reporte..."                            // M-01
            │  [v-else-if loadError]  LoadError  (.grm-error)  texto del backend (404/500) o "Error al cargar el reporte"
            │                         // sin pasos ni datos vacíos; el modal sigue siendo cerrable
            │  [v-if actionError]     ActionError  (.grm-action-error)  texto de la acción fallida (intento/cierre)
            │                         // no cambia de paso, reactiva los botones
            │  [v-if !canManage && step !== 'cerrado']
            │                         ReadonlyNotice  (.grm-readonly-notice)  "Solo un agente integral o supervisor puede gestionar este reporte."
            │                         // no se muestra en el paso 'cerrado': en un reporte cerrado canManage es false por definición
            │
            │  ┌─────────────────────────────────────────────────┐
            │  │ STEP 'validacionPrevia'  (v-if step===...)       │
            │  │ Visible cuando: close_reason vacío Y 0 intentos  │
            │  └─────────────────────────────────────────────────┘
            ├── [step === 'validacionPrevia']
            │   ├── InfoBlock  (.grm-info)
            │   │   ├── p  [v-if form.nombreReportante]  "Quien reporta: {nombre} {apellido} ({parentesco})"
            │   │   ├── p  "Teléfono de contacto: {telefono || 'No suministrado'}"
            │   │   ├── p  [v-if form.correoReportante]  "Correo: {correo}"
            │   │   ├── p  "Ajuste razonable solicitado: {texto}"  + Tooltip
            │   │   └── p  "Descripción de los hechos:"  + Box  form.descripcionHechos
            │   ├── HelpText  "Verifique si en la BD SALVIA existe información asociada..."
            │   └── ButtonGrid  (2 columnas)
            │       ├── Btn "No información mínima"  + Tooltip   → cerrarPrevio('ni')   // M-02
            │       └── Btn "Procede recontacto" (primario) + Tooltip → goTo('esquema3x3') // M-04
            │       // NO incluye "Reporte de prueba" — instrucción del usuario (2026-10-12):
            │       // los reportes de prueba se generan y descartan en el ambiente de
            │       // pruebas, no necesitan un botón propio en producción.
            │
            │  ┌─────────────────────────────────────────────────┐
            │  │ STEP 'esquema3x3'                                │
            │  │ Visible cuando: close_reason vacío Y >=1 intento,│
            │  │ o tras M-04                                      │
            │  └─────────────────────────────────────────────────┘
            ├── [step === 'esquema3x3']
            │   ├── Header3x3  (.grm-3x3-header)
            │   │   ├── Title  "Esquema de Recontacto 3x3"
            │   │   └── Counter  "Hoy: {intentosHoy}/3 · Días completos: {diasCompletos}/3"
            │   │       // diasCompletos = días con >= 3 fallidos cada uno (ver GAP corregido
            │   │       // en flow-M07) — texto literal del mockup, pero la cuenta detrás
            │   │       // exige 3 por día, no 1.
            │   ├── HelpText  "Diferencie horarios (franja 6:00–18:00h)..."
            │   ├── [v-if intentosHoy >= 3 && diasCompletos < 3]
            │   │   QuotaBanner  (naranja)  "Puede continuar hoy o retomar mañana..."
            │   ├── [v-if diasCompletos >= 3]
            │   │   CompleteBanner  (verde/rojo)
            │   │   ├── Text  "Se completó el esquema 3x3 (3 intentos en 3 días distintos)"
            │   │   └── Btn "Cerrar por no contacto"  → cerrarPorNoContacto()   // M-07
            │   ├── ActionButtons
            │   │   ├── Btn "Intento Fallido"  (rojo)   → registrarIntento('fallido')   // M-05
            │   │   └── Btn "Intento Exitoso" (verde)  → registrarIntento('exitoso')   // M-06
            │   └── AttemptsLog  (.grm-attempts-log)                              ← agregado por Felipe (2026-10-12)
            │       ├── EmptyRow  [v-if intentos.length===0]  "No hay intentos registrados"
            │       └── [v-else]
            │           DayGroup × N  [v-for intentos agrupados por día de Bogotá, orden desc]  (.grm-day-group)
            │           // Se agrupa y se muestra la hora con timeZone 'America/Bogota' (no la del navegador): así los
            │           // badges X/3 coinciden con "Días completos" aunque el usuario esté fuera de Colombia.
            │           ├── DayHeader  (.grm-day-header)
            │           │   ├── DateLabel  "{fecha del grupo}"
            │           │   └── DayBadge  (.grm-day-badge)
            │           │       :class="grm-day-badge--complete" si el grupo tiene >= 3 fallidos
            │           │       "{fallidosDelDía}/3"
            │           │       // Permite ver de un vistazo qué días ya cuentan como
            │           │       // "completos" para el esquema 3x3 y cuáles van a medias.
            │           └── Row × N  [v-for intentos del día]  hora | resultado | agente
            │
            │  ┌─────────────────────────────────────────────────┐
            │  │ STEP 'postContacto'  (tras M-06, auto)           │
            │  └─────────────────────────────────────────────────┘
            ├── [step === 'postContacto']
            │   ├── Title  "Contacto establecido"
            │   ├── HelpText  "Según lo conversado, seleccione el resultado..."
            │   ├── [v-if !pendingClose]  ChoiceGrid
            │   │   ├── Btn "Cierre por no competencia"       + Tooltip → pendingClose='competencia'
            │   │   └── Btn "Cierre por decisión de la persona" + Tooltip → pendingClose='decision'
            │   │   // NO incluye "Apertura de caso" en esta fase — ver nota abajo
            │   └── [v-if pendingClose in ('competencia','decision')]
            │       ConfirmBox
            │       ├── Label  "Registre el motivo de no competencia:" | "Registre la manifestación de la persona:"
            │       ├── Textarea  v-model="motivoTexto"  maxlength="500"  (obligatorio)
            │       ├── Btn "Cancelar"           → pendingClose=''
            │       └── Btn "Confirmar cierre"  :disabled="!motivoTexto.trim()" → confirmarCierrePostContacto()  // M-08 / M-09
            │
            │   // ── Diferido, no se implementa en esta fase (instrucción del usuario,
            │   //    2026-10-11 — coincide con la nota de FigJam: "aún no vamos a
            │   //    trabajar el registro de un nuevo caso formal a partir de un
            │   //    reporte") ──
            │   // [v-if pendingClose === 'apertura']  OpenCaseBox — botón "Apertura de
            │   // caso" + pregunta de Atención Especializada + confirmarApertura().
            │   // Diseño completo conservado en flow-M10 por si se retoma más adelante.
            │   // Mientras tanto, abrir un caso desde un reporte exitoso sigue el
            │   // camino YA EXISTENTE: /salvia/primer_contacto → botón "Crear caso".
            │
            │  ┌─────────────────────────────────────────────────┐
            │  │ STEP 'cerrado'  (solo lectura)                   │
            │  │ Visible cuando: close_reason ya tiene valor      │
            │  │ (o es un legado invalidado, status 'i') al       │
            │  │ abrir (M-01), o inmediatamente tras M-02/07/08/09│
            │  └─────────────────────────────────────────────────┘
            └── [step === 'cerrado']
                ├── Icon  ✓  (verde si categoria==='Válido', gris si 'Inválido')
                ├── Title  "Reporte {icode} cerrado"
                ├── p  "Categoría: {categoria}"
                ├── p  "Motivo: {motivoCierreLabel}"   // etiqueta + " — " + observación (co/dp); en un legado, su descripción
                └── Btn "Cerrar"  → close()   // M-11, siempre emite 'updated'
```

---

## Permiso de gestión — `op` y `sv` (decisión confirmada, GAP-2)

El usuario confirmó (2026-10-11) que **tanto `op` como `sv` pueden gestionar cualquier reporte** que puedan ver — no solo el agente con `victim_contact_agent_id` propio. No existe modo solo-lectura por pertenencia: `canManage` depende únicamente del rol de la sesión, no de a quién está asignado el reporte.

```
canManage = sesión.role === 'op' || sesión.role === 'sv'
```

El detalle devuelve `canManage = (rol op/sv) AND (el reporte no está cerrado)`: es `true` para cualquiera que pueda abrir el modal sobre un reporte abierto, y `false` en uno cerrado o en un legado invalidado (no hay nada que gestionar; por eso el aviso de solo lectura no se muestra en el paso `cerrado`). El chequeo de rol se mantiene en el backend de cada endpoint de escritura activo en esta fase (M-02, M-05, M-06, M-07, M-08, M-09) como defensa en profundidad ante una llamada directa al API (un rol sin permiso recibe 401), no como una restricción visible en la UI.

---

## Diferido — "Apertura de caso" (M-10)

> Desde 2026-10-13 el backend **rechaza `close_reason='ap'` con 400** (`deferredCloseReasonCodes`, decisión del usuario): no se puede cerrar como "Cerrado - válido" por API mientras este evento no exista.

El usuario pidió explícitamente (2026-10-11) no construir este evento por ahora — coincide con la propia nota de FigJam: *"En este punto aun no vamos a trabajar el registro de un nuevo caso formal a partir de un reporte."* En esta fase, el `ChoiceGrid` de `postContacto` solo ofrece "Cierre por no competencia" y "Cierre por decisión de la persona" (M-08/M-09). El diseño de M-10 (botón, pregunta de Atención Especializada, redirección por tipo) queda documentado en [`Flujos/flow-M10-cuando-confirma-apertura-caso.md`](Flujos/flow-M10-cuando-confirma-apertura-caso.md) para cuando se retome, pero **no forma parte del alcance a implementar**. Mientras tanto, si un contacto exitoso sí amerita abrir un caso, el camino sigue siendo el ya existente: `/salvia/primer_contacto` → botón "Crear caso".

Con M-10 fuera de alcance, GAP-6 (Atención Especializada) y GAP-7 (destino para Familiar/Hombre) — ambos preguntas que solo aplicaban a esa rama — quedan igualmente diferidos, no hace falta resolverlos para aprobar o implementar el resto de este plan.
