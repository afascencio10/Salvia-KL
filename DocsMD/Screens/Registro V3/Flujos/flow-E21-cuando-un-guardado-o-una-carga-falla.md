━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando un guardado o una carga falla
   Tipo: User Interaction (manejo de errores)
   Estado: ✅ Implementado — 2026-11-02 (WP-2 · REG-11)
   Funciones: leerRespuesta() · mensajeErrorGuardado() · mostrarErrorGuardado() · enviarSeccion() · guardarModulo() · guardarYSalir() · guardarTamizaje() · abrirCaso() · crearNuevoRegistro() · ejecutarBusqueda() → registro-v3-app.js · SaveSection() → victim_case_v3_controller.go
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

> Origen: comentario de Felipe en la planeación de WP-2 — «asegúrate de que si algo falla el usuario lo sepa; si es por el
> formato del dato, pedirle al usuario que lo corrija». Antes un 500 al guardar una sección se leía como JSON sin
> `triadaCompleta` y la pantalla igual decía «Sección guardada.»; «Finalizar llamada» ignoraba el error y salía; abrir o crear
> un registro que fallaba no mostraba nada; y el Tamizaje dejaba `{error}` como «resultado».

INPUT: {
  respuesta:  HTTP status + cuerpo JSON de la petición que falló   → fetch de buscar, crear borrador, abrir Caso, guardar sección o guardar Tamizaje
}

PASO 1 — Clasificar la respuesta
  leerRespuesta(r): lee el cuerpo como texto y lo intenta parsear; si no es JSON devuelve {} (no lanza)  → { ok, status, data }
  SI ok: sigue el flujo normal del evento que la disparó (E04, E06, E07, E11, E12)
  SI NO: PASO 2

PASO 2 — Texto para la persona (mensajeErrorGuardado, para guardar una sección)
  SEGÚN status:
    400 con `detalle` (dato inválido: formato, tamaño u opción):
        "No se guardó la sección. Corrige «<pregunta>»: <detalle>. Luego vuelve a guardar."
        // <pregunta> = etiqueta de camposPorId[data.campo]; el detalle empieza con mayúscula y termina en punto
        // sin `campo`: "No se guardó la sección. Corrige lo que escribiste: <detalle>. Luego vuelve a guardar."
        // ej.: «Corrige «¿A qué hora ocurrió el hecho?»: Escribe la hora en formato de 24 horas HH:MM (por ejemplo 19:30).»
        //      «Corrige «¿Cuál es la dirección donde ocurrieron los hechos…?»: Este dato admite máximo 150 caracteres y escribiste 151.»
    401 / 403:
        "No se guardó la sección: tu sesión venció o no tienes permiso. Vuelve a iniciar sesión; lo que escribiste en esta pantalla no se guardó."
    cualquier otro (500…):
        "No se pudo guardar la sección por un error del servidor. Lo que escribiste sigue en pantalla: intenta de nuevo y, si el problema continúa, avisa a soporte."
  La petición lanza (sin conexión):
        "No se pudo guardar la sección: no hay conexión con el servidor. Lo que escribiste sigue en pantalla; intenta de nuevo."
  Tamizaje (guardarTamizaje): 401/403 → "No se guardó el Tamizaje: tu sesión venció…"; 500 → "No se pudo guardar el Tamizaje por un error del servidor…"; sin conexión → "No se pudo guardar el Tamizaje: no hay conexión con el servidor…"

PASO 3 — Mostrar el aviso (mostrarErrorGuardado)
  errorGuardado = texto del PASO 2
  mensajeGuardado = "No se guardó."                       // junto al botón "Guardar sección"; nunca "Sección guardada." si no se guardó
  Vue muestra el aviso rojo (.alert-danger #registro-v3-error-guardado) bajo las pestañas, y se hace scroll hasta él
  El estado local NO se toca: lo escrito sigue en pantalla para corregirlo y volver a guardar

PASO 4 — Limpiar el aviso
  errorGuardado = '' al: empezar otro guardado · guardar bien · cambiar de módulo (irAModulo) · abrir un Caso (abrirCaso) · volver al Cabezote (volverACabezote)

OTROS PUNTOS DEL MISMO EVENTO

  Cabezote (errorBusqueda, aviso .alert-danger):
    ejecutarBusqueda:    !ok → data.error o "Error al buscar."
    crearNuevoRegistro:  !ok o sin icode → "No se pudo crear el registro. Intenta de nuevo; si el problema continúa, avisa a soporte."  (sin conexión: "…no hay conexión con el servidor…")
    abrirCaso:           !ok o sin `summary` → "No se pudo abrir el registro. Intenta de nuevo; si el problema continúa, avisa a soporte."  (flow-E07)

  "Finalizar llamada (guardar y salir)" (flow-E14), si el guardado falla:
    mostrarErrorGuardado() y window.confirm("<texto>\n\n¿Salir de todos modos? Lo que escribiste en este módulo se perdería.")
      Cancelar → se queda en el Formulario · Aceptar → volverACabezote()

  "Finalizar registro y enrutar a Seguimiento" (flow-E15): ya tenía su propio aviso (errorFinalizar, con los módulos faltantes).

  Servidor (SaveSection): cada rechazo deja una línea en el log — «[RegistroV3] SaveSection rechazado icode=… módulo=… campo=…: <motivo>» (400) o
    «[RegistroV3] SaveSection falló icode=… módulo=…: <error>» (500)

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión | Paso afectado |
|---|---|
| El respaldo `sendBeacon` de `beforeunload` (flow-E14) no puede mostrar nada porque la pestaña se está cerrando: si ese guardado es rechazado, se pierde sin aviso | — |
| El aviso nombra la pregunta pero no la resalta en el formulario (sin `is-invalid` por campo): evita tocar líneas del template compartidas con otras entregas; se puede agregar después | PASO 3 |
| Entidades y barreras (flow-E17, flow-E18) tienen sus propios avisos y no se revisaron en esta ronda | — |
