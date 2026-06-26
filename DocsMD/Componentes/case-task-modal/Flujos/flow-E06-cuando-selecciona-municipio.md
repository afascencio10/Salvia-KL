━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando selecciona municipio
   Tipo: User Interaction
   Función: onSelectMunicipio()
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Aplica a: formularios 'gestion_llamada' y 'proyectar_oficio' únicamente.

INPUT: {
  form.municipioId:  town_code del municipio seleccionado  → <select> Municipio
                     ej: "11001000"
}


PASO 1 — Resetear entidad

  form.entidadId     = null
  form.entidadNombre = ""
  entidades          = []


PASO 2 — Cargar sedes de entidades del municipio

  GET /api/v1/entity-branches?town_code={form.municipioId}
  → resultado: [{ id, icode, name }]
    id   = integer PK de salvia.entity_branch  (se usa como value del <option>)
    name = nombre de la sede visible en el dropdown

  SI error:
    → Log advertencia
    → entidades = []  (el dropdown queda vacío pero habilitado)

  SI ok:
    → entidades = resultado de la llamada


PASO 3 — Agregar opción fija "Otra entidad" al final

  entidades.push({ id: 'otra', name: 'Otra entidad' })


PASO 4 — Vue re-evalúa reactivamente

  → Dropdown Entidad muestra las sedes del municipio + "Otra entidad" al final
    (:disabled si !form.municipioId → ahora municipioId tiene valor)
  → Campo "Nombre de la entidad" permanece oculto
    ([v-if form.entidadId === 'otra'] → aún no seleccionado)

→ FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión                                    | Paso afectado |
|--------------------------------------------------------|---------------|
| Confirmar que el value del <option> Entidad sea        | PASO 2        |
| el id integer (PK) y no el icode                       |               |
