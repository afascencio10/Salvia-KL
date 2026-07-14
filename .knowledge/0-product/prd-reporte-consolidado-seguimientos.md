---
okf_version: "1.0"
type: Product_Requirement
title: "PRD: Reporte consolidado de seguimientos (Excel)"
description: "Exportar a Excel el historial consolidado de seguimientos de los casos creados en un rango de fechas, desde la pantalla Lista de Casos, reemplazando el script SQL directo que dejó de ser viable."
owner: "@salvia-product"
status: draft
tags: [feature, reportes]
last_updated: "2026-07-02"
---

# PRD: Reporte consolidado de seguimientos (Excel)

## Problema / Contexto

El equipo de SALVIA necesita consultar la evolución de la atención a las víctimas a lo largo del tiempo: detalle del caso, cada seguimiento realizado, la información capturada en el formulario de cada seguimiento y el timeline de eventos.

Antes esto se obtenía con un **script SQL directo a la base de datos**. Tras el cambio de URL de la infraestructura (ver [Matriz de Entornos](/.knowledge/4-operations/environment-matrix.md)) ese mecanismo **dejó de ser viable**, por lo que se requiere con **urgencia** una funcionalidad de exportación dentro del producto.

## Objetivo

Permitir descargar, desde la pantalla **Lista de Casos**, un archivo **Excel consolidado** con el historial completo de seguimientos de los casos **creados en un rango de fechas** seleccionado por el usuario.

## Alcance (decisiones cerradas con el negocio, 2026-07-02)

- **Un único botón consolidado** en Lista de Casos: "Descargar reporte consolidado". No hay botón por fila.
- Al pulsarlo se abre un **modal con un rango de fechas** (fecha desde / fecha hasta). Es el **único filtro**.
- El rango aplica sobre la **fecha de creación del caso** (`victim_case_creation_date`). Se incluyen los casos creados en ese rango y, para cada uno, **todo su historial** de seguimientos y timeline.
- **Máximo 1 año** de rango (para acotar volumen y tiempo de generación).
- El Excel se organiza **multi-hoja**: Casos, Seguimientos, Respuestas de formulario (formato largo/clave-valor), Timeline.
- El reporte respeta la **visibilidad por rol** del usuario (mismo alcance que ya ve en Lista de Casos).
- **UX Persuasiva (Simulación de progreso)**: Dado que el procesamiento es síncrono y puede tardar de 5 a 15 segundos para volúmenes altos (ej. 10k casos y 80k seguimientos), el loader en frontend simulará el progreso del backend cambiando dinámicamente sus mensajes descriptivos a intervalos regulares para mantener informado al usuario y evitar reintentos.
- **Límite de Escala**: El formato largo de la hoja "Respuestas de formulario" está sujeto al límite físico de Microsoft Excel de 1,048,576 filas. Si el volumen proyectado supera esta escala, el usuario deberá segmentar las descargas en rangos de fechas menores.

## Roles habilitados

Solo **Supervisor (`sv`)** y **Administrador (`ad`)** — datos sensibles de víctimas (PII). Ver [Matriz RBAC](/.knowledge/6-security/rbac-matrix.md), permiso `report_followups_consolidated`.

**Funcional:**
- El botón "Descargar reporte consolidado" está disponible en la pantalla Lista de Casos para los roles habilitados.
- Al pulsarlo, un modal solicita el rango de fechas (desde / hasta, máx. 1 año).
- Se genera un Excel enriquecido con los casos creados en el rango, estructurado en las siguientes hojas con sus respectivas columnas:
  - **Casos**: `"Código Caso"`, `"Fecha Creación"`, `"Documento Víctima"`, `"Nombre Víctima"`, `"Apellido Víctima"`, `"Edad"`, `"Género"` (normalizado a textos amigables like Mujer/Hombre/Otro), `"Teléfono Víctima"`, `"Municipio"` (nombre), `"Departamento"` (nombre), `"Equipo Asignado"`, `"Estado Caso"`.
  - **Registro**: `"Código Caso"`, `"Documento Víctima"`, `"Tipo Documento"`, `"Nombres Víctima"`, `"Apellidos Víctima"`, `"Edad"`, `"Género"`, `"Teléfono Víctima"`, `"Dirección Víctima"`, `"Correo Electrónico"`, `"Nombre Identitario"`, `"Orientación Sexual"`, `"Procedencia"`, `"Ocupación"`, `"Grupo Étnico"`, `"Si es Afrodescendiente"`, `"Si es Indígena"`, `"Si es Persona Campesina"`, `"Número de Hijos"`, `"Estado Civil"`, `"Discapacidad u Otras Condiciones"`, `"Municipio Residencia"`, `"Departamento Residencia"`, `"Nombres del Contacto"`, `"Teléfono del Contacto"`, `"Parentesco del Contacto"`, `"Descripción de los Hechos"`, `"Fecha Ocurrencia Hechos"`, `"Hora Inicio Hechos"`, `"Hora Fin Hechos"`, `"Día de la Semana Hechos"`, `"Ocurrencia de los Hechos"`, `"Ámbito de la Violencia"`, `"Escenario de Violencia"`, `"Existe Riesgo Feminicida"`, `"Nivel de Riesgo (Form 2)"`, `"Nombre Agresor"`, `"Tipo Documento Agresor"`, `"Documento Agresor"`, `"Dirección Agresor"`, `"Teléfono Agresor"`, `"Relación con Agresor"`, `"Violencia Física Aumentada (Último Año)"`, `"Separado de Pareja (Último Año)"`, `"Amenazado con Arma/Objeto"`, `"Amenaza de Muerte a Ella/Hijos"`, `"Celoso Violento Constantemente"`, `"Cree Capaz de Matarla"`, `"Formulario Registro Completado (Form 1)"`, `"Formulario Valoración Completado (Form 2)"`.
  - **Seguimientos**: `"Código Caso"`, `"Documento Víctima"`, `"ID Seguimiento"`, `"Número Secuencia"`, `"Fecha Planificada"`, `"Fecha Realizado"`, `"Estado"`, `"Equipo"`, `"Nivel Riesgo"`, `"Resumen / Notas"`.
  - **Respuestas de formulario**: `"Código Caso"`, `"Documento Víctima"`, `"ID Seguimiento"`, `"Fecha Realizado Seguimiento"`, `"ID Envío Formulario"`, `"Pregunta"`, `"Respuesta"`.
  - **Timeline**: `"Código Caso"`, `"Documento Víctima"`, `"Fecha Evento"`, `"Categoría"`, `"Tipo Evento"`, `"Descripción"`, `"Actor"`.
- Si no se encuentran casos en el rango de fechas seleccionado, el sistema no descarga un Excel vacío; en su lugar, el backend responde un error de negocio que el modal muestra como advertencia: `"No se encontraron casos creados en el rango de fechas seleccionado."`.
- Cada seguimiento queda claramente diferenciado (una fila por seguimiento, referenciado a su caso).

**Técnico:**
- El archivo refleja la información **más reciente** registrada en el sistema al momento de la descarga (sin caché).
- El endpoint construye el Excel únicamente con el rango de fechas seleccionado en el modal.
- Se valida el rango (obligatorio, `desde <= hasta`, ventana ≤ 366 días, y existencia de datos en el rango).

## Historia de Usuario

| ID | Historia |
| :--- | :--- |
| US-RPT-01 | Como **Supervisor/Administrador**, quiero descargar un Excel consolidado de los seguimientos de los casos creados en un rango de fechas, para consultar la evolución de la atención sin depender de un script SQL manual. |

## Fuera de alcance (por ahora)

- Filtros adicionales (rol, territorio, estado del caso) más allá del rango de fechas.
- Exportación por caso individual (botón por fila).
- Generación asíncrona / envío por correo. Ver [nota de rendimiento](/.knowledge/3-features/ReporteConsolidadoSeguimientos/backend/services/reporte-consolidado-service.md).
