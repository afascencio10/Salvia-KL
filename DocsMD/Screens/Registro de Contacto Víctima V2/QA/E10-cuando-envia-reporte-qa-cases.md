# QA — E-10 (Cuando presiona "Enviar Reporte") y navegación general del wizard

Ejecutado 2026-09-27, siguiendo `qa-process.md`. Cubre el único evento que persiste en backend (E-10) más verificación de la reactividad de UI de E-01 a E-09/E-11. **Revisado el 2026-10-31 contra el código:** la fila C6 no coincide con el comportamiento actual (ver la nota bajo la tabla de casos).

## Herramientas usadas

- **`playwright-cli`** (interactivo, contra `http://localhost:9090` con rebuild forzado) — para la UI: navegación entre pasos, renderizado condicional por tipo de reportante, y 3 flujos E2E completos (Tercero, Familiar, Hombre) hasta la pantalla de confirmación.
- **`curl` directo al backend** (`POST /salvia/public-v2`, con `User-Agent: flutter-client` para saltar el captcha en las pruebas que no lo estaban probando a propósito) — para los casos de validación puros, más rápidos de iterar que por UI.
- **`psql` contra la BD real** (Supabase) — para verificar que cada registro válido quedó persistido correctamente en las 3 tablas, y que los inválidos NO crearon nada.

> La sesión de 2026-09-27 no se guardó como `.spec.ts` (era exploratoria). Desde entonces `qa-salvia` tiene dos specs que sí recorren el wizard: `tests/registro-contacto-victima-v2/crear-reporte-victima.spec.ts` (C1, C3 y C16: radica un reporte tipo Víctima por la UI con el `User-Agent` `flutter-client` y lo borra con `resetVictimContactV2Fixture`) y `tests/gestion-reportes-e2e/ciclo-reporte-nuevo.spec.ts` (radica un Víctima y lo sigue hasta el cierre; ver [`../../Listado Reportes/QA/gestion-reportes-e2e-qa-cases.md`](../../Listado%20Reportes/QA/gestion-reportes-e2e-qa-cases.md)). Los tipos Tercero, Familiar y Hombre y la batería de validaciones siguen sin spec.

## Casos de prueba

| Caso | Descripción | Condición inicial | Resultado esperado | Método de verificación |
|---|---|---|---|---|
| C1 | Autorización = "Sí, autorizo" | Paso 'auth' | Botón "Continuar" se habilita | Playwright |
| C2 | Autorización = "No autorizo" | Paso 'auth' | Aviso de corte visible, "Continuar" deshabilitado | Playwright |
| C3 | Reporte tipo Víctima, flujo completo válido | Paso 'formulario', todos los campos requeridos llenos | 200, radicado devuelto, registro visible en `victim_contact`/`_form2`/`_form2_v2` | Playwright E2E + verificación en BD |
| C4 | Reporte tipo Tercero, anonimato = "Anónimo", sin nombre de reportante | Paso 'formulario' | 200 — el nombre del reportante NO es obligatorio cuando es anónimo | curl + verificación en BD |
| C5 | Reporte tipo Tercero, anonimato = "A nombre propio", flujo completo | Paso 'formulario' | 200, radicado, `reporter_names`/`reporter_last_names`/`reporter_email` persistidos | Playwright E2E + verificación en BD |
| C6 ⚠️ | Reporte tipo Familiar, disposición de recontacto = "Sí" | Paso 'formulario' | Aparece el campo "hora de contacto". **Resultado vigente:** al enviar con la hora digitada responde **400** `form2V2.reporterBestContactTime` (L1, ver abajo); con "Sí" y la hora vacía el botón queda deshabilitado. Lo que sí debe cumplirse una vez corregido L1: 200 y `victim_col_phone`/`best_contact_time` de la víctima en 0/vacío (no aplica) | Playwright E2E + verificación en BD |
| C7 | Reporte tipo Hombre, disposición de recontacto = "No" | Paso 'formulario' | Bloque "Datos de la persona afectada" NO se renderiza; al enviar, 200 y `victim_contact.names`/`last_names` quedan vacíos (decisión 7A) | Playwright E2E + verificación en BD |
| C8 | Captcha incorrecto, resto del formulario válido | POST directo | 400, error específico en `victimContact.captchaSolution`, no crea ningún registro | curl + verificación en BD (tabla sin la fila nueva) |
| C9 | `adjustmentsGBV` vacío (`[]`) | POST directo | 400, error "Por favor escoja una o más opciones" | curl |
| C10 | `reportType.icode` con un UUID que no existe en ninguna categoría de enum | POST directo | 400, error "Por favor escoja una opción" en `reportType` | curl |
| C11 | `factsDescription` de 5 caracteres ("corto") | POST directo | 400, error "Este campo es obligatorio" en `factsDescription` | curl |
| C12 | Tercero "a nombre propio" sin `reporterNames`/`reporterLastNames` | POST directo | 400, error en ambos campos | curl |
| C13 | `victimColPhone = 0` (tipo Tercero, donde sí aplica) | POST directo | 400, error "Digite un número válido" | curl |
| C14 | `authorizationAnswer` = código "n" (No autorizo) enviado directo al backend | POST directo | 400, error "Este campo es obligatorio" — el backend no confía solo en el guard del cliente | curl |
| C15 | Nombre con apóstrofe/ñ (`O'Brien Núñez`) y descripción con comillas dobles, comillas simples y emoji | POST directo | 200 — datos persistidos exactamente igual (round-trip de encoding verificado), tabla `victim_contact` intacta (sin riesgo de inyección SQL, queries parametrizadas) | curl + verificación en BD |
| C16 | "Volver al Inicio" tras una confirmación exitosa | Paso 'confirmacion' | Wizard vuelve a 'auth' con formulario en blanco; el popup de éxito compartido no bloquea el click (bug #6, corregido) | Playwright |
| C17 | Pantalla V1 (`/salvia/public/nuevo`) tras desplegar V2 | — | El `<select>` de "Tipo de reporte" sigue mostrando exactamente 2 opciones (Víctima/Tercero), sin errores de consola — confirma que V2 no filtró sus 2 tipos nuevos hacia V1 | Playwright (regresión) |

> **C6 — corrección (2026-10-31).** La versión original de esta fila decía "200". Al contrastar con el código: `reporterBestContactTime` viaja como `"HH:MM:00"`, `VictimContactForm2V2DTO` (`time.Time`, sin `UnmarshalJSON`) no lo deserializa, `JSONToStruct` descarta el error y el controller rechaza el campo en cero (`VictimContactControllerV2.go:164-168`). El mismo 400 se observó por API el 2026-10-13 con el formato que arma el wizard. Queda pendiente corregir el DTO y repetir C6 por el wizard. Los helpers de `qa-salvia` (`crearReporteQA`) usan disposición "No" para Familiar/Hombre precisamente por esto.

## Cobertura

| Sesión de prueba | Casos cubiertos |
|---|---|
| Playwright — flujo Víctima completo (UI) · hoy también `crear-reporte-victima.spec.ts` | C1, C2, C3 |
| Playwright — flujo Tercero completo (UI, anonimato = propio) | C5, C16 |
| Playwright — flujo Familiar completo (UI, disposición = Sí) | C6 (ver corrección: hoy responde 400) |
| Playwright — flujo Hombre completo (UI, disposición = No) | C7 |
| Playwright — regresión sobre V1 | C17 |
| curl — batería de validaciones de backend | C4, C8, C9, C10, C11, C12, C13, C14, C15 |

## Bugs encontrados durante esta ronda de QA

Ver el detalle completo (causa raíz + fix) en [`../related-tables.md`](../related-tables.md#bugs-encontrados-durante-implementación-y-qa-2026-09-27). Resumen: 7 bugs — 2 de frontend (pantalla invisible por `id` duplicado; popup de éxito bloqueando un botón), 4 de backend (decodificación JSON de `form2V2`; 2 variantes del mismo bug de Locale mal referenciado en captcha y en validaciones de campo requerido; una comparación de longitud siempre falsa), 1 de la migración (permisos faltantes para el rol con el que la app corre en runtime). Todos corregidos y re-verificados con los casos de esta tabla antes de cerrar la ronda.
