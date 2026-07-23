━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando presiona "Ver Detalle" de una entidad
   Tipo: User Interaction
   Función: onVerDetalle(entidad)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  entidad:  objeto de la card presionada   → ya está en memoria (viene de `entidades`, cargado en E01)
}

PASO 1 — Navegar directamente a la pantalla de detalle de entidad
  window.location.href = '/salvia/entidad/' + entidad.entityBranchICode

→ A diferencia del diseño original, el componente ya no emite un evento al padre —
  navega él mismo, por decisión del usuario.

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión                                                                                          | Paso afectado |
|------------------------------------------------------------------------------------------------------------------|---------------|
| La ruta `/salvia/entidad/:id` no existe hoy en el código real (solo existe `/salvia/sedes`, el listado de sedes). Falta crear esa pantalla/facade — queda fuera del alcance de este componente pero es una dependencia directa de este evento. | PASO 1 |
