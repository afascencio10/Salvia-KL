━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
📋 DECISIÓN DE PLAN: Filtro por profesional incluye remisiones vía dupla
   Código: DEC-E08-01
   Afecta: E-01, E-05, E-08 (+ stats con mismos filtros)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

## Problema detectado

Un mismo profesional del equipo psicosocial puede aparecer en remisiones de **dos formas distintas**:

| Vía | Campo en BD | Ejemplo |
|---|---|---|
| Asignación directa | `psychosocial_support.professional_id` | Remisión asignada solo a la psicóloga P1 |
| Asignación por dupla | `psychosocial_support.dupla_id` → `dupla.psychologist_id` o `dupla.social_worker_id` | Remisión asignada a la dupla donde P1 es psicóloga |

La implementación inicial de **E-08** filtraba solo:

```sql
AND ps.professional_id = {filter_professional_id}
```

**Consecuencia:** al seleccionar un profesional en el autocomplete, **no aparecían** las remisiones donde ese usuario participa únicamente como miembro de una dupla (sin `professional_id` en la fila).

---

## Regla de negocio acordada

`filter_professional_id` representa al **profesional seleccionado** (su `general_user_i_code`), no solo la asignación directa.

Debe devolver todas las remisiones donde el profesional está vinculado por **cualquiera** de estas vías:

1. **Directa:** `ps.professional_id = icode`
2. **Indirecta (dupla):** existe una dupla activa donde `psychologist_id = icode` **o** `social_worker_id = icode`, y `ps.dupla_id` apunta a esa dupla

La unión es **OR** (conjunto ampliado). Si una remisión cumple ambas condiciones, aparece una sola vez.

---

## Impacto en frontend (E-07 / E-08)

**Sin cambio de contrato HTTP.**

- E-07 sigue buscando profesionales en `GET /api/v1/agents/search-psicosocial`.
- E-08 sigue enviando un único param: `filter_professional_id={icode}`.
- No se envían ids de dupla desde el autocomplete; el backend resuelve las duplas del profesional.

El usuario no necesita saber si la remisión está asignada por dupla o por profesional directo: el filtro responde a **“todas las remisiones de este profesional”**.

---

## Impacto en backend

### Cláusula WHERE — `filter_professional_id`

Reemplazar la condición simple por:

```sql
AND (
  BTRIM(ps.professional_id::text) = BTRIM({filter_professional_id})
  OR ps.dupla_id IN (
    SELECT d.id
    FROM salvia.dupla d
    WHERE d.deleted_at IS NULL
      AND (
        BTRIM(d.psychologist_id::text) = BTRIM({filter_professional_id})
        OR BTRIM(d.social_worker_id::text) = BTRIM({filter_professional_id})
      )
  )
)
```

Aplicar la **misma condición** en:

- `GET /api/v1/psychosocial-support/list` (E-01, E-08, E-06, …)
- `GET /api/v1/psychosocial-support/stats` (cuando `filter_professional_id` está activo)

### Relación con `filter_dupla_id` (E-11)

Los filtros siguen siendo **combinativos (AND)** entre sí.

| Filtros activos | Comportamiento |
|---|---|
| Solo `filter_professional_id` | Remisiones directas **OR** vía duplas del profesional |
| Solo `filter_dupla_id` | Solo remisiones de esa dupla (E-11, sin cambio) |
| Ambos a la vez | Intersección: remisión debe cumplir dupla concreta **y** pertenecer al profesional (caso raro en UI) |

---

## Impacto en `defaultFilter` (E-01 / E-05)

Pantalla **“Mis remisiones (profesional)”** con:

```js
defaultFilter: { professional_id: sessionAgentIcode }
```

Debe usar la **misma regla ampliada** en backend, de modo que el usuario vea:

- Remisiones con `professional_id = su icode`
- Remisiones de duplas donde es `psychologist_id` o `social_worker_id`

Ya **no** es necesaria una pantalla separada “Mis remisiones por dupla” solo para cubrir este hueco, aunque `defaultFilter.dupla_id` sigue siendo válido si el producto quiere acotar a **una dupla concreta**.

---

## Modelo de datos — sin cambio de schema

No se agregan columnas ni tablas. Es un ajuste de **criterio de consulta** sobre:

- `psychosocial_support.professional_id`
- `psychosocial_support.dupla_id`
- `dupla.psychologist_id` / `dupla.social_worker_id`

---

## Checklist de implementación

- [ ] Actualizar `buildPsychosocialListWhere` — condición OR para `filter_professional_id`
- [ ] Verificar stats con el mismo filtro
- [ ] Probar profesional con solo asignaciones directas
- [ ] Probar profesional que solo aparece en duplas (sin `professional_id` en filas)
- [ ] Probar profesional con ambos tipos de asignación
- [ ] Probar combinación con otros filtros (estado, identidad, etc.)

---

## Documentos actualizados

| Archivo | Cambio |
|---|---|
| [flow-E08-cuando-selecciona-autocomplete.md](./flow-E08-cuando-selecciona-autocomplete.md) | WHERE ampliado |
| [flow-E01-cuando-carga-componente.md](./flow-E01-cuando-carga-componente.md) | Tabla `filter_professional_id` |
| [flow-E05-cuando-limpia-filtros.md](./flow-E05-cuando-limpia-filtros.md) | Alcance “Mis remisiones (profesional)” |
| [remisiones-psicosocial-component-interface.md](../remisiones-psicosocial-component-interface.md) | Criterio filtro + defaultFilter |
| [remisiones-psicosocial-component-events.md](../remisiones-psicosocial-component-events.md) | E-08, decisiones, tabla impacto |
