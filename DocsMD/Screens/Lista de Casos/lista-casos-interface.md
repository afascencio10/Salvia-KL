# Lista de casos — Interfaz

**Ruta:** `GET /salvia/lista-casos`  
**Template:** `src/frontend/html/salvia/list-cases/list_cases.html`  
**Facade:** `ListCasesGET` en `src/salvia/facades/VictimCaseFacade.go`

## Archivos relevantes

| Archivo | Rol |
|---|---|
| `src/frontend/html/salvia/list-cases/list_cases.html` | Pantalla padre |
| `src/frontend/js/components/casos-component.js` | Tabla y filtros |
| `src/frontend/html/salvia/list-cases/reasignar_casos_modal.html` | Modal reasignación (`sv`) |
| `src/frontend/html/salvia/list-cases/consolidated_report_modal.html` | Reporte consolidado (`sv`/`ad`) |

## Árbol de interfaz

```
list_cases  (.salvia-main > .container)
│
├── WelcomeSection  (.welcome-section)
│   ├── TitleWrap
│   │   ├── h2  "Lista de casos"
│   │   └── [v-if isEnDepartmentScope]
│   │       Subtitle  (.text-muted)  "Casos del departamento asignado al usuario"
│   └── Actions
│       ├── [v-if canDescargarReporte]  Btn  "Descargar reporte consolidado"
│       └── Link  /salvia/casos/nuevo  Btn  "Crear Nuevo caso"
│
├── [v-if enMissingDepartment]
│   Alert  (.alert.alert-danger)
│   └── "Tu usuario no tiene un departamento asignado..."
│
└── [v-else]
    ├── casos-component
    ├── [v-if canReasignarCasos]  reasignar-casos-modal
    └── [v-if canDescargarReporte]  consolidated-report-modal
```

`enMissingDepartment` = `currentRole === 'en'` y `assignedDepartmentId` vacío.  
`isEnDepartmentScope` = `currentRole === 'en'` y hay `assignedDepartmentId`.  
`canReasignarCasos` = rol `sv`.  
`canDescargarReporte` = rol `sv` o `ad`.
