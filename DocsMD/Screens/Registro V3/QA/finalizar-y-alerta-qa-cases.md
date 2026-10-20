# Alerta Inminente y Finalizar — Casos de prueba

Cubre E13 (ya no es botón de header, decisión C7 — es la pregunta GS_18/GS_19) y E15 (Finalizar). Endpoint de finalizar: `POST /api/v1/victim-cases-v3/:icode/finalizar`.

## Alerta Inminente (GS_18/GS_19)

| Caso | Descripción | Condición inicial | Resultado esperado | Método de verificación |
|---|---|---|---|---|
| A1 | Confirma que no existe ningún botón de "Alerta inminente" en el header/sidebar | Cualquier pantalla del Formulario | Solo existe la pregunta dentro de GS § E — ningún otro punto de entrada (C7) | Validar DOM en todas las pantallas del Formulario |
| A2 | GS_18 = "No" | — | GS_19 (justificación) permanece oculto | Validar DOM |
| A3 | GS_18 = "Sí" | GS_19 oculto | Aparece el textarea de justificación | Validar DOM |
| A4 | Guarda GS con GS_18="Sí" y GS_19 vacío | — | Confirmar si el backend exige contenido en GS_19 o lo acepta vacío (no hay validación explícita en `resolveFieldUpdates` para esto) | Request directo + leer BD |
| A5 | Guarda GS con GS_18="Sí" y una justificación de texto | — | `victim_case_form2_wants_imminent_alert` = id de "Sí"; `victim_case_form2_imminent_alert_note` con el texto exacto | Leer BD |

## Finalizar (E15) — sin validación de obligatorias (G2)

| Caso | Descripción | Condición inicial | Resultado esperado | Método de verificación |
|---|---|---|---|---|
| F1 | Finaliza un Caso que todavía está en `Borrador` (nunca pasó por `Iniciado`) | Caso recién creado, sin ningún `SaveSection` | `POST /finalizar` → 200; `victim_case.status` pasa directo de `bo` a `te`, sin pasar por `in` (no hay validación de tríada ni de obligatorias — comportamiento actual, G2) | Request directo + leer BD |
| F2 | Finaliza un Caso con la tríada completa (`in`) pero TZ y GS nunca tocados | Tras completar tríada (T4 de `triada-y-guardado-qa-cases.md`) | `POST /finalizar` → 200; status pasa a `te` igual, sin exigir TZ/GS completos | Request directo + leer BD |
| F3 | Finaliza dos veces el mismo Caso | Caso ya en `te` (tras F1 o F2) | Confirmar qué hace `UpdateStatus` al llamarse sobre un Caso ya `te` — no debería fallar (es un `UPDATE` incondicional), pero confirmar que no genera ningún efecto secundario raro | Request directo + leer BD |
| F4 | Finaliza un Caso inexistente (icode inválido) | — | Confirmar el código de respuesta (probablemente 500 "error al finalizar" en vez de 404, ya que `UpdateStatus` no valida existencia) | Request directo |

## Nota

F1/F2/F3 documentan el comportamiento **actual** de `Finalize` (sin ninguna validación), no el descrito en el MD original de planeación (que sí exigía completitud). Si el usuario decide implementar la validación real de obligatorias [1er contacto] + [máx. 2º contacto] de los 7 módulos antes de permitir Finalizar, estos casos hay que reescribirlos — por ahora reflejan el código tal como está.
