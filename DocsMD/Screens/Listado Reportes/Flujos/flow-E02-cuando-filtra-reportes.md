━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando filtra reportes
   Tipo: User Interaction
   Código: E-02
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Disparado por: clic en "Buscar" o "Limpiar" en FilterBarReportes

INPUT: {
  filterNames:      v-model input nombres
  filterLastNames:  v-model input apellidos
  filterPhone:      v-model input teléfono
  currentFilter:    fcv | fci  (tab activo)
}


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

PASO 1 — applyReportFilters()  [Botón Buscar]

  currentPage = 0
  → goToPage(0) usando buildContactListUrl() con query params

PASO 2 — clearReportFilters()  [Botón Limpiar]

  filterNames     = ""
  filterLastNames = ""
  filterPhone     = ""
  currentPage     = 0
  → goToPage(0) sin query params

PASO 3 — Llamada API

  GET /salvia/primer_contacto/f/{currentFilter}/p/{currentPage}
    ?names={filterNames}           // opcional
    &lastNames={filterLastNames}   // opcional
    &phone={filterPhone}           // opcional

  SI respuesta ok:
    victimContacts = response.data
    numPages       = response.numPages

  → FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  BACKEND — VictimContactFacade.VictimContactGET
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

PASO 4 — Leer query params

  names     = c.Query("names")
  lastNames = c.Query("lastNames")
  phone     = c.Query("phone")

PASO 5 — Llamar DAO extendido

  GetVictimContactsWithoutVictimCase(page, status, filters, ...)
  // status = "v" si filter fcv, "i" si filter fci


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  BACKEND — VictimContactDAO (cambio requerido)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

PASO 6 — Query base (sin cambio)

  LEFT JOIN victim_case ... WHERE victim_case IS NULL
  AND victim_contact_status = $status

PASO 7 — Filtros aditivos (AND)

  SI names != "":
    AND victim_contact_names ILIKE '%' || names || '%'

  SI lastNames != "":
    AND victim_contact_last_names ILIKE '%' || lastNames || '%'

  SI phone != "":
    AND EXISTS (
      SELECT 1 FROM salvia.victim_contact_form2 vcf2
      WHERE vcf2.victim_contact_form2_victim_contact = vc.victim_contact_id
      AND vcf2.victim_contact_form2_victim_col_phone::text ILIKE '%' || phone || '%'
    )

  // Requiere JOIN o subquery a victim_contact_form2 para teléfono

PASO 8 — Paginación y COUNT igual que hoy


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FIRMA PROPUESTA DAO
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

```go
type VictimContactListFilters struct {
    Names     string
    LastNames string
    Phone     string
}

func GetVictimContactsWithoutVictimCase(
    page int,
    status string,
    filters VictimContactListFilters,
    connData *db.ConnData,
    ...
) ([]VictimContactDTO, int, error)
```

Actualizar también la carga inicial en `VictimCaseFacade` (E-01) para usar la misma firma con filtros vacíos.


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  EJEMPLOS
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

| names | lastNames | phone | Resultado |
|---|---|---|---|
| `"María"` | *(vacío)* | *(vacío)* | Nombre contiene "María" |
| *(vacío)* | `"García"` | *(vacío)* | Apellido contiene "García" |
| *(vacío)* | *(vacío)* | `"300"` | Teléfono form2 contiene "300" |
| `"Ana"` | `"López"` | `"310"` | Cumple las tres condiciones (AND) |
