# Remisiones Medidas de Emergencia — Índice

Pantalla de gestión que monta [`remisiones-medidas-component`](../../Componentes/Remisiones%20Medidas%20component/remisiones-medidas-component-events.md). Ruta `GET /salvia/remisiones-medidas`.

La vista sale de `currentRole`. El listado todavía usa el seed del componente. El contrato de lectura y de los dos guardados está en los eventos del componente.

## Navegación

| Documento | Link |
|---|---|
| Carga de la pantalla | [Flujos/flow-E01-cuando-carga-pantalla.md](./Flujos/flow-E01-cuando-carga-pantalla.md) |
| Componente | [remisiones-medidas-component-events.md](../../Componentes/Remisiones%20Medidas%20component/remisiones-medidas-component-events.md) |
| Detalle | [detalle-remision-medida-index.md](../Detalle%20Remision%20Medida/detalle-remision-medida-index.md) |
| Interfaz del componente | [remisiones-medidas-component-interface.md](../../Componentes/Remisiones%20Medidas%20component/remisiones-medidas-component-interface.md) |

## Ruta, permiso y menú

| Aspecto | Valor |
|---|---|
| Ruta | `GET /salvia/remisiones-medidas` |
| Handler | `RemisionesMedidasGET` en `src/salvia/facades/RemisionesMedidasFacade.go` |
| Template | `src/frontend/html/salvia/remisiones-medidas/remisiones_medidas.html` |
| Título | Remisiones Medidas de Emergencia |
| Subtítulo | Gestión de remisiones al equipo de medidas de emergencia |
| Permiso | `get_remisiones_medidas`: `sm`, `em`, `sc`, `cr`, `sv` |
| Menú | Ítem **Remisiones Medidas** para esos cinco roles |

| Rol | Vista que monta | Reasignar |
|---|---|---|
| `sm` | `supervisor` | Sí (`puedeReasignar`, permiso `reassign_remisiones_medidas`) |
| `sv` | `supervisor` | Sí. También puede revisar (`review_remisiones_medidas`) |
| `em` | `equipo-tecnico` | No. Sí puede revisar |
| `sc` | `supervision-convenio` | No |
| `cr` | `cruz-roja` | No |

## Qué monta

El handler traduce el rol y el template recibe `vista`. Para `sm` y `sv` también envía `:puede-reasignar="true"`.

`onVerDetalle` navega a `/salvia/remision-medida/{remisionId}` (E-15). `onVerCaso` sigue vacío (E-14).

Filtros, cards, selección, reasignación y revisión técnica no se documentan otra vez: son E-01 a E-13 del componente.

## Checklist de completitud

- [x] ¿La carga de la pantalla está cubierta? → E-01 de esta pantalla
- [x] ¿El listado y sus acciones están cubiertos? → eventos del componente
- [x] ¿Cada rol entra con su vista? → tabla de arriba
- [x] ¿`onVerDetalle` navega? → `/salvia/remision-medida/:id`
- [ ] ¿`onVerCaso` navega? → Pendiente (E-14)
