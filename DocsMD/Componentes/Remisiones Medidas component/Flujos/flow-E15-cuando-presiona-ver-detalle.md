━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando el usuario abre el detalle de la remisión
   Tipo: User Interaction
   Código: E-15
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Disparado por:
  - el link "Ver detalle →" de la columna Remisión o Medidas
  - un botón de fila cuyo tipo es detalle

"Asignar agente" es E-08. "Revisar (Equipo Técnico)" es E-11.
Este evento no abre esos modales.


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

PASO 1 — Qué botón es detalle

  El link "Ver detalle →" siempre emite, en cualquier vista y estado.

  El botón de la fila emite solo en estas combinaciones. La etiqueta es la
  del botón; el payload es el mismo en todas.

  | Vista                  | Estado                                      | Etiqueta                      |
  |------------------------|---------------------------------------------|-------------------------------|
  | supervisor             | riesgo_inminente_por_decidir_via            | Decidir vía de atención       |
  | supervision-convenio   | riesgo_inminente_por_escalar_cruz_roja      | Escalar a Cruz Roja           |
  | cruz-roja              | riesgo_inminente_por_entregar_medida        | Entregar medida inmediata     |
  | supervisor             | caracterizado                               | Revisar caracterización       |
  | supervision-convenio   | por_revisar_trazabilidad                    | Revisar trazabilidad          |
  | supervision-convenio   | por_aprobar_cotizaciones                    | Aprobar cotizaciones          |
  | cruz-roja              | por_cotizar                                 | Cargar cotizaciones           |
  | agente                 | riesgo_inminente_por_confirmar_seguridad    | Confirmar seguridad           |

  Cualquier otra combinación no pinta botón de acción. En las vistas de
  gestión, la celda muestra el agente asignado o "Sin agente asignado",
  y el concepto técnico si ya hay revisión.

PASO 2 — Emitir

  $emit('ver-detalle', {
    remisionId: item.id,
    caseICode:  item.caseICode
  })

  El payload no distingue el botón. La pantalla de detalle lee el estado
  de la remisión y decide qué trabajo mostrar.

PASO 3 — El padre navega

  El componente no cambia de ruta. El padre escucha @ver-detalle.

  Remisiones Medidas y Mis Remisiones Medidas hacen:

    location = /salvia/remision-medida/{remisionId}

  Esa pantalla lee la remisión. Hoy es de lectura: cabecera, Info Medida,
  tareas y timeline del caso. Decidir la vía, escalar, entregar, revisar
  la caracterización, la trazabilidad, las cotizaciones y confirmar la
  seguridad siguen pendientes en esa pantalla, no en este listado.

  onVerCaso sigue vacío. E-14 no navega.

  → FIN EJECUCIÓN ✓
