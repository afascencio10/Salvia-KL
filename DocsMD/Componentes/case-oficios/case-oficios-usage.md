# `case-oficios` — Guía de uso

Lista los oficios (`entity_letter`) de un caso o de una barrera puntual, en formato de cards, con buscador, chips de tema y filtro por estado. Al hacer click en una card abre un modal con el detalle completo del oficio.

---

## Props

| Prop | Tipo | Requerido | Default | Descripción |
|---|---|---|---|---|
| `caseId` | String | Solo si no se pasa `barrierId` | `null` | ID del caso (`victim_case`) del cual se listan los oficios |
| `barrierId` | String | Opcional | `null` | ID de la barrera (`barrier_v2`). Si se pasa, tiene prioridad sobre `caseId`: filtra los oficios de esa barrera únicamente y oculta los chips de tema |

---

## Eventos emitidos

Ninguno. Es un componente de presentación pura — no emite eventos al padre.

---

## Qué hace el componente por sí solo

Al montarse, hace `GET /api/v1/entity-letters?barrierId={barrierId}` (si se pasó `barrierId`) o `?caseId={caseId}` y muestra los oficios como cards. Los filtros (buscador, chips de tema, estado) son enteramente internos — no requieren ninguna acción del padre más allá de pasarle el `caseId` y/o `barrierId`.

---

## Integración

**Detalle del Caso** (modo `caseId` — todos los oficios del caso, tab "Gestión institucional"):

```html
<case-oficios :case-id="caseICode"></case-oficios>
```

**Detalle de Barrera** (modo `barrierId` — solo los oficios de esa barrera, tab "Oficios"):

```html
<case-oficios :case-id="barrera.caseId" :barrier-id="barrierICode"></case-oficios>
```

---

## Relación con `oficios-list.js`

`case-oficios` es un componente **nuevo y separado** — no reemplaza ni modifica `oficios-list.js` (que sigue usándose tal cual en la pantalla de Notificaciones, en formato tabla). Ambos consumen el mismo endpoint (`GET /api/v1/entity-letters`) pero sirven layouts y necesidades de filtrado distintas. El modal de detalle de `case-oficios` reutiliza el mismo *diseño de contenido* que el modal ya probado en `oficios-list.js` (mismas secciones: Oficio, Radicación y Respuesta, Corrección, Barrera relacionada) — pero es una implementación propia dentro de este componente, no una dependencia compartida.

---

## Endpoint que consume

| Acción | Método | URL |
|---|---|---|
| Cargar oficios del caso | GET | `/api/v1/entity-letters?caseId={caseId}` |
| Cargar oficios de la barrera | GET | `/api/v1/entity-letters?barrierId={barrierId}` |

> Ambos endpoints ya existían (usados por `oficios-list.js` y por rutas de colección internas respectivamente) y devuelven el mismo shape `[]EntityLetter`, mismo orden (`created_at ASC`) — sin relaciones de caso/barrera enriquecidas (esas solo vienen cuando se filtra por `agentId`/`notificationUserId`). Ninguno soporta combinarse con `state` en la misma llamada — el filtro por estado se resuelve en el cliente.

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
