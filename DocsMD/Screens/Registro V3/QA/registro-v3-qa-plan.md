# Registro V3 — Plan de QA

Metodología: `qa-process.md`. Backend probado (Fase 2) contra la BD real de pruebas al diseñar este plan, y Fase 3 (Playwright CLI, navegador real, cuenta `claude.test.registrov3`/rol `op`) ejecutada 2026-10-15 sobre un recorrido end-to-end completo (Habeas Data → Cabezote → Formulario → Tamizaje → Gestión Salvia → Cierre, con `Finalize` incluido). Se encontraron y corrigieron **8 bugs en total** (H1-H8 abajo) — 3 antes de escribir los casos, 5 más durante la ejecución real en navegador. Ninguno de H5-H8 era detectable con pruebas de backend aisladas: dependían del payload real del frontend, de cómo Vue compila el template en el navegador, o de una colisión de valores que solo aparece con datos reales (un documento "cc" de verdad).

## Cobertura ejecutada (2026-10-15, Playwright CLI) — sesión completa

| Archivo de casos | Casos confirmados PASS | Pendiente (variantes negativas/de frontera, no exploran código nuevo) |
|---|---|---|
| `E04-E06-cabezote-qa-cases.md` | C1, C2, C5, C7, C8, C9, C11 | C3, C4, C6, C10, C12, C13 |
| `condicionales-formulario-qa-cases.md` | D1, D5, D7, D10, D13, D14, D18, D22 | D2-D4, D6, D8-D9, D11-D12, D15-D17, D19-D21, D23-D26 |
| `triada-y-guardado-qa-cases.md` | T1, T2, T3, T4, T14 (implícito vía reapertura) | T5-T13, T15 |
| `tamizaje-qa-cases.md` | TZ1, TZ2, TZ5, TZ10, TZ13 | TZ3-TZ4, TZ6-TZ9, TZ11-TZ12, TZ14, TZ16-TZ22 (mismo código, otros umbrales/ramas — ver nota) |
| `gestion-salvia-remision-qa-cases.md` | R1 (implícito), R2, R4, R7, R8, R12-R19 — **los 6 equipos y sus 4 variantes de criterio (obligatorio-solo, puntaje-mínimo, umbral-1-de-N, sin-umbral) quedaron confirmados** | R3, R5-R6, R9-R11 |
| `entidades-barreras-qa-cases.md` | E1, E2, E6, E7, E8, B1, B2, B3, B4, B6, B8, B9, B10, B12 | E3-E5, B5, B7, B11, B13 |
| `finalizar-y-alerta-qa-cases.md` | A1-A5, F2, F3 | F1, F4 |
| `registro-v3-e2e-qa-cases.md` | E2E-1 completo (Borrador→Iniciado→Terminado, con remisión/entidad/barrera/alerta/párrafo síntesis). E2E-2 confirmado como efecto colateral: el caso se reabrió correctamente **5 veces** tras otros tantos reinicios del servidor, siempre repoblando solo lo realmente respondido | E2E-3 (rama no-pareja + esOsigd + nivel extremo) — mismo código ya verificado con otra combinación de datos (TZ10/TZ11/TZ22), no se construyó un segundo Caso completo para esto |

Todos los casos marcados PASS se confirmaron contra el navegador real (Chrome vía `playwright-cli`), la respuesta real del endpoint, y/o una lectura directa de la base de datos de pruebas — no son inferencias de código. Caso de prueba: icode `01a13dd1-eb74-7392-9aa7-b7a345df7ab9`, documento `777888999`, estado final `te` (Terminado).

**Nota sobre lo que queda pendiente:** los casos sin ejecutar son, en su mayoría, la "otra mitad" de un mecanismo ya confirmado con datos reales (ej. probar que HV_12 se oculta cuando NO se marca "laboral", habiendo ya confirmado que aparece cuando sí se marca) o variantes de umbral sobre la misma tabla de puntajes ya verificada en código y en al menos un punto real. Se dejaron documentados en sus archivos para quien quiera correrlos como regresión, pero no se consideran de alto riesgo dado lo ya confirmado.

## Alcance real vs. planeado

`registro-v3-index.md` describe 18 eventos, pero varios se escribieron **antes** del pivote a "sin dinamic-form" (decisión A1) y de la implementación real. El código manda sobre el MD (`dev-docs-system.md`). Estos casos prueban lo que **existe hoy**:

| Implementado y cubierto aquí | Planeado en el MD pero NO implementado (no generar casos todavía) |
|---|---|
| E01 carga pantalla → Habeas Data | E08 "Nuevo Caso en Registro existente" (copiar perfil de un Caso previo) |
| E02 Habeas Data | E09 "Actualizar Caso" (reabrir un Caso Iniciado) |
| E04 Cabezote — búsqueda doc/teléfono (D7) | E16 accesos de sidebar Rutas/Barreras/Notificaciones (diferido, D2) |
| E06 "Nuevo Registro" (crear borrador) | Auto-guardado por cada tecla (E11 lo describe así; la implementación real guarda por **sección completa**, vía botón "Guardar sección") |
| E10/E11 navegación entre módulos + condicionales | Alertas automáticas de riesgo/ideación suicida mientras se escribe |
| E12 guardar sección + tríada → Borrador/Iniciado | Validación de obligatorias en Finalizar (ver hallazgo H4 abajo — no está implementada) |
| Tamizaje (gating, branching, score→nivel) | — |
| GS § A/B/C/D/E (gestión, remisión, entidades, barreras, alerta inminente) | — |
| E15 Finalizar → Terminado | — |

`E03` (tipo de contacto), `E05` (limpiar), `E07` (proceder al formulario), `E09`/`E14` (abandonar) tampoco tienen una pantalla de "Cabezote con resultados" standalone todavía — el flujo real actual es: Habeas Data → formulario de búsqueda+crear inline (una sola vista) → Formulario. Los casos de E04-E06 reflejan esto.

## Hallazgos corregidos durante la preparación de este plan (2026-10-14)

Encontrados probando el backend directamente contra la BD antes de escribir los casos (Fase 2 de `qa-process.md` exige esto). Quedan documentados porque explican por qué ciertos casos de abajo dicen "ya corregido" en vez de ser un caso de regresión nuevo.

| # | Archivo:línea | Problema | Corrección |
|---|---|---|---|
| H1 | `internal/repository/victim_case_v3_repository.go` (`insertDummyForm2Row`) | `CreateDraft` fallaba por violación NOT NULL — solo rellenaba columnas `bool_enum`/`enum1`, dejando sin valor `birth_date`, `facts_date`, `facts_description`, `residence_town`, etc. | Se rellenan todas las columnas NOT NULL por `Kind` (Text→"", Int→0, Date→"1900-01-01", Time→"00:00:00"), más una red de seguridad que audita `information_schema.columns` en caliente para cualquier columna NOT NULL sin default que ningún catálogo Go declare. |
| H2 | Mismo archivo | La red de seguridad forzaba `victim_case_form2_id=0` (pisando el `nextval` de la secuencia) y `victim_case_form2_i_code=''`, dejando la fila inubicable por `GetAnswers` (que interpreta `form2ID==0` como "no existe"). | Se excluyen del relleno las columnas con `column_default IS NOT NULL` (secuencias), y `i_code` se genera explícitamente con `utils.GetUUID()`. |
| H3 | `victim_case_v3_repository.go` (`SaveSection`) + `victim_case_v3_service.go` | La tríada (Borrador→Iniciado) se evaluaba leyendo los **valores** guardados (`GetAnswers`), pero un valor dummy (ej. "No" en un campo booleano) es idéntico a una respuesta real — la tríada se marcaba completa inmediatamente después de crear el Borrador, sin que el agente respondiera nada. | Se agregó `victim_case_v3_answered_keys` (columna nueva, `migrations/add_registro_v3_part4.sql`) que rastrea qué `FieldKey` llegaron alguna vez en un guardado real, independiente del valor. La tríada ahora se evalúa contra ese set (`GetAnsweredKeys`), no contra el contenido de BD. Verificado end-to-end: Borrador se mantiene tras guardar CI, DP, HV solos; pasa a Iniciado recién al completar AG (el último módulo de la tríada). |
| H4 | `victim_case_v3_repository.go` (`CreateDraft`) | El teléfono del criterio de búsqueda del Cabezote se perdía — nunca se escribía en `victim_case_form2_victim_phone` (quedaba en el dummy `0`), aunque D7/E06 dicen que el Borrador debe quedar "prediligenciado con los criterios de búsqueda". | Se agrega un `UPDATE` explícito del teléfono tras crear la fila dummy, y se marca `CI_4` como respondido. |

## Hallazgos corregidos durante la ejecución real en navegador (2026-10-15)

Estos NO se podían atrapar con pruebas de backend aisladas — dependían del payload real que arma el frontend o de cómo Vue compila el template en el navegador.

| # | Archivo:línea | Problema | Corrección |
|---|---|---|---|
| H5 | `internal/repository/victim_case_v3_repository.go` (`SaveSection`, sync a `victim_case`) | `CI_5` (tipo de documento) es `KindEnum1` — el frontend manda el **icode** (uuid), pero `victim_case_victim_doc_type` es `varchar(2)` y espera el **code** corto ("cc"). El guardado de la sección CI reventaba con `value too long for type character varying(2)` (HTTP 500) la primera vez que se probó con un payload real del navegador. | Se agregaron `resolveEnumCodeByICode`/`resolveEnumICodeByCategoryAndCode` y se resuelve en ambas direcciones (guardar → code; `GetAnswers` → icode, para que el dropdown se repueble bien al reabrir). |
| H6 | `frontend/html/salvia/registro_v3/registro_v3.html` + `registro-v3-app.js` | Vue, al compilar el template en el navegador (sin build step), reescribe cualquier identificador que no reconoce como dato/método/prop del componente a `_ctx.<nombre>` en vez de dejarlo caer al scope global — una referencia directa a una variable global del catálogo (`TAMIZAJE_VPP_V3.label`, usada sin pasar por un método) revienta el render completo con "Cannot read properties of undefined" y deja la pantalla en blanco. Afectaba Tamizaje (`TAMIZAJE_GENERALES_V3`, `TAMIZAJE_VPP_V3`) y Gestión Salvia §C/D (`SECTORES_ENTIDAD_BARRERA_V3`, `CANAL_ACTIVACION_V3`). | Se expusieron los 4 catálogos como datos del componente (`data()`) en vez de referencias sueltas al global. Cualquier catálogo nuevo que el HTML use directo (no envuelto en un método) debe agregarse ahí también. |
| H7 | `frontend/html/salvia/registro_v3/registro_v3.html` (botón "+ Siguiente barrera") | `victim_case_barrier_date`/`_official` son NOT NULL en BD, pero el botón solo exigía `SectorCode`/`Description` — enviar el formulario sin fecha reventaba `POST /barreras` con `invalid input syntax for type date: ""` (HTTP 500). | Se agregó `nuevaBarrera.Date`/`.Official` a la condición `:disabled` del botón, igual de exigente que lo que la BD realmente requiere. |
| H8 | `victim_case_v3_repository.go` (`GetAnswers`, sync desde `victim_case`) + `registro-v3-app.js` (`parrafoSintesis`) | Dos bugs relacionados encontrados en la misma verificación: (a) `GetAnswers` comparaba `core.DocType` contra el dummy `"cc"` para decidir si CI_5 "ya se respondió" — pero `"cc"` (Cédula de Ciudadanía) es también el tipo de documento real más común en Colombia, así que cualquier víctima con cédula real perdía ese campo al reabrir el Caso; (b) el párrafo síntesis mostraba el **icode** crudo del equipo de remisión en vez de su nombre ("Se recomienda remisión a: 6e9e4650-...") porque comparaba un icode contra el `code` de 2 letras de `EQUIPOS_REMISION_V3`. | (a) Se quitó la comparación contra el dummy — el filtro correcto de "¿lo respondió el agente?" ya lo hace `answeredKeys` (H3) en el frontend, no hace falta duplicarlo aquí y la duplicación era la fuente del bug. (b) `parrafoSintesis` ahora resuelve icode→code vía el catálogo de opciones antes de buscar en `EQUIPOS_REMISION_V3`, igual que el resto de la pantalla. |

Nota aparte (no bloqueante, se corrigió de pasada): el `<body>` del layout compartido (`header.html`) ya usa `id="app"` como mecanismo real de mostrar/ocultar la página — la plantilla de Registro V3 dupli­caba ese mismo id en su propio `<div>` contenedor, así que `Vue.mount("#app")` y el script de mostrar/ocultar apuntaban al `<body>` por casualidad de orden en el DOM, no al div. Se quitó el id duplicado del div (mismo patrón que usa `list_reportes.html`, que nunca puso `id="app"` en su contenedor).

## Hallazgos NO corregidos — documentados como comportamiento actual a validar

| # | Descripción | Por qué no se corrigió ahora |
|---|---|---|
| G1 | El mecanismo de empate Feminicidio/Homicidio abyecto (decisión B6, "tentativa de feminicidio" / "tentativa de homicidio abyecto o fútil" como subtipo de HV) es **inalcanzable**: el catálogo real `victim_case_form2_subtype_violence_experienced_fi` no tiene ninguna opción con ese texto (las 7 reales son Ataque con agente químico/Golpizas/Quemadura/Sofocamiento/Agresión con objetos/Tortura/Secuestro). El tipo de riesgo hoy se decide **solo** por `esOsigd` (D3), nunca por el override de subtipo. | Requiere decisión de producto: ¿se agregan esas 2 opciones al catálogo real (nueva migración + `Locale.go`), o se descarta el override y D3 queda solo por identidad/orientación? No es un bug de código, es un gap de datos/alcance. |
| G2 | `Finalize` (E15) no valida nada — cambia el estado a `Terminado` sin revisar que las obligatorias estén diligenciadas, y puede ejecutarse incluso con el Caso todavía en `Borrador` (nunca pasó por `Iniciado`). | El MD original (E15) describe una validación que nunca se implementó en esta ronda; no estaba en el alcance de "decide lo demás tú" de la sesión de implementación. Los casos de abajo prueban el comportamiento **actual** (sin validación), no el planeado. |
| G3 | `GS_17` (identifica barreras) oculta/muestra el formulario de agregar barrera, pero `GS_14` (reportó a entidad) **no** oculta el formulario de agregar entidad — siempre está visible. Inconsistencia entre los dos repeaters. | Es ambiguo si debería gatear o no (podría ser intencional: registrar una entidad "acudida" no depende de si ya se reportó formalmente). Se deja como caso de verificación de comportamiento, no de corrección. |
| G4 | `AG_2 = "No se puede determinar / no informa"` deja visible `AG_5` (depende económicamente del agresor) — el código solo oculta `AG_5` cuando la etiqueta contiene "desconocid". | Ambigüedad de negocio (ver caso específico abajo), no un bug claro. |

## Archivos de casos

| Archivo | Cobertura |
|---|---|
| [E04-E06-cabezote-qa-cases.md](E04-E06-cabezote-qa-cases.md) | Búsqueda (D7), crear borrador prediligenciado, duplicados |
| [condicionales-formulario-qa-cases.md](condicionales-formulario-qa-cases.md) | Los 13 condicionales de CI/DP/HV/AG — foco principal pedido por el usuario |
| [triada-y-guardado-qa-cases.md](triada-y-guardado-qa-cases.md) | Guardado parcial por sección, Borrador→Iniciado (H3), multi-valor (enumN) |
| [tamizaje-qa-cases.md](tamizaje-qa-cases.md) | Gating 5 preguntas, rama pareja/no-pareja, tipo/nivel/score, empate (G1) |
| [gestion-salvia-remision-qa-cases.md](gestion-salvia-remision-qa-cases.md) | Equipos de remisión + criterios con puntaje |
| [entidades-barreras-qa-cases.md](entidades-barreras-qa-cases.md) | Repeaters de Entidades y Barreras |
| [finalizar-y-alerta-qa-cases.md](finalizar-y-alerta-qa-cases.md) | Alerta inminente (ya no es botón, C7), Finalizar (G2) |
| [registro-v3-e2e-qa-cases.md](registro-v3-e2e-qa-cases.md) | 2 recorridos completos de punta a punta |

## Orden de ejecución (ver `qa-process.md`)

1. Backend ya probado para `CreateDraft`/`SaveSection`/tríada (ver hallazgos arriba) — falta probar `SaveTamizaje`, `Finalize`, `AddEntity`/`AddBarrier`, `GetEnumCatalog` endpoint por endpoint antes de pasar a Playwright.
2. Playwright CLI — explorar interactivamente cada caso de los archivos de abajo, confirmar selectores/labels reales contra el catálogo de enums (`GET /api/v1/victim-cases-v3/enums`) antes de fijarlos en un `.spec.ts`.
3. Formalizar como script solo los recorridos que se van a repetir (regresión) — ver criterio en `qa-process.md`.
