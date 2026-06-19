# Listado de Reportes — Inventario de Eventos

Pantalla: `GET /salvia/casos` (roles `sv`, `ro`)  
Templates: `get_victim_cases_sv.html`, `get_victim_cases_ro.html`  
API de listado paginado: `GET /salvia/primer_contacto/f/{fcv|fci}/p/{page}`

---

## Eventos identificados

### E-01 — Cuando carga la pantalla

📄 [Ver flujo → flow-E01-cuando-carga-pantalla.md](./Flujos/flow-E01-cuando-carga-pantalla.md)

```
Evento:       Cuando carga la pantalla
Tipo:         Lifecycle
Descripción:  El facade renderiza el template según rol (sv/ro).
              Carga por defecto reportes válidos sin caso (filter fcv) vía
              GetVictimContactsWithoutVictimCase(page, "v"). Inicializa
              victimContacts, paginación y menuToolsContactTable (Ver + Crear caso).
Requerido:    Sí
```

---

### E-02 — Cuando filtra reportes por nombre, apellido o teléfono

📄 [Ver flujo → flow-E02-cuando-filtra-reportes.md](./Flujos/flow-E02-cuando-filtra-reportes.md)

```
Evento:       Cuando filtra reportes
Tipo:         User Interaction
Descripción:  Usuario escribe en uno o más inputs (nombre, apellido, teléfono)
              y presiona "Buscar", o presiona "Limpiar" para resetear.
              Resetea currentPage a 0 y llama al backend con query params
              names, lastNames, phone combinados con AND.
Requerido:    Sí
```

---

### E-03 — Cuando cambia de tab o de página

📄 [Ver flujo → flow-E03-cuando-cambia-tab-o-pagina.md](./Flujos/flow-E03-cuando-cambia-tab-o-pagina.md)

```
Evento:       Cuando cambia de tab o de página
Tipo:         User Interaction
Descripción:  Usuario selecciona tab Recontacto (fcv) o Inválidos (fci),
              o navega paginación Anterior/Siguiente. Recarga victimContacts
              manteniendo filtros activos si los hay.
Requerido:    Sí
```

---

### E-04 — Cuando presiona Crear caso o Ver

📄 [Ver flujo → flow-E04-cuando-presiona-crear-caso.md](./Flujos/flow-E04-cuando-presiona-crear-caso.md)

```
Evento:       Cuando presiona una acción de fila
Tipo:         User Interaction
Descripción:  Enlaces de menuToolsContactTable:
              (A) Ver → detalle del reporte
              (B) Crear caso → formulario set_victim_case pre-cargado
Requerido:    Sí
```

---

## Eventos comentados (legacy — no aplican en solo-reportes)

| Evento legacy | Ubicación | Acción |
|---|---|---|
| `searchCase()` | HTML + script | Comentar — busca casos por documento |
| `goToTab(routing\|…)` | script | Comentar cases del switch |
| Paginación `victimCases` | `goToPage()` rama casos | Comentar rama |
| `assignOperator()` | sv only | Comentar overlay y método |
| `submit('openAssign'\|'assign')` | sv only | Comentar |

---

## Checklist de completitud

- [x] Carga inicial → E-01
- [x] Filtros nombre/apellido/teléfono → E-02
- [x] Tabs y paginación → E-03
- [x] Acciones Ver / Crear caso → E-04
- [x] Backend para filtros → E-02 (DAO)

---

## Resumen

| # | Evento | Tipo | Persiste en backend |
|---|---|---|---|
| E-01 | Cuando carga la pantalla | Lifecycle | No (solo lee) |
| E-02 | Cuando filtra reportes | User Interaction | No (solo lee) |
| E-03 | Cuando cambia tab o página | User Interaction | No (solo lee) |
| E-04 | Cuando presiona Crear caso o Ver | User Interaction | No (navegación) |
