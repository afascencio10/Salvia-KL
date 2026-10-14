# `case-gestion-institucional` — Guía de uso

Mini-dashboard de solo lectura con el resumen de gestión institucional de un caso (entidades, obligaciones, barreras, oficios). Cada indicador abre un modal con el detalle tabular filtrado. Puerto del componente `ResumenGestionInstitucional.jsx` del mockup React a datos reales.

---

## Props

| Prop | Tipo | Requerido | Default | Descripción |
|---|---|---|---|---|
| `caseId` | String | Sí | — | `victim_case_i_code` del caso a resumir |
| `barriers` | Array | Sí | — | Barreras del caso, **reutilizadas del padre** (mismo array que `get_case_detail_sv.html` ya carga para la lista de barreras del tab) — evita un fetch duplicado |

---

## Eventos emitidos

Ninguno. Toda la interacción (abrir/cerrar modal, filtrar por sector o estado) es interna al componente.

---

## Qué hace el componente por sí solo

Al montarse, dispara en paralelo tres llamadas: entidades del caso (`GET /api/v1/casos/:caseId/entidades`), resumen de obligaciones (`GET /api/v1/casos/:caseId/obligaciones-resumen` — agregado global + por sector + detalle, calculado en una sola consulta) y oficios (`GET /api/v1/entity-letters?caseId=`). Las barreras no las fetchea — usa directamente la prop `barriers`. Calcula y muestra 4 KPIs, el desglose por sector con anillos de progreso, y dos listas de badges (barreras y oficios por estado). Cada número es clickeable y abre un modal con el detalle fila por fila, filtrado según qué se clickeó. Maneja internamente sus propios estados de carga, error y vacío — no requiere que el padre intervenga.

---

## Integración

```html
{{ template "components/case_gestion_institucional.html" . }}
```

En el template del padre (`get_case_detail_sv.html`, tab "Gestión institucional" — **antes** de `<case-entities>`):

```html
<case-gestion-institucional
  :case-id="caseICode"
  :barriers="barriers">
</case-gestion-institucional>
```

---

## Endpoints que consume

| Acción | Método | URL | Estado |
|---|---|---|---|
| Listar entidades del caso | GET | `/api/v1/casos/:caseId/entidades` | Ya existe — usado también por `case-entities` |
| Resumen agregado de obligaciones | GET | `/api/v1/casos/:caseId/obligaciones-resumen` | **Nuevo** — a crear junto con este componente |
| Listar oficios del caso | GET | `/api/v1/entity-letters?caseId=` | Ya existe — usado también por `case-oficios` |

### Shape esperado — `GET /api/v1/casos/:caseId/obligaciones-resumen` (nuevo)

```json
{
  "global": { "completadas": 8, "total": 17 },
  "porSector": [
    { "sector": "he", "cantidadEntidades": 1, "completadas": 4, "total": 4 },
    { "sector": "js", "cantidadEntidades": 1, "completadas": 2, "total": 5 },
    { "sector": "pt", "cantidadEntidades": 1, "completadas": 2, "total": 8 }
  ],
  "items": [
    { "entidad": "FISCALÍA LOCAL DE ALEJANDRÍA", "sector": "js", "obligacion": "Asignar fiscal responsable del caso", "completada": false },
    { "entidad": "FISCALÍA LOCAL DE ALEJANDRÍA", "sector": "js", "obligacion": "Recibir la denuncia formalmente", "completada": true }
  ]
}
```

> `sector` viene como código de 2 caracteres (`entity.entity_sector`, igual que ya resuelve `case-entities`) — el frontend mapea a la etiqueta (`he`→Salud, `js`→Justicia, `pt`→Protección) con la misma tabla `CE_SECTOR_LABELS` que ya usa `case-entities`. `porSector` solo incluye sectores con al menos una entidad — el frontend completa con "Sin obligaciones" los que falten de la lista fija (Salud/Justicia/Protección). `items` trae una fila por cada obligación de cada entidad del caso (join `entity_case` → `entity_branch` → `entity` → `entity_case_obligation` → `entity_obligation`), calculado en una sola consulta.
