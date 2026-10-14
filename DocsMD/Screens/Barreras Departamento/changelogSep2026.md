# Changelog — Barreras Departamento — septiembre 2026

## 2026-09-02 — Ajustes de listado (doc/tel, ciudad select, tipografía, estados)

Columna CASO muestra documento y teléfono (no el id del caso). Ciudad pasa a select con todas las ciudades del departamento. Estados: Abierta / En gestión / Articulada / Gestionada. Tipografía de filas alineada a Lista de casos.

**Archivos**

| Archivo | Cambio |
|---|---|
| `barrier_v2_repository.go` / DTO | `victimPhone` + filtro `cityId` |
| `barrier_v2_gin_controller.go` | Query `cityId` |
| `BarrerasDepartamentoFacade.go` | `assignedDepartmentId` en config |
| `barreras_departamento_v2.html/js` | Select ciudad, Doc/Tel, labels de estado |
| `barreras_departamento.css` | Tamaños tipo `casos-component` |

## 2026-09-02 — Listado departamental v2 (mockup 2 columnas + filtros API)

Pantalla del Enlace territorial: todas las barreras cuyo `department_id` coincide con el departamento asignado al usuario. UI nueva de dos columnas; se conserva la UI de 6 columnas. El API pagina y filtra en servidor (documento, entidad, ciudad, estado). Se corrige el 500 por tipos `city_id` integer vs varchar.

**Archivos**

| Archivo | Cambio |
|---|---|
| `DocsMD/Screens/Barreras Departamento/*` | Documentación de pantalla |
| `src/internal/repository/barrier_v2_repository.go` | SQL count/list, casts, filtros, paginación |
| `src/salvia/service/barrier_v2_service.go` | Contrato paginado |
| `src/salvia/controller/barrier_v2_gin_controller.go` | Query params + JSON `{items,total}` |
| `src/frontend/html/salvia/barriers/barreras_departamento_v2.html` | Template 2 columnas |
| `src/frontend/js/components/barreras_departamento_v2.js` | Filtros debounce + página al API |
| `src/frontend/css/barreras_departamento.css` | Estilos `bd-` |
| `src/salvia/facades/BarrerasDepartamentoFacade.go` | Monta template v2 |
| `src/salvia/config/Enums.go` | Clave `barreras_departamento_v2` |
| `DocsMD/General/proyecto-description.md` | Entrada en índice |
