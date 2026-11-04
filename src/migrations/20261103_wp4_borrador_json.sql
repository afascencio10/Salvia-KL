-- REG-38* (WP-4 T9) · JSON del borrador de Registro V3.
--
-- Mientras el Caso está en Borrador ('bo') las entidades, las barreras y los
-- metadatos que Seguimiento no guarda en columnas (entityMeta / barrierMeta) viven
-- en esta columna. Cuando el Caso deja de ser borrador ('bo' -> 'in') el servicio
-- los materializa en entity_case, entity_case_obligation, barrier_v2, entity_letter,
-- case_task y en las tablas de los equipos, y deja aquí solo entityMeta / barrierMeta
-- (la sección queda como snapshot de solo lectura). Ver DocsMD/Screens/Registro V3.
--
-- Idempotente: se puede aplicar más de una vez. En el sandbox cloud se aplica a mano:
--   psql -h 127.0.0.1 -U postgres -d salvia_local -f src/migrations/20261103_wp4_borrador_json.sql
ALTER TABLE salvia.victim_case
    ADD COLUMN IF NOT EXISTS victim_case_v3_draft_data jsonb;

COMMENT ON COLUMN salvia.victim_case.victim_case_v3_draft_data IS
    'Registro V3: entidades y barreras en borrador (entities, barriers), metadatos sin columna (entityMeta, barrierMeta) y marca materializedAt. Nula si nunca se usó.';
