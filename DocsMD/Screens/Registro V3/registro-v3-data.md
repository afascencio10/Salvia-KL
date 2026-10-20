# Registro V3 — Decisiones de arquitectura y mapeo a Base de Datos

> Complementa `registro-v3-secciones.md` (catálogo de preguntas). Este
> documento resuelve **dónde cae cada cosa en BD** y las 3 decisiones de
> alcance confirmadas con el usuario antes de planear.

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

## Columnas nuevas necesarias en `victim_case_form2`

Todas las columnas existentes de V2-real (`victim_case_form_fields.go`) se
reutilizan donde el concepto coincide — ver la tabla de categorías de enum
más abajo para qué columnas cambian de contenido sin cambiar de esquema. Las
siguientes son **columnas que no existen hoy** y hay que agregar:

| Columna propuesta | Kind | Pregunta que alimenta |
|---|---|---|
| `victim_case_form2_vpp_prejudice_aggression` | `bool_enum` | TZ pregunta fuerte 8 (define RHAF) |
| `victim_case_form2_tz_v3_general_1` … `_7` | `bool_enum` × 7 | TZ preguntas generales 1-7 |
| `victim_case_form2_tz_v3_pareja_1` … `_24` | `bool_enum` × 24 | TZ rama pareja (incl. ideación suicida) |
| `victim_case_form2_tz_v3_nopareja_1` … `_24` | `bool_enum` × 24 | TZ rama no pareja (incl. ideación suicida) |
| `victim_case_form2_tamizaje_tipo_riesgo` | `enum1` | Resultado calculado: Feminicidio / Homicidio abyecto o fútil (persistido para alertas, KPIs y traspaso a Seguimiento) |
| `victim_case_form2_tamizaje_nivel_riesgo` | `enum1` | Resultado calculado: Bajo / Moderado / Alto / Extremo |
| `victim_case_form2_requires_referral` | `bool_enum` | GS pregunta 5 (¿requiere remisión a equipos SALVIA?) |
| `victim_case_form2_referral_justification` | `text` | GS pregunta 13 (elementos que evidencian la remisión) |
| `victim_case_form2_wants_imminent_alert` | `bool_enum` | GS pregunta 18 — marca si se activó Alerta Inminente desde GS (además del log de `flow-E13`) |
| `victim_case_form2_identifies_barriers` | `bool_enum` | GS pregunta 17 (¿identificó barreras? — GAP G-V3-19) |
| `victim_case_form2_reported_to_entity` | `bool_enum` | GS pregunta 14 (¿ya puso en conocimiento a alguna entidad?) |
| `victim_case_form2_salvia_activates_routes` | `bool_enum` | GS pregunta 16 (¿Salvia debe activar rutas?) |
| `victim_case_form2_agrees_route_activation` | `enum1` | GS pregunta 15 (¿de acuerdo con la activación?) |

> `¿Cuáles equipos?` (GS pregunta 6) y todos los bloques de criterios por
> equipo (preguntas 7, 9, 10, 11, 12) son `multiple` — no necesitan columna
> propia, van por la tabla de relación genérica (ver nota debajo). El
> banner de advertencia (6b) y el de info (8) no se persisten, son
> puramente de UI.

> Las columnas de Tamizaje V3 **no reemplazan** las 38 columnas booleanas del
> Tamizaje de V2-real (`aggressor_violence_physical_increase`, etc.) — se
> dejan intactas para no afectar Casos ya cerrados con el formulario V2. Los
> Casos nuevos creados con V3 simplemente no vuelven a escribir en esas
> columnas viejas. Ver GAP G-V3-03 sobre cómo deben leerse Casos mixtos
> (abiertos en V2, continuados en V3) una vez V3 esté en producción.

Las preguntas de tipo `multiple` (enumN) **no necesitan columna propia** —
usan la tabla de relación genérica `rel_victim_case_form2_enums_victim_case_form2`
ya existente; solo necesitan una categoría de enum nueva o ampliada (ver
abajo).

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
