# seed_test_data

Fase 2 del seed del sandbox cloud: datos que se aplican **después** de que la
app arranca, porque dependen de tablas que crea la app (GORM `AutoMigrate`:
`case_task`, `follow_up_v2`, `dupla`, `articulating_mechanism`...), y datos
leídos de Supabase (BD de pruebas) para probar requerimientos concretos:
casos, víctimas, seguimientos, barreras, remisiones, tareas, respuestas de
formularios, etc.

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

Hoy trae `20261101-base-mecanismos-y-duplas.sql`: el catálogo real de
mecanismos articuladores, el mecanismo del usuario `cloud.test.mc` y la dupla
de prueba `cloud.test.ps` + `cloud.test.ts`.

Procedimiento, plantilla del archivo y cómo aplicar uno nuevo en un sandbox
que ya está corriendo: [`../actualizar-seed-con-composio.md`](../actualizar-seed-con-composio.md).
