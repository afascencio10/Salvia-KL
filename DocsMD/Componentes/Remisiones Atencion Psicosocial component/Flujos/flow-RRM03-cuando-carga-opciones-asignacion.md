━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando carga opciones del select
   Tipo: Backend read
   ID: RRM-03
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Disparado por: RRM-01 (apertura) o RRM-02 (cambio de switch)

Precondiciones:
  visible === true


PASO 1 — Marcar carga

  loadingOptions = true
  optionsError = null


PASO 2 — Elegir endpoint según modo

  ┌─ SI asignarEnDupla === false ─────────────────────────────┐
  │  GET /api/v1/psychosocial-support/profesionales-reasignacion │
  │  → professionalGroups = response.groups                    │
  └─────────────────────────────────────────────────────────────┘

  ┌─ SI asignarEnDupla === true ───────────────────────────────┐
  │  GET /api/v1/duplas/reasignacion                             │
  │  → duplaOptions = response.duplas                            │
  └─────────────────────────────────────────────────────────────┘


PASO 3 — Procesar respuesta

  SI error HTTP o body.error:
    → optionsError = mensaje
    → loadingOptions = false
    → TERMINAR

  SI ok:
    → Poblar estado correspondiente
    → loadingOptions = false


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  BACKEND — Profesionales (modo individual)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Archivos sugeridos:
  src/internal/repository/psychosocial_professionals_repository.go
  src/salvia/service/psychosocial_reassign_service.go (o dedicado)
  src/salvia/controller/psychosocial_reassign_controller.go

Consulta (referencia):

```sql
SELECT
  gu.general_user_i_code AS icode,
  TRIM(CONCAT(gup.general_user_profile_names, ' ', gup.general_user_profile_last_names)) AS full_name,
  r.role_code AS role
FROM security.general_user gu
JOIN security.general_user_profile gup ON gup.general_user_i_code = gu.general_user_i_code
JOIN security.rel_role_general_user rrgu ON rrgu.general_user_i_code = gu.general_user_i_code
JOIN security.role r ON r.role_i_code = rrgu.role_i_code
WHERE gu.general_user_status = 'e'
  AND r.role_code IN ('ps', 'ts')
ORDER BY r.role_code ASC, full_name ASC
```

Agrupar en Go en dos bloques:
- `ps` → label "Psicólogas"
- `ts` → label "Trabajadoras Sociales"


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  BACKEND — Duplas (modo dupla)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Extender `DuplaRepository` o crear método `ListActiveEnriched`:

```sql
SELECT
  d.id,
  d.name,
  TRIM(CONCAT(ps_gup.general_user_profile_names, ' ', ps_gup.general_user_profile_last_names)) AS psychologist_name,
  TRIM(CONCAT(ts_gup.general_user_profile_names, ' ', ts_gup.general_user_profile_last_names)) AS social_worker_name
FROM salvia.dupla d
JOIN security.general_user ps_gu ON ps_gu.general_user_i_code = d.psychologist_id
JOIN security.general_user_profile ps_gup ON ps_gup.general_user_i_code = ps_gu.general_user_i_code
JOIN security.general_user ts_gu ON ts_gu.general_user_i_code = d.social_worker_id
JOIN security.general_user_profile ts_gup ON ts_gup.general_user_i_code = ts_gu.general_user_i_code
WHERE d.deleted_at IS NULL
  AND ps_gu.general_user_status = 'e'
  AND ts_gu.general_user_status = 'e'
ORDER BY d.name ASC
```

Construir en servicio:
`label = fmt.Sprintf("%s — %s + %s", name, psychologistName, socialWorkerName)`

Ejemplo: `Dupla 1 — Alejandra Mora + Valentina Ospina`
