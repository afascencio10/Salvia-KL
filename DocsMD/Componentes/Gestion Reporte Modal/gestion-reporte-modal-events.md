# `gestion-reporte-modal` — Inventario de Eventos

Componente: `GestionReporteModal`
Archivo fuente: `src/frontend/js/components/gestion-reporte-modal.js` (nuevo)
Usado por: `Listado Reportes`, `Mis Reportes`

> **Nota de alcance:** Toda la secuencia está tomada literal de `Mockup_Reportes_Registro.html` (PASO 6a–6e, 7). El mockup simula todo en memoria (`reporteActivo`, `intentos[]`); aquí cada transición que en el mockup es una asignación local se convierte en una llamada real a `GET/POST /api/v1/victim-contacts/...`.

---

## Eventos identificados

### M-01 — Cuando abre el modal

📄 [Ver flujo → flow-M01-cuando-abre-modal.md](./Flujos/flow-M01-cuando-abre-modal.md)

```
Evento:       Cuando abre el modal
Tipo:         Lifecycle
Descripción:  El padre llama open(report) tras el action-clicked 'ver'; el
              header se pinta de inmediato con los datos de la fila.
              Carga el detalle completo vía GET /api/v1/victim-contacts/:icode
              (incluye intentos ya registrados y canManage). Determina el
              step inicial: 'cerrado' si ya tiene close_reason o es un
              reporte legado invalidado (status 'i'); 'postContacto'
              si hubo un intento exitoso sin cerrar; 'esquema3x3' si tiene
              >=1 intento fallido sin exitoso ni cierre; 'validacionPrevia'
              si no tiene ningún intento.
Requerido:    Sí
```

---

### M-02 — Cuando presiona "No información mínima"

📄 [Ver flujo → flow-M02-cuando-cierra-sin-informacion-minima.md](./Flujos/flow-M02-cuando-cierra-sin-informacion-minima.md)

```
Evento:       Cuando presiona "No información mínima"
Tipo:         User Interaction / Backend write
Descripción:  Visible solo en step 'validacionPrevia' y solo si canManage.
              Llama POST /api/v1/victim-contacts/:icode/cerrar con
              close_reason='ni'. Pasa a step 'cerrado'. Si el reporte ya está
              cerrado (o es un legado invalidado) el backend responde 409 y el
              modal no cambia de paso (muestra el error).
Requerido:    Sí
```

---

### M-03 — removido

> Antes "Cuando presiona 'Reporte de prueba'". Instrucción del usuario (2026-10-12): no se construye — los reportes de prueba se generan y descartan en el ambiente de pruebas, no necesitan manejo dentro de producción. El `ButtonGrid` de `validacionPrevia` queda con 2 botones, no 3 (ver `gestion-reporte-modal-interface.md`). Se deja el número M-03 vacío en vez de renumerar M-04 en adelante, para no romper las referencias cruzadas ya escritas en el resto de esta planeación.

---

### M-04 — Cuando presiona "Procede recontacto"

📄 [Ver flujo → flow-M04-cuando-procede-recontacto.md](./Flujos/flow-M04-cuando-procede-recontacto.md)

```
Evento:       Cuando presiona "Procede recontacto"
Tipo:         User Interaction
Descripción:  Cambia el step local de 'validacionPrevia' a 'esquema3x3'.
              No llama al backend — el primer intento real ocurre en M-05/M-06.
Requerido:    Sí
```

---

### M-05 — Cuando registra un intento fallido

📄 [Ver flujo → flow-M05-cuando-registra-intento-fallido.md](./Flujos/flow-M05-cuando-registra-intento-fallido.md)

```
Evento:       Cuando registra un intento fallido
Tipo:         User Interaction / Backend write
Descripción:  Llama POST /api/v1/victim-contacts/:icode/intento con
              result='fallido'. El backend calcula y devuelve los contadores
              actualizados (intentosHoy, diasCompletos) y el nombre del agente
              que lo registró (agentName, para el log). Si diasCompletos
              llega a 3, se habilita (no se dispara sola) la opción de cierre
              M-07 — el esquema 3x3 completo NUNCA cierra el reporte sin que
              el agente confirme.
Requerido:    Sí
```

---

### M-06 — Cuando registra un intento exitoso

📄 [Ver flujo → flow-M06-cuando-registra-intento-exitoso.md](./Flujos/flow-M06-cuando-registra-intento-exitoso.md)

```
Evento:       Cuando registra un intento exitoso
Tipo:         User Interaction / Backend write
Descripción:  Llama el mismo endpoint con result='exitoso'. A diferencia de
              un intento fallido, este SÍ dispara una transición automática
              (con una breve pausa visual, igual al mockup) al step
              'postContacto'.
Requerido:    Sí
```

---

### M-07 — Cuando presiona "Cerrar por no contacto"

📄 [Ver flujo → flow-M07-cuando-cierra-por-no-contacto.md](./Flujos/flow-M07-cuando-cierra-por-no-contacto.md)

```
Evento:       Cuando presiona "Cerrar por no contacto"
Tipo:         User Interaction / Backend write
Descripción:  Visible solo cuando diasCompletos >= 3 (esquema 3x3 completo).
              Llama a cerrar con close_reason='nc'. Pasa a step 'cerrado'.
Requerido:    Sí
```

---

### M-08 — Cuando confirma cierre por no competencia

📄 [Ver flujo → flow-M08-cuando-confirma-cierre-no-competencia.md](./Flujos/flow-M08-cuando-confirma-cierre-no-competencia.md)

```
Evento:       Cuando confirma cierre por no competencia
Tipo:         User Interaction / Backend write
Descripción:  Step 'postContacto', rama 'competencia'. Requiere motivoTexto
              no vacío (botón deshabilitado si no; máx. 500 caracteres: el
              textarea lo limita y el backend responde 400 si se supera). Llama
              a cerrar con close_reason='co' y close_observation=motivoTexto.
Requerido:    Sí
```

---

### M-09 — Cuando confirma cierre por decisión de la persona

📄 [Ver flujo → flow-M09-cuando-confirma-cierre-decision-persona.md](./Flujos/flow-M09-cuando-confirma-cierre-decision-persona.md)

```
Evento:       Cuando confirma cierre por decisión de la persona
Tipo:         User Interaction / Backend write
Descripción:  Idéntico a M-08 con close_reason='dp'.
Requerido:    Sí
```

---

### M-10 — Cuando confirma apertura de caso ⏸️ DIFERIDO — no se implementa en esta fase

> **Desde 2026-10-13** el backend **rechaza `close_reason='ap'` con 400** (`deferredCloseReasonCodes`): tampoco se llega a este cierre por API. Quitar `ap` de ese mapa al construir el evento.
>
> **Instrucción del usuario (2026-10-11):** no construir este evento por ahora — coincide con la nota de FigJam ("aún no vamos a trabajar el registro de un nuevo caso formal a partir de un reporte"). El botón "Apertura de caso" no existe en el `ChoiceGrid` de `postContacto` en esta fase (ver `gestion-reporte-modal-interface.md`). Diseño conservado por si se retoma:

📄 [Ver flujo → flow-M10-cuando-confirma-apertura-caso.md](./Flujos/flow-M10-cuando-confirma-apertura-caso.md)

```
Evento:       Cuando confirma apertura de caso
Tipo:         User Interaction / Backend write
Descripción:  Step 'postContacto', rama 'apertura'. Llama a cerrar con
              close_reason='ap'. Tras éxito, navega a
              /salvia/casos/{icode}/nuevo (ya existente, GAP-7 sobre el
              destino para tipo Familiar/Hombre).
Requerido:    No por ahora — diferido
```

---

### M-11 — Cuando cierra el modal

📄 [Ver flujo → flow-M11-cuando-cierra-modal.md](./Flujos/flow-M11-cuando-cierra-modal.md)

```
Evento:       Cuando cierra el modal
Tipo:         User Interaction
Descripción:  Disparado por la "✕", el backdrop, o el botón "Cerrar" del
              step 'cerrado'. Limpia el estado interno, oculta el modal y
              emite 'updated' siempre (no distingue si hubo cambios).
Requerido:    Sí
```

---

## Checklist de completitud

- [x] ¿La apertura y la detección del step inicial están cubiertas? → M-01
- [x] ¿Las ramas de cierre previo al contacto activas en esta fase? → M-02 (M-03 removido, ver nota)
- [x] ¿El tránsito al esquema 3x3? → M-04
- [x] ¿Ambos resultados de intento? → M-05, M-06
- [x] ¿El cierre tras completar 3x3? → M-07
- [x] ¿Las ramas post-contacto activas en esta fase? → M-08, M-09 (M-10 diferido, ver nota)
- [x] ¿El cierre del modal? → M-11
- [x] ¿Hay lógica de backend desacoplada de la respuesta HTTP? → No
- [x] ¿Hay scheduled tasks, webhooks o sockets? → No

---

## Resumen

| # | Evento | Tipo | Persiste en backend |
|---|---|---|---|
| M-01 | Cuando abre el modal | Lifecycle | No (solo lee) |
| M-02 | Cuando presiona "No información mínima" | User Interaction / Backend write | Sí |
| ~~M-03~~ | ~~Cuando presiona "Reporte de prueba"~~ | — | ❌ Removido (2026-10-12) — no se construye |
| M-04 | Cuando presiona "Procede recontacto" | User Interaction | No (estado local) |
| M-05 | Cuando registra un intento fallido | User Interaction / Backend write | Sí |
| M-06 | Cuando registra un intento exitoso | User Interaction / Backend write | Sí |
| M-07 | Cuando presiona "Cerrar por no contacto" | User Interaction / Backend write | Sí |
| M-08 | Cuando confirma cierre por no competencia | User Interaction / Backend write | Sí |
| M-09 | Cuando confirma cierre por decisión de la persona | User Interaction / Backend write | Sí |
| ~~M-10~~ | ~~Cuando confirma apertura de caso~~ | — | ⏸️ Diferido, no se implementa en esta fase |
| M-11 | Cuando cierra el modal | User Interaction | No (emite hacia el padre) |

**Total a implementar en esta fase: 9 eventos — 1 Lifecycle, 8 User Interaction (6 de ellos con escritura en backend). M-03 removido, M-10 diseñado pero fuera de alcance.**
