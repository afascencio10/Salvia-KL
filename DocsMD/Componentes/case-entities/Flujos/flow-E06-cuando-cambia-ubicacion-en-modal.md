━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando cambia ubicación o sector en el modal
   Tipo: User Interaction
   Funciones: onCambiaDepartamento() · onCambiaCiudad() · onCambiaMunicipio() · onCambiaSector()
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  nivel:   cuál select cambió            → 'departamento' | 'ciudad' | 'municipio' | 'sector'
  valor:   el nuevo valor seleccionado    → v-model del <select> correspondiente
}

PASO 1 — Determinar qué nivel cambió y resetear los niveles dependientes (en cascada)

SEGÚN nivel:
  CASO 'departamento':
    → formNuevaEntidad.ciudad = ''
    → formNuevaEntidad.municipio = ''
    → entidadesDisponibles = []
    → CONTINÚA

  CASO 'ciudad':
    → formNuevaEntidad.municipio = ''
    → entidadesDisponibles = []
    → CONTINÚA

  CASO 'municipio':
    → CONTINÚA (no resetea nada más, es el nivel final de ubicación)

  CASO 'sector':
    → CONTINÚA (no afecta la cascada de ubicación)

PASO 2 — Verificar si ya hay municipio Y sector seleccionados (ambos obligatorios)

SI formNuevaEntidad.municipio está vacío O formNuevaEntidad.sector está vacío:
  → No se hace ninguna llamada — el <select> de Entidad permanece vacío/oculto
  → TERMINAR ejecución

SI ambos tienen valor:
  → cargandoEntidades = true
  → CONTINÚA FLUJO GENERAL

PASO 3 — Consultar sedes disponibles en esa ubicación y sector
  EntityBranch.buscar({
    town_code: formNuevaEntidad.municipio,     // código del municipio elegido
    sector:    formNuevaEntidad.sector          // obligatorio — ya no es opcional
  })

  [MÉTODO] GET /api/v1/entity-branches
  // filter:
  {
    "town_code": formNuevaEntidad.municipio,   // origen: cascada de ubicación del modal
    "sector":    formNuevaEntidad.sector        // origen: <select> sector del modal (obligatorio)
  }

  → resultado → entidadesDisponibles = resultado
  → cargandoEntidades = false

SI entidadesDisponibles.length === 0:
  → Muestra EmptyState "No hay sedes registradas en este municipio y sector. Puedes crear una desde Sedes."
  → TERMINAR ejecución (usuario no puede continuar sin elegir una entidad)

SI NO:
  → Muestra <select> Entidad con las opciones
  → CONTINÚA FLUJO GENERAL (usuario puede ahora elegir una entidad y luego escribir el objetivo)

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión                                                                          | Paso afectado |
|------------------------------------------------------------------------------------------------|---------------|
| Qué selector en cascada de departamento/ciudad/municipio reutilizar exactamente (probablemente el mismo que usa el registro de caso) — confirmar componente exacto antes de implementar. | PASO 1 |
| Si no existe ninguna sede en el municipio elegido, el flujo no ofrece "crear sede nueva desde aquí" — el agente debe salir al módulo de Sedes. A confirmar si se quiere una vía más rápida. | PASO 3 |
