# Remisiones Medidas de Emergencia — Índice

Pantalla de gestión que monta [`remisiones-medidas-component`](../../Componentes/Remisiones%20Medidas%20component/remisiones-medidas-component-events.md) con `vista="supervisor"`. Ruta `GET /salvia/remisiones-medidas`.

El listado todavía usa el seed del componente. El contrato de lectura y de los dos guardados está en los eventos del componente. El mapeo `currentRole` → `vista` queda pendiente: esta pantalla no cambia de vista según el rol.

## Navegación

| Documento | Link |
|---|---|
| Carga de la pantalla | [Flujos/flow-E01-cuando-carga-pantalla.md](./Flujos/flow-E01-cuando-carga-pantalla.md) |
| Componente | [remisiones-medidas-component-events.md](../../Componentes/Remisiones%20Medidas%20component/remisiones-medidas-component-events.md) |
| Interfaz del componente | [remisiones-medidas-component-interface.md](../../Componentes/Remisiones%20Medidas%20component/remisiones-medidas-component-interface.md) |

## Ruta, permiso y menú

| Aspecto | Valor |
|---|---|
| Ruta | `GET /salvia/remisiones-medidas` |
| Handler | `RemisionesMedidasGET` en `src/salvia/facades/RemisionesMedidasFacade.go` |
| Template | `src/frontend/html/salvia/remisiones-medidas/remisiones_medidas.html` |
| Título | Remisiones Medidas de Emergencia |
| Subtítulo | Gestión de remisiones al equipo de medidas de emergencia |
| Permiso | `get_remisiones_medidas`: solo `sv`. El resto recibe 401 |
| Menú | Ítem **Remisiones Medidas**, dentro del bloque `sv` de `sidebar.html` |

El acceso de `sv` es temporal, para poder abrir la pantalla mientras no existan los roles de Medidas.

## Qué monta

```html
<remisiones-medidas-component
  vista="supervisor"
  @ver-caso="onVerCaso"
  @ver-detalle="onVerDetalle"
></remisiones-medidas-component>
```

`onVerCaso` y `onVerDetalle` están vacíos. E-14 y E-15 del componente emiten y esta pantalla no redirige.

Filtros, cards, selección, reasignación y revisión técnica no se documentan otra vez: son E-01 a E-13 del componente.

## Checklist de completitud

- [x] ¿La carga de la pantalla está cubierta? → E-01 de esta pantalla
- [x] ¿El listado y sus acciones están cubiertos? → eventos del componente
- [ ] ¿El permiso deja de ser solo `sv` cuando exista el rol? → Pendiente
- [ ] ¿`onVerCaso` y `onVerDetalle` navegan? → Pendiente
