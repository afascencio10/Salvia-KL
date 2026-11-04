# seed_test_data

Fase 2 del seed del sandbox cloud: datos que se aplican **después** de que la
app arranca, porque dependen de tablas que crea la app (GORM `AutoMigrate`:
`case_task`, `follow_up_v2`, `dupla`, `articulating_mechanism`...), y datos
leídos de Supabase (BD de pruebas) o sintéticos para probar requerimientos
concretos: casos, reportes, víctimas, seguimientos, barreras, remisiones,
tareas, respuestas de formularios, etc.

**Solo datos.** Esta carpeta guarda `*.sql` con datos y este README; ni scripts
de QA ni nada que no sea data para el sandbox (ver "Alcance de esta carpeta" en
[`../README.md`](../README.md)).

**Es permanente.** Todo lo que una sesión agrega aquí se commitea y lo
encuentran todas las sesiones siguientes. Antes de leer de Supabase, revisa
si el dato ya está en esta carpeta.

- Un archivo por requerimiento: `AAAAMMDD-<slug>.sql`.
- `apply-test-seeds.sh` aplica los `*.sql` en orden alfabético, cada uno en
  una sola transacción, cuando la app ya respondió (lo lanza solo
  `./run_dev.sh`). Lleva la cuenta en `public.cloud_seed_applied`: aplica solo
  los archivos nuevos o modificados. Si un archivo falla, se revierte
  completo, se muestra el error de Postgres y se sigue con el siguiente.
- Cada archivo debe ser **idempotente** (`ON CONFLICT DO NOTHING`),
  **autocontenido** (incluye sus filas padre) y **sin `BEGIN`/`COMMIT`**.
- Puede asumir el esquema, `seed_reference_data.sql`, los usuarios
  `cloud.test.*` y las tablas que ya creó la app.
- Solo lectura sobre Supabase. No se copian credenciales (hashes, tokens);
  las FK hacia usuarios se remapean a los `cloud.test.*`.
- Trae lo mínimo que el requerimiento necesita, no la tabla completa.
- Si los datos son **sintéticos** (no salen de Supabase), dilo en la cabecera y
  hazlos reconocibles: ids fijos `5eed…`, radicados `SEED-…`, nombres «Semilla».
  Las fechas relativas ("hace N días") mantienen el ejemplo vigente.

Hoy trae:

| Archivo | Qué trae | Para qué |
|---|---|---|
| `20261101-base-mecanismos-y-duplas.sql` | Catálogo real de mecanismos articuladores (6), el mecanismo del usuario `cloud.test.mc` y la dupla `cloud.test.ps` + `cloud.test.ts` | Que las pantallas de esos roles tengan su configuración |
| `20261101-entidades-bogota.sql` | 7 sedes reales de Bogotá (2 por sector: salud, justicia, protección; 1 de otras instituciones) y las 65 obligaciones de sus entidades | Que las pantallas que buscan sedes por municipio y sector tengan opciones reales |
| `20261102-reportes-wp1.sql` | **Sintético.** 17 reportes con radicado `SEED-xxxx` (4 tipos; anónimos y con sesión `op`/`et`; sin asignar, asignado a `cloud.test.op`, con intentos y contacto logrado, borrador iniciado, cerrado válido y cerrados inválidos por cada motivo, incluido el esquema 3x3 completo y un estado heredado), 17 intentos, 49 hitos del historial y 2 casos de Registro V3 (borrador y terminado) con sus usuarios `us` | Listado de Reportes, Mis Reportes (`cloud.test.op`), modal de gestión con historial, consulta pública (`/salvia/public-consulta`) y «Registro de caso» sin radicar y gestionar reportes a mano. Qué radicado sirve para qué: cabecera del archivo |

Procedimiento, plantilla del archivo y cómo aplicar uno nuevo en un sandbox
que ya está corriendo: [`../actualizar-seed-con-composio.md`](../actualizar-seed-con-composio.md).
