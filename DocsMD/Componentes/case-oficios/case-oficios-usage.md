# `case-oficios` — Guía de uso

> ⚠️ Componente **no implementado aún** — MD generado antes del código para revisión previa.

Lista los oficios (`entity_letter`) de un caso en formato de cards, con buscador, chips de tema y filtro por estado. Al hacer click en una card abre un modal con el detalle completo del oficio.

---

## Props

| Prop | Tipo | Requerido | Default | Descripción |
|---|---|---|---|---|
| `caseId` | String | Sí | — | ID del caso (`victim_case`) del cual se listan los oficios |

---

## Eventos emitidos

Ninguno. Es un componente de presentación pura — no emite eventos al padre.

---

## Qué hace el componente por sí solo

Al montarse, hace `GET /api/v1/entity-letters?caseId={caseId}` y muestra los oficios como cards. Los filtros (buscador, chips de tema, estado) son enteramente internos — no requieren ninguna acción del padre más allá de pasarle el `caseId`.

---

## Integración

```html
<case-oficios :case-id="caseICode"></case-oficios>
```

Pensado para montarse dentro de la pantalla `Detalle del Caso`, probablemente en la pestaña de Derivaciones o Barreras (junto a `case-tasks`) — la ubicación exacta queda por confirmar.

---

## Relación con `oficios-list.js`

`case-oficios` es un componente **nuevo y separado** — no reemplaza ni modifica `oficios-list.js` (que sigue usándose tal cual en la pantalla de Notificaciones, en formato tabla). Ambos consumen el mismo endpoint (`GET /api/v1/entity-letters`) pero sirven layouts y necesidades de filtrado distintas. El modal de detalle de `case-oficios` reutiliza el mismo *diseño de contenido* que el modal ya probado en `oficios-list.js` (mismas secciones: Oficio, Radicación y Respuesta, Corrección, Barrera relacionada, Caso relacionado) — pero es una implementación propia dentro de este componente, no una dependencia compartida.

---

## Endpoint que consume

| Acción | Método | URL |
|---|---|---|
| Cargar oficios del caso | GET | `/api/v1/entity-letters?caseId={caseId}` |

> Este endpoint ya existe y está en uso por `oficios-list.js`. Devuelve `[]EntityLetter` (sin relaciones de caso/barrera enriquecidas — esas solo vienen cuando se filtra por `agentId`/`notificationUserId`). No soporta combinar `caseId` con `state` en la misma llamada — el filtro por estado se resuelve en el cliente.

### Shape esperado por oficio (respuesta del endpoint)

```json
{
  "id": "uuid",
  "barrierId": "uuid",
  "caseId": "uuid",
  "state": "por_proyectar | para_revisar | en_correccion | aprobacion_juridica | para_radicar | radicado | respondido",
  "priority": "normal | alta",
  "agentId": "icode | null",
  "notificationUserId": "icode | null",
  "entidad": "string | null",
  "nivel": "nacional | departamental | municipal | null",
  "urlKofax": "string | null",
  "subject": "string | null",
  "asuntoRadicado": "string | null",
  "numeroRadicado": "string | null",
  "correoEntidad": "string | null",
  "correoRemitente": "string | null",
  "asuntoRespuesta": "string | null",
  "responseDate": "timestamp | null",
  "reasonCorrection": "string | null",
  "createdAt": "timestamp",
  "updatedAt": "timestamp"
}
```

> ⚠️ **GAP de schema** (ver detalle en `case-oficios-interface.md`): para que el chip "Medidas de Emergencia" / "Apoyo Psicosocial" / "Estabilización Económica" funcione, este shape necesitaría además `emergencyMeasureId`, `psychosocialSupportId`, `economicStabilizationId` — columnas que hoy no existen en `entity_letter`.

---

## Estados posibles (`state`)

| Valor | Label mostrado |
|---|---|
| `por_proyectar` | Por proyectar |
| `para_revisar` | Para revisar |
| `en_correccion` | En corrección |
| `aprobacion_juridica` | Aprobación jurídica |
| `para_radicar` | Para radicar |
| `radicado` | Radicado |
| `respondido` | Respondido |
