# `barrier-follow-up-timeline` — Guía de uso

Componente de solo lectura que muestra el historial cronológico de seguimientos (`barrier_follow_up`) de una barrera específica.

---

## Props

| Prop | Tipo | Requerido | Default | Descripción |
|---|---|---|---|---|
| `barrierId` | `String` | Sí | — | UUID de la barrera (`barrier_v2.id`) cuyos seguimientos se quieren mostrar |

---

## Eventos emitidos

El componente no emite eventos. Es exclusivamente de lectura.

---

## Qué hace el componente por sí solo

Al montarse y cada vez que cambia `barrierId`, llama al API para cargar los seguimientos de la barrera. Maneja internamente los estados de carga, error y vacío. No requiere intervención del padre después del montaje.

---

## Integración en el HTML de la pantalla

```html
{{ template "components/barrier_follow_up_timeline.html" . }}
```

En el template del padre:

```html
<barrier-follow-up-timeline :barrier-id="barreraSeleccionada.id">
</barrier-follow-up-timeline>
```

---

## Ejemplo de uso

```javascript
// En el componente padre (ej: barrera_detalle.html):
data() {
  return {
    barrierId: '90e81aed-91e5-4ef0-bb20-cc4e03662e61'
  }
}

// Template:
// <barrier-follow-up-timeline :barrier-id="barrierId"></barrier-follow-up-timeline>
```

---

## Endpoint que consume

| Acción | Método | URL |
|---|---|---|
| Cargar seguimientos de la barrera | GET | `/api/v1/barriers-v2/:barrierId/follow-ups` |

### Respuesta esperada

```json
[
  {
    "id": "uuid",
    "barrierId": "uuid",
    "followUpId": "uuid",
    "createdById": "icode del usuario",
    "actorName": "Nombre del usuario",
    "persists": true,
    "institutionalResponse": "respuesta_oficial",
    "managementActions": "articulacion_institucional,gestion_llamada",
    "actions": "Texto libre de actuaciones realizadas",
    "closesBarrier": false,
    "closureReason": "",
    "createdAt": "2026-07-02T07:23:18Z"
  }
]
```

> **Nota:** El campo `actorName` debe ser resuelto por el backend a partir de `created_by_id`. El frontend no consulta usuarios por separado.
