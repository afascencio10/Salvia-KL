# Cabezote — Casos de prueba

Cubre E02 (Habeas Data), E04 (buscar) y E06 (nuevo registro). Endpoints: `GET /api/v1/victim-cases-v3/buscar`, `POST /api/v1/victim-cases-v3/draft`.

## Casos

| Caso | Descripción | Condición inicial | Resultado esperado | Método de verificación |
|---|---|---|---|---|
| C1 | Responde "No" en Habeas Data | Pantalla recién cargada, `etapa='habeas'` | No avanza a Cabezote; no se crea ningún `victim_case` ni `general_user` | Validar DOM (sigue en pantalla de Habeas Data) + leer BD (sin filas nuevas) |
| C2 | Responde "Sí" en Habeas Data | Igual que C1 | Avanza a `etapa='cabezote'` | Validar DOM |
| C3 | Busca solo con documento | `etapa='cabezote'`, campo teléfono vacío | `GET /buscar?documento=X` — 200, filtra por `victim_case_victim_doc_number` | Request directo al backend |
| C4 | Busca solo con teléfono | Campo documento vacío | `GET /buscar?telefono=X` — 200, filtra por `victim_case_form2_victim_phone` | Request directo |
| C5 | Busca sin documento ni teléfono | Ambos vacíos | 400 `{"error":"digite al menos un criterio..."}` — el botón Buscar debería estar deshabilitado o mostrar el mismo error en el DOM | Request directo + validar DOM |
| C6 | Busca por **nombre** (debe ignorarse) | Hay un registro cuyo nombre coincide pero documento/teléfono no | Query de nombre no existe en el endpoint — confirmar que no hay ningún input de "nombre" en el Cabezote (D7) | Validar DOM (no debe existir el campo) |
| C7 | Busca con documento que no tiene coincidencias | Documento nuevo, sin registros previos | `resultadosBusqueda=[]`, se muestra el botón "+ Nuevo registro" | Validar DOM |
| C8 | Busca con documento que **sí** tiene un Caso previo (evita duplicado) | Ya existe un Caso con ese documento (de C9 u otro caso previo) | El Caso existente aparece en la tabla de resultados con su `status` real; NO aparece botón para crear uno nuevo con ese documento sin antes ver el resultado | Validar DOM + comparar contra BD |
| C9 | Crea "Nuevo Registro" con documento + teléfono | Búsqueda sin resultados (C7) | `POST /draft` → 200 `{icode}`; en BD: `victim_case.status='bo'`, `victim_case_form2_victim_phone` = el teléfono digitado (**H4**, ya corregido — antes quedaba en 0) | Request directo + leer BD |
| C10 | Crea "Nuevo Registro" solo con teléfono (sin documento) | Búsqueda solo por teléfono, sin resultados | `POST /draft` con `documento=""` — el backend usa el dummy `"0000000000"` para doc_number; confirmar que NO revienta (antes de H1 sí reventaba) | Request directo + leer BD |
| C11 | Abre un Caso existente desde resultados | Resultado de C8 visible | `GET /:icode` → `etapa='formulario'`, `resumenCaso` con los datos reales (no "Pendiente") | Validar DOM + request directo |
| C12 | Reabre un Borrador recién creado sin ninguna respuesta | Caso de C9, sin ningún `SaveSection` todavía | `GET /:icode` → `answers` contiene solo lo prediligenciado (CI_4=teléfono); el resto de los 60+ campos NO aparece como respondido aunque `GetAnswers` devuelva valores dummy (ver H3 — la UI no debe mostrarlos como si el agente los hubiera llenado) | Request directo, revisar que el JSON de respuesta distinga lo prediligenciado del resto |
| C13 | Presiona "Limpiar" tras una búsqueda | Hay resultados visibles (C7/C8) | Se vacía `buscar.documento`/`buscar.telefono`, `resultadosBusqueda=[]` | Validar DOM |

## Notas para quien ejecute en Playwright

- El tipo de documento (`tipo_documento`) es opcional en la búsqueda — probar también C3/C4 variando si se incluye o no, ya que el backend solo aplica el filtro de tipo cuando hay documento (`internal/repository/victim_case_v3_repository.go`, `FindByDocumentOrPhone`).
- C9/C10 van a crear filas reales en la BD de pruebas — limpiar después (`victim_case`, `victim_case_form2`, `general_user`, `general_user_profile`, `rel_role_general_user`) para no ensuciar búsquedas de otros casos.
