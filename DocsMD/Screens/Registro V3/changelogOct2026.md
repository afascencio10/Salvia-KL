## 2026-10-15 — Implementación y QA de Registro V3

Se construyó `/salvia/registro-v3` como pantalla estática con app Vue a medida, sin `dinamic-form`: Habeas data, tipo de contacto, Cabezote por documento o teléfono, formulario de 7 módulos con guardado por sección, Tamizaje con nivel de riesgo, Gestión Salvia (remisión, entidades, barreras anidadas, alerta inminente), párrafo síntesis y Cierre. La API `/api/v1/victim-cases-v3` escribe en `victim_case`, `victim_case_form2`, `entity_case` y `barrier_v2`. La QA con Playwright del 15/10 corrigió bugs de borrador, tríada, tipo de documento y finalización.

### Archivos modificados

| Archivo | Cambio |
|---|---|
| `src/frontend/html/salvia/registro_v3/registro_v3.html` | Nuevo — template con las 6 etapas: Habeas data, tipo de contacto, Cabezote, Formulario y Cierre |
| `src/frontend/js/components/registro-v3-app.js` | Nuevo — app Vue: estado, validaciones, llamadas a la API, Tamizaje, entidades, barreras y párrafo síntesis |
| `src/frontend/js/components/registro-v3-campos.js` | Nuevo — catálogo de preguntas por módulo, máscaras de documento, textos de enums y criterios de remisión |
| `src/salvia/facades/RegistroV3Facade.go` | Nuevo — `RegistroV3GET`, valida `manage_registro_v3` y renderiza la página |
| `src/salvia/controller/victim_case_v3_controller.go` | Nuevo — API `/api/v1/victim-cases-v3/*` (buscar, draft, sección, tamizaje, finalizar, entidades, barreras) |
| `src/salvia/service/victim_case_v3_service.go` | Nuevo — tríada, Borrador→Iniciado→Terminado, riesgo del Tamizaje, finalización validada y modelo real de entidades/barreras |
| `src/internal/repository/victim_case_v3_repository.go · victim_case_v3_fields.go` | Nuevo — borrador con valores dummy, guardado parcial por módulo, keys respondidas y sincronización del Perfil |
| `src/salvia/facades/MainRouter.go · salvia/config/Menu.go · salvia/config/Enums.go` | Solo adiciones — ruta `/salvia/registro-v3`, permiso `manage_registro_v3` (op, sv), template y estados `in`/`te` |
| `src/main.go` | Solo adiciones — cableado del repositorio, servicio y controller V3 con los repositorios de entidades y barreras |
| `src/salvia/controller/entity_branch_api_controller.go` | Quita el registro duplicado de `/competences` que tumbaba el arranque del servidor |
| `src/migrations/add_registro_v3.sql (partes 1 a 4)` | Nuevo — columnas `tzv3_*` y de Gestión Salvia en `victim_case_form2`, categorías de enum y `answered_keys` (no está en git) |

> En el mismo commit (`2c78f3b`) se modificó `common/utils/CommonSession.go` y `security/controllers/GeneralUserController.go` para cerrar las sesiones activas de un usuario al inhabilitarlo; es un cambio ajeno a esta pantalla.
