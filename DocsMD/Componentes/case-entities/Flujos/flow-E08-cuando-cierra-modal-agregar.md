━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando cierra el modal de agregar entidad
   Tipo: User Interaction
   Función: cerrarModalAgregar()
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  (ninguno — se dispara desde el botón "Cancelar", la "×", o el click en el overlay)
}

PASO 1 — Ocultar el modal
  modalAgregarAbierto = false

PASO 2 — Descartar cualquier cambio no guardado en el formulario
  (no se persiste nada — el próximo evento E05 vuelve a resetear formNuevaEntidad de todas formas)

→ TERMINAR ejecución (no hay llamadas a sistemas externos, no se recarga la lista)

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión                                    | Paso afectado |
|----------------------------------------------------------|---------------|
| Ninguno — evento trivial, sin dependencias externas.      | —             |
