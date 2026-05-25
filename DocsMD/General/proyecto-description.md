# SOG Salvia — Descripción General

SOG Salvia (Sistema de Orientación y Gestión) es una plataforma web para organizaciones que atienden casos de violencia basada en género (VBG). Permite registrar víctimas, gestionar casos, realizar seguimientos programados, identificar barreras institucionales e intervenir con medidas de protección, apoyo psicosocial y estabilización económica.

El flujo principal comienza cuando un operador registra un caso: completa un formulario con datos de la víctima, los hechos, el agresor y un tamizaje de riesgo (sí/no) que el backend recalcula para asignar un nivel (bajo, moderado, alto o extremo). Al guardar, se genera automáticamente un calendario de seguimientos según el nivel de riesgo y se asigna un agente de forma balanceada entre el equipo.

Los agentes realizan los seguimientos desde la pantalla **Hacer Seguimiento**, que carga un formulario dinámico por secciones. Cada sección registra si la víctima visitó la entidad, si recibió atención y qué barreras encontró. Las barreras activas se gestionan desde **Mis Barreras**, donde el agente puede enviar oficios a las entidades y registrar las respuestas. El historial completo del caso se visualiza en el componente **Timeline**.

El sistema también gestiona casos de **feminicidio** y **riesgo de feminicidio** con formularios propios.

**Tecnologías:** Go (Gin) en el backend, Vue 3 (cargado vía `vue3-sfc-loader` dentro de templates Go), PostgreSQL en Supabase. Flutter para la aplicación móvil.

**Roles:** Supervisor (`sv`), Operador (`op`), Operador de Riesgo (`ro`), Administrador (`ad`), Operador de Feminicidio (`fo`), Operadores Territoriales (`no`, `do`), Entidad (`et`), Usuario externo (`us`).

---

## Índice de pantallas y módulos

| Pantalla / Módulo | Descripción |
|---|---|
| `Registro de Caso` | Formulario para registrar un caso VBG: datos de víctima, hechos, agresor y tamizaje de riesgo. |
| `Detalle del Caso` | Vista completa del caso con timeline, barreras, seguimientos, oficios e intervenciones activas. |
| `Hacer Seguimiento` | Pantalla donde el agente completa el formulario dinámico de un seguimiento programado. |
| `Mis Seguimientos` | Vista del agente con sus seguimientos del día: pendientes, prioritarios y completados. |
| `Seguimientos Área` | Vista del supervisor con todos los seguimientos del área organizados por estado y equipo. |
| `Mis Barreras` | Lista de barreras activas asignadas al agente con acciones para gestionar oficios y respuestas. |
| `Barrera Detalle` | Vista interna de una barrera: historial de actuaciones, oficios enviados y respuestas recibidas. |
| `Notificaciones` | Centro de notificaciones del sistema para alertas de casos y seguimientos. |
| `Registro de Contacto Víctima` | Formulario de ingreso inicial de la víctima como contacto antes de abrir un caso formal. |
| `Feminicidio` | Listado y registro de casos de feminicidio con formulario de datos de la víctima y familia. |
| `Riesgo de Feminicidio` | Listado y registro de casos de riesgo de feminicidio con formulario propio. |
| `Entidades` | Catálogo de entidades institucionales que participan en la ruta de atención. |
| `Sedes` | Gestión de sedes (branches) de entidades con ubicación geográfica y municipio. |
| `Asignar Operadores` | Herramienta para asignar manualmente operadores a casos. |
| `Gestión de Usuarios` | Creación y configuración de usuarios, roles y equipos del sistema. |
| `dinamic-form` | Componente reutilizable que renderiza formularios dinámicos por secciones con guardado parcial. |
| `case-timeline` | Componente reutilizable que muestra el historial cronológico de eventos de un caso. |
