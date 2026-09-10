# `case-mecanismos` — Guía de uso

Lista y registra los escalamientos del caso a usuarios con rol Mecanismo (`mc`). El select del modal carga usuarios activos cuyo `articulating_mechanism` coincide con el tipo y la ubicación elegidos.

---

## Props

| Prop | Tipo | Requerido | Descripción |
|---|---|---|---|
| `caseId` | String | Sí | `victim_case_i_code` |
| `userId` | String | Sí | icode del usuario en sesión |
| `userRole` | String | Sí | Controla el botón agregar (`sv` / `op` / `ro`). El backend valida aparte. |

## Eventos emitidos

Ninguno.

## Integración

```html
<case-mecanismos :case-id="caseICode" :user-id="userICode" :user-role="userRole"></case-mecanismos>
```

## API

| Método | Ruta |
|---|---|
| GET | `/api/v1/casos/:id/mecanismos` |
| POST | `/api/v1/casos/:id/mecanismos` |
| POST | `/api/v1/casos/:id/mecanismos/:mid/comentarios` (endpoint conservado; el input de la card está retirado) |
| GET | `/api/v1/mecanismos/usuarios?scope=&department_id=&city_id=&town_id=` |

`scope`: `municipal` | `departamental` | `nacional`.

Filtro de usuarios `mc` (tabla `salvia.articulating_mechanism` vía `general_user.articulating_mechanism_id`):
- Municipal: `type = municipal` y `town_code` = municipio.
- Departamental: `type = departamental` y `department_id` = departamento.
- Nacional: `type = nacional`.

Usuarios `mc` sin fila de mecanismo no aparecen hasta re-guardarlos en el alta/edición de usuario.

Al crear un registro se escribe un evento en el Timeline (categoría `Mecanismos`). Cada gestión (`POST /api/v1/case-tasks/gestion-mecanismo`) también escribe un evento (`Gestión del mecanismo`). La card lista esas gestiones en `gestiones` (no los comentarios viejos).

## Shape de un registro

```json
{
  "id": "...",
  "caseId": "...",
  "scope": "municipal",
  "assignedUserId": "...",
  "assignedUserName": "Ana Pérez",
  "departmentName": "San Salvador",
  "cityName": "San Salvador",
  "townName": "San Salvador Centro",
  "salviaCommitment": "...",
  "gestiones": [
    {
      "id": "...",
      "result": "Abordado en sesión: Sí, sesión ordinaria. Reporte de gestión: …",
      "formData": { "abordadoEnSesion": "si_ordinaria", "reporteGestion": "…" },
      "completedAt": "...",
      "assignedUserName": "Ana Pérez"
    }
  ]
}
```


---

## Props

| Prop | Tipo | Requerido | Descripción |
|---|---|---|---|
| `caseId` | String | Sí | `victim_case_i_code` |
| `userId` | String | Sí | icode del usuario en sesión |
| `userRole` | String | Sí | Controla el botón agregar (`sv` / `op` / `ro`). El backend valida aparte. |

## Eventos emitidos

Ninguno.

## Integración

```html
<case-mecanismos :case-id="caseICode" :user-id="userICode" :user-role="userRole"></case-mecanismos>
```

## API

| Método | Ruta |
|---|---|
| GET | `/api/v1/casos/:id/mecanismos` |
| POST | `/api/v1/casos/:id/mecanismos` |
| POST | `/api/v1/casos/:id/mecanismos/:mid/comentarios` |
| GET | `/api/v1/mecanismos/usuarios?scope=&department_id=&city_id=&town_id=` |

`scope`: `municipal` | `departamental` | `nacional`.

Filtro de usuarios `mc`:
- Municipal: `profile_town` = municipio (`town_id` es el `town_code` DANE).
- Departamental: `city.department_id` = departamento; `city_id` opcional.
- Nacional: todos los `mc` activos.

Al crear un registro se escribe un evento en el Timeline (categoría `Mecanismos`). Cada gestión también escribe un evento. Un comentario de la tabla vieja no genera evento.

## Shape de un registro

```json
{
  "id": "...",
  "caseId": "...",
  "scope": "municipal",
  "assignedUserId": "...",
  "assignedUserName": "Ana Pérez",
  "departmentName": "San Salvador",
  "cityName": "San Salvador",
  "townName": "San Salvador Centro",
  "salviaCommitment": "...",
  "gestiones": [
    {
      "id": "...",
      "result": "Abordado en sesión: Sí, sesión ordinaria. Reporte de gestión: …",
      "formData": { "abordadoEnSesion": "si_ordinaria", "reporteGestion": "…" },
      "completedAt": "...",
      "assignedUserName": "Ana Pérez"
    }
  ]
}
```
