# migrationsProd — migraciones SQL que hay que correr en producción

Aquí están las migraciones SQL que **se aplican a mano en la base de datos de
producción**. Cada una **ya se corrió en Supabase** (proyecto
`pwwelwfhauspznpatuqm`), que es el ambiente **develop** (el de pruebas). La base de
producción es otra y **no tiene estos cambios hasta que alguien los corra**; la
aplicación no los hace sola.

> **Hay que correr en producción todo archivo de esta carpeta que figure como
> ⏳ PENDIENTE en la tabla de abajo**, en orden por nombre (la fecha va primero) y
> **antes** de publicar el código que lo usa.

## Estado

| Archivo | Qué cambia | Develop (Supabase) | Producción |
|---|---|---|---|
| [`20261102_wp2_registro_v3_catalogos_y_columnas.sql`](20261102_wp2_registro_v3_catalogos_y_columnas.sql) | WP-2 (REG-07, REG-08, REG-11): 9 filas en el catálogo `salvia.victim_case_form2_enums` (zona Rural/Urbano y 7 tipos de discapacidad), 3 columnas nuevas en `victim_case_form2` y 3 columnas de municipio ensanchadas a `varchar(100)` | ✅ 2026-11-02 | ⏳ PENDIENTE |

Al correr una en producción, cambia su celda de «Producción» por `✅ AAAA-MM-DD
(quién la corrió)` en el mismo commit.

## Cómo correr una migración en producción

1. **Respaldo.** Confirma que hay un respaldo reciente de la base de producción.
2. **Chequeos previos** (solo lectura): los de la migración, más abajo.
3. **Aplicar**, con un rol dueño de las tablas del esquema `salvia` (la migración
   hace `ALTER TABLE`):

   ```bash
   psql "$URL_BD_PRODUCCION" -v ON_ERROR_STOP=1 -f migrationsProd/20261102_wp2_registro_v3_catalogos_y_columnas.sql
   ```

   También sirve pegar el archivo completo en el editor SQL del proveedor. Cada
   archivo trae su propio `BEGIN; … COMMIT;`: si algo falla, no queda nada a medias.
   Y son **idempotentes**: correr una dos veces no duplica filas ni falla.
4. **Chequeos posteriores**, más abajo.
5. **Publicar el código** que usa la migración.
6. **Anotar** el resultado en la tabla de «Estado».

En Supabase (develop) la migración de WP-2 se corrió con el `execute_sql` del MCP de Supabase.
Otra anterior (Registro de Contacto Víctima V2) se corrió con el rol dueño `salvia_gorm`, según
`DocsMD/Screens/Registro de Contacto Víctima V2/related-tables.md`.

## `20261102_wp2_registro_v3_catalogos_y_columnas.sql`

### Qué hace

- **REG-07 · zona de residencia.** Agrega las opciones `victim_case_form2_residence_zone_ru`
  (`ru`, «Rural») y `…_ur` (`ur`, «Urbano») al catálogo. La lista «Zona de residencia»
  salía vacía.
- **REG-11 · hechos y residencia.** Agrega `victim_case_form2.victim_case_form2_facts_department`
  y `victim_case_form2.victim_case_form2_residence_department` (`varchar(100)`, nulas), y
  amplía a `varchar(100)` `victim_case_form2.victim_case_form2_facts_town_code`,
  `victim_case_form2.victim_case_form2_residence_town` y `victim_case.victim_case_victim_town_code`
  (medían 8 caracteres: guardar un municipio escrito a mano fallaba con un 500).
- **REG-08 · tipo de discapacidad.** Agrega las 7 opciones
  `victim_case_form2_disability_type_01` … `_07` (códigos `01` … `07`) y la columna
  `victim_case_form2.victim_case_form2_disability_type` (`bigint`, nula, FK al catálogo).
- Adelanta la secuencia del id del catálogo si quedó por detrás del máximo (nunca la retrocede).

### Por qué va antes del código

El código de Registro V3 de WP-2 escribe esas columnas y lee esas opciones. Sin la
migración, guardar las secciones de datos personales y hechos falla, y las listas de
zona de residencia y tipo de discapacidad salen vacías.

### Antes de correrla (solo lectura)

**A. Opciones del catálogo.** En una base sin la migración da **0 filas**:

```sql
SELECT victim_case_form2_enums_name AS name, victim_case_form2_enums_code AS code,
       victim_case_form2_enums_category AS categoria
FROM salvia.victim_case_form2_enums
WHERE victim_case_form2_enums_name ~ '^victim_case_form2_(residence_zone_(ru|ur)|disability_type_0[1-7])$'
ORDER BY categoria, code;
```

Si salen filas (alguien las creó a mano), la migración las respeta —busca por `name`, no
las duplica—, pero revisa que su `categoria` y su `code` sean los de la lista de arriba:
la aplicación busca por categoría y código.

**B. Vistas o reglas que dependan de las columnas que se ensanchan.** Debe dar **0 filas**;
si sale alguna, `ALTER COLUMN … TYPE` falla (y la transacción se deshace): hay que
eliminar la vista, correr la migración y volver a crearla.

```sql
SELECT DISTINCT v.relname AS vista_o_regla
FROM pg_depend d
JOIN pg_rewrite r ON d.objid = r.oid
JOIN pg_class v ON r.ev_class = v.oid
JOIN pg_class t ON d.refobjid = t.oid
JOIN pg_attribute a ON a.attrelid = d.refobjid AND a.attnum = d.refobjsubid
WHERE t.relnamespace = 'salvia'::regnamespace
  AND t.relname IN ('victim_case_form2', 'victim_case')
  AND a.attname IN ('victim_case_form2_facts_town_code', 'victim_case_form2_residence_town', 'victim_case_victim_town_code')
  AND v.oid <> t.oid;
```

**C. Políticas RLS** sobre las dos tablas (en develop hay 0; si en producción hay, que ninguna
mencione las columnas de municipio):

```sql
SELECT tablename, policyname FROM pg_policies
WHERE schemaname = 'salvia' AND tablename IN ('victim_case_form2', 'victim_case');
```

**Bloqueo.** Ampliar un `varchar` no reescribe la tabla, pero toma un bloqueo corto sobre
`salvia.victim_case_form2` y `salvia.victim_case`: córrela en un momento de poca actividad.

### Después de correrla (solo lectura)

**A.** La consulta A de arriba ahora da **9 filas**: 7 de `victim_case_form2_disability_type`
(`01` … `07`) y 2 de `victim_case_form2_residence_zone` (`ru`, `ur`).

**B. Columnas** — deben salir **6 filas**:

```sql
SELECT table_name, column_name, data_type, character_maximum_length AS largo
FROM information_schema.columns
WHERE table_schema = 'salvia'
  AND ((table_name = 'victim_case_form2' AND column_name IN (
          'victim_case_form2_facts_department', 'victim_case_form2_residence_department',
          'victim_case_form2_disability_type', 'victim_case_form2_facts_town_code',
          'victim_case_form2_residence_town'))
    OR (table_name = 'victim_case' AND column_name = 'victim_case_victim_town_code'))
ORDER BY table_name, column_name;
```

| table_name | column_name | data_type | largo |
|---|---|---|---|
| `victim_case` | `victim_case_victim_town_code` | character varying | 100 |
| `victim_case_form2` | `victim_case_form2_disability_type` | bigint | — |
| `victim_case_form2` | `victim_case_form2_facts_department` | character varying | 100 |
| `victim_case_form2` | `victim_case_form2_facts_town_code` | character varying | 100 |
| `victim_case_form2` | `victim_case_form2_residence_department` | character varying | 100 |
| `victim_case_form2` | `victim_case_form2_residence_town` | character varying | 100 |

**C. Secuencia del catálogo** — `ultimo_valor_secuencia` debe ser **mayor o igual** que `max_id`
(en develop dio 923 y 923):

```sql
SELECT (SELECT max(victim_case_form2_enums_id) FROM salvia.victim_case_form2_enums) AS max_id,
       s.last_value AS ultimo_valor_secuencia
FROM pg_sequences s
WHERE s.schemaname || '.' || s.sequencename
      = pg_get_serial_sequence('salvia.victim_case_form2_enums', 'victim_case_form2_enums_id');
```

### Al publicar

- **Registro V3 no depende del catálogo en memoria:** lo lee de la base en cada consulta, tanto
  para listar (`GET /api/v1/victim-cases-v3/enums`) como para guardar, así que las opciones nuevas
  se ven sin reiniciar. (Lo que la aplicación sí carga a memoria una sola vez al arrancar es el
  catálogo de los formularios anteriores: si algún día corres una migración de catálogo con la
  aplicación ya corriendo, reiníciala.) Los navegadores con Registro V3 abierto deben recargar la
  pantalla, porque guardan el catálogo hasta recargar.
- Los textos que ve la persona («Rural», «Urbano», los 7 tipos de discapacidad) viajan en el
  código (`src/salvia/config/Locale.go`), así que llegan con la publicación.

### Marcha atrás

No hay script de reversa. La migración solo agrega (filas de catálogo y columnas nulas) y
ensancha columnas; el código anterior no usa nada de eso, así que si hay que desistir antes de
publicar el código, dejarla aplicada no lo afecta. Si aun así se necesita deshacerla (con
respaldo; **no se probó**): `DROP COLUMN` de las 3 columnas nuevas y, después, borrar las 9
filas del catálogo. El ensanchado a `varchar(100)` no se revierte.

## Reglas de esta carpeta

1. Cuando una migración se aplica en Supabase (develop) y producción todavía no la tiene,
   se copia aquí **idéntica** a la de `src/migrations/` y se agrega una fila a la tabla de
   «Estado» con Producción ⏳ PENDIENTE. Para comprobar que la copia es idéntica:
   `cmp src/migrations/<archivo> migrationsProd/<archivo>`.
2. Nombre `AAAAMMDD_tema.sql`; se corren en producción por orden de nombre.
3. Idempotentes y con su propio `BEGIN; … COMMIT;`, como las de `src/migrations/`.
4. **No se edita** una migración después de correrla en producción: el cambio va en una
   migración nueva, con otra fecha.
5. No entran aquí los seeds y datos de prueba (`cloudDevelopment/`, solo del sandbox) ni lo que
   la aplicación crea sola al arrancar (ver abajo).
6. El `.gitignore` ignora `*.sql`; esta carpeta tiene su excepción (`!migrationsProd/*.sql`),
   así que un archivo nuevo se agrega con `git add` normal.

## Lo que no se corre a mano

**Cambios de esquema que la aplicación hace sola al arrancar** (`AutoMigrate`, en `src/main.go`):
no hay SQL que correr, pero conviene revisar el log de arranque, porque si un `AutoMigrate`
falla la aplicación solo lo anota como `[WARN]` y sigue arrancando. Lo de WP-1 va por ahí:
`victim_contact_radicado` (con el índice único `uq_victim_contact_radicado`) y
`victim_contact_reporter_role` en `salvia.victim_contact` (`migrateLegacyTables`), y
`CaseTimelineEvent` (el evento ya no exige `case_id`: puede colgar de un reporte), que está
en la lista de modelos de `main.go`. Esto se leyó del código; no se probó contra producción.

## Por confirmar: migraciones anteriores

No tengo acceso a producción, así que no sé cuáles de estas ya se corrieron allá. Quien lo sepa
debe confirmarlo; las que falten se agregan a esta carpeta con su fecha y su fila en «Estado».

| Migración | Dónde está | Qué se sabe |
|---|---|---|
| `src/migrations/add_followup_v2_calendar_fields.sql` (HU-027) | en el repo | Agrega columnas del calendario a `salvia.follow_up_v2`. `FollowUpV2` está en el `AutoMigrate` de `main.go`, que probablemente ya las cubre. |
| `src/migrations/alter_followup_v2_scheduled_date_to_timestamp.sql` | en el repo | Cambia `scheduled_date` de `DATE` a `TIMESTAMPTZ`. |
| `add_victim_contact_management.sql`, `add_victim_contact_v2.sql`, `add_registro_v3.sql` (partes 1 a 4) | **no están en el repo**: se hicieron cuando `*.sql` estaba ignorado y los DocsMD dicen «no está en git» | Sin el archivo no se pueden copiar aquí. Qué hacían: `DocsMD/Screens/Registro V3/changelogOct2026.md`, `DocsMD/Screens/Listado Reportes/related-tables.md` y `DocsMD/Screens/Registro de Contacto Víctima V2/related-tables.md`. |
