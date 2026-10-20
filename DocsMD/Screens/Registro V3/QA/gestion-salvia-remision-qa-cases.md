# Gestión Salvia § B — Remisión a equipos — Casos de prueba

6 equipos (`EQUIPOS_REMISION_V3`): Medidas de Emergencia (`me`, sin criterios), Atención Psicosocial (`ap`, 1 obligatorio + 6 con puntaje, mínimo 3 pts), Estabilización (`es`, 3 criterios, mínimo 1), Atención Hombres (`ah`, 1 obligatorio), Salvia Dignidad (`sd`, sin criterios), Discapacidad (`di`, selección de servicios sin umbral).

## Grupo 1 — Visibilidad condicional (GS_5 → GS_6 → criterios)

| Caso | Descripción | Condición inicial | Resultado esperado | Método de verificación |
|---|---|---|---|---|
| R1 | GS_5 = "No" | — | El bloque de equipos (GS_6) permanece oculto | Validar DOM |
| R2 | GS_5 = "Sí" | GS_6 oculto | Aparece el checklist de 6 equipos | Validar DOM |
| R3 | Marca el equipo "Medidas de Emergencia" | — | Aparece el bloque de medidas de emergencia (GS_7, categoría `victim_case_form2_emergency_measure`, 6 opciones) — es el único equipo con un bloque adicional en vez de criterios de elegibilidad | Validar DOM |
| R4 | Marca "Atención Psicosocial" | — | Aparece el bloque de criterios (GS_9: 1 obligatorio "Mujeres en todas sus diversidades..." + 6 con puntaje) | Validar DOM |
| R5 | Desmarca "Atención Psicosocial" tras R4 | Criterios con selección | El bloque de criterios se oculta | Validar DOM |
| R6 | Marca "Salvia Dignidad" | — | No aparece ningún bloque adicional (no tiene `criterioFieldId`) | Validar DOM |

## Grupo 2 — Cumplimiento de criterios (`cumpleCriterioEquipo`)

### Atención Psicosocial (obligatorio + puntaje mínimo 3)

Puntos por criterio: conducta suicida=3, interseccionalidad=2, no activó ruta=1, condiciones territoriales=1, sin acceso=1, naturalización=1.

| Caso | Descripción | Condición inicial | Resultado esperado | Método de verificación |
|---|---|---|---|---|
| R7 | Marca solo el obligatorio, sin ningún criterio de puntaje | Equipo "Atención Psicosocial" marcado | Badge en rojo "Aún no cumple los criterios mínimos" (`puntos=0 < 3`) | Validar DOM |
| R8 | Marca el obligatorio + "conducta suicida" (3 pts) | — | Badge en verde "Cumple criterios de remisión" (3≥3) | Validar DOM |
| R9 | Marca el obligatorio + "no activó ruta" + "sin acceso" (1+1=2 pts) | — | Badge en rojo (2<3) | Validar DOM |
| R10 | Marca 3 criterios de 1 punto cada uno (no activó ruta + condiciones territoriales + naturalización = 3 pts) **sin marcar el obligatorio** | — | Badge en rojo — el obligatorio cuenta aparte del puntaje (`codigos.indexOf('ob') === -1 → return false` antes de evaluar puntos) | Validar DOM — este caso distingue "cumple puntaje" de "cumple obligatorio" |
| R11 | Marca el obligatorio + todos los 6 de puntaje (3+2+1+1+1+1=9 pts) | — | Badge en verde | Validar DOM |

### Atención Hombres (solo obligatorio, sin puntaje)

| Caso | Descripción | Condición inicial | Resultado esperado | Método de verificación |
|---|---|---|---|---|
| R12 | Marca "Atención Hombres" sin marcar el criterio obligatorio | — | Badge en rojo | Validar DOM |
| R13 | Marca el criterio obligatorio | — | Badge en verde (no exige puntaje adicional, `criterioCategory` de hombres no está en `CRITERIOS_PUNTAJE_V3`, cae en el `return codigos.length > 0`) | Validar DOM |

### Estabilización (3 criterios, mínimo 1, sin obligatorio)

| Caso | Descripción | Condición inicial | Resultado esperado | Método de verificación |
|---|---|---|---|---|
| R14 | Marca "Estabilización" sin marcar ningún criterio | — | Badge en rojo (0 pts < mínimo 1) | Validar DOM |
| R15 | Marca exactamente 1 de los 3 criterios | — | Badge en verde (1≥1, y no hay código "ob" que exigir en esta categoría) | Validar DOM |

### Discapacidad (sin umbral, cualquier selección basta)

| Caso | Descripción | Condición inicial | Resultado esperado | Método de verificación |
|---|---|---|---|---|
| R16 | Marca "Discapacidad" sin marcar ningún servicio | — | Badge en rojo (`codigos.length > 0` es `false`) | Validar DOM |
| R17 | Marca 1 de los 2 servicios | — | Badge en verde | Validar DOM |

## Grupo 3 — Persistencia

| Caso | Descripción | Condición inicial | Resultado esperado | Método de verificación |
|---|---|---|---|---|
| R18 | Guarda la sección GS con 3 equipos marcados (`me`, `ap` con criterios cumplidos, `sd`) | — | `GS_6` en BD como relación multi-valor de categoría `victim_case_form2_referral_team`, 3 filas; `GS_9` (criterios psicosocial) con sus filas propias de categoría `victim_case_form2_psycho_criteria` | Leer BD |
| R19 | Guarda GS_13 (justificación) junto con la remisión | — | `victim_case_form2_referral_justification` con el texto exacto | Leer BD |
