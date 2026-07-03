# `remisiones-psicosocial-component` — Inventario de Eventos

Componente: `RemisionesPsicosocialComponent`  
Archivo fuente: `src/frontend/js/components/remisiones-psicosocial-component.js`  
Usado por: Pantallas que necesitan visualizar, filtrar y reasignar remisiones de Atención Psicosocial

> **Nota de alcance:** Este componente es reutilizable. Consulta `salvia.psychosocial_support` con JOINs al caso, formulario de riesgo, dupla y perfiles de agentes. Los filtros se combinan con AND y todos disparan recarga paginada desde el backend. La interacción hacia afuera ocurre mediante tres eventos emitidos: `ver-caso`, `ver-remision` y `reasignar-remisiones`. La lógica de filtrado, paginación y selección masiva es interna al componente.
>
> **Prerequisitos de infraestructura:**
> - [M-01 — Migración schema dupla + psychosocial_support](./Flujos/flow-M01-migracion-schema-dupla-psychosocial-support.md)
> - [M-02 — submitted_by_team, professional_id y team_contact](./Flujos/flow-M02-migracion-submitted-by-team-professional-id-team-contact.md) *(post-reunión)*

Interfaz detallada: [remisiones-psicosocial-component-interface.md](./remisiones-psicosocial-component-interface.md)

---

## Cambios de plan (reunión) — resumen de impacto

| Cambio de schema / UI | Eventos afectados |
|---|---|
| Nuevo campo `submitted_by_team` en `psychosocial_support` | **E-01**, **E-13** |
| Renombrar `agent_id` → `professional_id` | **M-02**, **E-01**, **E-05**, **E-08** |
| Nueva tabla `team_contact` — sesiones con `is_psico_session = true AND is_completed = true` | **M-02**, **E-01**, **E-10** |
| Cards de resumen: **Total** + 4 estados `PsychosocialSupportStatus` | **E-01** |
| Columna REMISIÓN: quitar badge y tags | **E-01** |
| Columna ESTADO Y ASIGNACIÓN: label quemado `"Sesiones (4 - 6)"` + 6 puntos pintados por sesiones completadas | **E-01** |
| **DEC-E08-01:** `filter_professional_id` incluye remisiones vía dupla | **E-01**, **E-05**, **E-08** |

Eventos **sin cambio funcional**: E-02 (índice), E-03, E-04, E-06, E-07, E-09, E-11, E-12, E-14, E-15, E-16, E-17.

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
Requerido:    Sí (base; ver M-02 para ajustes post-reunión)
Estado:       Parcialmente implementado — M-02 corrige agent_id → professional_id
```

---

### M-02 — Migración submitted_by_team, professional_id y team_contact

📄 [Ver flujo → flow-M02-migracion-submitted-by-team-professional-id-team-contact.md](./Flujos/flow-M02-migracion-submitted-by-team-professional-id-team-contact.md)

```
Evento:       Migración de schema — submitted_by_team, professional_id y team_contact
Tipo:         Infrastructure / Backend
Descripción:  Agregar submitted_by_team a psychosocial_support (equipo remitente
              denormalizado). Renombrar agent_id → professional_id. Crear tabla
              salvia.team_contact para registrar sesiones psicosociales. Prerequisito
              de conteo de sesiones en E-01/E-10 y simplificación de E-13.
Requerido:    Sí (prerequisito de E-01 actualizado)
Depende de:   M-01
```

---

## Eventos identificados

---

### E-01 — Cuando carga el componente

📄 [Ver flujo → flow-E01-cuando-carga-componente.md](./Flujos/flow-E01-cuando-carga-componente.md)

```
Evento:       Cuando carga el componente
Tipo:         Lifecycle
Descripción:  Monta el componente, aplica defaultFilter (professional_id o dupla_id)
              si el padre lo envió, carga en paralelo listado paginado, catálogo de
              duplas, equipos remitentes y (si mostrarCards) stats de resumen.
              Renderiza cards encima de filtros cuando mostrarCards === true.
              Renderiza tabla con bloques CASO / REMISIÓN / ESTADO Y ASIGNACIÓN.

              CAMBIOS POST-REUNIÓN:
              • Cards: Total + 4 métricas por status (abierto, en_gestion,
                en_devolucion, cerrado) — 5 cards en fila
              • REMISIÓN: solo remitido por, equipo (submitted_by_team), fecha y links
                — sin badge tipo ni tags extra
              • ESTADO Y ASIGNACIÓN: label fijo "Sesiones (4 - 6)"; 6 puntos; cada
                punto pintado = 1 registro team_contact con is_psico_session=true
                AND is_completed=true para esa remisión (psicosocial_id)
              • submitted_by → JOIN solo para nombre del remitente
              • submitted_by_team → columna directa (sin JOIN para equipo)
              • professional_id reemplaza agent_id en filtros de alcance y asignación
Requerido:    Sí
Prerequisito: M-01 + M-02
```

---

### E-02 — Índice de filtros dropdown

📄 [Ver índice → flow-E02-cuando-selecciona-filtro.md](./Flujos/flow-E02-cuando-selecciona-filtro.md)

| Sub-evento | Filtro | Flujo | Cambio |
|---|---|---|---|
| **E-09** | Estado remisión | [flow-E09](./Flujos/flow-E09-cuando-filtra-estado-remision.md) | — |
| **E-10** | Sesiones completadas | [flow-E10](./Flujos/flow-E10-cuando-filtra-sesiones-completadas.md) | **Fuente: team_contact** |
| **E-11** | Dupla asignada | [flow-E11](./Flujos/flow-E11-cuando-filtra-dupla-asignada.md) | — |
| **E-12** | Nivel de riesgo | [flow-E12](./Flujos/flow-E12-cuando-filtra-nivel-riesgo.md) | — |
| **E-13** | Equipo remitente | [flow-E13](./Flujos/flow-E13-cuando-filtra-equipo-remitente.md) | **Fuente: submitted_by_team** |

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
              paginación. **Preserva defaultFilter** del prop (professional_id
              o dupla_id). Recarga listado y stats con solo el alcance fijo.

              CAMBIO: defaultFilter usa professional_id (antes agent_id).
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
Descripción:  Seleccionar activa `filter_professional_id` y recarga el listado. El backend
              devuelve remisiones donde el profesional está asignado directamente
              (`professional_id`) **o** participa en la dupla de la remisión
              (`dupla.psychologist_id` / `social_worker_id`). Ver DEC-E08-01.
              Limpiar quita el filtro y recarga combinando el resto de filtros activos.

              CAMBIO: filter_professional_id (antes filter_agent_id / agent_id).
              CAMBIO DEC-E08-01: ya no filtra solo `ps.professional_id`.
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
              (Alineado con las 4 cards de resumen en E-01)
Requerido:    Sí
```

---

### E-10 — Cuando filtra por sesiones completadas

📄 [Ver flujo → flow-E10-cuando-filtra-sesiones-completadas.md](./Flujos/flow-E10-cuando-filtra-sesiones-completadas.md)

```
Evento:       Cuando filtra por sesiones completadas
Tipo:         User Interaction
Descripción:  Coincidencia exacta con el conteo de sesiones completadas (0 … 6).
              El conteo proviene de team_contact WHERE psicosocial_id = remisión
              AND is_psico_session = true AND is_completed = true.

              CAMBIO: ya no usa psychosocial_support.session_count (campo inexistente
              o no confiable); fuente única = team_contact.
Requerido:    Sí
Prerequisito: M-02
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
Descripción:  Filtra por psychosocial_support.submitted_by_team (columna directa).
              Catálogo de opciones: DISTINCT submitted_by_team WHERE NOT NULL.
              submitted_by sigue usándose solo para el nombre del remitente (JOIN).

              CAMBIO: ya no resuelve equipo vía JOIN a general_user_team.
Requerido:    Sí
Prerequisito: M-02
```

---

### E-14 — Cuando el usuario selecciona remisiones para reasignación

📄 [Ver flujo → flow-E14-cuando-selecciona-remision-reasignacion.md](./Flujos/flow-E14-cuando-selecciona-remision-reasignacion.md)

```
Evento:       Cuando el usuario selecciona o deselecciona una remisión para reasignación
Tipo:         User Interaction
Descripción:  Solo si :reasignacion === true. Checkbox por fila y "Todos".
              Regla: excluir status cerrado. Cualquier combinación de estados permitida. Solo página actual.
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

- [x] ¿Migración de schema documentada? → M-01, M-02
- [x] ¿Ciclo de vida inicial cubierto? → E-01 (actualizado)
- [x] ¿Filtros cubiertos y combinables (AND)? → E-02 … E-13, E-03, E-04, E-05
- [x] ¿Paginación cubierta? → E-06
- [x] ¿Autocomplete profesional cubierto? → E-07, E-08 (professional_id)
- [x] ¿Selección masiva y emisión al padre? → E-14, E-15
- [x] ¿Navegación emitida al padre? → E-16, E-17
- [x] ¿Cards de resumen documentadas? → E-01 — **Total + 4 cards por PsychosocialSupportStatus**
- [x] ¿Interfaz actualizada (interface.md)? → Sí
- [x] ¿Filtro inicial professional_id / dupla_id? → prop `defaultFilter`, E-01, E-05
- [x] ¿Flujos detallados escritos? → M-01, M-02, E-01 … E-17
- [x] ¿Modal reasignar-remisiones? → [reasignar-remisiones-modal-events.md](./reasignar-remisiones-modal-events.md) (RRM-01 … RRM-05)
- [x] ¿Pantalla padre e integración de rutas? → Historial de Remisiones (sv)
- [ ] ¿Implementación M-02 en código? → Pendiente

---

## Resumen

| # | Evento | Tipo | Persiste en backend | Cambio post-reunión |
|---|---|---|---|---|
| M-01 | Migración schema dupla + psychosocial_support | Infrastructure | Sí (DDL) | Base; ver M-02 |
| M-02 | submitted_by_team, professional_id, team_contact | Infrastructure | Sí (DDL) | **Nuevo** |
| E-01 | Cuando carga el componente | Lifecycle | No (solo lee) | **Cards, REMISIÓN, sesiones** |
| E-02 | Índice de filtros dropdown | — | — | Notas E-10, E-13 |
| E-03 | Cuando escribe número de identidad | User Interaction | No (solo lee) | — |
| E-04 | Cuando escribe teléfono | User Interaction | No (solo lee) | — |
| E-05 | Cuando presiona "Limpiar filtros" | User Interaction | No (solo lee) | **professional_id** |
| E-06 | Cuando cambia de página | User Interaction | No (solo lee) | — |
| E-07 | Cuando escribe autocomplete profesional | User Interaction | No (solo lee) | — |
| E-08 | Cuando selecciona/limpia autocomplete | User Interaction | No (solo lee) | **professional_id** |
| E-09 | Cuando filtra por estado de remisión | User Interaction | No (solo lee) | Alineado a cards |
| E-10 | Cuando filtra por sesiones completadas | User Interaction | No (solo lee) | **team_contact** |
| E-11 | Cuando filtra por dupla asignada | User Interaction | No (solo lee) | — |
| E-12 | Cuando filtra por nivel de riesgo | User Interaction | No (solo lee) | — |
| E-13 | Cuando filtra por equipo remitente | User Interaction | No (solo lee) | **submitted_by_team** |
| E-14 | Cuando selecciona remisión para reasignación | User Interaction | No (estado interno) | — |
| E-15 | Cuando presiona "Reasignar" | User Interaction | No (emite hacia padre) | — |
| E-16 | Cuando presiona "Ver caso" | User Interaction | No (emite hacia padre) | — |
| E-17 | Cuando presiona "Ver remisión" | User Interaction | No (emite hacia padre) | — |

**Total: 19 eventos — 2 Infrastructure, 1 Lifecycle, 14 User Interaction, 1 Índice, 2 condicionales de reasignación**

---

## Decisiones aplicadas

| Tema | Resolución |
|---|---|
| Status en BD | Español snake_case: `abierto`, `en_gestion`, `en_devolucion`, `cerrado` |
| Cards de resumen | **Total remisiones** + **una card por cada status** (Abiertos, En gestión, En devolución, Cerrados) — 5 cards. Ya no: Pendiente asignación / En proceso |
| Autocomplete equipos | `general_user_team IN ('psicologia', 'trab. social')` |
| Barra de sesiones | Label UI quemado **"Sesiones (4 - 6)"**; **6 puntos** fijos; pintar = COUNT `team_contact` con `is_psico_session = true AND is_completed = true` |
| Badge tipo / tags REMISIÓN | **Eliminados** — columna solo muestra remitente, equipo, fecha y acciones |
| `submitted_by` | JOIN solo para **nombre** del remitente |
| `submitted_by_team` | Columna directa en `psychosocial_support` — render y filtro E-13 |
| `professional_id` | Asignación directa en remisión; filtro E-08 también incluye duplas del profesional (DEC-E08-01) |
| `defaultFilter.professional_id` | Alcance fijo “mis remisiones”: directas + vía dupla |
| `defaultFilter.dupla_id` | Alcance fijo a una dupla concreta (opcional; distinto de E-08) |
| `team_contact` | Fuente de verdad para sesiones psicosociales completadas |
| `mostrarCards` | Prop booleano; 5 cards (total + 4 status) encima de filtros |
| Filtro E-08 / `defaultFilter.professional_id` | Incluye asignación directa **OR** duplas del profesional — [DEC-E08-01](./Flujos/flow-decision-E08-filtro-profesional-incluye-dupla.md) |
| Botón reasignación | Label **"Reasignar"**; emite `reasignar-remisiones` al padre |
| Modal reasignación | `reasignar-remisiones-modal` — RRM-01 … RRM-05; padre consume E-15 |

---

## GAPS pendientes

| Tema | Impacto |
|---|---|
| Implementar DEC-E08-01 en backend (`buildPsychosocialListWhere`) | E-08, E-01 defaultFilter, stats |
| Implementar M-02 en código (modelos + AutoMigrate) | Prerequisito backend real |
| Población de `submitted_by_team` al crear remisión | E-01, E-13 — fuera del componente |
| Creación de registros `team_contact` al completar sesiones | E-01 barra de puntos, E-10 |
| Actualizar mock frontend existente | E-01 — alinear con nueva UI |
| Implementación `reasignar-remisiones-modal` + endpoints bulk | E-15 — padre; ver RRM-05 |
