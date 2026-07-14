---
okf_version: "1.0"
type: UI_Component
title: "Componente: Botón + modal de reporte consolidado"
description: "Botón 'Descargar reporte consolidado' en Lista de Casos que abre un modal de rango de fechas y dispara la descarga del Excel."
owner: "@frontend-squad"
status: active
tags: [frontend, component, reportes]
dependencies:
  - Endpoint: 3-features/ReporteConsolidadoSeguimientos/backend/endpoints/reporte-consolidado-seguimientos-post.md
last_updated: "2026-07-03"
code_refs:
  - src/frontend/html/salvia/list-cases/list_cases.html
  - src/frontend/html/salvia/list-cases/consolidated_report_modal.html
  - src/frontend/js/components/consolidated-report-modal.js
  - src/salvia/facades/VictimCaseFacade.go
---

# Componente: Botón + modal de reporte consolidado

Se integra en la pantalla **Lista de Casos** (`list_cases.html`, controlada por `casos-component.js`), siguiendo el patrón del modal ya existente `reasignar-casos-modal.js`.

## Comportamiento

1. En la barra de acciones de Lista de Casos se muestra el botón **"Descargar reporte consolidado"**.
2. El botón solo es visible para roles con el permiso `report_followups_consolidated` (**`ad`**, **`sv`**). La visibilidad en UI es conveniencia; la autorización real la impone el backend.
3. Al pulsarlo se abre un **modal** con dos campos de fecha: **Desde** y **Hasta**.
4. Al confirmar, se valida en cliente (ambas fechas, `desde <= hasta`, rango ≤ 366 días) y se hace `POST /api/v1/reportes/seguimientos-consolidado` con `{ start_date, end_date }`.
5. La respuesta (`.xlsx`, blob) se descarga en el navegador con el nombre `reporte-casos_<desde>_<hasta>.xlsx`.

## Props (API Pública)

```typescript
interface ReporteConsolidadoModalProps {
  visible: boolean;          // controla apertura del modal
  currentRole: string;       // rol de sesión (para mostrar/ocultar el botón)
}
```

## Estados del Componente

| Estado | Descripción |
| :--- | :--- |
| `idle` | Modal cerrado; botón visible según rol |
| `open` | Modal abierto, esperando selección de rango (muestra texto de advertencia sobre tiempo de generación y límite de filas en rangos amplios) |
| `validating` | Validación de rango en cliente |
| `loading` | Petición en curso; botón deshabilitado, spinner activo y textos informativos alternantes que simulan el progreso |
| `error` | Error de validación o del backend, visible en el modal |
| `success` | Descarga iniciada; el modal se cierra |

## Validación en cliente (espeja al backend)

- Ambas fechas requeridas.
- `start_date <= end_date`.
- Diferencia ≤ 366 días → si no, mensaje "El rango no puede superar 1 año".

## Manejo de la descarga

- Al iniciar la descarga y entrar al estado `loading`, se inicia un intervalo que actualiza un mensaje de texto descriptivo del progreso simulado en el modal:
  * **Segundos 0 a 3**: `"Consultando base de datos..."`
  * **Segundos 4 a 7**: `"Escribiendo registros en las hojas de cálculo..."`
  * **Segundos 8 a 12**: `"Comprimiendo archivo para descarga..."`
  * **Segundos > 12**: `"Iniciando descarga en el navegador..."`
- `fetch` con `POST` y cuerpo JSON; leer `response.blob()` en `200`.
- En error (`4xx/5xx`) la respuesta es **JSON**: detener el intervalo de progreso simulado y mostrar el mensaje mapeado del [Catálogo de Errores](/.knowledge/2-data-dictionary/error-catalog.md) dentro del modal (no descargar nada).
- Al recibir el blob exitosamente:
  1. Detener el intervalo de progreso.
  2. Crear un `objectURL` del blob y forzar la descarga en el navegador.
  3. Revocar el URL después de dispararse la descarga.
  4. Cambiar al estado `success` y cerrar el modal.

## Advertencia Visual en UI (Modal)
- El modal incluirá un aviso claro bajo los campos de fecha:
  > *"Nota: Para rangos de fecha amplios, la generación del reporte consolidado puede tomar hasta 20 segundos debido al volumen de datos. Por favor, no cierres este modal hasta que finalice."*

## Accesibilidad (a11y)

- Inputs de fecha con `aria-label` ("Fecha desde", "Fecha hasta").
- Foco atrapado dentro del modal; cierre con `Esc` y con el botón "Cancelar" (deshabilitado durante `loading`).
- El botón de descarga anuncia estado ocupado (`aria-busy`) durante `loading`.

## Eventos Emitidos

| Evento | Trigger | Payload |
| :--- | :--- | :--- |
| `open` | Click en "Descargar reporte consolidado" | — |
| `submit` | Confirmar rango válido | `{ start_date, end_date }` |
| `close` | Cancelar / Esc / éxito | — |
