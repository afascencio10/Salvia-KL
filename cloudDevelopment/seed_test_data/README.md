# seed_test_data

Datos leídos de Supabase (BD de pruebas) para probar requerimientos concretos
en el sandbox cloud: casos, víctimas, seguimientos, barreras, remisiones,
tareas, respuestas de formularios, etc.

**Es permanente.** Todo lo que una sesión agrega aquí se commitea y lo
encuentran todas las sesiones siguientes. Antes de leer de Supabase, revisa
si el dato ya está en esta carpeta.

- Un archivo por requerimiento: `AAAAMMDD-<slug>.sql`.
- `setup-local-postgres.sh` aplica los `*.sql` en orden alfabético, cada uno
  en una sola transacción, después de `seed_reference_data.sql` y
  `seed_local_dev.sql`. Si un archivo falla, se omite con un aviso y el
  arranque continúa.
- Cada archivo debe ser **idempotente** (`ON CONFLICT DO NOTHING`),
  **autocontenido** (incluye sus filas padre) y **sin `BEGIN`/`COMMIT`**.
- Solo lectura sobre Supabase. No se copian credenciales (hashes, tokens);
  las FK hacia usuarios se remapean a los `cloud.test.*`.
- Trae lo mínimo que el requerimiento necesita, no la tabla completa.

Procedimiento, plantilla del archivo y cómo aplicar uno nuevo en un sandbox
que ya está corriendo: [`../actualizar-seed-con-composio.md`](../actualizar-seed-con-composio.md).
