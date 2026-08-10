# Estructura del formulario "Registro de Caso" (dinamic-form)

> Nota: este documento describe el formulario **"Registro de Caso"** (form_id
> `0a24ab30-3cfc-4861-b74d-65d21524bc00`), creado en esta sesión de trabajo
> para migrar la pantalla `set_victim_case.html` al componente `dinamic-form`
> (ver [registro-caso-v2-index.md](registro-caso-v2-index.md)). No es el
> formulario de "Hacer Seguimiento" (`2d0aeb46-...`), que ya existía antes y
> solo se usó como referencia de patrón — avisa si en realidad buscabas ese.

**Motor:** `dinamic-form` (form / form_section / question / option /
visibility_condition — mismo motor DB-driven que Hacer Seguimiento).
**Guardado:** progresivo por sección (`POST /api/v1/forms/saveSection`),
igual que Seguimiento — ver [flow-E03](Flujos/flow-E03-cuando-se-guarda-primera-seccion.md).

**Totales:** 9 secciones · 128 preguntas · 36 opciones estáticas (escala de
1996) · 57 condiciones de visibilidad · 1 render_modification.

---

## Resumen de secciones

| # | Sección | Preguntas | Descripción |
|---|---------|:---------:|-------------|
| 1 | Datos de la Víctima | 16 | Identificación, residencia y accesibilidad de la víctima |
| 2 | Contacto de Apoyo | 4 | Persona de apoyo o confianza de la víctima |
| 3 | Hechos | 18 | Relato, ubicación y clasificación de los hechos de violencia |
| 4 | Agresor | 11 | Datos e identificación del presunto agresor |
| 5 | Tamizaje | 39 | Batería de preguntas para calcular el nivel de riesgo |
| 6 | Datos Personales | 34 | Discapacidad, nacionalidad, identidad y condición socioeconómica |
| 7 | Plan de Acción | 2 | Acciones acordadas con la víctima |
| 8 | Denuncia Fácil | 1 | Habilitación del canal de denuncia fácil |
| 9 | Lugar de Atención | 3 | Municipio donde se brinda la atención |

---

## 1 — Datos de la Víctima

| # | Pregunta | Tipo | Requerida |
|---|---|---|:---:|
| 1 | Nombres | text | Sí |
| 2 | Apellidos | text | Sí |
| 3 | Nombre identitario | text | No |
| 4 | Teléfono | number | Sí |
| 5 | Tipo de documento | dropdown | Sí |
| 6 | Número de documento | text | Sí |
| 7 | Zona de residencia | dropdown | Sí |
| 8 | Departamento de residencia | dropdown | Sí |
| 9 | Ciudad de residencia | dropdown | Sí |
| 10 | Municipio de residencia | dropdown | Sí |
| 11 | Dirección de residencia | text | Sí |
| 12 | Accesibilidad *(banner informativo)* | info | No |
| 13 | ¿Tiene alguna discapacidad? | boolean | Sí |
| 14 | Ajustes razonables requeridos (VBG) | multiple | Sí |
| 15 | ¿Requiere intérprete de idioma? | boolean | Sí |
| 16 | ¿Cuál idioma? | text | No |

Los 5 campos obligatorios que alimentan la creación inmediata del
`victim_case` (nombres, apellidos, tipo/número de documento, municipio de
residencia — preguntas 1, 2, 5, 6, 10) son los que reciben valores dummy si
aún no se han respondido — ver [flow-E03](Flujos/flow-E03-cuando-se-guarda-primera-seccion.md).
Departamento/Ciudad (8, 9) son solo cascadas de UI para resolver el
Municipio; no se proyectan sobre `victim_case_form2`.

## 2 — Contacto de Apoyo

| # | Pregunta | Tipo | Requerida |
|---|---|---|:---:|
| 1 | Nombres del contacto de apoyo | text | No |
| 2 | Teléfono del contacto de apoyo | number | No |
| 3 | Correo del contacto de apoyo | text | No |
| 4 | Parentesco del contacto de apoyo | dropdown | No |

Única sección completamente opcional.

## 3 — Hechos

| # | Pregunta | Tipo | Requerida |
|---|---|---|:---:|
| 1 | Relato de los hechos | text | Sí |
| 2 | Fecha de los hechos | date | Sí |
| 3 | Hora de inicio de los hechos | datetime | Sí |
| 4 | Zona de los hechos | dropdown | Sí |
| 5 | Departamento de los hechos | dropdown | Sí |
| 6 | Ciudad de los hechos | dropdown | Sí |
| 7 | Municipio de los hechos | dropdown | Sí |
| 8 | Dirección de los hechos | text | Sí |
| 9 | Escenario de la violencia | dropdown | Sí |
| 10 | Tipo de violencia experimentada | multiple | Sí |
| 11 | Subtipo de violencia experimentada | multiple | Sí |
| 12 | Ámbito de la violencia | multiple | Sí |
| 13 | Sector laboral de ocurrencia | dropdown | Sí |
| 14 | ¿La violencia fue motivada por género? | boolean | Sí |
| 15 | ¿Denunció previamente? | boolean | Sí |
| 16 | ¿A quién denunció? | multiple | Sí |
| 17 | ¿La atención recibida fue apropiada? | boolean | Sí |
| 18 | Recurrencia de la agresión | dropdown | Sí |

## 4 — Agresor

| # | Pregunta | Tipo | Requerida |
|---|---|---|:---:|
| 1 | Número de agresores | dropdown | Sí |
| 2 | Proximidad con el agresor principal | dropdown | Sí |
| 3 | Relación con el presunto agresor | dropdown | Sí |
| 4 | Ocupación del agresor | dropdown | Sí |
| 5 | ¿Depende económicamente del agresor? | boolean | Sí |
| 6 | Identidad de género del agresor | dropdown | Sí |
| 7 | Nombres del agresor | text | No |
| 8 | Tipo de documento del agresor | dropdown | No |
| 9 | Número de documento del agresor | text | No |
| 10 | Dirección del agresor | text | No |
| 11 | Teléfono del agresor | number | No |

La pregunta 3 (**Relación con el presunto agresor**) es la que determina
`wasPartner` (`código ∈ {'pi','ex'}` → pareja/expareja íntima) y con eso qué
rama de Tamizaje (Sección 5) se muestra.

## 5 — Tamizaje

39 preguntas, todas booleanas salvo el banner final. Se divide en 3 bloques
por `visibility_condition` (según el resultado de la pregunta 4.3):

| Bloque | Preguntas (orden) | Cantidad | Se muestra si... |
|---|---|:---:|---|
| Comunes | 1–6 | 6 | Siempre |
| Pareja íntima | 7–24 | 18 | `wasPartner == true` |
| No pareja | 25–38 | 14 | `wasPartner == false` |
| Banner de riesgo | 39 | 1 (info) | Siempre, al final |

Cada respuesta "Sí" suma 1 punto. El score se compara contra umbrales
distintos según la rama (`riskThresholdsPartner` / `riskThresholdsNonPartner`
en `victim_case_form_service.go`) para obtener `riskLevel` (1–4). Varias
preguntas de las ramas pareja/no-pareja comparten texto casi idéntico (ej.
"¿El agresor está desempleado?" aparece en ambas) — por eso el mapeo interno
usa `FieldKey(sección, orden)` y no el texto de la pregunta.

## 6 — Datos Personales

| # | Pregunta | Tipo | Requerida |
|---|---|---|:---:|
| 1 | Fecha de nacimiento | date | Sí |
| 2 | ¿Tiene dificultades físicas, mentales o sensoriales? | boolean | Sí |
| 3 | Dificultad para oír | single (escala) | Sí |
| 4 | Dificultad para hablar | single (escala) | Sí |
| 5 | Dificultad para ver | single (escala) | Sí |
| 6 | Dificultad para moverse o caminar | single (escala) | Sí |
| 7 | Dificultad para tomar objetos | single (escala) | Sí |
| 8 | Dificultad para entender o aprender | single (escala) | Sí |
| 9 | Dificultad para comer | single (escala) | Sí |
| 10 | Dificultad para relacionarse con otros | single (escala) | Sí |
| 11 | Dificultad para actividades cotidianas | single (escala) | Sí |
| 12 | Ley 1996 — apoyos para la toma de decisiones | multiple | Sí |
| 13 | Nacionalidad | dropdown | Sí |
| 14 | Nacionalidad específica | dropdown | No |
| 15 | Condición migratoria | dropdown | No |
| 16 | Identidad de género | dropdown | Sí |
| 17 | Orientación sexual | dropdown | Sí |
| 18 | Sexo asignado al nacer | dropdown | Sí |
| 19 | Población especialmente protegida | multiple | Sí |
| 20 | Afiliación étnica | dropdown | Sí |
| 21 | Pueblo indígena | dropdown | No |
| 22 | ¿Reconocimiento como campesina? | boolean | Sí |
| 23 | Estado civil | dropdown | Sí |
| 24 | Último nivel educativo | dropdown | Sí |
| 25 | Ocupación | dropdown | Sí |
| 26 | Forma de generación de ingresos | dropdown | Sí |
| 27 | Modalidad de actividad sexual pagada (ASP) | multiple | No |
| 28 | Relación laboral | dropdown | No |
| 29 | Inicio aproximado de la actividad (ASP) | text (timestamp) | No |
| 30 | Razón de la actividad (ASP) | multiple | No |
| 31 | Forma de tenencia de la vivienda | dropdown | Sí |
| 32 | Estrato de la vivienda | dropdown | Sí |
| 33 | Personas a cargo | multiple | Sí |
| 34 | ¿Está actualmente en embarazo? | boolean | Sí |

Preguntas 3–11 usan la escala de dificultad de 4 niveles (sin dificultad /
alguna dificultad / mucha dificultad / no puede — las 36 opciones estáticas
del seed). Preguntas 27–30 (ASP) solo son visibles condicionalmente.

## 7 — Plan de Acción

| # | Pregunta | Tipo | Requerida |
|---|---|---|:---:|
| 1 | Plan de acción | multiple | Sí |
| 2 | Explicación de la gestión Salvia | text | Sí |

## 8 — Denuncia Fácil

| # | Pregunta | Tipo | Requerida |
|---|---|---|:---:|
| 1 | ¿Permite el uso del canal de Denuncia Fácil? | boolean | Sí |

## 9 — Lugar de Atención

| # | Pregunta | Tipo | Requerida |
|---|---|---|:---:|
| 1 | Departamento de atención | dropdown | Sí |
| 2 | Ciudad de atención | dropdown | Sí |
| 3 | Municipio de atención | dropdown | Sí |

Solo el Municipio (3) se proyecta sobre `victim_case` (actualiza
`victim_case_victim_town_code` al activar el caso — ver
[flow-E04](Flujos/flow-E04-cuando-se-completa-formulario.md)); Departamento y
Ciudad (1, 2) son cascadas de UI. Esta sección **no incluye Asignación de
Ruta** — se excluyó a propósito del alcance de esta migración (decisión
tomada al planear el formulario).

---

## Notas generales

- **Tipos de pregunta usados:** `text`, `number`, `date`, `datetime`,
  `dropdown`, `boolean`, `multiple`, `single` (escala) e `info` (banners sin
  respuesta) — todos tipos ya existentes del motor `dinamic-form`, no se
  extendió el componente.
- **Opciones dinámicas:** la mayoría de los `dropdown`/`multiple` usan
  `state_options_path` para leer del catálogo compartido
  `victim_case_form2_enums` (vía `formState.enums.<categoría>`, calculado por
  `set_victim_case_v2.html`) en vez de duplicar `option` rows por pregunta.
- **Proyección a `victim_case_form2`:** 118 de las 128 preguntas mapean 1:1 a
  una columna (ver `victim_case_form_fields.go`); las ~10 restantes son
  cascadas de UI (Departamento/Ciudad) o banners informativos sin respuesta.
