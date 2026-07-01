# `remisiones-psicosocial-component` — Inventario de Eventos

Componente: `RemisionesPsicosocialComponent`  
Archivo fuente: `src/frontend/components/remisiones-psicosocial-component.js` *(pendiente de crear)*  
Usado por: Pantallas que necesitan visualizar, filtrar y reasignar remisiones de Atención Psicosocial

> **Nota de alcance:** Este componente es reutilizable. Consulta `salvia.psychosocial_support` con JOINs al caso, formulario de riesgo, dupla y perfiles de agentes. Los filtros se combinan con AND y todos disparan recarga paginada desde el backend. La interacción hacia afuera ocurre mediante tres eventos emitidos: `ver-caso`, `ver-remision` y `reasignar-remisiones`. La lógica de filtrado, paginación y selección masiva es interna al componente.
>
> **Prerequisito de infraestructura:** [M-01 — Migración schema dupla + psychosocial_support](./Flujos/flow-M01-migracion-schema-dupla-psychosocial-support.md)

Interfaz detallada: [remisiones-psicosocial-component-interface.md](./remisiones-psicosocial-component-interface.md)

---

## Eventos de infraestructura

---

### M-01 — Migración de schema (dupla + psychosocial_support)

📄 [Ver flujo → flow-M01-migracion-schema-dupla-psychosocial-support.md](./Flujos/flow-M01-migracion-schema-dupla-psychosocial-support.md)

```
Evento:       Migración de schema — tabla dupla y ajuste psychosocial_support
Tipo:         Infrastructure / Backend
Descripción:  Crear modelo y tabla salvia.dupla. Agregar submitted_by, dupla_id
              y agent_id a psychosocial_support. Definir constantes de status
              en español snake_case (patrón entity_letter.go). Cambiar default
              de status de 'ACTIVE' a 'abierto'. Registrar AutoMigrate en main.go.
              Migrar registros legacy ACTIVE → abierto.
Requerido:    Sí (prerequisito de E-01)
```

---

## Eventos identificados

---

### E-01 — Cuando carga el componente

📄 [Ver flujo → flow-E01-cuando-carga-componente.md](./Flujos/flow-E01-cuando-carga-componente.md)

```
Evento:       Cuando carga el componente
Tipo:         Lifecycle
Descripción:  Monta el componente, aplica defaultFilter (agent_id o dupla_id) si
              el padre lo envió, carga en paralelo listado paginado, catálogo de
              duplas, equipos remitentes y (si mostrarCards) stats de resumen.
              Renderiza cards encima de filtros cuando mostrarCards === true.
              Renderiza tabla con bloques CASO / REMISIÓN / ESTADO Y ASIGNACIÓN.
              submitted_by: solo lectura — JOIN para nombre y team del remitente.
Requerido:    Sí
```

---

### E-02 — Índice de filtros dropdown

📄 [Ver índice → flow-E02-cuando-selecciona-filtro.md](./Flujos/flow-E02-cuando-selecciona-filtro.md)

| Sub-evento | Filtro | Flujo |
|---|---|---|
| **E-09** | Estado remisión | [flow-E09](./Flujos/flow-E09-cuando-filtra-estado-remision.md) |
| **E-10** | Sesiones completadas | [flow-E10](./Flujos/flow-E10-cuando-filtra-sesiones-completadas.md) |
| **E-11** | Dupla asignada | [flow-E11](./Flujos/flow-E11-cuando-filtra-dupla-asignada.md) |
| **E-12** | Nivel de riesgo | [flow-E12](./Flujos/flow-E12-cuando-filtra-nivel-riesgo.md) |
| **E-13** | Equipo remitente | [flow-E13](./Flujos/flow-E13-cuando-filtra-equipo-remitente.md) |

---

### E-03 — Cuando el usuario escribe en el filtro de número de identidad

📄 [Ver flujo → flow-E03-cuando-escribe-numero-identidad.md](./Flujos/flow-E03-cuando-escribe-numero-identidad.md)

```
Evento:       Cuando el usuario escribe en el filtro de número de identidad
Tipo:         User Interaction
Descripción:  Debounce 400ms. Filtra por victim_case_victim_doc_number (ILIKE).
              Combinable con todos los demás filtros activos.
Requerido:    Sí
```

---

### E-04 — Cuando el usuario escribe en el filtro de teléfono

📄 [Ver flujo → flow-E04-cuando-escribe-telefono.md](./Flujos/flow-E04-cuando-escribe-telefono.md)

```
Evento:       Cuando el usuario escribe en el filtro de teléfono
Tipo:         User Interaction
Descripción:  Debounce 400ms. Filtra por teléfono de la víctima (COALESCE form2
              / contact form1). Combinable con todos los demás filtros activos.
Requerido:    Sí
```

---

### E-05 — Cuando el usuario presiona "Limpiar filtros"

📄 [Ver flujo → flow-E05-cuando-limpia-filtros.md](./Flujos/flow-E05-cuando-limpia-filtros.md)

```
Evento:       Cuando el usuario presiona "Limpiar filtros"
Tipo:         User Interaction
Descripción:  Resetea filtros de UI (búsqueda, dropdowns, autocomplete) y
              paginación. **Preserva defaultFilter** del prop (agent_id o
              dupla_id). Recarga listado y stats con solo el alcance fijo.
Requerido:    Sí
```

---

### E-06 — Cuando el usuario cambia de página

📄 [Ver flujo → flow-E06-cuando-cambia-pagina.md](./Flujos/flow-E06-cuando-cambia-pagina.md)

```
Evento:       Cuando el usuario cambia de página
Tipo:         User Interaction
Descripción:  PaginationBar. Mantiene filtros activos, limpia selectedRemisiones.
Requerido:    Sí
```

---

### E-07 — Cuando el usuario escribe en el autocomplete de profesional asignada

📄 [Ver flujo → flow-E07-cuando-escribe-autocomplete.md](./Flujos/flow-E07-cuando-escribe-autocomplete.md)

```
Evento:       Cuando el usuario escribe en el autocomplete de profesional asignada
Tipo:         User Interaction
Descripción:  Busca agentes activos con general_user_team IN ('psicologia',
              'trab. social'). Mínimo 3 caracteres, debounce 400ms. No recarga tabla.
Requerido:    Sí
```

---

### E-08 — Cuando el usuario selecciona o limpia el autocomplete

📄 [Ver flujo → flow-E08-cuando-selecciona-autocomplete.md](./Flujos/flow-E08-cuando-selecciona-autocomplete.md)

```
Evento:       Cuando el usuario selecciona o limpia el autocomplete
Tipo:         User Interaction
Descripción:  Seleccionar filtra por psychosocial_support.agent_id. Limpiar quita
              el filtro y recarga combinando el resto de filtros activos.
Requerido:    Sí
```

---

### E-09 — Cuando filtra por estado de remisión

📄 [Ver flujo → flow-E09-cuando-filtra-estado-remision.md](./Flujos/flow-E09-cuando-filtra-estado-remision.md)

```
Evento:       Cuando filtra por estado de remisión
Tipo:         User Interaction
Descripción:  Filtra por psychosocial_support.status. Valores en español snake_case:
              abierto | en_gestion | en_devolucion | cerrado
Requerido:    Sí
```

---

### E-10 — Cuando filtra por sesiones completadas

📄 [Ver flujo → flow-E10-cuando-filtra-sesiones-completadas.md](./Flujos/flow-E10-cuando-filtra-sesiones-completadas.md)

```
Evento:       Cuando filtra por sesiones completadas
Tipo:         User Interaction
Descripción:  Coincidencia exacta con session_count (0 … 6).
Requerido:    Sí
```

---

### E-11 — Cuando filtra por dupla asignada

📄 [Ver flujo → flow-E11-cuando-filtra-dupla-asignada.md](./Flujos/flow-E11-cuando-filtra-dupla-asignada.md)

```
Evento:       Cuando filtra por dupla asignada
Tipo:         User Interaction
Descripción:  Opciones de salvia.dupla (name → id). Filtra dupla_id.
Requerido:    Sí
```

---

### E-12 — Cuando filtra por nivel de riesgo

📄 [Ver flujo → flow-E12-cuando-filtra-nivel-riesgo.md](./Flujos/flow-E12-cuando-filtra-nivel-riesgo.md)

```
Evento:       Cuando filtra por nivel de riesgo
Tipo:         User Interaction
Descripción:  Filtra por victim_case_form2_risk_level del caso (1-4).
Requerido:    Sí
```

---

### E-13 — Cuando filtra por equipo remitente

📄 [Ver flujo → flow-E13-cuando-filtra-equipo-remitente.md](./Flujos/flow-E13-cuando-filtra-equipo-remitente.md)

```
Evento:       Cuando filtra por equipo remitente
Tipo:         User Interaction
Descripción:  Filtra por general_user_team del usuario en submitted_by. El
              componente solo lee submitted_by y resuelve nombre/equipo vía JOIN.
Requerido:    Sí
```

---

### E-14 — Cuando el usuario selecciona remisiones para reasignación

📄 [Ver flujo → flow-E14-cuando-selecciona-remision-reasignacion.md](./Flujos/flow-E14-cuando-selecciona-remision-reasignacion.md)

```
Evento:       Cuando el usuario selecciona o deselecciona una remisión para reasignación
Tipo:         User Interaction
Descripción:  Solo si :reasignacion === true. Checkbox por fila y "Todos".
              Regla: mismo status en todas las seleccionadas. Solo página actual.
Requerido:    Condicional
```

---

### E-15 — Cuando el usuario presiona "Reasignar"

📄 [Ver flujo → flow-E15-cuando-presiona-reasignar.md](./Flujos/flow-E15-cuando-presiona-reasignar.md)

```
Evento:       Cuando el usuario presiona "Reasignar"
Tipo:         User Interaction
Descripción:  Emite 'reasignar-remisiones' al padre con selectedRemisiones.
              No abre modal ni persiste — responsabilidad del padre.
Requerido:    Condicional
```

---

### E-16 — Cuando el usuario presiona "Ver caso"

📄 [Ver flujo → flow-E16-cuando-presiona-ver-caso.md](./Flujos/flow-E16-cuando-presiona-ver-caso.md)

```
Evento:       Cuando el usuario presiona "Ver caso"
Tipo:         User Interaction
Descripción:  Emite 'ver-caso' con caseICode y objeto remisión.
Requerido:    Sí
```

---

### E-17 — Cuando el usuario presiona "Ver remisión"

📄 [Ver flujo → flow-E17-cuando-presiona-ver-remision.md](./Flujos/flow-E17-cuando-presiona-ver-remision.md)

```
Evento:       Cuando el usuario presiona "Ver remisión"
Tipo:         User Interaction
Descripción:  Emite 'ver-remision' con remisionId, followUpId y objeto remisión.
Requerido:    Sí
```

---

## Checklist de completitud

- [x] ¿Migración de schema documentada? → M-01
- [x] ¿Ciclo de vida inicial cubierto? → E-01
- [x] ¿Filtros cubiertos y combinables (AND)? → E-02 … E-13, E-03, E-04, E-05
- [x] ¿Paginación cubierta? → E-06
- [x] ¿Autocomplete profesional cubierto? → E-07, E-08
- [x] ¿Selección masiva y emisión al padre? → E-14, E-15
- [x] ¿Navegación emitida al padre? → E-16, E-17
- [x] ¿Cards de resumen documentadas? → E-01, prop `mostrarCards`
- [x] ¿Filtro inicial agent_id / dupla_id? → prop `defaultFilter`, E-01, E-05
- [x] ¿Flujos detallados escritos? → M-01, E-01 … E-17
- [ ] ¿Modal reasignar-remisiones? → Planeación futura (padre consume E-15)
- [ ] ¿Pantalla padre e integración de rutas? → Pendiente

---

## Resumen

| # | Evento | Tipo | Persiste en backend |
|---|---|---|---|
| M-01 | Migración schema dupla + psychosocial_support | Infrastructure | Sí (DDL) |
| E-01 | Cuando carga el componente | Lifecycle | No (solo lee) |
| E-02 | Índice de filtros dropdown | — | — |
| E-03 | Cuando escribe número de identidad | User Interaction | No (solo lee) |
| E-04 | Cuando escribe teléfono | User Interaction | No (solo lee) |
| E-05 | Cuando presiona "Limpiar filtros" | User Interaction | No (solo lee) |
| E-06 | Cuando cambia de página | User Interaction | No (solo lee) |
| E-07 | Cuando escribe autocomplete profesional | User Interaction | No (solo lee) |
| E-08 | Cuando selecciona/limpia autocomplete | User Interaction | No (solo lee) |
| E-09 | Cuando filtra por estado de remisión | User Interaction | No (solo lee) |
| E-10 | Cuando filtra por sesiones completadas | User Interaction | No (solo lee) |
| E-11 | Cuando filtra por dupla asignada | User Interaction | No (solo lee) |
| E-12 | Cuando filtra por nivel de riesgo | User Interaction | No (solo lee) |
| E-13 | Cuando filtra por equipo remitente | User Interaction | No (solo lee) |
| E-14 | Cuando selecciona remisión para reasignación | User Interaction | No (estado interno) |
| E-15 | Cuando presiona "Reasignar" | User Interaction | No (emite hacia padre) |
| E-16 | Cuando presiona "Ver caso" | User Interaction | No (emite hacia padre) |
| E-17 | Cuando presiona "Ver remisión" | User Interaction | No (emite hacia padre) |

**Total: 18 eventos — 1 Infrastructure, 1 Lifecycle, 14 User Interaction, 1 Índice, 2 condicionales de reasignación**

---

## Decisiones aplicadas

| Tema | Resolución |
|---|---|
| Status en BD | Español snake_case: `abierto`, `en_gestion`, `en_devolucion`, `cerrado` (patrón `entity_letter.go`) |
| Autocomplete equipos | `general_user_team IN ('psicologia', 'trab. social')` |
| Máximo sesiones | Constante quemada `MAX_SESSIONS = 6` en frontend |
| Badge tipo / tags extra | Texto **"Por consultar"** en UI |
| `submitted_by` | Solo lectura: leer id guardado → JOIN nombre + team del remitente |
| `mostrarCards` | Prop booleano; cards de resumen encima de filtros |
| `defaultFilter` | `{ agent_id }` o `{ dupla_id }` — alcance fijo, persiste al limpiar filtros |
| Botón reasignación | Label **"Reasignar"**; emite `reasignar-remisiones` al padre |
| Modal reasignación | Fuera de alcance — planeación futura |

---

## GAPS pendientes

| Tema | Impacto |
|---|---|
| Origen real del máximo de sesiones (hoy quemado en 6) | E-01, barra de progreso |
| Fuente de badge tipo remisión y tags adicionales | Bloque REMISIÓN — hoy "Por consultar" |
| Planeación `reasignar-remisiones-modal` | E-15 — padre |
| Pantalla padre y rutas E-16 / E-17 | Integración |
