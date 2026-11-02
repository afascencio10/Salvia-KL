# QA — Registro V3 · E-12 Cuando guarda una sección (WP-2)

Casos de prueba de [`flow-E12`](../Flujos/flow-E12-cuando-guarda-una-seccion.md) para lo que cambió en WP-2 (2026-11-02): el guardado atómico y validado (**REG-11**), la hora del hecho (**REG-09**), la pregunta «Tipo de discapacidad» (**REG-08**) y la zona de residencia (**REG-07**). Endpoints: `POST /api/v1/victim-cases-v3/:icode/seccion/:modulo` y `GET /api/v1/victim-cases-v3/:icode` (y `GET /victim-cases-v3/enums` para el catálogo). Lo que se ve en pantalla —avisos, lista desplegable, control de hora— está en [`E11-actualizar-respuesta-y-tamizaje-qa-cases.md`](E11-actualizar-respuesta-y-tamizaje-qa-cases.md) y [`E21-guardado-o-carga-fallida-qa-cases.md`](E21-guardado-o-carga-fallida-qa-cases.md).

> **Estado: ejecutado el 2026-11-02** — resultados al final de este archivo.

**Método.** Backend primero (qa-process.md): se llama el endpoint con `page.evaluate(fetch)` desde una sesión `op` y se verifica **directo en la BD** con `psql`; se agregó un log temporal en el repositorio (`[WP2-QA]`, retirado antes del commit) y los logs permanentes del controlador (`[RegistroV3] SaveSection rechazado|falló`) se comparan con los 400/500 provocados. Cada script crea sus Casos de prueba por la API (documento `9xxxxxxxx`) y los borra al terminar. Los scripts se corrieron antes del arreglo (para confirmar que detectan el defecto) y después. Un mismo usuario no puede tener dos sesiones (el segundo login invalida la primera), por eso cada script usa una sola.

---

## Casos de prueba

### REG-11 — Guardado atómico y validado (`wp2-reg11-guardado-seccion.js`)

| Caso | Descripción | Condición inicial | Resultado esperado | Método de verificación |
|---|---|---|---|---|
| E12-C1 | Hechos completos se guardan y vuelven en el formato de la pantalla | Borrador nuevo; HV con fecha `2026-10-30`, hora `19:30`, departamento «Atlántico», municipio «Barranquilla» (12 letras; antes el límite era 8) | 200; BD: fecha `2026-10-30`, hora `19:30:00`, departamento y municipio guardados; `GET` devuelve `HV_2 = 2026-10-30`, `HV_3 = 19:30`, `HV_4` y `HV_5` (antes: fecha cruda `2026-10-30 00:00:00 +0000 UTC` y `HV_4` perdido) | API + BD |
| E12-C2 | Residencia: departamento y municipio largo | DP con `DP_27` «Atlántico», `DP_28` «Barranquilla», `DP_1` `1990-05-20` | 200; BD `residence_department`, `residence_town` y `victim_case_victim_town_code`; `GET` devuelve `DP_1` en `AAAA-MM-DD` | API + BD |
| E12-C3 | Datos inválidos → 400 con el campo y el motivo | Un Caso nuevo por entrada: `HV_3 = 7pm`, `HV_2 = 30/10/2026`, `HV_2 = 2026-02-30`, `HV_7` de 151 caracteres, `HV_6` con un icode inexistente, `CA_2 = abc`, `CI_6` de 33 caracteres, `CI_2` con un primer nombre de 33 caracteres | 400 `{error: "valor no válido", campo, detalle}` con el `campo` enviado y un `detalle` en español (antes: 500 genérico, y `CA_2 = abc` se guardaba como 0 con 200) | API |
| E12-C4 | Lo que se puede interpretar se acepta | `CA_2 = "300 123 4567"`; `HV_2 = "2026-10-31 00:00:00 +0000 UTC"` (formato viejo); `HV_3 = "08:05:30"` | 200; BD teléfono `3001234567`, fecha `2026-10-31`, hora `08:05:30`; `GET` devuelve `HV_3 = 08:05` (antes: el teléfono quedaba en `300` y la fecha vieja daba 500) | API + BD |
| E12-C5 | Atomicidad: un dato inválido junto a uno válido no guarda nada | Borrador nuevo; `{HV_1: "TEXTO NUEVO…", HV_3: "7pm"}` | 400; la descripción de los hechos no cambia; `victim_case_v3_answered_keys` no cambia; `HV_1` y `HV_3` no quedan respondidas (antes: 500 y las claves quedaban marcadas) | API + BD |
| E12-C6 | El Caso no pasa a Iniciado con los hechos sin guardar | CI y DP guardados; HV con `HV_3 = 7pm`; luego AG; luego HV corregido | HV → 400; AG → 200 con `triadaCompleta = false`; estado sigue `bo` y ninguna clave HV respondida; al corregir HV → 200, `triadaCompleta = true`, estado `in` (antes: AG daba `triadaCompleta = true` y el Caso pasaba a `in` con 12 claves HV marcadas y los hechos en blanco) | API + BD |
| E12-C7 | Valores vacíos ya no dan 500 ni borran lo guardado | Caso con hechos guardados; `{HV_2: "", HV_3: "", HV_6: ""}` | 200; fecha, hora y zona conservan su valor (antes: 500) | API + BD |
| E12-C8 | Sin respuesta ≠ valor marcador | Borrador sin responder | `GET` no devuelve `HV_2` ni `DP_1` (la fecha marcador `1900-01-01` no se expone; antes salía como respuesta) | API |
| E12-C9 | Departamento vaciado → NULL | Caso con `HV_4` guardado; `{HV_4: ""}` | 200; BD `NULL`; `GET` no devuelve `HV_4` | API + BD |
| E12-C10 | El fan-out del Perfil sigue funcionando tras mover el guardado a una transacción | Dos Casos con el mismo tipo y número de documento; en el primero se guarda `DP_29` y `DP_27` | 200; el segundo Caso recibe la dirección y el departamento | API + BD |
| E12-C11 | Logs del servidor | Log de la app durante la corrida | Una línea `[RegistroV3] SaveSection rechazado …` por cada 400 provocado (10) y ninguna `SaveSection falló` (500) | Logs |

### REG-09 — Hora del hecho, bordes en el API (`wp2-reg09-hora-hecho.js`)

| Caso | Descripción | Condición inicial | Resultado esperado | Método de verificación |
|---|---|---|---|---|
| E12-C12 | Medianoche y 23:59 se guardan como respuestas reales | Caso con hechos; `HV_3 = 00:00` y luego `23:59` | 200; BD `00:00:00` y `23:59:00`; `GET` devuelve `00:00` y `23:59` (la medianoche no se confunde con el valor marcador) | API + BD |
| E12-C13 | Hora fuera de rango | `HV_3 = 25:00` | 400 con `campo = HV_3`; la hora guardada no cambia | API + BD |

### REG-08 — Tipo de discapacidad `DP_51` (`wp2-reg08-tipo-discapacidad.js`)

| Caso | Descripción | Condición inicial | Resultado esperado | Método de verificación |
|---|---|---|---|---|
| E12-C14 | El catálogo trae las 7 opciones con su texto | `GET /victim-cases-v3/enums` | `victim_case_form2_disability_type` = 01 Discapacidad física · 02 auditiva · 03 visual · 04 Sordoceguera · 05 intelectual · 06 psicosocial (mental) · 07 múltiple, en ese orden | API |
| E12-C15 | La migración es idempotente | Reaplicar `20261102_wp2_registro_v3_catalogos_y_columnas.sql` | Sin error; sigue habiendo 7 filas y una sola columna `victim_case_form2_disability_type` | psql |
| E12-C16 | Un borrador nuevo no trae valor | Borrador sin responder | Columna `NULL` (no hay valor marcador); `GET` no devuelve `DP_51` | API + BD |
| E12-C17 | Guardar y cambiar la opción | `DP_2 = Sí`, `DP_51 = Discapacidad física`; luego `Discapacidad múltiple` | 200; BD con el id de la opción 01 y luego 07; `GET` devuelve el icode y `DP_51` queda en `answeredKeys` | API + BD |
| E12-C18 | Limpiar | `DP_2 = No`, `DP_51 = ""` | 200; columna `NULL`; `GET` no devuelve `DP_51` | API + BD |
| E12-C19 | Opción inexistente | `DP_51 = "no-existe"` | 400 con `campo = DP_51`; la columna no cambia | API + BD |
| E12-C20 | El Perfil se replica | Dos Casos con el mismo documento; `DP_51` en el primero | El segundo Caso recibe la opción (DP es Perfil) | API + BD |

### REG-07 — Zona de residencia `DP_30` (`wp2-reg07-zona-residencia.js`)

| Caso | Descripción | Condición inicial | Resultado esperado | Método de verificación |
|---|---|---|---|---|
| E12-C21 | El catálogo trae Rural y Urbano | `GET /victim-cases-v3/enums` | `victim_case_form2_residence_zone` = `ru` Rural · `ur` Urbano, con icode de 36 caracteres (antes: 0 opciones) | API |
| E12-C22 | La zona del hecho no cambió | Mismo catálogo | `victim_case_form2_facts_zone` sigue con 3 opciones (cm, cp, pr) con su texto | API |
| E12-C23 | La migración es idempotente | Reaplicar el SQL | Sin error; siguen 2 filas | psql |
| E12-C24 | Guardar y cambiar la zona | `DP_30 = Rural` y luego `Urbano` | 200; BD con el id de cada opción; `GET` devuelve el icode y la clave queda respondida | API + BD |

---

## Cobertura

| Flujo de prueba | Casos cubiertos | Script |
|---|---|---|
| `wp2-reg11` — comprobaciones C1, C2 | E12-C1, E12-C2 | `cloudDevelopment/e2e/wp2-reg11-guardado-seccion.js` |
| `wp2-reg11` — C3 (8 entradas inválidas) | E12-C3 | ídem |
| `wp2-reg11` — C4 | E12-C4 | ídem |
| `wp2-reg11` — C5, C6a, C6b | E12-C5, E12-C6 | ídem |
| `wp2-reg11` — C7, C8, C9 | E12-C7, E12-C8, E12-C9 | ídem |
| `wp2-reg11` — C10, C11 | E12-C10, E12-C11 | ídem |
| `wp2-reg09` — C1, C2, C3 | E12-C12, E12-C13 | `cloudDevelopment/e2e/wp2-reg09-hora-hecho.js` |
| `wp2-reg08` — C1…C8 | E12-C14 … E12-C20 | `cloudDevelopment/e2e/wp2-reg08-tipo-discapacidad.js` |
| `wp2-reg07` — C1…C5 | E12-C21 … E12-C24 | `cloudDevelopment/e2e/wp2-reg07-zona-residencia.js` |
| Unitarias Go — `victim_case_v3_validation_test.go` | Validadores de fecha, hora y número; traducción de errores `22xxx` (apoyan E12-C3 y E12-C4) | `src/internal/repository/` (`go test ./internal/repository/`) |

## Resultados de ejecución

**Fecha:** 2026-11-02 · **Build:** rama `claude/busy-edison-uk6g2a` en `http://127.0.0.1:9090` (sandbox cloud) · **BD:** Postgres local con la migración WP-2 aplicada · **Herramienta:** Playwright headless + `psql` + `go test`. Backend primero; la parte de pantalla de estos scripts está en E11 y E21.

| Casos | Resultado | Antes del arreglo |
|---|---|---|
| E12-C1 … E12-C11 (REG-11) | ✅ 11 de 11 | Fallaban C1-C7, C8 y C10; solo pasaban C9 y C11 (C11 con 0 líneas) |
| E12-C12, E12-C13 (REG-09) | ✅ 2 de 2 | Ya pasaban con el backend de REG-11 (el defecto de REG-09 es del control, ver E11) |
| E12-C14 … E12-C20 (REG-08) | ✅ 7 de 7 | La pregunta no existía |
| E12-C21 … E12-C24 (REG-07) | ✅ 4 de 4 | C21 y C23 fallaban (0 opciones), C22 pasaba y el resto no podía correr |

**Resumen:** 24 de 24 casos ✅. Hallazgos que quedaron fuera de alcance: [`registro-v3-data.md`](../registro-v3-data.md#estado-de-los-gaps-2026-10-31) G-V3-27 (opcionales NOT NULL que no se pueden vaciar) y G-V3-28 (clave sin traducir en CI_5).
