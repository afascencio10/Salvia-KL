# Detalle de Remisión — Medida de Emergencia — Índice

Pantalla de lectura de una remisión. Se llega desde "Ver detalle →" y desde los botones de fila de tipo detalle (E-15 del componente). Ruta `GET /salvia/remision-medida/:id`. El `:id` es el uuid de `emergency_measure_referral`.

No hay select de rol. `currentRole` decide el recorte de lectura y si aparece "Reasignar agente de medidas". El tab Tareas filtra por la remisión. Timeline sigue con el `caseICode`.

## Navegación

| Documento | Link |
|---|---|
| Carga de la pantalla | [Flujos/flow-E01-cuando-carga-pantalla.md](./Flujos/flow-E01-cuando-carga-pantalla.md) |
| Reasignar agente | [Flujos/flow-E02-cuando-reasigna-agente.md](./Flujos/flow-E02-cuando-reasigna-agente.md) |
| Interfaz | [detalle-remision-medida-interface.md](./detalle-remision-medida-interface.md) |
| Origen del clic | [flow-E15-cuando-presiona-ver-detalle.md](../../Componentes/Remisiones%20Medidas%20component/Flujos/flow-E15-cuando-presiona-ver-detalle.md) |
| Bandeja general | [remisiones-medidas-index.md](../Remisiones%20Medidas/remisiones-medidas-index.md) |
| Bandeja del agente | [mis-remisiones-medidas-index.md](../Mis%20Remisiones%20Medidas/mis-remisiones-medidas-index.md) |
| Changelog octubre 2026 | [changelogOct2026.md](./changelogOct2026.md) |

## Ruta y permiso

| Aspecto | Valor |
|---|---|
| Ruta | `GET /salvia/remision-medida/:id` |
| Handler | `RemisionMedidaDetalleGET` en `src/salvia/facades/RemisionMedidaDetalleFacade.go` |
| Template | `src/frontend/html/salvia/remision-medida/remision_medida_detalle.html` |
| API | `GET /api/v1/emergency-measure-referrals/:id` |
| Título | Detalle de Remisión — Medida de Emergencia |
| Permiso | `get_remision_medida_detalle`: `sm`, `em`, `sc`, `cr`, `sv`, `am` |
| Menú | No tiene ítem. Se entra desde las bandejas |

| Rol | Qué puede leer | Reasignar |
|---|---|---|
| `sm`, `sv` | Cualquier remisión | Sí (`reassign_remisiones_medidas`) |
| `am` | La suya, y solo los cinco estados de su bandeja | No |
| `em` | Solo `pendiente_de_revision` | No |
| `sc` | Trazabilidad, cotizaciones y escalar a Cruz Roja | No |
| `cr` | Aprobadas, por cotizar y entregar medida | No |

Si la remisión no entra en ese recorte, el API responde 404. El mismo `scopeClause` del listado.

Volver: `am` cae en `/salvia/mis-remisiones-medidas`. El resto, en `/salvia/remisiones-medidas`.

## Pendiente — ciclo de gestión

E-15 deja en esta pantalla decidir la vía, escalar, entregar la medida, revisar la caracterización, la trazabilidad, las cotizaciones y confirmar la seguridad. Esta entrega no las implementa. Info Medida solo lee. Tareas y Timeline no son ese ciclo: son `case-tasks` y `case-timeline` del caso.

Tampoco se pintan la descripción de la situación ni las observaciones del mockup: no tienen columna. "Ver caso" (E-14) sigue sin navegar.

## Checklist de completitud

- [x] ¿La carga está cubierta? → E-01
- [x] ¿El recorte por rol está cubierto? → E-01
- [x] ¿Reasignar desde el detalle? → E-02, solo `sm` y `sv`
- [x] ¿Tareas y Timeline reutilizan los del caso? → E-01
- [ ] ¿El ciclo de gestión de E-15 está implementado? → Pendiente
- [ ] ¿"Ver caso" navega? → Pendiente (E-14)
