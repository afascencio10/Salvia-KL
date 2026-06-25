# `reportes-component` — Interfaz del Componente

Tabla reutilizable para visualizar y filtrar reportes de violencia basada en género registrados por víctimas o terceros. Los datos provienen de `salvia.victim_contact` y `salvia.victim_contact_form2`. **Solo lista reportes que aún no tienen un caso asociado** (`victim_case` inexistente) y con estado válido (`victim_contact_status = 'v'`). Acepta props para configurar columnas visibles y botones de acción opcionales por fila. El filtro expone dos inputs independientes: nombre y teléfono de la víctima. Emite eventos hacia el padre cuando el usuario presiona un botón de acción en una fila.

## Archivos relevantes

| Archivo | Rol |
|---|---|
| `src/frontend/components/reportes-component.js` | Componente principal — template y lógica *(pendiente de crear)* |
| `src/frontend/css/reportes-component.css` | Estilos de layout y estados visuales *(pendiente de crear)* |

---

## Árbol de interfaz

```
reportes-component  (.rc-wrapper)
│
├── [v-if loading]
│   LoadingState  (.rc-loading)
│   └── Spinner + "Cargando reportes..."
│
├── [v-else-if loadError]
│   ErrorState  (.rc-error)
│   └── AlertText  (.alert.alert-danger)  loadError
│
└── [v-else]
    │
    ├── FilterBar  (.rc-filter-bar)
    │   │
    │   ├── FilterGroup  (.rc-filter-group)
    │   │   ├── NameFilter  (.rc-name-filter)
    │   │   │   // Input independiente — filtra por nombre o apellido de la víctima.
    │   │   │   // Llama al backend (→ E-02) con debounce 400ms.
    │   │   │   ├── Label  (.rc-filter-label)  "Nombre de la víctima"
    │   │   │   └── <input type="text">  (.rc-name-input)
    │   │   │       placeholder="Buscar por nombres o apellidos"
    │   │   │       v-model=searchName
    │   │   │       → onNameFilterInput()  [debounce 400ms → E-02]
    │   │   │
    │   │   └── PhoneFilter  (.rc-phone-filter)
    │   │       // Input independiente — filtra por teléfono de la víctima.
    │   │       // Combina con el filtro de nombre si ambos tienen valor (AND).
    │   │       // Llama al backend (→ E-02) con debounce 400ms.
    │   │       ├── Label  (.rc-filter-label)  "Teléfono de la víctima"
    │   │       └── <input type="text">  (.rc-phone-input)
    │   │           placeholder="Buscar por teléfono de contacto"
    │   │           v-model=searchPhone
    │   │           → onPhoneFilterInput()  [debounce 400ms → E-02]
    │   │
    │   └── SortSelector  (.rc-sort-selector)
    │       ├── Label  (.rc-filter-label)  "Ordenar por"
    │       └── <select>  (.rc-select)  :value="sortOption"
    │           ├── <option value="registration_date_desc">  "Fecha de registro DESC"
    │           └── <option value="registration_date_asc">   "Fecha de registro ASC"
    │           → onSortChange(event)   // → E-03 (recarga backend)
    │
    └── TableWrap  (.rc-table-wrap)
        │
        ├── [v-if filteredReports.length === 0]
        │   EmptyState  (.rc-empty)
        │   └── "No se encontraron reportes"
        │
        └── [v-else]
            │
            ├── <table.rc-table>
            │   │
            │   ├── <thead>
            │   │   └── <tr>
            │   │       ├── <th> × N  [v-for visibleColumns]  (.rc-th)
            │   │       │   └── column.label
            │   │       └── [v-if buttons.length > 0]
            │   │           <th>  (.rc-th)  "Acciones"
            │   │
            │   └── <tbody>
            │       └── <tr> × N  [v-for filteredReports]  (.rc-row)
            │           │
            │           ├── [col.key === 'registration_date']
            │           │   <td>  (.rc-cell)
            │           │   └── report.creationDate  (formateada DD/MM/YYYY HH:mm)
            │           │
            │           ├── [col.key === 'adjustments_gbv']
            │           │   <td>  (.rc-cell-adjustments)
            │           │   └── formatAdjustmentsGBV(report.adjustmentsGBV)
            │           │       // Etiquetas legibles separadas por coma
            │           │       // Ejemplo: "No informa" o "Intérprete de lengua de señas colombiana, Transporte seguro..."
            │           │
            │           ├── [col.key === 'report_type']
            │           │   <td>  (.rc-cell)
            │           │   └── report.reportTypeLabel || "—"
            │           │       // victim_contact_form2_report_type → enum name
            │           │       // Ej: "Usted es la víctima de violencia basada en género"
            │           │
            │           ├── [col.key === 'will_receive_call']
            │           │   <td>  (.rc-cell)
            │           │   ├── [report.reportTypeCode === 'v']
            │           │   │   report.willReceiveCallLabel   // "Sí" | "No" | "—"
            │           │   └── [v-else]
            │           │       "—"   // Solo aplica cuando el reportante es la víctima
            │           │
            │           ├── [col.key === 'victim_names']
            │           │   <td>  (.rc-cell)
            │           │   └── report.names || "—"
            │           │       // victim_contact.victim_contact_names
            │           │
            │           ├── [col.key === 'victim_last_names']
            │           │   <td>  (.rc-cell)
            │           │   └── report.lastNames || "—"
            │           │       // victim_contact.victim_contact_last_names
            │           │
            │           ├── [col.key === 'victim_phone']
            │           │   <td>  (.rc-cell)
            │           │   └── formatPhone(report.victimPhone)
            │           │       // victim_contact_form2_victim_col_phone
            │           │
            │           ├── [col.key === 'best_contact_time']
            │           │   <td>  (.rc-cell)
            │           │   └── report.bestContactTime  (formateada hh:mm)
            │           │       // victim_contact_form2_best_contact_time
            │           │
            │           ├── [col.key === 'facts_description']
            │           │   <td>  (.rc-cell-description)
            │           │   └── truncateText(report.factsDescription, 120)
            │           │       // Texto truncado con ellipsis; tooltip con texto completo
            │           │
            │           └── [v-if buttons.length > 0]
            │               <td>  (.rc-cell-actions)
            │               └── ActionBtn × N  [v-for buttons]  (.rc-action-btn)
            │                   btn.label
            │                   → emitActionClicked(btn.id, report)
            │
            └── PaginationBar  (.rc-pagination)
                ├── PrevBtn  (.rc-page-btn)  "← Anterior"
                │   :disabled si currentPage === 1
                │   → changePage(currentPage - 1)
                ├── PageInfo  (.rc-page-info)
                │   "Página {currentPage} de {totalPages}"
                └── NextBtn  (.rc-page-btn)  "Siguiente →"
                    :disabled si currentPage === totalPages
                    → changePage(currentPage + 1)
```

---

## Props del componente

| Prop | Tipo | Requerido | Default | Descripción |
|---|---|---|---|---|
| `columns` | `Array<Object>` | Sí | — | Columnas a mostrar. Cada objeto: `{ key: string, label: string }` |
| `hiddenColumns` | `Array<string>` | No | `[]` | Keys de columnas a ocultar aunque estén en `columns` |
| `buttons` | `Array<Object>` | No | `[]` | Botones de acción por fila *(opcional)*. Cada objeto: `{ id: string, label: string }`. El padre puede pasar tantos botones como necesite. Si el array está vacío, no se renderiza la columna "Acciones". |

> **Nota:** A diferencia de `casos-component`, este componente no recibe `defaultFilter`, `availableFilters` ni `reasignacion`. Los filtros de usuario son dos inputs independientes: nombre y teléfono (E-02).

### Botones de acción (`buttons`)

| Condición | Comportamiento |
|---|---|
| `buttons === []` (default) | No se muestra columna "Acciones" |
| `buttons.length >= 1` | Se renderiza una columna "Acciones" con un botón por cada entrada del array |
| Al presionar un botón | El componente **solo emite** `action-clicked` al padre con `{ buttonId, report }`. No navega, no abre modales ni llama al backend |

---

## Columnas sugeridas (configuración por defecto del padre)

Corresponden a los campos visibles del formulario público de reporte (`set_victim_contact.html`):

| Key | Label sugerido | Fuente de dato |
|---|---|---|
| `registration_date` | Fecha de registro | `victim_contact.victim_contact_creation_date` |
| `adjustments_gbv` | Ajuste razonable | `rel_victim_case_form2_enums_victim_contact_form2` → enum `victim_case_form2_adjustments_gbv` |
| `report_type` | Tipo de reporte | `victim_contact_form2.victim_contact_form2_report_type` → enum name |
| `will_receive_call` | Dispuesto(a) a recibir llamada | `victim_contact_form2.victim_contact_form2_will_receive_call` → enum yes/no |
| `victim_names` | Nombres de la víctima | `victim_contact.victim_contact_names` |
| `victim_last_names` | Apellidos de la víctima | `victim_contact.victim_contact_last_names` |
| `victim_phone` | Teléfono de contacto | `victim_contact_form2.victim_contact_form2_victim_col_phone` |
| `best_contact_time` | Mejor hora de contacto | `victim_contact_form2.victim_contact_form2_best_contact_time` |
| `facts_description` | Descripción de los hechos | `victim_contact_form2.victim_contact_form2_facts_description` |

---

## Formato de columnas especiales

### `adjustments_gbv`

| Condición | Texto en celda |
|---|---|
| Sin ajustes seleccionados | `"—"` |
| Un ajuste | Etiqueta del enum (ej. `"No informa"`) |
| Varios ajustes | Etiquetas separadas por coma (ej. `"Intérprete de lengua de señas colombiana, Transporte seguro..."`) |

Los ajustes razonables se almacenan en la tabla relacional `rel_victim_case_form2_enums_victim_contact_form2` vinculada al registro de `victim_contact_form2`.

### `will_receive_call`

| `reportTypeCode` | Comportamiento |
|---|---|
| `'v'` (víctima) | Mostrar etiqueta del enum yes/no (`willReceiveCallLabel`) |
| `'r'` (tercero) u otro | Mostrar `"—"` — el campo no aplica en el formulario para reportes de tercero |

### `facts_description`

| Condición | Comportamiento |
|---|---|
| Texto ≤ 120 caracteres | Mostrar completo |
| Texto > 120 caracteres | Truncar con `"..."` y mostrar texto completo en `title` (tooltip nativo) |

---

## Alcance de datos

El listado incluye únicamente reportes que cumplen **ambas** condiciones:

| Condición | Implementación |
|---|---|
| Sin caso asociado | `LEFT JOIN salvia.victim_case` sobre `victim_case_victim_contact = victim_contact_id` → `WHERE victim_case_victim_contact IS NULL` |
| Estado válido | `victim_contact_status = 'v'` |

Esta lógica replica `GetVictimContactsWithoutVictimCase(page, "v", ...)` del DAO existente, ampliada con JOIN a `victim_contact_form2` para los campos del formulario de reporte.

### Campo `victim_contact_status`

Código de un carácter en `salvia.victim_contact.victim_contact_status`. En el código actual de Salvia se usan **minúsculas** (`v`, `i`), no los valores documentados en `MODELO_DATOS.md` (`A`, `P`, `I`), que están desactualizados respecto a la implementación.

| Código | Etiqueta UI (`VICTIM_CONTACT_STATUS`) | Cuándo se asigna | Uso actual |
|---|---|---|---|
| `'v'` | Válido | Al crear el reporte — `SetVictimContactDefaults` en INSERT | Reporte activo pendiente de gestión. Es el filtro por defecto en pantallas legacy (`fcv` = primer contacto válido) |
| `'i'` | Inválido | Al invalidar — `InvalidateVictimContactByICode` vía controller | Reporte descartado manualmente; se guarda motivo en `victim_contact_status_description`. Listado separado en legacy con filtro `fci` |

**Flujo de vida de un reporte:**

```
INSERT (formulario público)
  → status = 'v'
  → aparece en listado de reportes sin caso (fcv / reportes-component)

Opcional: InvalidateVictimContactByICode
  → status = 'i'
  → deja de aparecer en listado de válidos; pasa al listado fci

Creación de victim_case vinculado al contacto
  → deja de aparecer en reportes-component (sin caso)
  → el status del contacto no cambia automáticamente a 'p' en el código actual
```

> **Decisión para reportes-component:** solo `'v'` sin caso. Los invalidados (`'i'`) quedan fuera de alcance; si se necesitan en el futuro, el padre podría pasar un prop `statusFilter` o usar otra instancia del componente.

---

## Eventos emitidos

| Evento | Cuándo se dispara | Payload |
|---|---|---|
| `action-clicked` | Usuario presiona uno de los botones definidos en `:buttons` | `{ buttonId: string, report: Object }` |

---

## Keys de columna disponibles

| Key | Fuente de dato | Descripción |
|---|---|---|
| `registration_date` | `VictimContact.creationDate` | Fecha y hora en que se registró el reporte |
| `adjustments_gbv` | `VictimContactForm2.adjustmentsGBV[]` → enums | Ajustes razonables solicitados |
| `report_type` | `VictimContactForm2.reportType.name` | Tipo de reporte (víctima o tercero) |
| `will_receive_call` | `VictimContactForm2.willReceiveCall.name` | Disposición a recibir llamada (solo tipo víctima) |
| `victim_names` | `VictimContact.names` | Nombres de la víctima |
| `victim_last_names` | `VictimContact.lastNames` | Apellidos de la víctima |
| `victim_phone` | `VictimContactForm2.victimColPhone` | Teléfono de contacto en Colombia |
| `best_contact_time` | `VictimContactForm2.bestContactTime` | Mejor hora para devolver la llamada |
| `facts_description` | `VictimContactForm2.factsDescription` | Descripción libre de los hechos |

---

## Objeto `report` en cada fila (shape de respuesta backend)

```json
{
  "id": 123,
  "icode": "abc-123-uuid",
  "creationDate": "18/06/2026 14:30:00",
  "names": "María",
  "lastNames": "García López",
  "adjustmentsGBV": [
    { "code": "ni", "name": "No informa" }
  ],
  "reportTypeCode": "v",
  "reportTypeLabel": "Usted es la víctima de violencia basada en género",
  "willReceiveCallLabel": "Sí",
  "victimPhone": "3001234567",
  "bestContactTime": "09:30",
  "factsDescription": "Descripción de los hechos..."
}
```
