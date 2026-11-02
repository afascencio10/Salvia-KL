# QA — Registro V3 · E-21 Cuando un guardado o una carga falla (WP-2 · REG-11)

Casos de prueba de [`flow-E21`](../Flujos/flow-E21-cuando-un-guardado-o-una-carga-falla.md): lo que ve la persona cuando guardar una sección o el Tamizaje, abrir un registro o salir de la llamada falla. Origen: comentario de Felipe en la planeación de WP-2 — «asegúrate de que si algo falla el usuario lo sepa; si es por el formato del dato, pedirle al usuario que lo corrija». El comportamiento del API (400 con `campo` y `detalle`, atomicidad, logs) está en [`E12-guardar-seccion-qa-cases.md`](E12-guardar-seccion-qa-cases.md) y no se repite aquí; esta lista es la parte de **pantalla** (`wp2-reg11-guardado-seccion.js`, fase 2, y `wp2-reg12-tamizaje-tablas.js`, caso U4).

> **Estado: ejecutado el 2026-11-02** — resultados al final de este archivo.

**Método.** Solo pantalla (lo que el backend no cubre). Las fallas del servidor se simulan interceptando la petición con `page.route` (500 con el cuerpo real del servidor, o `abort` para "sin conexión"); el dato inválido se provoca de verdad (una dirección de 151 caracteres) para que el 400 lo devuelva el backend. `window.confirm` se maneja con `page.once('dialog')`. Los casos de aviso se verifican en el DOM (`#registro-v3-error-guardado`) y la BD confirma que no se guardó nada.

---

## Casos de prueba

| Caso | Descripción | Condición inicial | Resultado esperado | Método de verificación |
|---|---|---|---|---|
| E21-C1 | Reabrir un Caso muestra lo guardado con el formato de los controles | Caso con hechos guardados (fecha `2026-10-30`, hora `19:30`, departamento «Atlántico», municipio «Barranquilla»), abierto desde el Cabezote | En Hechos de violencia: fecha `2026-10-30`, hora `19:30`, departamento «Atlántico», municipio «Barranquilla» (antes la fecha y el departamento salían en blanco) | DOM |
| E21-C2 | Dato inválido: aviso que nombra la pregunta; no se guarda; al corregir sí | En Hechos de violencia se escribe una dirección de 151 caracteres y se pulsa «Guardar sección»; luego se corrige a «Calle corregida 99» | Aviso rojo «No se guardó la sección. Corrige «¿Cuál es la dirección donde ocurrieron los hechos…?»: Este dato admite máximo 150 caracteres y escribiste 151. Luego vuelve a guardar.»; la BD conserva la dirección anterior. Al corregir: «Sección guardada.», el aviso rojo desaparece y la BD tiene la dirección nueva | DOM + BD |
| E21-C3 | Error del servidor (500) | `POST …/seccion/HV` responde 500 `{error:"error al guardar la sección"}` | Aviso rojo «No se pudo guardar la sección por un error del servidor. Lo que escribiste sigue en pantalla…»; no aparece «Sección guardada.» (antes la pantalla decía «Sección guardada.») | DOM |
| E21-C4 | Sin conexión | La petición de guardado se aborta | Aviso rojo «No se pudo guardar la sección: no hay conexión con el servidor. Lo que escribiste sigue en pantalla; intenta de nuevo.» | DOM |
| E21-C5 | «Finalizar llamada» con el guardado fallando | Misma falla; se pulsa «Finalizar llamada (guardar y salir)» | Aparece un diálogo con el motivo y «¿Salir de todos modos? Lo que escribiste en este módulo se perdería.»; al cancelar la persona sigue en el Formulario; al aceptar vuelve al Cabezote (antes salía sin avisar) | DOM (diálogo) |
| E21-C6 | No se puede abrir el Caso | `GET …/:icode` responde 500 al pulsar «Abrir» | El Cabezote muestra «No se pudo abrir el registro. Intenta de nuevo; si el problema continúa, avisa a soporte.» y no cambia de pantalla (antes no mostraba nada) | DOM |
| E21-C7 | Falla al guardar el Tamizaje | `POST …/tamizaje` responde 500 al pulsar «Calcular y guardar Tamizaje» | Aviso rojo «No se pudo guardar el Tamizaje por un error del servidor. Lo que respondiste sigue en pantalla…» (antes se mostraba `{error}` como resultado) | DOM |

---

## Cobertura

| Flujo de prueba | Casos cubiertos | Script |
|---|---|---|
| `wp2-reg11` — U0, U1 | E21-C1 | `scripts/e2e/wp2-reg11-guardado-seccion.js` |
| `wp2-reg11` — U2a, U2b | E21-C2 | ídem |
| `wp2-reg11` — U3, U4 | E21-C3, E21-C4 | ídem |
| `wp2-reg11` — U5 | E21-C5 | ídem |
| `wp2-reg11` — U6 | E21-C6 | ídem |
| `wp2-reg12` — U4 | E21-C7 | `scripts/e2e/wp2-reg12-tamizaje-tablas.js` |

## Resultados de ejecución

**Fecha:** 2026-11-02 · **Build:** rama `claude/busy-edison-uk6g2a` en `http://127.0.0.1:9090` (sandbox cloud) · **Herramienta:** Playwright headless (una sola sesión `op`: un segundo login del mismo usuario invalida la primera y daba 401 en la corrida de prueba del propio script).

| Caso | Resultado | Antes del arreglo |
|---|---|---|
| E21-C1 | ✅ | Fecha y departamento en blanco |
| E21-C2 | ✅ | Un 500 se leía como éxito: decía «Sección guardada.» sin haber guardado |
| E21-C3 | ✅ | Decía «Sección guardada.» |
| E21-C4 | ✅ | Texto genérico «Error al guardar la sección.» |
| E21-C5 | ✅ | Salía sin avisar |
| E21-C6 | ✅ | Sin aviso |
| E21-C7 | ✅ | `{error}` como resultado |

**Resumen:** 7 de 7 casos ✅. Contra el código anterior fallaba E21-C1 (fecha y departamento en blanco) y el script se detenía en E21-C2 porque no existía el aviso rojo.
