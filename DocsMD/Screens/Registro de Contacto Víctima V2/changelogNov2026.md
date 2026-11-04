## 2026-11-02 — WP-1: reportes (radicado, hora, captcha y reportante con sesión)

Cambios en el reporte ciudadano V2 (`/salvia/public-v2/nuevo`). Probados por API, SQL y Playwright contra el sandbox; ver `QA/`.

- **T1 — Hora de recontacto (REP-13\*).** Familiar y Hombre con «Sí» a la llamada ya pueden radicar con la hora digitada: `VictimContactForm2V2DTO.UnmarshalJSON` acepta `HH:MM:SS`, `HH:MM` y RFC 3339; una hora ilegible responde «Hora con formato inválido» en vez de descartarse en silencio. Se cierra la limitación L1.
- **T2 — Captcha nuevo (REP-14\*).** «Volver al Inicio» pide un captchaID nuevo (`refreshCaptcha()`); el segundo reporte seguido ya no falla por un captcha consumido. Se cierra L2.
- **T3 — Radicado (REP-01).** Cada reporte nuevo recibe un radicado `XXXX-XXXX` (columna `victim_contact_radicado`, única) que la confirmación muestra con un botón «Copiar radicado». El `icode` sigue siendo el id interno y viaja también en la respuesta.
- **T8 — Quién reporta (REP-07).** Con la sesión iniciada el reporte queda asociado al usuario y al rol (cualquier rol). Paso inicial «¿Tiene usuario?» (en, op, et) con regreso desde el login (`?next=/salvia/public-v2/nuevo`, lista blanca exacta).
- **T6 — Historial.** El reporte recibido deja el hito «Reporte recibido» en el timeline (después del commit, best-effort).

### Archivos modificados

| Archivo | Cambio |
|---|---|
| `src/frontend/html/salvia/victim_contact/set_victim_contact_v2.html` | `refreshCaptcha()`, radicado y «Copiar radicado», paso `usuario` y `elegirUsuario()` |
| `src/salvia/dao/VictimContactForm2V2DAO.go` | `UnmarshalJSON` del DTO y `ErrReporterBestContactTimeInvalid` |
| `src/salvia/controllers/VictimContactControllerV2.go` | `json.Unmarshal` de form2V2, parámetro `reporter`, radicado en la transacción, respuesta con `radicado` |
| `src/salvia/controllers/radicado.go`, `report_milestones.go` | Nuevos — generador del radicado y hito «Reporte recibido» |
| `src/salvia/dao/VictimContactRadicadoDAO.go` | Nuevo — `AssignVictimContactRadicado` |
| `src/salvia/facades/VictimContactFacadeV2.go` | Sesión opcional de quien reporta (GET y POST) |
| `src/main.go` | +2 columnas en `VictimContactSync` (radicado único y rol del reportante) |
| `src/salvia/config/Locale.go` | Solo adiciones — hora inválida y etiquetas de rol |
| `src/frontend/html/security/login.html` | `?next=` con lista blanca exacta |

Pruebas Go: `go test -vet=off ./salvia/dao -run VictimContactForm2V2DTO` y `./salvia/controllers -run Radicado` (con `-vet=off` porque hay un error de `go vet` previo en `VictimCaseDAO.go:1569`).
