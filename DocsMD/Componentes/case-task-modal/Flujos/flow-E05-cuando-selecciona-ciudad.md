━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando selecciona ciudad
   Tipo: User Interaction
   Función: onSelectCiudad()
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Aplica a: formularios 'gestion_llamada' y 'proyectar_oficio' únicamente.

INPUT: {
  form.ciudadId:  ID de la ciudad seleccionada  → <select> Ciudad
}


PASO 1 — Resetear en cascada los niveles inferiores

  form.municipioId   = null
  form.entidadId     = null
  form.entidadNombre = ""
  municipios         = []
  entidades          = []


PASO 2 — Cargar municipios de la ciudad seleccionada

  GET /api/v1/locations/towns?city_id={form.ciudadId}
  → municipios = [{ id, name }]
  // id es el town_code (DIVIPOLA), ej: "11001000"

  SI error:
    → Log advertencia
    → municipios = []  (el dropdown queda vacío pero habilitado)

  SI ok:
    → municipios disponibles en el estado del componente


PASO 3 — Vue re-evalúa reactivamente

  → Dropdown Municipio muestra los municipios cargados y queda habilitado
    (:disabled si !form.ciudadId → ahora ciudadId tiene valor)
  → Dropdown Entidad permanece deshabilitado (:disabled si !form.municipioId)
  → Campo "Nombre de la entidad" permanece oculto

→ FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión                                    | Paso afectado |
|--------------------------------------------------------|---------------|
| Shape exacto de la respuesta de /locations/towns       | PASO 2        |
