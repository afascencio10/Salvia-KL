## 2026-08-04 — Nueva tab "Oficios" (reutiliza case-oficios filtrado por barrera)

Se agregó una cuarta tab "📄 Oficios" a la pantalla Detalle de Barrera, entre "Tareas" y "Timeline", que muestra únicamente los oficios (`entity_letter`) de esa barrera puntual. Reutiliza el componente `case-oficios.js` — ya implementado y en uso en Detalle del Caso, pero hasta ahora solo sabía filtrar por caso completo. Se le agregó un prop opcional `barrierId` que, cuando está presente, tiene prioridad sobre `caseId`, usa `GET /api/v1/entity-letters?barrierId=` (endpoint que ya existía) y oculta los chips de tema (irrelevantes cuando ya está todo filtrado a una sola barrera). Sin cambios de backend ni de base de datos. Ver plan completo en `DocsMD/Otros/temp/req-oficios-en-detalle-barrera.md`.

### Archivos modificados

| Archivo | Cambio |
|---|---|
| `src/frontend/js/components/case-oficios.js` | Nuevo prop `barrierId` (opcional); `cargar()` usa `?barrierId=` si está presente; oculta `.co-chips-tema` en modo barrera |
| `src/frontend/html/salvia/barriers/barrera_detalle.html` | Nueva tab "📄 Oficios" (botón + contenido), monta `<case-oficios>`; carga el script del componente |
