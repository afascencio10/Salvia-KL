# Administrar Duplas — Index

Pantalla de administración de duplas del equipo de Atención Psicosocial. Permite crear, editar y eliminar (lógico) pares psicóloga (`ps`) + trabajadora social (`ts`). Una **psicóloga** solo puede pertenecer a **una** dupla activa; una **trabajadora social** puede pertenecer a **varias**.

> **Ruta:** `GET /salvia/administrar-duplas` (solo `sv`)  
> **Prerequisito de modelo:** tabla `salvia.dupla` (ver [M-01](../../Componentes/Remisiones%20Atencion%20Psicosocial%20component/Flujos/flow-M01-migracion-schema-dupla-psychosocial-support.md)). Esta pantalla es la que **persiste** las duplas; el listado de remisiones y el modal de reasignación solo las **leen**.  
> **Estado UI:** interfaz montada; carga/guardado/eliminación vía API pendiente (flujos E01–E07).

---

## Navegación

| Documento | Link |
|---|---|
| Interfaz | [administrar-duplas-interface.md](./administrar-duplas-interface.md) |
| Contexto remisiones / dupla | [reasignar-remisiones-modal-events.md](../../Componentes/Remisiones%20Atencion%20Psicosocial%20component/reasignar-remisiones-modal-events.md) |

---

## Roles con acceso

| Rol | Código | Alcance |
|---|---|---|
| Supervisor | `sv` | Único rol con acceso. Gestiona (CRUD lógico) todas las duplas |

---

## Reglas de negocio

1. Una dupla = 1 psicóloga (`role_code = 'ps'`) + 1 trabajadora social (`role_code = 'ts'`).
2. Solo usuarios con `general_user_status = 'e'` (activos) aparecen en listas y selects.
3. Una **psicóloga** solo puede estar en **una** dupla con `deleted_at IS NULL`.
4. Una **trabajadora social** puede pertenecer a **varias** duplas activas a la vez.
5. Al editar, la psicóloga actual de esa dupla sigue disponible en el select (aunque esté “ocupada” por sí misma).
6. Eliminar es **lógico** (`deleted_at`); libera a la psicóloga para otra dupla. No borra usuarios ni remisiones.
7. El **nombre es libre** (lo escribe el supervisor); máximo `varchar(36)`. Debe ser **único entre duplas activas** (`deleted_at IS NULL`). Si ya existe → error `"nombre de la dupla en uso"`.
8. **No se puede eliminar** una dupla si está en uso en una sesión/remisión no cerrada:
   - Existe `psychosocial_support` con `dupla_id` = esa dupla y `status <> 'cerrado'`, **o**
   - Existe `team_contact` con `dupla_id` = esa dupla cuyo `psychosocial_support` asociado (`psicosocial_id`) tiene `status <> 'cerrado'`.
   - En ese caso se muestra un **modal de error** (no se hace soft-delete).
   - Si solo hay referencias con remisión en `cerrado`, sí se puede eliminar (el `dupla_id` histórico puede quedar).

---

## Resumen de eventos

| # | Evento | Tipo | Flujo |
|---|---|---|---|
| E01 | Cuando carga la pantalla | Lifecycle | [📄 flow-E01-cuando-carga-pantalla.md](./Flujos/flow-E01-cuando-carga-pantalla.md) |
| E02 | Cuando abre modal crear o editar dupla | User Interaction | [📄 flow-E02-cuando-abre-modal-crear-o-editar.md](./Flujos/flow-E02-cuando-abre-modal-crear-o-editar.md) |
| E03 | Cuando cancela modal crear o editar | User Interaction | [📄 flow-E03-cuando-cancela-modal-crear-o-editar.md](./Flujos/flow-E03-cuando-cancela-modal-crear-o-editar.md) |
| E04 | Cuando guarda modal | User Interaction | [📄 flow-E04-cuando-guarda-modal.md](./Flujos/flow-E04-cuando-guarda-modal.md) |
| E05 | Cuando abre confirmación de eliminar | User Interaction | [📄 flow-E05-cuando-abre-confirmacion-eliminar.md](./Flujos/flow-E05-cuando-abre-confirmacion-eliminar.md) |
| E06 | Cuando cancela confirmación de eliminar | User Interaction | [📄 flow-E06-cuando-cancela-confirmacion-eliminar.md](./Flujos/flow-E06-cuando-cancela-confirmacion-eliminar.md) |
| E07 | Cuando confirma eliminar dupla | User Interaction | [📄 flow-E07-cuando-confirma-eliminar-dupla.md](./Flujos/flow-E07-cuando-confirma-eliminar-dupla.md) |
| E08 | Cuando cierra modal de error al eliminar | User Interaction | [📄 flow-E08-cuando-cierra-modal-error-eliminar.md](./Flujos/flow-E08-cuando-cierra-modal-error-eliminar.md) |

---

## Inventario de eventos

---

**Nombre del evento:** Cuando carga la pantalla  
**Tipo:** Lifecycle  
**Descripción:** Al montar el componente, valida acceso por rol, consulta profesionales activos `ps`/`ts` y las duplas activas (`deleted_at IS NULL`). Enriquece listas con badge “En dupla” / “Disponible” y renderiza las tres secciones de la UI.  
**Requerido:** Sí  

📄 [Ver flujo → flow-E01-cuando-carga-pantalla.md](./Flujos/flow-E01-cuando-carga-pantalla.md)

---

**Nombre del evento:** Cuando abre modal crear o editar dupla  
**Tipo:** User Interaction  
**Descripción:** Disparado por “+ Nueva dupla” o “Editar”. Determina modo `create` | `edit`, precarga el formulario si edita, calcula psicólogas disponibles (excluyendo las ya asignadas a otras duplas) y deja **todas** las TS seleccionables.  
**Requerido:** Sí  

📄 [Ver flujo → flow-E02-cuando-abre-modal-crear-o-editar.md](./Flujos/flow-E02-cuando-abre-modal-crear-o-editar.md)

---

**Nombre del evento:** Cuando cancela modal crear o editar  
**Tipo:** User Interaction  
**Descripción:** El usuario presiona “Cancelar”, ✕ o el backdrop. Cierra el modal sin persistir y limpia estado del formulario.  
**Requerido:** Sí  

📄 [Ver flujo → flow-E03-cuando-cancela-modal-crear-o-editar.md](./Flujos/flow-E03-cuando-cancela-modal-crear-o-editar.md)

---

**Nombre del evento:** Cuando guarda modal  
**Tipo:** User Interaction  
**Descripción:** Valida nombre + ps + ts. En backend verifica unicidad de **psicóloga** (no de TS) y unicidad de nombre entre duplas activas (`"nombre de la dupla en uso"`). Crea (`POST`) o actualiza (`PUT`) según el modo. Si ok, cierra el modal y recarga datos (E01).  
**Requerido:** Sí  

📄 [Ver flujo → flow-E04-cuando-guarda-modal.md](./Flujos/flow-E04-cuando-guarda-modal.md)

---

**Nombre del evento:** Cuando abre confirmación de eliminar  
**Tipo:** User Interaction  
**Descripción:** El usuario presiona “Eliminar” en una fila de duplas activas. Guarda la dupla objetivo y abre el modal de confirmación con el nombre dinámico.  
**Requerido:** Sí  

📄 [Ver flujo → flow-E05-cuando-abre-confirmacion-eliminar.md](./Flujos/flow-E05-cuando-abre-confirmacion-eliminar.md)

---

**Nombre del evento:** Cuando cancela confirmación de eliminar  
**Tipo:** User Interaction  
**Descripción:** El usuario presiona “Cancelar” o cierra el modal de confirmación. No modifica BD; limpia la dupla objetivo.  
**Requerido:** Sí  

📄 [Ver flujo → flow-E06-cuando-cancela-confirmacion-eliminar.md](./Flujos/flow-E06-cuando-cancela-confirmacion-eliminar.md)

---

**Nombre del evento:** Cuando confirma eliminar dupla  
**Tipo:** User Interaction  
**Descripción:** Valida en backend que la dupla no esté en uso en `psychosocial_support` / `team_contact` con remisión no `cerrado`. Si está en uso → modal de error. Si no → soft-delete, libera miembros y recarga listas.  
**Requerido:** Sí  

📄 [Ver flujo → flow-E07-cuando-confirma-eliminar-dupla.md](./Flujos/flow-E07-cuando-confirma-eliminar-dupla.md)

---

**Nombre del evento:** Cuando cierra modal de error al eliminar  
**Tipo:** User Interaction  
**Descripción:** El usuario cierra el modal de error (“no se puede eliminar porque se está utilizando en una sesión”). Solo limpia el estado del modal de error; no cambia BD.  
**Requerido:** Sí  

📄 [Ver flujo → flow-E08-cuando-cierra-modal-error-eliminar.md](./Flujos/flow-E08-cuando-cierra-modal-error-eliminar.md)

---

## Checklist de validación

- [x] ¿Cada acción del usuario puede ser manejada?
- [x] ¿La carga inicial de datos está cubierta?
- [x] ¿Crear / editar / eliminar (lógico) están incluidos?
- [x] ¿Cancelar modales está cubierto?
- [x] ¿Rol de acceso confirmado? — Solo `sv`
- [x] ¿Regla de eliminación con remisiones/sesiones en uso? — E07
- [x] ¿Unicidad de nombre? — E04 (`"nombre de la dupla en uso"`)
- [ ] ¿Hay actualizaciones en tiempo real? — No aplica

---

## Relación con otros módulos

```
Administrar Duplas (esta pantalla)
        │  escribe salvia.dupla
        ▼
salvia.dupla  ──lectura──►  Remisiones Psicosocial (filtro dupla)
                          ►  ReasignarRemisionesModal (RRM-03 GET /duplas/reasignacion)
                          ►  psychosocial_support.dupla_id  (bloqueo delete si status ≠ cerrado)
                          ►  team_contact.dupla_id          (bloqueo delete vía remisión asociada)
```
