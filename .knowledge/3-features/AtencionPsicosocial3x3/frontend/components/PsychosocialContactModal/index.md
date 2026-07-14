---
okf_version: "1.0"
type: UI_Component
title: "Componente: psychosocial-contact-modal"
description: "Modal Vue que orquesta el flujo 3x3 de contacto psicosocial: decisión ¿contestó?, registro de intento con fecha/hora editable, Consentimiento Informado, agendamiento de sesión (inmediata/flexible) y modal de acciones al 3er intento fallido del día (próximo intento, añadir intento, cierre por imposibilidad)."
owner: "@frontend-squad"
status: active
tags: [frontend, component, vue, psicosocial, 3x3]
dependencies:
  - Endpoint: 3-features/AtencionPsicosocial3x3/backend/endpoints/register-contact-attempt.md
  - Endpoint: 3-features/AtencionPsicosocial3x3/backend/endpoints/set-consent.md
  - Endpoint: 3-features/AtencionPsicosocial3x3/backend/endpoints/schedule-session.md
  - Endpoint: 3-features/AtencionPsicosocial3x3/backend/endpoints/set-next-attempt.md
  - Endpoint: 3-features/AtencionPsicosocial3x3/backend/endpoints/init-closure-form.md
  - Endpoint: 3-features/AtencionPsicosocial3x3/backend/endpoints/complete-closure.md
  - Model: 2-data-dictionary/closure-form-model.md
code_refs:
  - "src/frontend/js/components/psychosocial-contact-modal.js"
last_updated: "2026-07-09"
---

# Componente: `<psychosocial-contact-modal />`

Componente Vue (registrado con `app.component`, delimitadores `['${', '}']`), paralelo a `follow-up-contact-modal.js` pero para el módulo psicosocial. Opera sobre un proceso `psychosocial_support` (no sobre `follow_up_v2`).

Es el **flujo de modales** que hospeda el card [`psychosocial-contact-card`](../PsychosocialContactCard/index.md): el card lo compone e invoca por `ref` con `start(process)` desde su botón "Registrar intento", y escucha su evento `completed` para refrescar el historial.

## Props (API Pública)

```typescript
interface PsychosocialContactModalProps {
  currentUser: string;    // nombre de la profesional
  currentUserId: string;  // id de la profesional (ps/ts)
  userTeam: string;       // equipo
}
```

Se invoca por `ref` desde el padre: `start(process)` donde `process` incluye `{ id (psicosocialId), caseId, caseName, dailyFailedCount, totalCount, distinctDaysCount }`.

## Estados del Componente

| Estado (data flag) | Descripción |
| :--- | :--- |
| `modalContesto` | Diálogo "¿Contestó la llamada?" con banner de intentos previos, selector de fecha/hora del intento y **campo de nota** libre (aplica a ambos caminos). Es el **único** paso de captura del intento (no hay segundo modal). |
| `modalConsentimiento` | Consentimiento Informado (Acepta / No acepta). |
| `modalSesion` | Elegir sesión inmediata o reagendar (calendario flexible). |
| `modalAcciones` | Acciones al alcanzar 3 intentos fallidos del día (3a/3b/3c). |
| `modalNextAttempt` | Calendario/inputs para fijar `next_contact_attempt_at` (3a). |
| `modalCierre` | Formulario de cierre: hospeda `<dinamic-form>` con el `form_id` del [Cierre de proceso psicosocial](/.knowledge/2-data-dictionary/closure-form-model.md), Motivo preseleccionado (editable). Se **abre de inmediato** al disparar el cierre y muestra un **loader interno** (`closureLoading`) mientras `init-closure-form` responde; luego renderiza el `dinamic-form` (que a su vez tiene su propio loader al cargar la estructura). |
| `loading` | Operación en progreso. |
| `showSuccessOverlay` | Overlay de éxito temporal. |

## Lógica de apertura (`start`)

```text
initAttemptDateTime()           // default = ahora (datetime-local, editable)
if process.dailyFailedCount >= 3:
    modalAcciones = process      // salta la alerta, va directo a acciones
else:
    modalContesto = process
```

## Árbol de Renderizado

```text
<psychosocial-contact-modal />
│
├── [modalContesto] "¿Contestó la llamada?"  (layout = mockup aprobado)
│     ├── icono teléfono (círculo violeta)
│     ├── "Estamos contactando a" + ${caseName}
│     ├── banner rojo "${attempts} intento(s) sin respuesta"  (solo si attempts > 0)
│     ├── input datetime-local "Fecha y hora del intento" (default ahora, editable)
│     ├── textarea "Nota" (opcional) → v-model note  [aplica a ambos botones]
│     ├── botón "No contestó"  → handleNoAnswer() → POST attempt(was_answered=false, note) →
│     │        if dailyThresholdReached: modalAcciones  else: cerrar con éxito  (sin segundo modal)
│     ├── botón "Sí contestó ✓" → handleYesAnswer() → POST attempt(was_answered=true, note) → modalConsentimiento
│     └── enlace "CANCELAR" → closeModals()
│
├── [modalConsentimiento] "Consentimiento Informado"
│     ├── botón "No acepta" → PATCH consent(false) → POST init-closure-form(reason=no_consentimiento) → modalCierre
│     └── botón "Sí acepta" → PATCH consent(true)  → modalSesion
│
├── [modalSesion] "¿Cuándo se realiza la sesión?"
│     ├── botón "De inmediato" → POST sessions(immediate=true)  → redirect a formulario de atención (🔴 futuro)
│     └── botón "Agendar"      → inputs fecha/hora flexible → POST sessions(immediate=false) → éxito
│
├── [modalAcciones] "Límite de intentos del día"
│     ├── (3a) "Fijar próximo intento" → modalNextAttempt → PUT next-attempt
│     ├── (3b) "Añadir intento"        → habilitado si totalCount < 50 → vuelve a modalContesto
│     └── (3c) "Cerrar por imposibilidad" → visible si totalCount>=9 && distinctDays>=3 → POST init-closure-form(reason=imposibilidad_contacto_3x3) → modalCierre
│
└── [modalCierre] "Formulario de Cierre de proceso psicosocial"
      ├── <dinamic-form :form-id="fcc7dc8d…" :submission-id :form-state> (Motivo preseleccionado según reason)
      └── @form-completed → POST /psychosocial/{id}/close {submission_id} → status cerrado/en_devolucion → cerrar + éxito
```

## Diferencias clave vs `follow-up-contact-modal.js`

| Aspecto | Seguimientos | Psicosocial 3x3 |
| :--- | :--- | :--- |
| Tabla de intentos | `follow_up_attempts` | `contact_attempts` |
| Umbral que dispara acciones | Total `attempts >= 3` | **Fallidos del día `>= 3`** |
| Posponer | Automático +1/+2/+3 días | **Manual: fecha/hora exacta** (`next_contact_attempt_at`) |
| Tras "Sí contestó" | Redirect a `hacer-seguimiento` | Abre **Consentimiento Informado** |
| Tope de intentos | (sin tope explícito) | **50** |
| Cierre | Formulario de cierre de caso | Formulario de cierre **psicosocial** (`dinamic-form`, motivo preseleccionado (editable)) — cierra solo el proceso |

## Accesibilidad (a11y)

- `datetime-local` y textarea con `aria-label` descriptivo.
- Botones de acciones navegables por teclado; foco inicial en la acción primaria.
- Estados `disabled` con `aria-disabled` durante `loading`.

## Eventos Emitidos

| Evento | Trigger | Payload |
| :--- | :--- | :--- |
| `completed` | Intento fallido guardado | `{ type: 'attempt', psicosocialId, counters }` |
| `completed` | Contacto exitoso + en_gestion | `{ type: 'contacted', psicosocialId }` |
| `completed` | Sesión agendada | `{ type: 'scheduled', psicosocialId, sessionId }` |
| `completed` | Próximo intento fijado | `{ type: 'next-attempt', psicosocialId, nextContactAttemptAt }` |
| `completed` | Cierre disparado | `{ type: 'closure-triggered', psicosocialId, reason }` |
