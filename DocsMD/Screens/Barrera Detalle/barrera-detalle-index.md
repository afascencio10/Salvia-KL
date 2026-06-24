# Barrera Detalle — Index

## Navegación

| Documento | Link |
|---|---|
| Interfaz | [barrera-detalle-interface.md](barrera-detalle-interface.md) |

## Resumen de eventos

| # | Evento | Tipo | Flujo |
|---|---|---|---|
| E01 | Cuando carga la pantalla | Lifecycle | [📄 Ver flujo](Flujos/flow-E01-cuando-carga-pantalla.md) |
| E02 | Cuando cambia de tab | User Interaction | — |
| E03 | Cuando activa/desactiva el toggle de estado | User Interaction | — |

---

## Inventario de eventos

📄 [Ver flujo → flow-E01-cuando-carga-pantalla.md](Flujos/flow-E01-cuando-carga-pantalla.md)

**Evento:** Cuando carga la pantalla
**Tipo:** Lifecycle
**Descripción:** `mounted()` muestra el div `#app` (oculto durante la carga del template Go). No realiza llamadas a API — los datos se leen de los valores mock en `data()`. El barrierICode está disponible como variable Vue pero no se usa aún para consultar ningún endpoint.
**Requerido:** Sí

---

**Evento:** Cuando cambia de tab
**Tipo:** User Interaction
**Descripción:** Actualiza `activeTab` ('info' | 'timeline'). Vue muestra u oculta los bloques correspondientes. El tab de Timeline siempre muestra el estado vacío (no hay eventos cargados).
**Requerido:** Sí

---

**Evento:** Cuando activa/desactiva el toggle de estado
**Tipo:** User Interaction
**Descripción:** Toggle local: `barrera.active = !barrera.active`. Cambia el label y la clase visual del switch. No llama ningún endpoint — cambio solo en memoria del cliente.
**Requerido:** No (funcionalidad pendiente de implementar contra API)

---

## Checklist de completitud

- [x] ¿Se cubre la carga inicial de datos?
- [x] ¿Todos los botones de la UI tienen evento?
- [ ] ¿Hay llamadas a API? — No. Pantalla en estado mock, pendiente integración con backend.
