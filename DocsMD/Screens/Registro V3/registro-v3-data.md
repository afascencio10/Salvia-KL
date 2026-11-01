# Registro V3 — Decisiones de arquitectura y mapeo a Base de Datos

> Complementa `registro-v3-secciones.md` (catálogo de preguntas). Este
> documento resuelve **dónde cae cada cosa en BD** y las 3 decisiones de
> alcance confirmadas con el usuario antes de planear.
>
> **Estado a 2026-10-31.** Nació como planeación (D1-D9, antes de construir) y D10 pivotó la arquitectura. Las secciones
> de decisiones D1-D9 se conservan como historia y razonamiento; lo que **quedó en el código** está en las tablas
> "Estado de las decisiones", "Columnas reales en `victim_case_form2`" y "Estado de los GAPS". Donde difieren, manda el código.

## Estado de las decisiones (código real)

| Decisión | Qué quedó en el código |
|---|---|
| D1 — `victim_case` único | Sin tabla Registro. Se creó la **sincronización del Perfil**: al guardar CI/DP/CA el backend replica el cambio en los demás Casos con el mismo tipo+número de documento (`fanOutPerfilV3`, best-effort, sin control de concurrencia). "Registro completo/incompleto" se calcula en `Search` sobre el Caso más reciente del documento (CI + DP + DP_10/DP_15 + CA_1-4). **No** existe "Nuevo Caso en Registro existente" (copia del Perfil) ni "Actualizar Caso" |
| D2 — Expediente de víctima | Sigue diferido; el botón del Cierre se descartó |
| D3 — Tipo de Tamizaje | Implementado en `SaveTamizaje`: tipo por `esOsigd` (el cliente lo calcula con el texto de DP_8/DP_9) y override por subtipos de tentativa de HV_10; si vienen ambos gana Feminicidio y se guarda `tie_flag`. La pregunta VPP (TZ_8) solo suma punto y se muestra siempre que se revele el resto |
| D4 — 5 preguntas iniciales | TZ_1, 2, 3, 4 y 7; el resto se revela si **alguna** es "Sí". Brecha: con las 5 en "No" el Tamizaje no se puede guardar (ver G-V3-23) |
| D5 — Estados y observaciones del cliente | Estados `bo`/`in`/`te` en V3 (`cd` lo pone Seguimiento); "solo letras" y correo de la víctima (DP_50) implementados; DIVIPOLA **no** (G-V3-10) |
| D6 — Aportes del mockup | Adoptados: botón "Finalizar llamada (guardar y salir)" y pantalla de Cierre (simplificada). **No** adoptados: checkbox "otra persona", modales de motivo y de Registro Completo, trazabilidad del Perfil, letreros OSIGD/Debida diligencia, marca de gestación, "solo 1er contacto" |
| D7 — Búsqueda sin nombre | Implementada (`/buscar`); con documento ignora el teléfono; no consulta reportes |
| D8 — Borrador con valores dummy | Implementado (`insertDummyForm2Row`) y ampliado: "respondido" ya no se infiere del valor sino de `victim_case.victim_case_v3_answered_keys` |
| D9 — Modelo real de Seguimiento | Los 6 equipos y sus criterios, sí. Entidades en `entity_case`/`entity_case_obligation` y Barreras en `barrier_v2`, **no** como repeaters de `dinamic-form` (G-V3-21 resuelto); las barreras cuelgan de una entidad |
| D10 — Pivote (A1, A3, B4-B6, C7, C9, D11) | Todo implementado: sin `dinamic-form`; enums desde `victim_case_form2_enums` por `GET /victim-cases-v3/enums`; empate → Feminicidio; Alerta Inminente solo en GS (GS_18/GS_19); 5 opciones de tipo de contacto (4 con aviso); `general_user` + credenciales en texto plano en `victim_case_new_user_credentials` |

---

## Decisiones de arquitectura confirmadas

### D1 — Modelo de datos: se mantiene `victim_case` único (sin separar Registro/Caso)

El Excel fuente (pestaña "Ajustes generales Registro") pide una
reestructuración completa: una entidad **Registro** (identidad única de
víctima, PK = Tipo+Número de documento, estados Incompleto/Completo) separada
de una entidad **Caso** (1 Registro → N Casos, estados Borrador/Iniciado/
Terminado). Hoy esa separación **no existe** — todo vive en una fila de
`victim_case` (confirmado: no hay tabla `victim_registry` ni equivalente).

**Decisión:** no se hace esa separación en esta ronda. Se mantiene una fila
de `victim_case` por Caso. Mecanismo de compensación:

- **Perfil de víctima duplicado por Caso**: cuando el Cabezote detecta un
  Registro existente (mismo Tipo+Número de documento) y el agente elige
  "Nuevo Caso en Registro existente", se crea un `victim_case` nuevo y se
  **copian** (no se referencian) los campos de CI+DP+CA desde el `victim_case`
  más reciente con el mismo documento — ver
  [flow-E08](Flujos/flow-E08-cuando-presiona-nuevo-caso-en-registro.md).
- **Sincronización por "fan-out" al editar el Perfil**: si se edita un campo
  de CI/DP/CA en un Caso, el mismo cambio se replica en los demás `victim_case`
  que comparten Tipo+Número de documento — ver el sub-flujo correspondiente en
  [flow-E11](Flujos/flow-E11-cuando-actualiza-una-respuesta.md). No hay
  bloqueo de concurrencia entre agentes editando el mismo Perfil desde Casos
  distintos al mismo tiempo — ver GAPS.
- **"Registro completo/incompleto" es un valor calculado, no una columna**:
  se deriva en el momento de la consulta (Cabezote, badges) evaluando si el
  **Caso más reciente** con ese Tipo+Número de documento tiene diligenciadas
  todas las obligatorias de CI+DP+CA. No se persiste.

> **Recomendación a futuro (no se implementa ahora):** la separación real
> (tabla `victim_registry` + FK desde `victim_case`) sigue siendo el modelo
> correcto a mediano plazo — resuelve de raíz la duplicidad y el riesgo de
> desincronización del fan-out. Queda como GAP G-V3-01.

### D2 — Alcance: núcleo del formulario primero

Esta ronda cubre: Habeas data → Cabezote de búsqueda → Formulario de Ingreso
Registro (7 módulos, incl. Tamizaje) → Gestión Salvia. **Expediente de
víctima** (pantalla de consulta con sus 5 vistas: Perfil, Casos asociados,
Detalle de Caso, Personas relacionadas, Reportes de origen) se planea en una
ronda posterior, porque su diseño depende de cómo termine resuelto D1 (sin
tabla Registro propia, "Casos asociados" se consulta filtrando `victim_case`
por documento, lo cual es suficiente para esta fase pero debe revisarse si
D1 cambia más adelante).

También quedan fuera, por ser módulos/formularios distintos no tocados en
esta ronda: Formulario de Seguimiento (ya existe), Formulario de Reporte
(Tercero), Registro de feminicidio (equipo EAIF), módulo Rutas y Barreras
(incluye las 3 preguntas de "denunció previamente / entidad / atención
efectiva" que el propio Excel marcó como eliminadas de este formulario), y la
grilla de Asignación de Ruta (ya excluida desde V2-real).

### D3 — Tipo de Tamizaje: se decide por identidad/orientación directamente, NO por la pregunta VPP

> ⚠️ **Corregido por feedback del cliente** (`Observaciones y ajustes al
> MockUp de registro 2.0.pdf`, nota #5 — este PDF tiene prioridad sobre el
> Excel/Guion Funcional de Drive por ser retroalimentación reciente del
> cliente). Texto literal: *"Clasificación VPP prioriza criterios de
> identidad/orientación en las filas 19-20. La pregunta se suprime como
> criterio definitorio del árbol de decisiones del tipo de tamizaje."*
>
> Esto **reemplaza** la lógica que el Guion Funcional de Drive describía
> (donde la respuesta a la pregunta fuerte VPP de TZ decidía el tipo).

Regla implementada (versión corregida):

1. El tipo de Tamizaje (Feminicidio vs. Homicidio abyecto o fútil) se decide
   **directamente** a partir de `identidad_genero`, `orientacion_sexual` y
   `sexo_asignado_nacer` (módulo DP) — **no** a partir de la respuesta a
   ninguna pregunta de Tamizaje:
   - `esOsigd` (mismo flag que gatilla preguntas condicionales) `== true`
     ⇒ tipo por defecto = **Homicidio abyecto o fútil**
   - `esOsigd == false` ⇒ tipo por defecto = **Feminicidio**
   - Al ser un cálculo sobre campos de DP, el tipo puede determinarse **tan
     pronto el módulo DP esté diligenciado**, incluso antes de llegar a TZ.
2. La pregunta fuerte VPP ("¿Esta persona ha realizado amenazas, agresiones
   o comentarios motivados por prejuicio...?") **se mantiene** en el banco de
   Tamizaje, condicionada a `esOsigd` igual que antes, y **sigue sumando su
   punto al score** — pero ya **no determina el tipo**, solo es una pregunta
   de puntaje más.
3. **Override por subtipo de Hechos de violencia (decisión previa confirmada
   con el usuario, no contradicha por el PDF del cliente):** si en HV
   (pregunta 10 — subtipo) se marcó **"Tentativa Feminicidio"** ⇒ el tipo se
   fuerza a **Feminicidio** sin importar el resultado de DP. Si se marcó
   **"Tentativa Homicidio abyecto o fútil"** ⇒ se fuerza a ese tipo. Este
   override sigue aplicando **sobre** el nuevo criterio por defecto
   (identidad/orientación), en el mismo orden de prioridad que antes.
4. **Caso borde no resuelto por ninguna fuente — GAP G-V3-02:** si se marcan
   **ambos** subtipos de tentativa a la vez, no hay regla de prioridad
   definida. Mientras se confirma con el equipo, se propone: prioridad al
   tipo que coincida con el resultado de identidad/orientación (si uno
   coincide); si ninguno coincide, se asume Feminicidio por defecto y se
   marca el Caso con una bandera de "revisión manual de tipo de riesgo".

### D4 — Tamizaje: 5 preguntas obligatorias iniciales, el resto se revela solo si alguna fue "Sí"

> **Corregido por feedback del cliente** (observación #15). Texto literal:
> *"En el tamizaje, ajustar que se visualicen primero las cinco preguntas
> obligatorias en todos los tamizajes pues no en todos los casos se debe
> aplicar la totalidad de las preguntas. Cuando alguna de esas cinco
> iniciales se responde afirmativamente se despliegan las demás. Cuando las
> primeras 5 se responden 'No', no debería visualizarse el resto de
> preguntas."*
>
> Esto **reemplaza** el modelo "7 preguntas generales siempre visibles,
> después bifurca por pareja/no-pareja" que traía el Guion Funcional de
> Drive. El modelo corregido es: **5 preguntas siempre visibles primero** →
> si **todas** responden "No", el Tamizaje termina ahí (nivel Bajo, sin más
> preguntas) → si **alguna** responde "Sí", se revela el resto (las 2
> preguntas generales restantes + la pregunta VPP si aplica + la rama
> pareja/no-pareja completa).

**GAP G-V3-08 — cuáles son exactamente las 5 preguntas obligatorias:** ni
este PDF ni el Excel de Drive enumeran por nombre cuáles 5 de las 7
"generales" de Drive son las 5 iniciales. El prototipo React de una ronda
anterior (que ya había leído este mismo PDF) usó: *aumento de violencia
física, uso/amenaza con arma, amenaza de muerte, percepción de capacidad de
matar, e ideación de autolesión* — pero "ideación de autolesión" no es una
de las 7 generales de Drive (en Drive vive como ítem final de cada rama, no
como general). Mientras se confirma con el equipo, se propone usar estas 5
de las 7 generales de Drive como las iniciales obligatorias: **(1) aumento
de violencia física en 60 días, (2) uso de armas u objetos, (3) amenaza de
muerte, (4) asfixia o estrangulamiento, (7) acceso a armas** — dejando (5)
persigue/espía y (6) notas amenazantes como parte del bloque que se revela
después. Esta selección es una propuesta razonada, no una confirmación.

> También confirmado por el cliente en la misma ronda de observaciones: los
> rangos de ponderación **siguen pendientes de ajuste final** (nota #3) — no
> recalibrar, se mantienen los umbrales ya documentados en `flow-E11` tal
> cual hasta nueva instrucción.

### D5 — Otros ajustes confirmados por el feedback del cliente

| # obs. | Ajuste | Dónde queda reflejado |
|---|---|---|
| 1 | **Estados del Caso son 4, no 3:** Borrador → Iniciado → Terminado → **Cerrado**. El Guion Funcional de Drive solo documentaba 3 (sin Cerrado) — el cliente confirma que Cerrado sí existe. Esta pantalla (Registro V3) solo llega hasta **Terminado** (handoff a Seguimiento); la transición Terminado→Cerrado es propiedad del Módulo de Seguimiento, fuera de esta pantalla — pero el badge debe reconocer el valor Cerrado si el Caso ya pasó por Seguimiento y se reabre un link directo (flow-E01). | `registro-v3-interface.md` (BadgeEstadoCaso), `flow-E15` |
| 2 | Habilitación secuencial de TZ/GS **"pendiente incorporación de situación Alerta inminente"** — el propio cliente reconoce que todavía no está resuelto si una Alerta inminente debe saltarse el requisito de tríada completa. No se implementa ningún salto todavía. | GAP G-V3-09 |
| 6 | **Confirma** que el criterio de obligatoriedad por pregunta del Excel (1er contacto / máx. 2º contacto / Deseable) es el que manda — incluye explícitamente Contacto de apoyo y Denuncia fácil, resolviendo a favor del Excel la contradicción interna que tenía el Guion Funcional de Drive (que en un párrafo las llamaba "solo opcionales"). | Ya reflejado en `registro-v3-secciones.md` — no requirió cambio, solo queda confirmado |
| 10 | **DIVIPOLA**: en producción, Departamento/Municipio/Zona (residencia, hechos, atención) deben asociarse a las tablas maestras oficiales de DIVIPOLA, no a listas ad-hoc. V2-real ya tiene un mecanismo de catálogo geográfico (`departments`, `geo.*.towns`) — confirmar que es DIVIPOLA oficial antes de reutilizarlo. | GAP G-V3-10 |
| 11 | **Verificar alternativas de direccionamiento** desde Habeas data / Tipo de contacto hacia el Formulario Familiares (feminicidio) y el Formulario Hombres — ninguno de los dos aparece como opción en la pestaña Habeas data del Excel, pero sí existen como formularios hermanos mencionados en el Guion Funcional. Posible 4ª/5ª opción de "Tipo de contacto" pendiente de confirmar. | GAP G-V3-11, `flow-E03` |
| 12 | **Nuevo campo**: "Correo electrónico de la víctima", opcional, en el módulo DP. No está en el Excel de Drive — es una adición del cliente. | `registro-v3-secciones.md` § DP |
| 13 | **Validación "solo letras"** en campos de nombre (Nombres y apellidos, Nombre identitario, Nombres del contacto de apoyo, Nombres del agresor) — deben rechazar una entrada compuesta solo por números. | `registro-v3-secciones.md` (nota por campo), `flow-E11` |
| 16 | **Ubicación del botón de Alerta inminente** — el propio cliente pregunta dónde debería quedar; no está resuelto. Se mantiene la propuesta de este plan (header del Formulario) como punto de partida, no como respuesta confirmada. | GAP ya existente en `flow-E13`, reforzado aquí |
| 17 | **3er grupo del sidebar: "Registro de gestión"** (Rutas / Barreras / Notificaciones) — el mockup original tiene un tercer bloque de navegación paralelo a "Datos y Registro de la víctima" y "Datos del Caso de violencia" que esta planeación había omitido por completo. El cliente mismo pregunta si la conexión con esos módulos es vía las preguntas de GS o vía botones aparte — **no resuelto**. | `registro-v3-interface.md` (sidebar), GAP G-V3-12, nuevo evento E16 |
| 18 | **Párrafo síntesis** — funcionalidad completa que esta planeación no incluía. Ver `registro-v3-parrafo-sintesis.md` (nuevo). | `registro-v3-parrafo-sintesis.md`, `registro-v3-interface.md` § GS, `flow-E11` |

### D6 — Aportes de `Mockup_Reportes_Registro (1).html` (prioridad más baja — solo adiciones, nunca sobrescribe)

> Al revisar el mockup HTML original se confirmó que es la **versión previa**
> a `Observaciones y ajustes al MockUp de registro 2.0.pdf` — su propio
> arreglo interno `NOTAS_REVISION` contiene prácticamente las mismas 10
> notas que después se convirtieron en ese PDF de observaciones del cliente.
> Donde el mockup y el PDF de observaciones difieren, **gana el PDF** (ya
> aplicado en D3/D4/D5). Lo que sigue son únicamente adiciones que no
> chocan con nada ya definido; las adopté en los flujos correspondientes:
> checkbox "verifiqué que es otra persona" (E06), botón explícito "Finalizar
> llamada" (E14), modal "Editar con trazabilidad" + motivo de edición (E11),
> modal "Registro Completo" (E11), pantalla de Cierre con resumen antes de
> ir a Seguimiento (interface.md + E15), texto verbatim de los letreros
> Debida Diligencia/Respeto OSIGD (interface.md), banner de vulnerabilidad
> por embarazo (interface.md), toggle "solo 1er contacto" (interface.md),
> formato concreto de trazabilidad `{ts}·{agente}·{entidad}·{campo}:
> {antes}→{después} [{motivo}]` (reflejado en `flow-E11`).

**Explícitamente descartado por contradecir decisiones ya tomadas** (no se
adoptó nada de esto):
- El mockup asume en varios lugares que "el Perfil no se duplica" — contradice D1 (se mantiene la duplicación + fan-out ya decidida).
- La lógica de tamizaje del mockup todavía usa "VPP decide el tipo" — es la versión *anterior* a la corrección D3, no se adopta.

**Adoptado tras discutirlo con el usuario** (ya no es nota abierta):
- El `NOTAS_REVISION` del mockup (nota 1) proponía una interpretación
  distinta a la que este plan tenía originalmente en `flow-E12` para
  Borrador→Iniciado ("Iniciado en el primer guardado", sin exigir ningún
  campo mínimo). Se adoptó la propuesta del mockup — **Borrador = primer
  dato guardado; Iniciado = cuando la tríada obligatoria de 1er contacto
  queda completa** — unificándola además con el mismo umbral que desbloquea
  TZ/GS (antes eran dos criterios distintos: ver D1/D-tríada abajo). También
  resuelve de forma más limpia la ambigüedad "Caso Abierto" que se había
  detectado en el Guion Funcional. Ya aplicado en `flow-E11` (sub-flujo
  "Recalcular habilitación secuencial") y `flow-E12` (PASO 2).

**Adoptado como propuesta razonada para un GAP que seguía abierto** (ninguna
fuente de mayor prioridad lo resolvía):
- G-V3-09 (Alerta inminente vs. habilitación secuencial) seguía sin
  resolver tras el PDF de observaciones. El mockup propone una regla
  relacionada y más acotada: si el subtipo de violencia marcado es
  "Tentativa Feminicidio" o "Tentativa Homicidio abyecto o fútil", el módulo
  de **Tamizaje** se habilita de inmediato aunque la tríada no esté
  completa (Gestión Salvia sigue bloqueada hasta completar la tríada). Es
  razonable — permite correr la valoración de riesgo en una situación ya
  marcada como grave sin esperar a terminar de diligenciar Agresor — pero
  sigue siendo una propuesta del mockup, no una confirmación del cliente.
  Se agrega como GAP G-V3-14 en vez de implementarse directamente.

### D7 — Búsqueda del Cabezote: solo documento y teléfono, nunca por nombre

> **Corregido por un diagrama de flujo enviado directamente por el cliente**
> (imagen, 4 cajas): *"Antes de permitir crear un caso usuario debe Buscar
> por numero de telefono y cedula (no nombre y apellido)"* → "Plataforma le
> muestra casos existentes o borradores en progreso que coinciden" → el
> usuario puede ir al caso/borrador encontrado, **o** "Crear nuevo caso de
> todos modos".
>
> Esto **reemplaza** el diseño de búsqueda de 3 campos (nombre, documento,
> teléfono) que traía el Cabezote de búsqueda del Excel/Guion Funcional de
> Drive — se elimina **nombre** como criterio de búsqueda y de coincidencia.
> El nombre sigue mostrándose en los resultados (el agente necesita ver a
> quién encontró), solo deja de ser un campo de entrada de la búsqueda.
>
> Motivo razonable aunque no explicitado por el cliente: la búsqueda por
> nombre es inherentemente difusa (variantes de escritura, nombres
> repetidos entre personas distintas) y puede producir coincidencias falsas
> que mezclen a dos víctimas distintas bajo el mismo Registro — documento y
> teléfono son identificadores mucho más precisos.
>
> El resto del diagrama (mostrar casos/borradores coincidentes, poder
> continuar el borrador encontrado, o crear uno nuevo "de todos modos" aun
> habiendo coincidencias) ya estaba cubierto por este plan — ver `flow-E04`
> (búsqueda), `flow-E07` (proceder al borrador) y `flow-E06` (nuevo registro,
> con el checkbox de confirmación agregado en la ronda del mockup HTML, que
> cumple el mismo rol que "Crear nuevo caso de todos modos").
>
> **GAP G-V3-15** (no resuelto por el diagrama): qué hace el agente cuando
> no tiene ni documento ni teléfono de la víctima a mano — el diagrama no
> contempla ese caso límite. Mientras se confirma, el formulario no ofrece
> ninguna alternativa de búsqueda por nombre, ni siquiera como opción
> secundaria.

### D8 — Cómo persistir un Borrador parcial en `victim_case`/`victim_case_form2` sin violar NOT NULL (mecanismo ya implementado y probado en V2-real)

> El usuario pidió revisar si V2-real ya había resuelto esto. **Sí** — está
> implementado y en producción, no solo planeado. Fuente: código real
> (`src/salvia/service/victim_case_form_service.go:96` `UpdateCaseDraft`;
> `src/internal/repository/victim_case_form_repository.go:133` `CreateDraft`/
> `:207` `UpdateDraftCoreFields`/`:263` `UpsertForm2FromAnswers`) y su doc
> `DocsMD/Screens/Registro de Caso V2/Flujos/flow-E03-cuando-se-guarda-primera-seccion.md`.
> Este mecanismo resuelve exactamente el problema de D1 (mantener
> `victim_case`/`victim_case_form2` como tabla única, con decenas de
> columnas NOT NULL, mientras el agente todavía no ha respondido casi nada).

**Filosofía del mecanismo: "nunca posponer, usar valores dummy."** La
versión original de V2 posponía crear `victim_case` hasta que la Sección 1
estuviera completa, y capturaba el error de Postgres `SQLSTATE 23502` (NOT
NULL violation) como señal de "todavía incompleto, reintentar". Se
abandonó ese diseño por instrucción explícita del negocio — ahora la fila
se crea/actualiza en **cada** `saveSection`, sin importar cuál sección se
guardó primero, rellenando con un valor dummy cualquier columna que
todavía no tenga respuesta real.

**Tabla de valores dummy por tipo de columna** (de `UpsertForm2FromAnswers`,
`internal/repository/victim_case_form_repository.go:290-358` — aplica
igual de forma literal a las columnas nuevas de V3, ya que siguen el mismo
sistema `Kind*` de V2-real):

| Kind | Valor dummy si aún no hay respuesta |
|---|---|
| `KindText` | `""` (cadena vacía) |
| `KindInt` | `0` |
| `KindDate` / `KindTimestamp` | fecha actual (`nullableOrNow`) |
| `KindTime` | se extrae solo la hora del input, sin dummy especial |
| `KindBoolEnum` | el i_code de **"No"** (categoría `yes_no`) |
| `KindEnum1` | también el i_code de **"No"** — es un hack pragmático: cualquier fila válida de `victim_case_form2_enums` satisface la FK/NOT NULL, no importa que semánticamente no sea "la opción No" de esa pregunta en particular |
| `KindEnumN` | no necesita dummy — no tiene columna propia (va a la tabla de relación), simplemente se omite si no hay respuesta |
| `KindScale` | `0` |

Los valores reales sobrescriben el dummy en cuanto el agente responde esa
pregunta — en cualquier sección posterior, no hace falta volver a la
sección original.

**Cómo se dispara:** un hook genérico `OnSectionUpdate(formID, submissionID, actorID)`
corre después de **cada** `saveSection`, de cualquier formulario del
sistema, y hace dispatch por `formID` a una función específica
(`UpdateCaseDraft` para V2). Un error ahí se loggea pero **no aborta**
`saveSection` — la sección ya quedó guardada; el draft del Caso se
sincroniza en la siguiente llamada.

**Otros detalles reutilizables para V3:**
- `victim_case.victim_case_form_submission_id` (columna propia, no tabla de
  relación aparte) permite buscar "¿ya existe un `victim_case` para este
  submission?" sin ambigüedad — mismo patrón a replicar si V3 sigue sobre
  `dinamic-form`.
- `FieldKey` indexado por posición (`"S{sección}Q{orden}"`), no por texto de
  la pregunta — evita que preguntas con texto idéntico entre ramas (ej.
  Tamizaje pareja/no-pareja) se confundan. V3 debería adoptar el mismo
  patrón de indexación para sus propias columnas de Tamizaje nuevas.
- `wasPartner`/`riskScore`/`riskLevel` (el equivalente de V2 al tipo/nivel
  de riesgo de V3) se recalculan en **cada** `saveSection`, con lo que haya
  disponible hasta ese momento — igual a como V3 ya lo decidió para
  `tamizaje_tipo_riesgo`/`tamizaje_nivel_riesgo` en `flow-E11`.

**Lo que V2 crea además de `victim_case`, a confirmar si V3 también lo
necesita:** `createDraftShell` también crea un `general_user_profile` +
`general_user` + `rel_role_general_user` (rol `"us"`) con login/password
aleatorios para la víctima, guardados en texto plano en
`victim_case_new_user_credentials` (jsonb) — presumiblemente para que la
víctima pueda eventualmente acceder a algo como "Mi Salvia". **No se asume
que V3 necesita esto mismo** — ninguna fuente de V3 (Excel, Guion
Funcional, Observaciones, mockup) menciona credenciales de usuario para la
víctima. Se deja como GAP G-V3-16 en vez de replicarlo sin confirmar.

**Gap que V2-real dejó abierto y que V3 hereda igual:** no existe un job de
expiración para Casos abandonados en estado Borrador indefinidamente — sigue
sin resolver en el código real, y esta planeación tampoco lo resuelve.

### D9 — Remisión, Entidades y Barreras: se reutiliza el modelo real de Seguimiento, no el del Excel

> El usuario pidió revisar en BD (no solo en el MD, que estaba
> desactualizado) cómo el formulario de Seguimiento real
> (`2d0aeb46-1af3-4c47-a0d5-c5bfc4d549ff`, tablas `form_section`/`question`/
> `option`/`visibility_condition`/`repeater_group`, consultadas directo
> contra Supabase con las credenciales de solo lectura de
> `src/config/db_config.json`) resuelve remisión a equipos, entidades y
> barreras. El MD `DocsMD/Screens/hacer-seguimiento/form-seguimiento-data.md`
> estaba desactualizado: le faltaban **2 secciones completas**
> ("Identificación de Entidades" y "Seguimiento a Entidades", que
> reemplazaron a lo que el MD llamaba "Rutas Diferenciales"), y el total
> real es 7 secciones, no 6.

**Decisión:** Registro V3 reutiliza el modelo real (6 equipos de remisión
con criterios propios, repeater de Entidades con dos variantes
Acudidas/Activación, repeater de Barreras) en vez del modelo simplificado
de 4 preguntas que traía el Excel de Drive para GS. Aplica el mismo
criterio que D8: cuando V2-real/Seguimiento ya resolvió algo en producción,
se reutiliza en vez de reinventar un modelo propio de V3.

**Diferencias deliberadas frente a copiar Seguimiento tal cual:**
- Seguimiento dispara "Identificación de Barreras" desde una pregunta de
  su sección "Valoración del Riesgo" (que no existe en Registro V3 de la
  misma forma). Se adaptó a una pregunta propia dentro de GS (secciones.md
  § GS, pregunta 17) — GAP G-V3-19, sin confirmar con el cliente.
- No se incluye el mecanismo de *reasignación completa del Caso a otro
  equipo* (Seguimiento sección 5, preguntas 14-15) — reasignar presupone
  que el Caso ya tiene un equipo asignado, lo cual no aplica en el primer
  Registro — GAP G-V3-18.
- El botón de Alerta Inminente del header (flow-E13) se mantiene; se agrega
  una entrada adicional dentro de GS para que quede documentada como parte
  formal de la gestión del caso — no es parte del formulario real de
  Seguimiento, es una decisión de este plan.

**Arquitectura de almacenamiento — importante, no verificada del todo:**
Los repeaters (`repeater_group`) de `dinamic-form` **no se proyectan sobre
columnas de `victim_case_form2`** — por su naturaleza de cantidad variable
de entradas, se guardan con el mecanismo genérico de respuestas de
`dinamic-form` (submission/answer), el mismo que ya usa Seguimiento para
sus repeaters de Entidades y Barreras. Esto significa que **Registrar
Entidades y Registrar Barreras (D y C de GS) no necesitan columnas nuevas
en `victim_case_form2`** — solo necesitan sus propias filas de
`question`/`repeater_group`/`visibility_condition` en el seed de V3, igual
que Seguimiento. **No se inspeccionó en este análisis la tabla exacta
donde `dinamic-form` guarda las respuestas de un repeater** (no es
`victim_case_form2`, pero no se confirmó cuál es) — GAP G-V3-21, confirmar
antes de implementar los repeaters de V3.

Las preguntas de remisión (B) que SÍ son de tipo escalar/multiselect
simple (¿requiere remisión?, ¿cuáles equipos?, criterios por equipo,
justificación) sí siguen el patrón `victim_case_form_fields.go` normal —
ver columnas nuevas abajo.

### D10 — Decisiones finales confirmadas por el usuario (imagen de 7 cajas) — pivote de arquitectura

> **Supersede partes de D1-D9 relativas al motor de formulario.** El resto
> de D1-D9 (modelo de datos, tríada, tamizaje, búsqueda, remisión/entidades/
> barreras) se mantiene — lo único que cambia es que ya no se construye
> sobre `dinamic-form`.

| Ref. imagen | Decisión | Afecta |
|---|---|---|
| A1 | **No se usa el componente `dinamic-form`** para Registro V3 — motivo del usuario: no maneja los porcentajes que necesita Tamizaje/progreso, es más fácil construirlo sin ese motor. Registro V3 se implementa como pantalla estática + componentes Vue a medida, **mismo patrón ya usado en este repo** para `set_victim_case.html`, `set_victim_contact.html`, `casos-component.js`, `gestion-reporte-modal.js` — no uno nuevo. Esto reemplaza toda referencia a `<dinamic-form>`, `form`/`form_section`/`question`/`option`/`visibility_condition`/`repeater_group` en `registro-v3-interface.md` — el árbol de interfaz sigue siendo válido como estructura/orden de campos, pero el `// components/dinamic-form.js` de la raíz del formulario debe leerse como "componente Vue a medida nuevo", no como el motor genérico. |
| A3 | Los catálogos de opciones (enums) se guardan **donde convenga** — decisión delegada a mí. Se reutiliza `salvia.victim_case_form2_enums` (la misma tabla que ya usan `set_victim_case.html`/`set_victim_contact.html`), cargada a un objeto `formState.enums.*` en el frontend — es el patrón ya establecido en este repo para pantallas estáticas, no uno nuevo. G-V3-06 (verificar categorías contra BD) queda delegado a mí, se resuelve durante la implementación, no antes. |
| B4 | G-V3-08 resuelto: se implementan las 5 preguntas propuestas (aumento de violencia física 60 días, uso de armas, amenaza de muerte, asfixia/estrangulamiento, acceso a armas) como las obligatorias iniciales del Tamizaje. Queda sujeto a validación del equipo Salvia en pruebas — no es definitivo, es implementable. |
| B5 | G-V3-20 resuelto activamente: el catálogo real de opciones de Barreras por sector se lee directo de la BD de Seguimiento antes de implementar el repeater — ver extracción abajo. |
| B6 | G-V3-02 resuelto: si se marcan ambos subtipos "Tentativa Feminicidio" y "Tentativa Homicidio abyecto o fútil" a la vez, **Feminicidio tiene prioridad**. |
| C7 | G-V3-13 resuelto: Alerta Inminente **no es un botón de acceso rápido en el header** — es una pregunta dentro de la sección Gestión Salvia (ya existía como SubsecciónE, lo que cambia es que se **elimina el botón del header** de `flow-E13`/`interface.md`, deja de ser un punto de entrada dual). Esto también resuelve G-V3-09 por consecuencia: al vivir solo dentro de GS, que ya requiere tríada completa para habilitarse, Alerta Inminente **no puede** saltarse la habilitación secuencial — nunca es alcanzable antes de ese punto. |
| C9 | G-V3-11 resuelto: Tipo de contacto (Habeas data) **sí añade opciones para Familiares (feminicidio) y Hombres**, con el mismo patrón/estilo que las 3 opciones existentes (Víctima-Caso Nuevo / Víctima-Seguimiento / Tercero-Reporte) — ahora son 5 opciones. |
| D11 | G-V3-16 resuelto: **sí** se crea un `general_user` con credenciales para la víctima al crear el Borrador, igual que V2-real (`createDraftShell`). |

---

## Columnas reales en `victim_case_form2` y `victim_case` (BD, 2026-10-31)

`victim_case_form2` tiene hoy 183 columnas. Las propias de V3 (los nombres de la planeación —`tz_v3_general_N`, `tamizaje_tipo_riesgo`…— **no** se usaron):

| Columna(s) | Tipo | Pregunta |
|---|---|---|
| `victim_case_form2_tzv3_g1` … `_g7` | integer (FK `yes_no`) | TZ_1 … TZ_7 |
| `victim_case_form2_tzv3_vpp` | integer (FK `yes_no`) | TZ_8, pregunta de prejuicio |
| `victim_case_form2_tzv3_p1` … `_p24` | integer (FK `yes_no`) | Rama pareja (TZP_1-24, la 24 es ideación suicida) |
| `victim_case_form2_tzv3_np1` … `_np24` | integer (FK `yes_no`) | Rama no pareja (TZNP_1-24) |
| `victim_case_form2_tzv3_risk_type` · `_risk_level` | integer (FK enum) | Resultado: tipo `fe`/`ha` y nivel `ba`/`mo`/`al`/`ex` (categorías `victim_case_form2_tzv3_risk_type` / `_risk_level`) |
| `victim_case_form2_tzv3_risk_score` · `_risk_tie_flag` | integer | Puntaje y marca de empate de subtipos |
| `victim_case_form2_victim_pronoun` · `_victim_identity_name` · `_victim_email` | integer FK / varchar | CI_1, CI_3, DP_50 |
| `victim_case_form2_allows_easy_report` · `_interinstitutional_actions_note` · `_saliva_management_explanation` | bigint / varchar / text | GS_1, GS_3, GS_4 (la última reutiliza la columna de V2) |
| `victim_case_form2_requires_referral` · `_referral_justification` | integer / varchar | GS_5, GS_13 |
| `victim_case_form2_reported_to_entity` · `_salvia_activates_routes` · `_agrees_route_activation` | integer (FK `yes_no`) | GS_14, GS_16, GS_15 |
| `victim_case_form2_identifies_barriers` | integer (FK `yes_no`) | GS_17 |
| `victim_case_form2_wants_imminent_alert` · `_imminent_alert_note` | integer / varchar | GS_18, GS_19 (solo se guardan; nadie más las lee) |
| `victim_case_form2_aggressors_occupation` | bigint NOT NULL | AG_4 (categoría propia, corrige el G-12 de V2) |

Las preguntas multi-valor (`multiple`, p. ej. CI_8, DP_3, HV_9, GS_6) no tienen columna: van a `rel_victim_case_form2_enums_victim_case_form2`, acotadas por la categoría de enum de cada campo (`victimCaseV3EnumNCategories`). Las columnas V2 de Tamizaje (38) quedan intactas.

En `victim_case` (campos núcleo y control): `victim_case_victim_names` / `_last_names` (de CI_2, partido en el primer espacio), `_doc_type` (code corto de CI_5), `_doc_number` (CI_6), `_town_code` (DP_28, que hoy recibe el texto digitado), `agent_id` (agente que creó el Borrador), `victim_case_new_user_credentials` (jsonb con login/clave de la víctima) y `victim_case_v3_answered_keys` (text, CSV de las keys respondidas alguna vez; base de `triadaCompleta` y de la validación de Finalizar).

Entidades y barreras no usan estas tablas: `salvia.entity_case` + `entity_case_obligation` y `salvia.barrier_v2` (ver [`registro-v3-pendientes.md`](registro-v3-pendientes.md) § B).

**Conteos reales de enums (BD, 2026-10-31)** donde difieren de la tabla siguiente: `adjustments_gbv` 12 (plan: 11) · `victim_doc_type` 20 (19) · `scenario_violence` 31 (36) · `scope_of_violence` 13 (14) · `relationship_with_presumed_aggressor_01` 20 + `_02` 19 (18) · `aggressors_occupation` 16 (19). Coinciden: `specially_protected_population` 17, `type_violence_experienced` 7, `recurrence_aggression` 3, `referral_team` 6, `emergency_measure` 6, `psycho_criteria` 7, `men_criteria` 1, `disability_service` 2, `stabiliz_criteria` 3, `interinst_action` 5, `sector` 4, `activation_channel` 3, y las barreras (`health` 12, `justice` 18, `protect` 13, `institutional` 5, `economic` 2, `territorial` 1, `differential` 12, `management` 6).

---

## Categorías de enum (`victim_case_form2_enums_category`) — nuevas o a ampliar

> ⚠️ Esta tabla asume qué opciones existen hoy según los MDs de V2-real
> leídos en la sesión anterior — **ninguna fue verificada contra las filas
> reales de `victim_case_form2_enums`** antes de escribir este plan (mismo
> GAP ya señalado en su momento para V2-real, G-03). Confirmar contra BD
> antes de escribir el seed SQL de V3.

| Categoría | Acción | Detalle |
|---|---|---|
| Tipo de documento (`DOCUMENT_TYPE_FORM2`, mapa Go, no enum DB) | Verificar/ampliar | Confirmar que los 19 tipos del Excel están en el mapa Go; agregar los que falten (Certificado RAMV, PEP-FF, Visa de Refugiado son los más propensos a faltar) |
| `adjustments_gbv` (ajustes razonables) | Ampliar a 11 opciones | Ver lista en `registro-v3-secciones.md` § CI |
| `nationality` | Agregar valor | **Apátrida** |
| `migration_condition` | Reemplazar/ampliar a 11 opciones | |
| `gender_identity` | Alinear a 7 opciones | Mujer / Mujer Trans / Hombre / Hombre Trans / No binaria-asignado masculino / No binaria-asignado femenino / Otro |
| `sexual_orientation` | Alinear a 7 opciones | Separa Gay/Lesbiana (antes posiblemente "Homosexual" unificado — verificar) |
| `specially_protected_population` (enumN) | Ampliar a 17 opciones | Lista específica colombiana — ver secciones.md § DP |
| `ethnic_affiliation` | Alinear a 7 opciones | |
| `indigenous_people` | Verificar | Lista cerrada de pueblos indígenas de Colombia — confirmar si ya existe en algún catálogo del sistema (Feminicidio/EAIF podría ya tenerlo) |
| `marital_status` | Alinear a 7 opciones | Agrega "No informa" |
| `last_education_level` | Alinear a 9 opciones | |
| `occupation` (víctima) | Ampliar a 24 opciones | |
| `income_generation_method` | Alinear a 8 opciones | |
| `employment_relationship` | Alinear a 6 opciones | |
| `asp_mode` / `reason_asp` | Alinear a 4 / 6 opciones | |
| `housing_tenancy_form` | Alinear a 7 opciones | |
| `housing_stratum` | Alinear a 8 opciones | "Sin estrato" ahora distinto de "No aplica/No sabe" |
| `has_dependents` (enumN) | Alinear a 8 opciones | |
| `type_violence_experienced` (enumN) | Agregar valores | **Política, Reproductiva, Vicaria** |
| Subtipos dinámicos (`violenceSubtypes`) | Agregar valores | **Tentativa Feminicidio, Tentativa Homicidio abyecto o fútil** (bajo Física) + reconciliar el resto de subtipos contra la redacción exacta del Excel |
| `scope_of_violence` (enumN) | Ampliar a 14 opciones | + regla de UI "máximo 2 selecciones" (validación de formulario, no de enum) |
| `workplace_sector_occurrence` | Alinear a 4 opciones | |
| `recurrence_aggression` | **Reemplazar** de 4 a 3 opciones | Verificar que ningún Caso histórico dependa del 4º valor antes de quitarlo de las opciones activas |
| `num_agressors` | Alinear a 5 opciones | Agrega "No se puede determinar/no informa" |
| `proximity_principal_aggressor` | **Reemplazar contenido** | Nuevo set fusionado "conoce y convive" (4 opciones) — ya no es solo proximidad |
| `relationship_with_presumed_aggressor` | Ampliar a 18 opciones | |
| Ocupación del agresor | **Nueva categoría dedicada** | Corrige GAP G-12 de V2-real (hoy reutiliza por error `relationship_with_presumed_aggressor_02`) — seedear con los 19 valores reales del Excel bajo una categoría propia |
| `aggressor_gender_identity` | Alinear a 6 opciones | |
| `action_plan` (gestión general, GS pregunta 2) | **Reemplazar contenido** | Orientación y enrutamiento / Gestión administrativa mediante llamada / Activación de ruta interinstitucional / Articulación interinstitucional / Escalamiento a organismos de control — alineado con Seguimiento real (D9), reemplaza el set anterior |
| Equipos de remisión (GS pregunta 6) | **Nueva categoría**, resuelto por D9 | Medidas de Emergencia / Atención Psicosocial / Estabilización / Atención Hombres / Salvia Dignidad / Discapacidad — 6 equipos reales, no 4 |
| Medidas de emergencia (GS pregunta 7) | **Nueva categoría** | Alojamiento / Transporte / Alimentación / Vestuario / Apoyo psicosocial / Otras ME |
| Criterios — Atención Psicosocial (GS pregunta 9) | **Nueva categoría, resuelve G-V3-04** | 1 obligatorio (mujeres/personas no binarias ≥18, VBG/VpP, voluntariedad, post-emergencia) + 6 opcionales con puntaje (conducta suicida=3, interseccionalidad=2, sin ruta=1, condiciones territoriales=1, sin acceso psico=1, naturalización VBG=1) — umbral: obligatorio + mínimo 3 puntos |
| Criterios — Atención Hombres (GS pregunta 10) | **Nueva categoría** | 1 solo valor, obligatorio: "Hombres diversos victimarios o posibles victimarios" |
| Servicios — Discapacidad (GS pregunta 11) | **Nueva categoría** | Apoyo LSC / Apoyo con enfoque de discapacidad (ley de capacidad legal) |
| Criterios — Estabilización (GS pregunta 12) | **Nueva categoría** | 3 valores, mínimo 1: sin ingresos / vivienda temporal / personas a cargo |
| Sector de entidad/barrera (GS preguntas de Entidades y Barreras) | **Nueva categoría compartida** | Salud / Justicia / Protección / Otras instituciones |
| Canal de activación (repeater Entidades) | **Nueva categoría** | Notificación (tarea) / Llamada (tarea) / Mi Salvia — Gestión institucional externa (tarea) |
| Catálogo de entidades por sector (`entidadesDisponibles`) | Dinámico, no enum estático | Mismo mecanismo que Seguimiento — confirmar fuente exacta antes de implementar |
| Acciones institucionales completadas/pendientes (`obligacionesCatalogo`) | Dinámico, no enum estático | Mismo mecanismo que Seguimiento |
| Barreras por sector (Salud/Justicia/Protección) y categorías transversales | **Catálogo pendiente de extraer — GAP G-V3-20** | Existen como filas reales de `option` en BD, no se extrajeron en este análisis |
| `support_contact_kinship` | Ampliar a 9 opciones | Incluye "Pareja o expareja (solo si no representa riesgo)" — la condición entre paréntesis es una regla de negocio, no solo texto de opción; ver GAP G-V3-05 |
| `scenario_violence` | Ampliar a 36 opciones | Lista completa en secciones.md § HV |

---

## GAPS

| # | Descripción | Impacto |
|---|---|---|
| G-V3-01 | Separación Registro/Caso en BD quedó deferida (D1) — el mecanismo de fan-out al editar Perfil no tiene control de concurrencia entre agentes | Medio — bajo volumen de edición simultánea del mismo documento, riesgo real pero no bloqueante |
| G-V3-02 | Sin regla de prioridad si se marcan ambos subtipos "Tentativa Feminicidio" y "Tentativa Homicidio abyecto o fútil" a la vez | Bajo — caso infrecuente, se propuso default en D3 |
| G-V3-03 | Cómo leer/mostrar Casos "mixtos" (Tamizaje diligenciado parcialmente en V2, continuado en V3) una vez V3 esté en producción | Medio — depende de si se permite continuar Casos Borrador antiguos con el formulario nuevo |
| G-V3-04 | Catálogo de "criterios que justifican la remisión" por cada proceso especializado (Psicosocial/Medidas de emergencia/Estabilización/Atención a hombres) no está en ninguna fuente leída — el Excel solo dice "según cada caso" | Medio — bloquea implementar GS pregunta 7 tal cual; se puede lanzar con un campo de texto libre mientras se define el catálogo |
| G-V3-05 | "Pareja o expareja (solo si no representa riesgo)" como opción de parentesco del contacto de apoyo tiene una condición de negocio embebida en el texto — no está claro si debe ser una validación activa (bloquear la opción si hay riesgo marcado en Tamizaje) o solo una nota para el agente | Bajo |
| G-V3-06 | Ninguna categoría de enum fue verificada contra las filas reales de `victim_case_form2_enums` — confirmar antes de escribir el seed SQL | Alto — bloquea implementación, no planeación |
| G-V3-07 | No se confirmó si existe ya un catálogo de "pueblos indígenas de Colombia" reutilizable de otro módulo (p. ej. Feminicidio/EAIF) | Bajo |
| G-V3-08 | Cuáles son exactamente las 5 preguntas obligatorias iniciales del Tamizaje (observación #15 del cliente) — propuesta no confirmada en D4 | Alto — bloquea implementar correctamente la lógica de revelado del Tamizaje |
| G-V3-09 | Cómo debe la Alerta inminente interactuar con la habilitación secuencial de TZ/GS — el cliente mismo lo marca como pendiente (observación #2) | Medio |
| G-V3-10 | Confirmar que el catálogo geográfico ya usado por V2-real (`departments`, `geo.*.towns`) es DIVIPOLA oficial, o si hay que migrarlo (observación #10) | Medio — no bloquea esta ronda, sí producción |
| G-V3-11 | Si "Tipo de contacto" en Habeas data necesita opciones adicionales para enrutar a Formulario Familiares (feminicidio) y Formulario Hombres (observación #11) | Medio |
| G-V3-12 | Si las preguntas de GS (acciones de gestión) o botones aparte del sidebar son los que conectan con Rutas/Barreras/Notificaciones (observación #17) — pregunta del propio cliente, sin responder | Medio |
| G-V3-13 | Dónde ubicar exactamente el botón de Alerta inminente (observación #16) — pregunta del propio cliente, sin responder | Bajo |
| G-V3-14 | Si el subtipo "Tentativa Feminicidio"/"Tentativa Homicidio abyecto o fútil" debe habilitar el módulo de Tamizaje de inmediato sin esperar tríada completa — propuesta del mockup HTML (prioridad más baja), relacionada con G-V3-09 pero no confirmada por el cliente | Medio |
| G-V3-15 | Qué alternativa de búsqueda existe si el agente no tiene ni documento ni teléfono de la víctima — el diagrama del cliente que eliminó la búsqueda por nombre (D7) no contempla este caso límite | Medio |
| G-V3-16 | Si V3 necesita crear un `general_user` con credenciales para la víctima al crear el Borrador, igual que V2-real — ninguna fuente de V3 lo menciona, no se replica sin confirmar (D8) | Medio |
| G-V3-17 | No existe job de expiración para Casos abandonados en Borrador — gap heredado de V2-real, sigue sin resolver (D8) | Bajo |
| G-V3-18 | Mecanismo de reasignación completa del Caso a otro equipo (distinto de remisión) que sí existe en Seguimiento real — no incluido en Registro V3 por no tener aún equipo asignado en el primer contacto (D9) | Bajo |
| G-V3-19 | Pregunta "¿Identificó barreras institucionales durante esta atención?" (GS pregunta 17) es una adaptación propia, sin equivalente en ninguna fuente de V3 — en Seguimiento el disparador es una pregunta de Valoración del Riesgo que V3 no tiene (D9) | Medio |
| G-V3-20 | Catálogo exacto de opciones de barreras por sector (Salud/Justicia/Protección), categorías transversales y "Gestión de la barrera" — existen como filas reales de `option` en BD pero no se extrajeron en este análisis (D9) | Medio — bloquea implementar el repeater de Barreras tal cual |
| G-V3-21 | Tabla real donde `dinamic-form` persiste las respuestas de un `repeater_group` (confirmado que NO es `victim_case_form2`, pero no se identificó cuál es) — bloquea implementar los repeaters de Entidades y Barreras de V3 (D9) | Alto — bloquea implementación, no planeación |

---

## Estado de los GAPS (2026-10-31)

| GAP | Estado | Cómo quedó |
|---|---|---|
| G-V3-01 concurrencia del fan-out | Abierto | Sin control entre agentes |
| G-V3-02 ambos subtipos de tentativa | ✅ Resuelto | Gana Feminicidio y se guarda `tie_flag` |
| G-V3-03 Casos mixtos V2/V3 | Abierto | V3 no lee ni migra Casos de V2 |
| G-V3-04 criterios de remisión | ✅ Resuelto | Catálogo en `registro-v3-campos.js` (puntajes y obligatorios) |
| G-V3-05 parentesco "Pareja o expareja" | Abierto | Sin validación activa |
| G-V3-06 verificar enums contra BD | ✅ Resuelto | Se leen de `victim_case_form2_enums` por `GET /victim-cases-v3/enums` |
| G-V3-07 pueblos indígenas | ✅ Resuelto | DP_13 usa `victim_case_form2_indigenous_people` |
| G-V3-08 5 iniciales del Tamizaje | ✅ Resuelto (propuesta) | TZ_1, 2, 3, 4, 7; sujeto a validación del equipo |
| G-V3-09 / G-V3-13 Alerta inminente | ✅ Resuelto | Solo GS_18/GS_19; sin botón en el header |
| G-V3-10 DIVIPOLA | Abierto | Departamentos de DP_27/HV_4 son una lista fija de nombres; DP_28/HV_5 (municipio) es texto libre y DP_28 se guarda en `victim_case_victim_town_code` |
| G-V3-11 tipos de contacto | 🟡 Parcial | 5 opciones, pero 4 sin destino (solo aviso) |
| G-V3-12 Rutas/Barreras/Notificaciones | Abierto | No implementado (E-16) |
| G-V3-14 tentativa habilita TZ | Abierto | No implementado |
| G-V3-15 búsqueda sin documento ni teléfono | Abierto | La búsqueda exige uno de los dos |
| G-V3-16 `general_user` de la víctima | ✅ Resuelto | Se crea al crear el Borrador, con credenciales en texto plano |
| G-V3-17 expiración de Borradores | Abierto | Sin job |
| G-V3-18 reasignación de equipo | Abierto | No incluido |
| G-V3-19 pregunta GS_17 | ✅ Resuelto | Radio "¿Se identificaron barreras?" que habilita "+ Agregar barrera" |
| G-V3-20 catálogo de barreras | ✅ Resuelto | Leído de BD (categorías `victim_case_form2_barrier_*`) |
| G-V3-21 dónde se guardan los repeaters | ✅ Resuelto | No hay repeaters: `entity_case` y `barrier_v2` |
| **G-V3-22** eliminar una entidad no elimina sus barreras (el `confirm` dice que sí) | Abierto | Las barreras quedan en `barrier_v2` y siguen contando en el Cierre |
| **G-V3-23** con las 5 iniciales del Tamizaje en "No" no hay forma de guardar el Tamizaje y Finalizar lo rechaza | Abierto | Por lectura de código; no verificado en navegador |
| **G-V3-24** al guardar GS/TZ el cliente deja `triadaCompleta = false` (el backend responde false al no evaluarla) y reaparecen los 🔒 | Abierto | Por lectura de código; se corrige al reabrir el Caso |
| **G-V3-25** `enviarSeccion` no revisa `r.ok`: un 500 muestra "Sección guardada." | Abierto | Por lectura de código |
| **G-V3-26** el botón "Finalizar" no guarda GS antes de validar | Abierto | Hay que presionar "Guardar sección" primero |

