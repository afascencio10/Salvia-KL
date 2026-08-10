# flow-E07 — Cuando cambia entidad — ELIMINADO

> **Estado:** eliminado del alcance (jul 2026).

## Motivo

La relación usuario `et` ↔ sede quedó definida:

- `security.general_user.entity_branch_id` → sede fija del usuario
- Un usuario `et` **es** un `entity_branch`; no puede cambiar de entidad

El selector/dropdown de organización ya no aplica.

## Reemplazo

La sede se resuelve en **E01** desde la sesión. El header muestra la entidad padre (`entity`) de esa sede.
