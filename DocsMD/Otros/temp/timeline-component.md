# Timeline Component

Componente Vue reutilizable (`case-timeline`) que muestra el historial de eventos de un caso en orden cronológico descendente.

Recibe como props `caseId` (requerido) y `barrierId` (opcional para filtrar por barrera). Hace su propio fetch a `GET /api/v1/casos/:id/timeline-events`, resuelve el nombre del actor via JOIN a la tabla de usuarios y renderiza cada evento con su ícono Font Awesome, tipo, fecha y descripción.

Incluye filtros pill por categoría: General, Barreras, Seguimientos, Oficios, Medidas, Psicosocial y Estabilización.

> Doc completa: `DocsMD/Componentes/Timeline component/TimeLine Comp usage.md`
