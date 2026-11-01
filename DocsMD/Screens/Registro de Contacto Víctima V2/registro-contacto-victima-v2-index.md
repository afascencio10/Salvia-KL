# `registro-contacto-victima-v2` — Inventario de Eventos

> ✅ **Implementado y probado — 2026-09-27; MD sincronizado con el código el 2026-10-31.** Ver [`changelogSep2026.md`](changelogSep2026.md) para el resumen del cambio y [`QA/E10-cuando-envia-reporte-qa-cases.md`](QA/E10-cuando-envia-reporte-qa-cases.md) para los casos de prueba y su cobertura. Limitaciones conocidas (L1-L4) en [`registro-contacto-victima-v2-interface.md`](registro-contacto-victima-v2-interface.md#limitaciones-conocidas-verificadas-leyendo-el-código-2026-10-31).

## Archivos de esta pantalla

| Archivo | Descripción |
|---|---|
| [`registro-contacto-victima-v2-interface.md`](registro-contacto-victima-v2-interface.md) | Árbol de interfaz — wizard de 4 pasos y condiciones de render |
| [`related-tables.md`](related-tables.md) | Mapeo campo↔BD, columnas/enums nuevos, migración ejecutada y bugs encontrados |

### Changelogs

- [`changelogSep2026.md`](changelogSep2026.md)

### Otros MDs

- [`QA/E10-cuando-envia-reporte-qa-cases.md`](QA/E10-cuando-envia-reporte-qa-cases.md) — casos de prueba y cobertura del único evento que persiste en backend

### Resumen de eventos

| # | Evento | Tipo | Persiste en backend |
|---|---|---|---|
| E-01 | [Cuando carga la pantalla](Flujos/flow-E01-cuando-carga-pantalla.md) | Lifecycle | No |
| E-02 | [Cuando selecciona la autorización de datos personales](Flujos/flow-E02-cuando-selecciona-autorizacion.md) | User Interaction | No |
| E-03 | [Cuando presiona "Continuar" en Autorización](Flujos/flow-E03-cuando-continua-a-perfil.md) | User Interaction | No |
| E-04 | [Cuando selecciona el tipo de reporte](Flujos/flow-E04-cuando-selecciona-tipo-reporte.md) | User Interaction | No |
| E-05 | [Cuando selecciona el ajuste razonable](Flujos/flow-E05-cuando-selecciona-ajuste-razonable.md) | User Interaction | No |
| E-06 | [Cuando presiona "Atrás" en Perfil](Flujos/flow-E06-cuando-presiona-atras-perfil.md) | User Interaction | No |
| E-07 | [Cuando presiona "Continuar" en Perfil](Flujos/flow-E07-cuando-continua-a-formulario.md) | User Interaction | No |
| E-08 | [Cuando selecciona la forma del reporte (anonimato)](Flujos/flow-E08-cuando-selecciona-forma-reporte.md) | User Interaction | No |
| E-09 | [Cuando presiona "Atrás" en el Formulario](Flujos/flow-E09-cuando-presiona-atras-formulario.md) | User Interaction | No |
| E-10 | [Cuando presiona "Enviar Reporte"](Flujos/flow-E10-cuando-envia-reporte.md) | User Interaction | **Sí** (crea `victim_contact` + `victim_contact_form2`) |
| E-11 | [Cuando presiona "Volver al Inicio"](Flujos/flow-E11-cuando-vuelve-al-inicio.md) | User Interaction | No |

> **Nota de alcance:** este inventario cubre únicamente el **Flujo A** del prototipo (`Mockup_Reportes_Registro.html`, pasos `auth → perfil → formulario → confirmacion`). El Flujo B (bandeja de agente, recontacto 3x3, cierre del reporte) ya existe como pantallas aparte: [Listado Reportes](../Listado%20Reportes/listado-reportes-index.md) (asignación por el supervisor), [Mis Reportes](../Mis%20Reportes/mis-reportes-index.md) y el modal [`gestion-reporte-modal`](../../Componentes/Gestion%20Reporte%20Modal/gestion-reporte-modal-index.md). La apertura de caso (M-10) sigue diferida.

> **Nota de completitud:** todos los eventos —incluidos los de navegación trivial (E-03, E-06, E-07, E-09, E-11)— tienen flujo documentado, siguiendo la regla de `dev-docs-workflow.md` para ítems nuevos (categoría A: "sin excepciones, no solo los complejos").

---

## Eventos identificados

---

### E-01 — Cuando carga la pantalla

📄 [Ver flujo → flow-E01-cuando-carga-pantalla.md](Flujos/flow-E01-cuando-carga-pantalla.md)

```
Evento:       Cuando carga la pantalla
Tipo:         Lifecycle
Descripción:  Se ejecuta al abrir GET /salvia/public-v2/nuevo. El backend genera un
              captchaID nuevo y renderiza el template con los catálogos (sí/no,
              tipo de reporte V2, ajustes razonables, anonimato, parentesco);
              si la petición trae `Accept: application/json` devuelve esos
              catálogos y el captchaID en JSON (lo usan la app móvil y los
              helpers de QA). El app Vue arranca en step='auth' con un form en
              blanco (data()) y mounted() fija el título y muestra #app. El
              captcha se muestra recién en el paso 'formulario'.
Requerido:    Sí
```

---

### E-02 — Cuando selecciona la autorización de datos personales

📄 [Ver flujo → flow-E02-cuando-selecciona-autorizacion.md](Flujos/flow-E02-cuando-selecciona-autorizacion.md)

```
Evento:       Cuando selecciona la autorización de datos personales
Tipo:         User Interaction
Descripción:  Fija form.autorizacion (objeto enum, code 'y' = "Sí, autorizo" / 'n' =
              "No autorizo"). Si el code es 'n', muestra de forma reactiva el
              aviso de corte de flujo y deshabilita "Continuar"; si es 'y', lo
              habilita. Sin llamada a backend.
Requerido:    Sí
```

---

### E-03 — Cuando presiona "Continuar" en Autorización

📄 [Ver flujo → flow-E03-cuando-continua-a-perfil.md](Flujos/flow-E03-cuando-continua-a-perfil.md)

```
Evento:       Cuando presiona "Continuar" en el paso de Autorización
Tipo:         User Interaction
Descripción:  Navegación pura: step = 'perfil'. Solo alcanzable si
              form.autorizacion.code === 'y' (botón deshabilitado en caso contrario).
Requerido:    Sí
```

---

### E-04 — Cuando selecciona el tipo de reporte

📄 [Ver flujo → flow-E04-cuando-selecciona-tipo-reporte.md](Flujos/flow-E04-cuando-selecciona-tipo-reporte.md)

```
Evento:       Cuando selecciona el tipo de reporte
Tipo:         User Interaction
Descripción:  Fija form.tipoReporte (code v | r | fm | ho; textos largos en la tabla de
              etiquetas de la interfaz).
              No dispara backend, pero determina qué bloques condicionales del
              paso "formulario" serán visibles más adelante (efecto reactivo
              documentado en ese mismo evento, no como evento aparte).
Requerido:    Sí
```

---

### E-05 — Cuando selecciona el ajuste razonable

📄 [Ver flujo → flow-E05-cuando-selecciona-ajuste-razonable.md](Flujos/flow-E05-cuando-selecciona-ajuste-razonable.md)

```
Evento:       Cuando selecciona el ajuste razonable
Tipo:         User Interaction
Descripción:  Fija form.ajusteRazonable (12 opciones posibles). Junto con
              tipoReporte, habilita el botón "Continuar" del paso Perfil.
Requerido:    Sí
```

---

### E-06 — Cuando presiona "Atrás" en Perfil

📄 [Ver flujo → flow-E06-cuando-presiona-atras-perfil.md](Flujos/flow-E06-cuando-presiona-atras-perfil.md)

```
Evento:       Cuando presiona "Atrás" en el paso Perfil
Tipo:         User Interaction
Descripción:  Navegación pura: step = 'auth'. Los valores ya diligenciados en
              form no se limpian (permite volver a avanzar sin re-digitar).
Requerido:    Sí
```

---

### E-07 — Cuando presiona "Continuar" en Perfil

📄 [Ver flujo → flow-E07-cuando-continua-a-formulario.md](Flujos/flow-E07-cuando-continua-a-formulario.md)

```
Evento:       Cuando presiona "Continuar" en el paso Perfil
Tipo:         User Interaction
Descripción:  Navegación pura: step = 'formulario'. Solo alcanzable si
              tipoReporte y ajusteRazonable están diligenciados.
Requerido:    Sí
```

---

### E-08 — Cuando selecciona la forma del reporte (anonimato)

📄 [Ver flujo → flow-E08-cuando-selecciona-forma-reporte.md](Flujos/flow-E08-cuando-selecciona-forma-reporte.md)

```
Evento:       Cuando selecciona la forma del reporte (Anónimo / A nombre propio)
Tipo:         User Interaction
Descripción:  Solo visible si tipoReporte.code === 'r' (Tercero). Fija form.anonimato y,
              de forma reactiva, muestra u oculta los campos de identidad del
              reportante (nombre, apellido, correo) — mismo criterio que
              identifyTriggers.md sobre reactividad de framework: es el mismo
              evento, no uno nuevo por el re-render.
Requerido:    Sí
```

---

### E-09 — Cuando presiona "Atrás" en el Formulario

📄 [Ver flujo → flow-E09-cuando-presiona-atras-formulario.md](Flujos/flow-E09-cuando-presiona-atras-formulario.md)

```
Evento:       Cuando presiona "Atrás" en el paso Formulario
Tipo:         User Interaction
Descripción:  Navegación pura: step = 'perfil'. No limpia los campos ya
              diligenciados del formulario específico.
Requerido:    Sí
```

---

### E-10 — Cuando presiona "Enviar Reporte"

📄 [Ver flujo → flow-E10-cuando-envia-reporte.md](Flujos/flow-E10-cuando-envia-reporte.md)

```
Evento:       Cuando presiona "Enviar Reporte"
Tipo:         User Interaction
Descripción:  Único evento que persiste en backend. El botón solo se habilita si
              isFormValid (computed, según tipoReporte); al presionarlo arma el
              payload y hace POST a /salvia/public-v2
              (VictimContactPOST_Public_V2 → SetVictimContactV2, controller y
              DAO nuevos e independientes de V1 — decisión 2A). Escribe en
              victim_contact/victim_contact_form2 (columnas existentes, vía
              las funciones DAO de V1 sin modificarlas) y en la tabla de
              extensión nueva victim_contact_form2_v2 (campos exclusivos de
              V2). Si es exitoso, guarda el número de radicado devuelto y
              navega a step='confirmacion'. Es el evento con más ramas
              condicionales de todo el inventario — ver flujo detallado en
              Flujos/flow-E10-cuando-envia-reporte.md.
              Implementado y probado end-to-end (Playwright + curl contra
              backend real, ver QA/) — 7 bugs encontrados y corregidos en esta
              implementación, ver related-tables.md. Limitación vigente (L1):
              Familiar/Hombre con disposición "Sí" y hora digitada reciben 400
              (reporterBestContactTime no se deserializa).
Requerido:    Sí
```

---

### E-11 — Cuando presiona "Volver al Inicio"

📄 [Ver flujo → flow-E11-cuando-vuelve-al-inicio.md](Flujos/flow-E11-cuando-vuelve-al-inicio.md)

```
Evento:       Cuando presiona "Volver al Inicio" en la Confirmación
Tipo:         User Interaction
Descripción:  Resetea form a su estado en blanco (conservando el mismo captchaID) y
              step = 'auth'. Permite iniciar un segundo reporte sin recargar la
              página, pero ese captchaID ya fue consumido por el envío anterior
              (L2): el segundo envío recibe 400 de captcha y un id nuevo.
Requerido:    Sí
```

---

## Checklist de completitud

- [x] ¿El ciclo de vida inicial (carga de datos) está cubierto? → E-01
- [x] ¿Toda acción del usuario sobre la UI propia de la pantalla está cubierta? → E-02 a E-11
- [x] ¿Hay lógica de backend desacoplada de la respuesta HTTP? → No (E-10 es síncrono: valida, escribe, responde)
- [x] ¿Hay scheduled tasks o webhooks? → No
- [x] ¿Hay sockets o notificaciones en tiempo real? → No

**Total: 11 eventos — 1 Lifecycle, 10 User Interaction (de los cuales solo E-10 escribe en el backend).**
**Todos los eventos tienen flujo documentado (módulo nuevo → modo Metódico, sin excepciones).**
