-- HU-027: Agrega campos del calendario a follow_up_v2
-- Ejecutar UNA SOLA VEZ antes de arrancar la aplicación con la nueva versión.

ALTER TABLE salvia.follow_up_v2
    ADD COLUMN IF NOT EXISTS team                VARCHAR(50),
    ADD COLUMN IF NOT EXISTS risk_status         VARCHAR(20),
    ADD COLUMN IF NOT EXISTS scheduled_date      DATE,
    ADD COLUMN IF NOT EXISTS is_completed        BOOLEAN DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS sequence_number     INT DEFAULT 0,
    ADD COLUMN IF NOT EXISTS form_submission_id  VARCHAR(36);

-- Índices para consultas frecuentes del calendario
CREATE INDEX IF NOT EXISTS idx_fup2_status      ON salvia.follow_up_v2(status);
CREATE INDEX IF NOT EXISTS idx_fup2_case_status ON salvia.follow_up_v2(case_id, status);
