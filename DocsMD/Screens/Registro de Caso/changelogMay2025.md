# Changelog — Mayo 2025

## 2026-05-22 — Evento de hechos del caso en el timeline

Al registrar un caso se inserta automáticamente un evento de tipo **"Hechos del caso"** en el timeline, con los subtipos de violencia y el relato de los hechos.

El evento usa la fecha de los hechos (`factsDate`) como fecha del evento, no la fecha de creación.

Se ejecuta en una goroutine post-commit — si falla, no afecta la creación del caso.

---

### Archivos modificados

| Archivo | Cambio |
|---|---|
| `src/internal/models/case_timeline_event.go` | Nuevas constantes: `TimelineTypeHechosCaso`, `TimelineIconHechosCaso`, `TimelineColorLightRed` |
| `src/salvia/controllers/VictimCaseController.go` | Variable `CaseTimelineRepo` inyectable + goroutine post-commit que inserta el evento |
| `src/main.go` | Inyección de `CaseTimelineRepo` en el controller legacy |
| `DocsMD/Screens/Registro de Caso/registro-caso-flujo-guardar.md` | Flujo actualizado con el nuevo sub-flujo y nota en PASO 18 |
