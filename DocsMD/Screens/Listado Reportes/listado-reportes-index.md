# Listado de Reportes (roles `sv` y `ro`) — Índice

Adaptación de la pantalla legacy de casos para mostrar **únicamente reportes** (`victim_contact` sin `victim_case`), reutilizando la tabla existente en `get_victim_cases_sv.html` como referencia base.

## Alcance

| Incluido | Excluido (comentado, no eliminado) |
|---|---|
| Tab **Recontacto** (`fcv`) — reportes válidos sin caso | Tabs de casos: Enrutados, Por aprobar, Vencidos, Completados, Con novedades |
| Tab **Inválidos** (`fci`) — reportes invalidados sin caso | Tabla y paginación de `victimCases` |
| Filtros: nombre, apellido, teléfono | Buscador "Encontrar caso" por documento *(recomendado comentar)* |
| Botón **Crear caso** en acciones de fila (todos los roles) | — |

## Roles afectados

| Rol | Template actual | Acción |
|---|---|---|
| `sv` | `get_victim_cases_sv.html` | Adaptar (archivo base) |
| `ro` | `get_victim_cases_ro.html` | Replicar mismos cambios *(tiene header propio: Crear caso / Mis casos)* |

> **Nota:** No se unifica en un solo template porque `ro` tiene controles de cabecera distintos. Los cambios de tabla, tabs y filtros se aplican en **paralelo** en ambos archivos.

## Archivos de esta pantalla

| Archivo | Descripción |
|---|---|
| [`listado-reportes-interface.md`](listado-reportes-interface.md) | Árbol de interfaz y cambios por sección |
| [`listado-reportes-events.md`](listado-reportes-events.md) | Inventario de eventos |
| [`Flujos/flow-E01-cuando-carga-pantalla.md`](Flujos/flow-E01-cuando-carga-pantalla.md) | Carga inicial |
| [`Flujos/flow-E02-cuando-filtra-reportes.md`](Flujos/flow-E02-cuando-filtra-reportes.md) | Filtros nombre / apellido / teléfono |
| [`Flujos/flow-E03-cuando-cambia-tab-o-pagina.md`](Flujos/flow-E03-cuando-cambia-tab-o-pagina.md) | Tabs Recontacto/Inválidos y paginación |
| [`Flujos/flow-E04-cuando-presiona-crear-caso.md`](Flujos/flow-E04-cuando-presiona-crear-caso.md) | Acción Crear caso |

## Archivos de código a modificar

| Archivo | Cambio |
|---|---|
| `src/frontend/html/salvia/victim_case/get_victim_cases_sv.html` | UI: comentar casos, agregar filtros, simplificar tabs |
| `src/frontend/html/salvia/victim_case/get_victim_cases_ro.html` | Mismos cambios que `sv` (bloques equivalentes) |
| `src/salvia/config/Menu.go` | Agregar botón "Crear caso" en `menu_tool_get_victim_contacts` del rol `sv` |
| `src/salvia/dao/VictimContactDAO.go` | Extender `GetVictimContactsWithoutVictimCase` con filtros opcionales |
| `src/salvia/controllers/VictimContactController.go` | Pasar query params al DAO |
| `src/salvia/facades/VictimContactFacade.go` | Leer query params de Gin y propagarlos |

## Relación con `reportes-component`

Esta adaptación **no crea** el componente reutilizable planeado en `DocsMD/Componentes/Reportes component/`. Es un refactor incremental sobre HTML legacy. A futuro, la lógica de filtros y tabla podría migrarse a `reportes-component.js`.

## Orden de implementación sugerido

1. Backend — filtros en DAO/controller/facade (E-02)
2. `Menu.go` — botón Crear caso para `sv`
3. `get_victim_cases_sv.html` — comentar bloques de casos + UI filtros
4. `get_victim_cases_ro.html` — replicar cambios
5. Pruebas manuales por rol (`sv`, `ro`) en tabs `fcv` y `fci`
