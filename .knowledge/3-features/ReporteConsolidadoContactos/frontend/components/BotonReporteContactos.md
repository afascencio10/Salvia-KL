---
okf_version: "1.0"
type: UI_Component
title: "Componente: Botón + modal de descarga de reportes"
description: "Botón 'Descargar reportes' en la pantalla de reportes del supervisor que abre un modal de rango de fechas y dispara la descarga del Excel consolidado de contactos."
owner: "@frontend-squad"
status: active
tags: [frontend, component, reportes, contactos]
dependencies:
  - Endpoint: 3-features/ReporteConsolidadoContactos/backend/endpoints/reporte-contactos-consolidado-post.md
last_updated: "2026-07-17"
code_refs:
  - src/frontend/html/salvia/victim_case/get_victim_cases_sv.html
  - src/frontend/html/salvia/victim_case/contacts_report_modal.html
  - src/frontend/css/contacts-report-modal.css
  - src/frontend/js/components/contacts-report-modal.js
  - src/salvia/facades/VictimCaseFacade.go
---

# Componente: Botón + modal de descarga de reportes

Se integra en la **pantalla de reportes** del supervisor — "Reportes Salvia" (`get_victim_cases_sv.html`, servida por `VictimCaseFacade.go`) — siguiendo el patrón del componente ya implementado `consolidated-report-modal.js` (Lista de Casos).

> [!NOTE]
> Corrección 2026-07-21: la primera versión de este spec ubicaba el botón en `get_victim_contacts_sv.html`, pero la pantalla que el supervisor realmente ve como "Reportes Salvia" (con el botón "Crear Nuevo caso") es `get_victim_cases_sv.html`. El spec y el código se corrigieron en conjunto.

> [!IMPORTANT]
> **Estilos del modal**: los estilos van en la hoja externa `/static/css/contacts-report-modal.css` (prefijo `crm-`) cargada con `<link>` desde el template del modal — **no** en un bloque `<style>` inline. Vue monta sobre `body#app` y al re-renderizar descarta los `<style>` del template; el patrón de referencia es `reasignar_casos_modal.html` → `reasignar-casos-modal.css`.

## Comportamiento

1. En la pantalla de reportes se muestra el botón **"Descargar reportes"**, ubicado **junto al botón de crear nuevo caso** (decisión de la entrevista 2026-07-17).
2. El botón solo es visible para roles con el permiso `report_contacts_consolidated` (**`ad`**, **`sv`**). La visibilidad en UI es conveniencia; la autorización real la impone el backend.
3. Al pulsarlo se abre un **modal** con dos campos de fecha: **Desde** y **Hasta** (sobre la fecha de **creación** del reporte).
4. Al confirmar, se valida en cliente (ambas fechas, `desde <= hasta`, rango ≤ 366 días) y se hace `POST /api/v1/reportes/contactos-consolidado` con `{ start_date, end_date }`.
5. La respuesta (`.xlsx`, blob) se descarga en el navegador con el nombre `reporte-contactos_<desde>_<hasta>.xlsx`.
6. Si el backend responde `400 VALIDATION_FAILED` porque no hay reportes en el rango, el modal muestra **"No hay reportes en el rango seleccionado"** y no se descarga nada.

## Props (API Pública)

```typescript
interface ReporteContactosModalProps {
  visible: boolean;          // controla apertura del modal
  currentRole: string;       // rol de sesión (para mostrar/ocultar el botón)
}
```

## Estados del Componente

| Estado | Descripción |
| :--- | :--- |
| `idle` | Modal cerrado; botón visible según rol |
| `open` | Modal abierto, esperando selección de rango |
| `validating` | Validación de rango en cliente |
| `loading` | Petición en curso; botón deshabilitado, spinner activo |
| `error` | Error de validación o del backend, visible en el modal |
| `success` | Descarga iniciada; el modal se cierra |

## Validación en cliente (espeja al backend)

- Ambas fechas requeridas.
- `start_date <= end_date`.
- Diferencia ≤ 366 días → si no, mensaje "El rango no puede superar 1 año".

## Manejo de la descarga

- `fetch` con `POST` y cuerpo JSON; leer `response.blob()` en `200`.
- En error (`4xx/5xx`) la respuesta es **JSON**: mostrar el mensaje mapeado del [Catálogo de Errores](/.knowledge/2-data-dictionary/error-catalog.md) dentro del modal (no descargar nada). El caso "cero reportes" llega como `400 VALIDATION_FAILED` y se muestra con el texto específico "No hay reportes en el rango seleccionado".
- Al recibir el blob exitosamente:
  1. Crear un `objectURL` del blob y forzar la descarga en el navegador.
  2. Revocar el URL después de dispararse la descarga.
  3. Cambiar al estado `success` y cerrar el modal.

## Accesibilidad (a11y)

- Inputs de fecha con `aria-label` ("Fecha desde", "Fecha hasta").
- Foco atrapado dentro del modal; cierre con `Esc` y con el botón "Cancelar" (deshabilitado durante `loading`).
- El botón de descarga anuncia estado ocupado (`aria-busy`) durante `loading`.

## Eventos Emitidos

| Evento | Trigger | Payload |
| :--- | :--- | :--- |
| `open` | Click en "Descargar reportes" | — |
| `submit` | Confirmar rango válido | `{ start_date, end_date }` |
| `close` | Cancelar / Esc / éxito | — |
