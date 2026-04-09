-- HU-027 Ajuste: Riesgo Extremo tiene 5 seguimientos, S1 a las 4 horas.
-- Cambiar scheduled_date de DATE a TIMESTAMPTZ para almacenar la hora del S1 extremo.

ALTER TABLE salvia.follow_up_v2
    ALTER COLUMN scheduled_date TYPE TIMESTAMPTZ USING scheduled_date::TIMESTAMPTZ;
