# SOG Salvia — Descripción General

SOG Salvia (Sistema de Orientación y Gestión) es una plataforma web para organizaciones que atienden casos de violencia basada en género (VBG). Permite registrar víctimas, gestionar casos, realizar seguimientos programados, identificar barreras institucionales e intervenir con medidas de protección, apoyo psicosocial y estabilización económica.

El flujo principal comienza cuando un operador registra un caso: completa un formulario con datos de la víctima, los hechos, el agresor y un tamizaje de riesgo (sí/no) que el backend recalcula para asignar un nivel (bajo, moderado, alto o extremo). Al guardar, se genera automáticamente un calendario de seguimientos según el nivel de riesgo y se asigna un agente de forma balanceada entre el equipo.

Los agentes realizan los seguimientos desde la pantalla **Hacer Seguimiento**, que carga un formulario dinámico por secciones. Cada sección registra si la víctima visitó la entidad, si recibió atención y qué barreras encontró. Las barreras activas se gestionan desde **Mis Barreras**, donde el agente puede enviar oficios a las entidades y registrar las respuestas. El historial completo del caso se visualiza en el componente **Timeline**.

El sistema también gestiona casos de **feminicidio** y **riesgo de feminicidio** con formularios propios.

**Tecnologías:** Go (Gin) en el backend, Vue 3 (cargado vía `vue3-sfc-loader` dentro de templates Go), PostgreSQL en Supabase. Flutter para la aplicación móvil.

**Roles:** Supervisor (`sv`), Operador (`op`), Operador de Riesgo (`ro`), Administrador (`ad`), Operador de Feminicidio (`fo`), Operadores Territoriales (`no`, `do`), Entidad (`et`), Usuario externo (`us`).

---

## Índice de pantallas y módulos

| Pantalla / Módulo | Descripción | Archivo principal |
|---|---|---|
| [`Registro de Caso`](../Screens/Registro%20de%20Caso/registro-caso-index.md) | Formulario para registrar un caso VBG: datos de víctima, hechos, agresor y tamizaje de riesgo. | `src/frontend/html/salvia/victim_case/set_victim_case.html` |
| `Detalle del Caso` | Vista completa del caso con timeline, barreras, seguimientos, oficios e intervenciones activas. | `src/frontend/html/salvia/case_detail/get_case_detail_sv.html` |
| [`Hacer Seguimiento`](../Screens/hacer-seguimiento/hacer-seguimiento-index.md) | Pantalla donde el agente completa el formulario dinámico de un seguimiento programado. | `src/frontend/html/salvia/follow_up_v2/hacer_seguimiento.html` |
| `Mis Seguimientos` | Vista del agente con sus seguimientos del día: pendientes, prioritarios y completados. | `src/frontend/html/salvia/my_follow_ups/get_my_follow_ups.html` |
| `Seguimientos Área` | Vista del supervisor con todos los seguimientos del área organizados por estado y equipo. | `src/frontend/html/salvia/follow_up_v2/seguimientos_area.html` |
| `Mis Barreras` | Lista de barreras activas asignadas al agente con acciones para gestionar oficios y respuestas. | `src/frontend/html/salvia/barriers/barrieris.html` |
| `Barrera Detalle` | Vista interna de una barrera: historial de actuaciones, oficios enviados y respuestas recibidas. | `src/frontend/html/salvia/barriers/barrera_detalle.html` |
| `Notificaciones` | Centro de notificaciones del sistema para alertas de casos y seguimientos. | `src/frontend/html/salvia/notifications/notificaciones.html` |
| `Registro de Contacto Víctima` | Formulario de ingreso inicial de la víctima como contacto antes de abrir un caso formal. | `src/frontend/html/salvia/victim_contact/set_victim_contact.html` |
| `Feminicidio` | Listado y registro de casos de feminicidio con formulario de datos de la víctima y familia. | `src/frontend/html/salvia/feminicide/set_feminicide.html` |
| `Riesgo de Feminicidio` | Listado y registro de casos de riesgo de feminicidio con formulario propio. | `src/frontend/html/salvia/feminicide/set_feminicide_risk.html` |
| `Entidades` | Catálogo de entidades institucionales que participan en la ruta de atención. | — |
| `Sedes` | Gestión de sedes (branches) de entidades con ubicación geográfica y municipio. | `src/frontend/html/salvia/entity_branch/get_entity_branches.html` |
| `Asignar Operadores` | Herramienta para asignar manualmente operadores a casos. | `src/frontend/html/salvia/victim_case/assign_operators.html` |
| `Gestión de Usuarios` | Creación y configuración de usuarios, roles y equipos del sistema. | `src/frontend/html/security/general_user/set_general_user.html` |

## Índice de componentes reutilizables

| Componente | Descripción | Archivo principal |
|---|---|---|
| `dinamic-form` | Componente reutilizable que renderiza formularios dinámicos por secciones con guardado parcial. | `src/frontend/js/components/dinamic-form.js` |
| `case-timeline` | Componente reutilizable que muestra el historial cronológico de eventos de un caso. | `src/frontend/js/components/case-timeline.js` |
