# Endpoint de Simulación de Casos

Permite crear un caso completo en el sistema **sin usar el formulario del front-end** y **sin autenticación**. Útil para pruebas, QA y desarrollo.

---

## URL

```
POST /public/simular/caso
```

## Headers

```
Content-Type: application/json
```

---

## Body

Todos los campos son opcionales. El sistema rellena los datos faltantes con valores por defecto válidos.

| Campo | Tipo | Default | Descripción |
|---|---|---|---|
| `nombres` | string | `"María Prueba"` | Nombres de la víctima |
| `apellidos` | string | `"Simulación Test"` | Apellidos de la víctima |
| `tipoDoc` | string | `"cc"` | Tipo de documento: `cc`, `ce`, `cd`, `ti`, `pn`, etc. |
| `numDoc` | string | Timestamp aleatorio | Número de documento |
| `esPareja` | boolean | `false` | `true` si el agresor es pareja íntima (cambia el tamizaje) |
| `nivelRiesgo` | integer | `3` | Nivel de riesgo deseado: `1`=Bajo, `2`=Moderado, `3`=Alto, `4`=Extremo |
| `municipioAtencion` | string | `"11001000"` (Bogotá) | Código del municipio de atención |
| `municipioResidencia` | string | `"11001000"` | Código del municipio de residencia de la víctima |
| `municipioHechos` | string | `"11001000"` | Código del municipio donde ocurrieron los hechos |
| `operadorICode` | string | Andres Ascencio | ICode del operador que registra el caso |

---

## Cómo funciona el `nivelRiesgo`

El endpoint calcula automáticamente las respuestas del tamizaje (sí/no) para producir exactamente el nivel pedido, siguiendo la misma lógica del backend.

**Agresor es pareja (`esPareja: true`)**

| Nivel | Puntaje | Respuestas "sí" |
|---|---|---|
| 1 - Bajo | 0–4 | 0 |
| 2 - Moderado | 5–8 | 5 |
| 3 - Alto | 9–15 | 9 |
| 4 - Extremo | 16–24 | 16 |

**Agresor no es pareja (`esPareja: false`)**

| Nivel | Puntaje | Respuestas "sí" |
|---|---|---|
| 1 - Bajo | 0–2 | 0 |
| 2 - Moderado | 3–5 | 3 |
| 3 - Alto | 6–8 | 6 |
| 4 - Extremo | 9–20 | 9 |

---

## Ejemplos

### Caso mínimo (todo por defecto)
```bash
curl -X POST https://salvia.kloustr.com/public/simular/caso \
  -H "Content-Type: application/json" \
  -d '{}'
```

---

### Riesgo Bajo — No pareja
```bash
curl -X POST https://salvia.kloustr.com/public/simular/caso \
  -H "Content-Type: application/json" \
  -d '{
    "nombres": "Ana Lucía",
    "apellidos": "Martínez Gómez",
    "tipoDoc": "cc",
    "numDoc": "1012000001",
    "esPareja": false,
    "nivelRiesgo": 1
  }'
```

---

### Riesgo Moderado — Pareja íntima
```bash
curl -X POST https://salvia.kloustr.com/public/simular/caso \
  -H "Content-Type: application/json" \
  -d '{
    "nombres": "Laura Milena",
    "apellidos": "García Torres",
    "tipoDoc": "cc",
    "numDoc": "1098765432",
    "esPareja": true,
    "nivelRiesgo": 2
  }'
```

---

### Riesgo Alto — Pareja íntima, Cali
```bash
curl -X POST https://salvia.kloustr.com/public/simular/caso \
  -H "Content-Type: application/json" \
  -d '{
    "nombres": "Sofía",
    "apellidos": "Ramírez",
    "tipoDoc": "cc",
    "numDoc": "1234567890",
    "esPareja": true,
    "nivelRiesgo": 3,
    "municipioAtencion": "76001000",
    "municipioResidencia": "76001000",
    "municipioHechos": "76001000"
  }'
```

---

### Riesgo Extremo — Pareja íntima
```bash
curl -X POST https://salvia.kloustr.com/public/simular/caso \
  -H "Content-Type: application/json" \
  -d '{
    "nombres": "Valentina",
    "apellidos": "Pérez Cruz",
    "tipoDoc": "cc",
    "numDoc": "1099999001",
    "esPareja": true,
    "nivelRiesgo": 4
  }'
```

---

### Con operador específico
```bash
curl -X POST https://salvia.kloustr.com/public/simular/caso \
  -H "Content-Type: application/json" \
  -d '{
    "nombres": "Carolina",
    "apellidos": "Díaz",
    "tipoDoc": "cc",
    "numDoc": "1011000099",
    "esPareja": false,
    "nivelRiesgo": 3,
    "operadorICode": "019e07ac-5ff6-7956-b3a7-f22449c9a38f"
  }'
```

---

## Respuesta exitosa

```json
{
  "login": "usr_ab3x9k",
  "pass": "Kt7p"
}
```

Son las credenciales del usuario creado para la víctima. El caso queda registrado en la base de datos exactamente igual que si lo hubiera creado un operador desde el formulario.

## Respuesta con error

```json
{
  "form2": {
    "riskScore": "El puntaje de riesgo no coincide con el calculado por el servidor"
  }
}
```

Los errores siguen el mismo formato que el endpoint real — campo por campo.

---

## Notas

- El endpoint llama exactamente la misma lógica que el formulario del front-end (`SetVictimCase`), incluyendo creación de usuario, generación del calendario de seguimientos y asignación al operador.
- Cada llamada crea un caso real en la base de datos. No es un mock — úsalo con cuidado en producción.
- El `numDoc` debe ser único por `tipoDoc`. Si ya existe ese documento en el sistema, el caso se asocia al usuario existente.
- Los códigos de municipio (`11001000` = Bogotá, `76001000` = Cali, `05001000` = Medellín) vienen de la tabla `security.town`.
