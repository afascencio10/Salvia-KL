# flow-E05 — Cuando presiona "Ver caso"

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando presiona "Ver caso"
   Tipo: User Interaction
   Función: goToCase(item)
   Estado: implementado
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  item.caseId: entity_case.case_id (= victim_case_i_code)
}

PASO 1 — SI !item.caseId → toast error → TERMINAR

PASO 2 — Navegar al detalle moderno
  window.location.href = '/salvia/casos/' + encodeURIComponent(item.caseId) + '/detalle'
  → CaseDetailGET · permiso get_case_detail_sv (incluye rol et)
```
