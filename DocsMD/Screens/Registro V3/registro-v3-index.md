# `Registro V3` — Índice

## Navegación

| MD | Contenido |
|---|---|
| [registro-v3-interface.md](registro-v3-interface.md) | Árbol de interfaz completo |
| [registro-v3-data.md](registro-v3-data.md) | Decisiones de arquitectura, columnas/enum nuevos, GAPS |

### Otros MDs

| MD | Contenido |
|---|---|
| [registro-v3-secciones.md](registro-v3-secciones.md) | Catálogo completo de preguntas por módulo (CI/DP/CA/HV/AG/TZ/GS) |
| [registro-v3-parrafo-sintesis.md](registro-v3-parrafo-sintesis.md) | Reglas completas del párrafo síntesis auto-generado (módulo GS) |

### Changelogs

_Ninguno aún — primera versión de la planeación._

---

## Resumen de eventos

| Evento | Tipo | Flujo |
|---|---|---|
| Cuando carga la pantalla de Registro | Lifecycle | 📄 [Ver flujo](Flujos/flow-E01-cuando-carga-la-pantalla.md) |
| Cuando responde la autorización de datos personales | User Interaction | 📄 [Ver flujo](Flujos/flow-E02-cuando-responde-habeas-data.md) |
| Cuando selecciona el tipo de contacto | User Interaction | 📄 [Ver flujo](Flujos/flow-E03-cuando-selecciona-tipo-contacto.md) |
| Cuando busca en el Cabezote | User Interaction | 📄 [Ver flujo](Flujos/flow-E04-cuando-busca-en-cabezote.md) |
| Cuando presiona "Limpiar" en el Cabezote | User Interaction | 📄 [Ver flujo](Flujos/flow-E05-cuando-limpia-busqueda.md) |
| Cuando presiona "Nuevo Registro" | User Interaction | 📄 [Ver flujo](Flujos/flow-E06-cuando-presiona-nuevo-registro.md) |
| Cuando presiona "Proceder al formulario" | User Interaction | 📄 [Ver flujo](Flujos/flow-E07-cuando-procede-al-formulario.md) |
| Cuando presiona "Nuevo Caso en Registro existente" | User Interaction | 📄 [Ver flujo](Flujos/flow-E08-cuando-presiona-nuevo-caso-en-registro.md) |
| Cuando presiona "Actualizar Caso" | User Interaction | 📄 [Ver flujo](Flujos/flow-E09-cuando-presiona-actualizar-caso.md) |
| Cuando cambia de módulo en el sidebar | User Interaction | 📄 [Ver flujo](Flujos/flow-E10-cuando-cambia-de-modulo.md) |
| Cuando actualiza una respuesta del formulario | User Interaction | 📄 [Ver flujo](Flujos/flow-E11-cuando-actualiza-una-respuesta.md) |
| Cuando guarda una sección/módulo | User Interaction | 📄 [Ver flujo](Flujos/flow-E12-cuando-guarda-una-seccion.md) |
| Cuando presiona "Alerta inminente" | User Interaction | 📄 [Ver flujo](Flujos/flow-E13-cuando-presiona-alerta-inminente.md) |
| Cuando abandona la pantalla sin finalizar | User Interaction (botón "Finalizar llamada") + Lifecycle (respaldo) | 📄 [Ver flujo](Flujos/flow-E14-cuando-abandona-sin-finalizar.md) |
| Cuando presiona "Finalizar registro y enrutar a Seguimiento" | User Interaction | 📄 [Ver flujo](Flujos/flow-E15-cuando-finaliza-y-enruta-a-seguimiento.md) |
| Cuando presiona un acceso de Rutas/Barreras/Notificaciones | User Interaction | 📄 [Ver flujo](Flujos/flow-E16-cuando-presiona-rutas-barreras-notificaciones.md) |
| Cuando agrega o elimina una Entidad | User Interaction | 📄 [Ver flujo](Flujos/flow-E17-cuando-agrega-o-elimina-una-entidad.md) |
| Cuando agrega o elimina una Barrera | User Interaction | 📄 [Ver flujo](Flujos/flow-E18-cuando-agrega-o-elimina-una-barrera.md) |

---

## Inventario de eventos

**Evento:** Cuando carga la pantalla de Registro
**Tipo:** Lifecycle
**Descripción:** Monta la pantalla en `stage = 'habeasData'`. Si llega con un `icode` de Caso existente en la URL (p. ej. desde un link del Expediente u otra pantalla), salta directo a `stage = 'formulario'` con ese Caso cargado.
**Requerido:** Sí

---

**Evento:** Cuando responde la autorización de datos personales
**Tipo:** User Interaction
**Descripción:** Si responde Sí, avanza a `stage = 'tipoContacto'`. Si responde No, pasa a `stage = 'formularioCerrado'` y no se crea ningún registro ni se guarda ninguna respuesta.
**Requerido:** Sí

---

**Evento:** Cuando selecciona el tipo de contacto
**Tipo:** User Interaction
**Descripción:** Si elige "Víctima - Caso Nuevo", avanza a `stage = 'cabezote'`. Si elige "Víctima - Seguimiento" o "Tercero - Reporte", redirige a los formularios correspondientes (fuera de alcance de esta pantalla).
**Requerido:** Sí

---

**Evento:** Cuando busca en el Cabezote
**Tipo:** User Interaction
**Descripción:** Busca por número de documento y/o número de teléfono — **no por nombre** (corregido por diagrama del cliente, ver D7 en `registro-v3-data.md`). El nombre se muestra en los resultados pero no participa en la búsqueda. Clasifica los resultados en Reportes pendientes, Registros con sus Casos asociados, y Borradores.
**Requerido:** Sí

---

**Evento:** Cuando presiona "Limpiar" en el Cabezote
**Tipo:** User Interaction
**Descripción:** Resetea los criterios de búsqueda y oculta los resultados.
**Requerido:** No

---

**Evento:** Cuando presiona "Nuevo Registro"
**Tipo:** User Interaction
**Descripción:** Solo visible cuando la búsqueda no tuvo coincidencias. Crea un `victim_case` nuevo en estado Borrador prediligenciado con los criterios de búsqueda, y abre el Formulario en el módulo CI.
**Requerido:** Sí

---

**Evento:** Cuando presiona "Proceder al formulario"
**Tipo:** User Interaction
**Descripción:** Retoma un `victim_case` existente en Borrador. Deja trazabilidad del agente y la fecha de retome. Abre el Formulario en el último módulo con progreso.
**Requerido:** Sí

---

**Evento:** Cuando presiona "Nuevo Caso en Registro existente"
**Tipo:** User Interaction
**Descripción:** Crea un `victim_case` nuevo, copia (no referencia) el Perfil de víctima (CI+DP+CA) del `victim_case` más reciente con el mismo documento, genera nuevo Caso en Borrador y abre el Formulario en el módulo HV.
**Requerido:** Sí

---

**Evento:** Cuando presiona "Actualizar Caso"
**Tipo:** User Interaction
**Descripción:** Solo visible sobre Casos en estado Iniciado. Reabre el mismo `victim_case` (misma tríada víctima-hecho-agresor) para registrar un nuevo episodio, con trazabilidad.
**Requerido:** Sí

---

**Evento:** Cuando cambia de módulo en el sidebar
**Tipo:** User Interaction
**Descripción:** Valida si el módulo destino está desbloqueado (TZ/GS requieren tríada completa). Si está bloqueado, muestra el tooltip explicativo y no navega. Si está habilitado, cambia `moduloActual` y el `<dinamic-form>` renderiza la sección correspondiente.
**Requerido:** Sí

---

**Evento:** Cuando actualiza una respuesta del formulario
**Tipo:** User Interaction
**Descripción:** El evento central del formulario — guarda el valor en el estado local y recalcula, dentro de la misma ejecución: condicionales de visibilidad, auto-sugerencias (nacionalidad, tratamiento según pronombres), barra de avance del módulo, habilitación secuencial de TZ/GS, el resultado del Tamizaje (tipo por identidad/orientación + nivel de riesgo por score) si el campo es de DP o TZ, el párrafo síntesis (módulo GS), validación "solo letras" en campos de nombre, verificación de duplicidad si el campo es el número de documento, y las alertas de riesgo/ideación suicida. Ver el flujo para el detalle completo con sub-flujos.
**Requerido:** Sí

---

**Evento:** Cuando presiona un acceso de Rutas/Barreras/Notificaciones
**Tipo:** User Interaction
**Descripción:** Agregado tras re-leer el PDF de observaciones del cliente (obs. #17) — el mockup original tiene un 3er grupo de sidebar "Registro de gestión" con accesos a Rutas/Barreras/Notificaciones que esta planeación había omitido. **Parcialmente resuelto** tras revisar el formulario real de Seguimiento (D9): el registro de Entidades y Barreras de ESTE Caso ya vive dentro de GS (subsecciones C y D); estos accesos del sidebar quedan para la vista histórica/agregada en el Expediente de víctima (diferido, D2). Sigue sin confirmar si "Notificaciones" tiene un destino propio.
**Requerido:** No

---

**Evento:** Cuando agrega o elimina una Entidad
**Tipo:** User Interaction
**Descripción:** Repeater dentro de GS (Registrar Entidades) — agrega o quita una entrada de "Entidades Acudidas" o "Entidades Activación", y recalcula la cascada Sector→Entidad de cada ítem. Modela el repeater real de Seguimiento (D9).
**Requerido:** No

---

**Evento:** Cuando agrega o elimina una Barrera
**Tipo:** User Interaction
**Descripción:** Repeater dentro de GS (Registrar Barreras) — agrega o quita una entrada, y muestra el bloque de barreras correspondiente al sector elegido (Salud/Justicia/Protección/Otras). Modela el repeater real de Seguimiento (D9).
**Requerido:** No

---

**Evento:** Cuando guarda una sección/módulo
**Tipo:** User Interaction
**Descripción:** Persiste las respuestas del módulo actual vía `POST /api/v1/forms/saveSection`. En el primer guardado de cualquier Caso nuevo, transiciona su estado de Borrador a Iniciado.
**Requerido:** Sí

---

**Evento:** Cuando presiona "Alerta inminente"
**Tipo:** User Interaction
**Descripción:** Registra en la bitácora del Caso la activación de una alerta inmediata. Proveniente del prototipo React de la ronda anterior, no de las fuentes de Google Drive — criterios de activación y destino de la notificación quedan provisionales (ver flujo).
**Requerido:** No

---

**Evento:** Cuando abandona la pantalla sin finalizar
**Tipo:** Lifecycle
**Descripción:** Antes de desmontar la pantalla (navegación a otra ruta, cierre de pestaña detectable), asegura que el último estado del Caso (Borrador o Iniciado) quedó guardado.
**Requerido:** Sí

---

**Evento:** Cuando presiona "Finalizar registro y enrutar a Seguimiento"
**Tipo:** User Interaction
**Descripción:** Valida que todas las obligatorias [1er contacto] y [máx. 2º contacto] de los 7 módulos estén diligenciadas. Cambia el Caso a Terminado, genera la Denuncia fácil si corresponde, confirma la remisión a procesos especializados si aplica, y traspasa automáticamente Perfil + tríada + resultado de Tamizaje + acciones al Módulo de Seguimiento.
**Requerido:** Sí

---

## Checklist de validación de completitud

- [x] ¿Puede manejarse cada acción del usuario? Sí — 18 eventos cubren Habeas data, Cabezote (5 acciones), navegación/edición del formulario, Rutas/Barreras/Notificaciones, los repeaters de Entidades y Barreras, y cierre.
- [x] ¿Está cubierta la carga inicial de datos? Sí — E01.
- [x] ¿Están incluidos todos los envíos de formulario? Sí — E12 (por sección) y E15 (final).
- [ ] ¿Hay actualizaciones en tiempo real (sockets, notificaciones)? No aplica a esta pantalla — las alertas (riesgo, duplicidad) son parte de E11, no push externo.
- [ ] *(Backend)* ¿Hay tareas programadas? No aplica — esta pantalla no tiene cron jobs propios.
- [ ] *(Backend)* ¿Hay webhooks? No aplica — no hay integraciones externas que empujen datos a esta pantalla.
