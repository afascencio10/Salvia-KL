# Mis Remisiones — Medidas de Emergencia — Índice

Bandeja propia de medidas. La montan el agente de medidas (`am`) y el agente de seguimiento (`ro`), con el mismo [`remisiones-medidas-component`](../../Componentes/Remisiones%20Medidas%20component/remisiones-medidas-component-events.md) en `vista="agente"`. Ruta `GET /salvia/mis-remisiones-medidas`.

`am` ve lo asignado a su icode y solo cinco estados. `ro` ve las remisiones cuyo `submitted_by` es su icode, en cualquier estado. El API distingue esas dos con las vistas `agente` y `remitente`. Esta pantalla no vuelve a documentar filtros ni cards.

## Navegación

| Documento | Link |
|---|---|
| Carga de la pantalla | [Flujos/flow-E01-cuando-carga-pantalla.md](./Flujos/flow-E01-cuando-carga-pantalla.md) |
| Componente | [remisiones-medidas-component-events.md](../../Componentes/Remisiones%20Medidas%20component/remisiones-medidas-component-events.md) |
| Interfaz del componente | [remisiones-medidas-component-interface.md](../../Componentes/Remisiones%20Medidas%20component/remisiones-medidas-component-interface.md) |
| Pantalla general | [remisiones-medidas-index.md](../Remisiones%20Medidas/remisiones-medidas-index.md) |
| Detalle | [detalle-remision-medida-index.md](../Detalle%20Remision%20Medida/detalle-remision-medida-index.md) |

## Ruta, permiso y menú

| Aspecto | Valor |
|---|---|
| Ruta | `GET /salvia/mis-remisiones-medidas` |
| Handler | `MisRemisionesMedidasGET` en `src/salvia/facades/MisRemisionesMedidasFacade.go` |
| Template | `src/frontend/html/salvia/mis-remisiones-medidas/mis_remisiones_medidas.html` |
| Título | Mis Remisiones — Medidas de Emergencia |
| Permiso | `get_mis_remisiones_medidas`: `am` y `ro` |
| Menú | Ítem **Mis Remisiones Medidas**, para `am` y `ro` |
| Entrada al iniciar sesión | `/salvia/mis-remisiones-medidas` solo para `am`. El login de `ro` no cambia |
| Subtítulo | `am`: Remisiones asignadas a tu perfil. `ro`: Remisiones que enviaste |

## Qué monta

```html
<remisiones-medidas-component
  vista="agente"
  agente-id="{general_user_i_code de la sesión}"
  @ver-caso="onVerCaso"
  @ver-detalle="onVerDetalle"
  user-id="{icode de la sesión}"
  user-role="{role_code de la sesión}"
></remisiones-medidas-component>
```

El seed usa ids de mentira (`am-1`, `am-2`, `am-3`). Con el icode real de la sesión la bandeja abre vacía hasta que el listado lea `emergency_measure_referral`.

`onVerDetalle` navega a `/salvia/remision-medida/{remisionId}` (E-15). `onVerCaso` sigue vacío (E-14). "Tareas" abre el panel de la remisión (E-16). El agente ve la tarea `intentar_contacto_3x3` cuando está asignada a su usuario. La de decidir vía sigue siendo de `sm`.

## Checklist de completitud

- [x] ¿La carga de la pantalla está cubierta? → E-01 de esta pantalla
- [x] ¿El recorte del agente está cubierto? → E-01 del componente
- [x] ¿`agente-id` sale de la sesión? → `s.UserICode`
- [x] ¿`onVerDetalle` navega? → `/salvia/remision-medida/:id`
- [ ] ¿`onVerCaso` navega? → Pendiente (E-14)
