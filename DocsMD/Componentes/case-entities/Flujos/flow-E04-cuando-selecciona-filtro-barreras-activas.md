━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando selecciona el filtro "¿Con barreras activas?"
   Tipo: User Interaction
   Función: v-model="filtroBarrerasActivas"  (computed: entidadesFiltradas)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  valor:  elegido en el <select>   → "" (Todos) | "si" | "no"
}

PASO 1 — Actualizar variable reactiva `filtroBarrerasActivas`

PASO 2 — Vue re-evalúa el computed `entidadesFiltradas`, combinando este filtro con
  sector (E03) y búsqueda de texto (E02) — sobre los datos ya cargados en E01
  (cada entidad ya trae `barrerasActivasCount` resuelto por el backend).

SEGÚN filtroBarrerasActivas:
  CASO "":    → No aplica filtro adicional → CONTINÚA
  CASO "si":  → Mantiene solo entidades con barrerasActivasCount > 0 → CONTINÚA
  CASO "no":  → Mantiene solo entidades con barrerasActivasCount === 0 → CONTINÚA

PASO 3 — Vue re-renderiza el Grid con el resultado
  SI entidadesFiltradas.length === 0:
    → Muestra EmptyState "No hay entidades que coincidan con los filtros"

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión                                                                                       | Paso afectado |
|--------------------------------------------------------------------------------------------------------------|---------------|
| Mientras no se implemente el cambio en el formulario de registro de barrera (fuera de alcance), `barrerasActivasCount` será 0 en todas las entidades — el filtro "Sí" siempre devolverá una lista vacía en la práctica hasta entonces. | PASO 2 |
