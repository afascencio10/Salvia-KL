# Repeater de Seguimiento a Barreras — Estructura de Preguntas

**Sección:** Seguimiento a Barreras (Sección 3 del Formulario de Seguimiento)
**Section ID:** `81e75efd-a885-4fd7-af84-48c03238da67`

## Repeater Group

| Campo | Valor |
|---|---|
| Nombre | Seguimiento a Barreras Activas |
| item_name | Barrera |
| add_button_label | Siguiente Barrera |

> Cada entrada del repeater representa el seguimiento a una barrera registrada.

> ⚡ **Visibilidad de la sección:** Solo se muestra si hay una barrera activa registrada en el caso — condición basada en `formState` externo (`trigger_state_path`, path a definir).

---

## Estructura de preguntas por entrada

### 1. Seguimiento a Barrera
| Campo | Valor |
|---|---|
| Tipo | `info` |
| Requerido | ❌ |
| Condición | — siempre visible |

---

### 2. ¿Persiste la barrera?
| Campo | Valor |
|---|---|
| Tipo | `boolean` |
| Requerido | ✅ |
| Condición | — |

---

### 3. ¿Hubo respuesta institucional?
| Campo | Valor |
|---|---|
| Tipo | `single` |
| Requerido | ✅ |
| Condición | — |

| # | Label | Value |
|---|---|---|
| 1 | Sí. La entidad respondió de manera oficial a Salvia | `respuesta_oficial` |
| 2 | Sí. La entidad se contactó con la víctima | `contacto_victima` |
| 3 | No se obtuvo respuesta | `sin_respuesta` |

---

### 4. Gestión de la barrera
| Campo | Valor |
|---|---|
| Tipo | `multiple` |
| Requerido | ✅ |
| Condición | — |

| # | Label | Value |
|---|---|---|
| 1 | Orientación y enrutamiento - Llamada | `orientacion_llamada` |
| 2 | Gestión administrativa - Llamada | `gestion_llamada` |
| 3 | Activación de ruta interinstitucional | `activacion_ruta_interinstitucional` |
| 4 | Articulación institucional | `articulacion_institucional` |
| 5 | Escalamiento a organismo de control | `escalamiento_organismo_control` |
| 6 | Alerta por barreras | `alerta_barreras` |

---

### 5. Actuaciones realizadas y descripción de la gestión realizada con relación a las barreras
| Campo | Valor |
|---|---|
| Tipo | `text` |
| Requerido | ✅ |
| Condición | — |

---

### 6. ¿Se realiza cierre de la barrera?
| Campo | Valor |
|---|---|
| Tipo | `boolean` |
| Requerido | ✅ |
| Condición | — |

---

### 7. Motivo del cierre
| Campo | Valor |
|---|---|
| Tipo | `single` |
| Requerido | ✅ |
| Condición | Q6 = `true` \| EQUALS |

| # | Label | Value |
|---|---|---|
| 1 | Expresa no voluntad de accionar institucional | `no_voluntad` |
| 2 | Barrera no gestionable desde la competencia institucional | `no_gestionable` |
| 3 | Se clasificó incorrectamente la barrera | `clasificacion_incorrecta` |
| 4 | La barrera no fue resuelta directamente, pero quedó formalmente instalada en la entidad competente | `instalada_entidad` |
| 5 | La barrera fue resuelta de manera efectiva | `resuelta` |
| 6 | La barrera fue superada parcialmente y no requiere más gestión inmediata | `superada_parcialmente` |

---

## Resumen de VCs del repeater

| Target (Q#) | Trigger (Q#) | Trigger value | Operator |
|---|---|---|---|
| Q7 Motivo del cierre | Q6 ¿Se realiza cierre? | `true` | EQUALS |

---

## Pendientes / Decisiones abiertas

- [ ] Definir `trigger_state_path` para la visibilidad de la sección completa (ej. `currentCase.hasActiveBarrier`)
