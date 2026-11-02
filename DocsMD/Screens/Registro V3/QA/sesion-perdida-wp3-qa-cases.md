# QA — Registro V3 · «Finalizar registro» responde «no autenticado» (WP-3 · REG-21)

Reproducción de las causas posibles y verificación de lo que ve la persona. Eventos: [`flow-E15`](../Flujos/flow-E15-cuando-finaliza-y-enruta-a-seguimiento.md) y [`flow-E21`](../Flujos/flow-E21-cuando-un-guardado-o-una-carga-falla.md).

> **Estado: ejecutado el 2026-11-02** — `scripts/e2e/wp3-reg21-sesion.js`, **7 de 7 ✅**, más el experimento S6 (manual).

## Qué significa cada mensaje (código: `requireRegistroV3Access`)

| Respuesta | Qué pasó |
|---|---|
| 401 `{"error":"no autenticado"}` | La **cookie del navegador no trae `userData`** (se cerró sesión en otra pestaña o se borraron las cookies) |
| 401 `{"error":"sesión expirada"}` | La cookie sí lo trae, pero el servidor ya no tiene esa sesión (el mismo usuario inició sesión en otro navegador: el login borra la sesión previa) |
| 401 sin cuerpo | La sesión es válida pero el rol no tiene `manage_registro_v3` |

## Experimentos y resultado

| Caso | Hipótesis | Qué se hizo (con el formulario abierto) | Resultado |
|---|---|---|---|
| S1 | Línea base | Sesión vigente, Finalizar | ✅ Cierra el registro (`te`): **Finalizar funciona con sesión vigente** |
| S2 | H1 | Cerrar sesión en **otra pestaña** del mismo navegador y pulsar Finalizar | 401 `no autenticado` — **reproduce el síntoma del cliente**. Ahora la pantalla avisa y conserva los datos |
| S3 | H1 / H3 | Borrar las cookies del sitio (equivale a cookie vencida) | 401 `no autenticado`, mismo aviso |
| S7 | Sesión única por usuario | Iniciar sesión con el **mismo usuario** en otro navegador | 401 `sesión expirada`, mismo aviso |
| S6 | H4 | Reiniciar el servidor conservando la cookie | 200: la sesión **sobrevive** al reinicio (`sessions` se recarga). No es la causa en este entorno |
| H2 | Carrera del middleware que guarda la cookie en cada petición | Revisión de código | **No reproducida**: exigiría una respuesta tardía de una petición sin sesión; no se cambia `MainRouter.go` |

**Conclusión.** «No autenticado» no era un defecto de «Finalizar»: la cookie perdió la sesión mientras el formulario estaba abierto (cierre de sesión en otra pestaña o cookie borrada). El defecto real era que la pantalla mostraba el texto crudo del servidor, sin decir qué hacer, y que `finalizarRegistro` leía el cuerpo como JSON (un 401 sin cuerpo no mostraba nada). Ahora los 401 (con o sin cuerpo, `no autenticado` o `sesión expirada`) muestran «tu sesión venció o no tienes permiso. Inicia sesión de nuevo en otra pestaña y vuelve a pulsar…: lo que diligenciaste sigue en esta pantalla» y no se recarga ni se redirige.

## Cobertura

| Flujo de prueba | Casos cubiertos |
|---|---|
| `wp3-reg21-sesion.js` | S1, S2, S2b, S3, S3b, S7, S7b |
| Manual (`stop-app.sh` + `run_dev.sh`) | S6 |
| `wp3-guardado-cierre.js` (U3b/c, U4c/d, G6, G7) | Mensaje y datos intactos con 401 con y sin cuerpo |
