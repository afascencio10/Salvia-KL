# `consulta-reporte-publica` — Inventario de Eventos

> ✅ **Implementado y probado — WP-1 T4 (2026-11-02).** Pantalla pública `GET /salvia/public-consulta` (sin sesión) y consulta `POST /api/v1/public/consulta-reporte`. Resumen del cambio en [`changelogNov2026.md`](changelogNov2026.md); casos de prueba en [`QA/`](QA/).

## Archivos de esta pantalla

| Archivo | Descripción |
|---|---|
| [`consulta-reporte-publica-interface.md`](consulta-reporte-publica-interface.md) | Árbol de interfaz, estados visuales y mensajes |
| [`related-tables.md`](related-tables.md) | Tablas que lee la consulta, regla del estado público y protección contra consultas masivas |

### Changelogs

- [`changelogNov2026.md`](changelogNov2026.md)

### Otros MDs

- [`QA/CR-E03-consultar-qa-cases.md`](QA/CR-E03-consultar-qa-cases.md) — API de la consulta (estados, protección, entradas hostiles)
- [`QA/CR-E01-E02-E04-ui-qa-cases.md`](QA/CR-E01-E02-E04-ui-qa-cases.md) — máscara del campo, resultado, error y bloqueo en la interfaz

### Resumen de eventos

| # | Evento | Tipo | Persiste en backend |
|---|---|---|---|
| CR-E01 | [Cuando carga la pantalla](Flujos/flow-E01-cuando-carga-pantalla.md) | Lifecycle | No |
| CR-E02 | [Cuando escribe el radicado](Flujos/flow-E02-cuando-escribe-radicado.md) | User Interaction | No |
| CR-E03 | [Cuando presiona «Consultar»](Flujos/flow-E03-cuando-presiona-consultar.md) | User Interaction | No (solo lee; el limitador es en memoria) |
| CR-E04 | [Cuando presiona «Hacer otra consulta»](Flujos/flow-E04-cuando-presiona-hacer-otra-consulta.md) | User Interaction | No |

---

### CR-E01 — Cuando carga la pantalla

📄 [Ver flujo → flow-E01-cuando-carga-pantalla.md](Flujos/flow-E01-cuando-carga-pantalla.md)

```
Evento:       Cuando carga la pantalla
Tipo:         Lifecycle
Descripción:  GET /salvia/public-consulta renderiza la plantilla consulta_reporte.html
              sin exigir sesión (la URI contiene /public, AuthMiddleware la deja
              pasar). El app Vue appConsulta arranca en la vista 'consulta', con el
              campo vacío y el foco en el radicado.
Requerido:    Sí
```

---

### CR-E02 — Cuando escribe el radicado

📄 [Ver flujo → flow-E02-cuando-escribe-radicado.md](Flujos/flow-E02-cuando-escribe-radicado.md)

```
Evento:       Cuando escribe el radicado
Tipo:         User Interaction
Descripción:  onRadicadoInput() aplica la máscara XXXX-XXXX (mayúsculas, sin signos,
              guion tras el 4.º símbolo, máximo 8), limpia el error anterior y el
              botón «Consultar» se habilita solo si el formato es válido.
Requerido:    Sí
```

---

### CR-E03 — Cuando presiona «Consultar»

📄 [Ver flujo → flow-E03-cuando-presiona-consultar.md](Flujos/flow-E03-cuando-presiona-consultar.md)

```
Evento:       Cuando presiona «Consultar» (o Enter)
Tipo:         User Interaction
Descripción:  POST /api/v1/public/consulta-reporte {radicado}. El backend aplica el
              límite de fallos por IP (5 en 15 min → 429), normaliza el radicado,
              busca con un SELECT parametrizado y responde el estado público
              (recibido, en gestión o cerrado) sin datos personales. Formato inválido
              y radicado inexistente dan el MISMO 404.
Requerido:    Sí
```

---

### CR-E04 — Cuando presiona «Hacer otra consulta»

📄 [Ver flujo → flow-E04-cuando-presiona-hacer-otra-consulta.md](Flujos/flow-E04-cuando-presiona-hacer-otra-consulta.md)

```
Evento:       Cuando presiona «Hacer otra consulta»
Tipo:         User Interaction
Descripción:  nuevaConsulta() limpia el resultado, el campo y el error y vuelve a la
              vista 'consulta'. Si la IP está bloqueada, «Consultar» sigue
              deshabilitado hasta que pase la espera.
Requerido:    Sí
```

---

## Checklist de completitud

- [x] ¿El ciclo de vida inicial está cubierto? → CR-E01
- [x] ¿Toda acción del usuario sobre la UI propia de la pantalla está cubierta? → CR-E02 a CR-E04
- [x] ¿Hay lógica de backend desacoplada de la respuesta HTTP? → No (la consulta es síncrona)
- [x] ¿Hay scheduled tasks o webhooks? → No
- [x] ¿Hay sockets o notificaciones en tiempo real? → No

**Total: 4 eventos — 1 Lifecycle, 3 User Interaction (ninguno escribe en la BD).**
**Todos los eventos tienen flujo documentado (pantalla nueva → modo Metódico).**
