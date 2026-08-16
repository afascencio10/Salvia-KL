# QA — Respuestas del Formulario en el tab Resumen

Casos de prueba para la nueva sección "Respuestas del Formulario" del tab Resumen de `/salvia/seguimiento/:id` (ver [detalle-seguimiento-interface.md](../detalle-seguimiento-interface.md)): lee, por `formSubmissionId`, 6 respuestas puntuales del formulario de seguimiento (`loadFormAnswers` en `followup_v2_service.go`) y las muestra en texto legible, sin pasar por `dinamic-form`.

---

## Casos de prueba

| Caso | Descripción | Condición inicial | Resultado esperado | Método de verificación |
|---|---|---|---|---|
| C1 | Campos poblados se leen y muestran correctamente | `formSubmissionId` con `riskAnalysis`, `caseManagement` y 3 entries de `barrierIdentified` respondidas (multi-select "Gestión de la barrera" con las 6 opciones) | El backend arma `formAnswers` con el texto exacto de cada pregunta y resuelve las 6 opciones de gestión a sus labels legibles (incluyendo `orientacion_llamada`, que no está en el mapa de labels ya existente en `form_service.go`) | Request directo al endpoint (`GET /api/v1/cases/follow-ups/detail`) + query a `salvia.answer`/`salvia.repeater_entry` |
| C2 | Campos/listas vacías no se muestran | Mismo `formSubmissionId` de C1 — `referralEvidence` sin responder y `barrierFollowUps` (repeater Seguimiento a Barreras) sin entries | `formAnswers.referralEvidence = ""` y `formAnswers.barrierFollowUps = null` — el frontend oculta ambos bloques (`v-if`) sin afectar los campos poblados de C1 | Request directo al endpoint + Playwright (DOM) |
| C3 | El DOM renderiza lo que devuelve el backend | Mismo seguimiento — tab Resumen (activo por defecto) | La sección "Respuestas del Formulario" es visible; el análisis de riesgo y la gestión del seguimiento muestran el texto exacto; las 3 entradas de "Gestión — Identificación de Barreras" muestran el texto de las 6 opciones resueltas; "Gestión — Seguimiento a Barreras" y "Elementos que evidencian la remisión" no están en el DOM | Playwright (DOM) |

> C1 y C2 se verificaron a nivel de API/BD durante el desarrollo (no repetidas por Playwright — ver Fase 2). C3 es lo único que realmente necesita navegador: confirmar que el binding Vue (`v-if`/`v-for`/interpolación) refleja fielmente el JSON ya probado.

---

## Fase 2 — Backend (verificado durante implementación, 2026-08-14)

Contra el seguimiento de prueba `fd053f17-5614-4206-8c9f-1a61012eacc8` (caso E2E compartido `019eadd1-6df3-7728-b014-6db076cbbb1e`), `GET /api/v1/cases/follow-ups/detail?id=fd053f17-5614-4206-8c9f-1a61012eacc8` devolvió:

```json
"formAnswers": {
  "riskAnalysis": "Se identificaron factores de riesgo que requieren seguimiento continuo. La situación de riesgo se mantiene moderada.",
  "caseManagement": "Se realizó gestión integral con seguimiento a barreras institucionales identificadas en los sectores salud, justicia y protección.",
  "referralEvidence": "",
  "barrierFollowUps": null,
  "barrierIdentified": [
    { "gestion": "Orientación y enrutamiento - Llamada, Gestión administrativa - Llamada, Activación de ruta interinstitucional, Articulación institucional, Escalamiento a organismo de control, Alerta por barreras" },
    { "gestion": "Orientación y enrutamiento - Llamada, Gestión administrativa - Llamada, Activación de ruta interinstitucional, Articulación institucional, Escalamiento a organismo de control, Alerta por barreras" },
    { "gestion": "Orientación y enrutamiento - Llamada, Gestión administrativa - Llamada, Activación de ruta interinstitucional, Articulación institucional, Escalamiento a organismo de control, Alerta por barreras" }
  ]
}
```

Confirma C1 (texto exacto + las 6 labels resueltas, incluyendo `orientacion_llamada`) y C2 (`referralEvidence` vacío, `barrierFollowUps` nulo por no tener entries) al nivel de datos.

---

## Fase 3 — Frontend (Playwright, headless)

Test agregado en `qa-salvia/tests/hacer-seguimiento/seguimiento-detalle.spec.ts` — mismo archivo que ya cubre `/salvia/seguimiento/:id` para el tab Formulario, siguiendo su convención (sin page object dedicado, `page.goto` directo + `waitForURL` post-login).

```
CI=1 npx playwright test tests/hacer-seguimiento/seguimiento-detalle.spec.ts -g "Respuestas del Formulario"

✓ /salvia/seguimiento/:id tab Resumen muestra "Respuestas del Formulario" (7.4s)
1 passed
```

Corrido headless (`CI=1`, sin abrir navegador) contra el backend local en :9090 con los cambios de esta sesión. Login vía el fixture estándar del repo (`auth.fixture.ts`, credenciales `TEST_USER`/`TEST_PASSWORD` de `.env.local`).

---

## Cobertura

| Flujo de prueba | Casos cubiertos |
|---|---|
| Request directo a `GET /api/v1/cases/follow-ups/detail` + query BD (Fase 2, durante desarrollo) | C1, C2 |
| `seguimiento-detalle.spec.ts` › `tab Resumen muestra "Respuestas del Formulario"` (Playwright, headless) | C3 |

**Resultado:** 3/3 casos cubiertos y verdes.
