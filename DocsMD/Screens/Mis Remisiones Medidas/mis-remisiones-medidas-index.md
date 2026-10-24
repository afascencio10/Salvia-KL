# Mis Remisiones — Medidas de Emergencia — Índice

Bandeja del agente. Monta el mismo [`remisiones-medidas-component`](../../Componentes/Remisiones%20Medidas%20component/remisiones-medidas-component-events.md) con `vista="agente"`. Ruta `GET /salvia/mis-remisiones-medidas`.

El recorte (solo lo asignado a ese agente, y solo cinco estados) está en E-01 del componente. Esta pantalla no vuelve a documentar filtros ni cards.

## Navegación

| Documento | Link |
|---|---|
| Carga de la pantalla | [Flujos/flow-E01-cuando-carga-pantalla.md](./Flujos/flow-E01-cuando-carga-pantalla.md) |
| Componente | [remisiones-medidas-component-events.md](../../Componentes/Remisiones%20Medidas%20component/remisiones-medidas-component-events.md) |
| Interfaz del componente | [remisiones-medidas-component-interface.md](../../Componentes/Remisiones%20Medidas%20component/remisiones-medidas-component-interface.md) |
| Pantalla general | [remisiones-medidas-index.md](../Remisiones%20Medidas/remisiones-medidas-index.md) |

## Ruta, permiso y menú

| Aspecto | Valor |
|---|---|
| Ruta | `GET /salvia/mis-remisiones-medidas` |
| Handler | `MisRemisionesMedidasGET` en `src/salvia/facades/MisRemisionesMedidasFacade.go` |
| Template | `src/frontend/html/salvia/mis-remisiones-medidas/mis_remisiones_medidas.html` |
| Título | Mis Remisiones — Medidas de Emergencia |
| Subtítulo | Remisiones asignadas a tu perfil |
| Permiso | `get_mis_remisiones_medidas`: solo `sv`. El resto recibe 401 |
| Menú | Ítem **Mis Remisiones Medidas**, dentro del bloque `sv` de `sidebar.html` |

El acceso de `sv` es temporal. Cuando exista el rol de Agente de Medidas, el permiso y el menú pasan a ese rol, y el id del agente sale de la sesión.

## Qué monta

```html
<remisiones-medidas-component
  vista="agente"
  agente-id="am-1"
  @ver-caso="onVerCaso"
  @ver-detalle="onVerDetalle"
></remisiones-medidas-component>
```

`am-1` es el id del seed (Pedro Ramírez). En el seed, Pedro no tiene filas dentro de los cinco estados de la bandeja, así que la pantalla abre en "No tienes remisiones asignadas." Al cablear el API, `agente-id` pasa a ser `general_user_i_code` de la sesión y el query deja de aceptarlo desde el cliente (E-01 del componente).

`onVerCaso` y `onVerDetalle` están vacíos.

## Checklist de completitud

- [x] ¿La carga de la pantalla está cubierta? → E-01 de esta pantalla
- [x] ¿El recorte del agente está cubierto? → E-01 del componente
- [ ] ¿`agente-id` sale de la sesión? → Pendiente, junto con el rol
- [ ] ¿`onVerCaso` y `onVerDetalle` navegan? → Pendiente
