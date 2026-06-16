# Modelo `directories`

Tabla para almacenar entidades de contacto institucional que la app Flutter lista por ciudad.

---

## Ubicación

| Elemento | Valor propuesto |
|---|---|
| Schema | `salvia` |
| Tabla | `directories` |
| Modelo Go | `src/internal/models/directory.go` |
| Nombre struct | `Directory` |

---

## Diagrama de referencia (requerimiento inicial)

Campos según el diagrama entregado:

| Columna BD | Tipo diagrama | Nullable | Descripción |
|---|---|---|---|
| `id` | `uuid` / `varchar(36)` | No | PK (`gen_random_uuid()`) |
| `city_id` | `varchar(36)` | No | **Almacena `city_i_code`** de `security.city`, no el `city_id` numérico |
| `creation_date` | `timestamp` | No | Fecha de creación |
| `update_date` | `timestamp` | No | Fecha de última actualización |
| `name` | `varchar(12)` | No | Nombre de la entidad (ej. "Fiscalía Seccional Bogotá") |
| `address` | `varchar(128)` | No | Dirección física |
| `phone` | `varchar(2)` | Sí / No* | Teléfono de contacto |
| `email` | `varchar(20)` | Sí / No* | Correo electrónico |
| `opening_hours` | `varchar(20)` | Sí / No* | Horario de atención |
| `type` | `varchar(2)` | No | Tipo de entidad (catálogo por definir) |

\* Obligatoriedad de contacto (`phone`, `email`) pendiente de definir reglas de negocio.

---

## Propuesta de longitudes (recomendada para implementación)

Los límites del diagrama son demasiado cortos para los datos del mockup. Se sugiere ajustar en la implementación:

| Columna BD | Tipo propuesto | Motivo |
|---|---|---|
| `id` | `varchar(36) PRIMARY KEY DEFAULT gen_random_uuid()` | Convención GORM (como `entity_letter`) |
| `city_id` | `varchar(36) NOT NULL` | Coincide con `security.city.city_i_code` |
| `creation_date` | `timestamptz NOT NULL DEFAULT now()` | Convención del proyecto |
| `update_date` | `timestamptz NOT NULL DEFAULT now()` | Convención del proyecto |
| `name` | `varchar(128) NOT NULL` | "Fiscalía Seccional Bogotá" ≈ 26 caracteres |
| `address` | `varchar(256) NOT NULL` | Direcciones largas con edificio y numeración |
| `phone` | `varchar(32)` | Ej. `601-5701000` |
| `email` | `varchar(128)` | Ej. `fiscalia.bogota@fiscalia.gov.co` |
| `opening_hours` | `varchar(128)` | Ej. `Lunes a Viernes 8:00 AM - 5:00 PM` |
| `type` | `varchar(2) NOT NULL` | Código corto de categoría |

### SQL de referencia (propuesta)

```sql
CREATE TABLE salvia.directories (
    id              VARCHAR(36) PRIMARY KEY DEFAULT gen_random_uuid(),
    city_id         VARCHAR(36)  NOT NULL,  -- city_i_code de security.city
    creation_date   TIMESTAMPTZ  NOT NULL DEFAULT now(),
    update_date     TIMESTAMPTZ  NOT NULL DEFAULT now(),
    name            VARCHAR(128) NOT NULL,
    address         VARCHAR(256) NOT NULL,
    phone           VARCHAR(32),
    email           VARCHAR(128),
    opening_hours   VARCHAR(128),
    type            VARCHAR(2)   NOT NULL
);

CREATE INDEX idx_directories_city_id ON salvia.directories (city_id);
CREATE INDEX idx_directories_type ON salvia.directories (type);
```

---

## Relación con `security.city`

```
security.city
    city_id        (uint, PK interno)
    city_i_code    (varchar 36, UQ)  ← se persiste en directories.city_id
    city_name      (varchar 64)      ← se recibe en el endpoint como cityName
```

**Regla de negocio:** al crear o actualizar un directorio, el backend recibe el **nombre de la ciudad** (`cityName`), busca en `security.city` por coincidencia exacta (o normalizada) sobre `city_name`, y persiste el `city_i_code` resultante en `directories.city_id`.

**Consulta de referencia:**

```sql
SELECT city_i_code, city_name
FROM security.city
WHERE city_name = $1   -- ej. 'Bogotá D.C.'
LIMIT 1;
```

Si no hay coincidencia → error `404` o `400` con mensaje `"ciudad no encontrada"`.

> **Nota:** `CityLight` en `src/internal/models/location_light.go` ya proyecta `city_name`. El repositorio de directorios puede reutilizar GORM sobre `security.city` o extender `LocationRepository` con `GetCityICodeByName`.

---

## Mapeo UI Flutter → modelo

Referencia: tarjeta de directorio en la app.

| Elemento UI | Campo BD | Origen en respuesta API |
|---|---|---|
| Selector superior "Bogotá D.C." | — | Parámetro de filtro `cityName` en el `GET`; no se guarda en el registro |
| Título "Fiscalía Seccional Bogotá" | `name` | `name` |
| Badge "Bogotá D.C." | — | `cityName` resuelto desde `security.city` vía `city_id` |
| Ícono ubicación + dirección | `address` | `address` |
| Ícono reloj + horario | `opening_hours` | `openingHours` |
| Ícono sobre + correo (link) | `email` | `email` |
| Botón "Llamar: 601-5701000" | `phone` | `phone` (prefijo "Llamar:" es solo UI) |

---

## DTOs previstos

### Request — crear directorio (`CreateDirectoryRequest`)

Campos que envía quien carga el directorio (no incluye `city_id` resuelto):

```json
{
  "cityName": "Bogotá D.C.",
  "name": "Fiscalía Seccional Bogotá",
  "address": "Av. El Dorado No. 58-90, Edificio Fiscalía",
  "phone": "601-5701000",
  "email": "fiscalia.bogota@fiscalia.gov.co",
  "openingHours": "Lunes a Viernes 8:00 AM - 5:00 PM",
  "type": "FI"
}
```

| Campo JSON | Obligatorio | Validación |
|---|---|---|
| `cityName` | Sí | Debe existir en `security.city.city_name` |
| `name` | Sí | 1–128 caracteres |
| `address` | Sí | 1–256 caracteres |
| `phone` | Recomendado | 1–32 caracteres |
| `email` | Recomendado | Formato email, max 128 |
| `openingHours` | Recomendado | 1–128 caracteres |
| `type` | Sí | Código de 2 caracteres; ver catálogo abajo |

### Response — directorio persistido (`DirectoryResponse`)

```json
{
  "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "cityId": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "cityName": "Bogotá D.C.",
  "name": "Fiscalía Seccional Bogotá",
  "address": "Av. El Dorado No. 58-90, Edificio Fiscalía",
  "phone": "601-5701000",
  "email": "fiscalia.bogota@fiscalia.gov.co",
  "openingHours": "Lunes a Viernes 8:00 AM - 5:00 PM",
  "type": "FI",
  "creationDate": "2026-06-16T10:00:00Z",
  "updateDate": "2026-06-16T10:00:00Z"
}
```

`cityName` se **enriquece en lectura** haciendo JOIN o lookup a `security.city`; no se duplica en la tabla `directories`.

---

## Modelo GORM (implementado)

Archivo: `src/internal/models/directory.go`

- PK `id` uuid (`varchar(36)`, `gen_random_uuid()` — igual que `EntityLetter`)
- `city_id` → `security.city.city_i_code`
- Hooks `BeforeCreate` / `BeforeUpdate` para `creation_date` y `update_date`
- Validación de tipos vía `IsValidDirectoryType()` y constantes `DirectoryType*`
- Registrado en `AutoMigrate` de `main.go`

---

## Catálogo `type` (confirmado)

| Código | Etiqueta | Constante Go |
|---|---|---|
| `FI` | Fiscalías | `DirectoryTypeFiscalia` |
| `CF` | Comisarías de familia | `DirectoryTypeComisariaFamilia` |
| `RU` | Red de Urgencias | `DirectoryTypeRedUrgencias` |
| `LE` | Líneas de Emergencia | `DirectoryTypeLineasEmergencia` |

Mapa de etiquetas: `models.DirectoryTypeLabels`
