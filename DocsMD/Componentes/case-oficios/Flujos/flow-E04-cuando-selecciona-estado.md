━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando selecciona un estado en el filtro
   Tipo: User Interaction
   Función: reactividad de Vue sobre `estadoActivo` (v-model del <select>)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  valor:  'por_proyectar' | 'para_revisar' | 'en_correccion' |
          'aprobacion_juridica' | 'para_radicar' | 'radicado' |
          'respondido' | null (= "Todos")   → opción elegida en el <select>
}


PASO 1 — Actualizar el estado del filtro

  estadoActivo = valor


PASO 2 — Vue recalcula oficiosFiltrados (computed)

  cumpleEstado(o):
    SI estadoActivo === null:
      → true  (opción "Todos")
    SI NO:
      → o.state === estadoActivo


PASO 3 — Vue re-renderiza el grid

  → El grid se recalcula con el nuevo filtro combinado (buscador + tema + estado)
  → SI oficiosFiltrados.length === 0 → EmptyState

→ FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión                                                    | Paso afectado |
|--------------------------------------------------------------------------|---------------|
| Ninguno — filtro directo sobre un campo enum ya validado por el backend.| —             |
