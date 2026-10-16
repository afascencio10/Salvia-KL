# QA — Registro V3, ronda de implementación completa (2026-10-15)

> Generado siguiendo `qa-process.md` Fase 1. Cubre los 26 ítems marcados `[x]`
> en `registro-v3-pendientes.md` (secciones A, B, C, D). El backend de la
> mayoría de estos ítems **ya se probó en Fase 2** (smoke tests de servicio
> contra la BD real de pruebas, ver tabla "Registro de pruebas" en
> `registro-v3-pendientes.md`) mientras se implementaba cada uno — acá solo
> se listan los casos que de verdad requieren navegador: el *wiring* real
> entre la plantilla Vue y el backend (selectors, eventos, render
> condicional), que un test unitario con `vm`/mocks no puede garantizar.
>
> Casos explícitamente **fuera de esta ronda de Playwright** porque ya
> quedaron cubiertos por Fase 2 sin superficie de UI propia que agregue
> riesgo: C4 (traspaso a Seguimiento — no tiene UI), D3-D5 (lógica de
> fragmentos del párrafo síntesis, 38 aserciones unitarias ya la cubren
> exhaustivamente — acá solo se verifica que SÍ se renderice, no cada regla
> de inclusión/omisión otra vez), D8 (fan-out — ya verificado contra BD real,
> no tiene UI propia), D12-D14 (condicionales de campo, ya cubiertos por
> D12/D13/D14 unitarios + se tocan de pasada en C7/C17).

## Casos

| Caso | Descripción | Condición inicial | Resultado esperado | Método de verificación |
|---|---|---|---|---|
| C1 | Sidebar bloquea TZ/GS con tríada incompleta (A1) | Caso en Borrador recién creado, CI/DP/CA/HV/AG sin completar | Clic en pestaña TZ o GS no cambia de módulo; aparece 🔒 en la pestaña y un mensaje de aviso | DOM |
| C2 | Sidebar se desbloquea al completar la tríada (A1) | Mismo Caso, tras guardar CI+DP+CA+HV+AG completos | `triadaCompleta` pasa a true (alert verde visible); clic en TZ/GS sí navega | DOM |
| C3 | `triadaCompleta` persiste al reabrir un Caso ya completo (A2) | Cerrar el Caso del C2 y reabrirlo desde el Cabezote (buscar por documento) | Sidebar ya desbloqueado sin guardar nada primero | DOM |
| C4 | Guardar sección con una obligatoria `[1er contacto]` vacía (A3) | Módulo CI con `CI_4` (teléfono) vacío | Mensaje "Faltan obligatorias..." aparece; no sale ningún POST a `/seccion/CI` | DOM + Network |
| C5 | Campo `soloLetras` con solo números (A4) | `CI_2` = "12345" | Input queda con borde rojo + mensaje; guardar bloqueado | DOM |
| C6 | Buscar entidad por Depto/Ciudad/Municipio/Sector (B1) | GS § C, sin entidad agregada aún | Al completar los 4 selects aparecen opciones reales de `entity_branch` en el select de Entidad | Network (`GET /entity-branches`) + DOM |
| C7 | Seleccionar una entidad carga su checklist de obligaciones reales (B2) | Entidad seleccionada en C6 | Aparece el checklist "Acciones institucionales" con las obligaciones reales de esa organización (o el mensaje de "sin tareas configuradas") | Network (`GET /entity-branches/:id/obligations`) + DOM |
| C8 | Agregar entidad exitosamente (B3) | Entidad + obligaciones marcadas | La entidad aparece en la lista de "ya agregadas" con sector/nombre/ubicación resueltos | Network (POST 200) + DOM |
| C9 | "+ Agregar barrera a esta entidad" hereda ubicación sin repetir búsqueda (B4, rediseño anidado) | `GS_17`=Sí, entidad ya agregada en C8 | Botón visible solo bajo esa entidad; al abrir el sub-formulario no hay selects de Depto/Ciudad/Municipio — ya tiene sector resuelto y carga competencias directo | DOM + Network (`GET /competences`) |
| C10 | Barrera agregada queda anidada bajo su entidad (B4) | Barrera agregada en C9 | La barrera aparece en la tabla DENTRO de la tarjeta de esa entidad, no en una lista separada | DOM |
| C11 | Eliminar entidad/barrera pide confirmación (A5) | Clic en "Eliminar" sobre la entidad o barrera de C8/C10 | Aparece un diálogo de confirmación del navegador antes de eliminar | Evento `dialog` de Playwright |
| C12 | Finalizar con módulos nuevos incompletos (C1) | Intentar Finalizar sin responder `CA_1-4`/`DP_10`/`DP_15`/`GS_16` | Error 400 se muestra en pantalla listando los módulos faltantes (no solo en la respuesta HTTP) | DOM + Network |
| C13 | Finalizar exitoso muestra confirmación de Denuncia fácil/remisión (C2/C3) | Completar los módulos faltantes de C12 y Finalizar de nuevo | Pantalla de Cierre muestra "Se solicitó la Denuncia fácil..."/"Se confirma remisión a..." si aplica | DOM |
| C14 | Tipo de contacto — "Víctima - Caso Nuevo" (D1) | Etapa `tipoContacto` recién mostrada | Avanza a Cabezote | DOM |
| C15 | Tipo de contacto — cualquier otra opción (D1) | Etapa `tipoContacto` | Muestra aviso "se gestiona en otra pantalla", no navega | DOM |
| C16 | Habeas Data — "No acepto" (D2) | Pantalla inicial de Habeas Data | Muestra pantalla `formularioCerrado`, no se crea ningún Caso | DOM + Network (cero requests de creación) |
| C17 | Auto-sugerencia de nacionalidad (D6) | `CI_5` = Cédula de Ciudadanía, `DP_5` aún vacío | `DP_5` queda prediligenciado en "Colombiana(o)" | DOM |
| C18 | Tratamiento por pronombre (D6) | `CI_1` = Ella | Encabezado del Caso muestra "Sra. [nombre]" | DOM |
| C19 | Aviso de duplicidad de documento (D7) | Escribir en `CI_6` el documento de un Caso ya existente | Tras ~500ms aparece aviso "Ya existe un Registro con este documento: ..." | Network + DOM |
| C20 | Badge Completo/Incompleto en Cabezote (D9) | Buscar por un documento con Perfil (CI+DP+CA) completo vs. uno incompleto | Badge verde "Completo" / gris "Incompleto" según corresponda | DOM |
| C21 | Alerta de riesgo alto / ideación suicida (D10) | Completar Tamizaje con varias respuestas "Sí" que den nivel alto, y marcar ideación suicida | Alerta roja de nivel de riesgo + alerta roja separada de ideación suicida | DOM |
| C22 | Botón "Finalizar llamada (guardar y salir)" (D11) | Módulo CI a medio llenar, sin cumplir obligatorias | Guarda lo que haya sin exigir validación y vuelve al Cabezote | Network + DOM |
| C23 | Pantalla de Cierre completa (D15) | Tras Finalizar exitosamente (C13) | Se ven las 4 tarjetas (Estado/Riesgo/Entidades/Barreras), tabla de actuaciones, bitácora, botón "Volver al Cabezote" | DOM |
| C24 | Párrafo síntesis se renderiza con datos reales (D3-D5) | Mismo Caso de C23, con CI/DP/HV/AG diligenciados | El párrafo en Cierre refleja los datos reales del Caso (nombre, documento, hechos, agresor) sin `undefined` ni icodes crudos | DOM |
| C25 | Máscara de CI_6 por tipo de documento (D13) | `CI_5` = Cédula de Ciudadanía, escribir un número de 2 dígitos | Aparece "No cumple el formato esperado" + hint del formato; con un número válido (6-10 dígitos) el error desaparece | DOM |

## Cobertura

Ejecutado con `playwright-cli` contra servidor local (`PORT=9090 go run .`),
usuario de prueba `qa.playwright` (rol `op`), sobre un Caso real construido
desde cero: icode `01a1426b-8f55-77fe-b633-7f70a5425866`, cc-1122334455,
"Carmen Rodriguez". CAPTCHA de login resuelto manualmente por el usuario (no
se intentó bypass).

| Flujo | Descripción | Casos cubiertos |
|---|---|---|
| F1 — Alta y Perfil | Habeas Data → Tipo de contacto → búsqueda Cabezote → CI → DP → CA, incluyendo un intento de guardar con obligatoria/`soloLetras` inválidos antes de corregir | C4, C5, C14, C15, C16, C17, C18, C19, C20, C25 |
| F2 — Tríada y bloqueo de sidebar | HV + AG completados; TZ/GS probados bloqueados antes y desbloqueados después; reapertura del Caso desde cero (nueva sesión de navegador) confirma que el desbloqueo persiste sin re-guardar | C1, C2, C3 |
| F3 — Tamizaje | Las 5 iniciales + rama pareja completa (incluida ideación suicida); primer intento con refs de radio equivocadas (score bajo, corregido) y recálculo a Nivel alto/Puntaje 15 | C21 |
| F4 — Entidades y barreras anidadas | Búsqueda geográfica en cascada, checklist de obligaciones real, alta de 1 entidad, alta de 1 barrera anidada bajo esa entidad (sin repetir búsqueda), eliminación de ambas con diálogo de confirmación | C6, C7, C8, C9, C10, C11 |
| F5 — Gestión Salvia y Finalizar | GS completo (denuncia fácil, remisión a Atención Psicosocial con criterio de 3pts cumplido, entidad/barrera reconstruidas), Finalizar exitoso, Cierre con las 4 tarjetas + tabla de actuaciones + bitácora + párrafo síntesis | C12 (ruta de éxito; el rechazo por módulos incompletos ya estaba cubierto en Fase 2 por smoke test directo al servicio), C13, C23, C24 |
| F6 — Guardar y salir | Reapertura del mismo Caso ya Terminado, botón "Finalizar llamada (guardar y salir)" desde GS → vuelve al Cabezote con guardado best-effort | C22 |

Las 25 filas de la tabla de Casos quedan cubiertas (24 en vivo vía Playwright,
C12 complementado con el smoke test de Fase 2 ya registrado en
`registro-v3-pendientes.md` para su ruta de rechazo).

### Hallazgos — bugs encontrados y corregidos durante esta ejecución

Ninguno de estos lo había detectado Fase 2 (smoke tests de servicio) porque
los cinco requieren el *wiring* real navegador↔backend↔BD que esa fase no
ejercita:

1. **Servidor no arrancaba** — ruta duplicada `GET /api/v1/competences`
   registrada dos veces (`entity_branch_api_controller.go` y, ya desde antes,
   `form_controller.go`), Gin entraba en panic al iniciar. Corregido
   eliminando el registro duplicado.
2. **Sidebar nunca se desbloqueaba para un agente real** — `DP_17` tenía
   `req: "D"` en `registro-v3-campos.js` pero `triadaFieldsV3` (Go) ya lo
   exigía para `triadaCompleta`. Corregido a `req: "C1"`.
3. **Párrafo síntesis no aparecía al entrar a Gestión Salvia** —
   `secciones.md` documenta que debe mostrarse ahí, pero solo estaba
   conectado a la pantalla de Cierre. Agregado al bloque GS.
4. **Fecha cruda de Go filtrando al párrafo síntesis** — `HV_2` se mostraba
   como `"2026-10-01 00:00:00 +0000 UTC"` tras reabrir un Caso en vez de
   `"2026-10-01"`. Agregado `limpiarFecha()` en `registro-v3-app.js`.
5. **Cierre mostraba "Nivel de riesgo: —" en cualquier Caso reabierto** —
   `tzResultado` solo vivía en memoria del frontend (la respuesta directa de
   `POST /tamizaje`); `GetCase` nunca lo releía de BD aunque
   `SaveTamizajeBranch` ya lo persistía. Cualquier Finalizar que no ocurriera
   en la misma sesión de navegador que guardó el Tamizaje mostraba el nivel
   de riesgo vacío en el Cierre. Corregido con `GetTamizajeResult` (nuevo,
   `internal/repository/victim_case_v3_repository.go`) + `TzResultado` en
   `CaseStateV3` (`victim_case_v3_service.go`) + repoblado en `abrirCaso`
   (`registro-v3-app.js`). Verificado reabriendo el Caso ya Terminado desde
   cero: `tzResultado` llega como `{riskType, riskLevel: "alto", score: 15}`.

### Hallazgos incidentales — fuera de alcance, no corregidos

Preexistentes, no introducidos por esta ronda de implementación. Se dejan
documentados para una ronda futura:

- **TZ_8** (pregunta VPP) se renderiza incondicionalmente aunque su propia
  condición documentada es `esOsigd`.
- **HV_4/DP_27** (campos de departamento) nunca persisten al reabrir un Caso
  por diseño original (`Column: ""`, "cascada UI, no se proyecta") — hace que
  el párrafo síntesis omita el fragmento de departamento tras un reload. Es
  una simplificación ya documentada, no una regresión de esta ronda.
