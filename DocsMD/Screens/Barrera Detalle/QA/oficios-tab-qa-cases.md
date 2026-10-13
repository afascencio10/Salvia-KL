# QA — Tab "Oficios" en Detalle de Barrera

Casos de prueba derivados de `req-oficios-en-detalle-barrera.md` (`DocsMD/Otros/temp/`). Cubre el nuevo prop `barrierId` de `case-oficios.js` y su montaje en la nueva tab "Oficios" de `barrera_detalle.html`. No hay endpoint nuevo — reutiliza `GET /api/v1/entity-letters?barrierId=`, ya cubierto por el backend existente.

---

## Casos de prueba

| Caso | Descripción | Condición inicial | Resultado esperado | Método de verificación |
|---|---|---|---|---|
| C1 | Tab "Oficios" aparece en el orden correcto | Cualquier barrera | Tabs en orden: Información General, Tareas, Oficios, Timeline | Playwright — orden de `.brd-tab` |
| C2 | Click en la tab pide oficios filtrados por `barrierId` | Cualquier barrera | Request de red a `/api/v1/entity-letters?barrierId=<id>`, no `?caseId=` | Playwright — `page.on('request')` |
| C3 | Chips de "tema" ocultos en modo barrera | Tab Oficios activa | `.co-chips-tema` no se renderiza (`toHaveCount(0)`) | Playwright |
| C4 | Estado vacío si la barrera no tiene oficios | Barrera fixture sin `entity_letter` asociados | `.co-empty` visible, sin error en consola | Playwright |
| C5 | Cards visibles si la barrera tiene oficios | Barrera con ≥1 `entity_letter` | `.co-card` visible ≥1, cada card corresponde a un oficio de esa barrera | Playwright |
| C6 | Buscador y filtro de estado siguen operativos en modo barrera | Tab Oficios con datos cargados | Escribir en el buscador / cambiar `<select>` de estado filtra las cards visibles | Playwright |
| C7 | Click en card abre modal con "Barrera relacionada" | Card visible | Modal muestra sección "Barrera relacionada" con link a `/salvia/barreras/:id` | Playwright |
| C8 | Compatibilidad — Detalle del Caso no cambia | Tab "Gestión institucional" de Detalle del Caso (modo `caseId`, sin `barrierId`) | Chips de tema siguen visibles, request sigue usando `?caseId=` | Playwright |

---

## Cobertura

| Flujo de prueba | Casos cubiertos |
|---|---|
| `TC-OF-01` — tab nueva, filtro por barrierId, chips ocultos | C1, C2, C3 |
| `TC-OF-02` — estado vacío | C4 |
| `TC-OF-03` — cards + filtros + modal de detalle | C5, C6, C7 |
| `TC-OF-04` — compatibilidad con Detalle del Caso | C8 |

Archivo de tests: `qa-salvia/tests/barrera-detalle/oficios-tab.spec.ts`
Page object: `qa-salvia/pages/barrera-detalle.page.ts` (`switchToOficiosTab`, `oficiosChipsTema`, `oficiosCards`, `oficiosEmptyState`)

## Backend

Sin endpoints nuevos. `GET /api/v1/entity-letters?barrierId=` ya existe y se probó indirectamente al construir la planeación (mismo servicio/repositorio que `?caseId=`, mismo shape de respuesta) — no requirió pruebas de backend adicionales.

## Datos de prueba

| Variable (`.env.local`) | Uso |
|---|---|
| `GESTION_PROPIA_BARRIER_OTHER_DEPT` | Barrera fixture sin `entity_letter` propios → caso C4 (estado vacío) |
| `GESTION_PROPIA_CASE_SAME_DEPT` | Caso real usado para C8 (Detalle del Caso) |
| Descubrimiento dinámico vía `GET /api/v1/entity-letters?limit=1` | Se usa para encontrar un `barrierId` real con al menos un oficio existente → casos C1, C2, C3, C5, C6, C7 (evita depender de un fixture hardcodeado que pueda quedar sin datos) |
