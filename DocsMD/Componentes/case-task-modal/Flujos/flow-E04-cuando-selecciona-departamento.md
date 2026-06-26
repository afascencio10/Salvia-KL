━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando selecciona departamento
   Tipo: User Interaction
   Función: onSelectDepartamento()
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Aplica a: formularios 'gestion_llamada' y 'proyectar_oficio' únicamente.

INPUT: {
  form.departamentoId:  ID del departamento seleccionado  → <select> Departamento
}


PASO 1 — Filtrar ciudades del departamento seleccionado

  ciudades = allCities.filter(c => c.departmentId === form.departamentoId)
  // filtrado en memoria, sin fetch


PASO 2 — Resetear en cascada los niveles inferiores

  form.ciudadId      = null
  form.municipioId   = null
  form.entidadId     = null
  form.entidadNombre = ""
  municipios         = []
  entidades          = []


PASO 3 — Vue re-evalúa reactivamente

  → Dropdown Ciudad muestra las ciudades filtradas y queda habilitado
    (:disabled si !form.departamentoId → ahora departamentoId tiene valor)
  → Dropdown Municipio permanece deshabilitado (:disabled si !form.ciudadId)
  → Dropdown Entidad permanece deshabilitado (:disabled si !form.municipioId)
  → Campo "Nombre de la entidad" permanece oculto ([v-if form.entidadId === 'otra'])

→ FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión                                     | Paso afectado |
|---------------------------------------------------------|---------------|
| Nombre exacto del campo departmentId dentro de allCities | PASO 1        |
