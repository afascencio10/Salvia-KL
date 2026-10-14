# flow-E05 — Cuando presiona "Ver caso"

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando presiona "Ver caso"
   Tipo: User Interaction
   Función: goToCase(item)
   Estado: implementado
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  item.caseCode || item.caseId: victim_case_i_code
}

PASO 1 — SI falta código → toast "No se encontró el identificador del caso." → TERMINAR

PASO 2 — Navegar al detalle
  window.location.href = '/salvia/casos/' + encodeURIComponent(code) + '/detalle'
  → CaseDetailGET · permiso get_case_detail_sv (incluye rol mc)
```
