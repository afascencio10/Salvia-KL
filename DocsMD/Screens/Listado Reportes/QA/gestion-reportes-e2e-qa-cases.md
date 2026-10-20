# QA — Gestión de reportes · E2E de un reporte nuevo

Un reporte **nuevo** recorre de punta a punta las tres piezas de "Gestión de reportes": nace en el formulario público (`/salvia/public-v2`, PSWYA-755), el supervisor lo ve "Sin asignar" en [Listado Reportes](../listado-reportes-index.md) y lo reasigna ([E-01](../Flujos/flow-E01-cuando-carga-pantalla.md), [E-04](../Flujos/flow-E04-cuando-presiona-ver-reporte.md), [E-05](../Flujos/flow-E05-cuando-presiona-reasignar-reportes.md), [E-06](../Flujos/flow-E06-cuando-confirma-reasignacion.md)), el agente lo gestiona en [Mis Reportes](../../Mis%20Reportes/mis-reportes-index.md) con el modal [`gestion-reporte-modal`](../../../Componentes/Gestion%20Reporte%20Modal/gestion-reporte-modal-index.md) ([M-04](../../../Componentes/Gestion%20Reporte%20Modal/Flujos/flow-M04-cuando-procede-recontacto.md), [M-05](../../../Componentes/Gestion%20Reporte%20Modal/Flujos/flow-M05-cuando-registra-intento-fallido.md), [M-06](../../../Componentes/Gestion%20Reporte%20Modal/Flujos/flow-M06-cuando-registra-intento-exitoso.md), [M-09](../../../Componentes/Gestion%20Reporte%20Modal/Flujos/flow-M09-cuando-confirma-cierre-decision-persona.md)) y el supervisor ve el resultado.

> **Estado: ejecutado el 2026-10-13** — resultados al final de este archivo. Datos, usuarios, helpers y decisiones: [`gestion-reportes-qa-plan.md`](gestion-reportes-qa-plan.md).

**Método.** Un solo spec con un solo test: los cuatro pasos dependen del mismo reporte, así que van como `test.step` y no como tests sueltos. No repite validaciones ya cubiertas en los otros specs (permisos, filtros, esquema 3x3 completo, límites): comprueba que **las piezas se encadenan** y que el dato que deja una lo ve la siguiente. Es la única prueba de la feature que crea el reporte por el wizard (los demás lo crean por el `POST` real de `/salvia/public-v2`).

**Cómo se mantiene repetible** (se puede ejecutar dos, o veinte, veces seguidas sin preparar nada):

| Riesgo al repetir | Mecanismo |
|---|---|
| Datos que chocan entre corridas | El apellido lleva un id de corrida (`Ciclo-<runId>`) y el reporte se busca por él: la tabla muestra una sola fila. El teléfono es aleatorio |
| Reporte, intentos y formularios que quedan en la BD | El fixture `reportes` borra el reporte en cascada al terminar el test, aunque falle (`limpiarReporteQA`). Se registra apenas el wizard devuelve el radicado |
| Corrida interrumpida (Ctrl+C, `kill -9`, caída del equipo) | `beforeAll` barre los reportes `QAGRE2E%` huérfanos; `afterAll` vuelve a barrer y **exige que no haya nada** (`GR-E2E-C5`) |
| Hace falta un supervisor y un agente | El spec crea `qagr.sv` y `qagr.op` (`scripts/usuarios-qagr.sh crear`, con contraseña aleatoria solo en `.env.local`) y las borra al terminar. Si ya hay cuentas `GR_*` creadas a mano, las usa y no las toca |
| Cuentas de una corrida colgada | El bloque `GR_*` que deja el spec lleva la nota `QAGR_AUTOGEN=1`: la siguiente corrida las recrea limpias y `afterAll` las retira |
| El servidor no deja iniciar sesión | `beforeAll` verifica el login de ambas cuentas por API y reintenta hasta 90 s (`esperarInicioSesion`); si no, falla ahí y con el motivo (ver H9 del plan) |

**No ejecutar dos copias a la vez:** iniciar sesión con un login que ya tiene sesión activa cierra la anterior.

---

## Casos de prueba

| Caso | Descripción | Condición inicial | Resultado esperado | Método de verificación |
|---|---|---|---|---|
| GR-E2E-C1 | Radicación por el formulario público | Sin sesión; wizard de 4 pasos con un reporte tipo Víctima (autoriza, "No requiere" ajuste, nombres `QAGRE2E`, apellido único, teléfono, "Sí" a la llamada, 10:00, descripción, captcha `0000` que se omite por el `User-Agent` de la app móvil) | `POST /salvia/public-v2` responde 200 con el radicado y la pantalla final lo muestra. BD: el reporte existe con esos nombres y apellido, sin agente y sin motivo de cierre | UI + respuesta interceptada + BD |
| GR-E2E-C2 | Supervisor: aparece "Sin asignar" y se reasigna | Sesión `SV` en Listado Reportes; búsqueda por el apellido único | Una sola fila con el apellido y `ID: <radicado>`, tipo "Víctima", "Sin intentos", "Asignado a" y estado "Sin asignar". Marcada la fila, el botón dice "Reasignar Reportes"; elegir a `OP1` y confirmar emite un `POST reasignar-bulk` con `{report_icodes:[radicado], agent_icode:<OP1>}` (200); la fila pasa a "QAGR Agente QA" y "Asignado". BD: agente = `OP1` | UI + request interceptada + BD |
| GR-E2E-C3 | Agente: lo gestiona y lo cierra | Sesión `OP1` en Mis Reportes | El reporte está en su lista como "Asignado" y "Sin intentos". El modal abre con el radicado y el apellido; "Procede recontacto" muestra el contador `Hoy: 0/3 · Días completos: 0/3`; un intento fallido lo deja en `Hoy: 1/3` con una fila en el log; el exitoso pasa a "Contacto establecido" con solo dos formas de cerrar; "Cierre por decisión de la persona" pide el texto (el botón de confirmar sigue deshabilitado hasta escribirlo) y emite `POST cerrar` con `{close_reason:'dp', close_observation}` (200). El modal dice "Reporte <radicado> cerrado" y "Inválido por Cierre por decisión de la persona — <observación>"; la fila queda "Cerrado - inválido" y "1/9 intentos · 0/3 días". BD: motivo `dp`, observación, fecha de cierre y dos intentos del agente (fallido, exitoso), ninguno borrado | UI + requests interceptadas + BD |
| GR-E2E-C4 | Supervisor ve el resultado | La misma sesión `SV`; recarga Listado Reportes y busca el apellido | La fila está "Cerrado - inválido", "Asignado a" `QAGR Agente QA` y "1/9 intentos · 0/3 días". "Ver Reporte" abre en `cerrado` con el mismo motivo que vio el agente y un solo botón (solo lectura). Ninguna de las tres páginas (público, `SV`, `OP1`) lanzó una excepción de JavaScript | UI + eventos `pageerror` |
| GR-E2E-C5 | Repetible: sin residuo | Terminada la corrida (y tras la anterior) | `afterAll` barre `QAGRE2E%` y exige 0. Además, tras cada corrida no quedan cuentas `qagr.*` ni bloque `GR_*` (salvo las manuales), ni intentos huérfanos; el conteo de filas de todas las tablas legibles de `salvia` y `security` es idéntico al de antes | Aserción de `afterAll` + censo de filas por tabla (ver resultados) |

---

## Cobertura

| Flujo de prueba | Casos cubiertos | Spec |
|---|---|---|
| `TU-GR-E2E-1` — ciclo de vida de un reporte nuevo: radicación → reasignación → gestión 3x3 y cierre → lectura del supervisor → sin residuo | GR-E2E-C1, C2, C3, C4, C5 | `gestion-reportes-e2e/ciclo-reporte-nuevo.spec.ts` |

## Cómo ejecutarlo

```bash
cd qa-salvia
# backend en :9090:  cd ../SOG_SALVIA-develop/src && PORT=9090 go run -a .
CI=1 npx playwright test tests/gestion-reportes-e2e --retries=0     # una vez
CI=1 npx playwright test tests/gestion-reportes-e2e --retries=0     # y otra vez, seguida
```

Sin `CI=1` el navegador se abre visible (`playwright.config.ts`). No hace falta crear cuentas ni limpiar nada antes ni después. Dura ≈ 25 s.

**Huella que sí deja:** una fila de auditoría en `salvia.app_log` por corrida (`Reasignación manual de reporte (bulk): agente actualizado`, escrita por el servicio de reasignación). El rol de la BD de pruebas no tiene permiso sobre esa tabla, así que el test no la puede ver ni borrar; no afecta a las corridas siguientes. Para limpiarlas (con un rol que tenga permiso):

```sql
DELETE FROM salvia.app_log
WHERE message LIKE 'Reasignación manual de reporte (bulk)%' AND data->>'actorName' = 'QAGR Supervisor QA';
```

---

## Resultados de ejecución

**Fecha:** 2026-10-13 · **Build:** `SOG_SALVIA-develop` (árbol de trabajo) en `http://localhost:9090` · **BD:** Supabase de pruebas · **Herramienta:** Playwright headless (`CI=1`), 1 worker, sin reintentos · **Repo de pruebas:** `qa-salvia/`

| Caso | Resultado | Prueba (spec en `qa-salvia/tests/`) | Hallazgo |
|---|---|---|---|
| GR-E2E-C1 | ✅ | `gestion-reportes-e2e/ciclo-reporte-nuevo.spec.ts` | — |
| GR-E2E-C2 | ✅ | `gestion-reportes-e2e/ciclo-reporte-nuevo.spec.ts` | — |
| GR-E2E-C3 | ✅ | `gestion-reportes-e2e/ciclo-reporte-nuevo.spec.ts` | — |
| GR-E2E-C4 | ✅ | `gestion-reportes-e2e/ciclo-reporte-nuevo.spec.ts` | — |
| GR-E2E-C5 | ✅ | `gestion-reportes-e2e/ciclo-reporte-nuevo.spec.ts` | H9 |

**Resumen:** 5 de 5 casos ✅ (1 test, 4 pasos y la aserción de `afterAll`).

### Corridas verificadas

Cada corrida terminó con 0 cuentas `qagr.*`, 0 perfiles `QAGR`, 0 reportes `QAGR%`, 0 intentos huérfanos y sin bloque `GR_*` en `.env.local`. En las marcadas con "censo", además, el conteo exacto de filas de las **67 tablas legibles** de `salvia` y `security` fue idéntico antes y después.

| # | Escenario | Resultado |
|---|---|---|
| 1 | Primer intento, 19:10 | ❌ **Fallo del entorno, no del spec:** el login de `qagr.sv` fue rechazado. Ninguna cuenta (ni siquiera las existentes) podía entrar: el pool `salvia_legacy` del pooler estaba lleno (`EMAXCONNSESSION`, 15 de 15; el servidor local tenía solo 2 sockets hacia el pooler, así que lo ocupaban otros clientes) y el backend lo reporta como "Usuario o contraseña incorrectos". Se liberó solo a las 19:17 (H9). Aun así terminó sin residuo: el fixture y `afterAll` borraron reporte y cuentas |
| 2 | Cuentas automáticas, 1.ª de dos seguidas | ✅ 23,8 s |
| 3 | Cuentas automáticas, 2.ª seguida (sin pasos intermedios) | ✅ 29,3 s |
| 4 | Con la verificación previa de login · censo | ✅ 23,5 s |
| 5 | **Corrida interrumpida** con `kill -9` al detectar el reporte en la BD (8 procesos): quedan 2 cuentas, el bloque `GR_*` y 1 reporte huérfano. La corrida siguiente · censo | ✅ 26,3 s — recreó las cuentas, barrió el huérfano y dejó todo en 0 |
| 6 | Con cuentas **manuales** (`usuarios-qagr.sh crear`), dos seguidas | ✅ 20,9 s y ✅ 20,6 s — las cuentas siguen ahí después, como corresponde |
| 7 | **Par final** con el código definitivo, dos invocaciones seguidas · censo tras cada una | ✅ 23,5 s y ✅ 21,0 s |

La columna *Hallazgo* indica el hallazgo no corregido que este caso dejó a la vista; el detalle está en [`gestion-reportes-qa-plan.md`](gestion-reportes-qa-plan.md#hallazgos-que-no-se-corrigieron).
