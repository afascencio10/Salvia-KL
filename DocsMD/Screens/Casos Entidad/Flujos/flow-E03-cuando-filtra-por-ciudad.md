# flow-E03 — Cuando filtra por ciudad — ELIMINADO

> **Estado:** eliminado del alcance (jul 2026).

## Motivo

Cada usuario `et` es una sede (`entity_branch`) con `entity_branch_town_code` (ciudad implícita). No tiene sentido filtrar por ciudad: el listado ya está acotado a esa sede.

## Reemplazo

Ninguno. El pin de ciudad en la fila del caso puede seguir mostrándose como **display** (ciudad del caso), pero no hay input/select de filtro.
