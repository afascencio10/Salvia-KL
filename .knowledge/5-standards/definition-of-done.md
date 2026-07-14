---
okf_version: "1.0"
type: Coding_Standard
title: "Definition of Done (DoD)"
description: "Criterios objetivos y verificables que todo feature, bugfix o mejora debe cumplir antes de ser considerado terminado y listo para producción."
owner: "@platform-team"
status: active
tags: [standards, process, quality, dod]
last_updated: "2026-06-23"
---

# Definition of Done (DoD)

Este documento define los criterios **no negociables** que todo cambio debe cumplir antes de hacer merge a `main`. Es el contrato entre el developer y el equipo.

## 1. Especificación (Spec-First)

| # | Criterio | Verificación |
| :--- | :--- | :--- |
| DoD-S1 | Existe un documento OKF que describe el cambio **antes** de escribir código | `git log` muestra commit de spec antes de commit de code |
| DoD-S2 | El frontmatter OKF del spec pasa validación CI | Pipeline verde en stage `validate` |
| DoD-S3 | Las `dependencies:` del spec apuntan a documentos que existen | Script `validate_okf.py --check-refs` pasa |
| DoD-S4 | Si es un endpoint nuevo, tiene JSON Schema formal en el spec | Revisión manual en PR |

## 2. Código

| # | Criterio | Verificación |
| :--- | :--- | :--- |
| DoD-C1 | El código sigue los [coding-standards.md](/.knowledge/5-standards/coding-standards.md) | Linter configurado pasa sin warnings |
| DoD-C2 | No hay `TODO`, `HACK`, o `FIXME` sin issue vinculado | `grep -r "TODO\|HACK\|FIXME" --include="*.ts"` revisado |
| DoD-C3 | No hay secretos hardcodeados (API keys, passwords) | Secret scanner en CI |
| DoD-C4 | Los commits siguen [commit-convention.md](/.knowledge/5-standards/commit-convention.md) | commitlint en CI |

## 3. Testing

| # | Criterio | Verificación |
| :--- | :--- | :--- |
| DoD-T1 | Tests unitarios escritos para la lógica nueva (≥80% coverage del cambio) | Coverage report en CI |
| DoD-T2 | Tests de integración para endpoints nuevos o modificados | CI stage `test:integration` verde |
| DoD-T3 | Tests E2E actualizados si el cambio afecta flujos críticos | CI stage `test:e2e` verde |
| DoD-T4 | Todos los tests existentes siguen pasando (zero regressions) | CI completo verde |

## 4. Documentación

| # | Criterio | Verificación |
| :--- | :--- | :--- |
| DoD-D1 | El spec OKF fue actualizado si el comportamiento cambió | Diff del PR incluye cambios en `.knowledge/` |
| DoD-D2 | Si es un endpoint público, el contrato de API está documentado | Revisión en PR |
| DoD-D3 | Si hay breaking changes, están documentados en el spec con `> [!WARNING]` | Revisión en PR |

## 5. Code Review

| # | Criterio | Verificación |
| :--- | :--- | :--- |
| DoD-R1 | Al menos 1 approval de un peer | GitHub/GitLab PR rules |
| DoD-R2 | Todos los comentarios del reviewer están resueltos | PR sin threads abiertos |
| DoD-R3 | El reviewer verificó el [code-review-checklist.md](/.knowledge/5-standards/code-review-checklist.md) | Checklist template en PR |

## 6. Despliegue

| # | Criterio | Verificación |
| :--- | :--- | :--- |
| DoD-P1 | El pipeline completo (validate → build → test → security) pasa en CI | Pipeline verde |
| DoD-P2 | Si hay migración de base de datos, fue probada en staging | Confirmación del developer |
| DoD-P3 | El runbook fue actualizado si el cambio afecta operaciones | Revisión en PR |

> [!IMPORTANT]
> Si un criterio no aplica a un cambio específico (ej. un fix de typo no necesita tests E2E), el developer debe **explícitamente** justificarlo en la descripción del PR con la frase: `DoD-XX: N/A — [razón]`.
