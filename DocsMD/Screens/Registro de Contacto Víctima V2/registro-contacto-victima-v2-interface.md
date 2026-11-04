# `/salvia/public-v2/nuevo` — Interfaz de la Pantalla

> ✅ **Implementado y probado — 2026-09-27; MD sincronizado con el código el 2026-10-31.** Ruta pública `GET /salvia/public-v2/nuevo` (wizard) y `POST /salvia/public-v2` (radicación), sin sesión ni permisos. Casos de prueba en [`QA/`](QA/). Limitaciones conocidas al final de este archivo. **WP-1 (2026-11-02):** paso inicial «¿Tiene usuario?» (T8), radicado `XXXX-XXXX` con botón «Copiar radicado» (T3), captcha nuevo al volver al inicio (T2) y hora de recontacto `HH:MM:SS` aceptada (T1) — ver [`changelogNov2026.md`](changelogNov2026.md).

Segunda versión pública del reporte ciudadano de hechos de VBG/VpP. Coexiste con la pantalla anterior (`set_victim_contact.html`) sin modificarla. Es un wizard de 5 pasos (¿Tiene usuario? → Autorización → Perfil de ingreso → Formulario específico → Confirmación con radicado) con 4 tipos de reportante: Víctima (`v`), Tercero (`r`), Familiar (`fm`) y Hombre (`ho`). Lo que se radica entra a [Listado Reportes](../Listado%20Reportes/listado-reportes-index.md) como "Sin asignar" y desde ahí lo gestionan [Mis Reportes](../Mis%20Reportes/mis-reportes-index.md) y el modal [`gestion-reporte-modal`](../../Componentes/Gestion%20Reporte%20Modal/gestion-reporte-modal-index.md).

Archivo fuente: `src/frontend/html/salvia/victim_contact/set_victim_contact_v2.html` (app Vue 3 `appV2`, Options API, delimitadores `${ }`).

## Decisiones de arquitectura confirmadas (2026-09-25)

| # | Decisión | Elegida |
|---|---|---|
| 1 | Frontend | Vue wizard estático a la medida — **no** `dinamic-form` (el formulario es pequeño; el motor DAG está pensado para formularios grandes) |
| 2 | Backend | Controller y DAO **nuevos y separados**, sin compartir código con `SetVictimContact` (V1) |
| 3 | Tipos de reporte nuevos | Nueva clave de enum (`victim_case_form2_report_type_v2`) sobre la misma columna `report_type` |
| 4 | Captcha | Se mantiene, igual que V1 |
| 5 | Anonimato | Campo nuevo y dedicado, no reutiliza `report_type_details` |
| 6 | Disposición/hora de llamada | Se separan: par propio para el reportante, par existente para la víctima |
| 7 | Nombre/apellido sin víctima nombrada (Hombre) | Se dejan vacíos, a completar por el agente |

Detalle completo de cada decisión y su impacto en el esquema: [`related-tables.md`](related-tables.md).

---

## Archivos relevantes

| Archivo | Rol |
|---|---|
| `src/frontend/html/salvia/victim_contact/set_victim_contact_v2.html` | Pantalla anfitriona — wizard de 4 pasos, Vue standalone (sin `dinamic-form`) |
| `src/frontend/templates/layouts/header.html` · `footer.html` *(compartidos)* | Layout: `<body id="app" style="display:none">`, franja GOV.CO, banner de SALVIA y pie GOV.CO. La pantalla no dibuja encabezado propio |
| `src/frontend/templates/layouts/overlay.html` *(compartido)* | Popups de éxito/error globales (`#successOverlay`) |
| `src/salvia/facades/VictimContactFacadeV2.go` | Handlers públicos `VictimContactPOST_Public_V2` (POST) y `VictimContactPOST_GET_Public_V2` (GET: HTML, o catálogos en JSON si `Accept: application/json`) |
| `src/salvia/controllers/VictimContactControllerV2.go` | `SetVictimContactV2` — validación y creación propias, independientes de `SetVictimContact` (V1) |
| `src/salvia/dao/VictimContactForm2V2DAO.go` | Escribe `victim_contact_form2_v2` (tabla de extensión nueva) e invoca —sin modificarlas— las funciones DAO existentes para las columnas que sí reutiliza. Cero cambios en `VictimContactDAO.go` / `VictimContactForm2DAO.go` |
| `src/salvia/config/Locale.go` *(existente, solo adiciones)* | Textos de los enums de V2 (`victim_case_form2_report_type_v2_*`, `_anonimato_*`, `_parentesco_*`) y de los ajustes razonables |
| `src/salvia/facades/MainRouter.go` · `salvia/config/Enums.go` *(existentes)* | Registran las 2 rutas (`MainRouter.go:63-64`) y sus `FormPaths` (`Enums.go:96-97`) |

---

## Por qué una ruta y no solo un archivo nuevo

El handler actual `VictimContactPOST_GET_Public` (V1) inyecta las opciones de `victim_case_form2_report_type` desde la BD hacia un `<select>` genérico y reutilizable (`entityAttr: "reportType", enumsAttr: "victim_case_form2_report_type"`, ver `set_victim_contact.html:36-37`). Si V2 necesitara los 2 tipos de reportante nuevos (Familiar, Hombre) agregando filas a esa misma clave de enum, **aparecerían automáticamente en el dropdown de V1** sin tocar ni una línea de su código — una regresión silenciosa. Por eso V2 usa una clave de enum propia (`victim_case_form2_report_type_v2`, decisión 3) y handlers propios: aísla el cambio de datos, no solo el de código.

La decisión 2 (backend totalmente separado) extiende este mismo criterio a todo el resto del esquema: en vez de agregar columnas nullable a `victim_contact_form2` (el planteamiento inicial de esta planeación), los campos genuinamente nuevos —anonimato, parentesco, correo y apellido del reportante, disposición/hora del reportante— viven en una tabla de extensión nueva (`victim_contact_form2_v2`). Resultado: ni una tabla ni un archivo Go existente cambia. Ver el detalle completo en [`related-tables.md`](related-tables.md).

---

## Árbol de interfaz

> Los valores de los `<select>` son objetos enum `{icode, code, name}` que entrega el backend; las condiciones del template comparan por `code`. Códigos: autorización `y`/`n`; tipo de reporte `v`/`r`/`fm`/`ho`; anonimato `an`/`pr`; sí/no `y`/`n`.

```
Layout compartido  (header.html)  <body id="app">
├── Franja GOV.CO  +  Banner "Ministerio de la igualdad"
│
└── Page  (.container.mt-5 > .row.justify-content-center > .col-lg-10)
    │
    ├── [v-if errors.default && errors.default.global_msg]
    │   AlertaGlobal  (#alert .alert-danger)  ${errors.default.global_msg}
    ├── RegionAria  (#global-errors-live .sr-only, aria-live="assertive")   ← resumen de errores por campo para lectores de pantalla
    │
    ├── 🆕 [v-if step === 'usuario']                                     (T8 · paso inicial)
    │   SeccionUsuario
    │   ├── Titulo  "¿Tiene usuario en Salvia?"  (.tituloform, ícono fa-user-circle)
    │   ├── [v-if reporterName]                 // solo en, op y et con sesión iniciada
    │   │   ├── Aviso  (#reporter-banner .alert-info)  "Está reportando como ${reporterName} (${reporterRoleLabel}). El reporte quedará asociado a su usuario."
    │   │   └── BtnContinuar  "Continuar"  (#btn-usuario-continuar)  → goTo('auth')
    │   └── [v-else]
    │       ├── Texto  "Si es enlace territorial, agente integral o entidad, inicie sesión para reportar con su usuario."
    │       ├── BtnNo  "No, continuar sin usuario"  (#btn-usuario-no)  → elegirUsuario('no')  → goTo('auth')
    │       └── BtnSi  "Sí, iniciar sesión"  (#btn-usuario-si)  → elegirUsuario('si')  → /seguridad/login?next=/salvia/public-v2/nuevo
    │
    ├── [v-else-if step === 'auth']
    │   SeccionAutorizacion
    │   ├── Titulo  "Autorización de datos personales"  (.tituloform, ícono fa-list-alt)
    │   ├── Texto legal  (Ley 1581/2012 · derechos del titular · art. 56 Ley 1437/2011 mod. Ley 2080/2021 — notificaciones electrónicas · Resolución 772/2024)
    │   ├── <select>  form.autorizacion   ["Sí, autorizo" (code y) | "No autorizo" (code n)]
    │   ├── [v-if form.autorizacion.code === 'n']
    │   │   └── Aviso  (.alert-danger)  "No es posible continuar con el reporte sin la autorización de tratamiento de datos. Hasta aquí llega el formulario."
    │   └── BtnContinuar  "Continuar"  :disabled="form.autorizacion?.code !== 'y'"  → goTo('perfil')
    │
    ├── [v-else-if step === 'perfil']
    │   SeccionPerfil  "Perfil de Ingreso"
    │   ├── <select>  form.tipoReporte  × 4  [v-for enums.reportTypeV2]  (texto largo, ver tabla de etiquetas)
    │   ├── <select>  form.ajusteRazonable  × 12  [v-for enums.adjustmentsGBV]  (ver tabla de etiquetas)
    │   ├── BtnAtras  "Atrás"  → goTo('auth')
    │   └── BtnContinuar  "Continuar"  :disabled="!form.tipoReporte || !form.ajusteRazonable"  → goTo('formulario')
    │
    ├── [v-else-if step === 'formulario']
    │   SeccionFormulario  "Datos del Reporte (${tipoReporteTexto})"   ← tipoReporteTexto = nombre largo del tipo elegido
    │   │
    │   ├── [v-if tipoReporte.code !== 'v']
    │   │   BloqueQuienReporta  "Datos de quien reporta"
    │   │   ├── [v-if code === 'r']  <select>  form.anonimato  ["Anónimo" (an) | "A nombre propio" (pr)]
    │   │   ├── [v-if mostrarNombreReportante]       ← computed: fm, ho, o (r y anonimato = pr)
    │   │   │   ├── <input>  form.nombreReportante   "Nombres *"
    │   │   │   ├── <input>  form.apellidoReportante  "Apellidos *"
    │   │   │   └── <input type="email">  form.correoReportante  "Correo electrónico" (opcional)
    │   │   ├── [v-if code === 'fm']  <select>  form.parentesco  ["Familia" | "Familia social" | "Otra"]
    │   │   ├── [v-if code === 'fm' || 'ho']  <input type="tel">  form.telefonoReportante
    │   │   ├── [v-if code === 'r' || 'fm']  <select>  form.rolCuidado  [Sí | No]
    │   │   ├── [v-if code === 'r']  <select>  form.victimaConoce  [Sí | No]
    │   │   └── [v-if code === 'fm' || 'ho']
    │   │       ├── <select>  form.disposicionReportante  [Sí | No]
    │   │       └── [v-if disposicionReportante.code === 'y']  <input type="time">  form.horaReportante
    │   │
    │   ├── [v-if tipoReporte.code !== 'ho']
    │   │   BloquePersonaAfectada  "Datos de la persona afectada"
    │   │   ├── <input>  form.nombreVictima  "Nombres *"
    │   │   ├── <input>  form.apellidoVictima  "Apellidos *"
    │   │   ├── [v-if code === 'r' || 'v']
    │   │   │   ├── <input type="tel">  form.telefonoVictima
    │   │   │   ├── <select>  form.disposicionLlamada  [Sí | No]
    │   │   │   └── [v-if disposicionLlamada.code === 'y']  <input type="time">  form.horaContacto
    │   │   └── [v-if code === 'fm']
    │   │       └── Aviso  (.alert-warning)  "Para reportes de Familiar no se solicita teléfono ni disposición de la víctima (fallecida) — el recontacto se hace al teléfono de quien reporta, capturado arriba."
    │   │
    │   ├── BloqueDescripcion  "Descripción"
    │   │   └── <textarea>  form.descripcionHechos  "Descripción detallada de los hechos *"  (todos los tipos)
    │   │
    │   ├── BloqueCaptcha
    │   │   ├── <img>  #captchaIMG  src=/seguridad/login/captcha/${form.captchaID}/img
    │   │   ├── <audio controls>  #captchaAudio  src=/seguridad/login/captcha/${form.captchaID}/audio
    │   │   └── <input>  form.captchaSolution  "Escriba el texto de la imagen *"
    │   │
    │   ├── BtnAtras  "Atrás"  → goTo('perfil')
    │   └── BtnEnviar  "Enviar Reporte"  :disabled="!isFormValid"  → submitForm()
    │
    └── [v-else-if step === 'confirmacion']
        SeccionConfirmacion  (.text-center)
        ├── Icono  (fa-check-circle, verde)
        ├── Titulo  "¡Reporte Enviado!"
        ├── Texto  "Salvia le agradece por aportar a la identificación de hechos de violencia. Su reporte ha quedado ingresado en el sistema."
        ├── DatosRadicado  (.border.rounded.p-4)
        │   ├── "Fecha:"  ${fechaActual}  (toLocaleDateString('es-CO'))
        │   ├── "Hora:"  ${horaActual}  (toLocaleTimeString('es-CO'))
        │   └── "Número de Radicado:"  ${radicado}  (radicado «XXXX-XXXX» devuelto por el backend; si faltara, el icode)   ✏️ T3
        │       └── 🆕 BtnCopiar  "Copiar radicado"  (#btn-copiar-radicado)  → copiarRadicado()   // «Copiado» 3 s; si el navegador lo impide: aviso para copiarlo a mano
        ├── Texto  "Un agente integral intentará comunicarse al número de contacto proporcionado para la validación y posible apertura del caso."
        └── BtnVolver  "← Volver al Inicio"  (.btn-link)  → resetForm()   // pide un captchaID nuevo (T2) y vuelve al paso inicial
```

El paso inicial es `usuario`, salvo con una sesión de otro rol (sv, ad, ro…): ahí arranca en `auth`, sin el paso ni el mensaje, y el reporte igual queda asociado a su usuario (decisión de la revisión de WP-1).

### Etiquetas reales de los catálogos

| Catálogo | Código → texto que ve la persona |
|---|---|
| Tipo de reporte (`reportTypeV2`) | `v` "Usted es la víctima de violencia basada en género y/o por prejuicio" · `r` "Usted quiere reportar un hecho de violencia basada en género y/o por prejuicio experimentado por otra persona" · `fm` "Usted es un familiar de víctima de feminicidio o de homicidio abyecto o fútil" · `ho` "Usted es un hombre con posibilidad de ejercer violencia de género" |
| Ajuste razonable (`adjustmentsGBV`, 12) | `sc` Intérprete de lengua de señas colombiana · `gi` Guía–intérprete para personas sordociegas · `id` Intérprete de idiomas / traducción · `ap` Apoyos para la comunicación · `in` Información en formatos accesibles · `ad` Adecuaciones de accesibilidad física · `dt` Dispositivos o ayudas técnicas · `pr` Presencia de una persona de confianza · `tr` Transporte seguro o acompañamiento · `ni` No informa · `nr` No requiere · `ns` No requiere apoyos especiales |
| Anonimato | `an` Anónimo · `pr` A nombre propio |
| Parentesco | `fa` Familia · `fs` Familia social · `ot` Otra |

### Campos que exige cada tipo (`isFormValid`, habilita "Enviar Reporte")

Siempre: tipo de reporte, ajuste razonable, descripción de ≥ 10 caracteres (sin espacios sobrantes) y texto del captcha.

| Tipo | Además exige |
|---|---|
| Víctima (`v`) | nombres y apellidos de la víctima, teléfono, disposición de llamada (y hora si es "Sí") |
| Tercero (`r`) | lo de Víctima + forma del reporte; si es "A nombre propio": nombres y apellidos del reportante; rol de cuidado; si la víctima conoce el reporte |
| Familiar (`fm`) | nombres, apellidos y teléfono del reportante; disposición del reportante (y hora si es "Sí"); nombres y apellidos de la víctima; parentesco; rol de cuidado |
| Hombre (`ho`) | nombres, apellidos y teléfono del reportante; disposición del reportante (y hora si es "Sí") |

El teléfono no se valida por formato en el cliente: se envía como `parseInt(valor.replace(/\D/g, ''), 10)`. El botón solo se habilita; **no se pintan errores por campo en el cliente**. Los errores por campo que devuelve el backend (400) se reflejan solo en el alert global y en la región `aria-live`.

---

## Diferencias clave vs. la pantalla pública actual (`set_victim_contact.html`)

| Aspecto | V1 (`set_victim_contact.html`, actual) | V2 (`set_victim_contact_v2.html`, nueva) |
|---|---|---|
| Estructura | Formulario único, todos los campos visibles a la vez | Wizard de 4 pasos secuenciales, con navegación Atrás/Continuar |
| Tipos de reportante | 2 implícitos (Víctima directa / Tercero) vía `victim_case_form2_report_type` | 4 explícitos: Víctima, Tercero, Familiar, Hombre — con bloques condicionales propios cada uno |
| Anonimato del reportante | No contemplado | Nuevo — solo para Tercero (`form.anonimato`) |
| Parentesco | No contemplado | Nuevo — solo para Familiar (`form.parentesco`) |
| Correo del reportante | No contemplado | Nuevo — opcional, para Familiar/Tercero/Hombre con identidad |
| Ajuste razonable | Multi-select (`VictimContactForm2AdjustmentsGBV []enum`) | Single-select en UI; se envía como arreglo de 1 elemento — mismo campo, sin cambio de esquema |
| Captcha | Visible y requerido | Visible y requerido igual que V1 (decisión 4), con imagen y audio; el cliente móvil lo omite con el `User-Agent` `flutter-client` |
| Radicado / confirmación | Respuesta JSON simple | Pantalla de confirmación dedicada con fecha, hora y número de radicado, y aviso del esquema de recontacto |

---

## Esquema de datos

Ver detalle completo de mapeo a base de datos, la migración ejecutada y los ajustes que tomó correrla contra la BD real en [`related-tables.md`](related-tables.md).

---

## Bugs encontrados y corregidos durante implementación/QA (2026-09-27)

Ninguno de estos estaba en la planeación original — surgieron al implementar y probar. Detalle completo de cada uno, con causa raíz y fix, en [`related-tables.md`](related-tables.md#bugs-encontrados-durante-implementación-y-qa-2026-09-27) y en [`QA/E10-cuando-envia-reporte-qa-cases.md`](QA/E10-cuando-envia-reporte-qa-cases.md).

| # | Bug | Capa |
|---|---|---|
| 1 | `id="app"` duplicado — el `<body>` del layout compartido ya usa ese id; el `<div>` propio de esta pantalla quedaba invisible para siempre | Frontend |
| 2 | Embeber `VictimContactDTO` (que tiene `UnmarshalJSON` propio) dentro del struct de request rompía la decodificación del campo hermano `form2V2` | Backend |
| 3 | Locale equivocado (`salvia_config.Locale` en vez de `security_config.Locale`) para el mensaje de error de captcha — un captcha incorrecto no bloqueaba el envío | Backend |
| 4 | Mismo bug sistémico en ~8 validaciones de campo requerido (`salvia_config.Locale` en vez de `common_config.Locale`) — nombres vacíos, teléfono en 0, descripción muy corta, "No autorizo" se aceptaban en silencio | Backend |
| 5 | Faltaban permisos `GRANT` para `salvia_legacy` (el rol con el que la app corre queries) sobre la tabla nueva | Migración / BD |
| 6 | El popup de éxito compartido (`overlay.html`) no se auto-oculta — bloqueaba el botón "Volver al Inicio" en la confirmación | Frontend |
| 7 | Validación de longitud de `factsDescription` comparaba contra un `MinSize` que siempre es 0 — nunca se disparaba | Backend |

---

## Limitaciones conocidas (verificadas leyendo el código, 2026-10-31)

| # | Qué pasa | Dónde |
|---|---|---|
| L1 | ✅ **Resuelta (WP-1 T1, 2026-11-02).** `VictimContactForm2V2DTO.UnmarshalJSON` acepta `HH:MM:SS`, `HH:MM` y RFC 3339; una hora ilegible responde 400 «Hora con formato inválido» en `form2V2.reporterBestContactTime` en vez de descartarse | `VictimContactForm2V2DAO.go`, `VictimContactControllerV2.go` |
| L2 | ✅ **Resuelta (WP-1 T2, 2026-11-02).** `resetForm()` llama a `refreshCaptcha()`, que pide un captchaID nuevo a la propia página con `Accept: application/json`; si la petición falla se conserva el comportamiento anterior y se deja un aviso en consola | `set_victim_contact_v2.html` |
| L3 | El botón "Enviar Reporte" se deshabilita sin decir qué falta: no hay errores por campo del lado del cliente | `isFormValid` (`set_victim_contact_v2.html:280-304`) |
| L4 | Cualquier 400 del backend (no solo el de captcha) devuelve un `captchaID` nuevo y el cliente lo asigna, porque el controller lo regenera siempre que hay errores | `VictimContactControllerV2.go:199-201` |

