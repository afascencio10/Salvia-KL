━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando cierra el modal de detalle
   Tipo: User Interaction
   Función: cerrarModal()
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  // Ninguno — el usuario dispara este evento de 4 formas distintas, todas
  // llaman al mismo método sin argumentos:
}

Disparadores posibles:
  • Click en el botón "×" del header del modal
  • Click en el botón "Cerrar" del footer del modal
  • Click en el overlay oscuro (fuera del ModalBox) — @click.self
  • Tecla Escape presionada mientras el modal está abierto (listener en window, agregado en mounted y removido en unmounted)

PASO 1 — Limpiar el estado del modal

  modal = null

PASO 2 — Re-renderizar

  → Vue re-evalúa `v-if modal` → false
  → El ModalBox y el overlay se desmontan

→ FIN EJECUCIÓN ✓
