# `case-mecanismos` — Guía de uso

Lista y registra los escalamientos del caso a usuarios con rol Mecanismo (`mc`). El select del modal no es un catálogo de nombres: carga usuarios activos filtrados por la cadena `profile.town → city → department`.

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

Al crear un registro se escribe un evento en el Timeline (categoría `Mecanismos`). Un comentario no genera evento.

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
  "comments": [
    { "id": "...", "authorName": "...", "body": "...", "createdAt": "..." }
  ]
}
```
