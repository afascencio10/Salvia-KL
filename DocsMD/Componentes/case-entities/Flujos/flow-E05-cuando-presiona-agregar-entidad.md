━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando presiona "+ Agregar entidad"
   Tipo: User Interaction
   Función: abrirModalAgregar()
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  (ninguno — el botón no depende de datos previos)
}

PASO 1 — Resetear el formulario del modal
  formNuevaEntidad = { departamento: '', ciudad: '', municipio: '', sector: '', entidadId: '', objetivo: '' }

PASO 2 — Resetear estado auxiliar
  entidadesDisponibles = []
  cargandoEntidades = false
  errores = {}
  guardando = false

PASO 3 — Mostrar el modal
  modalAgregarAbierto = true

→ Ver flujo: Cuando cambia ubicación o sector en el modal (se dispara después, cuando el usuario
  empieza a elegir departamento/ciudad/municipio)

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión                                    | Paso afectado |
|----------------------------------------------------------|---------------|
| Ninguno — evento trivial, sin llamadas a sistemas externos. | —          |
