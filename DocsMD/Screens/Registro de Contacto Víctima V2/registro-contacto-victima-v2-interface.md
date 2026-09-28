# `/salvia/public-v2/nuevo` — Interfaz de la Pantalla

> ✅ **Implementado y probado — 2026-09-27.** Ruta y nombre de archivo confirmados y en funcionamiento: `set_victim_contact_v2.html` en `/salvia/public-v2/nuevo` (GET) / `/salvia/public-v2` (POST), siguiendo la convención `_v2` ya usada en este repo (`follow_up_v2`, `set_victim_case_v2.html`). QA completo (Playwright + backend) documentado en [`QA/E10-cuando-envia-reporte-qa-cases.md`](QA/E10-cuando-envia-reporte-qa-cases.md). Ver razones de aislamiento de ruta en la sección "Por qué una ruta y no solo un archivo nuevo" más abajo.

Segunda versión pública (sin autenticación) del reporte ciudadano de hechos de VBG/VpP. Coexiste con la pantalla actual (`Registro de Contacto Víctima` / `set_victim_contact.html`) sin reemplazarla ni modificarla. Implementa el **Flujo A** del prototipo `Mockup_Reportes_Registro.html`: un wizard de 4 pasos (Autorización → Perfil de ingreso → Formulario específico por tipo de reportante → Confirmación con radicado), con 4 tipos de reportante (Víctima, Tercero, Familiar, Hombre) en vez de los 2 que distingue implícitamente la versión actual.

Archivo fuente (nuevo): `src/frontend/html/salvia/victim_contact/set_victim_contact_v2.html`

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
| `src/frontend/html/salvia/victim_contact/set_victim_contact_v2.html` *(nuevo)* | Pantalla anfitriona — wizard de 4 pasos, Vue standalone (sin `dinamic-form`; ver nota de arquitectura) |
| `src/salvia/facades/VictimContactFacadeV2.go` *(nuevo)* | Handlers públicos `VictimContactPOST_Public_V2` / `VictimContactPOST_GET_Public_V2` |
| `src/salvia/controllers/VictimContactControllerV2.go` *(nuevo)* | `SetVictimContactV2` — validación y creación propias, independientes de `SetVictimContact` (V1) |
| `src/salvia/dao/VictimContactForm2V2DAO.go` *(nuevo)* | Escribe `victim_contact_form2_v2` (tabla de extensión nueva) e invoca —sin modificarlas— las funciones DAO existentes para las columnas que sí reutiliza. Cero cambios en `VictimContactDAO.go` / `VictimContactForm2DAO.go` |
| `src/salvia/config/Locale.go` *(existente, solo adiciones)* | Nuevas keys de enum/labels exclusivas de V2 — ninguna entrada existente se modifica |
| `src/salvia/facades/MainRouter.go` *(existente, solo adiciones)* | Registra las 2 rutas nuevas en un grupo `/salvia/public-v2` separado del grupo `/salvia/public` actual |

---

## Por qué una ruta y no solo un archivo nuevo

El handler actual `VictimContactPOST_GET_Public` (V1) inyecta las opciones de `victim_case_form2_report_type` desde la BD hacia un `<select>` genérico y reutilizable (`entityAttr: "reportType", enumsAttr: "victim_case_form2_report_type"`, ver `set_victim_contact.html:36-37`). Si V2 necesitara los 2 tipos de reportante nuevos (Familiar, Hombre) agregando filas a esa misma clave de enum, **aparecerían automáticamente en el dropdown de V1** sin tocar ni una línea de su código — una regresión silenciosa. Por eso V2 usa una clave de enum propia (`victim_case_form2_report_type_v2`, decisión 3) y handlers propios: aísla el cambio de datos, no solo el de código.

La decisión 2 (backend totalmente separado) extiende este mismo criterio a todo el resto del esquema: en vez de agregar columnas nullable a `victim_contact_form2` (el planteamiento inicial de esta planeación), los campos genuinamente nuevos —anonimato, parentesco, correo y apellido del reportante, disposición/hora del reportante— viven en una tabla de extensión nueva (`victim_contact_form2_v2`). Resultado: ni una tabla ni un archivo Go existente cambia. Ver el detalle completo en [`related-tables.md`](related-tables.md).

---

## Árbol de interfaz

```
Page  (#app)
│
├── Header  (.bg-salvia)
│   ├── Titulo  "Sistema SALVIA"
│   └── Subtitulo  {{ tituloPaso }}
│
└── Main
    │
    ├── [step === 'auth']
    │   SeccionAutorizacion
    │   ├── Titulo  "Autorización de datos personales"
    │   ├── TextoLegal  (Ley 1581 de 2012 · derechos ARCO · Ley 1437/2011 art.56 mod. Ley 2080/2021 — notificaciones electrónicas)
    │   ├── SelectSingle  form.autorizacion  ["Sí, autorizo" | "No autorizo"]
    │   ├── [v-if form.autorizacion === 'No']
    │   │   └── AvisoCorte  (.bg-red-50)  "No es posible continuar sin autorización"
    │   └── BtnContinuar  :disabled="form.autorizacion !== 'Si'"  → goTo('perfil')
    │
    ├── [v-else-if step === 'perfil']
    │   SeccionPerfil
    │   ├── SelectSingle  form.tipoReporte  ["Víctima" | "Tercero" | "Familiar" | "Hombre"]
    │   ├── SelectSingle  form.ajusteRazonable  (10 opciones — Guía-intérprete, Traducción, Apoyos, Formatos accesibles, Física, Dispositivos, Confianza, Transporte, No informa, No requiere)
    │   ├── BtnAtras  → goTo('auth')
    │   └── BtnContinuar  :disabled="!form.tipoReporte || !form.ajusteRazonable"  → goTo('formulario')
    │
    ├── [v-else-if step === 'formulario']
    │   SeccionFormulario  "Datos del Reporte ({{ tipoReporteTexto }})"
    │   │
    │   ├── [v-if form.tipoReporte !== 'Victima']
    │   │   BloqueDatosReportante
    │   │   ├── [v-if form.tipoReporte === 'Tercero']
    │   │   │   └── SelectSingle  form.anonimato  ["Anónimo" | "A nombre propio"]
    │   │   │
    │   │   ├── [v-if mostrarNombreReportante]        ← computed: Tercero+Propio, o Familiar, o Hombre
    │   │   │   ├── InputText  form.nombreReportante
    │   │   │   ├── InputText  form.apellidoReportante
    │   │   │   └── InputEmail  form.correoReportante  (opcional)
    │   │   │
    │   │   ├── [v-if form.tipoReporte === 'Familiar']
    │   │   │   └── SelectSingle  form.parentesco  ["Familia" | "Familia social" | "Otra"]
    │   │   │
    │   │   ├── [v-if form.tipoReporte === 'Familiar' || form.tipoReporte === 'Hombre']
    │   │   │   └── InputTel  form.telefonoReportante
    │   │   │
    │   │   ├── [v-if form.tipoReporte === 'Tercero' || form.tipoReporte === 'Familiar']
    │   │   │   └── SelectSingle  form.rolCuidado  ["Sí" | "No"]
    │   │   │
    │   │   ├── [v-if form.tipoReporte === 'Tercero']
    │   │   │   └── SelectSingle  form.victimaConoce  ["Sí" | "No"]
    │   │   │
    │   │   └── [v-if form.tipoReporte === 'Familiar' || form.tipoReporte === 'Hombre']
    │   │       ├── SelectSingle  form.disposicionReportante  ["Sí" | "No"]
    │   │       └── [v-if form.disposicionReportante === 'Si']
    │   │           └── <input type="time">  form.horaReportante
    │   │
    │   ├── [v-if form.tipoReporte !== 'Hombre']
    │   │   BloqueDatosVictima  "Datos de la persona afectada"
    │   │   ├── InputText  form.nombreVictima
    │   │   ├── InputText  form.apellidoVictima
    │   │   │
    │   │   ├── [v-if form.tipoReporte === 'Tercero' || form.tipoReporte === 'Victima']
    │   │   │   ├── InputTel  form.telefonoVictima
    │   │   │   ├── SelectSingle  form.disposicionLlamada  ["Sí" | "No"]
    │   │   │   └── [v-if form.disposicionLlamada === 'Si']
    │   │   │       └── <input type="time">  form.horaContacto
    │   │   │
    │   │   └── [v-if form.tipoReporte === 'Familiar']
    │   │       └── NotaInformativa  "no se pide teléfono/disposición de la víctima fallecida"
    │   │
    │   ├── BloqueDescripcion
    │   │   └── <textarea>  form.descripcionHechos  (todos los tipos de reporte)
    │   │
    │   ├── BloqueCaptcha                              ← decisión 4: se mantiene
    │   │   ├── <img>  src="/captcha/{{ form.captchaID }}.png"
    │   │   └── <input>  form.captchaSolution
    │   │
    │   ├── BtnAtras  → goTo('perfil')
    │   └── BtnEnviar  "Enviar Reporte"  :disabled="!isFormValid"  → submitForm()
    │
    └── [v-else-if step === 'confirmacion']
        SeccionConfirmacion
        ├── IconoCheck  (SVG verde)
        ├── Titulo  "¡Reporte Enviado!"
        ├── DatosRadicado  (.bg-gray-50)
        │   ├── Fecha
        │   ├── Hora
        │   └── NumeroRadicado
        ├── NotaRecontacto  "Un agente intentará comunicarse al número de contacto proporcionado"
        └── BtnVolver  "← Volver al Inicio"  → resetForm()
```

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
| Captcha | Visible y requerido | Visible y requerido igual que V1 (decisión 4) — el prototipo del Flujo A no lo dibujaba, pero se mantiene por ser un endpoint público |
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
