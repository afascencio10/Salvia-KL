# QA — Registro V3 · Tríada, bloqueo, guardado y cierre (WP-3 · REG-35, REG-13, REG-37, REG-21, REG-05, REG-36)

Casos de [`flow-E12`](../Flujos/flow-E12-cuando-guarda-una-seccion.md), [`flow-E10`](../Flujos/flow-E10-cuando-cambia-de-modulo.md) y [`flow-E15`](../Flujos/flow-E15-cuando-finaliza-y-enruta-a-seguimiento.md). El manejo de errores de guardado lo implementó WP-2 (REG-11, [`flow-E21`](../Flujos/flow-E21-cuando-un-guardado-o-una-carga-falla.md), casos en [`E21-…`](E21-guardado-o-carga-fallida-qa-cases.md)): aquí REG-36 solo se **verifica**.

> **Estado: ejecutado el 2026-11-02** — `scripts/e2e/wp3-guardado-cierre.js`, **29 de 29 ✅**.

**Método.** Backend primero (API con la sesión del navegador + lectura de la BD), después pantalla solo para lo que la API no cubre (Playwright; las fallas se simulan con `page.route`). Casos creados con `POST /victim-cases-v3/draft` y borrados al terminar.

## Casos de prueba

| Caso | REG | Descripción | Condición inicial | Resultado esperado | Método |
|---|---|---|---|---|---|
| G1 | 35 | Guardar GS con la tríada incompleta | Borrador sin CI/DP/HV/AG | 200 `triadaCompleta:false`; estado `bo` | API + BD |
| G5 | 35 | CA no aporta a la tríada | Solo CA guardado | `false`; sigue `bo` | API |
| G4a / G4 | 35 | Regresión: AG completa la tríada | CI, DP y HV guardados | Antes de AG `false`; con AG `true` y estado `bo → in` | API + BD |
| G2 | 35 | Guardar GS con la tríada completa | Caso `in` | `triadaCompleta:true`; el estado sigue `in` | API + BD |
| G3 | 35 | `/seccion/TZ` con la tríada completa | Caso `in` | `triadaCompleta:true` | API |
| N4 | 13 | Finalizar sin haber guardado GS | Todo lo demás completo | 400 con `modulosIncompletos = [GS]` (GS sigue siendo obligatorio) | API |
| C6 | 37 | Contenido de `camposFaltantes` | Mismo caso | Ordenado y coherente con los módulos: GS_1, GS_14, GS_16, GS_17, GS_18, GS_2, GS_4, GS_5 | API |
| C0 | 37 | Con GS guardado, finalizar | GS guardado | 200 y estado `te` | API + BD |
| G6 | 21 | Guardar sin cookie de sesión | Sin login | 401 `{"error":"no autenticado"}` | API |
| G7 | 21 | Rol sin permiso (`ro`) | Sesión de `ro` | 401 **sin cuerpo** | API |
| U1a–d | 13 | GS sin candado en un Borrador; TZ con candado; GS abre con CI vacío sin errores; aviso solo de Tamizaje | Borrador nuevo | Sin 🔒 en GS, 🔒 en TZ, sin errores de consola ni «undefined», aviso que nombra solo Tamizaje | Pantalla |
| U2a | 05 | Recordatorio de faltantes en rojo | GS_4 vacío en GS | Mismo texto, color `rgb(176, 42, 55)`, contraste 6,50:1 (≥ 4,5) | Pantalla (color calculado) |
| U2b | 05 | «Sección guardada.» conserva su estilo | GS_4 corregido | Sin rojo | Pantalla |
| U2c / U2d | 35 | Guardar GS con la tríada completa | Caso `in` | `triadaCompleta` sigue `true`, sin candados, aviso verde permanece; GS (último módulo) no cambia de módulo | Pantalla |
| U3a | 36 | El servidor responde 500 | `/seccion/GS` → 500 | Aviso rojo «…error del servidor…» y nunca «Sección guardada.» | Pantalla |
| U3b / U3c | 21 | 401 con y sin cuerpo al guardar | `/seccion/GS` → 401 | «…tu sesión venció… lo que escribiste sigue en esta pantalla»; datos intactos; sin excepciones | Pantalla |
| U4a | 37 | GS § A incompleto al finalizar | GS_4 vacío | Recordatorio junto al botón y **0** peticiones a `/finalizar` | Pantalla (conteo de peticiones) |
| U4b | 37 | Faltan GS_5, 14, 16, 17, 18 | § A completo, el resto sin responder | El aviso dice «Gestión Salvia, falta: …» con sus etiquetas; el Caso sigue `in` | Pantalla + BD |
| U4c / U4d | 21 | 401 con y sin cuerpo al finalizar | `/finalizar` → 401 | Aviso claro, sigue en el formulario, datos intactos, sin excepciones | Pantalla |
| U4e | 37 | Finalizar sin pulsar «Guardar sección» antes | GS escrito, sin guardar | Cierre; estado `te`; GS_4 guardado | Pantalla + BD |
| U4f | 37 | Doble clic en Finalizar | — | Una sola petición a `/finalizar` | Pantalla |
| U4g | 06 | El cierre y la pestaña no muestran «V3» | Cierre | Sin «V3» | Pantalla |

## Cobertura

| Flujo de prueba | Casos cubiertos |
|---|---|
| `wp3-guardado-cierre.js` — fase API | G1, G2, G3, G4a, G4, G5, N4, C6, C0, G6, G7 |
| `wp3-guardado-cierre.js` — fase pantalla | U1a–d, U2a–d, U3a–c, U4a–g |

## Resultados de ejecución

**Fecha:** 2026-11-02 · **Build:** rama `claude/bold-thompson-jbze8u` sobre `develop` @ 74933b1 · **Resultado:** 29/29 ✅. Regresión de lo que comparte pantalla con WP-2: `wp2-reg11` 26/26 (tras volver a abrir Hechos después de guardar, por REG-01), `wp2-reg12` 9/9 y `wp2-reg34` 7/7.
