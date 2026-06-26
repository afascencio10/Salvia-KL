━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando carga el componente
   Tipo: Lifecycle
   Función: mounted()
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  (ninguno — el componente no recibe props ni config al montarse)
}

PASO 1 — Inicializar estado del componente

  visible          = false
  tarea            = null
  taskId           = null
  form             = {}
  saveError        = null
  errorTarea       = null
  cargandoTarea    = false
  guardando        = false
  allDepartamentos = []
  allCities        = []
  municipios       = []
  entidades        = []


PASO 2 — Cargar datos de ubicación en paralelo

  Los dropdowns Departamento → Ciudad → Municipio → Entidad
  solo aparecen en los formularios 'gestion_llamada' y 'proyectar_oficio'.
  El tipo de tarea no se conoce hasta que open(taskId) sea llamado,
  por lo que los datos de ubicación se pre-cargan siempre al montar.

  [Llamada A] GET /api/v1/locations/departments
  → allDepartamentos = [{ id, name }]

  [Llamada B] GET /api/v1/locations/cities
  → allCities = [{ id, name, departmentId }]

  Ambas llamadas corren en paralelo (Promise.all o similar).

  SI error en cualquier llamada:
    → Log advertencia (no bloquea el montaje)
    → Los dropdowns de ubicación quedarán vacíos si el usuario abre una tarea de ese tipo

  SI ok:
    → allDepartamentos y allCities disponibles en el estado del componente
    → Vue pre-pobla el dropdown de Departamento cuando el modal se abra

→ FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión                                             | Paso afectado |
|-----------------------------------------------------------------|---------------|
| Endpoints exactos para cargar departamentos y ciudades          | PASO 2        |
| ¿Se cargan en una llamada combinada o en dos llamadas separadas?| PASO 2        |
