# Formulario de seguimiento V.4 — lectura del Excel

Fuente: `3. Formulario seguimiento - V.4. .xlsx`, compartido en Drive el 14 de octubre de 2026.
Título interno de la hoja de orden: **Formulario seguimiento V4. (integrando Salvia Dignidad)**.

Cada hoja visible es un módulo del formulario nuevo. La hoja `Final` está oculta y conserva una compilación anterior (rutas diferenciales, empalme, cierre). No reemplaza a las hojas visibles. Sus comentarios útiles están al final de este documento.

Las columnas del Excel no siempre coinciden con el encabezado: en varias filas “Si”, el tipo y las opciones están corridos. Aquí se lee por el contenido de la fila, no por la letra de la columna.

La lectura del Excel queda debajo. Donde choque con las revisiones del prototipo, manda la segunda revisión. Donde la segunda no hable, manda la primera.

---

## Segunda revisión del prototipo

Registrada el 20 de octubre de 2026. Sustituye lo de la primera revisión en estos puntos.

**Menú**

Se quita Barreras transversales. El menú queda así:

1. Valoración del riesgo
2. Seguimiento a medidas de emergencia, solo si el caso tiene medidas
3. Seguimiento a entidades, solo si el caso ya tiene entidades
4. Identificación de entidades
5. Seguimiento al plan de acción
6. Conclusión del ciclo de atención

**Activación de ruta, dentro de cada entidad**

No vuelve la sección Activación de rutas y no bloquea el cierre. Sí se recuperan las dos preguntas que no estaban en Identificación ni en Seguimiento:

- ¿La persona está de acuerdo con que SALVIA active la Ruta de Atención y las gestiones institucionales correspondientes? (Sí / No, obligatoria)
- ¿Qué gestiones se realizaron desde Salvia para la activación de ruta? (texto, obligatoria)

Van en cada entidad. En Identificación de entidades, después de los canales y antes de la barrera nueva. En Seguimiento a entidades, en el bloque de registro de barrera nueva y activación de ruta.

**Orden en Seguimiento a entidades**

1. Acciones institucionales que ya se completaron
2. Seguimiento a las barreras que la entidad ya trae (solo si tiene)
3. Acciones institucionales pendientes
4. Describa cualquier actualización de la ruta desde el último seguimiento
5. Registro de barrera nueva y activación de ruta: la pregunta de barrera nueva (si es Sí, la tarjeta), el acuerdo y qué gestiones se realizaron
6. Canales, con la opción Ninguna
7. Describa la gestión de las rutas y barreras (texto obligatorio)

**Canales**

La lista de canales incluye Ninguna. Si se marca, se quitan los demás. Aplica en Identificación de entidades, Seguimiento a entidades y en la conclusión.

**Barreras transversales dentro de la entidad**

En la tarjeta de barrera nueva, «Sector en el que se presentó la barrera» incluye Transversales. Si se elige, se muestran talento humano, económicas, territoriales y enfoque diferencial, más fecha, funcionario, territorio, descripción y gestión. No se muestran las listas de salud, justicia o protección, ni la institución. Salud, justicia y protección siguen como en la primera revisión.

**Puntaje psicosocial**

«Sin acceso a la atención psicosocial o psicológica» vale 1 punto. El aviso de puntaje sigue sumando lo marcado.

---

## Primera revisión del prototipo

Registrada el 16 de octubre de 2026, a partir de la revisión del prototipo en `/salvia/prototipo-seguimiento-v4`. El prototipo no guarda respuestas.

Estas decisiones sustituyen, en el formulario que se va a construir, las partes del Excel que se indican.

**Menú**

1. Valoración del riesgo
2. Seguimiento a medidas de emergencia, solo si el caso tiene medidas
3. Seguimiento a entidades, solo si el caso ya tiene entidades
4. Identificación de entidades
5. Barreras transversales
6. Seguimiento al plan de acción
7. Conclusión del ciclo de atención

Se puede abrir cualquier sección del menú sin haber completado las anteriores.

**Barrera asociada a una entidad**

No hay sección propia de identificación de barreras de entidad. Tampoco se usa un modal: al marcar Sí, el registro se abre en una tarjeta dentro de esa entidad. Pueden agregarse varias.

La tarjeta pide sector (salud, justicia o protección), tipos de ese sector, institución, otra institución, fecha, funcionario, territorio, descripción y gestión. No incluye las cuatro listas transversales.

La gestión «Activación de ruta» se registra ahí y no abre otro módulo ni impide cerrar. «Alerta por barreras» pide la perspectiva profesional y, en la versión definitiva, notificaría al equipo técnico.

En Seguimiento a entidades, el orden de cada entidad es:

1. Acciones institucionales que ya se completaron
2. Seguimiento a las barreras que la entidad ya trae (solo si tiene)
3. Acciones institucionales pendientes
4. ¿Se identifica alguna nueva barrera con esta entidad? Si es Sí, la tarjeta anterior
5. Canales de la ruta
6. Actualización de la ruta desde el último seguimiento

**Identificación de entidades**

Se quita de Valoración del riesgo la pregunta sobre barreras institucionales.

Al final de cada entidad de esta sección va la misma pregunta. Si es Sí, se abre la misma tarjeta, ligada a esa entidad.

**Barreras transversales**

La antigua sección «Identificación de barreras» pasa a llamarse «Barreras transversales». Siempre está en el menú y no se asocia a una entidad. Equivale a `barrera_transversal` del formulario actual: no despliega salud, justicia ni protección.

La pregunta «¿Se identifican barreras transversales en este seguimiento?» habilita el resto: talento humano, económicas, territoriales, enfoque diferencial, fecha, funcionario, territorio, descripción y gestión.

**Activación de rutas**

Se elimina la sección y la regla de no cerrar el formulario sin diligenciarla. Acuerdo, sector, lugar, entidad y acciones pendientes de ese módulo ya se cubren en Identificación de entidades.

En Conclusión, si «¿Se realizó alguna acción institucional en la conclusión del caso?» es Sí, además de la entidad se preguntan los canales de la ruta. No bloquea el cierre.

**Descripción de la conclusión**

El enunciado visible es el del Excel, completo:

Registrar de forma clara la descripción del criterio y la orientación brindada sobre continuidad de derechos y rutas disponibles:

- Intentos de contacto (fechas, medios, número de intentos y tiempo transcurrido sin respuesta)
- Solicitud de la persona (qué solicita, razones expresadas)
- Avances y acciones completadas con evidencia y los criterios que permiten concluir que el plan de atención se cumplió
- No corresponsabilidad (explique la ausencia de gestiones o de compromiso con el acompañamiento)
- No clasificación como VBG/VpP (explicar por qué no se clasifica)

---

## Orden de los módulos

La hoja `Orden` define qué se muestra y cuándo. No es una sección que el agente diligencie.

| Bloque | Componente | Qué es | Cuándo aparece |
|---|---|---|---|
| 1.0 Seguimiento al riesgo de feminicidio, homicidio abyecto o fútil | Nuevos hechos | | Obligatorio |
| | Identificación de factores protectores | | |
| | Identificación de factores de riesgo | | |
| | Identificación de factores de riesgo extremo | | |
| | Descripción | | |
| Seguimiento a medidas de emergencia | Indaga el estado de entrega de la medida | | Solo cuando hay medidas de emergencia, posterior a la caracterización |
| 2.0 Seguimiento a la ruta de atención integral | A. Seguimiento a entidades | Seguimiento a las entidades ya creadas, incluidas las acciones de las entidades que tienen barreras | Solo si hay entidades registradas antes |
| | B. Identificación de entidades | Entidades a las que acudió la persona, o a las que otras entidades activaron ruta. ¿A dónde ha ido? | Obligatorio |
| | C. Identificación de barreras | Registro de barreras nuevas | Solo si en valoración del riesgo, o en identificación o seguimiento de entidades, refiere tener barreras |
| | D. Activación de rutas | Acciones institucionales que hará Salvia en la ruta y en las barreras | |
| Seguimiento al plan de acción | B. Remisiones Salvia | Equipos Salvia con los que pedirá articulación según lo identificado | |
| | C. Elementos de acompañamiento psicosocial | Elementos de acompañamiento a víctimas de VBG usados en el seguimiento | |
| | D. Descripción del seguimiento | Integración de lo realizado. Campo de texto | |
| Conclusión del ciclo de atención | A. Criterio y motivo | | |
| 6.0 | B. Realiza alguna acción | Si marca sí, el formulario no se puede cerrar sin diligenciar activación de rutas | Obligatorio |

**Comentario de Laura Daniela González (23 sep 2026), celda D11:** el seguimiento a entidades integra el formulario de seguimiento a barreras y el de entidades.

---

## 1. Análisis del riesgo

Hoja: `1. Analisis del riesgo`. Sección: Valoración del riesgo. Siempre visible.

| Pregunta | Tipo | Obligatorio | Visibilidad | Opciones y notas |
|---|---|---|---|---|
| Indique quién brindó la información durante este seguimiento | Única | Sí | Siempre | Contacto directo con la víctima. Contacto con institución. Contacto indirecto (familiar o persona de confianza) |
| Texto de instrucción | Recordatorio del módulo | No | Siempre | Los factores son indicadores técnicos de riesgo de feminicidio, homicidio abyecto o fútil y otras violencias graves. No son preguntas directas a la víctima. Solo hechos desde el último contacto (registro o seguimiento) hasta la fecha |
| ¿Se registraron nuevos hechos de violencia desde el último seguimiento? | Única | Sí | Siempre | Sí / No |
| Registre los nuevos hechos: tiempo, modo, lugar y relación de la víctima con la persona agresora | Abierta | Sí | Solo si la anterior es Sí | Campo de texto |
| Fecha en la que se presentaron los hechos | Fecha | Sí | Solo si la anterior es Sí | Campo de fecha |
| Indique si estos factores protectores se presentan en el caso | Múltiple | Sí | Siempre | Ver lista. La propia celda de tipo trae la regla: con 3 factores protectores y ningún riesgo extremo, considerar remisión a seguimiento general; con 4 factores de riesgo, considerar remisión a riesgo |
| Indique si estos factores de riesgo se presentan en el caso | Múltiple | Sí | Siempre | Ver lista |
| Si identifica algún factor de riesgo extremo, remita al equipo de riesgo. Verifique si ya hay articulación interinstitucional y priorice la protección inmediata | Múltiple | Sí | Siempre | Ver lista |
| ¿Desde la atención anterior se han identificado barreras institucionales que hayan contribuido a mantener o incrementar el riesgo? | Única | No | Siempre | Sí / No. Recordatorio si es Sí: “Registre estas en el módulo de barreras” |
| Describa la situación de riesgo actual y cómo los factores de riesgo y protectores influyen en que este se mantenga o modifique | Abierta | Sí | Siempre | Campo de texto |

**Protectores** (columna de recordatorio: Protector)

- Medidas de protección vigentes, monitoreadas, efectivas y adecuadas a las condiciones diferenciales de la persona
- Permanencia en un lugar seguro o traslado reciente a un entorno protegido y respetuoso de las condiciones diferenciales de la persona
- Separación física y sostenida de la persona agresora
- Redes de apoyo activas, seguras y respetuosas de las condiciones diferenciales de la persona (familiares, afectivas, comunitarias, institucionales u organizaciones sociales)
- Acompañamiento psicosocial o psicoterapéutico
- Cuenta con autonomía económica
- Acciones institucionales con efecto protector tangible en la percepción de seguridad de la persona
- Ninguna

Nota en la primera opción protectora: **“NECESARIO ADAPTAR REDACCIÓN PARA INCLUIR SITUACIONES PERSONAS OSIG”**.

**Riesgo**

- Condiciones habitacionales inseguras o inestables (incluye mudanzas forzadas o frecuentes asociadas a la violencia)
- Agresor pertenece a un grupo organizado al margen de la ley
- Agresor con consumo problemático de sustancias psicoactivas ilegales y legales
- Miedo a denunciar o separarse por represalias
- Bajo reconocimiento del riesgo que implica ausencia de implementación de esquema de autoprotección
- Embarazo con riesgo o gestación vulnerable
- Discapacidad física, cognitiva y/o psicosocial
- Residencia rural aislada o instituciones distantes
- Persona agresora con poder institucional o territorial
- Agresor pertenece a la fuerza pública
- Dependencia respecto de personas o entornos que rechazan, sancionan o buscan ocultar la OSIGD
- Aislamiento o debilitamiento de redes por rechazo relacionado con la OSIGD
- Permanencia en un territorio o entorno que controla, sanciona o reprime las OSIGD
- Ninguna

**Riesgo extremo**

- Amenazas de muerte explícitas
- Escalamiento en frecuencia o gravedad de la violencia
- Contacto sexual no consentido (último mes)
- Aislamiento, retención o privación de la libertad
- Acceso reciente del agresor a armas
- Destrucción o control de bienes
- Aumento de vigilancia, control o monitoreo
- Tentativa de feminicidio u homicidio abyecto o fútil (daño corporal grave, intento de asfixia, ataque con agente químico u otros líquidos hirvientes)

---

## 2. Seguimiento a entidades

Hoja: `2. Seguimiento a entidades`. Solo si el caso ya tiene entidades (hoja Orden).

**Comentario de Laura (23 sep 2026), celda F1.** Este apartado debe integrar tres escenarios, y la barrera queda asociada a la entidad para no repetir el formulario:

1. Nuevas acciones de la entidad
2. Identificación de barreras
3. Seguimiento a barreras

| Pregunta | Tipo | Obligatorio | Visibilidad | Opciones y notas |
|---|---|---|---|---|
| Listado de las entidades y barreras ya creadas | Lo trae el sistema | | Siempre en la sección | Comentario en B2: “Esto lo trae el sistema” |
| Selección de la entidad a actualizar | Múltiple | Sí | Siempre | Entidad creada previamente |
| Acciones institucionales que ya se completaron | Múltiple | | | Competencias del Excel de Rutas |
| Acciones institucionales pendientes | Múltiple + texto | | | Competencias del Excel de Rutas, más “Otro” y un campo de texto |
| ¿Se identifica alguna nueva barrera en este seguimiento con relación a esta entidad? | Única | Sí | Siempre | Sí. Comentario en D6 y nota de la hoja: si es Sí, despliega el formulario de Identificación de barreras. En la hoja no aparece la opción No |
| Seleccione los canales utilizados para la activación o gestión de la Ruta de Atención | Múltiple | Sí | | Notificación SALVIA. Articulación institucional telefónica o por mensajería. Denuncia fácil. Mi Salvia |
| Describa cualquier actualización de la ruta desde el último seguimiento | Abierta | Sí | | Campo de texto. Nota: el agente podría escribir “No” o “no aplica” y así se quita una pregunta |
| Texto informativo de las barreras de la entidad a la que se hace seguimiento | Info | | | |
| ¿Las barreras persisten? | Única | | Solo si hay barreras asociadas a la entidad | Sí / No. Comentario en B13: “Arrastre la barrera en este punto que se haya creado previamente” |
| Gestión de la barrera por parte de Salvia | | | | Activación de ruta: despliega y obliga a diligenciar la activación de ruta en Seguimiento al plan de acción. Alerta por barreras: notifica al rol de equipo técnico |
| ¿Cuáles barreras se cerraron? | Múltiple | | | Por barrera |
| Motivo del cierre | Única | | | La persona no desea que Salvia gestione la barrera. Barrera no gestionable desde la competencia institucional. Se clasificó incorrectamente la barrera. La barrera fue superada de manera efectiva |
| Describa la gestión de SALVIA frente a esta barrera. Si elige “Alerta por barreras”, incluya los elementos para el análisis. Si cierra, mencione qué permite concluir el cierre y que no quedan gestiones pendientes | Abierta | | | Campo de texto |

En la fila 24 de la hoja queda un “No” suelto, sin pregunta ancla.

---

## 3. Identificación de entidades

Hoja: `2. Identificación de entidades`. Obligatorio según la hoja Orden. Es el registro de entidades a las que la persona ya acudió.

| Pregunta | Tipo | Obligatorio | Visibilidad | Opciones y notas |
|---|---|---|---|---|
| ¿Ha puesto en conocimiento de alguna entidad/institución los hechos de violencia reportados? | Única | Sí | Siempre | Sí / No |
| Entidades a las que ya acudió | Repeater | | Solo si la anterior es Sí | El bloque siguiente es una entidad |
| Sector de la entidad | Única | | Dentro del repeater | Salud. Justicia. Protección. Educación. Ministerio Público. Entidad territorial. Organización de la sociedad civil, comunitaria o internacional de Naciones Unidas. Entidades clave con acciones complementarias y/o especializadas frente a VBG/VPP. Otra |
| Departamento | Única | Sí | | |
| Municipio | Única | Sí | | |
| Ciudad | Única | Sí | | |
| Entidad | Única | Sí | | Las del mapeo, más Otra |
| Otra ¿Cuál? | Texto | | Solo si en Entidad marcó Otra | |
| Número de radicado (NUNC, medida de protección). Escriba tipo y número. Ejemplo: NUNC 2026120956789 | Texto | Sí | | |
| Fecha en la que reportó a la entidad | Fecha | Sí | | |
| Acciones institucionales completadas | Múltiple | Sí | | Lista del Excel de Rutas, según la entidad |
| Acciones institucionales pendientes | Múltiple | Sí | | Lista del Excel de Rutas, según la entidad |
| Otra información asociada a la Ruta | Texto | Sí | | |
| Seleccione los canales utilizados para la activación o gestión de la Ruta de Atención | Múltiple | Sí | | Notificación SALVIA. Articulación institucional telefónica o por mensajería. Denuncia fácil. Mi Salvia |

---

## 4. Identificación de barreras

Hoja: `3. Identificación de barreras`. Sección: Registro de barreras.

Visible solo si en valoración del riesgo, o en identificación o seguimiento de entidades, la persona refiere barreras (hoja Orden). Desde seguimiento a entidades, la pregunta “¿Se identifica alguna nueva barrera…?” = Sí también abre este formulario.

Instrucción fija: “Registre a continuación las barreras identificadas en el acceso a la ruta de atención integral por parte de la persona.”

La columna F trae, para Salud, Justicia, Protección y barreras estructurales, un texto largo de definición del sector (ámbito, instituciones y marco normativo). Es recordatorio del módulo, no una pregunta.

Notas sueltas en el margen derecho de la hoja, sin pregunta dueña:

- Diferenciar visualmente las preguntas
- Funcionario/a o dependencia donde se presentó la barrera: “Dependencia y/o funcionario con el que se presentó la barrera”. “Opciones de despliegue todas”. “Opción cuál - entidad”. “Entidades - lista”
- Alerta por barreras (barreras sistemáticas): “Resuma las actuaciones realizadas por Salvia” y “¿Cuál es su perspectiva profesional frente a la barrera?”
- Tipología: quitar la transversal

| Pregunta | Tipo | Obligatorio | Visibilidad | Opciones y notas |
|---|---|---|---|---|
| Salud | Múltiple | | Solo si marcó barrera en el sector | 13 opciones, abajo |
| ¿A qué institución acudió? (salud) | | | | Hospital Público. Hospital Privado. IPS. EPS. Puesto de Salud. Consultorio médico. No sabe / No responde |
| Justicia | Múltiple | | Solo si marcó barrera en el sector | 14 opciones, abajo |
| Otra barrera en sector Salud | Texto | | | En la hoja esta fila está debajo del bloque de justicia. Parece un desorden de filas: el texto corresponde a Salud |
| ¿A qué institución acudió? (justicia) | | | | CAI-Policía. CAIVAS. Casas de Justicia. CAVIV. Comisaría de Familia. Estación de Policía (carácter administrativo y de convivencia). Fiscalía General de la Nación. Inspección de Policía. Instituto Nacional de Medicina Legal. Jueces Civiles o Promiscuos Municipales. Jueces de Control de Garantías o Jueces Penales. Policía Judicial. Unidad de reacción inmediata. No sabe / No responde |
| Otra barrera en sector Justicia | Texto | No | | |
| Protección | Múltiple | | | 8 opciones, abajo. La hoja no repite la condición “solo si marcó el sector”, a diferencia de Salud y Justicia |
| Otra barrera en el marco del sector Protección | Texto | No | | |
| ¿A qué institución acudió? (protección) | | | | Comisarías de familia. Fiscalía (funciones de protección). ICBF. Unidad Nacional de Protección. Jueces Civiles o Promiscuos Municipales. Inspección de Policía. Defensoría de Familia. No sabe / No responde. Ejército Nacional. Armada Colombiana |
| ¿La barrera se presentó en otra institución interna (empleador o empleadora) o externa? | Texto | Sí | | |
| Barreras estructurales y/o transversales | Grupo | | | Los cuatro bloques siguientes |
| Barreras institucionales y de talento humano | Múltiple | No | | 5 opciones |
| Barreras económicas y socioeconómicas | Múltiple | No | | 2 opciones |
| Barreras territoriales y geográficas | Múltiple | No | | 1 opción |
| Barreras por ausencia de enfoque diferencial | Múltiple | No | | 12 opciones |
| Fecha en la que se presentó la barrera (aproximada) | Texto | Sí | | La hoja dice “campo de texto”, no fecha |
| Funcionario/a o dependencia donde se presentó la barrera | Texto | Sí | | |
| Departamento / municipio / vereda en la que se presentó la barrera | Texto | Sí | | Nota: desplegar las opciones por territorio |
| Descripción de la barrera (tiempo / modo / lugar) | Texto | Sí | | |
| Gestión de la barrera | Múltiple | Sí | | Activación de ruta interinstitucional: despliega el módulo de activación de rutas. Alerta por barreras: informa al equipo técnico |
| Descripción de la gestión realizada con relación a las barreras | Texto | Sí | | |

**Salud**

- Dificultades de aseguramiento o afiliación en salud
- Negaciones en los servicios de urgencias
- Fallas en la calidad del servicio de urgencias
- Demoras en la asignación de citas y continuidad de tratamientos
- Faltas de activación de protocolos para violencia sexual
- Faltas de activación de protocolos para ataques con agentes químicos
- Demoras y dificultades en la valoración medicolegal
- Negaciones en el acceso a la Interrupción Voluntaria del Embarazo (IVE)
- Exigencias indebidas de autorización para procedimientos en personas con discapacidad cognitiva y/o psicosocial
- Atenciones en salud no centradas en cosmovisiones y prácticas culturales propias
- Negaciones o imposiciones de procedimientos a personas con diversidad sexual o identidad de género diversa
- Negación, restricción o condicionamiento de servicios y procedimientos por falta de adecuación a las necesidades diferenciales de la persona

**Justicia**

- Negativa institucional para recibir la denuncia
- Denegación al derecho a la no confrontación
- Deficiencias en la notificación y recolección de datos del agresor
- Falta de entrega de información y documentación procesal
- Imposición de cargas probatorias injustificadas
- Falta de celeridad en la investigación judicial
- Tipificación errónea del delito
- Traslado inadecuado o incompleto del caso / expediente judicial
- Ausencia o dificultad de acceso a representación judicial
- Riesgo de vencimiento de términos
- Medidas preventivas de libertad dirigidas al presunto perpetrador o condenado que ponen en riesgo la vida de la mujer y su núcleo familiar
- Divulgación de información confidencial del proceso
- Omisión o minimización del posible móvil de prejuicio en la investigación

**Protección**

- Ausencia o demoras en la adopción de medidas de protección urgentes o definitivas en Comisaría de Familia
- Ausencia o dilación en adopción de medidas de protección físicas
- Demora o ausencia en la activación de medidas de atención integral
- Incumplimiento de medidas de protección sin respuesta institucional
- Medidas de protección ineficaces o mal implementadas
- Fallas en la valoración y actualización del riesgo
- Omisión de restricción de visitas ante riesgo de feminicidio
- Omisión de apoyo policial en residencia o lugar de trabajo

**Institucionales y de talento humano:** falta de aplicación del enfoque o perspectiva de género / revictimización o violencia institucional. Falta de aplicación o desactualización de protocolos de atención en VBG. Capacidad institucional limitada. Desarticulación institucional y ausencia de gestión integral de casos. Insuficiencia en la orientación en derechos a las víctimas.

**Económicas:** falta de autonomía económica. Pobreza y desigualdad estructural.

**Territoriales:** ausencia institucional, infraestructura y conectividad insuficiente o instalaciones distantes.

**Enfoque diferencial:** OSIGD. Salud mental. Origen étnico (indígena, afrodescendiente, raizal, palenquera). Ciclo de vida (niñez, adolescencia, adultez mayor). Discapacidad física o cognitiva. Personas migrantes y refugiadas. Personas privadas de la libertad. Situación de calle. Víctimas de trata. Actividades sexuales pagas. Lideresas y defensoras de derechos humanos. Víctimas del conflicto armado.

---

## 5. Seguimiento a medidas de emergencia

Hoja: `Seguimiento a medidas de emerge`. Solo cuando la persona tiene medidas de emergencia solicitadas y/o activas, después de la caracterización.

Instrucción del módulo: verificar la trazabilidad de la entrega, el uso para acceder a justicia y a la ruta, y la preparación para el egreso. Las preguntas se formulan en lenguaje sencillo, sin revictimizar y sin pedir que la persona repita los hechos.

**Periodicidad (regla de sistema, no es una pregunta del agente)**

- Primer seguimiento: a la semana de la caracterización
- Segundo: a los 15 días
- Tercero: al mes
- Cuarto: a los dos meses, según la evolución del caso

Nota técnica de la hoja: la plataforma debe poder cambiar esa periodicidad si cambia el riesgo, hay nuevos hechos, fallas en la entrega o efectividad de la medida, u otras necesidades de protección. Debe existir un “Nuevo seguimiento” extemporáneo: elegir profesional (puede ser distinto de quien caracterizó), justificación, fecha y trazabilidad.

La hoja detalla tres momentos. El cuarto seguimiento queda nombrado en la periodicidad y no tiene preguntas propias.

### Primer seguimiento — primera semana

Objetivo: verificar la seguridad inmediata, confirmar entrega y funcionamiento de las medidas e identificar necesidades urgentes no cubiertas.

| Pregunta | Tipo | Visibilidad del texto adicional |
|---|---|---|
| En una escala de 1 a 5 (1 muy bajo, 5 muy alto), ¿cómo califica hoy el riesgo para usted y su familia? | Única: 1 a 5 | Texto obligatorio siempre, para justificar la puntuación. La pregunta está repetida en las filas 9 y 14-15; es la misma escala |
| En una escala de 1 a 5 (1 completamente segura, 5 peligro permanente), ¿cómo se siente actualmente? | Única: 1 a 5 | Texto obligatorio siempre |
| Actualmente considera que sus necesidades básicas están | Única | Totalmente cubiertas / Parcialmente cubiertas / No cubiertas. Texto siempre, sea cual sea la opción |
| ¿Se siente segura en el lugar donde está viviendo actualmente? | Única Sí/No | Texto siempre |
| ¿Ha tenido algún contacto con el agresor desde que iniciaron las medidas? | Única Sí/No | Texto siempre |
| ¿Ha recibido amenazas, llamadas, mensajes o intentos de ubicarla? | Única Sí/No | Texto siempre |
| ¿Ha ocurrido alguna situación nueva que aumente su preocupación por su seguridad? | Única Sí/No | Texto siempre |
| ¿Sus hijas, hijos o personas a cargo se encuentran seguras? | Única Sí/No | Texto siempre |

**Estado de las medidas.** La hoja pide repetir estas preguntas por cada medida otorgada.

| Pregunta | Visibilidad |
|---|---|
| ¿La recibió oportunamente? | Sí/No. Si es No, texto obligatorio de por qué no se entregó a tiempo |
| ¿Ha podido utilizarla sin inconvenientes? | Solo si “¿La recibió oportunamente?” es Sí. Si es No, texto obligatorio |
| ¿Ha presentado alguna dificultad para hacer efectiva la medida? | Solo si la primera pregunta de estado es Sí |
| ¿Existe alguna necesidad urgente que aún no haya podido resolver? | Sí/No, después de revisar cada medida. Si es Sí, texto |

**Redes de apoyo**

| Pregunta | Visibilidad |
|---|---|
| ¿Cuenta con alguna persona con la que puede contar en este momento? | Si es Sí, texto obligatorio |
| Además de familiares o amistades, ¿hay alguna persona, organización, comunidad, iglesia, lideresa, colectivo o institución que le esté brindando apoyo? | Si es Sí, texto obligatorio |
| Si ocurriera una emergencia hoy, ¿sabe a quién acudir primero? | Si es Sí, texto obligatorio |
| ¿Cuenta con un lugar seguro al que pueda trasladarse si fuera necesario? | Si es Sí, texto obligatorio |

### Segundo seguimiento — 15 días

Objetivo: verificar la activación de las rutas institucionales y fortalecer las redes de apoyo.

Repite las dos escalas de riesgo y de sensación de seguridad, cada una con texto obligatorio de justificación.

| Pregunta | Visibilidad |
|---|---|
| ¿Las medidas continúan respondiendo a sus necesidades? | Sí/No. Texto obligatorio siempre, explicando si cubrió o no |
| ¿Ha requerido algún ajuste? | Sí/No. Si es Sí, texto obligatorio con los ajustes |
| ¿Ha presentado denuncia? | Sí/No. Si es No, texto obligatorio con los motivos. Si es Sí, lista de instituciones a las que acudió |
| ¿Ha acudido a la Comisaría de Familia? | Solo si presentó denuncia = Sí. Si es Sí, texto obligatorio del tipo de atención y de las medidas otorgadas |
| ¿Ha recibido atención por parte de la Fiscalía? | Igual: solo si denunció. Si es Sí, texto obligatorio |
| ¿La Policía ha intervenido en su caso? | Igual. Si es Sí, texto obligatorio |
| ¿Ha recibido acompañamiento de la Defensoría del Pueblo? | Igual. Si es Sí, texto obligatorio |
| ¿Ha recibido atención en salud física o mental? | Igual. Si es Sí, texto obligatorio |
| ¿Ha recibido apoyo de la Alcaldía, Gobernación u otra entidad? | Igual. Si es Sí, texto obligatorio |
| ¿Está recibiendo apoyos adicionales de otras instituciones? | Igual. Si es Sí, texto |
| ¿Le asignaron número de radicado o noticia criminal? | Solo si denunció. Si es Sí, texto |
| ¿Le otorgaron medidas de protección? | Aparece dos veces, con la misma condición. Si es Sí, texto obligatorio del número de medidas de emergencia |

**Redes y recursos del segundo seguimiento**

| Pregunta | Visibilidad |
|---|---|
| Redes personales: ¿quiénes la están acompañando actualmente? | Texto siempre |
| Redes comunitarias: ¿hace parte de algún grupo de mujeres, organización, iglesia u otro espacio? | Sí/No. Si es Sí, texto obligatorio |
| Redes institucionales: ¿conoce qué entidades puede contactar si el riesgo aumenta? | Sí/No. Si es Sí, texto |
| Recursos personales: ¿qué fortalezas le han ayudado hasta ahora? | Texto siempre |
| Actualmente, ¿cuenta con…? | Sí/No. Si es Sí, despliega las siete preguntas siguientes |
| ¿Un lugar seguro? | Si es Sí, texto |
| ¿Alguien que cuide a sus hijos si fuera necesario? | Si es Sí, texto |
| ¿Acceso a salud? | Si es Sí, texto |
| ¿Teléfono o internet? | Si es Sí, texto |
| ¿Transporte? | Si es Sí, texto del tipo de transporte |
| ¿Documentos personales? | Si es No, texto de qué documentos no tiene |
| ¿Apoyo emocional? | Si es No, texto de quién puede brindarlo |

### Tercer seguimiento — primer mes

Objetivo: evaluar el avance hacia la estabilización y el funcionamiento de las rutas.

Repite las dos escalas, con texto obligatorio. Repite “¿Las medidas continúan respondiendo?” (texto siempre) y “¿Ha requerido algún ajuste?” (texto si es Sí).

| Pregunta | Visibilidad |
|---|---|
| ¿Cómo han avanzado sus procesos? | Texto |
| ¿Ha tenido audiencias? | Sí/No. Si es Sí, texto del proceso jurídico |
| ¿Ha sentido revictimización? | Sí/No. Si es Sí, texto de las situaciones |
| ¿Cuenta actualmente con representación jurídica? | Sí/No. Si es Sí, texto |
| ¿Actualmente tiene ingresos? | Sí/No. Si es Sí, texto del tipo de ingreso |
| ¿Ha podido conseguir empleo? | Sí/No. Si es Sí, texto |
| ¿Está desarrollando alguna actividad productiva? | Sí/No. Si es Sí, texto |
| ¿Cuáles siguen siendo sus principales necesidades? | Texto |
| ¿Han cambiado sus redes de apoyo? | Aparece dos veces. Sí/No. Si es Sí, texto |
| ¿Qué recursos siente que hoy tiene para afrontar una nueva situación de riesgo? | Texto |
| ¿Qué recursos considera que aún necesita fortalecer? | Texto |

---

## 6. Seguimiento al plan de acción

Hoja: `Seguimiento al plan de acción`. Agrupa activación de ruta, remisiones, acompañamiento psicosocial y la descripción.

### Activación de ruta

| Pregunta | Tipo | Obligatorio | Opciones y notas |
|---|---|---|---|
| ¿La persona está de acuerdo con que SALVIA active la Ruta de Atención y las gestiones institucionales correspondientes? | Única | Sí | Sí / No |
| ¿Qué gestiones se realizaron desde Salvia para la activación de ruta? | | Sí | Sin opciones en la hoja |
| Sector de la entidad | | Sí | Comentario de Laura (23 sep 2026), celda A6: si aquí se diligencian sector, entidad y lugar, no se debe volver a pedir lo mismo en Notificaciones; hay que arrastrarlo |
| Departamento | | Sí | |
| Municipio | | Sí | |
| Ciudad | | Sí | |
| Entidad | | Sí | Mapeo, más Otro |
| Acciones institucionales pendientes | | Sí | Link al catálogo de competencias de Rutas |
| Seleccione los canales utilizados para la activación o gestión de la Ruta de Atención | Múltiple | Sí | Notificación SALVIA. Articulación institucional telefónica o por mensajería. Denuncia fácil. Mi Salvia |

La hoja Orden dice que, si en la conclusión se marca que sí se realizó una acción institucional, el formulario no se cierra sin diligenciar esta activación.

### Remisiones Salvia

| Pregunta | Tipo | Obligatorio | Visibilidad | Opciones y notas |
|---|---|---|---|---|
| Durante la atención se han generado necesidades que requieran remisión a los equipos SALVIA | Única | Sí | | Sí / No |
| Equipos | Múltiple | | Si la anterior es Sí | Medidas de Emergencia. Atención psicosocial. Psicosocial Personas OSIGNH. Atención Hombres. Equipo Discapacidad. Estabilización. Regla de la hoja: si remite a medidas, no puede remitir a psicosocial |
| Medidas de emergencia | Múltiple | | Si eligió Medidas de Emergencia | Alojamiento. Transporte. Alimentación. Vestuario. Apoyo psicosocial. Otras ME |
| Criterios de remisión a psicosocial | Múltiple | | Si eligió Atención psicosocial | Ver lista y puntajes |
| Criterios OSIGNH | | | Si eligió Psicosocial Personas OSIGNH | Dos ítems, abajo |
| Atención Hombres | | | Si eligió ese equipo | Hombres diversos victimarios o posibles victimarios |
| Equipo Discapacidad | Múltiple | | Si eligió ese equipo | Apoyo comunicativo en LSC. Apoyo con enfoque de discapacidad (ley de capacidad legal) |
| Estabilización | Múltiple | | Si eligió ese equipo | Ver lista |

**Criterios psicosociales**

- Mujeres en todas sus diversidades, mayores de 18 años, víctimas de VBG o VpP, con voluntariedad y en post-emergencia (no riesgo inminente, no desborde emocional)
- Conducta suicida derivada de VBG o VpP (riesgo suicida, último mes, sin atención) — 3 puntos
- Factores relacionados con el temor frente a la activación de la ruta (miedo, represalias, desconfianza institucional) — 1 punto
- Interseccionalidad (gestación, étnico, recursos limitados, etc.) — 2 puntos
- Condiciones territoriales y de orden público (zona apartada, sin acceso a ruta) — 1 punto
- Sin acceso a la atención psicosocial o psicológica (barreras EPS, falta de tiempo o recursos). La hoja no le asigna puntos en esta fila
- Naturalización de la VBG y roles de género (minimización, dependencia, múltiples agresores) — 1 punto

**Psicosocial Personas OSIGNH**

- Ser OSIGNH, mayor de 18 años, víctima de VpP, en post-emergencia, con voluntariedad y disponibilidad. Obligatorio
- No aplica si ya recibe atención psicosocial por EPS, salvo que existan barreras por orientación sexual o identidad de género

**Estabilización**

- No tiene ninguna forma de ingresos o tiene ingresos de forma inestable
- La persona reside actualmente en vivienda temporal o familiar
- La persona tiene personas a cargo que limitan su tiempo para trabajar o estudiar
- Depende económicamente de la persona agresora

### Elementos de acompañamiento psicosocial

Pregunta: “Seleccione los elementos psicosociales aplicados en este seguimiento.” Múltiple, obligatoria. Lista desplegable.

- Escucha activa y contención emocional en crisis
- Psicoeducación del ciclo de la violencia, telaraña del abuso y/o impactos de la violencia
- Apoyo en identificación de necesidades específicas y riesgo
- Despatologización de experiencias y reacciones asociadas a la violencia
- Reconocimiento y comprensión de derechos de víctimas
- Orientación sobre acciones disponibles y responsabilidad del Estado
- Sensibilización y orientación sobre servicios de salud mental
- Co-construcción del Plan de Seguridad
- Identificación de señales tempranas de riesgo y estrategias de respuesta
- Fortalecimiento de la corresponsabilidad de la persona en su seguridad
- Orientación para acceso a servicios, programas o apoyos institucionales
- Identificación de redes seguras y afirmativas (familia, comunidad, instituciones)
- Identificación y fortalecimiento de recursos personales
- Fortalecimiento de la autonomía económica y vinculación a proyecto de vida
- Cuestionamiento de mandatos culturales que sostienen la violencia
- Seguimiento a la ruta, actualización de riesgo y revisión de barreras
- Acompañamiento en la toma de decisiones informadas respetando la autonomía
- Promoción de autocuidado y estrategias de protección digital y física
- Promoción de reflexión crítica sobre creencias que justifican control y dependencia
- Identificación de creencias que sostienen la violencia y proponer narrativas alternativas

### Descripción

| Pregunta | Tipo | Obligatorio | Nota |
|---|---|---|---|
| Describa las acciones realizadas en el seguimiento | Texto | Sí | |
| Describa los ajustes que realiza al plan de acción de Salvia | Texto | Sí | La hoja la describe como el acumulado: se va sumando al plan de acción desde el registro, de forma similar a los nuevos hechos |

---

## 7. Conclusión del ciclo de atención

Hoja: `Conclusión del ciclo de atenció`.

| Pregunta | Tipo | Obligatorio | Visibilidad | Opciones y notas |
|---|---|---|---|---|
| ¿Se concluye el ciclo de atención? | Única | Sí | Siempre | Sí despliega el resto. No |
| Recordatorio | Info | | Si concluyó | Registrar de forma clara y sustentada las razones. Si hubo acción institucional asociada al cierre, registrarla en el módulo de rutas |
| Motivo de la conclusión del ciclo de atención | Única | Sí, si concluyó | Si la primera es Sí | Pérdida de contacto. Solicitud expresa de la persona para finalizar el proceso. Cumplimiento del plan de acción Salvia. No corresponsabilidad por parte de la persona. La información disponible no permite clasificar el caso como VBG o VpP |
| Descripción del criterio y de la orientación brindada sobre continuidad de derechos y rutas disponibles | Texto | Sí | Si concluyó | La hoja lista qué debe cubrir el texto: intentos de contacto (fechas, medios, número de intentos y tiempo sin respuesta); qué solicita la persona y por qué; avances y acciones completadas con evidencia; no corresponsabilidad; por qué no se clasifica como VBG/VpP. En la misma fila aparece “Otro ¿cuál?” como opción, sin motivo “Otro” en la lista de arriba |
| ¿Se realizó alguna acción institucional en la conclusión del caso? | Única | | | Sí / No |
| Seleccione la entidad con la que realizará acción institucional | | | | Sin opciones. La hoja Orden dice que, si marca sí, no se puede cerrar el formulario sin diligenciar activación de rutas |

---

## Hoja oculta `Final`

Está oculta. Es una compilación anterior del formulario completo: valoración del riesgo, rutas diferenciales, registro de barreras, seguimiento a barreras, seguimiento del caso, remisiones, empalme y cierre. Las hojas visibles de V.4 la sustituyen en esos temas. Tres comentarios de esa hoja siguen siendo indicaciones de diseño:

| Celda | Autor y fecha | Comentario | A qué apunta en esa hoja |
|---|---|---|---|
| D4 | Grupo fortalecimiento Salvia, 13 may 2026 | “¿Esta regla o indicación la podría implementar el sistema?” | La regla de 3 protectores / 4 riesgos para sugerir el equipo |
| B44 | Grupo fortalecimiento Salvia, 13 may 2026 | “No sé si exista forma de que sea solo desplegable, pues para quitar preguntas” | “¿A qué institución acudió?” dentro de rutas diferenciales |
| B247 | Grupo fortalecimiento Salvia, 14 may 2026 | “Esta pregunta se omitiría si el sistema lo hace de forma automática” | Está en el bloque de seguimiento a barreras de esa hoja oculta |

En la misma hoja oculta, la celda B40 pide que las sugerencias de remisión a riesgo o a seguimiento general salgan como alertas según lo seleccionado, o que el sistema sugiera el empalme de equipo. La celda B41 anota que rutas diferenciales estaba ahí porque aún no se integraba el módulo de agente integral, y que lo importante es que el sistema muestre en qué ruta está la víctima, a qué instituciones accedió y si la ruta la activó Salvia o la víctima.
