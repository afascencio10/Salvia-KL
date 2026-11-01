> ⛔ **No implementado (verificado en el código el 2026-10-31).** El Formulario no tiene el grupo "Registro de gestión" (Rutas / Barreras / Notificaciones): fuera de alcance de la ronda ("Notificaciones" se ignoró por indicación del usuario el 2026-10-15). El registro de entidades y barreras de ESTE Caso ocurre en GS § C/D (flow-E17, flow-E18). Se conserva el diseño original por si se retoma.

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando presiona un acceso de Rutas/Barreras/Notificaciones
   Tipo: User Interaction
   Función: abrirRutas() / abrirBarreras() / abrirNotificaciones() → set_victim_case_v3.html
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

> Evento agregado al releer `Observaciones y ajustes al MockUp de registro
> 2.0.pdf` (obs. #17). El mockup original incluye un 3er bloque de
> navegación en el sidebar — "3. Registro de gestión" con sub-ítems a.Rutas,
> b.Barreras, c.Notificaciones — paralelo a los bloques de Víctima y Caso.
> La primera versión de este plan lo había omitido por completo.
>
> El propio cliente pregunta en su observación: *"Son esas preguntas [las de
> Gestión Salvia] o los botones que aparecen en '3. Registro de gestión' los
> que conectan con el resto de los módulos? Precisar"*
>
> **Parcialmente resuelto** tras revisar en BD el formulario real de
> Seguimiento (`registro-v3-data.md` § D9): GS ya tiene su propio registro
> estructurado de Entidades y Barreras **para este Caso** (subsecciones C y
> D de `registro-v3-secciones.md`, repeaters — ver `flow-E17`/`flow-E18`),
> replicando el mismo modelo que usa Seguimiento. Estos botones del sidebar
> ya no necesitan ser el mecanismo de REGISTRO — tiene más sentido que
> naveguen a una vista de consulta/histórico (Expediente de víctima,
> diferido en D2) en vez de duplicar el registro. "Notificaciones" sigue
> sin un destino claro — no hay un equivalente real encontrado en Seguimiento.

INPUT: {
  icode:   identificador del Caso actual   → estado de la pantalla
  destino: 'rutas' | 'barreras' | 'notificaciones'   → qué botón se presionó
}

PASO 1 — Navegar al módulo correspondiente

SEGÚN destino:
  CASO 'rutas':
    → router.push(`/salvia/rutas?icode=${icode}`)  // vista histórica de rutas/entidades, fuera de esta pantalla — el registro de ESTE Caso ya ocurre en GS § C
  CASO 'barreras':
    → router.push(`/salvia/barreras?icode=${icode}`)  // vista histórica de barreras, ya existe (`barrera_detalle.html`) — el registro de ESTE Caso ya ocurre en GS § D
  CASO 'notificaciones':
    → router.push(`/salvia/notificaciones?icode=${icode}`)  // módulo de Notificaciones/oficios, fuera de esta pantalla

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión                                                                     | Paso afectado |
|------------------------------------------------------------------------------------------------|---------------|
| G-V3-12 — Parcialmente resuelto (D9): el registro ya ocurre en GS § C/D; falta confirmar si "Rutas"/"Barreras" aquí deben ser solo vista histórica (propuesta de este flujo) o si el cliente prefiere otra cosa, y qué es "Notificaciones" | Todo el evento |
| Si "navega" implica salir del Formulario de Registro (perdiendo el módulo actual) o abre una vista superpuesta/modal sin perder el progreso | PASO 1 |
| Rutas de destino (`/salvia/rutas`, `/salvia/notificaciones`) son propuestas — confirmar si ya existen o si estos módulos no están construidos todavía | PASO 1 |
