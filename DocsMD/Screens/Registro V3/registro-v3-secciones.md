# Registro V3 — Catálogo de preguntas por módulo

> **Fuente de la verdad:** `Registro - Formulario Kreivo v5.0.xlsx` (copia Google
> Sheets) y `Registro - Guion Funcional v3.0.pdf`, ambos en Google Drive
> (carpeta `RootAF/.../MockUps Andres` del usuario — ver referencia en memoria
> de sesión). El prototipo React (`Salvia/src/data/registroSalviaSchema.js`) y
> el `Registro de Caso V2` real de este repo se usan solo como guía de forma,
> **nunca** para resolver una duda de contenido — ante cualquier contradicción
> gana el Excel/PDF de Drive.
>
> Agrupación en 7 módulos (CI/DP/CA/HV/AG/TZ/GS): heredada del prototipo
> React, que a su vez replica el agrupamiento real del mockup. Útil porque
> coincide con los dos bloques que el propio Excel exige en su pestaña
> "Ajustes generales Registro": bloque **Víctima** (CI+DP+CA) y bloque
> **Caso** (HV+AG+TZ+GS).
>
> **Decisión de alcance (confirmada con el usuario):** esta ronda cubre el
> núcleo del formulario (Habeas data → Cabezote → 7 módulos). El Expediente de
> víctima (pantalla de consulta) queda para una ronda posterior — ver
> `registro-v3-data.md` § Decisiones de arquitectura.

---

## Divergencias entre este catálogo y el código (2026-10-31)

Este MD es la **fuente de contenido** (preguntas, opciones, reglas del cliente); el catálogo que ejecuta la pantalla es `src/frontend/js/components/registro-v3-campos.js`. Comparadas las 82 filas de módulos contra el catálogo, **el código manda** en lo siguiente:

| Tema | Este MD dice | El código hace |
|---|---|---|
| Obligatoriedad (`req`) | "máx. 2º contacto" en DP_4, DP_16, DP_18, DP_23, DP_24, DP_25, DP_30, HV_3, HV_7, AG_4 y AG_6 | Las marca **Deseable** (`D`); no se exigen ni al guardar ni al finalizar |
| Obligatoriedad de DP_26 (gestación) | "1er contacto · condicional" | `C2`, y la condición devuelve siempre `true` (se muestra a todas las personas) |
| [máx. 2º contacto] al finalizar | "Todas las [máx. 2º contacto] de los 7 módulos" | Solo DP_10, DP_15 y CA_1-4; el resto de C2 y todas las condicionales no se exigen (`obligatoriasFinalizarV3`) |
| DP_3 (actividades con dificultad) | 9 actividades, cada una en escala 1-4 | Checklist `multiple` de la categoría `victim_case_form2_activities_unable_to_perform`; sin escala |
| Departamento / municipio (DP_27-28, HV_4-5) | Desplegables DIVIPOLA encadenados | DP_27/HV_4: lista fija de nombres de departamento; DP_28/HV_5: **texto libre** (DP_28 se guarda en `victim_case_victim_town_code`) |
| HV_5 | "dropdown · condicional (depende de Q4)" | `text` sin condición |
| Tamizaje — preguntas condicionales | Preguntas OSIGD y de ámbito laboral ocultas según el caso; VPP solo si OSIGD | Se muestran las 24 preguntas de la rama y TZ_8 siempre que se revele el resto; el código declara los índices condicionales pero no los usa |
| Tamizaje — rama pareja | Encabezado dice "23 preguntas" | Son 24 (la 24 es ideación suicida, igual que en la rama no pareja) |
| GS_5-GS_8, GS_13-GS_17 | Preguntas numeradas del catálogo | No están en `registro-v3-campos.js`: se dibujan directamente en `registro_v3.html` (por eso `guardarModulo('GS')` no las valida) |
| GS_18 | "Botón Activar Alerta Inminente" | Radio "¿Activar alerta inminente?" (GS_18) + textarea de justificación (GS_19, opcional) |
| GS § C/D | Repeaters con 11 y 15 campos por entrada | Entidad: tipo, ubicación en cascada, sector, sede, radicado, canal, objetivo y acciones; Barrera anidada en la entidad: barreras específicas, categorías transversales, competencias, fecha, funcionario, descripción y gestión |
| Opciones | Conteos de la planeación | Ver "Conteos reales de enums" en [`registro-v3-data.md`](registro-v3-data.md) |

Coinciden con el MD: CI completo (incluida la máscara de CI_6 y la regla de CI_5 = "Persona sin identificación"), DP_17 como 1er contacto, DP_19-22 con sus condicionales, CA_1-4 como [máx. 2º contacto], HV_11 con máximo 2 y los 6 equipos de remisión con sus criterios y puntajes.

---

## Pre-sección — Habeas data (gate, antes de los módulos)

| # | Pregunta | Tipo | Requerida | Notas |
|---|---|---|:---:|---|
| 1 | Autorización para el tratamiento de datos personales *(banner legal largo + pregunta)* | boolean | Sí | Si responde **No**, el formulario se cierra ahí — no se crea ningún registro. Ley 1581 de 2012 / Resolución 772 de 2024. |
| 2 | Tipo de contacto | single | Sí | Opciones: **Víctima - Caso Nuevo** (abre este formulario) · Víctima - Seguimiento (abre el formulario de Seguimiento, ya existente, fuera de alcance) · Tercero - Reporte (abre el formulario de Reporte, fuera de alcance de esta ronda). |

---

## Módulo CI — Contacto e identificación (bloque Víctima)

| # | Pregunta | Tipo | Requerida | Opciones / Notas |
|---|---|---|:---:|---|
| 1 | ¿Cómo prefieres que nos refiramos a ti? (pronombres) | dropdown | 1er contacto | Ella · Él · Elle — determina el tratamiento mostrado en cabezote (Sra./Sr./Sre.) |
| 2 | ¿Cuáles son sus nombres y apellidos? | text | 1er contacto | Validación "solo letras" — rechaza entradas compuestas solo por números (obs. #13 del cliente) |
| 3 | ¿Con qué nombre se identifica o prefiere que le llamemos? (nombre identitario) | text | 1er contacto | Validación "solo letras" (obs. #13) |
| 4 | ¿Cuál es su número de teléfono de contacto en Colombia? | text | 1er contacto | Alfanumérico, admite símbolos |
| 5 | ¿Qué tipo de documento de identificación tiene? | dropdown | 1er contacto | 19 opciones: Carné Diplomático · Cédula de Ciudadanía · Cédula de Extranjería · Certificado de Nacido Vivo (país origen) · Certificado de Nacido Vivo en Colombia · Certificado RAMV · Documento Extranjero/DNI · Otro no especificado · Partida de Nacimiento país origen · Pasaporte · PEP · PEP-FF · Permiso por Protección Temporal · Persona sin identificación · Registro Civil de Nacimiento · Salvoconducto · Tarjeta Identidad · Visa · Visa de Refugiado |
| 6 | ¿Cuál es su número del documento de identificación? | text | 1er contacto · condicional | Formato de validación **depende del tipo de documento** (ver tabla de validaciones abajo). Prediligencia Nacionalidad (colombiano/extranjero) y, si es colombiano, inhabilita "condición migratoria". "Persona sin identificación" no pide número. |
| 7 | ¿Se reconoce como una persona con discapacidad? | boolean | 1er contacto | |
| 8 | ¿Requiere de algún ajuste razonable en el marco de la atención a su caso de VBG y/o VpP? | multiple | 1er contacto · condicional (Q7=Sí) | 11 opciones: Intérprete lengua de señas · Guía-intérprete sordociegos · Intérprete de idiomas/traducción · Apoyos para comunicación (pictogramas, tableros) · Formatos accesibles (lectura fácil, macrotipo, audio, braille) · Adecuaciones accesibilidad física · Dispositivos/ayudas técnicas · Persona de confianza durante la atención · Transporte seguro/acompañamiento · No informa · No requiere / No requiere apoyos especiales |
| 9 | ¿Requiere intérprete de idiomas y/o traducción en otra lengua distinta al español? | boolean | 1er contacto | |

### Validaciones de número de documento por tipo (Q6)

| Tipo de documento | Formato |
|---|---|
| Cédula de Ciudadanía | Numérico, 6-10 dígitos |
| Cédula de Extranjería | Numérico, 6-10 dígitos |
| Tarjeta Identidad | Numérico, 10-11 dígitos |
| Pasaporte | Alfanumérico, 6-9 caracteres |
| Carné Diplomático | Numérico |
| Permiso por Protección Temporal | Numérico, 7-8 dígitos |
| Persona sin identificación | No pide número |
| Certificado Nacido Vivo (origen/Colombia), RAMV, Documento Extranjero/DNI, Otro, Partida Nacimiento, Registro Civil, Salvoconducto, Visa, Visa Refugiado | Alfanumérico, hasta 15 caracteres |
| PEP, PEP-FF | Alfanumérico |

> **GAP:** el Excel no especifica formato para Carné Diplomático más allá de "Admite Numérico" — confirmar longitud con el equipo antes de implementar la máscara de validación.

---

## Módulo DP — Datos personales (bloque Víctima)

> Incluye los campos de **residencia**, que el Excel/PDF agrupan aquí (no en
> CI como hacían el prototipo y V2-real) — ver nota de reorganización en
> `registro-v3-data.md`. Incluye además "Correo electrónico de la víctima"
> (Q5b), campo nuevo agregado por retroalimentación directa del cliente
> (obs. #12), que no existe en el Excel de Drive.

| # | Pregunta | Tipo | Requerida | Opciones / Notas |
|---|---|---|:---:|---|
| 1 | ¿Cuál es su fecha de nacimiento? | date | 1er contacto | |
| 2 | "Dada su condición física, mental y/o sensorial... tiene dificultades para: ¿oír, hablar, ver, mover su cuerpo, agarrar objetos, entender, comer o vestirse, interactuar?" | boolean | 1er contacto | Gate — habilita Q3 |
| 3 | ¿Qué actividades no puede o presenta dificultades para realizarlas? | multiple + escala | máx. 2º contacto · condicional (Q2=Sí) | 9 actividades, cada una en escala 1-4: Oír la voz/sonidos · Hablar/conversar · Ver de cerca/lejos/alrededor · Mover el cuerpo/caminar/subir-bajar escaleras · Agarrar/mover objetos con las manos · Entender/aprender/recordar/decidir · Comer/vestirse/bañarse · Relacionarse/interactuar · Actividades diarias sin problemas cardiacos/respiratorios |
| 4 | Según la Ley 1996 de 2019 cuenta usted con | multiple | máx. 2º contacto · **condicional** | Valoración de apoyos · Acuerdos de apoyo · Directivas anticipadas · Ninguna de las anteriores — **visible solo si CI.Q7 "¿Se reconoce como persona con discapacidad?" = Sí** (corregido por obs. #7 del cliente; el Excel de Drive no marcaba esta condicionalidad explícitamente) |
| 5 | ¿Cuál es su nacionalidad? | dropdown | 1er contacto · condicional | Colombiana(o) *(prediligenciado si tipo doc. es colombiano)* · Extranjero(a) *(prediligenciado si no)* · **Apátrida** (persona sin nacionalidad) |
| 5b | Correo electrónico de la víctima | text | **Deseable (opcional)** | **Campo nuevo, no está en el Excel de Drive** — agregado por obs. #12 del cliente ("Incluir en Datos personales pregunta opcional por 'Correo electrónico de la víctima'") |
| 6 | ¿De qué país es su nacionalidad? | dropdown | máx. 2º contacto · condicional (Q5=Extranjero) | Lista de países |
| 7 | ¿Cuál es tu condición migratoria actualmente? | dropdown | Deseable · condicional (no se muestra si colombiano o Nacionalidad=Colombiana) | Visa vigente · Cédula de extranjería · PPT · PEP vencido · Salvoconducto SC-2 (refugio) · Salvoconducto SC-1 (trámite Migración) · Ingreso sin control migratorio · Documento vencido · Permanencia superior al tiempo autorizado · Refugiada(o) reconocida(o) · Colombiano(a) retornado(a) · No informa |
| 8 | ¿Cuál es su identidad de género? | dropdown | 1er contacto | Mujer · Mujer Trans · Hombre · Hombre Trans · No binaria - asignado masculino al nacer · No binaria - asignado femenino al nacer · Otro |
| 9 | ¿Cuál es su orientación sexual? | dropdown | 1er contacto · condicional | Heterosexual · Gay *(inhabilitada si identidad=Mujer/Mujer Trans)* · Lesbiana *(inhabilitada si identidad=Hombre/Hombre Trans)* · Bisexual · Otro · No estoy segura/o/e · No informa/prefiere no responder |
| 10 | ¿Cuál fue su sexo asignado al nacer? | dropdown | máx. 2º contacto · condicional | Femenino · Masculino *(inhabilitado si identidad=No binaria-asignado femenino)* · Intersexual o indeterminado *(Masculino inhabilitado si identidad=No binaria-asignado masculino)* |
| 11 | ¿Pertenece a alguna población de especial protección? | multiple | 1er contacto | 17 opciones: Mujer buscadora de víctimas de desaparición forzada · Adolescentes/jóvenes en SRPA · Adultas(os) mayores · Líderes(as) sociales · NNAJ desvinculados del conflicto armado · Migrantes colombianos repatriados · Migrantes deportados/expulsados de Venezuela · NNAJ en restablecimiento de derechos · Programa de protección a testigos · Habitante de calle · Privada de la libertad · Desmovilizada/firmante Acuerdo de Paz · Víctimas del conflicto armado · Veteranas(os) Fuerza Pública · Voluntarias(os) Defensa Civil/Cruz Roja/Bomberos · Mujer cabeza de familia · Ninguno |
| 12 | ¿Se reconoce como perteneciente a algún grupo étnico? | dropdown | 1er contacto | Afrocolombiana(o) · Negra(o) · Palenquera(o) · Raizal (San Andrés/Providencia/Santa Catalina) · Indígena · Pueblo ROM/Gitano(a) · No se autorreconoce/No informa |
| 13 | ¿Se reconoce como perteneciente a un pueblo indígena? | dropdown (lista cerrada estandarizada) | 1er contacto · condicional (Q12=Indígena) | Lista de pueblos indígenas de Colombia |
| 14 | ¿Se reconoce como persona campesina? | boolean | 1er contacto | |
| 15 | ¿Cuál es su estado civil? | dropdown | máx. 2º contacto | Soltera(o) · Unión libre · Casada(o) · Separada(o) · Divorciada(o) · Viuda(o) · No informa |
| 16 | ¿Cuál es el último nivel de escolaridad que ha finalizado? | dropdown | máx. 2º contacto | Sin educación formal · Primaria (hasta 5°) · Secundaria básica (hasta 9°) · Bachillerato/Media (11°) · Técnico o tecnológico · Profesional · Especialización o maestría · Doctorado · No informa |
| 17 | ¿Cuál es su ocupación? | dropdown | 1er contacto | 24 opciones (ver lista completa abajo) |
| 18 | ¿Cuál ha sido su forma de generación de ingresos actualmente o en los últimos 3 meses? | dropdown | máx. 2º contacto | Empleada(o) contrato formal · Empleada(o) sin contrato/informal · Desempleada(o) · Pensionada(o)/jubilada(o) · Trabajadora independiente · Sin forma de generación de ingresos · Persona en ASP · No informa |
| 19 | ¿Cuál es su tipo de vinculación laboral? | dropdown | Deseable · condicional (Q18=contrato formal o independiente) | Planta · Provisional · Contratista/OPS · Tercerización · Contrato de aprendizaje · Pasantía |
| 20 | ¿Cuál es la modalidad en la que desarrolla actividades sexuales pagas (ASP)? | multiple | 1er contacto · condicional (Q17 o Q18 = Persona en ASP) | Servicio sexual presencial con contacto físico · Servicio sexual virtual · Producción/venta de contenido sexual pago · Acompañamiento/compañía remunerada |
| 21 | ¿Hace cuánto tiempo inició en las ASP aproximadamente? | date | 1er contacto · condicional (igual que Q20) | |
| 22 | ¿Cuál fue el principal motivo por el cual inició en las ASP? | multiple | 1er contacto · condicional (igual que Q20) | Decisión propia/elección autónoma · Motivos económicos · Presión o influencia de terceros · Migración/ausencia de redes de apoyo · No informa · Otro |
| 23 | ¿Cuál es la forma de tenencia de la vivienda donde reside actualmente? | dropdown | máx. 2º contacto | Propia · Arrendada · Familiar (sin pago) · En ocupación/usufructo · Albergue temporal/refugio · Sin vivienda (calle/desplazamiento) · No informa |
| 24 | ¿Cuál es el estrato de la vivienda en la que reside actualmente? | dropdown | máx. 2º contacto | Sin estrato · 1 · 2 · 3 · 4 · 5 · 6 · No aplica/No sabe/No informa |
| 25 | ¿Tiene personas o animales a su cargo? | multiple | máx. 2º contacto | No, nadie · Sí hijos(as) <18 · Sí hijos(as) ≥18 · Sí padre/madre · Sí pareja · Sí otros familiares · Sí otras personas no familiares · Sí animales |
| 26 | ¿Actualmente se encuentra en estado de gestación (embarazo)? | dropdown | 1er contacto · condicional | Sí · No · No informa — **visible si `sexo_asignado_nacer = Femenino` O `identidad_genero ∈ {Mujer, Hombre Trans, No binaria-asignado femenino al nacer}`** (condición corregida/ampliada por obs. #8 del cliente — el Excel de Drive solo mencionaba el criterio de sexo asignado) |
| 27 | ¿En qué departamento reside actualmente? | dropdown | 1er contacto | Lista departamentos |
| 28 | ¿En qué municipio reside actualmente? | dropdown | 1er contacto · condicional (depende de Q27) | Lista municipios/área no municipalizada según departamento |
| 29 | ¿Cuál es la dirección de su lugar de residencia? | text | 1er contacto | |
| 30 | ¿En qué zona se encuentra su lugar de residencia? | dropdown | máx. 2º contacto | Cabecera municipal · Parte rural (vereda/campo) · Centro poblado (corregimiento/inspección/caserío) |

### Lista completa — Ocupación (Q17)

Servicios generales y aseo · Cocina y alimentos · Belleza y cuidado personal · Ama de casa · Cuidado (niñera, cuidado adultos mayores/discapacidad, madres comunitarias) · Trabajo doméstico remunerado · Comercio y ventas · Atención al cliente/call center · Seguridad · Transporte y mensajería · Construcción y oficios de obra · Manufactura y producción · Administrativo y oficina · Salud y cuidado · Educación · Agro/rural · Tecnología/plataformas · Estudiante · Persona en ASP · Defensora de derechos humanos · Funcionaria(o) pública(o) · Periodista · Turista/visitante temporal · Empresaria(o) · Dirigente sindical · Afiliada(o) sindical · No informa · Otra

---

## Módulo CA — Contacto de apoyo o emergencia (bloque Víctima)

> **Cambio de criterio de obligatoriedad vs. V2-real:** el V2-real tenía esta
> sección como completamente opcional. El Excel/PDF la marca **obligatoria
> máx. 2º contacto** en sus 4 preguntas — ya no es opcional, solo tiene gracia
> hasta el segundo contacto.

| # | Pregunta | Tipo | Requerida | Opciones / Notas |
|---|---|---|:---:|---|
| 1 | ¿Podría indicarme el nombre de una persona de contacto de apoyo o emergencia en Colombia? | text | máx. 2º contacto | Validación "solo letras" (obs. #13 del cliente) |
| 2 | ¿Cuál es el número de teléfono de ese contacto de apoyo o emergencia? | text | máx. 2º contacto | |
| 3 | ¿Cuenta con un correo electrónico de ese contacto de apoyo o emergencia? | text | máx. 2º contacto | |
| 4 | ¿Cuál es el parentesco o relación con esa persona de contacto? | dropdown | máx. 2º contacto | Madre o padre · Hermano(a) · Hijo(a) · Pareja o expareja *(solo si no representa riesgo)* · Otro familiar · Amigo(a)/conocido(a) · Jefe(a)/colega de trabajo · Otro · No informa |

---

## Módulo HV — Hechos de violencia (bloque Caso)

> **"¿Cuál es el número de su documento de identificación?"** aparecía en el
> Excel original de este módulo pero está marcada explícitamente **"Se
> elimina de Formulario de Ingreso Registro, va para el módulo de Barreras"**
> — no se incluye en Registro V3.

| # | Pregunta | Tipo | Requerida | Opciones / Notas |
|---|---|---|:---:|---|
| 1 | ¿Podría describir brevemente los hechos ocurridos? | text | 1er contacto | |
| 2 | ¿En qué fecha ocurrió el hecho? | date | 1er contacto | Admite fecha aproximada |
| 3 | ¿A qué hora ocurrió el hecho? | time | máx. 2º contacto | Admite hora aproximada |
| 4 | ¿En qué departamento ocurrieron los hechos? | dropdown | 1er contacto | |
| 5 | ¿En qué municipio ocurrieron los hechos? | dropdown | 1er contacto · condicional (depende de Q4) | |
| 6 | ¿En qué zona ocurrió el hecho? | dropdown | 1er contacto | Cabecera municipal · Parte rural · Centro poblado |
| 7 | ¿Cuál es la dirección donde ocurrieron los hechos (o nombre de la finca si es rural)? | text | máx. 2º contacto | |
| 8 | ¿En qué escenario ocurrieron los hechos? | dropdown | 1er contacto | 36 opciones — ver lista completa abajo |
| 9 | ¿Qué tipo o tipos de violencia ha experimentado? | multiple | 1er contacto | Física · Psicológica/emocional · Sexual · Patrimonial o económica · **Política** · **Reproductiva** · **Vicaria** (a través de hijos/as u otras personas dependientes) |
| 10 | ¿Qué subtipo o subtipos de violencia se presentaron? | multiple | máx. 2º contacto · condicional (depende de Q9) | Ver tabla de subtipos abajo — incluye **"Tentativa Feminicidio"** y **"Tentativa Homicidio abyecto o fútil"** como habilitantes directos del tipo de Tamizaje (ver `registro-v3-data.md` § Tamizaje) |
| 11 | ¿En qué ámbito o ámbitos se presentó la violencia? | multiple (máx. 2) | 1er contacto | Educativo · Conflicto armado · Salud · Institucional · Intrafamiliar (conviviente) · Laboral · Político · Reclusión intramural · Relacional (no convive) · Digital · Comunitario/vecinal/organizativo · Deportivo · Cultural/entretenimiento/ocio · Espacios sindicales |
| 12 | ¿En qué sector se presentó la violencia en el ámbito laboral? | dropdown | 1er contacto · condicional (Q11 incluye Laboral) | Sector público · Sector privado · Economía informal · Modalidades de tercerización |
| 13 | ¿Considera que la violencia se motivó en razón de su género, identidad de género y/o orientación sexual? | boolean | 1er contacto | |
| 14 | ¿Cuántas veces ha ocurrido esta agresión por parte de la principal persona agresora? | dropdown | 1er contacto | Es la primera vez · Es la segunda vez · Ha sido recurrente (más de 2 veces) |

### Lista completa — Escenario de la violencia (Q8)

Ambulancia/transporte sanitario · Áreas deportivas/recreativas · Calle · Carretera · Centro de atención médica · Centros de reclusión · Centros educativos · Espacios acuáticos al aire libre · Espacios terrestres al aire libre · Establecimiento comercial · Establecimiento industrial/obras · Establecimientos de comidas · Administración pública · Establecimientos financieros · Estaciones de servicio · Guarniciones militares/Policía · Actividades culturales · Lugares de cuidado de personas · Esparcimiento con alcohol · Explotación de minas/canteras · Hospedaje · Oficinas de empresas privadas · Parqueaderos · Sitio de culto · Terminales/aeropuertos · Transporte masivo · Vehículo particular · Vehículo público · Vivienda · Zonas agropecuarias · Espacios sindicales · Teletrabajo/trabajo remoto

### Tabla de subtipos de violencia (Q10) — condicional por tipo seleccionado en Q9

| Si Tipo incluye... | Subtipos disponibles |
|---|---|
| Física | Ataque con agente químico · Golpizas · Quemadura · Sofocamiento · Agresión con objetos · Tortura · Secuestro · **Tentativa Feminicidio** · **Tentativa Homicidio abyecto o fútil** |
| Psicológica/emocional | Subvaloración/descalificación/humillación/violencia verbal · Intimidación y amenaza · Acoso/hostigamiento/vigilancia/seguimiento · Control coercitivo por celos · Aislamiento · Acoso laboral |
| Sexual | Acceso carnal violento · Abuso sexual · Trata con fines de explotación sexual · Mutilación de órganos sexuales · Acoso sexual |
| Patrimonial o económica | Inasistencia alimentaria · Vigilancia/manipulación de dinero · Negación/retención/destrucción de bienes · Remuneración inequitativa |
| Política | Abuso de autoridad con fines políticos · Censura y persecución política |
| Reproductiva | Aborto forzado · Esterilización/planificación forzada · Denegación de servicios IVE · Matrimonio infantil/uniones tempranas · Violencia obstétrica |
| Vicaria | Amenazas/hostigamientos contra hijas/os · Privación de libertad y secuestro parental |

---

## Módulo AG — Características del agresor (bloque Caso)

| # | Pregunta | Tipo | Requerida | Opciones / Notas |
|---|---|---|:---:|---|
| 1 | ¿Cuántas personas agresoras participaron en los hechos? | dropdown | 1er contacto | 1 · 2 · 3 · Más de 3 · No se puede determinar/no informa |
| 2 | ¿Conoce y convive con la principal persona agresora? | dropdown | 1er contacto | Persona conocida con la que convive · Persona conocida con la que no convive · Persona desconocida · No se puede determinar/no informa |
| 3 | ¿Qué relación tiene o tenía con la persona agresora? | dropdown | 1er contacto | 18 opciones (ver abajo) — **Pareja íntima** y **Expareja/antigua pareja** son las que habilitan la rama "pareja" del Tamizaje |
| 4 | ¿Qué ocupación tiene la persona agresora? | dropdown | máx. 2º contacto | 19 opciones (ver abajo) |
| 5 | ¿Actualmente depende económicamente de la persona agresora? | boolean | 1er contacto · condicional (Q2 = conoce, con o sin convivencia) | |
| 6 | ¿Cuál es la identidad de género percibida de la principal persona agresora? | dropdown | máx. 2º contacto | Hombre · Hombre Trans · Mujer · Mujer Trans · Otro · No identifica |
| 7 | ¿Cuáles son los nombres y apellidos de la persona agresora? | text | Deseable | Validación "solo letras" (obs. #13 del cliente) |
| 8 | ¿Qué tipo de documento de identificación tiene la persona agresora? | dropdown | Deseable | |
| 9 | ¿Cuál es el número del documento de identificación de la persona agresora? | text | Deseable | |
| 10 | ¿Cuál es la dirección de la persona agresora (o finca, si rural)? | text | Deseable | |
| 11 | ¿Cuál es el número de teléfono de la persona agresora? | text | Deseable | No admite letras |

### Lista completa — Relación con el presunto agresor (Q3)

Amigo(a) · Encargado(a) del cuidado · Hermano(a) · Madre · Otros familiares civiles/consanguíneos (primo, abuelo, tío, cuñado, suegro, madrastra, padrastro) · Padre · **Pareja íntima** · **Expareja o antigua pareja** · Representante entidad contratante/interventoría (contratación tercerizada) · Supervisor(a)/empleador(a)/jefe(a) · Compañero(a) de trabajo · Profesor(a)/instructor(a)/coach · Compañero(a) de estudios · Proveedor(a) de servicios · Coinquilino(a) · Amigo(a) de la familia · Vecino(a) · Cliente(s) · Sin relación · Otro

### Lista completa — Ocupación del agresor (Q4)

Delincuencia común · Persona desmovilizada · Miembro de grupos alzados al margen de la ley · Personal de seguridad privada · Personal de custodia · Servidor(a) público(a) · Líder(esa) religioso(a)/espiritual · Representante entidad contratante/interventoría · Miembro Policía Nacional · Miembro Ejército Nacional · Miembro Armada Nacional · Miembro Fuerza Aérea · Docente/instructor(a)/coach · Personal institución de salud · Líder(esa) social · Afiliado(a) sindical · Dirigente sindical · Otro · No conoce/No informa

---

## Módulo TZ — Tamizaje de riesgo (bloque Caso)

> Ver la lógica completa de bifurcación en `registro-v3-data.md` § D3/D4.
> **Corregido por feedback directo del cliente** sobre dos puntos que
> cambian respecto al Guion Funcional de Drive:
> - El **tipo** de riesgo (Feminicidio vs. Homicidio abyecto o fútil) se
>   decide por `identidad_genero`/`orientación_sexual`/`sexo_asignado_nacer`
>   de DP, **no** por la respuesta a la pregunta fuerte VPP (pregunta 8 de
>   esta lista) — esa pregunta se mantiene, pero solo suma al score (D3).
> - **No todas las preguntas se muestran siempre:** de las 7 "generales", 5
>   se muestran primero; si las 5 responden "No", el Tamizaje termina ahí
>   (nivel Bajo) sin mostrar el resto. Si alguna es "Sí", se revela todo lo
>   demás (D4 — cuáles 5 exactamente es GAP G-V3-08, propuesta marcada abajo).

### Preguntas generales (7 — las marcadas 🔵 son la propuesta de "5 iniciales obligatorias", GAP G-V3-08)

| # | Pregunta | ¿Inicial obligatoria? (propuesta GAP G-V3-08) |
|---|---|:---:|
| 1 | ¿Ha aumentado la violencia física en severidad o frecuencia en los últimos 60 días? | 🔵 Sí |
| 2 | ¿Le ha agredido o amenazado utilizando algún arma u objeto en su contra? | 🔵 Sí |
| 3 | ¿Ha amenazado con matarle? | 🔵 Sí |
| 4 | ¿Ha intentado asfixiarle, estrangularle o impedirle respirar? | 🔵 Sí |
| 5 | ¿Esta persona le persigue, espía o llega a lugares donde no le está esperando? | Se revela después |
| 6 | ¿Le deja notas o mensajes amenazantes en diferentes medios? | Se revela después |
| 7 | ¿La persona agresora tiene acceso a armas? | 🔵 Sí |

### Pregunta fuerte VPP (condicional — solo si la víctima manifestó OSIGD no-hegemónica)

| # | Pregunta | Rol |
|---|---|---|
| 8 | ¿Esta persona ha realizado amenazas, agresiones o comentarios motivados por prejuicio hacia su orientación sexual, identidad o expresión de género? | Solo suma 1 punto al score — **ya no decide el tipo de riesgo** (corregido por obs. #5 del cliente, ver D3) |

### Rama Pareja/Expareja (23 preguntas — si relación con agresor = Pareja íntima o Expareja/antigua pareja)

| # | Pregunta |
|---|---|
| 1 | ¿Está su pareja/expareja desempleado(a) o sin trabajo actualmente? |
| 2 | ¿Sabe si su pareja/expareja ha tenido algún proceso con la justicia (lesiones, amenazas, porte ilegal de armas, hurto, homicidio, abuso sexual, etc.)? |
| 3 | ¿Sabe si ha tenido denuncias/procesos jurídicos por violencia intrafamiliar o incumplimiento de medidas de protección? |
| 4 | ¿Su pareja/expareja consume drogas (anfetamina, cocaína, heroína u otras)? |
| 5 | *(condicional a la 4)* ¿Le persuade u obliga a consumirlas? |
| 6 | ¿Su pareja/expareja es alcohólica o tiene problemas con el alcohol? |
| 7 | ¿Es celoso/a de manera constante o violenta y/o controla sus actividades, decisiones, desplazamientos o lugares? |
| 8 | ¿Limita, impide o amenaza su contacto con familiares, amistades, redes de apoyo o instituciones? |
| 9 | ¿La ha forzado a mantener relaciones sexuales con ella o con otras sin que usted quiera? |
| 10 | ¿Le vigila, persigue, acosa, contacta insistentemente o utiliza terceros para intimidar/amenazar/presionar? |
| 11 | ¿La acosa sexualmente o remite mensajes de índole sexual no pedidos e incómodos? |
| 12 | ¿La ha amenazado con suicidarse o lo ha intentado? |
| 13 | ¿Ha difundido, amenazado con difundir o utilizado fotos/videos/mensajes/redes/información personal para exponerle, humillarle, vigilarle, extorsionarle o ponerle en riesgo? |
| 14 | ¿Le ha impedido acceder a servicios de salud o salud mental? |
| 15 | ¿Utiliza su estado de salud o diagnósticos (VIH/ITS) o el acceso a tratamientos para chantajearla o amenazar con revelarlos? |
| 16 | ¿La ha agredido físicamente mientras estaba en enfermedad, recuperación médica, posoperatorio, embarazo u otra vulnerabilidad física/médica? |
| 17 | ¿La ha amenazado con quitarle, restringir comunicación/visitas o hacer daño a hijos, hijas, mascotas u otros familiares? |
| 18 | ¿Ha destruido o retenido objetos personales, documentos, medicamentos, tratamientos u otras pertenencias como forma de intimidación/humillación/control? |
| 19 | ¿Le ha expuesto intencionalmente a ITS o le ha impedido tomar medidas de protección en relaciones sexuales? |
| 20 | ¿Actualmente convive con su pareja/expareja o comparte espacios frecuentes? |
| 21 | ¿Ha tenido que cambiar rutinas, rutas, residencia o espacios frecuentes por temor a agresiones/represalias/amenazas? |
| 22 | *(condicional — solo si OSIGD)* ¿Ha amenazado con revelar su orientación sexual/identidad/expresión de género sin su consentimiento? |
| 23 | *(condicional — solo si OSIGD)* ¿Ha dejado de buscar ayuda institucional tras una agresión por temor a discriminación o no recibir protección por su orientación sexual/identidad de género? |
| 24 | **Ideación suicida** — ¿como consecuencia de la violencia ha tenido pensamientos para terminar con su vida? | Marcaje y abordaje de remisión psicosocial |

> Nota: la numeración de la fuente llega a 24 ítems en esta rama contando la
> condicional de consumo de sustancias como ítem propio — el Guion Funcional
> narrativo dice "23 ítems"; se documenta la lista completa tal como aparece
> en el Excel (24 filas) y se deja la discrepancia de conteo como detalle
> menor, no afecta la lógica.

### Rama No Pareja (24 preguntas — cualquier otra relación con el agresor)

| # | Pregunta |
|---|---|
| 1 | ¿La persona agresora está desempleada o sin trabajo actualmente? |
| 2 | ¿Sabe si tiene antecedentes penales, denuncias por otros delitos, procesos judiciales o antecedentes de incumplimiento de medidas de protección? |
| 3 | ¿Consume drogas (anfetamina, cocaína, heroína u otras)? |
| 4 | *(condicional)* ¿Le persuade u obliga a consumirlas? |
| 5 | ¿Es alcohólica o tiene problemas con el alcohol? |
| 6 | ¿Es celoso/a de manera constante o violenta y/o controla sus actividades, decisiones, desplazamientos o lugares? |
| 7 | ¿Limita, impide o amenaza su contacto con familiares, amistades, redes de apoyo o instituciones? |
| 8 | ¿La ha forzado a mantener relaciones sexuales con ella o con otras sin que usted quiera? |
| 9 | ¿Le vigila, persigue, acosa, contacta insistentemente o utiliza terceros para intimidar/amenazar/presionar? |
| 10 | ¿La acosa sexualmente o remite mensajes de índole sexual no pedidos e incómodos? |
| 11 | ¿Ha difundido, amenazado con difundir o utilizado fotos/videos/mensajes/redes/información personal para exponerle, humillarle, vigilarle, extorsionarle o ponerle en riesgo? |
| 12 | ¿Le ha impedido acceder a servicios de salud o salud mental? |
| 13 | ¿Ha aprovechado su estado de vulnerabilidad física (embarazo, posoperatorio de afirmación de género, condiciones crónicas) o mental para agredirla/hostigarla/amenazarla? |
| 14 | ¿La ha amenazado con quitarle, restringir comunicación/visitas o hacer daño a hijos, hijas, mascotas u otros familiares? |
| 15 | ¿Ha destruido o retenido objetos personales, documentos, medicamentos, tratamientos u otras pertenencias como forma de intimidación/humillación/control? |
| 16 | ¿Le ha expuesto intencionalmente a ITS o le ha impedido tomar medidas de protección en relaciones sexuales? |
| 17 | ¿Comparte espacios frecuentes con esta persona (barrio, trabajo, estudio, transporte, residencia, espacios comunitarios)? |
| 18 | ¿Ha tenido que cambiar rutinas, rutas, residencia o espacios frecuentes por temor a agresiones/represalias/amenazas? |
| 19 | *(condicional — solo si OSIGD)* ¿Ha amenazado con revelar su orientación sexual/identidad/expresión de género sin su consentimiento? |
| 20 | *(condicional — solo si OSIGD)* ¿La violencia es ejercida por miembros de su familia (expulsión de hogar, amenazas, control económico, daño a pertenencias, intentos de "corrección")? |
| 21 | *(condicional — solo si OSIGD)* ¿Ha dejado de buscar ayuda institucional tras una agresión por temor a discriminación o no recibir protección por su orientación sexual/identidad de género? |
| 22 | *(condicional — solo si ámbito=laboral)* ¿Existe alguna relación de jerarquía entre usted y la persona agresora (profesor, jefe, superior, líder religioso, policía, militar, actor armado, etc.)? |
| 23 | *(condicional — solo si ámbito=laboral)* ¿Ha utilizado su posición de autoridad (jefe, arrendador, líder) para amenazarla con pérdida de empleo, vivienda o acceso a servicios si no accede a sus exigencias? |
| 24 | **Ideación suicida** — ¿como consecuencia de la violencia ha tenido pensamientos para terminar con su vida? | Marcaje y abordaje de remisión psicosocial |

---

## Módulo GS — Gestión Salvia (bloque Caso: Denuncia virtual + Remisión + Entidades + Barreras + Alerta)

> Al entrar a este módulo se muestra primero el **párrafo síntesis** —
> narrativa auto-generada de la víctima y el caso. Ver
> [registro-v3-parrafo-sintesis.md](registro-v3-parrafo-sintesis.md)
> (agregado por obs. #18 del cliente, no estaba en el Excel de Drive).
>
> **Reestructurado a partir del formulario de Seguimiento real** (BD de
> producción — form_id `2d0aeb46-1af3-4c47-a0d5-c5bfc4d549ff`, secciones
> "Seguimiento de Caso", "Identificación de Entidades" e "Identificación de
> Barreras`). El Excel/Guion Funcional de Drive solo daba un modelo
> simplificado de 4 preguntas para remisión y nada para entidades/barreras
> — eso se reemplaza aquí por el modelo real ya probado en producción, con
> prioridad sobre el Excel (mismo criterio que D8: cuando el V2-real/
> Seguimiento ya resolvió algo en código real, se reutiliza en vez de
> reinventar). Ver `registro-v3-data.md` § D9 para el detalle de la
> investigación en BD y el mapeo.
>
> Resuelve GAP G-V3-04 (catálogo de criterios de remisión) y GAP G-V3-12
> (conexión con Rutas/Barreras/Notificaciones): los botones del sidebar
> ahora enlazan a la vista histórica/agregada de Entidades y Barreras del
> Expediente de víctima (diferido, D2); el registro de ESTE Caso se hace
> aquí, dentro de GS.

### A — Denuncia virtual y gestión general

| # | Pregunta | Tipo | Requerida | Opciones / Notas |
|---|---|---|:---:|---|
| 1 | ¿Usted permitiría que el Sistema SALVIA realice el proceso de denuncia fácil ante la FGN? | boolean | 1er contacto | |
| 2 | ¿Es necesario realizar gestiones interinstitucionales a partir de los hallazgos de esta atención? | multiple | 1er contacto | Orientación y enrutamiento · Gestión administrativa mediante llamada · Activación de ruta interinstitucional · Articulación interinstitucional · Escalamiento a organismos de control *(alineado con la pregunta real de Seguimiento — reemplaza el set de opciones anterior "Orientación/Activación de ruta/Articulación/Gestión de barreras", que mezclaba esta pregunta-resumen con el repeater C de abajo)* |
| 3 | En caso afirmativo, describa las acciones previstas, las entidades involucradas y su propósito | text | condicional (Q2 no vacío) | |
| 4 | Describa el acompañamiento realizado y las acciones que se desarrollarán para la gestión del caso | text | 1er contacto | Se guarda como la **primera entrada de la bitácora** del Caso en Seguimiento |

### B — Remisión a equipos SALVIA

> Reemplaza el modelo anterior (4 "procesos" genéricos sin criterios). Son
> **6 equipos reales**, cada uno con su propia regla de habilitación.

| # | Pregunta | Tipo | Requerida | Opciones / Notas |
|---|---|---|:---:|---|
| 5 | ¿Durante la atención se han generado necesidades que requieran remisión a los equipos SALVIA? | boolean | 1er contacto | |
| 6 | ¿Cuáles equipos? | multiple | condicional (Q5=Sí) | Medidas de Emergencia · Atención Psicosocial · Estabilización · Atención Hombres · Salvia Dignidad · Discapacidad |
| 6b | *(banner info)* | info | — | Si Q6 incluye **Medidas de Emergencia y Atención Psicosocial a la vez** → advertencia: "no se puede remitir a estos dos equipos simultáneamente" (regla real de Seguimiento, no bloquea — solo advierte) |
| 7 | ¿Cuáles medidas de emergencia? | multiple | condicional (Q6 incluye Medidas de Emergencia) | Alojamiento · Transporte · Alimentación · Vestuario · Apoyo psicosocial · Otras ME |
| 8 | *(banner info)* | info | — | Visible si Q6 incluye Atención Psicosocial: "Para que la remisión sea válida se debe marcar el criterio obligatorio y sumar mínimo 3 puntos" |
| 9 | Criterios de remisión — Atención Psicosocial | multiple | condicional (Q6 incluye Atención Psicosocial) | **Obligatorio**: "Mujeres en todas sus diversidades y personas no binarias, mayores de 18 años, víctimas de VBG o VpP con voluntariedad para la atención y en post-emergencia" · Opcionales con puntaje: Conducta suicida derivada de VBG/VpP (3 pts) · Interseccionalidad (2 pts) · Que no haya activado la ruta (1 pt) · Condiciones territoriales/orden público que impiden la atención (1 pt) · Sin acceso a atención psicosocial/psicológica (1 pt) · Naturalización de la VBG y roles de género (1 pt) |
| 10 | Criterios de remisión — Atención Hombres | multiple | condicional (Q6 incluye Atención Hombres) | **Obligatorio, único**: "Hombres diversos victimarios o posibles victimarios" |
| 11 | Seleccione los servicios del Equipo de Discapacidad | multiple | condicional (Q6 incluye Discapacidad) | Apoyo comunicativo en Lengua de Señas Colombiana (LSC) · Apoyo con enfoque de discapacidad (ley de capacidad legal) |
| 12 | Criterios de remisión — Estabilización | multiple | condicional (Q6 incluye Estabilización) | Mínimo 1: No tiene ingresos o son inestables · Reside en vivienda temporal o familiar · Tiene personas a cargo que limitan su tiempo para trabajar/estudiar |
| 13 | Describa los elementos que evidencia para realizar la remisión | text | condicional (Q5=Sí) | |

> **Salvia Dignidad** no tiene preguntas de criterios propias en el
> formulario real — se remite solo con la selección en Q6.
> **GAP G-V3-18**: el formulario real de Seguimiento también tiene un
> mecanismo de *reasignación completa del Caso a otro equipo* (distinto de
> remisión), con sus propias 4 opciones (Riesgo/General/Hombres/Salvia
> Dignidad) — no se incluye en Registro V3 porque reasignar implica que el
> Caso ya tiene un equipo asignado (ocurre en Seguimiento, no en el primer
> Registro) — confirmar si aplica igual desde el primer contacto.

### C — Registrar Entidades (repeater)

> Modela el repeater real "Identificación de Entidades" de Seguimiento —
> dos repeaters paralelos según el disparador.

| # | Pregunta | Tipo | Requerida | Opciones / Notas |
|---|---|---|:---:|---|
| 14 | ¿Ha puesto en conocimiento de alguna entidad/institución los hechos de violencia reportados? | boolean | 1er contacto | Habilita el repeater **Entidades Acudidas** |
| 15 | ¿Está usted de acuerdo con que SALVIA realice la activación de la Ruta de Atención? | single | condicional (Q14=Sí o Q16=Sí) | Sí · No |
| 16 | ¿Salvia debe activar rutas en el marco de la debida diligencia? | boolean | 1er contacto | Habilita el repeater **Entidades Activación** (puede ser Sí aunque Q14 sea No — es la propia SALVIA actuando de oficio) |

**Por cada entidad agregada** (repeater Acudidas o Activación, mismos campos en ambos):

| Campo | Tipo | Requerida | Opciones / Notas |
|---|---|:---:|---|
| Sector de la entidad | dropdown | Sí | Salud · Justicia · Protección · Otras instituciones |
| Departamento | dropdown | Sí | Lista departamentos |
| Ciudad | dropdown | Sí | Cascada dependiente del departamento |
| Municipio | dropdown | Sí | Cascada dependiente de la ciudad |
| Entidad | dropdown | Sí | Catálogo dinámico por sector (`entidadesDisponibles`, igual que Seguimiento) |
| N° de radicado, NUNC, MP, etc. | text | No | |
| Fecha de reporte | date | Sí | |
| Canal de activación | multiple | Sí (solo repeater Activación) | Notificación (tarea) · Llamada (tarea) · Mi Salvia — Gestión institucional externa (tarea) |
| Acciones institucionales completadas | multiple | No | Catálogo `obligacionesCatalogo` + "Otra" (condicional: texto libre) |
| Acciones institucionales pendientes | multiple | No | Mismo catálogo + "Otra" (condicional: texto libre) |
| Otra información asociada a la Ruta | text | Sí (solo repeater Activación) | |

### D — Registrar Barreras (repeater)

> Modela el repeater real "Identificación de Barreras" de Seguimiento. En
> Seguimiento el disparador es una pregunta de "Valoración del Riesgo" que
> Registro V3 no tiene — se adapta a una pregunta propia (Q17, nueva, sin
> equivalente directo en ninguna fuente — ver GAP).

| # | Pregunta | Tipo | Requerida | Opciones / Notas |
|---|---|---|:---:|---|
| 17 | ¿Identificó barreras institucionales durante esta atención? | boolean | 1er contacto | **Nueva, sin equivalente en Excel/Guion Funcional/Observaciones — adaptada de Seguimiento, ver GAP G-V3-19.** Habilita el repeater de barreras |

**Por cada barrera agregada:**

| Campo | Tipo | Requerida | Opciones / Notas |
|---|---|:---:|---|
| Sector de la barrera | dropdown | Sí | Salud · Justicia · Protección · Otras instituciones |
| Barreras identificadas en Salud | multiple | No | Condicional (sector=Salud) — catálogo propio del sector |
| Otra barrera en Salud ¿cuál? | text | No | Condicional (incluye "Otra") |
| Barreras identificadas en Justicia | multiple | No | Condicional (sector=Justicia) — catálogo propio del sector |
| Otra barrera en Justicia ¿cuál? | text | No | Condicional (incluye "Otra") |
| Barreras identificadas en Protección | multiple | No | Condicional (sector=Protección) — catálogo propio del sector |
| Otra barrera en Protección ¿cuál? | text | No | Condicional (incluye "Otra") |
| Entidad a la que se asignará la barrera | dropdown | No | Catálogo dinámico, igual que Registrar Entidades |
| Barreras institucionales y de talento humano | multiple | No | Categoría transversal (aplica a cualquier sector) |
| Barreras económicas y socioeconómicas | multiple | No | Categoría transversal |
| Barreras territoriales y geográficas | multiple | No | Categoría transversal |
| Barreras por ausencia de enfoque diferencial | multiple | No | Categoría transversal |
| Fecha en la que se presentó la barrera (aproximada) | date | Sí | |
| Funcionario/a o dependencia donde se presentó la barrera | text | Sí | |
| Descripción de la barrera | text | Sí | |
| Gestión de la barrera | multiple | Sí | |

> **GAP G-V3-20**: no se obtuvo el catálogo exacto de opciones de cada
> multiselect (barreras por sector, categorías transversales, gestión de
> la barrera) — están en la tabla `option` real pero no se extrajeron en
> este análisis; están disponibles en `form-barreras-repeater.md` /
> `form-seguimiento-barreras-repeater.md` de la carpeta de Seguimiento si
> ya fueron documentados ahí, o se consultan directo en BD antes de
> implementar.

### E — Activar Alerta Inminente

> Complementa (no reemplaza) el botón de acceso rápido del header — ver
> `flow-E13`. Aquí queda como una acción formal y documentada dentro de la
> gestión del caso, con justificación.

| # | Pregunta | Tipo | Requerida | Opciones / Notas |
|---|---|---|:---:|---|
| 18 | Botón "Activar Alerta Inminente" | acción (no es pregunta de formulario) | — | Abre el mismo flujo que el botón del header (`flow-E13`) — aquí queda registrado dentro de la bitácora de gestión de este módulo, con campo de justificación |

---

## Totales

| Módulo | Preguntas | Notas |
|---|:---:|---|
| CI — Contacto e identificación | 9 | |
| DP — Datos personales (incl. residencia) | 31 | +9 preguntas de escala dentro de la Q3. Incluye correo víctima (Q5b, nueva por obs. #12 del cliente) |
| CA — Contacto de apoyo | 4 | Ya no es opcional — obligatoria máx. 2º contacto |
| HV — Hechos de violencia | 14 | "Número de documento" eliminado (va a Barreras) |
| AG — Características del agresor | 11 | |
| TZ — Tamizaje | 7 generales + 1 VPP condicional + 23 (pareja) **o** 24 (no pareja) | Ramas excluyentes |
| GS — Gestión Salvia | 18 campos propios + 2 repeaters (Entidades: 11 campos por entrada; Barreras: 15 campos por entrada) | Reestructurado sobre el Seguimiento real — remisión a 6 equipos con criterios, registro de entidades y de barreras, Alerta Inminente formal |
| **Total (rama pareja)** | **~106** | |
| **Total (rama no pareja)** | **~107** | |
