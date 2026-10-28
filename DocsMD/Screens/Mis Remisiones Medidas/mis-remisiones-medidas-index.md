# Mis Remisiones — Medidas de Emergencia — Índice

Bandeja del agente de medidas (`am`). Monta el mismo [`remisiones-medidas-component`](../../Componentes/Remisiones%20Medidas%20component/remisiones-medidas-component-events.md) con `vista="agente"`. Ruta `GET /salvia/mis-remisiones-medidas`.

El recorte (solo lo asignado a ese agente, y solo cinco estados) está en E-01 del componente. Esta pantalla no vuelve a documentar filtros ni cards.

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
| Subtítulo | Remisiones asignadas a tu perfil |
| Permiso | `get_mis_remisiones_medidas`: solo `am` |
| Menú | Ítem **Mis Remisiones Medidas**, solo para `am` |
| Entrada al iniciar sesión | `/salvia/mis-remisiones-medidas` |

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
