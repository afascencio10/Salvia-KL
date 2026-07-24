# `/salvia/casos/nuevo-v2` — Interfaz de la Pantalla

Versión de "Registro de Caso" que usa `dinamic-form` en vez del formulario estático `set_victim_case.html`. Muestra un formulario dinámico por secciones (9, 128 preguntas) que el operador completa y guarda progresivamente. El `victim_case` se crea desde la primera sección guardada (estado `Borrador`) y se activa (`ra`) solo al completar todas las secciones — ver [`registro-caso-v2-index.md`](registro-caso-v2-index.md) E-03/E-04.

> **Nombre y ruta confirmados** con el usuario: `set_victim_case_v2.html` / `/salvia/casos/nuevo-v2`, siguiendo la convención `follow_up_v2` ya usada en este repo para migraciones a `dinamic-form`.

Archivo fuente: `src/frontend/html/salvia/victim_case/set_victim_case_v2.html`
Componente central: `src/frontend/js/components/dinamic-form.js`
Form ID: `0a24ab30-3cfc-4861-b74d-65d21524bc00` — ver [`form-registro-caso-v2-data.md`](form-registro-caso-v2-data.md)

---

## Archivos relevantes

| Archivo | Rol |
|---|---|
| `src/frontend/html/salvia/victim_case/set_victim_case_v2.html` *(nuevo)* | Pantalla anfitriona — gate de autorización, carga de enums/geografía, cálculo de `formState` |
| `src/frontend/js/components/dinamic-form.js` | Componente reutilizable que renderiza el formulario por secciones |
| `salvia.form` / `form_section` / `question` / `option` / `visibility_condition` / `render_modification` | Estructura completa del formulario (ver `form-registro-caso-v2-data.md`) |
| `salvia.victim_case_form2_enums` | Catálogo de opciones reutilizado vía `state_options_path` |
| `security.department` / `security.city` / `security.town` | Geografía para las 3 cadenas encadenadas |

---

## Árbol de interfaz

```
Page  (.container.mt-4)
│
├── [v-if !authorizationAnswer || authorizationAnswer.code != 'y']
│   SeccionAutorizacion
│   ├── Titulo  "Autorización de datos personales"
│   ├── TextoLegal  (Ley 1581 de 2012, derechos ARCO)
│   └── SelectSingle  authorizationAnswer  (yes_no)
│
└── [v-else]  #registro-caso-v2-app  (Vue mount)
    │
    ├── [v-if loadError]
    │   ErrorAlert  (.alert.alert-danger)  loadError
    │
    └── [v-else]
        DinamicForm  (form-id="0a24ab30-..." :form-state="formState" @answers-updated @form-completed)
        │
        ├── [v-if loading]      LoadingSpinner
        ├── [v-else-if error]   ErrorPanel
        └── [v-else]            .df-wrapper  (idéntico a hacer-seguimiento — ver dinamic-form-interface.md)
            ├── Sidebar (9 secciones)
            └── MainPanel
                ├── Sección 1 — Datos de la Víctima  (incl. banner "Accesibilidad")
                ├── Sección 2 — Contacto de Apoyo
                ├── Sección 3 — Hechos  (incl. subtipos dinámicos, sector laboral condicional)
                ├── Sección 4 — Agresor
                ├── Sección 5 — Tamizaje  (bloque pareja/no-pareja condicional + banner de riesgo)
                ├── Sección 6 — Datos Personales  (9 preguntas de escala + condicionales)
                ├── Sección 7 — Plan de Acción
                ├── Sección 8 — Denuncia Fácil
                └── Sección 9 — Lugar de Atención  (solo cadena geográfica — SIN grilla de ruta)
        │
        └── [v-if formCompleted]
            CompletedOverlay  (análogo a hs-completed-backdrop)
            ├── Icon  SVG checkmark verde
            ├── Titulo  "Caso registrado exitosamente"
            ├── Credenciales
            │   ├── Usuario:  newUser.login
            │   └── Clave:    newUser.pass
            └── BtnFinalizar  → goToCase()  [redirige a /salvia/casos/:caseId/detalle]
```

> El árbol interno de cada sección (tipos de input, chips, dropdowns, errores inline) es **idéntico** al de `dinamic-form` documentado en `dinamic-form-interface.md` — no se repite aquí.

---

## Diferencias clave vs. la pantalla estática actual

| Aspecto | `set_victim_case.html` (actual) | `set_victim_case_v2.html` (nuevo) |
|---|---|---|
| Render | Template Vue estático, ~140 campos hardcoded | `dinamic-form`, 128 preguntas dirigidas por BD |
| Guardado | Atómico — 1 POST al final crea todo | Progresivo — `saveSection` por sección; el caso se crea en `Borrador` desde la 1ª sección (E-03) y se activa al completar (E-04) |
| Enums | Inyectados server-side en el `data()` inicial | Cargados por el host y expuestos vía `formState.enums.*` + `state_options_path` |
| Cadenas geográficas | 3 lógicas independientes en el Vue del host | 3 cadenas vía `state_options_path` (`geo.<chain>.cities/towns`), mismo patrón que Barreras en Seguimiento |
| Tamizaje | Calculado inline en el template al vuelo | Calculado por el host en `answers-updated`, expuesto en `formState` (banner + visibilidad) |
| Asignación de Ruta | Grid de sedes por momento/sector | **Excluida de esta fase** — ver GAPS en `form-registro-caso-v2-data.md` |
| GPS | Capturado en `mounted()`, anexado antes del POST | Sigue siendo responsabilidad del host — no es una pregunta del form |
| Credenciales nuevo usuario | Vienen en la respuesta del único POST | Se generan al activar el caso (E-04); el host hace un GET puntual tras `form-completed` (E-05) para recuperarlas — ver GAPS sobre la tabla de resultado temporal |

---

## GAPS

| # | Descripción |
|---|---|
| G-07 | Si viene de un contacto previo (`:id` como en la pantalla actual), cómo se prellenan las 128 respuestas en el `form_submission` antes de que el usuario abra el formulario |
| G-11 | Agregar el código `bo` a `VICTIM_CASE_STATUS` y auditar los listados/KPIs que no filtran explícitamente por status para excluir casos en Borrador (ver `form-registro-caso-v2-data.md`) |
