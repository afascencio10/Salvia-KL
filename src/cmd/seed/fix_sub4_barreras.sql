-- Agrega iteración 2 (Justicia) e iteración 3 (Protección) al repeater de sub4

DO $$
DECLARE
  v_sub4   varchar := 'ed09276b-47b5-4ea5-9a43-49afb45f9d7b';
  v_rg_id  varchar := '5fd3ecdc-2e5f-4b31-97ef-8a994580586a';

  q_sector           varchar := 'f19378b6-55c5-4fdf-b765-7ebcc3978741';
  q_barreras_justicia  varchar := '2bec977e-97c7-42d7-a00a-b536af8038eb';
  q_barreras_proteccion varchar := '66c9fc1e-9b5e-4ad4-999f-7aeb483d84dc';

  v_entry2 varchar;
  v_entry3 varchar;
BEGIN

  -- Iteración 2 — Justicia
  INSERT INTO salvia.repeater_entry (id, form_submission_id, repeater_group_id, iteration)
  VALUES (gen_random_uuid()::varchar, v_sub4, v_rg_id, 2)
  RETURNING id INTO v_entry2;

  INSERT INTO salvia.answer (id, form_submission_id, question_id, repeater_entry_id, value) VALUES
    (gen_random_uuid()::varchar, v_sub4, q_sector,             v_entry2, 'justicia'),
    (gen_random_uuid()::varchar, v_sub4, q_barreras_justicia,  v_entry2, 'justicia_denuncia,justicia_celeridad');

  -- Iteración 3 — Protección
  INSERT INTO salvia.repeater_entry (id, form_submission_id, repeater_group_id, iteration)
  VALUES (gen_random_uuid()::varchar, v_sub4, v_rg_id, 3)
  RETURNING id INTO v_entry3;

  INSERT INTO salvia.answer (id, form_submission_id, question_id, repeater_entry_id, value) VALUES
    (gen_random_uuid()::varchar, v_sub4, q_sector,               v_entry3, 'proteccion'),
    (gen_random_uuid()::varchar, v_sub4, q_barreras_proteccion,  v_entry3, 'prot_incumplim');

  RAISE NOTICE 'entry2 (justicia): % | entry3 (proteccion): %', v_entry2, v_entry3;

END $$;
