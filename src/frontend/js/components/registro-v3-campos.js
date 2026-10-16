// Catálogo de campos de Registro V3 — DocsMD/Screens/Registro V3/registro-v3-secciones.md.
// Sin dinamic-form (decisión A1): config JS plana consumida directamente
// por registro-v3-app.js. Los dropdowns/checkboxes que vienen de
// victim_case_form2_enums declaran enumCategory/enumCategories (coinciden
// con las claves que expone GET /victim-cases-v3/enums) en vez de una
// lista fija de 'opciones' -- así el icode que se envía en guardarSeccion
// siempre es el real de BD, nunca inventado en el frontend.

// D13: máscaras de validación de CI_6 (número de documento) por tipo de
// documento seleccionado en CI_5 -- tabla completa de secciones.md. Code =
// sufijo de victim_case_form2_victim_doc_type_* (Locale.go). null = no pide
// número (Persona sin identificación). Sin `regex` = sin restricción de
// formato más allá de "no vacío" (Carné Diplomático, PEP, PEP-FF -- la
// fuente solo dice "admite numérico"/no especifica longitud).
var MASCARA_DOC_NUMERO_V3 = {
  cc: { regex: /^[0-9]{6,10}$/, hint: "Numérico, 6 a 10 dígitos" },
  ce: { regex: /^[0-9]{6,10}$/, hint: "Numérico, 6 a 10 dígitos" },
  ti: { regex: /^[0-9]{10,11}$/, hint: "Numérico, 10 a 11 dígitos" },
  ps: { regex: /^[A-Za-z0-9]{6,9}$/, hint: "Alfanumérico, 6 a 9 caracteres" },
  pt: { regex: /^[0-9]{7,8}$/, hint: "Numérico, 7 a 8 dígitos" },
  si: null, // Persona sin identificación -- no pide número
  no: { regex: /^[A-Za-z0-9]{1,15}$/, hint: "Alfanumérico, hasta 15 caracteres" },
  nc: { regex: /^[A-Za-z0-9]{1,15}$/, hint: "Alfanumérico, hasta 15 caracteres" },
  mv: { regex: /^[A-Za-z0-9]{1,15}$/, hint: "Alfanumérico, hasta 15 caracteres" },
  dn: { regex: /^[A-Za-z0-9]{1,15}$/, hint: "Alfanumérico, hasta 15 caracteres" },
  oe: { regex: /^[A-Za-z0-9]{1,15}$/, hint: "Alfanumérico, hasta 15 caracteres" },
  pn: { regex: /^[A-Za-z0-9]{1,15}$/, hint: "Alfanumérico, hasta 15 caracteres" },
  rc: { regex: /^[A-Za-z0-9]{1,15}$/, hint: "Alfanumérico, hasta 15 caracteres" },
  sc: { regex: /^[A-Za-z0-9]{1,15}$/, hint: "Alfanumérico, hasta 15 caracteres" },
  vs: { regex: /^[A-Za-z0-9]{1,15}$/, hint: "Alfanumérico, hasta 15 caracteres" },
  vr: { regex: /^[A-Za-z0-9]{1,15}$/, hint: "Alfanumérico, hasta 15 caracteres" },
  cd: { regex: /^[0-9]+$/, hint: "Numérico" },
  pp: { regex: /^[A-Za-z0-9]+$/, hint: "Alfanumérico" },
  pf: { regex: /^[A-Za-z0-9]+$/, hint: "Alfanumérico" },
};

var DEPARTAMENTOS_V3 = [
  "Amazonas","Antioquia","Arauca","Atlántico","Bolívar","Boyacá","Caldas","Caquetá","Casanare",
  "Cauca","Cesar","Chocó","Córdoba","Cundinamarca","Guainía","Guaviare","Huila","La Guajira",
  "Magdalena","Meta","Nariño","Norte de Santander","Putumayo","Quindío","Risaralda",
  "San Andrés y Providencia","Santander","Sucre","Tolima","Valle del Cauca","Vaupés","Vichada",
  "Bogotá, D.C.",
];

// Etiquetas propias en español para categorías de victim_case_form2_enums
// que son EXCLUSIVAS de Registro V3 (se sembraron con name = "{categoria}_
// {codigo}", un slug sin traducción en Locale.go -- a diferencia de las
// categorías reutilizadas de V2, que sí tienen traducción real ahí). El
// renderer usa esto en vez del label que devuelve /victim-cases-v3/enums
// cuando existe una entrada para (categoria, code).
var ENUM_LABELS_OVERRIDE_V3 = {
  victim_case_form2_pronoun: { el: "Él", la: "Ella", le: "Elle" },
  victim_case_form2_referral_team: {
    me: "Medidas de Emergencia", ap: "Atención Psicosocial", es: "Estabilización",
    ah: "Atención Hombres", sd: "Salvia Dignidad", di: "Discapacidad",
  },
  victim_case_form2_emergency_measure: {
    al: "Alojamiento", tr: "Transporte", am: "Alimentación", ve: "Vestuario", ps: "Apoyo psicosocial", ot: "Otras ME",
  },
  victim_case_form2_psycho_criteria: {
    ob: "Mujeres en todas sus diversidades y personas no binarias, mayores de 18 años, víctimas de VBG o VpP con voluntariedad para la atención y en post-emergencia (obligatorio)",
    cs: "Conducta suicida derivada de VBG o VpP (3 pts)",
    in: "Interseccionalidad (2 pts)",
    sr: "Que no haya activado la ruta (1 pt)",
    ct: "Condiciones territoriales y de orden público que impiden la atención (1 pt)",
    sa: "Sin acceso a la atención psicosocial o psicológica (1 pt)",
    nv: "Naturalización de la Violencia Basada en Género y de los roles de género (1 pt)",
  },
  victim_case_form2_men_criteria: { ob: "Hombres diversos victimarios o posibles victimarios (obligatorio)" },
  victim_case_form2_disability_service: {
    ls: "Apoyo comunicativo en Lengua de Señas Colombiana (LSC)",
    ed: "Apoyo con enfoque de discapacidad (ley de capacidad legal)",
  },
  victim_case_form2_stabiliz_criteria: {
    si: "No tiene ninguna forma de ingresos o tiene ingresos de forma inestable",
    vt: "La persona reside actualmente en vivienda temporal o familiar",
    pc: "La persona tiene personas a cargo que limitan su tiempo para trabajar/estudiar",
  },
  victim_case_form2_interinst_action: {
    oe: "Orientación y enrutamiento", ga: "Gestión administrativa mediante llamada",
    ar: "Activación de ruta interinstitucional", ai: "Articulación interinstitucional",
    ec: "Escalamiento a organismos de control",
  },
  victim_case_form2_sector: { sa: "Salud", ju: "Justicia", pr: "Protección", ot: "Otras instituciones" },
  victim_case_form2_activation_channel: {
    no: "Notificación (tarea)", ll: "Llamada (tarea)", ms: "Mi Salvia (Gestión institucional externa) (tarea)",
  },
  victim_case_form2_barrier_health: {
    "01": "Dificultades de aseguramiento o afiliación en salud", "02": "Negaciones en los servicios de urgencias",
    "03": "Fallas en la calidad del servicio de urgencias", "04": "Demoras en la asignación de citas y continuidad de tratamientos",
    "05": "Faltas de activación de protocolos para violencia sexual", "06": "Faltas de activación de protocolos para ataques con agentes químicos",
    "07": "Demoras y dificultades en la valoración medicolegal", "08": "Negaciones en el acceso a la IVE",
    "09": "Exigencias indebidas de autorización (discapacidad)", "10": "Atenciones no centradas en cosmovisiones culturales propias",
    "11": "Negaciones o imposiciones de procedimientos a personas OSIGD", "12": "Otras barreras en salud",
  },
  victim_case_form2_barrier_justice: {
    "01": "Negativa institucional para recibir la denuncia", "02": "Dilación en la recepción y registro de la denuncia",
    "03": "Denegación al derecho a la no confrontación", "04": "Falta de entrega de información y documentación procesal",
    "05": "Formalismo excesivo en los trámites judiciales", "06": "Imposición de cargas probatorias excesivas",
    "07": "Falta de celeridad en la investigación judicial", "08": "Incumplimiento de órdenes judiciales",
    "09": "Tipificación errónea del delito", "10": "Traslado inadecuado del expediente judicial",
    "11": "Dificultades en la reasignación de procesos", "12": "Ausencia de acceso a representación judicial",
    "13": "Vacíos normativos en VBG", "14": "Riesgo de vencimiento de términos",
    "15": "Medidas de libertad que ponen en riesgo a la víctima", "16": "Falta de protección de información confidencial",
    "17": "Deficiencias en la notificación del agresor", "18": "Otras barreras en justicia",
  },
  victim_case_form2_barrier_protect: {
    "01": "Ausencia o demoras en medidas de protección en Comisaría de Familia", "02": "Ausencia o dilación en medidas de protección físicas",
    "03": "Demora en la activación de atención integral en Comisaría de Familia", "04": "Incumplimiento de medidas de protección",
    "05": "Divulgación de información confidencial", "06": "Medidas de protección ineficaces",
    "07": "Falta de seguimiento institucional a las medidas", "08": "Fallas en la valoración del riesgo",
    "09": "Medidas insuficientes", "10": "Omisión de restricción de visitas ante riesgo",
    "11": "Omisión de apoyo policial", "12": "Otras barreras en protección", "13": "Otras barreras en protección",
  },
  victim_case_form2_barrier_institutional: {
    "01": "Falta de aplicación del enfoque de género / Revictimización",
    "02": "Falta de aplicación o desactualización de protocolos de atención en VBG",
    "03": "Capacidad institucional limitada", "04": "Desarticulación institucional y ausencia de gestión integral de casos",
    "05": "Insuficiencia en la orientación en derechos a las víctimas",
  },
  victim_case_form2_barrier_economic: { "01": "Falta de autonomía económica", "02": "Pobreza y desigualdad estructural" },
  victim_case_form2_barrier_territorial: { "01": "Ausencia institucional, infraestructura y conectividad insuficiente" },
  victim_case_form2_barrier_differential: {
    "01": "Barreras OSIGD", "02": "Barreras salud mental", "03": "Barreras población étnica", "04": "Barreras ciclo de vida",
    "05": "Barreras discapacidad", "06": "Barreras personas migrantes y refugiadas", "07": "Barreras personas privadas de la libertad",
    "08": "Barreras situación de calle", "09": "Barreras víctimas de trata", "10": "Barreras actividades sexuales pagas",
    "11": "Barreras lideresas y defensoras", "12": "Barreras víctimas del conflicto armado",
  },
  victim_case_form2_barrier_management: {
    "01": "Orientación y enrutamiento - Llamada", "02": "Gestión administrativa - Llamada", "03": "Activación de ruta interinstitucional",
    "04": "Articulación institucional", "05": "Escalamiento a organismo de control", "06": "Alerta por barreras",
  },
  victim_case_form2_aggressors_occupation: {
    "01": "Delincuencia común", "02": "Persona desmovilizada", "03": "Miembro de grupos alzados al margen de la ley",
    "04": "Personal de seguridad privada", "05": "Personal de custodia", "06": "Servidor(a) público(a)",
    "07": "Líder(esa) religioso(a)/espiritual", "08": "Miembro Policía Nacional", "09": "Miembro Ejército Nacional",
    "10": "Docente/instructor(a)/coach", "11": "Personal institución de salud", "12": "Líder(esa) social",
    "13": "Afiliado(a) sindical", "14": "Dirigente sindical", "15": "Otro", "16": "No conoce/No informa",
  },
};

// Metadata de puntaje para los criterios de remisión con puntaje (GS § B).
// code -> puntos; el obligatorio ('ob') se valida aparte, no suma puntos.
var CRITERIOS_PUNTAJE_V3 = {
  victim_case_form2_psycho_criteria: { minimo: 3, puntos: { cs: 3, in: 2, sr: 1, ct: 1, sa: 1, nv: 1 } },
  victim_case_form2_stabiliz_criteria: { minimo: 1, puntos: { si: 1, vt: 1, pc: 1 } },
};

var EQUIPOS_REMISION_V3 = [
  { value: "me", label: "Medidas de Emergencia", criterioCategory: null, criterioFieldId: null },
  { value: "ap", label: "Atención Psicosocial", criterioCategory: "victim_case_form2_psycho_criteria", criterioFieldId: "GS_9" },
  { value: "es", label: "Estabilización", criterioCategory: "victim_case_form2_stabiliz_criteria", criterioFieldId: "GS_12" },
  { value: "ah", label: "Atención Hombres", criterioCategory: "victim_case_form2_men_criteria", criterioFieldId: "GS_10" },
  { value: "sd", label: "Salvia Dignidad", criterioCategory: null, criterioFieldId: null },
  { value: "di", label: "Discapacidad", criterioCategory: "victim_case_form2_disability_service", criterioFieldId: "GS_11" },
];

// GS_9..GS_12 (criterios por equipo) no están en CAMPOS_GS_A_V3 porque su
// visibilidad depende de qué equipos se marcaron en GS_6 (registro-v3-app.js
// los renderiza aparte, en la sección de remisión), pero sí necesitan su
// definición de enumCategory para resolver opciones/icode igual que
// cualquier otro campo "multiple".
var CAMPOS_GS_CRITERIOS_V3 = [
  { id: "GS_9", tipo: "multiple", label: "Criterios Atención Psicosocial", enumCategory: "victim_case_form2_psycho_criteria" },
  { id: "GS_10", tipo: "multiple", label: "Criterio Atención Hombres", enumCategory: "victim_case_form2_men_criteria" },
  { id: "GS_11", tipo: "multiple", label: "Servicios Discapacidad", enumCategory: "victim_case_form2_disability_service" },
  { id: "GS_12", tipo: "multiple", label: "Criterios Estabilización", enumCategory: "victim_case_form2_stabiliz_criteria" },
];

// ── CI — Contacto e identificación ──────────────────────────────────────
var CAMPOS_CI_V3 = [
  { id: "CI_1", tipo: "dropdown", label: "¿Cómo prefieres que nos refiramos a ti?", enumCategory: "victim_case_form2_pronoun", req: "C1" },
  { id: "CI_2", tipo: "text", label: "¿Cuáles son sus nombres y apellidos?", req: "C1", soloLetras: true },
  { id: "CI_3", tipo: "text", label: "¿Con qué nombre se identifica o prefiere que le llamemos? (nombre identitario)", req: "C1", soloLetras: true },
  { id: "CI_4", tipo: "text", label: "¿Cuál es su número de teléfono de contacto en Colombia?", req: "C1" },
  { id: "CI_5", tipo: "dropdown", label: "¿Qué tipo de documento de identificación tiene?", enumCategory: "victim_case_form2_victim_doc_type", req: "C1" },
  // D13: "Persona sin identificación" no pide número (secciones.md) --
  // antes CI_6 era incondicional, así que una víctima sin documento no
  // podía completar nunca el módulo CI (A3 bloqueaba el guardado esperando
  // un número que no existe).
  { id: "CI_6", tipo: "text", label: "¿Cuál es su número del documento de identificación?", req: "C1", condicional: { campo: "__tieneNumeroDocumento", valor: true } },
  { id: "CI_7", tipo: "boolean", label: "¿Se reconoce como una persona con discapacidad?", req: "C1" },
  { id: "CI_8", tipo: "multiple", label: "¿Requiere de algún ajuste razonable en el marco de la atención a su caso de VBG y/o VpP?", req: "C1", condicional: { campo: "CI_7", valor: true }, enumCategory: "victim_case_form2_adjustments_gbv" },
  { id: "CI_9", tipo: "boolean", label: "¿Requiere intérprete de idiomas y/o traducción en otra lengua distinta al español?", req: "C1" },
];

// ── DP — Datos personales (incluye residencia) ───────────────────────────
var CAMPOS_DP_V3 = [
  { id: "DP_1", tipo: "date", label: "¿Cuál es su fecha de nacimiento?", req: "C1" },
  { id: "DP_2", tipo: "boolean", label: "¿Tiene dificultades físicas, mentales y/o sensoriales en su vida diaria?", req: "C1" },
  { id: "DP_3", tipo: "multiple", label: "¿Qué actividades no puede o presenta dificultades para realizarlas?", req: "C2", condicional: { campo: "DP_2", valor: true }, enumCategory: "victim_case_form2_activities_unable_to_perform" },
  { id: "DP_4", tipo: "multiple", label: "Según la Ley 1996 de 2019 cuenta usted con", req: "D", condicional: { campo: "CI_7", valor: true }, enumCategory: "victim_case_form2_law_1996" },
  { id: "DP_5", tipo: "dropdown", label: "¿Cuál es su nacionalidad?", req: "C1", enumCategory: "victim_case_form2_nationality" },
  { id: "DP_50", tipo: "text", label: "Correo electrónico de la víctima (opcional)", req: "D" },
  { id: "DP_6", tipo: "dropdown", label: "¿De qué país es su nacionalidad?", req: "C2", condicionalEnumCode: { campo: "DP_5", enumCategory: "victim_case_form2_nationality", code: "ex" }, enumCategory: "victim_case_form2_specified_nationality" },
  { id: "DP_7", tipo: "dropdown", label: "¿Cuál es tu condición migratoria actualmente?", req: "D", enumCategory: "victim_case_form2_migration_condition" },
  { id: "DP_8", tipo: "dropdown", label: "¿Cuál es su identidad de género?", req: "C1", enumCategory: "victim_case_form2_gender_identity" },
  { id: "DP_9", tipo: "dropdown", label: "¿Cuál es su orientación sexual?", req: "C1", enumCategory: "victim_case_form2_sexual_orientation" },
  { id: "DP_10", tipo: "dropdown", label: "¿Cuál fue su sexo asignado al nacer?", req: "C2", enumCategory: "victim_case_form2_assigned_sex_at_birth" },
  { id: "DP_11", tipo: "multiple", label: "¿Pertenece a alguna población de especial protección?", req: "C1", enumCategory: "victim_case_form2_specially_protected_population" },
  { id: "DP_12", tipo: "dropdown", label: "¿Se reconoce como perteneciente a algún grupo étnico?", req: "C1", enumCategory: "victim_case_form2_ethnic_affiliation" },
  { id: "DP_13", tipo: "dropdown", label: "¿Se reconoce como perteneciente a un pueblo indígena?", req: "C1", condicionalEnumLabel: { campo: "DP_12", contains: "Indígena" }, enumCategory: "victim_case_form2_indigenous_people" },
  { id: "DP_14", tipo: "boolean", label: "¿Se reconoce como persona campesina?", req: "C1" },
  { id: "DP_15", tipo: "dropdown", label: "¿Cuál es su estado civil?", req: "C2", enumCategory: "victim_case_form2_marital_status" },
  { id: "DP_16", tipo: "dropdown", label: "¿Cuál es el último nivel de escolaridad que ha finalizado?", req: "D", enumCategory: "victim_case_form2_last_education_level" },
  // FIX (encontrado en vivo con Playwright, 2026-10-15): req era "D" acá
  // pero triadaFieldsV3 (victim_case_v3_service.go) YA exigía DP_17 para
  // triadaCompleta/Borrador->Iniciado desde antes de esta ronda. El agente
  // podía completar CI+DP+CA+HV+AG según la UI (DP_17 nunca aparecía en
  // "Faltan obligatorias") y el sidebar de TZ/GS se quedaba bloqueado sin
  // ninguna pista de por qué -- A3 (guardar sección) y A1 (bloqueo de
  // sidebar) usaban dos listas de obligatorias distintas para el mismo
  // campo. secciones.md documenta DP_17 como "Deseable", así que el
  // desacuerdo real está ahí -- se resuelve a favor del backend/triadaFieldsV3
  // (ya probado para Borrador->Iniciado) en vez de arriesgar un Caso
  // permanentemente bloqueado.
  { id: "DP_17", tipo: "dropdown", label: "¿Cuál es su ocupación?", req: "C1", enumCategory: "victim_case_form2_occupation" },
  { id: "DP_18", tipo: "dropdown", label: "¿Cuál ha sido su forma de generación de ingresos actualmente o en los últimos 3 meses?", req: "D", enumCategory: "victim_case_form2_income_generation_method" },
  // D14: le faltaba el condicional documentado en secciones.md -- "condicional
  // (Q18=contrato formal o independiente)" -- antes era visible siempre.
  { id: "DP_19", tipo: "dropdown", label: "¿Cuál es su tipo de vinculación laboral?", req: "D", condicionalEnumCodeIn: { campo: "DP_18", enumCategory: "victim_case_form2_income_generation_method", codes: ["em", "ta"] }, enumCategory: "victim_case_form2_employment_relationship" },
  // D12: ausentes hasta esta ronda -- secciones.md Q20/21/22, condicional
  // (Q17 o Q18 = Persona en ASP).
  { id: "DP_20", tipo: "multiple", label: "¿Cuál es la modalidad en la que desarrolla actividades sexuales pagas (ASP)?", req: "C1", condicional: { campo: "__personaEnASP", valor: true }, enumCategory: "victim_case_form2_asp_mode" },
  { id: "DP_21", tipo: "date", label: "¿Hace cuánto tiempo inició en las ASP aproximadamente?", req: "C1", condicional: { campo: "__personaEnASP", valor: true } },
  { id: "DP_22", tipo: "multiple", label: "¿Cuál fue el principal motivo por el cual inició en las ASP?", req: "C1", condicional: { campo: "__personaEnASP", valor: true }, enumCategory: "victim_case_form2_reason_asp" },
  { id: "DP_23", tipo: "dropdown", label: "¿Cuál es la forma de tenencia de la vivienda donde reside actualmente?", req: "D", enumCategory: "victim_case_form2_housing_tenancy_form" },
  { id: "DP_24", tipo: "dropdown", label: "¿Cuál es el estrato de la vivienda en la que reside actualmente?", req: "D", enumCategory: "victim_case_form2_housing_stratum" },
  { id: "DP_25", tipo: "multiple", label: "¿Tiene personas o animales a su cargo?", req: "D", enumCategory: "victim_case_form2_has_dependents" },
  { id: "DP_26", tipo: "boolean", label: "¿Actualmente se encuentra en estado de gestación (embarazo)?", req: "C2", condicional: { campo: "__gestacionAplica", valor: true } },
  { id: "DP_27", tipo: "dropdown", label: "¿En qué departamento reside actualmente?", opciones: DEPARTAMENTOS_V3, req: "C1" },
  { id: "DP_28", tipo: "text", label: "¿En qué municipio reside actualmente?", req: "C1" },
  { id: "DP_29", tipo: "text", label: "¿Cuál es la dirección de su lugar de residencia?", req: "C1" },
  { id: "DP_30", tipo: "dropdown", label: "¿En qué zona se encuentra su lugar de residencia?", req: "D", enumCategory: "victim_case_form2_residence_zone" },
];

// ── CA — Contacto de apoyo ────────────────────────────────────────────────
// req: "C2" en las 4 -- FIX: estaban como "D" (Deseable), pero
// registro-v3-secciones.md documenta un cambio de criterio explícito vs.
// V2-real: "ya no es opcional, solo tiene gracia hasta el segundo contacto".
// Quedaban como Deseable por error desde la implementación original.
var CAMPOS_CA_V3 = [
  { id: "CA_1", tipo: "text", label: "¿Podría indicarme el nombre de una persona de contacto de apoyo o emergencia en Colombia?", req: "C2", soloLetras: true },
  { id: "CA_2", tipo: "text", label: "¿Cuál es el número de teléfono de ese contacto de apoyo o emergencia?", req: "C2" },
  { id: "CA_3", tipo: "text", label: "¿Cuenta con un correo electrónico de ese contacto de apoyo o emergencia?", req: "C2" },
  { id: "CA_4", tipo: "dropdown", label: "¿Cuál es el parentesco o relación con esa persona de contacto?", req: "C2", enumCategory: "victim_case_form2_support_contact_kinship" },
];

// Tipo de violencia (HV_9, código 2-letras) -> sufijo de categoría de
// subtipo (HV_10). Los códigos coinciden 1:1 (ver migrations/add_registro_v3.sql).
var TIPO_A_SUFIJO_SUBTIPO_V3 = { fi: "fi", pl: "pl", po: "po", ps: "ps", re: "re", se: "se", vi: "vi" };

// ── HV — Hechos de violencia ──────────────────────────────────────────────
var CAMPOS_HV_V3 = [
  { id: "HV_1", tipo: "textarea", label: "¿Podría describir brevemente los hechos ocurridos?", req: "C1" },
  { id: "HV_2", tipo: "date", label: "¿En qué fecha ocurrió el hecho?", req: "C1" },
  { id: "HV_3", tipo: "text", label: "¿A qué hora ocurrió el hecho?", req: "D", placeholder: "Ej: 19:30" },
  { id: "HV_4", tipo: "dropdown", label: "¿En qué departamento ocurrieron los hechos?", opciones: DEPARTAMENTOS_V3, req: "C1" },
  { id: "HV_5", tipo: "text", label: "¿En qué municipio ocurrieron los hechos?", req: "C1" },
  { id: "HV_6", tipo: "dropdown", label: "¿En qué zona ocurrió el hecho?", req: "C1", enumCategory: "victim_case_form2_facts_zone" },
  { id: "HV_7", tipo: "text", label: "¿Cuál es la dirección donde ocurrieron los hechos (o nombre de la finca si es rural)?", req: "D" },
  { id: "HV_8", tipo: "dropdown", label: "¿En qué escenario ocurrieron los hechos?", req: "C1", enumCategory: "victim_case_form2_scenario_violence" },
  { id: "HV_9", tipo: "multiple", label: "¿Qué tipo o tipos de violencia ha experimentado?", req: "C1", enumCategory: "victim_case_form2_type_violence_experienced" },
  { id: "HV_10", tipo: "multiple", label: "¿Qué subtipo o subtipos de violencia se presentaron?", req: "C2", subtipoDe: "HV_9" },
  { id: "HV_11", tipo: "multiple", label: "¿En qué ámbito o ámbitos se presentó la violencia?", req: "C1", max: 2, enumCategory: "victim_case_form2_scope_of_violence" },
  { id: "HV_12", tipo: "dropdown", label: "¿En qué sector se presentó la violencia en el ámbito laboral?", req: "C1", condicional: { campo: "__ambitoLaboral", valor: true }, enumCategory: "victim_case_form2_workplace_sector_occurrence" },
  { id: "HV_13", tipo: "boolean", label: "¿Considera que la violencia se motivó en razón de su género, identidad de género y/o orientación sexual?", req: "C1" },
  { id: "HV_14", tipo: "dropdown", label: "Recurrencia de la agresión (¿cuántas veces ha ocurrido?)", req: "C1", enumCategory: "victim_case_form2_recurrence_aggression" },
];

// ── AG — Características del agresor ─────────────────────────────────────
var CAMPOS_AG_V3 = [
  { id: "AG_1", tipo: "dropdown", label: "¿Cuántas personas agresoras participaron en los hechos?", req: "C1", enumCategory: "victim_case_form2_num_agressors" },
  { id: "AG_2", tipo: "dropdown", label: "¿Conoce y convive con la principal persona agresora?", req: "C1", enumCategory: "victim_case_form2_proximity_principal_aggressor" },
  { id: "AG_3", tipo: "dropdown", label: "¿Qué relación tiene o tenía con la persona agresora?", req: "C1", enumCategories: ["victim_case_form2_relationship_with_presumed_aggressor_01", "victim_case_form2_relationship_with_presumed_aggressor_02"] },
  { id: "AG_4", tipo: "dropdown", label: "¿Qué ocupación tiene la persona agresora?", req: "D", enumCategory: "victim_case_form2_aggressors_occupation" },
  { id: "AG_5", tipo: "boolean", label: "¿Depende económicamente del agresor?", req: "C1", condicional: { campo: "__agresorConocido", valor: true } },
  { id: "AG_6", tipo: "dropdown", label: "Identidad de género del agresor", req: "D", enumCategory: "victim_case_form2_aggressor_gender_identity" },
  { id: "AG_7", tipo: "text", label: "Nombres del agresor (si se conoce)", req: "D", soloLetras: true },
  { id: "AG_8", tipo: "dropdown", label: "Tipo de documento de identificación del agresor", req: "D", enumCategory: "victim_case_form2_victim_doc_type" },
  { id: "AG_9", tipo: "text", label: "Número de documento de identificación del agresor", req: "D" },
  { id: "AG_10", tipo: "text", label: "Dirección del agresor (si se conoce)", req: "D" },
  { id: "AG_11", tipo: "text", label: "Teléfono del agresor (si se conoce)", req: "D" },
];

// ── TZ — Tamizaje ──────────────────────────────────────────────────────────
// D4: 5 obligatorias iniciales (propuesta B4), el resto se revela solo si
// alguna de esas 5 = Sí.
var TAMIZAJE_GENERALES_V3 = [
  { id: "TZ_1", label: "¿Ha aumentado la violencia física en severidad o frecuencia en los últimos 60 días?", inicial: true },
  { id: "TZ_2", label: "¿Le ha agredido o amenazado utilizando algún arma u objeto en su contra?", inicial: true },
  { id: "TZ_3", label: "¿Ha amenazado con matarle?", inicial: true },
  { id: "TZ_4", label: "¿Ha intentado asfixiarle, estrangularle o impedirle respirar?", inicial: true },
  { id: "TZ_5", label: "¿Esta persona le persigue, espía o llega a lugares donde no le está esperando?", inicial: false },
  { id: "TZ_6", label: "¿Le deja notas o mensajes amenazantes en diferentes medios?", inicial: false },
  { id: "TZ_7", label: "¿La persona agresora tiene acceso a armas?", inicial: true },
];
var TAMIZAJE_VPP_V3 = { id: "TZ_8", label: "¿Esta persona ha realizado amenazas, agresiones o comentarios motivados por prejuicio hacia su orientación sexual, identidad o expresión de género?" };

var TAMIZAJE_PAREJA_V3 = [
  "¿Está su pareja/expareja desempleado(a) o sin trabajo actualmente?",
  "¿Sabe si ha tenido algún proceso con la justicia?",
  "¿Sabe si ha tenido denuncias/procesos por violencia intrafamiliar o incumplimiento de medidas de protección?",
  "¿Consume drogas (anfetamina, cocaína, heroína u otras)?",
  "¿Le persuade u obliga a consumirlas?",
  "¿Es alcohólica o tiene problemas con el alcohol?",
  "¿Es celoso/a de manera constante o violenta y/o controla sus actividades, decisiones, desplazamientos o lugares?",
  "¿Limita, impide o amenaza su contacto con familiares, amistades, redes de apoyo o instituciones?",
  "¿La ha forzado a mantener relaciones sexuales sin que usted quiera?",
  "¿Le vigila, persigue, acosa, contacta insistentemente o utiliza terceros para intimidar/amenazar?",
  "¿La acosa sexualmente o remite mensajes de índole sexual no pedidos e incómodos?",
  "¿La ha amenazado con suicidarse o lo ha intentado?",
  "¿Ha difundido, amenazado con difundir o utilizado fotos/videos/mensajes/redes/información personal para exponerle o extorsionarle?",
  "¿Le ha impedido acceder a servicios de salud o salud mental?",
  "¿Utiliza su estado de salud o diagnósticos para chantajearla o amenazar con revelarlos?",
  "¿La ha agredido físicamente mientras estaba en enfermedad, recuperación médica, posoperatorio, embarazo u otra vulnerabilidad?",
  "¿La ha amenazado con quitarle, restringir comunicación/visitas o hacer daño a hijos, hijas, mascotas u otros familiares?",
  "¿Ha destruido o retenido objetos personales, documentos, medicamentos u otras pertenencias como forma de intimidación?",
  "¿Le ha expuesto intencionalmente a ITS o le ha impedido tomar medidas de protección?",
  "¿Actualmente convive con su pareja/expareja o comparte espacios frecuentes?",
  "¿Ha tenido que cambiar rutinas, rutas, residencia o espacios frecuentes por temor a agresiones?",
  "¿Ha amenazado con revelar su orientación sexual/identidad/expresión de género sin su consentimiento?",
  "¿Ha dejado de buscar ayuda institucional tras una agresión por temor a discriminación por su orientación sexual/identidad de género?",
  "Ideación suicida — ¿como consecuencia de la violencia ha tenido pensamientos para terminar con su vida?",
];
var TAMIZAJE_PAREJA_OSIGD_V3 = [22, 23]; // índices 1-based condicionales a esOsigd

var TAMIZAJE_NO_PAREJA_V3 = [
  "¿La persona agresora está desempleada o sin trabajo actualmente?",
  "¿Sabe si tiene antecedentes penales, denuncias por otros delitos o procesos judiciales?",
  "¿Consume drogas (anfetamina, cocaína, heroína u otras)?",
  "¿Le persuade u obliga a consumirlas?",
  "¿Es alcohólica o tiene problemas con el alcohol?",
  "¿Es celoso/a de manera constante o violenta y/o controla sus actividades, decisiones, desplazamientos o lugares?",
  "¿Limita, impide o amenaza su contacto con familiares, amistades, redes de apoyo o instituciones?",
  "¿La ha forzado a mantener relaciones sexuales sin que usted quiera?",
  "¿Le vigila, persigue, acosa, contacta insistentemente o utiliza terceros para intimidar/amenazar?",
  "¿La acosa sexualmente o remite mensajes de índole sexual no pedidos e incómodos?",
  "¿Ha difundido, amenazado con difundir o utilizado fotos/videos/mensajes/redes/información personal para exponerle o extorsionarle?",
  "¿Le ha impedido acceder a servicios de salud o salud mental?",
  "¿Ha aprovechado su estado de vulnerabilidad física o mental para agredirla/hostigarla?",
  "¿La ha amenazado con quitarle, restringir comunicación/visitas o hacer daño a hijos, hijas, mascotas u otros familiares?",
  "¿Ha destruido o retenido objetos personales, documentos, medicamentos u otras pertenencias como forma de intimidación?",
  "¿Le ha expuesto intencionalmente a ITS o le ha impedido tomar medidas de protección?",
  "¿Comparte espacios frecuentes con esta persona (barrio, trabajo, estudio, transporte, residencia)?",
  "¿Ha tenido que cambiar rutinas, rutas, residencia o espacios frecuentes por temor a agresiones?",
  "¿Ha amenazado con revelar su orientación sexual/identidad/expresión de género sin su consentimiento?",
  "¿La violencia es ejercida por miembros de su familia (expulsión, amenazas, control económico, daño a pertenencias)?",
  "¿Ha dejado de buscar ayuda institucional tras una agresión por temor a discriminación?",
  "¿Existe alguna relación de jerarquía entre usted y la persona agresora?",
  "¿Ha utilizado su posición de autoridad para amenazarla con pérdida de empleo, vivienda o servicios?",
  "Ideación suicida — ¿como consecuencia de la violencia ha tenido pensamientos para terminar con su vida?",
];
var TAMIZAJE_NOPAREJA_OSIGD_V3 = [18, 19, 20]; // índices 1-based condicionales a esOsigd
var TAMIZAJE_NOPAREJA_LABORAL_V3 = [21, 22]; // índices 1-based condicionales a ámbito laboral

// ── GS — Gestión Salvia ───────────────────────────────────────────────────
var CAMPOS_GS_A_V3 = [
  { id: "GS_1", tipo: "boolean", label: "¿Usted permitiría que el Sistema SALVIA realice el proceso de denuncia fácil ante la FGN?", req: "C1" },
  { id: "GS_2", tipo: "multiple", label: "¿Es necesario realizar gestiones interinstitucionales a partir de los hallazgos de esta atención?", req: "C1", enumCategory: "victim_case_form2_interinst_action" },
  { id: "GS_3", tipo: "textarea", label: "En caso afirmativo, describa las acciones previstas, las entidades involucradas y su propósito", req: "D" },
  { id: "GS_4", tipo: "textarea", label: "Describa el acompañamiento realizado y las acciones que se desarrollarán para la gestión del caso", req: "C1" },
];

// GS § B — Remisión a otros equipos. GS_6/GS_7 se renderizan en una
// sección propia del Formulario (con criterios de cada equipo), no en el
// listado genérico de CAMPOS_GS_A_V3, pero necesitan su propia definición
// para que opcionesPara/etiqueta (registro-v3-app.js) resuelvan sus icode.
var CAMPOS_GS_B_V3 = [
  { id: "GS_6", tipo: "multiple", label: "¿A qué equipo(s) de Salvia se remite el caso?", enumCategory: "victim_case_form2_referral_team" },
  { id: "GS_7", tipo: "multiple", label: "¿Cuáles medidas de emergencia?", enumCategory: "victim_case_form2_emergency_measure" },
];

// Vocabulario real de entity_branch/barrier_v2 -- igual que
// CE_SECTORES_MODAL en case-entities.js. NO es el enum victim_case_form2_sector
// (sa/ju/pr/ot): ver entity_branch_api_controller.go#ListByTownCode, que
// traduce estas mismas 4 palabras a su propio código de 2 letras interno
// (he/js/pt/os). Compartido por el buscador de Entidades y de Barreras,
// porque ambos resuelven la misma sede real (entity_branch).
var SECTORES_ENTIDAD_BARRERA_V3 = [
  { value: "salud", label: "Salud" }, { value: "justicia", label: "Justicia" },
  { value: "proteccion", label: "Protección" }, { value: "otras_instituciones", label: "Otras instituciones" },
];
var CANAL_ACTIVACION_V3 = [
  { value: "no", label: "Notificación (tarea)" }, { value: "ll", label: "Llamada (tarea)" },
  { value: "ms", label: "Mi Salvia (Gestión institucional externa) (tarea)" },
];

// Sector (vocabulario real, ver SECTORES_ENTIDAD_BARRERA_V3) -> categoría de
// barreras específicas. Solo salud/justicia/protección tienen catálogo
// propio; "otras instituciones" usa directamente las transversales
// (institutional/economic/territorial/differential), sin lista "specific".
var SECTOR_A_CATEGORIA_BARRERA_V3 = {
  salud: "victim_case_form2_barrier_health", justicia: "victim_case_form2_barrier_justice", proteccion: "victim_case_form2_barrier_protect",
};
var CATEGORIAS_BARRERA_TRANSVERSAL_V3 = [
  { campo: "institutional", enumCategory: "victim_case_form2_barrier_institutional", label: "Barreras institucionales y de talento humano" },
  { campo: "economic", enumCategory: "victim_case_form2_barrier_economic", label: "Barreras económicas" },
  { campo: "territorial", enumCategory: "victim_case_form2_barrier_territorial", label: "Barreras territoriales" },
  { campo: "differential", enumCategory: "victim_case_form2_barrier_differential", label: "Barreras por ausencia de enfoque diferencial" },
];
var GESTION_BARRERA_ENUM_CATEGORY_V3 = "victim_case_form2_barrier_management";
