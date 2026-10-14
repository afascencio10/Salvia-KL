# Listado Reportes / Mis Reportes / Gestión Reporte Modal — Tablas de BD

Documento de BD **compartido** por las tres piezas de "Gestión de reportes" (FigJam, sección homónima): la pantalla `Listado Reportes`, la pantalla `Mis Reportes` y el componente `gestion-reporte-modal`. Las tres leen/escriben las mismas tablas, así que se documentan juntas aquí en lugar de triplicar el contenido.

> Fuente de la verdad: esquema real verificado en la BD de pruebas (Supabase, ver `DocsMD/General/tables-database.md`) el 2026-10-10, no solo el script `db_postgres_creation.sql` (varias columnas reales como `agent_id`, `victim_case_team`, `follow_up_v2` no están en ese script porque llegaron por migraciones `.sql` sueltas — están gitignored en este repo, igual que la migración de Registro V2).

---

## Backend existente a reutilizar — búsqueda y filtros de `victim_contact`

Hallazgo de la auditoría del 2026-10-11 (contra la imagen de "Funcionalidades del componente" de FigJam): `/salvia/primer_contacto` **ya tiene** backend de filtrado probado en producción, que el endpoint nuevo `GET /api/v1/victim-contacts/list` debe adaptar en vez de reinventar:

| Pieza existente | Archivo | Qué resuelve |
|---|---|---|
| `VictimContactListFilters` (struct: `Names`, `LastNames`, `Phone`, `Document`) | `VictimContactDAO.go:22-27` | Forma de los 4 filtros de búsqueda |
| `buildVictimContactsWithoutCaseWhere(...)` | `VictimContactDAO.go:316-372` | WHERE real: nombres/apellidos en columnas directas, **teléfono en `victim_contact_form2` O `victim_contact_form1`** (comentario en el código documenta un bug ya corregido — antes solo miraba Form2), **documento solo en `victim_contact_form1`** (Form2/tercero no tiene esa columna) |

El endpoint nuevo combina estos 4 campos con **OR** en un solo parámetro `search` (los 4 inputs separados de hoy se combinan con AND) — ver [flow-E07 de Listado Reportes](./Flujos/flow-E07-cuando-busca-por-texto.md).

---

## Tabla existente que se reutiliza sin cambios — `victim_contact`

Ya tiene todo lo necesario para identificar y listar un reporte. Columnas relevantes hoy:

| Columna | Tipo | Null | Uso en Gestión de reportes |
|---|---|---|---|
| `victim_contact_i_code` | `character varying(36)` | No | Identificador del reporte en todas las URLs/endpoints nuevos |
| `victim_contact_names` / `victim_contact_last_names` | `character varying(32)` | Sí | Columna "Víctima / ID Reporte" |
| `victim_contact_status` | `character varying(1)` | No | `'v'` vigente / `'i'` invalidado — mecanismo **legado**, ver nota de convivencia abajo |
| `victim_contact_status_description` | `character varying(256)` | Sí | Motivo de invalidación legado (texto libre) |
| `victim_contact_facts_description` | `text` | Sí | "Descripción de los hechos" en el popup de detalle |
| `victim_contact_general_user` | `character varying(36)` | **No** | Quién/qué registró el contacto (no es el agente gestor — ver nota) |

**Nota — `victim_contact_general_user` no es el campo de asignación.** Es `NOT NULL` en todas las filas existentes porque guarda quién originó el registro (backend genérico de creación). Reutilizarlo como "agente asignado" rompería la semántica para las filas ya existentes. Por eso se propone una columna nueva dedicada (`victim_contact_agent_id`, ver abajo) — exactamente el mismo patrón que ya existe en `victim_case`: ahí tampoco se reutilizó ningún campo de auditoría, se usa una columna propia (`agent_id`).

**Nota — convivencia con el botón "Invalidar" legado de `/salvia/primer_contacto`.** Ese flujo ya existente sigue intacto y no se toca. Una fila con `victim_contact_status = 'i'` se interpreta en las pantallas nuevas como "Cerrado" de solo lectura, con el motivo mostrado como `victim_contact_status_description` (sin intentar mapearlo a los 6 motivos nuevos). No hay backfill — es una regla de lectura, no una migración de datos.

**Implementada (2026-10-13):** `reportStatus = cerrado_invalido`; el detalle devuelve `closeReasonCode = 'legado'` (código sintético, no existe como fila de enum) con `victim_contact_status_description` como motivo; los filtros de estado la respetan (un legado no es "asignado" ni "sin asignar") y `POST /intento` y `POST /cerrar` responden **409**. Afecta a 2 269 de 3 503 reportes (65 %) a esa fecha.

---

## Columnas nuevas en `victim_contact`

Mismo patrón y mismos tipos que las columnas equivalentes ya probadas en `victim_case` (verificadas en BD: `agent_id character varying(64)`, `victim_case_team character varying(64)`).

| Columna nueva | Tipo | Null | Equivalente en `victim_case` | Descripción |
|---|---|---|---|---|
| `victim_contact_agent_id` | `character varying(64)` | Sí | `agent_id` | `general_user_i_code` del agente integral (`op`) con el reporte asignado. `NULL` = sin asignar. |
| `victim_contact_team` | `character varying(64)` | Sí | `victim_case_team` | Equipo del agente asignado, copiado al reasignar (ver GAP-3). `NULL` mientras no haya asignación. Si el agente destino no tiene equipo queda `''` (cadena vacía), no `NULL`. |
| `victim_contact_close_reason` | `bigint` | Sí | — (nuevo patrón, ver abajo) | FK a `victim_case_form2_enums.victim_case_form2_enums_id`, categoría `victim_contact_close_reason`. `NULL` = reporte abierto. |
| `victim_contact_close_observation` | `character varying(500)` | Sí | — | Observación libre, solo para los motivos que la piden (`ncm`, `dec` — ver tabla de motivos). |
| `victim_contact_closed_date` | `timestamp without time zone` | Sí | — | Cuándo se cerró. `NULL` mientras esté abierto. |

**Por qué no hay una columna "estado" separada:** el estado del reporte (Sin asignar / Asignado / Cerrado) es 100% derivable de las columnas de arriba — no hace falta guardarlo aparte y arriesgar que quede desincronizado:

```
SI victim_contact_close_reason IS NOT NULL  → "Cerrado"  (categoría Válido/Inválido según el motivo, tabla abajo)
SI NO Y victim_contact_agent_id IS NOT NULL → "Asignado"
SI NO                                       → "Sin asignar"
```

### Motivos de cierre (`victim_case_form2_enums`, categoría nueva `victim_contact_close_reason`)

Mismo mecanismo ya usado para `anonimato`/`parentesco`/`report_type_v2` en Registro V2 — filas con `ON CONFLICT (...name) DO NOTHING`, no tocan ninguna categoría existente. Motivos tomados literal del mockup ya validado (`Mockup_Reportes_Registro.html`, función `cerrarPrevio`/`cerrarPorNoContacto`/`confirmarCierrePostContacto`/`confirmarApertura`):

| `code` | `name` | Texto en UI | Categoría del cierre | Pide observación |
|---|---|---|---|---|
| `ni` | `victim_contact_close_reason_no_info_minima` | Inválido por No información Mínima | Inválido | No |
| `nc` | `victim_contact_close_reason_no_contacto` | Inválido por Cierre por no contacto | Inválido | No |
| `co` | `victim_contact_close_reason_no_competencia` | Inválido por Cierre por no competencia | Inválido | **Sí** |
| `dp` | `victim_contact_close_reason_decision_persona` | Inválido por Cierre por decisión de la persona | Inválido | **Sí** |
| ~~`ap`~~ | ~~`victim_contact_close_reason_apertura_caso`~~ | Válido - Apertura de Caso Salvia | Válido | ⏸️ **Diferido** — fila de enum sin daño en dejarla migrada, pero ningún evento la escribe en esta fase (M-10 diferido, ver `gestion-reporte-modal-interface.md`). Desde 2026-10-13 el backend **rechaza `ap` con 400** (`deferredCloseReasonCodes`, decisión del usuario); quitarlo de ese mapa al construir M-10 |

> Códigos reales confirmados en BD (2 caracteres, no 3 — ver "Migración — ejecutada" abajo).

---

## Tabla nueva — `victim_contact_attempt`

No existe ningún mecanismo hoy para registrar intentos de contacto (confirmado: `grep` sobre todo `*.go` no encuentra nada tipo "attempt" ligado a `victim_contact`). Se modela igual que `follow_up_v2` (la tabla más reciente y ya probada del proyecto para "eventos puntuales ligados a una entidad por i_code", con soft-delete):

```sql
CREATE TABLE salvia.victim_contact_attempt (
    victim_contact_attempt_id SERIAL NOT NULL,
    victim_contact_attempt_i_code CHARACTER VARYING(36) NOT NULL,
    victim_contact_attempt_victim_contact CHARACTER VARYING(36) NOT NULL,  -- victim_contact_i_code (como case_id en follow_up_v2, no FK numérica)
    victim_contact_attempt_result CHARACTER VARYING(10) NOT NULL,          -- 'fallido' | 'exitoso'
    victim_contact_attempt_date TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    victim_contact_attempt_agent CHARACTER VARYING(64) NOT NULL,          -- agente que registró el intento
    victim_contact_attempt_creation_date TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    victim_contact_attempt_deleted_at TIMESTAMP WITHOUT TIME ZONE,
    CONSTRAINT victim_contact_attempt_pkey PRIMARY KEY (victim_contact_attempt_id),
    CONSTRAINT victim_contact_attempt_i_code_key UNIQUE (victim_contact_attempt_i_code)
);
```

### Regla del esquema 3x3 (verbatim del mockup ya validado, no inventada)

> **Corrección 2026-10-12** (respuesta a comentario de Felipe en `flow-M07`): la versión anterior de esta regla decía "3 días distintos con intentos", sin exigir cuántos intentos por día — eso permitía cerrar con solo 3 intentos totales. La fuente real del mockup (`diasCompletosDe()`) exige **3 por día**, no 1.

- "Día completo" = un día con **al menos 3** intentos fallidos (zona `America/Bogota`, igual que `casos_nuevos` en `casos-component`). Al llegar al 3° intento fallido del día, se avisa la cuota pero **no bloquea** seguir intentando ese mismo día — un 4° o 5° intento el mismo día no suma un segundo "día completo".
- El cierre por no-contacto se **habilita** (no se dispara solo) cuando hay **3 días completos distintos** — es decir, mínimo 9 intentos fallidos en total, repartidos en al menos 3 días con 3 cada uno. No importa si los días son consecutivos.
- Un intento `exitoso` interrumpe el esquema inmediatamente y pasa a las 2 opciones de post-contacto activas en esta fase (no competencia / decisión de la persona — "abrir caso" queda diferido, ver `gestion-reporte-modal-interface.md`).

> Nota de implementación: estos cálculos (intentos de hoy, días distintos completos) se hacen en el backend a partir de `victim_contact_attempt`, igual que `next_follow_up_date`/`completed_follow_ups_count` se calculan hoy en `cases_list_repository.go` — no se guardan como columnas derivadas.

> **Convención de zona horaria (verificada en la implementación, 2026-10-13):** `victim_contact_attempt_date` es `timestamp without time zone` y guarda **UTC** (`time.Now().UTC()`; con `time.Now()` a secas Postgres ignora el offset y la hora quedaba en la zona del servidor). Para agrupar por día de Bogotá: `(col AT TIME ZONE 'UTC' AT TIME ZONE 'America/Bogota')::date` — `col AT TIME ZONE 'America/Bogota'` a secas interpreta la hora como ya-Bogotá y desplaza el día. El log del modal también formatea con `timeZone: 'America/Bogota'`.

---

## Migración — ejecutada (2026-10-12)

`src/migrations/add_victim_contact_management.sql`. Notas de ejecución real:

- **`victim_contact` es owned por `postgres`** (no por `salvia_gorm` como las tablas creadas en este proyecto más recientemente) — mismo caso que `security.general_user` y, de hecho, el mismo que ya enfrentó `victim_case.agent_id`/`victim_case_team` en su momento (ver `tables-database.md` → "AutoMigrate para tablas legacy"). Los 5 `ALTER TABLE ADD COLUMN` + el `ADD CONSTRAINT` se corrieron con el superusuario documentado (`postgres.pwwelwfhauspznpatuqm`, puerto 6543) en vez de `salvia_gorm` — no se tocó `db_config.json` ni `migrateLegacyTables()` en `main.go`, porque esto es un cambio puntual de una vez, no un sync continuo.
- **Códigos de enum de 2 caracteres, no 3** — `victim_case_form2_enums_code` es `character varying(2)` (verificado en BD, no está en `db_postgres_creation.sql`). Los códigos reales quedaron: `ni` (no información mínima), `nc` (no contacto), `co` (no competencia), `dp` (decisión de la persona), `ap` (apertura de caso) — reemplazan a `nim`/`noc`/`ncm`/`dec` usados en versiones anteriores de este documento.
- `victim_contact_attempt` sí se creó como `salvia_gorm` sin problema (tabla nueva, no legacy) — mismo patrón que `victim_contact_form2_v2` en la migración de Registro V2.

### Código final de los motivos de cierre (corregido)

| Motivo | Code |
|---|---|
| No información mínima | `ni` |
| No contacto (3x3) | `nc` |
| No competencia | `co` |
| Decisión de la persona | `dp` |
| Apertura de caso (diferido) | `ap` |

## Tablas que alimentan el popup de detalle (sin cambios)

| Tabla | Campo(s) usados |
|---|---|
| `victim_contact_form1` | Teléfono de contacto (según flujo que originó el reporte) |
| `victim_contact_form2` | Datos de la persona afectada, ajuste razonable |
| `victim_contact_form2_v2` | Reportante (nombre/apellido/correo), anonimato, parentesco — **solo si el reporte vino de V2** |

Ver GAP-1 abajo sobre cómo el detalle debe funcionar para reportes de cualquier origen (V1, V2, o registrado por agente).

---

## Decisiones confirmadas por el usuario (2026-10-11)

| # | Pregunta | Decisión | Aplicado en |
|---|---|---|---|
| GAP-2 | ¿Quién puede registrar intentos y cerrar un reporte? | **B) `op` y `sv` pueden gestionar** cualquier reporte visible para ellos — no solo el agente con `victim_contact_agent_id` propio. `sv` deja de estar limitado a "solo ver + reasignar". | `gestion-reporte-modal-interface.md` (se retira el modo solo-lectura), M-01…M-10 |
| GAP-3 | ¿Bloqueo de equipo al reasignar reportes? | **A) Sin bloqueo** — cualquier agente integral activo es candidato; el reporte hereda el equipo del agente destino. | `reasignar-casos-modal-interface.md` (sección Extensión), `Listado Reportes` flow-E06 |
| GAP-4 | ¿`/salvia/primer_contacto` sigue viva en paralelo? | **B) Coexisten + enlace** desde `/salvia/primer_contacto` hacia `/salvia/reportes`. No se retira ni se modifica su lógica — solo se agrega un link de salida. | Ver nota al final de `listado-reportes-index.md` |

### Decisiones confirmadas el 2026-10-13 (tras la ejecución de QA)

| Tema | Decisión | Aplicado en |
|---|---|---|
| `POST /cerrar` con `ap` | Bloquear con 400 mientras M-10 esté diferido | `deferredCloseReasonCodes` (servicio de gestión) |
| `co`/`dp` sin intento exitoso previo | Dejar como está: el backend no exige contacto previo (solo la UI encadena los pasos) | — |
| `/salvia/primer_contacto` en blanco | Arreglar (`victimCaseRes = "[]"`) | `VictimContactFacade.go` |
| Menú lateral | Agregar "Listado Reportes" (op, sv) y "Mis Reportes" (op) | `sidebar.html` |
| Ruta de Mis Reportes | Dejar `/salvia/mis-reportes` | `MainRouter.go` |

## GAPS — aún pendientes de confirmar (no son decisión técnica, requieren info de negocio)

| # | Pregunta | Impacto si no se resuelve |
|---|---|---|
| GAP-1 | FigJam dice que el componente debe procesar reportes "de todas las funcionalidades donde hoy entren reportes (creo que son 2 o 3)". ¿Cuáles son exactamente esas 2-3 fuentes además de V1/V2 público? | El popup "Ver Reporte" podría quedar vacío para reportes de un origen no contemplado — el endpoint de detalle debe leer `form1`, `form2` o `form2_v2` según cuál exista, sin asumir una sola fuente. **Estado 2026-10-13:** cubierto por la implementación (2 951 de 3 503 reportes son solo-`form1`; 552 tienen `form2`); falta confirmar con negocio si existe un 4.º origen |
| GAP-6 | ¿"Remitir a Atención Especializada" (reportes Hombre/Familiar, step post-contacto) es un requisito real? El propio mockup lo marca como no confirmado con el cliente. | Ver `gestion-reporte-modal-interface.md` |
| GAP-7 | ¿Qué debe pasar si se abre caso desde un reporte Familiar/Hombre, dado que ese formulario destino no existe hoy en la aplicación real (ni en el mockup)? | Ver `gestion-reporte-modal/Flujos/flow-M10` |
