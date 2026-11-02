# `/salvia/public-consulta` — Interfaz de la Pantalla

> ✅ **Implementado y probado — WP-1 T4 (2026-11-02).** Pública, sin sesión. Casos de prueba en [`QA/`](QA/).

Pantalla donde una persona escribe el número de radicado que recibió al enviar su reporte ciudadano V2 y ve, en lenguaje simple, en qué estado va: **recibido**, **en gestión** o **cerrado**. No muestra datos personales ni detalles de la gestión (agente, intentos, motivo de cierre). El campo aplica una máscara `XXXX-XXXX` y el botón se habilita solo con un formato válido. Un radicado mal escrito o inexistente da siempre el mismo mensaje y, tras 5 fallos desde una misma IP en 15 minutos, la consulta se bloquea temporalmente. Es una app Vue 3 standalone, igual que el wizard V2, montada sobre el layout compartido GOV.CO.

## Archivos relevantes

| Archivo | Rol |
|---|---|
| `src/frontend/html/salvia/victim_contact/consulta_reporte.html` | Plantilla y app Vue 3 `appConsulta` (Options API, delimitadores `${ }`) |
| `src/salvia/facades/ReportStatusFacade.go` | `ReportStatusGET_Public`: renderiza la página |
| `src/salvia/controller/report_status_controller.go` | `POST /api/v1/public/consulta-reporte` |
| `src/salvia/service/report_status_service.go` | Normalización, traducción de estados y limitador por IP |
| `src/internal/repository/report_status_repository.go` | `FindByRadicado`: SELECT parametrizado |
| `src/salvia/facades/MainRouter.go`, `src/salvia/config/Enums.go`, `src/salvia/config/Locale.go`, `src/main.go` | Ruta, plantilla, título y cableado (solo adiciones) |

## Árbol de interfaz

```
/salvia/public-consulta  (appConsulta · Vue 3)
Layout compartido  (header.html)  <body id="app">
└── Page  (.container.mt-5 > .row.justify-content-center > .col-lg-8)
    ├── Titulo  "Consulte el estado de su reporte"  (.tituloform, ícono fa-search)
    ├── Texto  "Escriba el número de radicado que recibió al enviar su reporte."
    │
    ├── [v-if vista === 'consulta']
    │   FormularioConsulta  <form @submit.prevent="consultar">
    │   ├── <label for="radicado">  "Número de radicado *"  (.fw-bold)
    │   ├── <input id="radicado">  v-model="radicado"  maxlength="9"  autocomplete="off"  placeholder="XXXX-XXXX"  @input="onRadicadoInput"
    │   ├── Ayuda  (.form-text)  "Ejemplo: K7QM-4R2X. Los radicados no usan las letras I ni O ni los números 0 y 1."
    │   ├── [v-if mensajeError]  Alerta  (#consulta-error .alert-danger, role="alert", aria-live="assertive")
    │   └── BtnConsultar  "Consultar"  (#btn-consultar)  type="submit"  :disabled="!radicadoValido || consultando || bloqueado"
    │       └── [v-if consultando]  "Consultando..."
    │
    ├── [v-else-if vista === 'resultado']
    │   TarjetaResultado  (#consulta-resultado .border.rounded.p-4)
    │   ├── "Número de radicado:"  ${resultado.radicado}
    │   ├── "Fecha de radicación:"  ${fechaFormateada}      // toLocaleDateString('es-CO'), fecha de Colombia
    │   ├── "Estado:"  Badge  (.badge .ms-2)  ${resultado.estadoTexto}
    │   │             recibido → text-bg-info · en_gestion → text-bg-warning · cerrado → text-bg-success   (Bootstrap 5.3)
    │   ├── Descripcion  (.text-muted)  ${resultado.descripcion}
    │   └── BtnOtra  "Hacer otra consulta"  (#btn-otra-consulta .btn-link)  → nuevaConsulta()
    │
    └── Pie  (.text-muted .mt-4)  "¿Aún no ha hecho su reporte? Hágalo aquí"  → /salvia/public-v2/nuevo
```

## Estados visuales del resultado

| Estado | Badge | Texto del estado | Descripción mostrada |
|---|---|---|---|
| `recibido` | azul (`text-bg-info`) | Recibido | Su reporte fue recibido. Un agente integral lo revisará y se comunicará con usted. |
| `en_gestion` | amarillo (`text-bg-warning`) | En gestión | Su reporte está siendo gestionado por el equipo de Salvia. |
| `cerrado` | verde (`text-bg-success`) | Cerrado | La gestión de su reporte finalizó. Gracias por comunicarse con Salvia. |

## Mensajes del sistema

| Situación | Mensaje |
|---|---|
| Radicado mal escrito o inexistente (404) | No encontramos un reporte con ese radicado. Verifique que esté bien escrito. |
| Demasiados intentos (429) | Demasiados intentos. Intente de nuevo en N minutos. |
| Error de red o del servidor | No fue posible consultar en este momento. Intente de nuevo. |
