# `case-mecanismos` — Interfaz del Componente

Tarjeta de la cadena de escalamiento del caso ante mecanismos articuladores (Municipal, Departamental, Nacional). Cada registro asigna un usuario con rol `mc` cuyo **tipo de mecanismo** fue registrado en el alta del usuario (`salvia.articulating_mechanism`).

Montado en el tab **Gestión institucional** del Detalle del Caso, después de `case-entities` y antes de la lista de barreras.

**Permisos:**
- Ver: cualquier rol con `get_case_detail_sv`.
- Agregar registro: solo `sv`, `op`, `ro`.
- El input de comentarios está retirado; el listado existente se mantiene hasta definir la lógica nueva.

## Archivos relevantes

| Archivo | Rol |
|---|---|
| `src/frontend/js/components/case-mecanismos.js` | Componente Vue (card + modal) |
| `src/frontend/html/salvia/case_detail/get_case_detail_sv.html` | Montaje en el tab Gestión institucional |
| `src/salvia/controller/case_articulating_mechanism_controller.go` | REST |
| `src/salvia/service/case_articulating_mechanism_service.go` | Validación por tipo + timeline |
| `src/internal/repository/case_articulating_mechanism_repository.go` | Queries |
| `src/internal/models/articulating_mechanism.go` | Registro 1:1 del usuario mc |
| `src/internal/models/case_articulating_mechanism.go` | Escalamiento del caso |

## Árbol de interfaz

```
case-mecanismos
│
├── Header
│   ├── Titulo  "Mecanismos"  (icono balanza)
│   └── [v-if puedeAgregar]  BtnAgregar  "+ Agregar mecanismo"  → abrirModal()
│
├── [v-if cargando]  "Cargando mecanismos..."
├── [v-else-if error]  Error + Reintentar
├── [v-else-if mecanismos.length === 0]
│   └── EmptyState  "Este caso no ha sido escalado a ningún mecanismo articulador"
│
└── [v-else]  Lista
    └── Item × N
        ├── BadgeAlcance  Municipal | Departamental | Nacional
        ├── Nombre  assignedUserName
        ├── [ubicación]  pin + municipio, ciudad (departamento)
        ├── [compromiso]  "Compromiso de Salvia"
        └── Comentarios
            ├── Titulo  "Comentarios del mecanismo (N)"
            ├── [vacío]  "Sin comentarios aún"
            └── Comentario × N  body + authorName + fecha
            (sin input de alta; lógica pendiente)

─── MODAL: AGREGAR MECANISMO ARTICULADOR ──────────────────────────────────

Overlay  @click.self="cerrarModal"
  ├── Titulo  "Agregar mecanismo articulador"
  ├── Subtitulo  "Instancia a la que se escaló o se va a escalar este caso"
  ├── <select> Tipo mecanismo*  municipal | departamental | nacional
  ├── [municipal]  Departamento* + Ciudad* + Municipio*
  ├── [departamental]  Departamento*
  ├── <select> Mecanismo*  usuarios mc con articulating_mechanism de ese tipo y ubicación
  ├── <textarea> Compromiso de Salvia (opcional)
  └── Cancelar | Agregar mecanismo
```


## Archivos relevantes

| Archivo | Rol |
|---|---|
| `src/frontend/js/components/case-mecanismos.js` | Componente Vue (card + modal) |
| `src/frontend/html/salvia/case_detail/get_case_detail_sv.html` | Montaje en el tab Gestión institucional |
| `src/salvia/controller/case_articulating_mechanism_controller.go` | REST |
| `src/salvia/service/case_articulating_mechanism_service.go` | Validación por alcance + timeline |
| `src/internal/repository/case_articulating_mechanism_repository.go` | Queries |
| `src/internal/models/case_articulating_mechanism.go` | Modelos GORM |

## Árbol de interfaz

```
case-mecanismos
│
├── Header
│   ├── Titulo  "Mecanismos"  (icono balanza)
│   └── [v-if puedeAgregar]  BtnAgregar  "+ Agregar mecanismo"  → abrirModal()
│
├── [v-if cargando]  "Cargando mecanismos..."
├── [v-else-if error]  Error + Reintentar
├── [v-else-if mecanismos.length === 0]
│   └── EmptyState  "Este caso no ha sido escalado a ningún mecanismo articulador"
│
└── [v-else]  Lista
    └── Item × N
        ├── BadgeAlcance  Municipal | Departamental | Nacional
        ├── Nombre  assignedUserName
        ├── [ubicación]  pin + municipio, ciudad (departamento)
        ├── [compromiso]  "Compromiso de Salvia"
        └── Comentarios
            ├── Titulo  "Comentarios del mecanismo (N)"
            ├── [vacío]  "Sin comentarios aún"
            ├── Comentario × N  body + authorName + fecha
            └── [v-if puedeAgregar]  textarea + "Agregar comentario"

─── MODAL: AGREGAR MECANISMO ARTICULADOR ──────────────────────────────────

Overlay  @click.self="cerrarModal"
  ├── Titulo  "Agregar mecanismo articulador"
  ├── Subtitulo  "Instancia a la que se escaló o se va a escalar este caso"
  ├── <select> Alcance*  municipal | departamental | nacional
  ├── [municipal]  Departamento* + Ciudad* + Municipio*
  ├── [departamental]  Departamento* + Ciudad (opcional)
  ├── <select> Mecanismo*  usuarios mc filtrados por alcance/territorio
  ├── <textarea> Compromiso de Salvia (opcional)
  └── Cancelar | Agregar mecanismo
```
