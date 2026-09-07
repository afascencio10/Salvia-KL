#!/usr/bin/env python3
"""Dump psychosocial form catalog from test DB into INSERT SQL. Not for production load by itself."""

from __future__ import annotations

import csv
import io
import json
import os
import subprocess
from datetime import datetime, timezone

def _conn_from_config() -> tuple[str, str]:
    cfg_path = os.path.join(os.path.dirname(__file__), "..", "..", "config", "db_config.json")
    with open(cfg_path, encoding="utf-8") as f:
        cfg = json.load(f)
    conn = (
        f"host={cfg['hostname']} port={cfg['port']} dbname={cfg['database']} "
        f"user={cfg['user']} sslmode={cfg.get('sslmode', 'require')}"
    )
    return conn, cfg["password"]

FORM_IDS = [
    "439b57e6-07ea-4da4-9721-8ed28c6ca43f",
    "501ab3d7-8382-4447-96a6-f463152693c6",
    "7a7b61b3-7bc3-4088-a74d-0974e54a3563",
    "c31026f8-7ce2-49b9-8990-04b44ed4513c",
]
IN_LIST = ",".join(f"'{x}'" for x in FORM_IDS)

BOOL_COLS = {"required"}
INT_COLS = {"order", "min_repetitions", "max_repetitions"}
JSON_COLS = {"metadata", "options"}
RESERVED = {"order", "options"}


def psql_copy(select_sql: str) -> list[list[str]]:
    conn, password = _conn_from_config()
    env = os.environ.copy()
    env["PGPASSWORD"] = password
    cmd = [
        "psql",
        conn,
        "-v",
        "ON_ERROR_STOP=1",
        "-c",
        f"COPY ({select_sql}) TO STDOUT WITH (FORMAT csv, NULL '\\N', ENCODING 'UTF8')",
    ]
    r = subprocess.run(cmd, capture_output=True, env=env)
    if r.returncode != 0:
        raise SystemExit(r.stderr.decode("utf-8", errors="replace"))
    text = r.stdout.decode("utf-8")
    if not text.strip():
        return []
    reader = csv.reader(io.StringIO(text))
    return [row for row in reader]


def ident(col: str) -> str:
    return f'"{col}"' if col in RESERVED else col


def lit(col: str, raw: str | None) -> str:
    if raw is None or raw == r"\N":
        return "NULL"
    if col in BOOL_COLS:
        return "TRUE" if raw in ("t", "true", "1") else "FALSE"
    if col in INT_COLS:
        return raw
    if col in JSON_COLS:
        escaped = raw.replace("'", "''")
        return f"'{escaped}'::jsonb"
    escaped = raw.replace("'", "''")
    return f"'{escaped}'"


def inserts(table: str, cols: list[str], rows: list[list[str]]) -> list[str]:
    col_sql = ", ".join(ident(c) for c in cols)
    out = []
    for row in rows:
        vals = ", ".join(lit(c, v) for c, v in zip(cols, row))
        out.append(
            f"INSERT INTO salvia.{table} ({col_sql}) VALUES ({vals}) ON CONFLICT (id) DO NOTHING;"
        )
    return out


def main() -> None:
    form_cols = [
        "id",
        "name",
        "description",
        "status",
        "created_at",
        "updated_at",
        "deleted_at",
        "campaign_id",
    ]
    section_cols = ["id", "form_id", "name", "description", "order", "created_at", "updated_at", "deleted_at"]
    rg_cols = [
        "id",
        "form_section_id",
        "name",
        "item_name",
        "order",
        "min_repetitions",
        "max_repetitions",
        "state_items",
        "created_at",
        "updated_at",
        "deleted_at",
        "add_button_label",
    ]
    q_cols = [
        "id",
        "form_id",
        "form_section_id",
        "repeater_group_id",
        "question_type",
        "description",
        "required",
        "state_options_path",
        "order",
        "created_at",
        "updated_at",
        "deleted_at",
        "metadata",
        "options",
    ]
    opt_cols = ["id", "question_id", "label", "value", "order", "created_at", "updated_at", "deleted_at"]
    vc_cols = [
        "id",
        "target_type",
        "target_id",
        "trigger_question_id",
        "trigger_state_path",
        "trigger_value",
        "operator",
        "created_at",
        "updated_at",
        "trigger_option_id",
    ]
    rm_cols = [
        "id",
        "target_type",
        "target_id",
        "target_field",
        "modification_type",
        "state_path",
        "search_string",
        "created_at",
        "updated_at",
    ]

    forms = psql_copy(
        f"SELECT {', '.join(ident(c) for c in form_cols)} FROM salvia.form WHERE id::text IN ({IN_LIST}) ORDER BY name"
    )
    sections = psql_copy(
        f"SELECT {', '.join(ident(c) for c in section_cols)} FROM salvia.form_section WHERE form_id::text IN ({IN_LIST}) ORDER BY form_id, \"order\", name"
    )
    rgs = psql_copy(
        f"""
        SELECT {", ".join("rg." + ident(c) for c in rg_cols)}
        FROM salvia.repeater_group rg
        JOIN salvia.form_section s ON rg.form_section_id::text = s.id::text
        WHERE s.form_id::text IN ({IN_LIST})
        ORDER BY rg.id
        """
    )
    questions = psql_copy(
        f"SELECT {', '.join(ident(c) for c in q_cols)} FROM salvia.question WHERE form_id::text IN ({IN_LIST}) ORDER BY form_id, \"order\", id"
    )
    options = psql_copy(
        f"""
        SELECT {", ".join("o." + ident(c) for c in opt_cols)}
        FROM salvia.option o
        JOIN salvia.question q ON o.question_id = q.id
        WHERE q.form_id::text IN ({IN_LIST})
        ORDER BY o.question_id, o.\"order\", o.id
        """
    )
    vcs = psql_copy(
        f"""
        SELECT {", ".join("vc." + ident(c) for c in vc_cols)}
        FROM salvia.visibility_condition vc
        WHERE vc.target_id IN (SELECT id FROM salvia.question WHERE form_id::text IN ({IN_LIST}))
           OR vc.target_id IN (SELECT id::text FROM salvia.form_section WHERE form_id::text IN ({IN_LIST}))
           OR vc.trigger_question_id IN (SELECT id FROM salvia.question WHERE form_id::text IN ({IN_LIST}))
        ORDER BY vc.id
        """
    )
    rms = psql_copy(
        f"""
        SELECT {", ".join("rm." + ident(c) for c in rm_cols)}
        FROM salvia.render_modification rm
        WHERE rm.target_id IN (SELECT id FROM salvia.question WHERE form_id::text IN ({IN_LIST}))
           OR rm.target_id IN (SELECT id::text FROM salvia.form_section WHERE form_id::text IN ({IN_LIST}))
           OR rm.target_id IN (
                SELECT rg.id FROM salvia.repeater_group rg
                JOIN salvia.form_section s ON rg.form_section_id::text = s.id::text
                WHERE s.form_id::text IN ({IN_LIST})
           )
        ORDER BY rm.id
        """
    )

    counts = {
        "form": len(forms),
        "form_section": len(sections),
        "repeater_group": len(rgs),
        "question": len(questions),
        "option": len(options),
        "visibility_condition": len(vcs),
        "render_modification": len(rms),
    }

    now = datetime.now(timezone.utc).strftime("%Y-%m-%d %H:%M UTC")
    lines = [
        "-- =============================================================================",
        "-- BACKUP catálogo formularios psicosociales (solo definición, sin submissions)",
        f"-- Origen: BD de pruebas (Supabase pooler)  Extraído: {now}",
        "-- Forms:",
        "--   439b57e6-07ea-4da4-9721-8ed28c6ca43f  Primer Contacto Psicosocial",
        "--   501ab3d7-8382-4447-96a6-f463152693c6  Primera Atención Psicosocial",
        "--   7a7b61b3-7bc3-4088-a74d-0974e54a3563  Atención Psicosocial",
        "--   c31026f8-7ce2-49b9-8990-04b44ed4513c  Cierre Psicosocial",
        f"-- Conteos: {json.dumps(counts)}",
        "-- Idempotente: ON CONFLICT (id) DO NOTHING",
        "-- NO ejecutar en prod sin revisar. No incluye answer/form_submission.",
        "-- =============================================================================",
        "",
        "BEGIN;",
        "",
        "-- 1. form",
        *inserts("form", form_cols, forms),
        "",
        "-- 2. form_section",
        *inserts("form_section", section_cols, sections),
        "",
        "-- 3. repeater_group",
        *inserts("repeater_group", rg_cols, rgs),
        "",
        "-- 4. question",
        *inserts("question", q_cols, questions),
        "",
        "-- 5. option",
        *inserts("option", opt_cols, options),
        "",
        "-- 6. visibility_condition",
        *inserts("visibility_condition", vc_cols, vcs),
        "",
        "-- 7. render_modification",
        *inserts("render_modification", rm_cols, rms),
        "",
        "COMMIT;",
        "",
    ]

    out = os.path.join(os.path.dirname(__file__), "backup_psicosocial_catalog_from_test.sql")
    with open(out, "w", encoding="utf-8") as f:
        f.write("\n".join(lines))
    print(json.dumps({"wrote": out, **counts}, indent=2))


if __name__ == "__main__":
    main()
