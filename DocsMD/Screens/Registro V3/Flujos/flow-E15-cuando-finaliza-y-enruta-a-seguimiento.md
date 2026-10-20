━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando presiona "Finalizar registro y enrutar a Seguimiento"
   Tipo: User Interaction
   Función: finalizarRegistro() → dinamic-form.js / set_victim_case_v3.html
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  icode:     identificador del Caso                    → estado de la pantalla
  respuestas: todas las respuestas acumuladas de los 7 módulos → estado local
}

Validar formulario:
  • Todas las obligatorias [1er contacto] Y [máx. 2º contacto] de CI, DP, CA, HV, AG, TZ y GS deben estar respondidas
  • Las [Deseable] quedan opcionales incluso para finalizar
  SI hay obligatorias sin responder:
    → Mostrar error señalando en qué módulo(s) falta diligenciamiento
    → Mostrar toast "No se puede finalizar — hay módulos incompletos"
    → TERMINAR ejecución

  SI validación ok:
    → CONTINÚA FLUJO GENERAL

PASO 1 — Persistir el último módulo (GS) si tiene cambios sin guardar
  → Ver flujo: Cuando guarda una sección/módulo

PASO 2 — Cambiar el estado del Caso a Terminado

[PUT] /api/v1/victim-cases/{icode}/estado
// payload:
{ "estado": "Terminado" }
→ caso.estado = 'Terminado'
→ Nota: el ciclo completo de estados confirmado por el cliente es Borrador →
  Iniciado → Terminado → **Cerrado** (obs. #1 de `Observaciones y ajustes al
  MockUp de registro 2.0.pdf`). Esta pantalla solo llega hasta Terminado —
  la transición a Cerrado es responsabilidad del Módulo de Seguimiento, fuera
  de este flujo.

PASO 3 — Generar la Denuncia fácil si corresponde

SI respuestas.denuncia_virtual === true:
  → SUB-FLUJO: Generar Denuncia fácil
    S1. DenunciaFacilService.generar({
          icode,
          datosVictima:  { CI, DP },      // de respuestas
          datosHechos:   { HV },
          datosAgresor:  { AG }
        })
    S2. → respuesta.idDenuncia → se asocia al Caso
  → FIN SUB-FLUJO → CONTINÚA FLUJO GENERAL

SI NO:
  → CONTINÚA FLUJO GENERAL

PASO 4 — Confirmar remisión a equipos SALVIA si corresponde

> Reestructurado sobre el modelo real de 6 equipos (D9) — ya no son 4
> "procesos" genéricos con un catálogo de criterios pendiente (G-V3-04
> resuelto).

SI respuestas.requiere_remision === true:
  → SUB-FLUJO: Confirmar remisión
    S1. [POST] /api/v1/remisiones
        // payload:
        {
          "icode":   icode,
          "equipos": respuestas.equipos_remision,   // multi — hasta 6: medidas_emergencia, atencion_psico, estabilizacion, atencion_hombres, salvia_dignidad, discapacidad
          "detalle_por_equipo": {
            "medidas_emergencia": respuestas.cuales_medidas_emergencia,      // si aplica
            "atencion_psico":     respuestas.criterios_psicosocial,         // si aplica — incluye el criterio obligatorio + los de puntaje marcados
            "atencion_hombres":   respuestas.criterio_hombres,              // si aplica
            "discapacidad":       respuestas.servicios_discapacidad,        // si aplica
            "estabilizacion":     respuestas.criterios_estabilizacion       // si aplica
          },
          "justificacion": respuestas.elementos_evidencia_remision
        }
    S2. PARA CADA equipo EN respuestas.equipos_remision:
          → Asignar el Caso a la bandeja del equipo especializado correspondiente
          → NotificationService.notificarEquipo({ equipo, icode })
    S3. El estado de cada remisión queda visible en Mis Seguimientos y (cuando exista) en el Expediente de víctima
  → FIN SUB-FLUJO → CONTINÚA FLUJO GENERAL

SI NO:
  → CONTINÚA FLUJO GENERAL

PASO 4b — Confirmar entidades y barreras registradas

SI hay entradas en el repeater de Entidades (Acudidas y/o Activación):
  → [POST] /api/v1/entidades/registrar-lote  { icode, entidades: respuestas.entidadesAcudidas.concat(respuestas.entidadesActivacion) }
  → Cada entrada queda asociada al Caso y visible luego en "Seguimiento a Entidades" del Módulo de Seguimiento

SI hay entradas en el repeater de Barreras:
  → [POST] /api/v1/barreras/registrar-lote  { icode, barreras: respuestas.barreras }
  → Cada entrada queda asociada al Caso y visible luego en "Seguimiento a Barreras" del Módulo de Seguimiento

SI ninguno de los dos repeaters tiene entradas:
  → CONTINÚA FLUJO GENERAL (ninguno es obligatorio)

PASO 5 — Traspasar el Caso al Módulo de Seguimiento

┌─────────────────────────────────────────┐
│  SUB-FLUJO: Handoff a Seguimiento  │
└─────────────────────────────────────────┘

  S1. SeguimientoService.crearDesdeRegistro({
        icode,
        perfilVictima:  { CI, DP, CA },
        triada:         { HV, AG },
        tamizaje:       { tipoRiesgo: victim_case_form2.tamizaje_tipo_riesgo, nivelRiesgo: victim_case_form2.tamizaje_nivel_riesgo },
        gestiones:      respuestas.gestiones_interinstitucionales,   // GS pregunta 2 (antes "acciones_gestion") — opciones realineadas con Seguimiento real en D9
        entidades:      ya persistidas en PASO 4b, Seguimiento las ve directo en su sección "Seguimiento a Entidades"
        barreras:       ya persistidas en PASO 4b, Seguimiento las ve directo en su sección "Seguimiento a Barreras"
      })

  S2. PARA CADA gestion EN respuestas.gestiones_interinstitucionales:
        → Crear una actuación en Seguimiento: { tipo: gestion, responsable: agentId, fecha: now(), estado: 'pendiente' }
      // Orientación y enrutamiento / Gestión administrativa mediante llamada / Activación de ruta interinstitucional / Articulación interinstitucional / Escalamiento a organismos de control — cada una su propia actuación

  S3. Guardar la bitácora inicial del Caso en Seguimiento:
        bitacora.primeraEntrada = {
          texto:       respuestas.descripcion_acompanamiento,   // GS pregunta 4
          agente:      agentId,
          fecha:       now()
        }
      → Esta entrada NUNCA se sobrescribe — las anotaciones posteriores en Seguimiento se acumulan cronológicamente después de ella

  S4. Conectar cada pieza de GS con su funcionalidad real (ya no depende de "medio_gestion", que se eliminó de GS — ver D9):
        SI respuestas.denuncia_virtual === true: → vincula con idDenuncia del PASO 3
        SI hubo entidades registradas (PASO 4b): → cada una trae su propio "Canal de activación" (Notificación/Llamada/Mi Salvia) y ya quedó persistida — Seguimiento solo la hereda, no la vuelve a crear
        SI hubo barreras registradas (PASO 4b): → Seguimiento las hereda directo en su sección "Seguimiento a Barreras", con la misma lógica de repeater (`state_items`) que ya usa en producción

  → FIN SUB-FLUJO → CONTINÚA FLUJO GENERAL

PASO 6 — Mostrar la pantalla de Cierre (agregado del mockup, ver registro-v3-interface.md)
  ses.fase = 'cierre'
  → Se muestra CierreRegistro con los 4 stat cards, tabla de actuaciones, bitácora y remisiones generadas en los pasos anteriores
  → El agente decide desde ahí: "Ver Expediente de víctima" (placeholder, diferido), "Ir al Cabezote de búsqueda" (stage = 'cabezote') o "Mis Reportes" (fuera de alcance)
  → La redirección automática a Seguimiento NO ocurre sola — el acceso al Módulo de Seguimiento queda disponible desde ahí mismo o desde Mis Seguimientos, no es forzoso salir de esta pantalla

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión                                                                        | Paso afectado |
|-------------------------------------------------------------------------------------------------|---------------|
| Endpoint real del Módulo de Seguimiento para recibir el handoff (`SeguimientoService.crearDesdeRegistro` es propuesto) — confirmar si ya existe algo equivalente | PASO 5 |
| ~~Catálogo de criterios de remisión por proceso (G-V3-04)~~ — **resuelto**, ver D9 y `registro-v3-secciones.md` § GS-B | PASO 4 |
| G-V3-21 — tabla exacta donde `dinamic-form` guarda las respuestas de un repeater (no es `victim_case_form2`) — confirmar antes de implementar PASO 4b | PASO 4b |
| G-V3-20 — catálogo exacto de opciones de barreras por sector y categorías transversales, no extraído de BD en este análisis | PASO 4b |
| Qué pasa si `requiere_remision` o `denuncia_virtual` fallan a mitad de camino (ej. el servicio de Denuncia fácil no responde) — ¿se bloquea la finalización completa del Caso o se marca para reintento? | PASO 3, PASO 4 |
| Si el registro de entidades/barreras (PASO 4b) debe bloquear la finalización si falla, igual que la pregunta anterior | PASO 4b |
