# QA — Registro V3 · E-11 Cuando actualiza una respuesta del formulario (WP-2)

Casos de prueba de [`flow-E11`](../Flujos/flow-E11-cuando-actualiza-una-respuesta.md) para lo que cambió en WP-2 (2026-11-02): etiquetas de los subtipos (**REG-10**), lo que se ve de la zona de residencia (**REG-07**), la hora del hecho (**REG-09**) y la pregunta «Tipo de discapacidad» (**REG-08**), y el Tamizaje en tablas (**REG-12**) y por etapas (**REG-34**). Los casos de Tamizaje anteriores (gating, ramas, umbrales) siguen en [`tamizaje-qa-cases.md`](tamizaje-qa-cases.md) y no se repiten aquí. El guardado en el API de cada requerimiento está en [`E12-guardar-seccion-qa-cases.md`](E12-guardar-seccion-qa-cases.md); los avisos de error, en [`E21-guardado-o-carga-fallida-qa-cases.md`](E21-guardado-o-carga-fallida-qa-cases.md).

> **Estado: ejecutado el 2026-11-02** — resultados al final de este archivo.

**Método.** Backend primero cuando lo hay (catálogo `GET /victim-cases-v3/enums`, `POST …/tamizaje`, `GET …/:icode`, verificado en la BD con `psql`); después Playwright solo para lo que el backend no cubre. Antes de escribir los scripts se exploró la pantalla con `playwright-cli` (autorización → tipo de contacto → «+ Nuevo registro» → módulo) para fijar los selectores. Las preguntas se ubican por el texto exacto de su etiqueta (`div.mb-3` con `label.form-label`); las filas del Tamizaje, por el texto de la pregunta. Cada script se corrió contra el código anterior para comprobar que detecta el defecto.

---

## Casos de prueba

### REG-10 — Subtipos de violencia con su texto (`wp2-reg10-subtipos-violencia.js`)

| Caso | Descripción | Condición inicial | Resultado esperado | Método de verificación |
|---|---|---|---|---|
| E11-C1 | El catálogo de subtipos está completo y es inequívoco | `GET /victim-cases-v3/enums` | Las 7 categorías `victim_case_form2_subtype_violence_experienced_*` tienen opciones con texto en español (ninguna con la clave cruda) y los códigos de cada una son únicos; el catálogo de tipos (HV_9) trae los 7 códigos | API |
| E11-C2 | Un solo tipo: Violencia psicológica | Borrador abierto en Hechos de violencia; se marca solo «Psicológica» | HV_10 muestra los 6 subtipos con el texto del catálogo, entre ellos «Acoso laboral» y «Aislamiento» (antes salían «Alojamiento» y «Articulación interinstitucional») | DOM |
| E11-C3 | Los 7 tipos a la vez | Se marcan los 7 tipos | 31 subtipos, cada uno con su texto (antes 6 con el texto de otra categoría: «Trata de personas…» → «Transporte», «Esterilización…» → «Estabilización», «Amenazas y hostigamientos contra hijas/os» → «Atención Hombres», «Privación de libertad…» → «Protección») | DOM |

### REG-07 — Zona de residencia en pantalla (`wp2-reg07-zona-residencia.js`)

| Caso | Descripción | Condición inicial | Resultado esperado | Método de verificación |
|---|---|---|---|---|
| E11-C4 | La lista desplegable tiene las dos opciones | Módulo Datos personales de un borrador nuevo | `DP_30` ofrece «Seleccione…», «Rural» y «Urbano» (antes solo «Seleccione…») | DOM |
| E11-C5 | La selección se conserva | Se elige «Rural», se pasa a otro módulo y se vuelve | El valor del campo es el icode de Rural; sigue seleccionado | DOM |

### REG-09 — Hora del hecho como selector (`wp2-reg09-hora-hecho.js`)

| Caso | Descripción | Condición inicial | Resultado esperado | Método de verificación |
|---|---|---|---|---|
| E11-C6 | El control es un selector de hora | Caso con hechos guardados (hora 19:30), módulo Hechos de violencia | `HV_3` es `<input type="time" step="60">` sin placeholder, con valor `19:30` | DOM |
| E11-C7 | Elegir una hora y guardar | Se escribe `07:05` y se guarda | «Sección guardada.»; BD `07:05:00`; al reabrir el control muestra `07:05` | DOM + BD |
| E11-C8 | Medianoche por pantalla | Se elige `00:00` y se guarda | BD `00:00:00`; al reabrir el control muestra `00:00` (no queda en blanco) | DOM + BD |
| E11-C9 | No se admite texto libre | Se intenta escribir «7pm» en el control | El navegador lo rechaza (`Malformed value`): no hay forma de enviar una hora con formato inválido | DOM |

### REG-08 — Tipo de discapacidad en pantalla (`wp2-reg08-tipo-discapacidad.js`)

| Caso | Descripción | Condición inicial | Resultado esperado | Método de verificación |
|---|---|---|---|---|
| E11-C10 | La pregunta aparece solo con DP_2 = Sí | Datos personales de un Caso con DP completo y DP_2 = No | Con No no se ve «Tipo de discapacidad»; al marcar Sí aparece con «Seleccione…» y las 7 opciones; al volver a No se oculta | DOM |
| E11-C11 | Elegir, guardar y reabrir | DP_2 = Sí; «Sordoceguera»; «Guardar sección» | «Sección guardada.»; BD con la opción 04; al reabrir DP_2 = Sí y la lista en «Sordoceguera» | DOM + BD |
| E11-C12 | Pasar DP_2 a No limpia la respuesta | Se marca No y se guarda | La pregunta se oculta, la BD queda en `NULL` y al reabrir sigue oculta | DOM + BD |

### REG-12 — Tamizaje en tablas «Pregunta | Sí | No» (`wp2-reg12-tamizaje-tablas.js`)

| Caso | Descripción | Condición inicial | Resultado esperado | Método de verificación |
|---|---|---|---|---|
| E11-C13 | Las generales respondidas se guardan | Caso con la tríada completa; `POST /tamizaje` con TZ_1, 3, 5, 7, 8 en Sí y la rama no pareja con 2 Sí | 200, puntaje 7, nivel alto; BD `tzv3_g1…g7` y `vpp` = `y,n,y,n,y,n,y,y`; rama `np1,np2,np24` = `n,y,y` (antes las generales quedaban todas en `n`) | API + BD |
| E11-C14 | Al reabrir, las generales vuelven como se respondieron | `GET /victim-cases-v3/:icode` | `TZ_1…TZ_8` iguales a lo enviado, todas en `answeredKeys`; `tzResultado` con puntaje 7 y nivel alto (antes: todas `false`) | API |
| E11-C15 | Volver a guardar con otras respuestas | `POST /tamizaje` con otras generales | BD `n,y,y,n,y,n,y,n` y `GET` lo refleja | API + BD |
| E11-C16 | Rama pareja | `wasPartner = true`, TZP_3 y TZP_24 en Sí | Puntaje 2, nivel bajo; `tzv3_p3` y `p24` = `y`, `p1` = `n`; solo se escribe la rama activa | API + BD |
| E11-C17 | Caso inexistente | `POST /tamizaje` a un icode que no existe | Error (500), no 200 con un resultado (antes: 200 y las claves "marcadas") | API |
| E11-C18 | Estructura de las tablas | Pestaña «Tamizaje de riesgo» | 2 tablas con encabezado «Pregunta · Sí · No», columnas 70/15/15; 7 filas generales y 25 del resto (TZ_8 + 24); cada fila con 2 radios del mismo `name` y `aria-label` «Sí — …» / «No — …»; la segunda tabla oculta | DOM |
| E11-C19 | Revelado | Se marca «Sí» en TZ_1 | Aparecen TZ_5, la segunda tabla y el botón (TZ_5 estaba oculta) | DOM |
| E11-C20 | Guardar desde la pantalla y reabrir | TZ_1, TZ_8 y la primera pregunta de la rama en Sí; «Calcular y guardar Tamizaje» | Resultado «Puntaje: 3»; BD de las generales `y,n,n,n,n,n,n,y`; al reabrir TZ_1 en Sí y TZ_2 en No | DOM + BD |

### REG-34 — Tamizaje por etapas (`wp2-reg34-tamizaje-por-etapas.js`)

| Caso | Descripción | Condición inicial | Resultado esperado | Método de verificación |
|---|---|---|---|---|
| E11-C21 | Solo las 5 iniciales | `POST /tamizaje` con TZ_1, 2, 3, 4, 7 en No | 200, puntaje 0, nivel bajo; esas 5 en la BD (`n`) y en `answeredKeys`; TZ_5, TZ_6, TZ_8 y la rama **no** quedan respondidas | API + BD |
| E11-C22 | El resto se reinicia | Antes: TZ_1, TZ_5, TZ_8 y TZNP_2/24 en Sí; después se guardan solo las 5 iniciales en No | Las columnas del resto (`g5`, `vpp`, `np2`, `np24`…) quedan en `n`; `GET` devuelve TZ_5 y TZNP_2 en `false` | API + BD |
| E11-C23 | Finalizar acepta el Tamizaje con solo las iniciales | Caso con todo lo que pide Finalizar salvo el Tamizaje | Sin Tamizaje: 400 con `modulosIncompletos = ["TZ"]`; con las 5 iniciales en No: 200 y estado `te` (antes el botón no aparecía y no había forma de guardar) | API + BD |
| E11-C24 | El botón aparece con las 5 iniciales respondidas | Tamizaje abierto en un Caso con la tríada completa | Oculto sin responder y con 4 iniciales; visible con las 5 en No, con la nota «Las cinco preguntas iniciales están en «No»…»; el resto sigue oculto | DOM |
| E11-C25 | Solo viajan las 5 iniciales | Se guarda con todas en No | El `POST /tamizaje` lleva exactamente TZ_1, 2, 3, 4, 7; resultado «Nivel: bajo — Puntaje: 0»; TZ_5 y la rama sin responder en la BD | DOM + red + BD |
| E11-C26 | Reabrir el Tamizaje guardado | Se reabre el Caso | Las cinco en «No», botón visible y resto oculto; un «Sí» revela el resto | DOM |
| E11-C27 | Un Caso con una inicial en Sí muestra el resto al abrir | Tamizaje guardado con TZ_1 = Sí | Al abrir la pestaña ya están visibles TZ_5, la segunda tabla y el botón, con TZ_1 en Sí (antes había que volver a marcar) | DOM |

---

## Cobertura

| Flujo de prueba | Casos cubiertos | Script |
|---|---|---|
| `wp2-reg10` — C1, C2, C3 | E11-C1 | `cloudDevelopment/e2e/wp2-reg10-subtipos-violencia.js` |
| `wp2-reg10` — C4, C5 | E11-C2, E11-C3 | ídem |
| `wp2-reg07` — C6, C7, C8 | E11-C4, E11-C5 | `cloudDevelopment/e2e/wp2-reg07-zona-residencia.js` |
| `wp2-reg09` — C4 … C7 | E11-C6 … E11-C9 | `cloudDevelopment/e2e/wp2-reg09-hora-hecho.js` |
| `wp2-reg08` — U1, U2, U3 | E11-C10, E11-C11, E11-C12 | `cloudDevelopment/e2e/wp2-reg08-tipo-discapacidad.js` |
| `wp2-reg12` — C1, C2, C3, C4, C5 | E11-C13 … E11-C17 | `cloudDevelopment/e2e/wp2-reg12-tamizaje-tablas.js` |
| `wp2-reg12` — U1, U2, U3 | E11-C18, E11-C19, E11-C20 | ídem |
| `wp2-reg34` — C1, C2, C3 | E11-C21, E11-C22, E11-C23 | `cloudDevelopment/e2e/wp2-reg34-tamizaje-por-etapas.js` |
| `wp2-reg34` — U1 … U4 | E11-C24 … E11-C27 | ídem |
| Unitarias Go — `victim_case_v3_fields_test.go` | Columnas generales y del «resto» del Tamizaje (apoyan E11-C13 y E11-C22) | `src/internal/repository/` |

## Resultados de ejecución

**Fecha:** 2026-11-02 · **Build:** rama `claude/busy-edison-uk6g2a` en `http://127.0.0.1:9090` (sandbox cloud) · **BD:** Postgres local con la migración WP-2 aplicada · **Herramienta:** Playwright headless + `psql` + `go test`.

| Casos | Resultado | Antes del arreglo |
|---|---|---|
| E11-C1 … E11-C3 (REG-10) | ✅ 3 de 3 | C2 y C3 fallaban (6 subtipos con el texto de otra categoría) |
| E11-C4, E11-C5 (REG-07) | ✅ 2 de 2 | La lista salía vacía |
| E11-C6 … E11-C9 (REG-09) | ✅ 4 de 4 | C6 y C9 fallaban (campo de texto con placeholder, aceptaba «7pm») |
| E11-C10 … E11-C12 (REG-08) | ✅ 3 de 3 | La pregunta no existía |
| E11-C13 … E11-C20 (REG-12) | ✅ 8 de 8 | Fallaban C13, C14, C15, C17 y todos los de pantalla (no había tablas); solo pasaba C16 |
| E11-C21 … E11-C27 (REG-34) | ✅ 7 de 7 | Fallaban C22 y los de pantalla (el botón no aparecía); C21 y C23 ya pasaban en el API |

**Resumen:** 27 de 27 casos ✅. El hallazgo de la corrida —cómo el resultado de un guardado anterior se queda en pantalla y confunde una comprobación— se resolvió esperando la respuesta del POST antes de leer el resultado (no era un defecto del producto).
