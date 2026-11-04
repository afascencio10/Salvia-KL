# Consulta Reporte Pública — Tablas y reglas

**Esta pantalla no crea tablas ni columnas.** Solo lee, con un SELECT parametrizado (`FindByRadicado`), lo mínimo de un reporte V2: `victim_contact_radicado` (creada por WP-1 T3, ver [`Registro de Contacto Víctima V2/related-tables.md`](../Registro%20de%20Contacto%20V%C3%ADctima%20V2/related-tables.md)).

## Tablas leídas

| Tabla | Para qué |
|---|---|
| `salvia.victim_contact` | radicado (búsqueda exacta), día de radicación (`TO_CHAR(victim_contact_creation_date, 'YYYY-MM-DD')`; el servidor fija la hora local en America/Bogota), agente asignado, estado legado (`'i'` = invalidado) |
| `salvia.victim_case_form2_enums` | código del motivo de cierre (solo para saber si existe; nunca se devuelve) |
| `salvia.victim_contact_attempt` | cantidad de intentos no eliminados (cualquier resultado) |
| `salvia.victim_case` | caso de Registro V3 vinculado más reciente (`victim_case_victim_contact`) en estado `bo`, `in`, `te` o `cd` |

## Regla del estado público

| Estado | Cuándo |
|---|---|
| `cerrado` | con motivo de cierre, invalidado en el flujo legado, o con el caso vinculado en `te`/`cd` (derivado, «Cerrado — válido») |
| `en_gestion` | sin cierre y con agente asignado, con al menos un intento, o con un borrador de caso iniciado (`bo`/`in`) |
| `recibido` | el resto (sin agente, sin intentos y sin borrador) |

Es la misma regla del estado derivado del reporte que usan el modal y el listado (ver [`Listado Reportes/related-tables.md`](../Listado%20Reportes/related-tables.md)); un caso de otro flujo (estados distintos de `bo`, `in`, `te`, `cd`) no cambia el estado.

## Qué devuelve (y qué no)

`{radicado, estado, estadoTexto, descripcion, fechaRadicacion}`. Nunca nombres, teléfono, descripción, tipo de reporte, agente, equipo, motivo de cierre ni `icode`.

## Protección contra consultas masivas

- Formato inválido y radicado inexistente devuelven el **mismo** 404 con el mismo cuerpo.
- Limitador **en memoria por IP**: 5 fallos en 15 minutos → 429 con `Retry-After`; los aciertos no cuentan y un 429 no consulta la BD. Se reinicia con el servicio y no se comparte entre réplicas.
- La IP sale de `c.ClientIP()`: detrás de un proxy depende de los proxies de confianza de Gin (con la configuración por defecto, un cliente podría falsear `X-Forwarded-For`; configurar `SetTrustedProxies` en el despliegue).
- El radicado tiene 8 símbolos de 32, 2 de ellos derivados de la fecha: unas 10⁹ combinaciones por día; con el límite de intentos adivinarlo es inviable. Sin captcha (decisión D1 de la planeación).
