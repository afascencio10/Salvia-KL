# Gestión de Usuarios — SOG Salvia

Documento que describe el modelo de datos de seguridad y los pasos para crear y configurar usuarios en el sistema.

---

## Modelo de datos (`schema: security`)

### `general_user`
Tabla principal del usuario. Contiene credenciales, estado y equipo.

| Columna | Tipo | Descripción |
|---|---|---|
| `general_user_id` | uuid | PK |
| `general_user_i_code` | varchar | Código único legible (ej. `user-icode-op-test-00001`) |
| `general_user_login` | varchar | Login de acceso |
| `general_user_password` | varchar | Hash bcrypt (cost 10) |
| `general_user_status` | char | Estado: `e` = activo, `i` = inactivo |
| `general_user_language` | varchar | Idioma (ej. `sp`) |
| `general_user_general_user_profile` | uuid | FK → `general_user_profile.general_user_profile_id` |
| `general_user_team` | varchar | Nombre del equipo (ej. `Riesgo alto`, `Riesgo bajo`) |

---

### `general_user_profile`
Datos personales del usuario. Debe existir y estar vinculado en `general_user` para que el usuario aparezca en las consultas de la app (los JOINs son INNER).

| Columna | Tipo | Descripción |
|---|---|---|
| `general_user_profile_id` | uuid | PK |
| `general_user_profile_i_code` | varchar | Código único legible |
| `general_user_profile_names` | varchar | Nombres |
| `general_user_profile_last_names` | varchar | Apellidos |
| `general_user_profile_gender` | char | Género |
| `general_user_profile_nick` | varchar | Apodo o nombre corto |
| `general_user_profile_description` | text | Descripción libre |
| `general_user_profile_doc_type` | varchar | Tipo de documento |
| `general_user_profile_doc_number` | varchar | Número de documento |
| `general_user_profile_town` | varchar | Código DANE del municipio (ej. `11001000` = Bogotá) |

> **Importante:** Si `general_user_general_user_profile` es NULL, el usuario queda excluido de todos los listados de la app que usan INNER JOIN con el perfil.

---

### `role`
Catálogo de roles del sistema.

| `role_code` | `role_name` |
|---|---|
| `sv` | Supervisor |
| `op` | Operador |
| `ro` | Operador de riesgo |
| `ad` | Administrador |
| `fo` | Operador de feminicidio |
| `no` | Operador territorial nacional |
| `do` | Operador territorial departamental |
| `et` | Entidad |
| `an` | A notificaciones |
| `us` | Usuario externo |

---

### `rel_role_general_user`
Relación N:N entre usuarios y roles.

| Columna | Descripción |
|---|---|
| `rel_role_general_user_id` | PK |
| `role_id` | FK → `role.role_id` |
| `general_user_id` | FK → `general_user.general_user_id` |
| `rel_role_general_user_creation_date` | Fecha de asignación |

---

## Equipos configurados

| Equipo | Descripción |
|---|---|
| `Riesgo alto` | Casos de alto riesgo |
| `Riesgo bajo` | Casos de bajo riesgo |

El equipo se guarda como texto libre en `general_user.general_user_team`.

---

## Usuarios activos por equipo y rol

```json
{
  "equipos": {
    "Riesgo alto": {
      "op": [
        { "login": "agente.seguimiento", "nombres": "Carlos",   "apellidos": "Ramírez Torres" },
        { "login": "agenteriesgo",        "nombres": "Agente",   "apellidos": "Riesgo" },
        { "login": "aseguimiento",        "nombres": "Ricardo",  "apellidos": "Contreras" }
      ],
      "ro": [
        { "login": "Operador.V5",      "nombres": "Operador", "apellidos": "riesgov5" },
        { "login": "operadorRiesgo1",  "nombres": "operador", "apellidos": "riesgo1" },
        { "login": "operadorRiesgo3",  "nombres": "operador", "apellidos": "riesgo3" }
      ],
      "sv": [
        { "login": "felipe10",         "nombres": "Andres",     "apellidos": "Ascencio" },
        { "login": "super.fernandez",  "nombres": "supervisor", "apellidos": "fernandez" },
        { "login": "supervisorriesgo", "nombres": "supervisor", "apellidos": "riesgo" }
      ]
    },
    "Riesgo bajo": {
      "op": [
        { "login": "agentegeneral", "nombres": "Agente", "apellidos": "General" }
      ],
      "ro": [
        { "login": "afascencio10",    "nombres": "Andres",        "apellidos": "Ascencio" },
        { "login": "Nicolas.Ruiz",    "nombres": "Nicolas",       "apellidos": "Ruiz" },
        { "login": "operadorRiesgo2", "nombres": "operador",      "apellidos": "riesgo2" },
        { "login": "operadorRiesgo4", "nombres": "operador",      "apellidos": "riesgo4" },
        { "login": "riesgo.bajo",     "nombres": "Agente riesgo", "apellidos": "bajo" }
      ],
      "sv": [
        { "login": "Supervisorbajo",    "nombres": "Supervisor", "apellidos": "Buitrago" },
        { "login": "supervisorgeneral", "nombres": "supervisor", "apellidos": "general" }
      ]
    }
  }
}
```

> Contraseña de prueba para todos los usuarios: `Admin1234!`

---

## Pasos para crear un usuario nuevo

### 1. Generar el hash de la contraseña

Usar bcrypt con cost 10. En Go:

```go
hash, _ := bcrypt.GenerateFromPassword([]byte("Admin1234!"), 10)
```

O con `htpasswd` / cualquier generador bcrypt online con cost = 10.

---

### 2. Crear el perfil (`general_user_profile`)

```sql
INSERT INTO security.general_user_profile (
    general_user_profile_id,
    general_user_profile_i_code,
    general_user_profile_names,
    general_user_profile_last_names,
    general_user_profile_town
) VALUES (
    gen_random_uuid(),
    'profile-icode-<login>',   -- código único legible
    'Nombres',
    'Apellidos',
    '11001000'                 -- Bogotá por defecto
);
```

---

### 3. Crear el usuario (`general_user`)

```sql
INSERT INTO security.general_user (
    general_user_id,
    general_user_i_code,
    general_user_login,
    general_user_password,
    general_user_status,
    general_user_language,
    general_user_general_user_profile,
    general_user_team
) VALUES (
    gen_random_uuid(),
    'user-icode-<rol>-<login>',
    '<login>',
    '<bcrypt_hash>',
    'e',
    'sp',
    (SELECT general_user_profile_id FROM security.general_user_profile
     WHERE general_user_profile_i_code = 'profile-icode-<login>'),
    'Riesgo alto'  -- o 'Riesgo bajo'
);
```

---

### 4. Asignar el rol (`rel_role_general_user`)

```sql
INSERT INTO security.rel_role_general_user (
    rel_role_general_user_id,
    role_id,
    general_user_id
) VALUES (
    gen_random_uuid(),
    (SELECT role_id FROM security.role WHERE role_code = 'op'),  -- sv | op | ro
    (SELECT general_user_id FROM security.general_user WHERE general_user_login = '<login>')
);
```

---

### 5. Verificar el usuario creado

```sql
SELECT
    gu.general_user_login   AS login,
    gup.general_user_profile_names      AS nombres,
    gup.general_user_profile_last_names AS apellidos,
    r.role_code             AS rol,
    gu.general_user_team    AS equipo
FROM security.general_user gu
INNER JOIN security.general_user_profile gup
        ON gup.general_user_profile_id = gu.general_user_general_user_profile
INNER JOIN security.rel_role_general_user rr
        ON rr.general_user_id = gu.general_user_id
INNER JOIN security.role r
        ON r.role_id = rr.role_id
WHERE gu.general_user_login = '<login>';
```

---

## Cambiar el equipo de un usuario

```sql
UPDATE security.general_user
SET general_user_team = 'Riesgo bajo'   -- o 'Riesgo alto'
WHERE general_user_login = '<login>';
```

---

## Cambiar la contraseña de un usuario

Generar nuevo hash bcrypt (cost 10) y actualizar:

```sql
UPDATE security.general_user
SET general_user_password = '<nuevo_bcrypt_hash>'
WHERE general_user_login = '<login>';
```

---

## Desactivar un usuario

```sql
UPDATE security.general_user
SET general_user_status = 'i'
WHERE general_user_login = '<login>';
```

---

## Notas importantes

- El campo `general_user_team` es texto libre — no tiene tabla de catálogo. Los valores válidos actualmente son `Riesgo alto` y `Riesgo bajo` (respetar mayúsculas y espacios).
- Si un usuario no tiene perfil vinculado (`general_user_general_user_profile IS NULL`) el sistema lo filtra silenciosamente en la mayoría de vistas.
- Al arrancar el servidor se ejecuta automáticamente un `UPDATE` que asigna `11001000` (Bogotá) como municipio por defecto a los usuarios activos con perfil sin municipio, para evitar errores 500 en reasignación de casos.
