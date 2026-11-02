# QA — Registro V3 · Encabezado, estados legibles y nombre de la pantalla (WP-3 · REG-03, REG-04a/b, REG-06)

Casos de [`flow-E04`](../Flujos/flow-E04-cuando-busca-en-cabezote.md) (columna Estado) y de la parte de encabezado de [`flow-E12`](../Flujos/flow-E12-cuando-guarda-una-seccion.md).

> **Estado: ejecutado el 2026-11-02** — `scripts/e2e/wp3-navegacion-encabezado.js` (comprobaciones T7, T8 y T9), todas ✅.

**Método.** Los estados del Cabezote se fabrican por SQL sobre casos con el mismo documento (`bo`, `in`, `te`, `cd` y uno inventado); el encabezado se verifica en pantalla con un Borrador sin datos y luego guardando CI por la pantalla.

## Casos de prueba

| Caso | REG | Descripción | Condición inicial | Resultado esperado | Método |
|---|---|---|---|---|---|
| T8a | 04a | Columna Estado del Cabezote | 5 casos del mismo documento en `bo`, `in`, `te`, `cd` y `zz` | Borrador, Iniciado, Terminado, Cerrado y «Desconocido»; nunca los códigos | DOM + BD |
| T7a | 03 | Borrador sin datos | Creado sin documento ni nombre | «Persona sin identificar — Documento sin registrar», insignia «Borrador», sin «Pendiente» | DOM |
| T7b | 03 | Tras guardar CI | CI_2 «Ana Perez», CI_3 «Anita», CI_5/6 | «Sra. Ana Perez (Anita) — cc-<documento>» sin recargar | DOM |
| T8b1 / T8b2 | 04b | Insignia del encabezado | Tríada incompleta → completa | «Borrador» y luego «Iniciado» sin recargar | DOM |
| T7c | 03 | Reabrir el Caso desde el Cabezote | Caso con CI guardado | Mismo encabezado e insignia | DOM |
| T9a | 06 | Etapa de autorización | `GET /salvia/registro-v3` | Título de la pestaña y h2 «Registro de víctima»; sin «V3» | DOM |
| T9b | 06 | Aviso de tipo de contacto no disponible | «Víctima - Seguimiento» | El aviso dice «Registro de víctima»; sin «V3» | DOM |
| T9c | 06 | La ruta no cambia | `GET /salvia/registro-v3` | 200 | API |

## Cobertura

| Flujo de prueba | Casos cubiertos |
|---|---|
| `wp3-navegacion-encabezado.js` | T7a–c, T8a, T8b1–2, T9a–c |

## Resultados de ejecución

**Fecha:** 2026-11-02 · **Resultado:** todos ✅. El refresco del encabezado es asíncrono (un `GET` tras guardar): el script espera a que aparezca el nombre en lugar de leerlo al instante.
