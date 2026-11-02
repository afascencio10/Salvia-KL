# `gestion-reporte-modal` — Guía de uso

## Props

| Prop | Tipo | Requerido | Default | Descripción |
|---|---|---|---|---|
| — | — | — | — | Sin props. Igual que `reasignar-casos-modal`, se controla por método público (`open`), no por props reactivas. |

## Método público

| Método | Parámetros | Descripción |
|---|---|---|
| `open(report)` | `ReportListItem` | Abre el modal para el reporte recibido. Dispara M-01 (carga el detalle completo vía backend — la fila de la tabla no trae todos los campos). Mientras llega el detalle, el header (título y tipo) se pinta con los datos de la fila. |

El padre lo invoca así:

```javascript
onActionClicked: function({ buttonId, case: report }) {
    if (buttonId === 'ver') {
        this.$refs.gestionModal.open(report);
    }
}
```

## Requisitos de la página anfitriona (WP-1)

El modal usa `<case-timeline>` para el «Historial del reporte»: la página que lo incluye debe cargar `/static/js/components/case-timeline.js` **antes** de `gestion-reporte-modal.js` (y definir `var app = home;` como ya hace). `list_reportes.html` y `get_my_reportes.html` lo hacen.

## Eventos emitidos

| Evento | Cuándo se dispara | Payload |
|---|---|---|
| `updated` | Al cerrarse el modal (M-11), **siempre** — haya habido o no un cambio real. El padre recarga `casos-component` sin necesidad de rastrear si algo cambió. | `{}` |

## Qué hace el componente por sí solo

- Decide su propia vista inicial al abrir (`validacionPrevia` / `esquema3x3` / `postContacto` / `cerrado`) a partir del detalle que trae del backend — el padre nunca le indica en qué paso empezar.
- Muestra `intentosHoy` y `diasCompletos` tal como los calcula el backend (vienen en el detalle y en la respuesta de cada intento; el cliente no los recalcula). Lo que sí hace en el cliente es agrupar el log de intentos por día de Bogotá y formatear las horas con `timeZone: 'America/Bogota'`.
- Muestra el encabezado (nombre, tipo, teléfono y radicado), la sección «Hechos del reporte» y el «Historial del reporte» (recarga con `historyVersion` tras cada intento o cierre).
- Tras un contacto exitoso de un reporte Víctima o Tercero ofrece «Registro de caso» (crea el borrador de Registro V3 y navega a `/salvia/registro-v3?caso=`); si el reporte ya tiene borrador abre en el paso «Borrador iniciado» con «Ir al borrador del caso».
- Hace sus propias llamadas HTTP de escritura (intento, cierre, conversión) — el padre no interviene en el guardado, solo recibe `updated` al final.
- Valida permisos de gestión (`op`/`sv`, GAP-2) en el backend de cada endpoint de escritura — no depende de ningún prop que el padre deba calcular.

## Ejemplos de uso

### Caso 1 — Pantalla de gestión (Listado Reportes / Mis Reportes)

```html
<gestion-reporte-modal ref="gestionModal" @updated="onReporteActualizado"></gestion-reporte-modal>
```

```javascript
onReporteActualizado: function() {
    this.$refs.casosComponent.reload();
}
```

No hay un "Caso 2" — a diferencia de `reasignar-casos-modal`, este componente no tiene variantes de configuración; su comportamiento cambia solo según los datos del reporte que recibe.
