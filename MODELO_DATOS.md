# Modelo de Datos - Sistema SALVIA
## Base de Datos PostgreSQL

---

## ESQUEMA: salvia (Gestión de Casos)

### Diagrama Entidad-Relación Principal

```
┌──────────────────────────────────────────────────────────────────┐
│                    MÓDULO DE CASOS DE VÍCTIMAS                   │
└──────────────────────────────────────────────────────────────────┘

┌──────────────────┐
│ victim_contact   │  Contacto inicial de víctima
│ ──────────────── │
│ PK contact_id    │
│    i_code        │
│    names         │
│    last_names    │
│    phone         │
│    doc_type      │
│    doc_number    │
│    birth_date    │
│    address       │
│    latitude      │
│    longitude     │
│    status        │  (A=Activo, I=Inactivo, P=Procesado)
│    facts_desc    │
│    general_user  │
└────────┬─────────┘
         │ 1
         │
         │ 0..1
         │
┌────────▼─────────────────────────────────────────────────┐
│ victim_case                                              │  Caso completo
│ ──────────────────────────────────────────────────────── │
│ PK case_id                                               │
│    i_code                                                │
│    status                (A=Activo, C=Cerrado, etc.)     │
│    general_user                                          │
│                                                          │
│ DATOS VÍCTIMA:                                           │
│    victim_names, victim_last_names                       │
│    victim_doc_type, victim_doc_number                    │
│    victim_birth_date, victim_age                         │
│    victim_phone, victim_email                            │
│    victim_address, victim_town_code                      │
│    victim_latitude, victim_longitude                     │
│                                                          │
│ IDENTIDAD:                                               │
│    victim_gender                                         │
│    victim_gender_identity                                │
│    victim_sexual_orientation                             │
│    victim_nationality                                    │
│    victim_ethnic_group                                   │
│                                                          │
│ SITUACIÓN SOCIAL:                                        │
│    victim_occupation                                     │
│    victim_marital_status                                 │
│    victim_children_number                                │
│    victim_disability                                     │
│    victim_dependents                                     │
│                                                          │
│ HECHOS:                                                  │
│    facts_date, facts_start_time, facts_end_time          │
│    facts_description                                     │
│    facts_occurrence                                      │
│    violence_town_code                                    │
│    violence_experienced                                  │
│    violence_scope                                        │
│    violence_scene                                        │
│                                                          │
│ AGRESOR:                                                 │
│    aggressor (S/N)                                       │
│    aggressor_name                                        │
│    aggressor_doc_type, aggressor_doc_number              │
│    aggressor_address, aggressor_phone                    │
│    relationship_with_aggressor                           │
│    aggressor_has_weapons                                 │
│                                                          │
│ RIESGO:                                                  │
│    femicide_risk                                         │
│    imminent_risk                                         │
│    death_threats                                         │
│    experienced_violence_before                           │
│    previously_reported_situation                         │
│                                                          │
│ RELACIONES:                                              │
│ FK victim_contact_id                                     │
│ FK approved_by (case_owner_id)                           │
└────────┬─────────────────────────────────────────────────┘
         │
         │ 1
         │
         │ N
┌────────▼─────────┐
│     moment       │  Estados del proceso de atención
│ ──────────────── │
│ PK moment_id     │
│    i_code        │
│    code          │  (01=Recepción, 02=Valoración, etc.)
│    status        │  (A=Activo, C=Completado, X=Cancelado)
│    approval_     │
│    cancellation  │
│    approval_     │
│    source        │
│    approval_     │
│    description   │
│ FK victim_case_id│
│ FK approval_     │
│    owner_id      │
│ FK entity_       │
│    branch_id     │
└────────┬─────────┘
         │
         │ 1
         │
         │ N
┌────────▼─────────┐
│    case_log      │  Registro de actividades
│ ──────────────── │
│ PK log_id        │
│    i_code        │
│    creation_date │
│    description   │
│    general_user  │
│ FK moment_id     │
└──────────────────┘


┌──────────────────┐         ┌──────────────────┐
│   case_owner     │    N    │  attention_line  │
│ ──────────────── │◀────────│ ──────────────── │
│ PK owner_id      │    1    │ PK line_id       │
│    i_code        │         │    i_code        │
│    general_user  │         │    name          │
│    num_cases     │         │    description   │
│ FK attention_    │         └──────────────────┘
│    line_id       │
│ FK entity_       │
│    branch_id     │
└────────┬─────────┘
         │
         │ N
         │
         │ N
┌────────▼─────────────────────────────────┐
│ rel_case_owner_victim_case               │
│ ──────────────────────────────────────── │
│ PK rel_id                                │
│ FK case_owner_id                         │
│ FK victim_case_id                        │
│    creation_date                         │
└──────────────────────────────────────────┘


┌──────────────────┐         ┌──────────────────┐
│     entity       │    1    │  entity_branch   │
│ ──────────────── │◀────────│ ──────────────── │
│ PK entity_id     │    N    │ PK branch_id     │
│    i_code        │         │    i_code        │
│    name          │         │    name          │
│    description   │         │    description   │
│    sector        │         │    address       │
│    is_           │         │    town_code     │
│    interoperable │         │    latitude      │
│    response_time │         │    longitude     │
│    start_edu_    │         │ FK entity_id     │
│    content       │         └──────────────────┘
│    fail_edu_     │
│    content       │
└────────┬─────────┘
         │
         │ N
         │
         │ N
┌────────▼─────────┐
│ rel_entity_      │  Relación entidad-momento
│    moment        │
│ ──────────────── │
│ PK rel_id        │
│    moment_code   │  (01, 02, 03, etc.)
│ FK entity_id     │
└──────────────────┘


┌──────────────────┐         ┌──────────────────┐
│      alert       │    N    │   victim_case    │
│ ──────────────── │◀────────│                  │
│ PK alert_id      │    N    │                  │
│    i_code        │         │                  │
│    type          │         │                  │
│    priority      │         │                  │
│    code          │         │                  │
└────────┬─────────┘         └──────────────────┘
         │
         │ N
         │
         │ N
┌────────▼─────────────────────────────────┐
│ rel_alert_victim_case                    │
│ ──────────────────────────────────────── │
│ PK rel_id                                │
│ FK alert_id                              │
│ FK victim_case_id                        │
│    creation_date                         │
│    data                                  │
└──────────────────────────────────────────┘
```

---

## TABLAS PRINCIPALES

### 1. victim_contact (Contacto Inicial)

**Propósito:** Almacenar información de contactos iniciales de víctimas antes de crear un caso completo.

```sql
CREATE TABLE salvia.victim_contact (
    victim_contact_id SERIAL PRIMARY KEY,
    victim_contact_i_code VARCHAR(36) UNIQUE NOT NULL,
    victim_contact_creation_date TIMESTAMP NOT NULL,
    victim_contact_update_date TIMESTAMP NOT NULL,
    victim_contact_general_user VARCHAR(36) NOT NULL,
    
    -- Datos personales
    victim_contact_names VARCHAR(32),
    victim_contact_last_names VARCHAR(32),
    victim_contact_nick VARCHAR(32),
    victim_contact_doc_type VARCHAR(2),
    victim_contact_doc_number VARCHAR(32),
    victim_contact_birth_date TIMESTAMP,
    
    -- Ubicación
    victim_contact_living_zone VARCHAR(2),
    victim_contact_address VARCHAR(120),
    victim_contact_living_latitude DOUBLE PRECISION,
    victim_contact_living_longitude DOUBLE PRECISION,
    
    -- Contacto
    victim_contact_phone VARCHAR(10),
    
    -- Identidad
    victim_contact_gender_identity VARCHAR(2),
    victim_contact_sexual_orientation VARCHAR(2),
    victim_contact_origin VARCHAR(1),
    victim_contact_occupation VARCHAR(2),
    victim_contact_occupation_other VARCHAR(32),
    
    -- Estado
    victim_contact_status VARCHAR(1) NOT NULL,
    victim_contact_status_description VARCHAR(256),
    victim_contact_facts_description TEXT
);
```

**Estados:**
- `A` = Activo (pendiente de procesar)
- `P` = Procesado (ya tiene caso asociado)
- `I` = Invalidado (descartado)

**Índices:**
```sql
CREATE INDEX idx_victim_contact_status ON salvia.victim_contact(victim_contact_status);
CREATE INDEX idx_victim_contact_user ON salvia.victim_contact(victim_contact_general_user);
```

---

### 2. victim_case (Caso de Víctima)

**Propósito:** Almacenar información completa de casos de víctimas de violencia de género.

```sql
CREATE TABLE salvia.victim_case (
    victim_case_id SERIAL PRIMARY KEY,
    victim_case_i_code VARCHAR(36) UNIQUE NOT NULL,
    victim_case_creation_date TIMESTAMP NOT NULL,
    victim_case_update_date TIMESTAMP NOT NULL,
    victim_case_status VARCHAR(2) NOT NULL,
    victim_case_general_user VARCHAR(36) NOT NULL,
    
    -- Datos de la víctima (50+ campos)
    -- Ver script SQL completo para todos los campos
    
    -- Relaciones
    victim_case_victim_contact BIGINT REFERENCES victim_contact(victim_contact_id),
    victim_case_approved_by BIGINT REFERENCES case_owner(case_owner_id)
);
```

**Estados:**
- `A` = Activo
- `C` = Cerrado
- `S` = Suspendido
- `T` = Transferido

**Campos Clave:**
- **Datos Personales**: nombres, apellidos, documento, edad, teléfono, email
- **Ubicación**: dirección, municipio, coordenadas
- **Identidad**: género, orientación sexual, etnia, nacionalidad
- **Situación Social**: ocupación, estado civil, hijos, discapacidad
- **Hechos**: fecha, descripción, tipo de violencia, ámbito
- **Agresor**: nombre, documento, dirección, relación con víctima
- **Riesgo**: riesgo de feminicidio, amenazas, armas, violencia previa

---

### 3. moment (Momentos del Proceso)

**Propósito:** Registrar los diferentes estados/momentos por los que pasa un caso.

```sql
CREATE TABLE salvia.moment (
    moment_id SERIAL PRIMARY KEY,
    moment_i_code VARCHAR(36) NOT NULL,
    moment_creation_date TIMESTAMP NOT NULL,
    moment_update_date TIMESTAMP NOT NULL,
    moment_code VARCHAR(2) NOT NULL,
    moment_status VARCHAR(1) NOT NULL,
    moment_approval_cancellation VARCHAR(1) NOT NULL,
    moment_approval_source VARCHAR(1),
    moment_approval_description TEXT,
    
    -- Relaciones
    moment_victim_case BIGINT NOT NULL REFERENCES victim_case(victim_case_id),
    moment_approval_owner BIGINT NOT NULL REFERENCES case_owner(case_owner_id),
    moment_entity_branch BIGINT NOT NULL REFERENCES entity_branch(entity_branch_id)
);
```

**Códigos de Momento:**
- `01` = Recepción
- `02` = Valoración inicial
- `03` = Atención psicológica
- `04` = Atención jurídica
- `05` = Atención médica
- `06` = Protección
- `07` = Seguimiento
- `08` = Cierre

**Estados:**
- `A` = Activo
- `C` = Completado
- `X` = Cancelado

---

### 4. case_owner (Propietario de Caso)

**Propósito:** Representar operadores que gestionan casos.

```sql
CREATE TABLE salvia.case_owner (
    case_owner_id SERIAL PRIMARY KEY,
    case_owner_i_code VARCHAR(36) UNIQUE NOT NULL,
    case_owner_creation_date TIMESTAMP NOT NULL,
    case_owner_update_date TIMESTAMP NOT NULL,
    case_owner_general_user VARCHAR(32) NOT NULL,
    case_owner_num_cases INTEGER NOT NULL,
    
    -- Relaciones
    attention_line_id BIGINT REFERENCES attention_line(attention_line_id),
    entity_branch_id BIGINT REFERENCES entity_branch(entity_branch_id)
);
```

**Propósito:**
- Asignar casos a operadores específicos
- Controlar carga de trabajo (num_cases)
- Vincular con línea de atención y entidad

---

### 5. entity (Entidad de Atención)

**Propósito:** Representar instituciones que brindan atención.

```sql
CREATE TABLE salvia.entity (
    entity_id SERIAL PRIMARY KEY,
    entity_i_code VARCHAR(36) UNIQUE NOT NULL,
    entity_creation_date TIMESTAMP NOT NULL,
    entity_update_date TIMESTAMP NOT NULL,
    entity_name VARCHAR(254) NOT NULL,
    entity_description VARCHAR(512) NOT NULL,
    entity_sector VARCHAR(2) NOT NULL,
    entity_is_interoperable VARCHAR(1) NOT NULL,
    entity_interoperability_code VARCHAR(32),
    entity_response_time BIGINT NOT NULL,
    entity_start_edu_content TEXT NOT NULL,
    entity_fail_edu_content TEXT NOT NULL
);
```

**Sectores:**
- `PS` = Psicológico
- `JU` = Jurídico
- `ME` = Médico
- `PR` = Protección
- `SO` = Social

**Interoperabilidad:**
- `S` = Sí (puede recibir derivaciones automáticas)
- `N` = No (requiere gestión manual)

---

### 6. entity_branch (Sede de Entidad)

**Propósito:** Representar sedes físicas de entidades.

```sql
CREATE TABLE salvia.entity_branch (
    entity_branch_id SERIAL PRIMARY KEY,
    entity_branch_i_code VARCHAR(36) UNIQUE NOT NULL,
    entity_branch_creation_date TIMESTAMP NOT NULL,
    entity_branch_update_date TIMESTAMP NOT NULL,
    entity_branch_name VARCHAR(254) NOT NULL,
    entity_branch_description VARCHAR(512),
    entity_branch_address VARCHAR(254) NOT NULL,
    entity_branch_town_code VARCHAR(8) NOT NULL,
    entity_branch_latitude DOUBLE PRECISION,
    entity_branch_longitude DOUBLE PRECISION,
    
    -- Relación
    entity_id BIGINT REFERENCES entity(entity_id)
);
```

**Uso:**
- Geolocalización de servicios
- Asignación de casos por proximidad
- Mapas interactivos en frontend

---

### 7. alert (Alertas)

**Propósito:** Generar alertas automáticas basadas en criterios de riesgo.

```sql
CREATE TABLE salvia.alert (
    alert_id SERIAL PRIMARY KEY,
    alert_i_code VARCHAR(36) UNIQUE NOT NULL,
    alert_creation_date TIMESTAMP NOT NULL,
    alert_type VARCHAR(1) NOT NULL,
    alert_priority INTEGER NOT NULL,
    alert_code VARCHAR(64) NOT NULL
);
```

**Tipos:**
- `R` = Riesgo (alto riesgo de feminicidio)
- `T` = Tiempo (excede tiempo de respuesta)
- `S` = Sistema (error o inconsistencia)

**Prioridades:**
- `1` = Crítica (atención inmediata)
- `2` = Alta (atención en 24h)
- `3` = Media (atención en 72h)
- `4` = Baja (seguimiento normal)

---

### 8. case_log (Registro de Actividades)

**Propósito:** Auditoría y trazabilidad de acciones en casos.

```sql
CREATE TABLE salvia.case_log (
    case_log_id SERIAL PRIMARY KEY,
    case_log_i_code VARCHAR(36) UNIQUE NOT NULL,
    case_log_creation_date TIMESTAMP NOT NULL,
    case_log_description TEXT NOT NULL,
    case_log_general_user VARCHAR(36) NOT NULL,
    
    -- Relación
    moment_id BIGINT REFERENCES moment(moment_id)
);
```

**Ejemplos de Logs:**
- "Caso creado por operador Juan Pérez"
- "Momento 03 (Atención psicológica) iniciado"
- "Derivación a entidad XYZ aprobada"
- "Caso cerrado exitosamente"


---

## ESQUEMA: security (Autenticación y Autorización)

### Diagrama Entidad-Relación

```
┌──────────────────────────────────────────────────────────────────┐
│                    MÓDULO DE SEGURIDAD                           │
└──────────────────────────────────────────────────────────────────┘

┌──────────────────┐         ┌──────────────────┐
│ general_user     │    N    │      role        │
│ ──────────────── │◀───────▶│ ──────────────── │
│ PK user_id       │    N    │ PK role_id       │
│    i_code        │         │    code          │
│    login         │         │    name          │
│    password      │         │    description   │
│    status        │         └──────────────────┘
│    language      │                 ▲
│ FK profile_id    │                 │
└────────┬─────────┘                 │
         │ 1                         │
         │                           │
         │ 1                         │
┌────────▼─────────┐                 │
│ general_user_    │                 │
│    profile       │                 │
│ ──────────────── │                 │
│ PK profile_id    │                 │
│    i_code        │                 │
│    names         │                 │
│    last_names    │                 │
│    doc_type      │                 │
│    doc_number    │                 │
│    gender        │                 │
│    nick          │                 │
│    description   │                 │
│ FK town_code     │                 │
└────────┬─────────┘                 │
         │ 1                         │
         │                           │
         │ N                         │
┌────────▼─────────┐                 │
│      email       │                 │
│ ──────────────── │                 │
│ PK email_id      │                 │
│    i_code        │                 │
│    data          │                 │
│ FK profile_id    │                 │
└──────────────────┘                 │
                                     │
┌────────────────┐                   │
│  phone_number  │                   │
│ ────────────── │                   │
│ PK phone_id    │                   │
│    i_code      │                   │
│    data        │                   │
│ FK profile_id  │                   │
└────────────────┘                   │
                                     │
┌────────────────────────────────────┘
│
│  rel_role_general_user
│  (Tabla de relación N:M)
│
└─────────────────────────────────────┐
                                      │
                                      ▼
                            ┌──────────────────┐
                            │ Permisos por Rol │
                            │ (En código Go)   │
                            └──────────────────┘
```

### Tablas de Seguridad

#### 1. general_user (Usuario)

```sql
CREATE TABLE security.general_user (
    general_user_id SERIAL PRIMARY KEY,
    general_user_i_code VARCHAR(36) UNIQUE NOT NULL,
    general_user_creation_date TIMESTAMP NOT NULL,
    general_user_update_date TIMESTAMP NOT NULL,
    general_user_login VARCHAR(128) UNIQUE NOT NULL,
    general_user_password VARCHAR(512) NOT NULL,  -- Hash bcrypt
    general_user_status VARCHAR(10) NOT NULL,
    general_user_language VARCHAR(2) NOT NULL,
    
    -- Relación
    general_user_general_user_profile BIGINT REFERENCES general_user_profile(general_user_profile_id)
);
```

**Estados:**
- `A` = Activo
- `I` = Inactivo
- `B` = Bloqueado
- `P` = Pendiente de activación

**Idiomas:**
- `sp` = Español
- `en` = Inglés

#### 2. general_user_profile (Perfil de Usuario)

```sql
CREATE TABLE security.general_user_profile (
    general_user_profile_id SERIAL PRIMARY KEY,
    general_user_profile_i_code VARCHAR(36) UNIQUE NOT NULL,
    general_user_profile_creation_date TIMESTAMP NOT NULL,
    general_user_profile_update_date TIMESTAMP NOT NULL,
    general_user_profile_names VARCHAR(128) NOT NULL,
    general_user_profile_last_names VARCHAR(128),
    general_user_profile_gender VARCHAR(2),
    general_user_profile_nick VARCHAR(120),
    general_user_profile_description VARCHAR(1024),
    general_user_profile_doc_type VARCHAR(2),
    general_user_profile_doc_number VARCHAR(32),
    general_user_profile_town VARCHAR(8)
);
```

#### 3. role (Rol)

```sql
CREATE TABLE security.role (
    role_id SERIAL PRIMARY KEY,
    role_i_code VARCHAR(36) UNIQUE NOT NULL,
    role_creation_date TIMESTAMP NOT NULL,
    role_update_date TIMESTAMP NOT NULL,
    role_code VARCHAR(32) UNIQUE NOT NULL,
    role_name VARCHAR(128) NOT NULL,
    role_description VARCHAR(256) NOT NULL
);
```

**Roles del Sistema:**
- `ad` = Administrador (acceso total)
- `sv` = Supervisor (gestión de operadores y casos)
- `op` = Operador (gestión de casos asignados)

#### 4. rel_role_general_user (Relación Usuario-Rol)

```sql
CREATE TABLE security.rel_role_general_user (
    rel_role_general_user_id SERIAL PRIMARY KEY,
    rel_role_general_user_creation_date TIMESTAMP NOT NULL,
    role_id BIGINT NOT NULL REFERENCES role(role_id),
    general_user_id BIGINT NOT NULL REFERENCES general_user(general_user_id)
);
```

**Nota:** Un usuario puede tener múltiples roles.

---

## RELACIONES ENTRE ESQUEMAS

### Conexión security ↔ salvia

```
security.general_user
         │
         │ (Referencia por i_code, no FK)
         │
         ▼
salvia.victim_contact.victim_contact_general_user
salvia.victim_case.victim_case_general_user
salvia.case_owner.case_owner_general_user
salvia.case_log.case_log_general_user
```

**Nota:** La relación se hace por `i_code` (VARCHAR) en lugar de FK para permitir flexibilidad entre esquemas.

### Conexión security.town ↔ salvia

```
security.town
         │
         │ (Referencia por town_code)
         │
         ▼
salvia.victim_case.victim_case_victim_town_code
salvia.victim_case.victim_case_violence_town_code
salvia.entity_branch.entity_branch_town_code
```

---

## ÍNDICES Y OPTIMIZACIONES

### Índices Críticos

```sql
-- Security
CREATE INDEX idx_general_user_login ON security.general_user(general_user_login);
CREATE INDEX idx_general_user_i_code ON security.general_user(general_user_i_code);
CREATE INDEX idx_general_user_status ON security.general_user(general_user_status);

-- Salvia
CREATE INDEX idx_victim_contact_status ON salvia.victim_contact(victim_contact_status);
CREATE INDEX idx_victim_contact_user ON salvia.victim_contact(victim_contact_general_user);
CREATE INDEX idx_victim_case_status ON salvia.victim_case(victim_case_status);
CREATE INDEX idx_victim_case_user ON salvia.victim_case(victim_case_general_user);
CREATE INDEX idx_victim_case_town ON salvia.victim_case(victim_case_victim_town_code);
CREATE INDEX idx_moment_case ON salvia.moment(moment_victim_case);
CREATE INDEX idx_moment_status ON salvia.moment(moment_status);
CREATE INDEX idx_case_log_moment ON salvia.case_log(moment_id);
```

### Índices Geoespaciales (Opcional)

```sql
-- Para búsquedas por proximidad
CREATE INDEX idx_victim_case_location 
    ON salvia.victim_case 
    USING GIST (
        ll_to_earth(victim_case_victim_living_latitude, 
                    victim_case_victim_living_longitude)
    );

CREATE INDEX idx_entity_branch_location 
    ON salvia.entity_branch 
    USING GIST (
        ll_to_earth(entity_branch_latitude, 
                    entity_branch_longitude)
    );
```

---

## CONSTRAINTS Y VALIDACIONES

### Check Constraints

```sql
-- Validar estados
ALTER TABLE salvia.victim_contact 
    ADD CONSTRAINT chk_victim_contact_status 
    CHECK (victim_contact_status IN ('A', 'P', 'I'));

ALTER TABLE salvia.victim_case 
    ADD CONSTRAINT chk_victim_case_status 
    CHECK (victim_case_status IN ('A', 'C', 'S', 'T'));

ALTER TABLE salvia.moment 
    ADD CONSTRAINT chk_moment_status 
    CHECK (moment_status IN ('A', 'C', 'X'));

-- Validar rangos
ALTER TABLE salvia.victim_case 
    ADD CONSTRAINT chk_victim_case_age 
    CHECK (victim_case_age >= 0 AND victim_case_age <= 120);

ALTER TABLE salvia.alert 
    ADD CONSTRAINT chk_alert_priority 
    CHECK (alert_priority >= 1 AND alert_priority <= 4);

-- Validar coordenadas
ALTER TABLE salvia.victim_case 
    ADD CONSTRAINT chk_victim_case_latitude 
    CHECK (victim_case_victim_living_latitude >= -90 
           AND victim_case_victim_living_latitude <= 90);

ALTER TABLE salvia.victim_case 
    ADD CONSTRAINT chk_victim_case_longitude 
    CHECK (victim_case_victim_living_longitude >= -180 
           AND victim_case_victim_living_longitude <= 180);
```

### Triggers

```sql
-- Actualizar fecha de modificación automáticamente
CREATE OR REPLACE FUNCTION update_modified_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.victim_case_update_date = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_victim_case_modtime
    BEFORE UPDATE ON salvia.victim_case
    FOR EACH ROW
    EXECUTE FUNCTION update_modified_column();
```

---

## VISTAS ÚTILES

### Vista: Casos con Información Completa

```sql
CREATE VIEW salvia.v_victim_case_full AS
SELECT 
    vc.*,
    t.town_name,
    c.city_name,
    d.department_name,
    co.case_owner_general_user as assigned_to,
    COUNT(m.moment_id) as total_moments,
    COUNT(CASE WHEN m.moment_status = 'C' THEN 1 END) as completed_moments
FROM salvia.victim_case vc
LEFT JOIN security.town t ON vc.victim_case_victim_town_code = t.town_code
LEFT JOIN security.city c ON t.city_id = c.city_id
LEFT JOIN security.department d ON c.department_id = d.department_id
LEFT JOIN salvia.case_owner co ON vc.victim_case_approved_by = co.case_owner_id
LEFT JOIN salvia.moment m ON vc.victim_case_id = m.moment_victim_case
GROUP BY vc.victim_case_id, t.town_name, c.city_name, d.department_name, co.case_owner_general_user;
```

### Vista: Alertas Activas

```sql
CREATE VIEW salvia.v_active_alerts AS
SELECT 
    a.*,
    vc.victim_case_i_code,
    vc.victim_case_victim_names,
    vc.victim_case_victim_last_names,
    rac.rel_alert_victim_case_creation_date as alert_date
FROM salvia.alert a
INNER JOIN salvia.rel_alert_victim_case rac ON a.alert_id = rac.rel_alert_victim_case_alert_id
INNER JOIN salvia.victim_case vc ON rac.rel_alert_victim_case_victim_case_id = vc.victim_case_id
WHERE vc.victim_case_status = 'A'
ORDER BY a.alert_priority ASC, rac.rel_alert_victim_case_creation_date DESC;
```

### Vista: Carga de Trabajo por Operador

```sql
CREATE VIEW salvia.v_operator_workload AS
SELECT 
    co.case_owner_general_user,
    co.case_owner_num_cases,
    COUNT(DISTINCT rcvc.victim_case_id) as active_cases,
    COUNT(DISTINCT m.moment_id) as pending_moments
FROM salvia.case_owner co
LEFT JOIN salvia.rel_case_owner_victim_case rcvc ON co.case_owner_id = rcvc.case_owner_id
LEFT JOIN salvia.victim_case vc ON rcvc.victim_case_id = vc.victim_case_id AND vc.victim_case_status = 'A'
LEFT JOIN salvia.moment m ON vc.victim_case_id = m.moment_victim_case AND m.moment_status = 'A'
GROUP BY co.case_owner_id, co.case_owner_general_user, co.case_owner_num_cases;
```

---

## DATOS DE EJEMPLO

### Insertar Roles

```sql
INSERT INTO security.role (role_i_code, role_creation_date, role_update_date, role_code, role_name, role_description)
VALUES 
    (gen_random_uuid(), NOW(), NOW(), 'ad', 'Administrador', 'Acceso total al sistema'),
    (gen_random_uuid(), NOW(), NOW(), 'sv', 'Supervisor', 'Gestión de operadores y casos'),
    (gen_random_uuid(), NOW(), NOW(), 'op', 'Operador', 'Gestión de casos asignados');
```

### Insertar Usuario de Prueba

```sql
-- Crear perfil
INSERT INTO security.general_user_profile (
    general_user_profile_i_code,
    general_user_profile_creation_date,
    general_user_profile_update_date,
    general_user_profile_names,
    general_user_profile_last_names
) VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    'Admin',
    'Sistema'
) RETURNING general_user_profile_id;

-- Crear usuario (usar el ID del perfil anterior)
INSERT INTO security.general_user (
    general_user_i_code,
    general_user_creation_date,
    general_user_update_date,
    general_user_login,
    general_user_password,  -- Hash de "admin123"
    general_user_status,
    general_user_language,
    general_user_general_user_profile
) VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    'admin',
    '$2a$10$...',  -- Generar con bcrypt
    'A',
    'sp',
    1  -- ID del perfil
);
```

---

## BACKUP Y MANTENIMIENTO

### Script de Backup

```bash
# Backup completo
pg_dump -U salvia_admin -d salvia -F c -f salvia_backup_$(date +%Y%m%d).dump

# Backup solo esquema
pg_dump -U salvia_admin -d salvia -s -f salvia_schema_$(date +%Y%m%d).sql

# Backup solo datos
pg_dump -U salvia_admin -d salvia -a -f salvia_data_$(date +%Y%m%d).sql
```

### Script de Restauración

```bash
# Restaurar desde dump
pg_restore -U salvia_admin -d salvia salvia_backup_20240223.dump

# Restaurar desde SQL
psql -U salvia_admin -d salvia -f salvia_backup_20240223.sql
```

### Mantenimiento Regular

```sql
-- Vacuum y análisis
VACUUM ANALYZE salvia.victim_case;
VACUUM ANALYZE salvia.moment;
VACUUM ANALYZE salvia.case_log;

-- Reindexar
REINDEX TABLE salvia.victim_case;
REINDEX TABLE salvia.moment;

-- Estadísticas
ANALYZE salvia.victim_case;
```

---

**Documento creado:** 23 de febrero de 2026  
**Versión:** 1.0  
**Autor:** Equipo de Desarrollo SALVIA
