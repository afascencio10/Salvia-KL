# Hacer Seguimiento

Ruta: `GET /salvia/hacer-seguimiento/:id`

Carga un seguimiento por ID. Muestra la tarjeta de la víctima en la parte superior y dos tabs: **Información del caso** (componente `case-info`) y **Formulario de seguimiento** (componente `dinamic-form`).

El formulario se completa por secciones. Cada sección se guarda individualmente vía `POST /api/v1/forms/saveSection`. Al responder todas las secciones visibles, se emite `form-completed` y redirige al detalle del caso.

`canEdit` se vuelve `false` si el seguimiento tiene status `REALIZADO` con más de 5 días desde `completed_at`.
