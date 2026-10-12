## 2026-09-02 — Corrección: alcance de departamento para rol en, no et

El recorte de Lista de casos por departamento asignado aplica al **Enlace territorial** (`en`). El rol entidad (`et`) vuelve a ver el listado completo; su recorte por sede sigue en Casos Entidad.

### Archivos modificados

| Archivo | Cambio |
|---|---|
| `src/salvia/controller/cases_list_controller.go` | Scope de sesión si rol `en` |
| `src/salvia/config/Menu.go` | Permiso `get_list_cases` para `en` |
| `src/frontend/html/salvia/list-cases/list_cases.html` | Subtítulo y alerta para `en` |
| `src/frontend/html/security/general_user/set_general_user.html` | Departamento solo para `en` (como antes) |
| `src/frontend/html/security/general_user/set_general_user_do.html` | Sin selector de departamento (alta et) |
| `src/security/controllers/GeneralUserController.go` | Sin exigir departamento al guardar `et` |

## 2026-09-01 — Lista de casos acotada por departamento (implementación inicial; rol corregido el 2026-09-02)

Implementación del filtro de departamento en `/salvia/lista-casos` vía sesión. El rol destino se corrigió a `en` en la entrada del 2026-09-02.
