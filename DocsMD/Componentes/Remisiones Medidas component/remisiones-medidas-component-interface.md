# `remisiones-medidas-component` — Interfaz del Componente

Tabla reutilizable de remisiones de Medidas de Emergencia. La prop `vista` decide el recorte, las cards o el banner, los filtros y el botón de cada fila. El padre escucha `ver-caso` y `ver-detalle`.

La pantalla ya renderiza con data quemada en el componente. El listado, los catálogos y los dos guardados que describen estos documentos son el contrato al que se va a cablear. Ese contrato todavía no está implementado.

**Prerequisito:** [M-01 — Cabecera de remisión y `referral_id`](./Flujos/flow-M01-migracion-emergency-measure-referral.md)

Eventos: [remisiones-medidas-component-events.md](./remisiones-medidas-component-events.md)

---

## Archivos relevantes

| Archivo | Rol |
|---|---|
| `src/frontend/js/components/remisiones-medidas-component.js` | Componente. Hoy clona un seed de 15 remisiones |
| `src/frontend/html/salvia/remisiones-medidas/remisiones_medidas_component.html` | Template Vue (`#tpl-remisiones-medidas-component`) |
| `src/frontend/css/remisiones-medidas-component.css` | Estilos, prefijo `rme-` |
| `src/internal/models/emergency_measure.go` | Medida por tipo. M-01 le agrega `referral_id` |
| `src/salvia/facades/RemisionesMedidasFacade.go` | Pantalla general, `vista="supervisor"` |
| `src/salvia/facades/MisRemisionesMedidasFacade.go` | Pantalla del agente, `vista="agente"` |

Los dos modales (reasignar y revisión técnica) viven en este mismo componente. No hay un segundo componente de modal.

---

## Props

| Prop | Tipo | Requerido | Default | Descripción |
|---|---|---|---|---|
| `vista` | `String` | No | `supervisor` | Recorte y presentación. Valores: `supervisor`, `equipo-tecnico`, `supervision-convenio`, `cruz-roja`, `agente` |
| `agenteId` | `String` | No | `am-1` | Solo lo usa `vista="agente"`. Es el `general_user_i_code` del agente. Hoy Mis Remisiones envía `am-1` |

El mapeo `currentRole` → `vista` queda pendiente. Estas dos pantallas fijan la prop. Las vistas `equipo-tecnico`, `supervision-convenio` y `cruz-roja` están implementadas en el componente y ninguna pantalla las monta todavía.

```html
<!-- Remisiones Medidas de Emergencia -->
<remisiones-medidas-component
  vista="supervisor"
  @ver-caso="onVerCaso"
  @ver-detalle="onVerDetalle"
></remisiones-medidas-component>

<!-- Mis Remisiones — Medidas de Emergencia -->
<remisiones-medidas-component
  vista="agente"
  agente-id="am-1"
  @ver-caso="onVerCaso"
  @ver-detalle="onVerDetalle"
></remisiones-medidas-component>
```

`agenteId` no se limpia con "Limpiar". Forma parte del recorte de la vista, igual que `defaultFilter` en el componente de Atención Psicosocial.

---

## Emits

| Emit | Payload | Cuándo |
|---|---|---|
| `ver-caso` | `{ caseICode, remisionId }` | Clic en "Ver caso →" (E-14) |
| `ver-detalle` | `{ remisionId, caseICode }` | Clic en "Ver detalle →" o en un botón de fila cuyo tipo es detalle (E-15) |

El componente no navega. Los padres actuales dejan los dos handlers vacíos: el detalle de la medida y los casos del seed no existen.

Reasignar y guardar la revisión no se emiten. El componente llama al API y vuelve a cargar el listado (E-10, E-13).

---

## Qué muestra cada vista

| `vista` | Encabezado | Filtro estado | Filtro agente | Checkbox |
|---|---|---|---|---|
| `supervisor` | 8 cards | Los 16 estados | Sí, incluye "Sin asignar" | Sí |
| `agente` | 6 cards | Los 5 estados de su bandeja | No | No |
| `equipo-tecnico` | Banner violeta | No | No | No |
| `supervision-convenio` | Banner fucsia | No | No | No |
| `cruz-roja` | Banner verde | No | No | No |

Búsqueda, urgencia y el botón de limpiar están en las cinco vistas.

Placeholder de búsqueda: "Nombre o número de identidad" en agente; "Nombre o doc. identidad" en el resto. Texto del botón: "Limpiar" en agente; "Limpiar filtros" en el resto.

---

## Cards

Se pintan encima de los filtros solo en `supervisor` y `agente`. El valor sale del recorte de la vista (`alcance`), antes de búsqueda, estado, agente y urgencia. Filtrar la tabla no mueve las cards.

### Supervisor (8)

| Key | Label | Cuenta en el alcance |
|---|---|---|
| `total` | Total | Todas |
| `riesgo` | Riesgo inminente | Los cuatro estados `riesgo_inminente_*` |
| `revision` | Pendiente de Revisión | `pendiente_de_revision` |
| `porAsignar` | Por asignar | `por_asignar` |
| `asignadas` | Asignadas | `asignada` |
| `gestion` | En gestión documental | `en_gestion_documental` |
| `caracterizado` | Caracterizado | `caracterizado` |
| `aprobadas` | Aprobadas | `aprobadas` |

### Agente (6)

| Key | Label | Cuenta en el alcance |
|---|---|---|
| `total` | Mis remisiones | Todas las de su recorte |
| `riesgo` | Riesgo inminente | Solo `riesgo_inminente_por_confirmar_seguridad` |
| `asignadas` | Asignadas | `asignada` |
| `gestion` | En gestión documental | `en_gestion_documental` |
| `aprobadas` | Aprobadas | `aprobadas` |
| `seguimiento` | En Seguimiento | `en_seguimiento` |

Las otras tres vistas no llevan cards. Llevan un banner:

| Vista | Tono | Texto |
|---|---|---|
| `equipo-tecnico` | `equipo` | Bandeja "Pendiente de Revisión": casos remitidos a Medidas de Emergencia a la espera del concepto cualificado del Equipo Técnico. |
| `supervision-convenio` | `convenio` | Mostrando casos "Por revisar trazabilidad", "Por aprobar cotizaciones" y "Riesgo inminente: por escalar a Cruz Roja" (prioritario) a cargo de Supervisión del Convenio. |
| `cruz-roja` | `cruz-roja` | Mostrando casos "Por cotizar", "Riesgo inminente: por entregar medida" (prioritario) y "Aprobadas" pendientes de otorgamiento. |

---

## Árbol de interfaz

```
RemisionesMedidasComponent  (.rme-wrapper)
├── ReassignBar                         v-if supervisor y hay selección
│   └── Button  "↺ Reasignar N seleccionada(s)"     → E-08
├── Cards                               v-if supervisor o agente
│   └── Card × 8 o × 6                  label + valor
├── Banner                              v-if equipo-tecnico, convenio o cruz-roja
├── Filters
│   ├── Input búsqueda                  → E-02
│   ├── Select estado                   → E-03, solo supervisor y agente
│   ├── Select agente                   → E-04, solo supervisor
│   ├── Select urgencia                 → E-05
│   └── Button limpiar                  → E-06
├── Table
│   ├── Head                            checkbox solo supervisor
│   ├── Empty                           ver mensaje vacío
│   └── Row × N
│       ├── Checkbox                    → E-07, solo supervisor
│       ├── Columna Caso                "Ver caso →" → E-14
│       ├── Columna Remisión o Medidas  "Ver detalle →" → E-15
│       └── Columna Estado              botón de acción, agente o concepto
└── Footer  "{visibles} de {alcance} remisiones"

Modal Reasignar                         E-08 abre, E-09 cancela, E-10 confirma
Modal Revisión Equipo Técnico          E-11 abre, E-12 cancela, E-13 confirma
```

No hay paginación.

---

## Columnas

La grilla del agente tiene tres columnas: Caso, Medidas, Estado. La del resto: Caso, Remisión, Estado y asignación. El supervisor suma la columna del checkbox.

### Caso

| Vista agente | Resto de vistas |
|---|---|
| Avatar con la inicial, nombre, documento | Nombre y badge de riesgo en la misma línea |
| Riesgo y municipio en la misma línea | `{iCode} · {documento}` y municipio debajo |
| "Ver caso →" | "Ver caso →" |

El badge de riesgo usa el texto ya resuelto del caso: `victim_case_form2_risk_level` con el mismo mapa que el detalle del caso (`1` Bajo, `2` Moderado, `3` Alto, `4` Extremo; si no hay form 2 y el form 1 marca riesgo de feminicidio, "Alto (feminicidio)"). El seed pinta "Crítico" y "Medio"; al cablear el API se muestran las etiquetas reales. Municipio: `town_name` del municipio de la víctima, el mismo join del listado de Atención Psicosocial.

### Remisión (gestión) o Medidas (agente)

Chips con cada `emergency_measure.type` de la remisión, badge de urgencia ("Riesgo inminente" o "Normal") y "Ver detalle →".

En gestión, debajo de los chips: "Por: {nombre del remitente}" y la fecha `referred_at`. En agente: "desde {fecha}" y no se muestra el remitente.

### Estado y asignación (gestión) o Estado (agente)

Badge del estado. Si `accionDe` devuelve un botón, se pinta ese botón y no el agente. Si no hay botón y la vista no es agente: avatar y nombre del agente, o el texto "Sin agente asignado". Si hay revisión técnica, debajo se muestra "Concepto técnico ( {fecha} ): {resultado}" y "Por: {profesional}". El texto del concepto no se pinta en la fila.

La columna del agente no muestra agente, "Sin agente asignado" ni concepto.

---

## Botón de acción de la fila

Una fila tiene como máximo un botón. La regla está en `accionDe` y queda partida en tres eventos:

| Tipo | Evento | Qué hace |
|---|---|---|
| `asignar` | E-08 | Abre el modal de reasignar con esa sola fila |
| `revisar` | E-11 | Abre el modal de revisión técnica |
| `detalle` | E-15 | Emite `ver-detalle` |

Etiquetas de tipo `detalle` (el negocio de cada una queda en la pantalla de detalle):

| Vista | Estado | Etiqueta |
|---|---|---|
| `supervisor` | `riesgo_inminente_por_decidir_via` | Decidir vía de atención |
| `supervision-convenio` | `riesgo_inminente_por_escalar_cruz_roja` | Escalar a Cruz Roja |
| `cruz-roja` | `riesgo_inminente_por_entregar_medida` | Entregar medida inmediata |
| `supervisor` | `caracterizado` | Revisar caracterización |
| `supervision-convenio` | `por_revisar_trazabilidad` | Revisar trazabilidad |
| `supervision-convenio` | `por_aprobar_cotizaciones` | Aprobar cotizaciones |
| `cruz-roja` | `por_cotizar` | Cargar cotizaciones |
| `agente` | `riesgo_inminente_por_confirmar_seguridad` | Confirmar seguridad |

"Revisar (Equipo Técnico)" aparece en `supervisor` y `equipo-tecnico` cuando el estado es `pendiente_de_revision`. "Asignar agente" aparece solo en `supervisor` cuando el estado es `por_asignar`.

---

## Modal Reasignar remisiones

Título "Reasignar remisiones". Texto con la cantidad seleccionada. Select "Agente de Medidas" cargado con `GET .../agentes`. "Reasignar" queda deshabilitado hasta elegir un agente. Cancelar, o confirmar (E-10).

Abrirlo desde "Asignar agente" reemplaza la selección por esa única fila. Abrirlo desde la barra usa la selección que ya había.

---

## Modal Revisión Equipo Técnico

Muestra el nombre y el `i_code` del caso. Select de profesional (`GET .../profesionales-revision`), textarea de concepto y tres resultados:

| Valor que viaja | Label |
|---|---|
| `aprobado_para_caracterizacion_me` | Aprobado para caracterización ME |
| `devuelto_para_ajuste` | Devuelto para ajuste |
| `devuelto_para_acompanamiento` | Devuelto para acompañamiento |

"Guardar revisión" se habilita con profesional, concepto con texto y un resultado. Al guardar, el estado de la remisión pasa a `por_asignar`, `devuelta_para_justificacion` o `devuelta_a_seguimiento`, en ese mismo orden.

---

## Vacío y pie

| Condición | Mensaje |
|---|---|
| Vista agente y el alcance está vacío | No tienes remisiones asignadas. |
| El alcance tiene filas y los filtros no dejan ninguna | No hay remisiones que coincidan con los filtros. |

Pie: `{cantidad visible} de {alcance} remisiones`. `alcance` es el recorte de la vista, no el total filtrado.

---

## Orden

Riesgo inminente primero (`urgency = riesgo_inminente`). Empate: `referred_at` descendente.

---

## Modelo (resumen)

La fila de la tabla es una remisión, no una medida. M-01 crea `salvia.emergency_measure_referral` y agrega `referral_id` a `salvia.emergency_measure`. El `status` de la medida sigue siendo el de hoy (`ACTIVE` al crearla) y no alimenta esta bandeja.

Estados de la remisión, urgencia y el mapa de la revisión están en [M-01](./Flujos/flow-M01-migracion-emergency-measure-referral.md). Los códigos de rol de Agente de Medidas, Equipo Técnico de medidas, Supervisión del Convenio y Cruz Roja no se definen aquí.
