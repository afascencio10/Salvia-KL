━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando selecciona un chip de tema
   Tipo: User Interaction
   Función: toggleTema(tema)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  tema:  'barrera' | 'emergencia' | 'psicosocial' | 'estabilizacion'   → chip presionado
}


PASO 1 — Verificar si el chip está habilitado

  SI tema !== 'barrera':
    → No hacer nada (chip renderizado con :disabled — ver GAP de schema
      en case-oficios-interface.md)
    → TERMINAR ejecución

  SI tema === 'barrera':
    → CONTINÚA FLUJO GENERAL


PASO 2 — Alternar el filtro

  SI temaActivo === tema:
    → temaActivo = null   // el usuario volvió a presionar el chip activo → lo desactiva
  SI NO:
    → temaActivo = tema   // activa el filtro (single-select: reemplaza cualquier tema previo)


PASO 3 — Vue recalcula oficiosFiltrados (computed)

  cumpleTema(o):
    SI temaActivo === null:
      → true  (sin filtro activo)
    SI temaActivo === 'barrera':
      → o.barrierId != null
    SI temaActivo === 'emergencia' | 'psicosocial' | 'estabilizacion':
      → false  (no hay campo que evaluar todavía — GAP de schema)

  → El chip presionado cambia visualmente a estado "activo" (:class="chip-activo")
  → El grid se recalcula con el nuevo filtro combinado (buscador + tema + estado)

→ FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión                                                    | Paso afectado |
|--------------------------------------------------------------------------|---------------|
| Los 3 chips no-barrera están deshabilitados hasta que se agreguen        | PASO 1 / 3    |
| emergency_measure_id / psychosocial_support_id /                        |               |
| economic_stabilization_id a entity_letter (ver case-oficios-interface). |               |
| ¿El filtro de tema es single-select (como se documentó aquí) o          | PASO 2        |
| multi-select (un oficio podría relacionarse a más de un tema a la vez)? |               |
| Depende de si en la práctica un entity_letter puede tener más de un    |               |
| ID de tema seteado simultáneamente.                                     |               |
