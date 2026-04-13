-- ============================================================================
-- DATOS SEMILLA: Seguimientos del Área (demo para desarrollo)
-- Ejecutar después de las migraciones de estructura.
-- Estos datos simulan el mockup con nombres legibles.
-- ============================================================================

BEGIN;

-- Limpiar datos de demo anteriores (si existen)
DELETE FROM salvia.follow_up_v2 WHERE case_id LIKE 'SAL-%';

-- Insertar seguimientos de ejemplo
-- Agentes: Agente López, Agente Ramírez, Agente Díaz, Agente Torres, Agente Mendoza
-- Casos con nombres legibles en case_id para el demo

INSERT INTO salvia.follow_up_v2
    (case_id, agent_id, team, risk_status, scheduled_date, status, sequence_number, is_completed)
VALUES
    -- Ejecutados (pasados)
    ('SAL-001', 'Agente López',   'RIESGO_BAJO', 'ALTO',     '2026-03-24 10:00', 'REALIZADO', 1, true),
    ('SAL-002', 'Agente Ramírez', 'RIESGO_BAJO', 'EXTREMO',  '2026-03-25 08:00', 'REALIZADO', 1, true),
    ('SAL-005', 'Agente Vásquez', 'RIESGO_BAJO', 'BAJO',     '2026-03-26 14:00', 'REALIZADO', 1, true),
    ('SAL-004', 'Agente Díaz',    'RIESGO_BAJO', 'ALTO',     '2026-03-27 09:00', 'REALIZADO', 1, true),
    ('SAL-003', 'Agente Torres',  'RIESGO_BAJO', 'MODERADO', '2026-03-28 11:00', 'REALIZADO', 1, true),
    ('SAL-008', 'Agente Jiménez', 'RIESGO_BAJO', 'ALTO',     '2026-03-29 16:00', 'REALIZADO', 1, true),

    -- Pendientes (fecha actual/futura)
    ('SAL-002', 'Agente Ramírez', 'RIESGO_BAJO', 'EXTREMO',  '2026-03-30 08:00', 'PENDIENTE', 2, false),
    ('SAL-004', 'Agente Díaz',    'RIESGO_BAJO', 'ALTO',     '2026-03-30 09:00', 'PENDIENTE', 2, false),
    ('SAL-006', 'Agente Ramírez', 'RIESGO_BAJO', 'EXTREMO',  '2026-03-30 09:30', 'PENDIENTE', 1, false),
    ('SAL-001', 'Agente López',   'RIESGO_BAJO', 'ALTO',     '2026-03-30 10:00', 'PENDIENTE', 2, false),
    ('SAL-008', 'Agente Díaz',    'RIESGO_BAJO', 'ALTO',     '2026-03-30 10:30', 'PENDIENTE', 2, false),
    ('SAL-003', 'Agente Torres',  'RIESGO_BAJO', 'MODERADO', '2026-03-30 11:00', 'PENDIENTE', 2, false),
    ('SAL-007', 'Agente López',   'RIESGO_BAJO', 'MODERADO', '2026-03-30 12:00', 'PENDIENTE', 1, false),
    ('SAL-009', 'Agente Mendoza', 'RIESGO_BAJO', 'ALTO',     '2026-03-30 13:00', 'PENDIENTE', 1, false),
    ('SAL-001', 'Agente López',   'RIESGO_BAJO', 'ALTO',     '2026-03-31 10:00', 'PENDIENTE', 3, false),
    ('SAL-002', 'Agente Ramírez', 'RIESGO_BAJO', 'EXTREMO',  '2026-03-31 08:00', 'PENDIENTE', 3, false);

COMMIT;
