# Registro V3 — Párrafo síntesis

> Fuente: `Observaciones y ajustes al MockUp de registro 2.0.pdf`, observación
> #18 (páginas 2-4), la de mayor prioridad de las fuentes leídas en esta
> planeación por ser retroalimentación directa y reciente del cliente. No
> existe ningún equivalente en el Excel/Guion Funcional de Drive — es una
> funcionalidad enteramente nueva frente a esas fuentes.
>
> **Qué es:** al avanzar el registro, en la parte superior del módulo
> **GS — Gestión Salvia**, el agente visualiza un párrafo auto-generado que
> resume a la víctima y el caso a partir de las respuestas ya diligenciadas.
> Se arma en 4 secciones (A-D). Cada sección se omite por completo si su
> condición disparadora no tiene datos. Dentro de cada sección, cada
> fragmento tiene su propia regla de inclusión/omisión — estas reglas están
> transcritas literales del PDF del cliente, sin reinterpretar.
>
> **Dónde se recalcula:** es un valor derivado, no una pregunta, y no tiene
> su propio evento. En el código (2026-10-31) lo calcula `parrafoSintesis()`
> (`sintesisA` a `sintesisD` en `registro-v3-app.js`) en cada render de Vue
> a partir de `answers`: se ve al entrar al módulo GS y en la pantalla de
> Cierre, y **no se guarda** en la BD. Ver
> [flow-E11](Flujos/flow-E11-cuando-actualiza-una-respuesta.md) paso 7.

---

## A. Contacto e identificación

**Plantilla:**
> La/El/Le Señora/e/_ [NOMBRE IDENTITARIO] [Nombre Apellido], identificada/e/o
> con [Tipo de documento] Número [Número de documento de identificación].
> [reconoce discapacidad]. [requiere intérprete de idiomas].

**Reglas por fragmento:**

| Fragmento | Regla |
|---|---|
| Tratamiento (La/El/Le Señora/Señor/Señore) | Según pronombres (Ella→Señora, Él→Señor, Elle→Señore) |
| Nombre identitario | Si hay valor, se incluye antes del nombre jurídico; si no hay, se omite (no deja espacio vacío) |
| "reconoce discapacidad" | Solo si CI.Q7 (¿Se reconoce con discapacidad?) = Sí; si No, no se muestra |
| "requiere intérprete de idiomas" | Solo si CI.Q9 (¿Requiere intérprete?) = Sí; si No, no se muestra |

Se omite la sección completa si no hay nombres/apellidos diligenciados.

---

## B. Datos personales

**Plantilla (bloque 1):**
> [Nombre Apellido] tiene XX años [calculado de fecha de nacimiento], nacional
> [Nacionalidad]. Reside en [Departamento de residencia] en [Municipio de
> residencia] en [Zona donde se encuentra lugar de residencia].

| Fragmento | Regla |
|---|---|
| Edad | Calculada de `fecha_nacimiento`; se omite el bloque si no hay fecha |
| Nacionalidad | "Colombiana(o)" · si Extranjero: "de [país de nacionalidad]" · o "Apátrida" |
| Departamento/Municipio/Zona de residencia | Respuesta directa |

**Plantilla (bloque 2):**
> Se reconoce [Identidad de Género], fue asignado al nacer [Sexo asignado al
> nacer] con orientación sexual [Orientación sexual].

| Fragmento | Regla |
|---|---|
| Identidad de género | Respuesta directa, **excepto si = "Otro" → no se muestra el bloque** |
| Orientación sexual | Respuesta directa, **excepto si ∈ {"Otro", "No estoy segura/o/e", "No informa"} → no se muestra esa parte** |

**Plantilla (bloque 3):**
> Se identifica [Grupo étnico] [reconoce persona campesina]. Pertenece a
> población de especial protección constitucional [Poblaciones a las que
> pertenece].

| Fragmento | Regla |
|---|---|
| Grupo étnico | Respuesta directa, **excepto "No se auto reconoce/no informa" → no se despliega esta parte** |
| Reconoce persona campesina | Si Sí → agrega "y se reconoce como persona campesina"; si No → no se muestra |
| Población de especial protección | Enlista todas las respuestas marcadas, **excepto si = "Ninguno" → no se muestra** |

**Plantilla (bloque 4):**
> Persona [Estado Civil] [último nivel de escolaridad]. [Forma de generación
> de ingresos] en su ocupación [Ocupación]. Su vivienda es [Forma de tenencia
> de la vivienda] en estrato [Estrato de la vivienda]. Tiene personas a su
> cargo [Personas o animales a su cargo].

| Fragmento | Regla |
|---|---|
| Estado civil | Respuesta directa |
| Último nivel de escolaridad | Si = "Sin educación formal" → se muestra tal cual; para cualquier otra respuesta → "con [respuesta]"; si = "No informa" → no se muestra |
| Forma de generación de ingresos + Ocupación | Respuesta directa de ambas |
| Forma de tenencia de vivienda | Respuesta directa, **excepto "No informa" → no se muestra** |
| Estrato | Respuesta directa, **excepto "No informa/No sabe" → no se muestra** |
| Personas/animales a cargo | Respuesta directa |

---

## C. Hechos de violencia

**Plantilla:**
> Hecho de violencia de Tipo(s) [Tipo o Tipos de violencia], [Subtipo(s) de
> violencia], ocurrido en ámbito(s) [Ámbito o ámbitos de violencia] en
> escenario [Escenario de ocurrencia]. Sucedido el [Fecha del hecho + Hora
> del hecho] en [Departamento] en [Municipio] en [Zona del hecho]. [Motivado
> por razones de género]. [Cuántas veces ha ocurrido esta agresión].

| Fragmento | Regla |
|---|---|
| Tipo(s)/Subtipo(s)/Ámbito(s)/Escenario | Respuesta directa |
| Fecha/Hora/Departamento/Municipio/Zona del hecho | Respuesta directa |
| Motivado por razón de género | Si Sí → "Motivado por razones de género"; si No → no se muestra |
| Cuántas veces ha ocurrido la agresión | Respuesta directa |

Se omite la sección completa si no hay tipo de violencia ni relato diligenciados.

---

## D. Características del agresor

**Plantilla:**
> Responsabiliza a [Cuántas personas agresoras participaron] personas de la
> agresión. Principalmente a [Conoce y convive con la principal persona
> agresora], de nombre [Nombres y apellidos de la persona agresora],
> [Identidad de género de la persona agresora], con quien tiene una relación
> de [Relación con la persona agresora], cuya ocupación es [Ocupación de la
> persona agresora].

| Fragmento | Regla |
|---|---|
| Número de personas agresoras | Respuesta directa, **excepto "No se puede determinar/No informa" → no se muestra** |
| Conoce y convive | Respuesta directa |
| Nombres y apellidos del agresor | Respuesta directa (si hay) |
| Identidad de género del agresor | Respuesta directa |
| Relación con la persona agresora | Respuesta directa, **excepto {"Sin relación", "Otro"} → no se muestra** |
| Ocupación del agresor | Respuesta directa, **excepto {"Otro", "No conoce o no informa"} → no se muestra** |

Se omite la sección completa si no hay relación con agresor ni número de
agresores diligenciados.

---

## Notas de implementación

- Las 4 secciones son independientes — cada una se evalúa y se omite por su
  cuenta si no tiene datos suficientes, no existe una condición "todo o nada".
- El párrafo se recalcula automáticamente con cada respuesta nueva (vive
  dentro del evento `Cuando actualiza una respuesta del formulario`), no
  requiere que el agente lo regenere manualmente.
- No se persiste en BD — se computa en el cliente a partir de las respuestas
  ya guardadas, igual que el resultado del Tamizaje.
