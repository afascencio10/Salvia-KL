# Gestión de Reportes — Plan de QA (compartido)

Plan de pruebas común a las tres piezas de "Gestión de reportes": el componente [`gestion-reporte-modal`](../../../Componentes/Gestion%20Reporte%20Modal/gestion-reporte-modal-index.md) y las pantallas [`Listado Reportes`](../listado-reportes-index.md) y [`Mis Reportes`](../../Mis%20Reportes/mis-reportes-index.md). Metodología: `qa-process.md` (casos → backend primero → frontend con Playwright → cobertura).

> **Estado: ejecutado el 2026-10-13.** Se corrieron todos los casos contra el backend real (`localhost:9090`) y la BD de pruebas; 155 de 155 pruebas pasan. Los bugs que aparecieron se corrigieron (tabla de abajo) y las decisiones de negocio se confirmaron con el usuario. Después se agregó un E2E repetible de un reporte nuevo, [`gestion-reportes-e2e-qa-cases.md`](gestion-reportes-e2e-qa-cases.md) (1 prueba, también en verde).

---

## Resultado de la ejecución

| Grupo | Specs | Pruebas | ✅ | ⏭️ |
|---|---|---|---|---|
| Backend (API) | `tests/gestion-reportes-api/` · 5 archivo(s) | 74 | 74 | 0 |
| Modal (UI) | `tests/gestion-reporte-modal/` · 5 archivo(s) | 42 | 42 | 0 |
| Listado Reportes (UI) | `tests/listado-reportes/` · 3 archivo(s) | 30 | 30 | 0 |
| Mis Reportes (UI) | `tests/mis-reportes/` · 1 archivo(s) | 9 | 9 | 0 |
| E2E de un reporte nuevo | `tests/gestion-reportes-e2e/` · 1 archivo(s) | 1 | 1 | 0 |
| **Total** | | **156** | **156** | **0** |

Duración acumulada de las pruebas: ≈ 15 min las 155 primeras (1 worker) y ≈ 25 s el E2E, que se agregó y se corrió aparte; tras tocar los helpers compartidos para él, `tests/mis-reportes/` volvió a pasar 9 de 9. **Casos diseñados: 176** (166 originales + 5 que surgieron al ejecutar: `M02-C12`, `M05-C20`, `M05-C21`, `LR-E01-C13`, `MR-E01-C9` + 5 del E2E: `GR-E2E-C1…C5`); el resultado de cada uno está al final de su archivo.

Además del suite se hizo una **comprobación de humo de solo lectura sobre 550 reportes reales** (los más recientes, abiertos y páginas antiguas): 550/550 detalles en 200, ninguno con más de 2,5 s, sin valores `null`/`undefined`, y el estado del listado coincide con el del detalle en todos. Los filtros y búsquedas sobre datos reales no produjeron ningún 5xx (incluidos `page=0`, `page=-3`, `page_size=0`, `page=abc`).

### Cómo repetir la ejecución

```bash
cd qa-salvia
scripts/usuarios-qagr.sh crear           # cuentas temporales qagr.sv (sv) y qagr.op (op, sin equipo) → GR_* en .env.local
# backend en :9090:  cd ../SOG_SALVIA-develop/src && PORT=9090 go run -a .
CI=1 npx playwright test tests/gestion-reportes-api tests/gestion-reporte-modal tests/listado-reportes tests/mis-reportes --retries=0
scripts/usuarios-qagr.sh eliminar        # borra las cuentas y quita GR_* de .env.local
```

Esas cuatro carpetas contienen **solo** specs de esta feature. `CI=1` fuerza headless (sin él, `playwright.config.ts` abre el navegador visible). Cada test borra lo que creó; `barrerReportesQAGR()` limpia huérfanos marcados `QAGR` al empezar los specs con pool.

El E2E de un reporte nuevo se ejecuta aparte y **no necesita nada de lo anterior**: si no hay `GR_*` crea y borra sus propias cuentas, y se puede lanzar dos veces seguidas (ver [`gestion-reportes-e2e-qa-cases.md`](gestion-reportes-e2e-qa-cases.md)):

```bash
CI=1 npx playwright test tests/gestion-reportes-e2e --retries=0
```

## Bugs encontrados y corregidos

**Cómo se detectó:** *roja* = una prueba falló y mostró el bug; *revisión* = lo vi leyendo el código/datos al preparar la prueba y se corrigió antes de la primera corrida (la prueba hoy lo vigila). Todos los arreglos están en `SOG_SALVIA-develop/src/`.

| ID | Qué pasaba | Causa raíz | Corrección (archivos) | Detección |
|---|---|---|---|---|
| BUG-01 | `GET /api/v1/victim-contacts/:icode` respondía **500** en todos los reportes | El SQL del detalle leía `vc.victim_contact_facts_description`, columna que no existe en `victim_contact` | Se quitó; la descripción sale de `form2`/`form1` (`internal/repository/victim_contacts_management_repository.go`) | Roja: 13 de 15 casos `M01-C*` |
| BUG-02 | El "Tipo" del reporte salía vacío en el listado y en el modal | El JOIN usaba `victim_contact_form2_report_type_details` (otra columna: "Usted es una persona…") y devolvía la clave del enum, no una etiqueta | JOIN por `victim_contact_form2_report_type` (V1 y V2) y etiqueta corta Víctima/Tercero/Familiar/Hombre (`internal/repository/victim_contacts_list_repository.go`, `salvia/service/victim_contacts_labels.go`, servicios de listado y gestión) | Roja: `M01-C1` ("esperaba Víctima, llegó vacío") |
| BUG-03 | Parentesco y ajuste razonable llegaban como clave (`victim_case_form2_parentesco_familia`) | El detalle devolvía el nombre del enum sin traducir | Se devuelve `categoría_código` y el servicio lo traduce con `Locale` (repositorio + `victim_contacts_labels.go`) | Roja: `M01-C4` |
| BUG-04 | Teléfono "0" en reportes Familiar/Hombre | La víctima no tiene teléfono en esos tipos (se guarda 0 por la restricción NOT NULL) y el detalle no usaba el del reportante | `NULLIF(…, 0)` y respaldo al teléfono del reportante / `form1` (repositorio de gestión) | Revisión (lo vigila `M01-C4`) |
| BUG-05 | Un mismo día de Bogotá se contaba como dos (o al revés): listado y detalle no coincidían (`diasCompletos` 1 vs 2) | `columna AT TIME ZONE 'America/Bogota'` sobre una columna `timestamp` naive que guarda UTC interpreta la hora como ya-Bogotá y desplaza el día | `AT TIME ZONE 'UTC' AT TIME ZONE 'America/Bogota'` en el listado y en `computeCounters` | Roja: `M01-C11` ("listado=1, detalle=2") |
| BUG-06 | Intentos y cierres se guardaban con la hora local del servidor | `time.Now()` en una columna `timestamp without time zone` que el sistema lee como UTC; Postgres ignora el offset, así que en un host que no está en UTC (esta máquina: EDT) quedaba desfasada. En Docker (UTC) no se nota | `time.Now().UTC()` en `InsertAttempt` y `Close` | Revisión (la vigila `M05-C1`: la fecha guardada debe ser ≈ ahora en UTC) |
| BUG-07 | **El 65 % de los reportes existentes (2 269 de 3 503, `victim_contact_status='i'`)** aparecían como "Sin asignar" y se podían gestionar | La regla de `related-tables.md` ("un legado invalidado se lee como Cerrado, solo lectura") no estaba implementada | Se leen como "Cerrado - inválido" en estado, filtros, listado y detalle (`closeReasonCode:'legado'`, motivo = su descripción); intento y cierre responden 409 | Roja: `M01-C15` |
| BUG-08 | Una observación de más de 500 caracteres daba **500** | `close_observation` es `varchar(500)` y el servicio no validaba la longitud | Validación: 400 "la observación no puede superar 500 caracteres" (`salvia/service/victim_contacts_management_service.go`, controller) | Revisión (la vigilan `M08-C3` / `M09-C3`) |
| BUG-09 | Al reabrir un reporte cerrado por no competencia / decisión, el motivo salía sin la observación | `applyDetail` no la concatenaba; al cerrar sí | `frontend/js/components/gestion-reporte-modal.js` | Revisión (la vigilan `M08-C10`, `M09-C4`) |
| BUG-10 | En el log de intentos, el agente del intento recién registrado salía vacío hasta recargar | El cliente armaba la fila con `agentName: ''` | `POST /intento` devuelve `agentName`; el modal lo usa (`victim_contacts_management_controller.go`, `gestion-reporte-modal.js`) | Revisión (la vigilan `M05-C10`, `M05-C19`) |
| BUG-11 | Header del modal vacío mientras carga el detalle; y el aviso "Solo un agente integral o supervisor puede gestionar…" salía en reportes cerrados | El header leía solo del detalle; `canManage` es `false` en cerrados y el aviso no distinguía | El header se precarga con la fila de la tabla; el aviso no se muestra en el paso `cerrado` (`gestion-reporte-modal.js` y su plantilla) | Revisión (`M01-C16`, `M01-C21`) |
| BUG-12 | La columna "Víctima / ID Reporte" no mostraba el ID; los reportes Hombre (sin nombres) quedaban en blanco | La celda `victim_info` era la de casos (nombre, documento, teléfono) | Para reportes muestra nombre (o "Sin nombre") e `ID: <radicado>` (`case_component.html`) | Revisión (`LR-E01-C6`, `LR-E01-C12`) |
| BUG-13 | **`/salvia/mis-reportes` respondía 400** con la página "Se ha producido un error inesperado" | `MyReportsGET` buscaba el template en la carpeta `my_reportes/`; el archivo está en `reportes/` | `"reportes/"` (`salvia/facades/VictimContactFacade.go`) | Roja: `LR-E06-C12` (captura de pantalla) |
| BUG-14 | En las pantallas de **casos existentes** el contador de selección decía "1 casos(s) seleccionado(s)" | El cambio `itemLabelPlural` reemplazó "caso(s)" por `${itemLabelPlural}(s)` | Singular derivado `itemLabelSingular` (`casos-component.js`, `case_component.html`). No se corrió ninguna pantalla de casos (fuera de alcance): se verificó leyendo el diff | Revisión (`LR-E05-C1`) |
| BUG-15 | El log del modal agrupaba y mostraba horas con la zona del **navegador**, no la de Bogotá: fuera de Colombia los badges `X/3` no cuadraban con "Días completos" | `toLocaleDateString('es-CO')` sin `timeZone` | `timeZone: 'America/Bogota'` (mismo criterio que `formatFechaHora` de Detalle del Caso) en `gestion-reporte-modal.js` | Roja: `M05-C20` (navegador en Tokio: 2 grupos para un único día Bogotá) |
| BUG-16 | `/salvia/primer_contacto` **en blanco para todos los roles** (previo a esta feature) → el enlace nuevo no se veía | `VictimContactGET` declara `victimCaseRes` pero nunca lo asigna y los templates hacen `JSON.parse({{.victimCases}})` → `JSON.parse("")` | `victimCaseRes = "[]"` (decisión del usuario, 2026-10-13) | Roja: `LR-E01-C11` (clic) |
| BUG-17 | `POST /cerrar` aceptaba `close_reason='ap'`: dejaba el reporte "Cerrado - válido" **sin abrir ningún caso** | `ap` estaba en la lista de motivos válidos aunque M-10 está diferido | `deferredCloseReasonCodes`: 400 mientras M-10 no exista (decisión del usuario, 2026-10-13) | Sonda de comportamiento; la vigila `M02-C12` |

## Decisiones del usuario (2026-10-13)

| Tema | Decisión | Estado |
|---|---|---|
| `POST /cerrar` con `ap` | Bloquear con 400 hasta que se construya M-10 | Hecho (BUG-17, `M02-C12`) |
| `co`/`dp` sin intento exitoso previo | Dejar como está (el backend no lo exige; probado: 200) | Sin cambios |
| `/salvia/primer_contacto` en blanco | Arreglar con `victimCaseRes = "[]"` | Hecho (BUG-16) |
| Menú lateral | Agregar "Listado Reportes" (op, sv) y "Mis Reportes" (solo op) | Hecho en `sidebar.html` (`LR-E01-C13`, `MR-E01-C9`) |
| Ruta de Mis Reportes | Dejar `/salvia/mis-reportes` (el diseño decía `/salvia/reportes/mis-reportes`) | Sin cambios |

## Hallazgos que NO se corrigieron

| # | Hallazgo | Por qué no se tocó |
|---|---|---|
| H1 | **Formulario público V2:** crear un reporte Familiar u Hombre con "¿disposición del reportante? = Sí" devuelve 400 "Este campo es obligatorio" en `reporterBestContactTime`. El wizard manda `"10:00:00"` y `VictimContactForm2V2DTO` lo declara `time.Time` sin parser de hora. Es igual en `SOG_SALVIA` y en `SOG_SALVIA-develop`. Visto por API con el mismo formato que arma el wizard; **no lo verifiqué por el wizard** | Es otra feature (PSWYA-755). Los helpers de QA usan "No" para evitarlo |
| H2 | **Riesgo latente:** si `REASSIGN_CROSS_TEAM_CONTINGENCY` pasa a `false` (`casos-component.js:43` y `reasignar-casos-modal.js:31`), no se podría marcar ningún reporte (la regla de equipo lee `caseObj.caseTeam`; el reporte trae `team`) ni abrir el modal (`resolveTeam`). Hoy el flag es `true` | Depende de una decisión del equipo sobre el flag. `flow-E05` ya recomendaba desactivar la regla para reportes |
| H3 | Al arrancar, `SyncAllSequences` registra "permission denied" para las secuencias de `victim_contact_attempt` y `victim_contact_form2_v2`: la migración da `USAGE, SELECT` y `setval` necesita `UPDATE` | Solo un warning; no afecta el funcionamiento. Mismo patrón que V2 |
| H4 | Una petición de listado abortada por el navegador (cierre de página) se registra como `context canceled` y el handler responde 500 (nadie la recibe) | Cosmético; mismo patrón en el resto de controllers |
| H5 | Una vez, una consulta de detalle se quedó 30 s y respondió 500: el keepalive de GORM había fallado con `context deadline exceeded` (estancamiento de la conexión al pooler de Supabase). El mismo reporte responde al instante después, y las 550 consultas del humo fueron rápidas | Transitorio de red/infraestructura, no reproducible |
| H6 | Reasignar un reporte ya **cerrado** se permite (probado: 200); y si el agente destino no tiene equipo, `victim_contact_team` queda `''` y no `NULL` | Sin regla de negocio que lo prohíba (Q10); sin efecto en los filtros |
| H7 | Datos reales: 2 951 de 3 503 reportes son "solo `form1`" (los registró un agente) y no tienen tipo de reporte: su columna "Tipo" muestra "—" | Es la forma de los datos, no un bug |
| H8 | **La reasignación masiva no es atómica:** `POST /reasignar-bulk` actualiza los reportes uno por uno, sin transacción. Si un `UPDATE` falla a mitad de la lista, los anteriores quedan reasignados (y con su fila de auditoría) y la respuesta es 400 con el mensaje del error, sin indicar cuántos se alcanzaron a cambiar. Todo error del servicio sale como 400, también los de BD (`victim_contacts_reassign_controller.go`). Visto leyendo el código; no se provocó una falla real a mitad de lista | Mismo patrón que la reasignación de casos. Hacerla atómica (transacción y 500 para fallas de BD) cambia el contrato de la ruta: decisión del equipo |
| H9 | **El login presenta un fallo de la BD como credenciales incorrectas:** `LoginGeneralUser` devuelve 400 "Usuario o contraseña incorrectos" también cuando no logra consultar al usuario. El 2026-10-13, de 19:10 a 19:17, ninguna cuenta (ni las ya existentes) pudo entrar: el pool `salvia_legacy` del pooler de Supabase estaba lleno (`FATAL: (EMAXCONNSESSION) max clients reached in session mode … pool_size: 15`) y el servidor local solo tenía 2 sockets abiertos hacia el pooler, así que lo ocupaban otros clientes de la misma BD (no se identificaron). Quien lo ve en pantalla cree que se equivocó de clave; la causa solo está en el log del servidor. Se liberó sola | Es el módulo de seguridad y distinguir el error (503 en vez de 400) cambia el contrato del login; no es de esta feature. Las pruebas lo mitigan con `esperarInicioSesion` (reintenta hasta 90 s y falla con el motivo) |
| H10 | **Huella en `salvia.app_log`:** cada reasignación (por UI o API) escribe una fila de auditoría por reporte (`Reasignación manual de reporte (bulk): agente actualizado`). El rol de la BD de pruebas no tiene permiso sobre esa tabla, así que los tests no pueden borrarlas: quedan las de las pruebas de reasignación y una más por cada corrida del E2E | Es el comportamiento previsto de la auditoría. Limpieza con un rol autorizado: `DELETE FROM salvia.app_log WHERE message LIKE 'Reasignación manual de reporte (bulk)%' AND data->>'actorName' = 'QAGR Supervisor QA';` |

## Documentos de diseño actualizados (2026-10-13)

Terminada la ejecución se alinearon con lo implementado; cada pieza tiene además su `changelogOct2026.md`.

| Pieza | Qué se actualizó |
|---|---|
| `Listado Reportes` | `flow-E02` reescrito (códigos `ni/nc/co/dp/ap` y filtro nombrado y combinable `filter_estado_reporte`); `flow-E01`, `E03`, `E05`, `E06`, `E07` y `E08`; `listado-reportes-interface.md`, `-events.md`, `related-tables.md` e `-index.md` (estado ✅, menú lateral, `/salvia/primer_contacto` ya carga: BUG-16, enlaces a QA) |
| `Mis Reportes` | Ruta real `/salvia/mis-reportes` en `-index`, `-events`, `-interface` y `flow-E01` (el diseño decía `/salvia/reportes/mis-reportes`) |
| `Gestion Reporte Modal` | `/intento` devuelve `agentName`; el detalle devuelve `closeReasonCode:'legado'` y `/cerrar` rechaza `ap`; `canManage` es `false` en reportes cerrados; el header se precarga con la fila; el log usa la zona de Bogotá. En `-index`, `-events`, `-interface`, `-usage`, `related-tables.md` y los flujos M01, M02, M05, M06, M07, M08 y M10 |
| `Casos component` | Props `agentId`, `listEndpoint`, `sortOptions`, `itemLabelPlural` / `itemLabelSingular` (derivado) y la celda `victim_info` para reportes; `reasignar-casos-modal` con `entityKind`. En las `-interface` y `-events` de ambos |
| `General/` | `proyecto-description.md` (3 filas nuevas) y `tables-database.md` (`victim_contact_attempt`) |

## Archivos de casos

| Pieza | Archivo | Eventos |
|---|---|---|
| Modal | [`M01-abre-modal-qa-cases.md`](../../../Componentes/Gestion%20Reporte%20Modal/QA/M01-abre-modal-qa-cases.md) | M-01 |
| Modal | [`M02-M04-validacion-previa-qa-cases.md`](../../../Componentes/Gestion%20Reporte%20Modal/QA/M02-M04-validacion-previa-qa-cases.md) | M-02, M-04 |
| Modal | [`M05-M07-esquema-3x3-qa-cases.md`](../../../Componentes/Gestion%20Reporte%20Modal/QA/M05-M07-esquema-3x3-qa-cases.md) | M-05, M-06, M-07 |
| Modal | [`M08-M09-cierre-post-contacto-qa-cases.md`](../../../Componentes/Gestion%20Reporte%20Modal/QA/M08-M09-cierre-post-contacto-qa-cases.md) | M-08, M-09 |
| Modal | [`M11-cierra-modal-qa-cases.md`](../../../Componentes/Gestion%20Reporte%20Modal/QA/M11-cierra-modal-qa-cases.md) | M-11 |
| Listado | [`E01-carga-pantalla-qa-cases.md`](E01-carga-pantalla-qa-cases.md) | E-01 |
| Listado | [`E02-E03-E07-E08-filtros-busqueda-orden-qa-cases.md`](E02-E03-E07-E08-filtros-busqueda-orden-qa-cases.md) | E-02, E-03, E-07, E-08 |
| Listado | [`E04-E05-E06-ver-y-reasignar-qa-cases.md`](E04-E05-E06-ver-y-reasignar-qa-cases.md) | E-04, E-05, E-06 |
| Mis Reportes | [`mis-reportes-qa-cases.md`](../../Mis%20Reportes/QA/mis-reportes-qa-cases.md) | E-01, E-02 |
| E2E | [`gestion-reportes-e2e-qa-cases.md`](gestion-reportes-e2e-qa-cases.md) | Un reporte nuevo: formulario público → Listado → Mis Reportes → modal |

**Convenciones.** Casos: `M01-C1` (modal), `LR-E02-C1` (Listado Reportes), `MR-E01-C1` (Mis Reportes), `GR-E2E-C1` (E2E de un reporte nuevo). Flujos de prueba: `TB-…` (backend, vía API con la sesión de Playwright) y `TU-…` (UI). Las dudas abiertas se citan como `(Q#)` — ver la última sección.

## Fuera de alcance

- **M-03** (removido) y **M-10** (diferido, junto con GAP-6/GAP-7): no hay casos. El botón "Apertura de caso" y "Reporte de prueba" **no deben existir** — se verifica su ausencia en M01-C17 y M08-C5.
- **Regresión** de los consumidores existentes de `casos-component` / `reasignar-casos-modal` (`list_cases.html`, `get_my_cases.html`, `case_component.html`). Los cambios son aditivos, pero el workflow indica no correr regresión salvo que el usuario la pida y defina el alcance.
- Lógica legada de `/salvia/primer_contacto` (solo se prueba el enlace nuevo hacia `/salvia/reportes`, `LR-E01-C11`).
- Texto de tooltips (los MDs no lo definen), colores/estilos exactos, rendimiento.

## Orden de ejecución (seguido en esta ronda)

**Paso 0 — Pre-chequeos (una vez)**

1. Backend corriendo en `http://localhost:9090` con el código nuevo — `playwright.config.ts` no levanta servidor.
2. Migración `add_victim_contact_management.sql` aplicada: columnas `victim_contact_agent_id/team/close_reason/close_observation/closed_date`, tabla `victim_contact_attempt` y enums `ni/nc/co/dp/ap` en `victim_case_form2_enums` (categoría `victim_contact_close_reason`).
3. **Convención de `victim_contact_attempt_date`** (`timestamp without time zone`): UTC naive, igual que el resto del proyecto, que convierte a `America/Bogota` solo al mostrar. Lo vigila `M05-C1` (la fecha guardada ≈ `now()` en UTC). Ojo: el endpoint guardaba la hora **local** del servidor (BUG-06); se corrigió antes de la primera corrida, así que no hay una prueba roja que lo muestre.
4. Cuentas de sesión: `scripts/usuarios-qagr.sh crear` (ver "Usuarios"). Los logins `CASE_ENTITIES_*` ya no existen en la BD.

**Paso 1 — Backend, endpoint por endpoint** (`qa-process.md` Fase 2)

Orden: `GET detalle` → `POST intento` → `POST cerrar` → `GET list` → `POST reasignar-bulk`. Por cada uno: agregar logs estratégicos en el controller/service (inicio de función, parámetros, resultado), correr sus flujos `TB-*`, revisar logs y verificar la BD. No se pasa a UI hasta cerrar el backend.

**Paso 2 — Frontend con Playwright** (`qa-process.md` Fase 3)

Los selectores salieron de los templates reales (`gestion_reporte_modal.html`, `case_component.html`, `reasignar_casos_modal.html`) y se formalizaron directo en `.spec.ts`; no se usó `playwright-cli` interactivo. Las capturas de los fallos (`test-results/`) sirvieron para diagnosticar BUG-13. Los `TU-*` no repiten validaciones ya cubiertas por la API: verifican render, transiciones, requests emitidas y persistencia visible. Orden: modal (apertura → validación previa → 3x3 → post-contacto → cierre) → Listado → Mis Reportes.

> Los `TU-*` de Listado y Mis Reportes aíslan sus fixtures buscando el marcador `QAGR` en la caja de búsqueda, así que `TB-LR-E07-1` (búsqueda por API) se ejecuta **antes** de cualquier `TU-LR-*`/`TU-MR-*`: si la búsqueda falla, el fallo se atribuye a ella y no a un flujo de UI.

**Paso 3 — Documentar**: resultados de cada caso en su MD (como en `E10-cuando-envia-reporte-qa-cases.md`), hallazgos y bugs.

> Solo se corren los specs de estas tres piezas, archivo por archivo. Sin barridos de carpetas ni regresión.

## Endpoints bajo prueba

| Endpoint | Evento(s) | Casos |
|---|---|---|
| `GET /api/v1/victim-contacts/:icode` | M-01 | `M01-C1…C15` |
| `POST /api/v1/victim-contacts/:icode/intento` | M-05, M-06 | `M05-C1…C9`, `M05-C21`, `M06-C1…C2` |
| `POST /api/v1/victim-contacts/:icode/cerrar` | M-02, M-07, M-08, M-09 | `M02-C1…C6`, `M02-C12`, `M07-C1…C4`, `M08-C1…C4`, `M09-C1…C3` |
| `GET /api/v1/victim-contacts/list` | E-01, E-02, E-03, E-07, E-08 (y Mis Reportes) | `LR-E01-C4/C5/C10`, `LR-E02-C1…C6`, `LR-E03-C1…C3`, `LR-E07-C1…C7`, `LR-E08-C1…C2`, `MR-E01-C7` |
| `POST /api/v1/victim-contacts/reasignar-bulk` | E-06 | `LR-E06-C1…C8` |
| `GET /salvia/reportes`, `GET /salvia/mis-reportes`, menú lateral | E-01 | `LR-E01-C1…C3`, `LR-E01-C13`, `MR-E01-C1…C2`, `MR-E01-C9` |
| `GET /api/v1/agents/by-role?role=op` (existente; lo usa el autocomplete de asignación) y `GET /api/v1/equipo-operadores?role=op` (candidatos del modal de reasignación) | E-03, E-06 | `LR-E03-C5`, `LR-E06-C11` |
| `POST /salvia/public-v2` (existente; lo usa el wizard público) | Alta del reporte del E2E | `GR-E2E-C1` |

Los nombres de campo del JSON se toman del modelo en implementación, `src/internal/models/report_list_item.go` (`ReportListItem`, `ReportDetail`, forma **plana**). Se valida por contenido, no por anidamiento (Q14).

## Usuarios

| Alias | Rol | Cuenta | Para qué |
|---|---|---|---|
| `SV` | `sv` | `GR_SV_LOGIN` (**`qagr.sv`**, temporal) | Reasignar; gestionar cualquier reporte (GAP-2) |
| `OP1` | `op` | `GR_OP_LOGIN` (**`qagr.op`**, temporal, **sin equipo**) | Dueño de "Mis Reportes"; gestión de reportes propios y ajenos |
| `OP2` | `op` (otro agente activo) | **Sin login:** se lee de la BD un `op` con `status='e'` distinto de `qagr.*` (`getOpExistente`) | Dueño "ajeno" de reportes; destino alternativo de reasignaciones |
| `RO` | `ro` | `TEST_USER` / `TEST_PASSWORD` (`yiceneth.zamudio`, la cuenta de QA del repo) | Rechazo por rol |

Las cuentas `qagr.*` las crea y las borra `qa-salvia/scripts/usuarios-qagr.sh` (a mano para las demás suites; el E2E de un reporte nuevo las crea y las retira solo, salvo que ya existan cuentas `GR_*` manuales, que respeta). **Se crearon porque los logins `CASE_ENTITIES_*` de `.env.local` ya no existen en la BD y no había contraseña de ningún `sv`/`op`** (el usuario aprobó crearlas el 2026-10-13). No se usa `auth.fixture` (cada test elige su rol): `helpers/sesion.ts` hace login por API (con el header `User-Agent: flutter-client`, que salta el captcha como la app móvil) y por pantalla (`LoginPage`, navegador en `America/Bogota`, locale `es-CO`).

**Limitación a respetar:** iniciar sesión con un login que ya tiene sesión activa **cierra la anterior** (`GeneralUserLOGIN_POST`). Por eso ningún test combina `api*` y `ui*` del mismo rol.

## Datos de prueba

La BD de pruebas es compartida con datos reales: **todo lo creado se marca, se acota por el marcador y se borra**.

### A. Un reporte nuevo por test (modal y backend)

`crearReporteQA(tipo)` crea el reporte por el `POST /salvia/public-v2` real (header `User-Agent: flutter-client`, que salta el captcha, igual que `crear-reporte-victima.spec.ts`) y devuelve el radicado (`icode`). Tipos: `victima`, `tercero-propio`, `tercero-anonimo`, `familiar`, `hombre`. El payload sale de la request real del wizard (capturable con `playwright-cli`); los 4 tipos ya están validados en `E10-cuando-envia-reporte-qa-cases.md` (C3–C7). El reporte se borra en `afterEach` con `limpiarReporteQA(icode)`.

Cerrar un reporte es irreversible por UI, por eso cada test que lo cierra o le registra intentos usa su propio reporte.

### B. Pool compartido para Listado y Mis Reportes (`beforeAll`)

Nombres marcados en BD con `QAGR-<runId>-R#` (así el marcador aísla los fixtures de los datos reales y no depende del tipo de reporte). Fechas de reporte fijas para que el orden sea determinista: `R# → hoy − # días, 10:00 Bogotá`. Los intentos se siembran siempre **después** de la fecha de creación del reporte.

Cada spec que usa el pool lo crea en su propio `beforeAll` con un `runId` distinto y lo borra en `afterAll`. Los tests que modifican asignaciones o cierres (reasignar, `MR-E01-C8`, `MR-E02-*` y todo el modal) usan reportes propios por test, nunca el pool: así ningún conjunto esperado de la tabla de abajo se contamina.

| Alias | Tipo | Estado esperado | Agente | Siembra | Cierre |
|---|---|---|---|---|---|
| R1 | Víctima (V2) | Sin asignar | — | — | — |
| R2 | Tercero | Asignado | `OP1` | fallidos: 3 en D−2 y 2 en D−1 → "5/9 intentos · 1/3 días" | — |
| R3 | Familiar | Asignado | `OP2` | — | — |
| R4 | Víctima | Cerrado - inválido | `OP1` | — | `ni` |
| R5 | Tercero | Cerrado - inválido | `OP1` | 3 fallidos en cada uno de D−4, D−3 y D−2 → "9/9 intentos · 3/3 días" | `nc` |
| R6 | Víctima | Cerrado - válido | `OP1` | 1 exitoso en D−5 | `ap` (por BD: ningún evento lo escribe en esta fase, Q2) |
| R7 | Víctima | Legado inválido (`victim_contact_status='i'`, sin `close_reason`) | — | — | — (Q9) |
| R8 | Tercero | Cerrado - inválido, **sin agente** | — | — | `co` + observación |
| R9 | V1 con `victim_contact_form1` (teléfono + documento) | Sin asignar | — | — | — |
| R10 | Víctima | Cerrado - inválido | `OP1` | — | `dp` + observación |

Teléfonos conocidos para las pruebas de búsqueda: R1 (Form2) y R9 (Form1 + documento, único reporte que se encuentra por documento).

Conjuntos esperados (se asierta pertenencia / no pertenencia de estos alias, **no** igualdad exacta, para tolerar datos reales):

| Filtro | Incluye | No incluye |
|---|---|---|
| `estado = asignado` | R2, R3 | R1, R4, R5, R6, R8, R9, R10 |
| `estado = sin_asignar` | R1, R9 | R2, R3, R4, R5, R6, R8, R10 |
| `estado = cerrado_invalido` | R4, R5, R8, R10 | R6, R1, R2, R3, R9 |
| `estado = cerrado_valido` | R6 | todos los demás |
| `asignacion = OP1` | R2, R4, R5, R6, R10 | R3, R1, R8, R9 |
| `asignacion = __unassigned__` | R1, R8, R9 | R2, R3, R4, R5, R6, R10 |
| Mis Reportes de `OP1` | R2, R4, R5, R6, R10 | R3, R1, R8, R9 |

R7 queda fuera de estos conjuntos: su clasificación depende de la regla de lectura de reportes legados (Q9) y tiene su propio caso (`LR-E01-C8`).

### C. Fechas y zona horaria

- Los intentos se siembran por BD con `sembrarIntentos(icode, [{ resultado, utc, agenteIcode }])` — la API siempre usa `NOW()`, así que el esquema 3x3 (3 días × 3 fallidos) **solo es alcanzable sembrando días pasados**.
- Las fechas se construyen con `fechaBogota(offsetDias, 'HH:MM')`, que devuelve el instante UTC de esa hora local de Bogotá (UTC−5). Siempre relativas a hoy → no dependen de cuándo se corra el test.
- **No correr** los casos de frontera de "hoy" (`M01-C12`) entre las 00:00 y las 00:10 hora Bogotá.
- Fronteras que se prueban a propósito (clase de bug ya vista en `intento-hora-zona-horaria.spec.ts`): 23:50 Bogotá y 00:10 Bogotá del día siguiente son el **mismo día UTC** pero **dos días Bogotá**; 02:15 UTC es 21:15 Bogotá del día anterior.

### D. Limpieza

`afterEach`/`afterAll`: `limpiarReporteQA(icode)` borra primero `victim_contact_attempt` del reporte y luego aplica el cascade de `resetVictimContactV2Fixture` (que hoy **no** conoce la tabla de intentos). Para reutilizar un reporte entre casos del mismo archivo: `resetGestionReporte(icode)` (borra intentos y deja `agent_id/team/close_*` en `NULL`).

### E. Un reporte por el wizard (E2E)

El E2E de un reporte nuevo (`GR-E2E-*`) es el único que crea el reporte **por la pantalla** del wizard público y no con `crearReporteQA`. Nombres `QAGRE2E`, apellido `Ciclo-<runId>` y teléfono aleatorio; la descripción empieza con `QAGRE2E`. Lo borra el fixture `reportes` (que lo registra en cuanto el wizard devuelve el radicado) y, como red de seguridad, `barrerReportesQAGR('QAGRE2E')` antes y después: el prefijo acota la barrida a lo suyo y no toca los pools de las otras suites.

## Helpers y page objects creados en `qa-salvia/`

| Archivo | Contenido |
|---|---|
| `helpers/reportes.ts` | SQL por `psql` sin shell (`sql`, `lit`); `crearReporteQA(tipo)` por el POST real de `/salvia/public-v2`; siembra (`sembrarIntentos`, `fallidosEnDia`, `asignarReporte`, `cerrarReporteDirecto`, `marcarReporteLegado`, `agregarForm1`, `quitarTelefonos`, `marcarReporte`); lectura (`getReporteEstado`, `getIntentosBD`, `contarIntentosBD`, `getUserICode`, `getOpExistente`, `getOpConEquipo`, `getOpInactivo`); limpieza (`limpiarReporteQA`, `barrerReportesQAGR(prefijo = 'QAGR')`, `limpiarPool`); `crearPoolListado` |
| `helpers/reportes-api.ts` | Atajos HTTP: `getDetalle`, `postIntento`, `postCerrar`, `postReasignar`, `listar`, `listarTodo` |
| `helpers/reportes-ui.ts` | `abrirModal`, `abrirModalEnTabla`, `textoDeBusqueda`, `contarRequests`, `aliasesEnTabla` |
| `helpers/sesion.ts` | Credenciales por rol; login por API (`apiLogin`) y por pantalla (`uiLogin`, con `timezoneId` opcional); `esperarInicioSesion` (verifica el login y reintenta hasta 90 s, ver H9) |
| `helpers/cuentas-qagr.ts` | Alta y baja de `qagr.sv` / `qagr.op` desde un spec (`prepararCuentasQA`, `retirarCuentasQA`, `hayCuentasAutogeneradas`), envolviendo `scripts/usuarios-qagr.sh` |
| `fixtures/gestion-reportes.fixture.ts` | Fixtures `apiSv/apiOp/apiRo/apiAnon`, `uiSv/uiOp/uiRo` y `reportes` (gestor que borra lo que creó el test) |
| `utils/fechas.ts` | `fechaBogota(offsetDias, 'HH:MM')`, `diaBogota`, `horaBogotaAhora` |
| `pages/gestion-reporte-modal.page.ts` | Locators del modal (`.grm-*`) |
| `pages/listado-reportes.page.ts` | `TablaReportesPage` (base de `ListadoReportesPage` y `MisReportesPage`) y `ReasignarModalPage` |
| `scripts/usuarios-qagr.sh` | Crea / elimina las cuentas `qagr.sv` y `qagr.op`; con `QAGR_AUTOGEN=1` marca el bloque `GR_*` como de una corrida automática |

Specs:

```
tests/gestion-reportes-api/   detalle · intentos · cerrar · listado · reasignar
tests/gestion-reporte-modal/  apertura · validacion-previa · esquema-3x3 · post-contacto · cierre-modal
tests/listado-reportes/       carga · filtros-busqueda-orden · ver-y-reasignar
tests/mis-reportes/           mis-reportes
tests/gestion-reportes-e2e/   ciclo-reporte-nuevo
```

## Dudas (Q#) — resolución tras la ejecución

| # | Duda | Resolución |
|---|---|---|
| Q1 | `flow-E02` desactualizado (códigos y nombre del filtro; su PASO 1 reemplazaría el scope) | Se implementó como filtro **nombrado y combinable** `filter_estado_reporte` (migración + `listado-reportes-interface.md`); lo confirman `LR-E02-C6` y `MR-E01-C5` (el scope de Mis Reportes sobrevive al filtro). **`flow-E02` se reescribió** con esos nombres (2026-10-13) |
| Q2 | `cerrado_valido` inalcanzable desde la UI | Se siembra por BD (R6) y solo se valida filtro y badge. Desde BUG-17 tampoco se llega por API (`ap` → 400) |
| Q3 | Qué muestra el modal ante 404/409/500 | Muestra el texto de error del servidor (`.grm-action-error` en acciones, `.grm-error` en la carga), no cambia de paso, reactiva los botones y se puede cerrar. Cumple el mínimo exigido |
| Q4 | Validación de `result` en `/intento` | 400 si falta, está vacío o es desconocido (`M05-C8`) |
| Q5 | `ap` y precondición del exitoso en `/cerrar` | **Decidido:** `ap` → 400 (`M02-C12`); `co`/`dp` **no** exigen exitoso previo (probado: 200) |
| Q6 | Longitud de `close_observation` | 500 → 200; 501 → 400 (BUG-08) |
| Q7 | ¿Gestionar asigna el reporte? | No: `intento` y `cerrar` no tocan `victim_contact_agent_id` ni el equipo; `M05-C9` verifica que el agente del intento es quien lo registró |
| Q8 | 401 o 403 por rol | **401** siempre (el middleware `CheckPermission` escribe 401); sin sesión también 401. No se observó ningún 403. Las pruebas aceptan 401/403 |
| Q9 | Reportes legados (`status='i'`) | Implementado según `related-tables.md` (BUG-07) |
| Q10 | Reasignación: rol, cerrados, etiqueta | Solo `sv` (op/ro → 401, `LR-E06-C5`); etiqueta correcta "Reasignar Reportes" (`LR-E05-C2`); reasignar un **cerrado** se permite (H6) |
| Q11 | Familiar/Hombre sin nombres | La tabla muestra "Sin nombre" + `ID: <radicado>` (BUG-12); el título del modal queda en "Reporte" sobre el radicado; el teléfono es el del reportante (BUG-04); ningún `0`/`null`/`undefined` visible |
| Q12 | `agents/search` sin filtro de rol | El autocomplete usa `GET /api/v1/agents/by-role?role=op` con caché en el cliente; solo agentes `op` (`LR-E03-C5`) |
| Q13 | ¿Rechazar fallido sobre un reporte con exitoso sin cerrar? | No se probó |
| Q14 | Forma del JSON de detalle | Plana, como el modelo en código (`ReportDetail`) |
