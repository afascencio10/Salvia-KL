# Documentación Completa de APIs - Sistema SALVIA

## Tabla de Contenidos

1. [Introducción](#introducción)
2. [Arquitectura de APIs](#arquitectura-de-apis)
3. [Autenticación y Autorización](#autenticación-y-autorización)
4. [Módulo Security](#módulo-security)
   - [Autenticación](#autenticación)
   - [Usuarios Generales](#usuarios-generales)
   - [Ciudades y Municipios](#ciudades-y-municipios)
5. [Módulo Salvia](#módulo-salvia)
   - [Contactos de Víctimas](#contactos-de-víctimas)
   - [Casos de Víctimas](#casos-de-víctimas)
   - [Alertas](#alertas)
   - [Momentos](#momentos)
   - [Ramas de Entidades](#ramas-de-entidades)
   - [Logs de Casos](#logs-de-casos)
   - [Asignación de Operadores](#asignación-de-operadores)
   - [Carga de Archivos](#carga-de-archivos)
6. [Consumo desde Frontend](#consumo-desde-frontend)
7. [Códigos de Respuesta](#códigos-de-respuesta)

---

## Introducción

Este documento describe todas las APIs REST del sistema SALVIA (Sistema de Atención a Víctimas de Violencia de Género).
El sistema está desarrollado en Go utilizando el framework Gin y sigue una arquitectura de 3 capas:

- **Capa de Presentación (Facades)**: Maneja las peticiones HTTP
- **Capa de Lógica (Controllers)**: Procesa la lógica de negocio
- **Capa de Datos (DAOs)**: Gestiona el acceso a la base de datos PostgreSQL

### Características Generales

- **Formato de Respuesta**: JSON o HTML según el header `Accept`
- **Autenticación**: Basada en sesiones con cookies
- **Autorización**: Sistema de roles y permisos
- **Idioma**: Español (código "sp")
- **Encoding**: UTF-8

---

## Arquitectura de APIs

### Formato de URLs

Las URLs siguen el patrón:
```
https://{host}:{port}/{módulo}/{recurso}/{acción}/{parámetros}
```

Ejemplo:
```
https://localhost:8443/salvia/contacto-victima/obtener/ABC123
```

### Tipos de Respuesta

#### Respuesta JSON
```json
{
  "icode": "ABC123",
  "name": "Nombre",
  "status": "A"
}
```

#### Respuesta HTML
Renderiza una plantilla Go con los datos embebidos.


---

## Autenticación y Autorización

### Sistema de Sesiones

El sistema utiliza sesiones basadas en cookies con el middleware `gin-contrib/sessions`.

**Flujo de Autenticación:**
1. Usuario envía credenciales a `/security/usuario-general/login`
2. Sistema valida credenciales y crea sesión
3. Se almacena `sessionID` en cookie
4. Cada petición incluye la cookie con el `sessionID`
5. Middleware valida la sesión antes de procesar la petición

### Roles del Sistema

| Código | Nombre | Descripción |
|--------|--------|-------------|
| `ad` | Administrador | Acceso completo al sistema |
| `op` | Operador | Gestiona casos de víctimas |
| `et` | Entidad Territorial | Acceso limitado a su municipio |
| `us` | Usuario Víctima | Acceso a su propio caso |
| `sv` | Supervisor | Supervisión de operadores |

### Middlewares Aplicados

1. **AuthMiddleware**: Verifica que el usuario esté autenticado
2. **LimitMiddleware**: Rate limiting para prevenir abuso
3. **HTTPS Redirect**: Fuerza conexiones seguras

---

## Módulo Security

### Autenticación

#### POST /security/usuario-general/login
**Descripción**: Autentica un usuario en el sistema

**Middlewares**: Ninguno (endpoint público)

**Request Body**:
```json
{
  "login": "usuario@example.com",
  "password": "contraseña",
  "captchaId": "captcha-id",
  "captchaSolution": "solución"
}
```

**Response (200 OK)**:
```json
{
  "role": "op",
  "path": "/salvia/contacto-victima",
  "title": "Contactos de Víctimas"
}
```

**Response (401 Unauthorized)**:
```json
{
  "error": "Credenciales inválidas"
}
```

**Consumo desde Frontend**:
```javascript
// Archivo: frontend/js/login.js
async function login(formData) {
  const response = await fetch('/security/usuario-general/login', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(formData)
  });
  return await response.json();
}
```

---

#### POST /security/usuario-general/logout
**Descripción**: Cierra la sesión del usuario actual

**Middlewares**: AuthMiddleware

**Request Body**: Ninguno

**Response (200 OK)**:
```json
{}
```

**Consumo desde Frontend**:
```javascript
async function logout() {
  await fetch('/security/usuario-general/logout', {
    method: 'POST'
  });
  window.location.href = '/security/usuario-general/login';
}
```

---

#### GET /security/usuario-general/login
**Descripción**: Muestra el formulario de login

**Middlewares**: Ninguno (endpoint público)

**Response**: HTML con formulario de login

---

#### POST /security/usuario-general/forgot
**Descripción**: Solicita reseteo de contraseña

**Request Body**:
```json
{
  "email": "usuario@example.com",
  "captchaId": "captcha-id",
  "captchaSolution": "solución"
}
```

**Response (200 OK)**:
```json
{
  "message": "Se ha enviado un correo con instrucciones"
}
```

---

### Usuarios Generales

#### GET /security/usuario-general
**Descripción**: Obtiene la lista de todos los usuarios

**Middlewares**: AuthMiddleware

**Permisos Requeridos**: `get_general_users`

**Roles Permitidos**: `ad`, `sv`

**Response (200 OK)**:
```json
[
  {
    "icode": "USER001",
    "login": "operador1@salvia.gov",
    "profile": {
      "names": "Juan",
      "lastNames": "Pérez",
      "docType": "CC",
      "docNumber": "12345678"
    },
    "roles": ["op"],
    "status": "A"
  }
]
```

**Consumo desde Frontend**:
```javascript
// Archivo: frontend/js/dao.js
async function getData(entity) {
  const response = await fetch(`/security/${entity}`, {
    headers: { 'Accept': 'application/json' }
  });
  return await response.json();
}

// Uso
const users = await getData('usuario-general');
```

---

#### GET /security/usuario-general/:id
**Descripción**: Obtiene un usuario específico por su código interno

**Middlewares**: AuthMiddleware

**Permisos Requeridos**: `get_general_user`

**Parámetros de Ruta**:
- `id` (string): Código interno del usuario

**Response (200 OK)**:
```json
{
  "icode": "USER001",
  "login": "operador1@salvia.gov",
  "profile": {
    "icode": "PROF001",
    "names": "Juan",
    "lastNames": "Pérez",
    "docType": "CC",
    "docNumber": "12345678",
    "gender": "M",
    "town": {
      "code": "05001",
      "name": "Medellín"
    }
  },
  "roles": ["op"],
  "status": "A",
  "creationDate": "2024-01-15T10:30:00Z"
}
```


---

#### POST /security/usuario-general
**Descripción**: Crea un nuevo usuario en el sistema

**Middlewares**: AuthMiddleware

**Permisos Requeridos**: `set_general_user`

**Request Body**:
```json
{
  "login": "nuevo@salvia.gov",
  "password": "contraseña",
  "profile": {
    "names": "María",
    "lastNames": "González",
    "docType": "CC",
    "docNumber": "87654321",
    "gender": "F",
    "townCode": "05001"
  },
  "roles": ["op"],
  "entityBranchICode": "ENT001"
}
```

**Response (201 Created)**:
```json
{
  "icode": "USER002",
  "message": "Usuario creado exitosamente"
}
```

**Consumo desde Frontend**:
```javascript
async function newEntity(entity, data) {
  const response = await fetch(`/security/${entity}`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data)
  });
  return await response.json();
}
```

---

#### PUT /security/usuario-general/:id
**Descripción**: Actualiza un usuario existente

**Middlewares**: AuthMiddleware

**Permisos Requeridos**: `update_general_user`

**Parámetros de Ruta**:
- `id` (string): Código interno del usuario

**Request Body**: Similar al POST

**Response (200 OK)**:
```json
{
  "message": "Usuario actualizado exitosamente"
}
```

---

#### PUT /security/usuario-general/:id/deshabilitar
**Descripción**: Deshabilita un usuario

**Middlewares**: AuthMiddleware

**Permisos Requeridos**: `update_general_user_disable`

**Response (200 OK)**:
```json
{
  "message": "Usuario deshabilitado"
}
```

---

#### PUT /security/usuario-general/:id/habilitar
**Descripción**: Habilita un usuario previamente deshabilitado

**Middlewares**: AuthMiddleware

**Permisos Requeridos**: `update_general_user_enable`

**Response (200 OK)**:
```json
{
  "message": "Usuario habilitado"
}
```

---

#### DELETE /security/usuario-general/:id
**Descripción**: Elimina un usuario del sistema

**Middlewares**: AuthMiddleware

**Permisos Requeridos**: `remove_general_user`

**Response (200 OK)**:
```json
{
  "message": "Usuario eliminado"
}
```


---

### Ciudades y Municipios

#### GET /security/ciudad/por-departamento/:id
**Descripción**: Obtiene las ciudades de un departamento específico

**Middlewares**: AuthMiddleware

**Permisos Requeridos**: `get_city_by_department`

**Parámetros de Ruta**:
- `id` (string): Código del departamento

**Response (200 OK)**:
```json
[
  {
    "icode": "CITY001",
    "code": "05",
    "name": "Antioquia"
  }
]
```

**Consumo desde Frontend**:
```javascript
// Cuando el usuario selecciona un departamento
async function loadCities(departmentCode) {
  const response = await fetch(
    `/security/ciudad/por-departamento/${departmentCode}`,
    { headers: { 'Accept': 'application/json' } }
  );
  return await response.json();
}
```

---

#### GET /security/municipio/por-ciudad/:id
**Descripción**: Obtiene los municipios de una ciudad específica

**Middlewares**: AuthMiddleware

**Permisos Requeridos**: `get_town_by_city_code`

**Parámetros de Ruta**:
- `id` (string): Código de la ciudad

**Response (200 OK)**:
```json
[
  {
    "icode": "TOWN001",
    "code": "05001",
    "name": "Medellín",
    "type": "M"
  }
]
```

---

## Módulo Salvia

### Contactos de Víctimas

#### GET /salvia/contacto-victima
**Descripción**: Obtiene la lista de contactos de víctimas sin caso asignado

**Middlewares**: AuthMiddleware

**Permisos Requeridos**: `get_victim_contacts`

**Roles Permitidos**: `op`, `sv`

**Response (200 OK)**:
```json
[
  {
    "icode": "CONT001",
    "names": "Ana",
    "lastNames": "Martínez",
    "docType": "CC",
    "docNumber": "11223344",
    "phone": "3001234567",
    "email": "ana@example.com",
    "status": "P",
    "creationDate": "2024-01-20T14:30:00Z"
  }
]
```

**Consumo desde Frontend**:
```javascript
// Archivo: frontend/js/victim_contact.js
async function loadContacts() {
  const response = await fetch('/salvia/contacto-victima', {
    headers: { 'Accept': 'application/json' }
  });
  return await response.json();
}
```

---

#### GET /salvia/contacto-victima/:id
**Descripción**: Obtiene un contacto específico por su código

**Middlewares**: AuthMiddleware

**Permisos Requeridos**: `get_victim_contact`

**Parámetros de Ruta**:
- `id` (string): Código interno del contacto

**Response (200 OK)**:
```json
{
  "icode": "CONT001",
  "names": "Ana",
  "lastNames": "Martínez",
  "docType": "CC",
  "docNumber": "11223344",
  "phone": "3001234567",
  "email": "ana@example.com",
  "gender": "F",
  "livingZone": "U",
  "occupation": "EMP",
  "status": "P"
}
```


---

#### POST /salvia/contacto-victima
**Descripción**: Crea un nuevo contacto de víctima

**Middlewares**: AuthMiddleware

**Permisos Requeridos**: `set_victim_contact`

**Request Body**:
```json
{
  "names": "Ana",
  "lastNames": "Martínez",
  "docType": "CC",
  "docNumber": "11223344",
  "phone": "3001234567",
  "email": "ana@example.com",
  "gender": "F",
  "livingZone": "U",
  "genderIdentity": "F",
  "sexualOrientation": "H",
  "origin": "COL",
  "occupation": "EMP",
  "language": "ES",
  "ethnicGroup": "N"
}
```

**Response (201 Created)**:
```json
{
  "icode": "CONT002",
  "message": "Contacto creado exitosamente"
}
```

**Consumo desde Frontend**:
```javascript
async function createContact(formData) {
  const response = await fetch('/salvia/contacto-victima', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(formData)
  });
  return await response.json();
}
```

---

#### PUT /salvia/contacto-victima/:id
**Descripción**: Invalida (marca como procesado) un contacto de víctima

**Middlewares**: AuthMiddleware

**Permisos Requeridos**: `put_victim_contact`

**Parámetros de Ruta**:
- `id` (string): Código interno del contacto

**Request Body**:
```json
{
  "reason": "Caso creado"
}
```

**Response (200 OK)**:
```json
{
  "message": "Contacto invalidado"
}
```

---

### Casos de Víctimas

#### GET /salvia/caso-victima
**Descripción**: Obtiene la lista de casos según el rol del usuario

**Middlewares**: AuthMiddleware

**Permisos Requeridos**: `get_victim_cases`

**Roles y Comportamiento**:
- `op`: Casos asignados al operador
- `et`: Casos del municipio y entidad territorial
- `us`: Caso propio del usuario víctima
- `sv`: Todos los casos del sistema

**Response (200 OK)**:
```json
[
  {
    "icode": "CASE001",
    "victimContact": {
      "icode": "CONT001",
      "names": "Ana",
      "lastNames": "Martínez"
    },
    "status": "A",
    "femicideRisk": "S",
    "townCode": "05001",
    "creationDate": "2024-01-21T09:00:00Z"
  }
]
```

**Consumo desde Frontend**:
```javascript
async function loadCases() {
  const response = await fetch('/salvia/caso-victima', {
    headers: { 'Accept': 'application/json' }
  });
  return await response.json();
}
```

---

#### GET /salvia/caso-victima/:id
**Descripción**: Obtiene un caso específico con toda su información detallada

**Middlewares**: AuthMiddleware

**Permisos Requeridos**: `get_victim_case`

**Parámetros de Ruta**:
- `id` (string): Código interno del caso

**Response (200 OK)**:
```json
{
  "icode": "CASE001",
  "victimContact": {
    "icode": "CONT001",
    "names": "Ana",
    "lastNames": "Martínez",
    "docType": "CC",
    "docNumber": "11223344"
  },
  "status": "A",
  "femicideRisk": "S",
  "townCode": "05001",
  "moments": [
    {
      "icode": "MOM001",
      "code": "01",
      "status": "A",
      "entityBranch": {
        "icode": "ENT001",
        "name": "Comisaría de Familia"
      }
    }
  ],
  "violenceExperienced": ["FIS", "PSI"],
  "creationDate": "2024-01-21T09:00:00Z"
}
```


---

#### GET /salvia/caso-victima/:docType/:id
**Descripción**: Busca casos por tipo y número de documento

**Middlewares**: AuthMiddleware

**Permisos Requeridos**: `get_victim_case_by_document`

**Parámetros de Ruta**:
- `docType` (string): Tipo de documento (CC, TI, CE, etc.)
- `id` (string): Número de documento

**Response (200 OK)**:
```json
[
  {
    "icode": "CASE001",
    "victimContact": {
      "docType": "CC",
      "docNumber": "11223344"
    },
    "status": "A"
  }
]
```

---

#### POST /salvia/caso-victima
**Descripción**: Crea un nuevo caso de víctima

**Middlewares**: AuthMiddleware

**Permisos Requeridos**: `set_victim_case`

**Request Body**:
```json
{
  "victimContactICode": "CONT001",
  "townCode": "05001",
  "departmentId": 5,
  "cityId": 1,
  "femicideRisk": "S",
  "violenceExperienced": ["FIS", "PSI"],
  "violenceScope": "DOM",
  "aggressor": "PAR",
  "relationshipWithAggressor": "ESP",
  "moments": [
    {
      "code": "01",
      "entityBranchId": 1,
      "status": "P"
    }
  ]
}
```

**Response (201 Created)**:
```json
{
  "icode": "CASE002",
  "message": "Caso creado exitosamente"
}
```

**Consumo desde Frontend**:
```javascript
async function createCase(caseData) {
  const response = await fetch('/salvia/caso-victima', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(caseData)
  });
  return await response.json();
}
```

---

#### POST /salvia/caso-victima/reporte
**Descripción**: Genera un reporte de casos con filtros

**Middlewares**: AuthMiddleware

**Permisos Requeridos**: `report_victim_cases`

**Request Body**:
```json
{
  "startDate": "2024-01-01",
  "endDate": "2024-01-31",
  "townCode": "05001",
  "status": "A",
  "femicideRisk": "S"
}
```

**Response (200 OK)**:
```json
[
  {
    "icode": "CASE001",
    "victimContact": {
      "names": "Ana",
      "lastNames": "Martínez"
    },
    "status": "A",
    "femicideRisk": "S",
    "creationDate": "2024-01-21T09:00:00Z"
  }
]
```

---

### Alertas

#### GET /salvia/alerta
**Descripción**: Obtiene las alertas según el rol del usuario

**Middlewares**: AuthMiddleware

**Permisos Requeridos**: `get_alerts_by_all`

**Roles y Comportamiento**:
- `op`: Alertas de casos asignados
- `et`: Alertas de casos asignados
- `sv`: Todas las alertas del sistema

**Response (200 OK)**:
```json
[
  {
    "icode": "ALERT001",
    "type": "RISK",
    "priority": 1,
    "code": "FEMICIDE_RISK",
    "data": "Riesgo alto de feminicidio",
    "creationDate": "2024-01-22T10:00:00Z",
    "victimCaseICode": "CASE001"
  }
]
```

**Consumo desde Frontend**:
```javascript
async function loadAlerts() {
  const response = await fetch('/salvia/alerta', {
    headers: { 'Accept': 'application/json' }
  });
  return await response.json();
}
```

---

#### GET /salvia/alerta/por-caso-victima/:id
**Descripción**: Obtiene las alertas de un caso específico

**Middlewares**: AuthMiddleware

**Permisos Requeridos**: `get_alerts_by_victim_case`

**Parámetros de Ruta**:
- `id` (string): Código interno del caso

**Response (200 OK)**:
```json
[
  {
    "icode": "ALERT001",
    "type": "RISK",
    "priority": 1,
    "code": "FEMICIDE_RISK",
    "data": "Riesgo alto de feminicidio",
    "creationDate": "2024-01-22T10:00:00Z"
  }
]
```


---

#### GET /salvia/alerta/por-municipio/:townCode
**Descripción**: Obtiene las alertas de un municipio específico

**Middlewares**: AuthMiddleware

**Permisos Requeridos**: `get_alerts_by_town_code`

**Parámetros de Ruta**:
- `townCode` (string): Código del municipio

**Response (200 OK)**: Similar al anterior

---

### Momentos

#### PUT /salvia/momento/:id/:momentCode/:entityBranchIcode
**Descripción**: Actualiza el estado de un momento en la ruta de atención

**Middlewares**: AuthMiddleware

**Permisos Requeridos**: `update_moment`

**Parámetros de Ruta**:
- `id` (string): Código interno del caso
- `momentCode` (string): Código del momento (01-08)
- `entityBranchIcode` (string): Código de la rama de entidad

**Request Body**:
```json
{
  "status": "A",
  "approvalSource": "E",
  "approvalDescription": "Atención realizada exitosamente",
  "canceled": "N"
}
```

**Response (200 OK)**:
```json
{
  "message": "Momento actualizado exitosamente"
}
```

**Consumo desde Frontend**:
```javascript
async function updateMoment(caseId, momentCode, branchCode, data) {
  const response = await fetch(
    `/salvia/momento/${caseId}/${momentCode}/${branchCode}`,
    {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data)
    }
  );
  return await response.json();
}
```

---

### Ramas de Entidades

#### GET /salvia/rama-entidad/:townCode/con-momentos
**Descripción**: Obtiene las ramas de entidades de un municipio con sus momentos

**Middlewares**: AuthMiddleware

**Permisos Requeridos**: `get_entity_branches_with_moments`

**Parámetros de Ruta**:
- `townCode` (string): Código del municipio

**Response (200 OK)**:
```json
{
  "01": {
    "SALUD": {
      "ENT001": [
        {
          "icode": "BRANCH001",
          "name": "Hospital San Juan",
          "address": "Calle 10 # 20-30",
          "townCode": "05001",
          "entityIcode": "ENT001",
          "entityName": "Salud",
          "sector": "SALUD",
          "moment": "01"
        }
      ]
    }
  }
}
```

**Consumo desde Frontend**:
```javascript
async function loadBranchesWithMoments(townCode) {
  const response = await fetch(
    `/salvia/rama-entidad/${townCode}/con-momentos`,
    { headers: { 'Accept': 'application/json' } }
  );
  return await response.json();
}
```

---

#### GET /salvia/rama-entidad/:entityICode/:townCode
**Descripción**: Obtiene las ramas de una entidad específica en un municipio

**Middlewares**: AuthMiddleware

**Permisos Requeridos**: `get_entity_branches_by_towncode`

**Parámetros de Ruta**:
- `entityICode` (string): Código de la entidad
- `townCode` (string): Código del municipio

**Response (200 OK)**:
```json
[
  {
    "icode": "BRANCH001",
    "name": "Hospital San Juan",
    "address": "Calle 10 # 20-30",
    "latitude": 6.2442,
    "longitude": -75.5812
  }
]
```

---

### Logs de Casos

#### GET /salvia/log-caso/por-momento/:id
**Descripción**: Obtiene los logs de un momento específico

**Middlewares**: AuthMiddleware

**Permisos Requeridos**: `get_case_logs`

**Parámetros de Ruta**:
- `id` (string): Código interno del momento

**Response (200 OK)**:
```json
[
  {
    "icode": "LOG001",
    "description": "Se realizó atención inicial",
    "creationDate": "2024-01-22T11:00:00Z",
    "user": "USER001",
    "roles": "Operador",
    "profile": {
      "names": "Juan",
      "lastNames": "Pérez"
    }
  }
]
```


---

#### POST /salvia/log-caso/:id
**Descripción**: Crea un nuevo log para un caso

**Middlewares**: AuthMiddleware

**Permisos Requeridos**: `set_case_log`

**Parámetros de Ruta**:
- `id` (string): Código interno del caso

**Request Body**:
```json
{
  "momentId": 1,
  "description": "Se realizó seguimiento telefónico"
}
```

**Response (201 Created)**:
```json
{
  "icode": "LOG002",
  "message": "Log creado exitosamente"
}
```

**Consumo desde Frontend**:
```javascript
async function createLog(caseId, logData) {
  const response = await fetch(`/salvia/log-caso/${caseId}`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(logData)
  });
  return await response.json();
}
```

---

### Asignación de Operadores

#### POST /salvia/asignar-operadores/todos/:p1/:p2
**Descripción**: Asigna múltiples operadores a casos

**Middlewares**: AuthMiddleware

**Permisos Requeridos**: `assign_operators`

**Parámetros de Ruta**:
- `p1` (string): Código del op
erador
- `p2` (string): Filtro de casos

**Response (200 OK)**:
```json
{
  "assigned": 5,
  "message": "5 casos asignados exitosamente"
}
```

---

#### POST /salvia/asignar-operadores/individual/:p1/:p2
**Descripción**: Asigna un operador a un caso específico

**Middlewares**: AuthMiddleware

**Permisos Requeridos**: `assign_operators`

**Parámetros de Ruta**:
- `p1` (string): Código del caso
- `p2` (string): Código del operador

**Response (200 OK)**:
```json
{
  "message": "Operador asignado exitosamente"
}
```

---

### Carga de Archivos

#### POST /salvia/archivos-planos/victim_service
**Descripción**: Carga un archivo plano para actualizar momentos

**Middlewares**: AuthMiddleware

**Permisos Requeridos**: `load_plain_files`

**Request**: Multipart form-data con archivo

**Response (200 OK)**:
```json
{
  "processed": 10,
  "errors": 0,
  "message": "Archivo procesado exitosamente"
}
```

**Consumo desde Frontend**:
```javascript
async function uploadFile(file, service) {
  const formData = new FormData();
  formData.append('file', file);
  
  const response = await fetch(`/salvia/archivos-planos/${service}`, {
    method: 'POST',
    body: formData
  });
  return await response.json();
}
```

---

#### POST /salvia/archivos-planos/entity_branch/:id
**Descripción**: Carga un archivo para actualizar ramas de entidades

**Middlewares**: AuthMiddleware

**Permisos Requeridos**: `load_plain_files`

**Parámetros de Ruta**:
- `id` (string): Código de la entidad

**Request**: Multipart form-data con archivo

**Response (200 OK)**: Similar al anterior

---

## Consumo desde Frontend

### Arquitectura del Frontend

El frontend utiliza un enfoque híbrido:
- **Go Templates**: Renderizado inicial del HTML
- **Vue.js**: Interactividad y gestión de estado
- **JavaScript Vanilla**: Capa de acceso a datos (DAO)

### Capa de Acceso a Datos (dao.js)

```javascript
// Archivo: frontend/js/dao.js

// Función genérica para obtener datos
async function getData(entity, params = '') {
  const url = params ? `/${entity}/${params}` : `/${entity}`;
  const response = await fetch(url, {
    headers: { 'Accept': 'application/json' }
  });
  
  if (!response.ok) {
    throw new Error(`Error ${response.status}: ${response.statusText}`);
  }
  
  return await response.json();
}

// Función genérica para crear entidades
async function newEntity(entity, data) {
  const response = await fetch(`/${entity}`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data)
  });
  
  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.message || 'Error al crear');
  }
  
  return await response.json();
}

// Función genérica para actualizar entidades
async function updateEntity(entity, id, data) {
  const response = await fetch(`/${entity}/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data)
  });
  
  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.message || 'Error al actualizar');
  }
  
  return await response.json();
}

// Función genérica para eliminar entidades
async function removeEntity(entity, id) {
  const response = await fetch(`/${entity}/${id}`, {
    method: 'DELETE'
  });
  
  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.message || 'Error al eliminar');
  }
  
  return await response.json();
}
```

### Ejemplo de Uso en Vue.js

```javascript
// Componente Vue para gestión de contactos
export default {
  data() {
    return {
      contacts: [],
      loading: false,
      error: null
    }
  },
  
  async mounted() {
    await this.loadContacts();
  },
  
  methods: {
    async loadContacts() {
      this.loading = true;
      this.error = null;
      
      try {
        this.contacts = await getData('salvia/contacto-victima');
      } catch (error) {
        this.error = error.message;
        console.error('Error loading contacts:', error);
      } finally {
        this.loading = false;
      }
    },
    
    async createContact(formData) {
      try {
        const result = await newEntity('salvia/contacto-victima', formData);
        await this.loadContacts(); // Recargar lista
        return result;
      } catch (error) {
        this.error = error.message;
        throw error;
      }
    }
  }
}
```

---

## Códigos de Respuesta

### Códigos de Éxito

| Código | Descripción | Uso |
|--------|-------------|-----|
| 200 | OK | Operación exitosa (GET, PUT, DELETE) |
| 201 | Created | Recurso creado exitosamente (POST) |

### Códigos de Error del Cliente

| Código | Descripción | Causa Común |
|--------|-------------|-------------|
| 400 | Bad Request | Datos inválidos en el request |
| 401 | Unauthorized | No autenticado |
| 403 | Forbidden | Sin permisos para la operación |
| 404 | Not Found | Recurso no encontrado |

### Códigos de Error del Servidor

| Código | Descripción | Causa Común |
|--------|-------------|-------------|
| 500 | Internal Server Error | Error en el servidor |
| 503 | Service Unavailable | Base de datos no disponible |

### Formato de Errores

```json
{
  "error": "Descripción del error",
  "field": "nombre_campo",
  "code": "ERROR_CODE"
}
```

---

## Notas Finales

### Seguridad

- Todas las APIs (excepto login público) requieren autenticación
- Las sesiones expiran después de 30 minutos de inactividad
- Se implementa rate limiting para prevenir abuso
- Todas las comunicaciones deben ser por HTTPS

### Paginación

Actualmente el sistema no implementa paginación. Todas las consultas devuelven el conjunto completo de resultados.

### Versionado

El sistema no implementa versionado de APIs. Cualquier cambio breaking debe ser comunicado con anticipación.

### Soporte

Para reportar problemas o solicitar nuevas funcionalidades, contactar al equipo de desarrollo.

---

**Documento generado**: 2024-01-23  
**Versión del Sistema**: 1.0  
**Última actualización**: 2024-01-23
