# `Ruta de atención integral` — Interfaz del módulo

> ✅ Módulo `RA` de `/salvia/registro-v3` (REG-17, WP-4). Pestaña entre «Tamizaje de riesgo» y «Plan de acción Salvia»; `moduloBloqueado('RA')` es siempre `false`. Un solo bloque `v-if="moduloActivo==='RA'"` en `registro_v3.html`; estado y métodos en `registro-v3-app.js`; catálogo de preguntas en `registro-v3-campos.js` (`CAMPOS_RA_V3`, `ENTIDAD_V4_V3`).

## Archivos relevantes

| Archivo | Rol |
|---|---|
| `src/frontend/html/salvia/registro_v3/registro_v3.html` | Bloque del módulo y fragmentos compartidos de la plantilla: `rv3_aviso_borrador`, `rv3_entidad_buscador`, `rv3_tarjeta_entidad`, `rv3_form_barrera` |
| `src/frontend/js/components/registro-v3-app.js` | `modulos[]` (entrada `RA`), `entidadesAcudidas`, `agregarEntidad()`, `entidadFaltantes()`, `activacionAplica()`, `guardarRutaAtencion()`, solo lectura (`entitiesReadOnly`) |
| `src/frontend/js/components/registro-v3-campos.js` | `CAMPOS_RA_V3` (`GS_14`, `GS_17`), `ENTIDAD_V4_V3` y `CANAL_ACTIVACION_V3` (generados del seed de Seguimiento v4), `BARRERA_V4_V3` |
| `scripts/gen_registro_v3_barrera_catalog.py` | Genera esos catálogos desde `src/internal/seed/seguimiento_v4_form.sql` |

## Árbol de interfaz

```
[moduloActivo === 'RA']
RutaAtencionIntegral
├── Texto guía  "Registra a qué rutas o entidades acudió la persona y si encontró barreras en cada una. Es informativo: las activaciones de ruta se registran en el Plan de acción Salvia."
├── Aviso de borrador (.alert-warning, [v-if !entitiesReadOnly])  |  Aviso de solo lectura (.alert-info, [v-else])      ← fragmento rv3_aviso_borrador
├── Radio  GS_14  "¿Ha puesto en conocimiento de alguna entidad o institución los hechos de violencia reportados? Sí | No"   (obligatoria, 1er contacto)
├── [v-if GS_14 === 'true' || (GS_17 === 'true' && !entidadesAcudidas.length)]
│   ├── TarjetaEntidad × N  [v-for entidadesAcudidas]  (.card)                                      ← fragmento rv3_tarjeta_entidad
│   │   ├── Badge sector · nombre de la sede · ubicación (departamento · ciudad · municipio)
│   │   ├── [v-if canales]  "Canales de activación: …"  ·  [v-if acciones]  "Acciones institucionales:" chips "✓ Completada: …" / "Pendiente: …"
│   │   ├── [v-if !entitiesReadOnly]  BtnEliminarEntidad  "Eliminar entidad"  → eliminarEntidad(relId)
│   │   ├── [v-if barreras]  Tabla  Descripción · Gestión (etiqueta) · "Eliminar" [v-if !entitiesReadOnly]
│   │   ├── [v-if GS_17 === 'true' && !entitiesReadOnly]  BtnAgregarBarrera  "+ Agregar barrera a esta entidad"
│   │   └── [v-if barreraParaEntidad === e.entityBranchId && !entitiesReadOnly]  FormularioBarrera (rv3_form_barrera: formulario de barrera de Seguimiento v4)
│   └── [v-if !entitiesReadOnly]  FormularioEntidadAcudida  (.card)  "Entidad a la que ya acudió"
│       ├── Cascada y búsqueda  (fragmento rv3_entidad_buscador)
│       │   ├── <select> Departamento*  →  Ciudad*  →  Municipio*  →  Sector*   (cada uno deshabilitado hasta elegir el anterior)
│       │   ├── "Buscando sedes disponibles…"  |  "No hay sedes registradas en este municipio y sector."
│       │   ├── <input> "Buscar entidad"*  + lista de sugerencias  |  "Ninguna entidad coincide con lo escrito."      ← lista cargada una vez; filtro local con debounce de 250 ms
│       │   └── [v-if EntityBranchID]  "Acciones institucionales"  Completadas (chips) · Pendientes (chips)  ← GET /entity-branches/:id/obligations (sin «otra»; una acción no puede estar en ambas)
│       ├── [v-if EntityBranchID]  <input> "Número de radicado (NUNC, medida de protección). Escriba tipo y número"*  ·  <input type=date> "Fecha en la que reportó a la entidad"*  ·  <textarea> "Otra información asociada a la Ruta"*
│       ├── [v-if PendingActions.length]  Radio RA_acuerdo  "¿La persona está de acuerdo con que SALVIA active la Ruta de Atención y las gestiones institucionales correspondientes?"*  Sí | No
│       │   ├── [v-if Agreement === true]  Radio RA_debida  "¿Salvia debe activar rutas en el marco de la debida diligencia?"*  Sí | No
│       │   └── [v-if Agreement === true && DueDiligence === true]  Checkbox × 5 "Seleccione los canales utilizados para la activación o gestión de la Ruta de Atención"* (Notificación SALVIA · Articulación institucional telefónica o por mensajería · Denuncia fácil · Mi Salvia · Ninguna, excluyente)  +  <textarea> "¿Qué gestiones se realizaron desde Salvia para la activación de ruta?"*
│       ├── [v-if EntityBranchID && faltan obligatorias]  "Falta: …"  ·  [v-if errorEntidad]  texto de error
│       └── BtnAgregarEntidad  "+ Siguiente entidad"  :disabled="!EntityBranchID || entidadFaltantes().length"  → agregarEntidad()
├── Radio  GS_17  "¿Se identificaron barreras? Sí | No"  (obligatoria, 1er contacto)
│   └── [v-if GS_17 === 'true']  Nota  "Registra cada barrera con el botón «+ Agregar barrera a esta entidad» de la tarjeta de la entidad donde ocurrió." (+ " Primero agrega la entidad." si no hay entidades)
└── BtnGuardar  "Guardar sección"  → guardarRutaAtencion()  +  Mensaje ${mensajeGuardado}  (rojo si mensajeGuardadoEsError)
```

## Estados y reglas

| Elemento | Regla |
|---|---|
| Qué entidades se ven | Solo las de tipo `acudida` (`entityMeta[sede].kind`; en datos anteriores, `solicitud` que no empieza por «Canal de activación:») |
| Obligatorias al agregar | Radicado, fecha de reporte, otra información, al menos una acción (completada o pendiente) si la sede tiene acciones, y —con acciones pendientes— acuerdo, debida diligencia y, con «Sí» y «Sí», canales y gestiones. **No** entran en Finalizar |
| Borrador | Con el Caso en Borrador las entidades y barreras se guardan en el JSON y se materializan al pasar a Iniciado (REG-38\*) |
| Solo lectura (D20) | Con `entitiesReadOnly`: sin formulario de entidad, sin «Eliminar» y sin «+ Agregar barrera»; `GS_14` y `GS_17` siguen respondibles (Finalizar las exige) |
| Finalizar | `GS_14` y `GS_17` siguen siendo obligatorias y, si faltan, se informan como módulo GS («Plan de acción Salvia») |
