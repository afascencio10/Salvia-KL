# Tablas de Base de Datos

Índice de todas las tablas de la BD. Schemas: `salvia` (dominio principal) y `security` (usuarios y geografía).

---

## Conexión a la base de datos

**Proveedor:** Supabase (PostgreSQL)  
**Config file:** `src/config/db_config.json`

### Pooler (aplicación + psql general)

| Campo | Valor |
|---|---|
| Host | `aws-1-us-west-1.pooler.supabase.com` |
| Puerto | `5432` (transaction mode) · `6543` (session mode) |
| Base de datos | `postgres` |
| SSL | `require` |

> **Puerto 5432 vs 6543:** usar `5432` para queries normales desde la aplicación. Usar `6543` (session mode) cuando se necesiten DDL (`CREATE TABLE`, `ALTER TABLE`) desde psql — el transaction mode del puerto 5432 no garantiza persistencia de DDL.

### Usuarios

| Usuario | Password | Uso | Permisos en `salvia` |
|---|---|---|---|
| `salvia_legacy.pwwelwfhauspznpatuqm` | `LegacySalvia2026@` | Consultas legacy / seed | SELECT, INSERT, UPDATE, DELETE + CREATE en schema |
| `salvia_gorm.pwwelwfhauspznpatuqm` | `GormSalvia2026@` | Backend Go (GORM) + AutoMigrate | Owner de todas las tablas + CREATE en schema |
| `postgres.pwwelwfhauspznpatuqm` | `Salvia2026@` | Superusuario — solo para operaciones administrativas | Superusuario |

> **AutoMigrate:** `salvia_gorm` es owner de todas las tablas del schema `salvia`. GORM AutoMigrate puede crear tablas nuevas y agregar columnas sin intervención manual. Miembro de `salvia_gorm`: también `postgres` (para poder transferir ownership).

### Conectarse con psql

```bash
# Usuario legacy (consultas, seed)
PGPASSWORD='LegacySalvia2026@' psql \
  -h aws-1-us-west-1.pooler.supabase.com \
  -U "salvia_legacy.pwwelwfhauspznpatuqm" \
  -d postgres

# Usuario GORM (queries + DDL sobre tablas que posee)
PGPASSWORD='GormSalvia2026@' psql \
  -h aws-1-us-west-1.pooler.supabase.com \
  -p 6543 \
  -U "salvia_gorm.pwwelwfhauspznpatuqm" \
  -d postgres

# Superusuario (operaciones administrativas, GRANT, cambio de ownership)
PGPASSWORD='Salvia2026@' psql \
  -h aws-1-us-west-1.pooler.supabase.com \
  -p 6543 \
  -U "postgres.pwwelwfhauspznpatuqm" \
  -d postgres
```

---

## Schema: `salvia`

### Casos y víctimas

| Tabla | Descripción |
|---|---|
| `victim_contact` | Datos de contacto de la víctima: identidad, documento, dirección, teléfono y género. |
| `victim_case` | Caso de violencia basado en género asociado a una víctima contacto. |
| `victim_case_form1` | Formulario 1 del caso: datos completos de víctima, hechos y agresor (legacy). |
| `victim_case_form2` | Formulario 2 del caso: datos ampliados con score y nivel de riesgo calculado. |
| `victim_contact_form1` | Formulario 1 del contacto: datos iniciales de la víctima reportados por un tercero. |
| `victim_contact_form2` | Formulario de reporte de tercero: datos del reportante, tipo de reporte y mejor horario. |
| `victim_case_form2_enums` | Catálogo de opciones compartidas por formularios del caso, contacto y feminicidio. |

### Organización institucional

| Tabla | Descripción |
|---|---|
| `entity` | Entidad institucional que atiende casos VBG: sector, interoperabilidad y tiempos de respuesta. |
| `entity_branch` | Sede o sucursal de una entidad con dirección, municipio y coordenadas geográficas. |
| `attention_line` | Línea de atención asociada a un profesional responsable de casos. |
| `case_owner` | Profesional responsable de casos, vinculado a una entidad y línea de atención. |

### Seguimiento

| Tabla | Descripción |
|---|---|
| `follow_up` | Evaluación de riesgo del proceso: factores de peligrosidad del agresor (legacy). |
| `follow_up_entry` | Entrada de seguimiento: visita o contacto con la víctima o una entidad (legacy). |
| `follow_up_entry_acting` | Actuación registrada para atender una barrera en una entrada de seguimiento (legacy). |
| `follow_up_v2` | Seguimiento v2: cita programada con la víctima, estado, equipo y número de secuencia. |
| `follow_up_attempts` | Intentos de contacto de un seguimiento v2: fecha, resultado y notas del agente. |

### Barreras institucionales

| Tabla | Descripción |
|---|---|
| `sector_barrier` | Sector institucional que agrupa barreras (justicia, salud, protección, etc.). |
| `barrier` | Catálogo maestro de barreras institucionales agrupadas por sector. |
| `barrier_v2` | Barrera identificada en un seguimiento: sector, descripción y estado actual. |
| `barrier_update` | Gestión de una barrera: oficio enviado, URL Kofax y respuesta recibida. |
| `rel_barrier_follow_up_entry` | Relación entre una barrera del catálogo y una entrada de seguimiento (legacy). |

### Formularios dinámicos

| Tabla | Descripción |
|---|---|
| `form` | Definición de un formulario dinámico: estructura base y metadatos de configuración. |
| `form_section` | Sección dentro de un formulario dinámico con orden y título. |
| `question` | Pregunta de una sección: tipo, texto, validaciones y orden de presentación. |
| `option` | Opción de respuesta para preguntas de selección simple o múltiple. |
| `visibility_condition` | Condición que controla si una pregunta/sección/repeater se muestra según respuestas previas o estado externo (`formState`). |
| `render_modification` | Modificación declarativa de textos visibles (labels, títulos, descripciones) según el estado externo (`formState`). Soporta `SET` (sobreescritura) y `REPLACE` (buscar y reemplazar). |
| `repeater_group` | Grupo de preguntas repetibles dentro de un formulario dinámico. |
| `repeater_entry` | Fila o instancia concreta de un repeater_group en un form_submission. |
| `form_submission` | Envío de un formulario: vincula respuestas a un caso, agente y seguimiento. |
| `answer` | Respuesta individual a una pregunta dentro de un form_submission. |

### Gestión del caso

| Tabla | Descripción |
|---|---|
| `moment` | Momento institucional del flujo: aprobación o cancelación de un paso del caso. |
| `case_log` | Bitácora de eventos importantes del ciclo de vida de un caso (legacy). |
| `case_timeline_event` | Evento del timeline del caso visible en la UI del historial del agente. |
| `case_task` | Tarea asignada dentro de un caso: tipo, estado, usuario y recurso relacionado. |
| `entity_letter` | Oficio enviado a una entidad para gestionar una barrera: estado y número de radicado. |

### Intervenciones

| Tabla | Descripción |
|---|---|
| `emergency_measure` | Medida de emergencia adoptada para proteger a la víctima: autoridad, vigencia y estado. |
| `psychosocial_support` | Apoyo psicosocial brindado a la víctima: tipo, institución, fechas y estado. |
| `economic_stabilization` | Proceso de estabilización económica de la víctima: beneficio, institución y estado. |

### Alertas

| Tabla | Descripción |
|---|---|
| `alert` | Alerta del sistema: tipo, prioridad y código identificador. |
| `rel_alert_victim_case` | Relación entre una alerta y un caso de víctima con datos adicionales. |
| `rel_alert_follow_up` | Relación entre una alerta y una entrada de seguimiento (legacy). |

### Feminicidio

| Tabla | Descripción |
|---|---|
| `feminicide` | Caso de feminicidio: identidad, documento y estado del registro. |
| `feminicide_risk` | Caso de riesgo de feminicidio: identidad y documento de la víctima en riesgo. |
| `feminicide_form1` | Formulario 1 de feminicidio: datos completos de la víctima y la familia afectada. |
| `feminicide_risk_form1` | Formulario 1 de riesgo de feminicidio: datos de la víctima en situación de riesgo. |

### Tablas relacionales

| Tabla | Descripción |
|---|---|
| `rel_entity_moment` | Relación entre una entidad y un momento del flujo del caso. |
| `rel_case_owner_victim_case` | Asignación de un profesional responsable a un caso (legacy). |
| `rel_case_owner_feminicide` | Asignación de un profesional a un caso de feminicidio. |
| `rel_case_owner_feminicide_risk` | Asignación de un profesional a un caso de riesgo de feminicidio. |
| `rel_victim_case_form2_enums_victim_case_form2` | Relación muchos-a-muchos entre enums y formulario 2 del caso. |
| `rel_victim_case_form2_enums_victim_contact_form2` | Relación entre enums y formulario 2 del contacto víctima. |
| `rel_victim_case_form2_enums_feminicide_form1` | Relación entre enums y formulario 1 de feminicidio. |
| `rel_victim_case_form2_enums_feminicide_risk_form1` | Relación entre enums y formulario 1 de riesgo de feminicidio. |

---

## Schema: `security`

### Usuarios y autenticación

| Tabla | Descripción |
|---|---|
| `general_user` | Usuario del sistema: login, contraseña hasheada y estado de autenticación. |
| `general_user_profile` | Perfil del usuario: nombre, documento, género y municipio de residencia. |
| `reset_password` | Solicitudes activas de restablecimiento de contraseña por usuario. |
| `role` | Roles del sistema: código, nombre y descripción de permisos. |
| `rel_role_general_user` | Asignación de roles a usuarios del sistema. |
| `email` | Correos electrónicos asociados al perfil de un usuario. |
| `phone_number` | Teléfonos asociados al perfil de un usuario. |

### Geografía

| Tabla | Descripción |
|---|---|
| `country` | Catálogo de países con código y nombre. |
| `department` | Catálogo de departamentos con código y nombre. |
| `city` | Catálogo de ciudades agrupadas por departamento. |
| `town` | Catálogo de municipios con código, nombre, tipo y coordenadas geográficas. |
