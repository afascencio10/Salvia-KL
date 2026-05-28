# Registro de Caso — Index

Pantalla: `set_victim_case.html`  
Ruta: `/salvia/casos/nuevo` (o `/salvia/casos/:id/nuevo` desde contacto previo)

## Archivos de esta pantalla

| Archivo | Descripción |
|---|---|
| [`registro-caso-interface.md`](registro-caso-interface.md) | Árbol de interfaz — estructura visual y condiciones de render |

### Changelogs

| Archivo | Descripción |
|---|---|
| [`changelogMay2025.md`](changelogMay2025.md) | Cambios realizados en mayo 2025 |

### Resumen de eventos

| # | Evento | Tipo |
|---|---|---|
| 1 | Cuando carga la pantalla | Lifecycle |
| 2 | Cuando selecciona autorizacion de datos personales | User Interaction |
| 3 | Cuando cambia departamento de residencia (cadena 3) | User Interaction |
| 4 | Cuando cambia ciudad de residencia (cadena 3) | User Interaction |
| 5 | Cuando cambia municipio de residencia (cadena 3) | User Interaction |
| 6 | Cuando cambia departamento de hechos (cadena 2) | User Interaction |
| 7 | Cuando cambia ciudad de hechos (cadena 2) | User Interaction |
| 8 | Cuando cambia municipio de hechos (cadena 2) | User Interaction |
| 9 | Cuando cambia departamento de atencion (cadena 1) | User Interaction |
| 10 | Cuando cambia ciudad de atencion (cadena 1) | User Interaction |
| 11 | Cuando cambia municipio de atencion (cadena 1) | User Interaction |
| 12 | Cuando cambia tipo de violencia experimentada | User Interaction |
| 13 | Cuando cambia ambito de la violencia | User Interaction |
| 14 | Cuando cambia proximidad con el agresor principal | User Interaction |
| 15 | Cuando cambia relacion con el presunto agresor | User Interaction |
| 16 | Cuando cambia cualquier pregunta comun del tamizaje | User Interaction |
| 17 | Cuando cambia cualquier pregunta especifica del tamizaje | User Interaction |
| 18 | [Cuando presiona boton "Guardar"](Flujos/registro-caso-flujo-guardar.md) | User Interaction |
| 19 | Cuando presiona boton "Volver" | User Interaction |
| 20 | Cuando presiona boton "Finalizar" | User Interaction |
| 21 | Cuando cambian los errores de validacion | Lifecycle |

---

## Eventos

---

**Evento:** Cuando carga la pantalla  
**Tipo de trigger:** Lifecycle  
**Descripcion:** Se ejecuta al navegar a la ruta. El backend (VictimCasePOST_GET) carga todos los enums del formulario, departamentos, reglas de navegacion y datos pre-existentes si viene de un contacto (:id). Vue monta el componente: setea document.title, hace visible el #app, solicita geolocalizacion del navegador (latitude/longitude), y si hay datos previos de departamento/ciudad/municipio los pre-selecciona en las 3 cadenas de selects encadenados (atencion, hechos, residencia).  
**Requerido:** Si

---

**Evento:** Cuando selecciona autorizacion de datos personales  
**Tipo de trigger:** User Interaction  
**Descripcion:** El usuario selecciona "Si" o "No" en el select authorizationAnswer. Si selecciona "Si" (code=='y') se renderiza todo el formulario principal. Si selecciona "No" o no ha seleccionado, solo se muestra el bloque de autorizacion con el texto legal. No ejecuta logica adicional, es un v-if reactivo de Vue.  
**Requerido:** Si

---

**Evento:** Cuando cambia departamento de residencia (cadena 3)  
**Tipo de trigger:** User Interaction  
**Descripcion:** El usuario selecciona un departamento en el select de residencia. Ejecuta changeDepartment3: limpia ciudades3, municipios3 y las selecciones previas. Llama a la API GET /seguridad/ciudades/{deptoCode}/departamento para cargar las ciudades del departamento seleccionado. Si status 200, puebla el select de ciudades3.  
**Requerido:** Si

---

**Evento:** Cuando cambia ciudad de residencia (cadena 3)  
**Tipo de trigger:** User Interaction  
**Descripcion:** El usuario selecciona una ciudad en el select de residencia. Ejecuta changeCity3: limpia municipios3. Llama a la API GET /seguridad/municipios/{cityCode}/codigo_ciudad para cargar los municipios de esa ciudad. Si status 200, puebla el select de municipios3 (towns3).  
**Requerido:** Si

---

**Evento:** Cuando cambia municipio de residencia (cadena 3)  
**Tipo de trigger:** User Interaction  
**Descripcion:** El usuario selecciona un municipio en el select de residencia. Ejecuta changeTown3: actualmente vacio (no ejecuta logica adicional). El v-model ya vincula el valor a form2.residenceTownCode.  
**Requerido:** Si

---

**Evento:** Cuando cambia departamento de hechos (cadena 2)  
**Tipo de trigger:** User Interaction  
**Descripcion:** El usuario selecciona un departamento en la seccion de hechos. Ejecuta changeDepartment2: limpia ciudades2, municipios2, city2Selected y form2.factsTownCode. Llama a la API para cargar ciudades del departamento. Si status 200, puebla el select de ciudades2.  
**Requerido:** Si

---

**Evento:** Cuando cambia ciudad de hechos (cadena 2)  
**Tipo de trigger:** User Interaction  
**Descripcion:** El usuario selecciona una ciudad en la seccion de hechos. Ejecuta changeCity2: limpia municipios2. Llama a la API para cargar municipios de esa ciudad. Si status 200, puebla el select de municipios2 (towns2).  
**Requerido:** Si

---

**Evento:** Cuando cambia municipio de hechos (cadena 2)  
**Tipo de trigger:** User Interaction  
**Descripcion:** El usuario selecciona un municipio en la seccion de hechos. Ejecuta changeTown2: actualmente vacio. El v-model vincula el valor a form2.factsTownCode.  
**Requerido:** Si

---

**Evento:** Cuando cambia departamento de atencion (cadena 1)  
**Tipo de trigger:** User Interaction  
**Descripcion:** El usuario selecciona un departamento en la seccion de lugar de atencion. Ejecuta changeDepartment: limpia ciudades, municipios, entitiesByMomentAndSector (la grilla de ruta) y townCode. Llama a la API para cargar ciudades. Si status 200, puebla el select de ciudades.  
**Requerido:** Si

---

**Evento:** Cuando cambia ciudad de atencion (cadena 1)  
**Tipo de trigger:** User Interaction  
**Descripcion:** El usuario selecciona una ciudad en la seccion de lugar de atencion. Ejecuta changeCity: limpia municipios y entitiesByMomentAndSector. Llama a la API para cargar municipios. Si status 200, puebla el select de municipios.  
**Requerido:** Si

---

**Evento:** Cuando cambia municipio de atencion (cadena 1)  
**Tipo de trigger:** User Interaction  
**Descripcion:** El usuario selecciona un municipio en la seccion de lugar de atencion. Ejecuta changeTown: llama a la API GET /salvia/sede/{townCode}/by/por_momentos para cargar las sedes (entity branches) organizadas por momento y sector. Si status 200, puebla entitiesByMomentAndSector lo que renderiza la seccion de "Asignacion de Ruta" con los selects de sedes por momento/sector.  
**Requerido:** Si

---

**Evento:** Cuando cambia tipo de violencia experimentada  
**Tipo de trigger:** User Interaction  
**Descripcion:** El usuario selecciona/deselecciona opciones en el multi-select typeViolenceExperienced. Ejecuta violenceTypeChanged: limpia los subtipos seleccionados (subtypeViolenceExperienced = []) y reconstruye la lista de subtipos disponibles (violenceSubtypes) concatenando los enums de cada tipo seleccionado (victim_case_form2_subtype_violence_experienced_{code}). Esto permite que los subtipos sean dinamicos segun los tipos elegidos.  
**Requerido:** Si

---

**Evento:** Cuando cambia ambito de la violencia  
**Tipo de trigger:** User Interaction  
**Descripcion:** El usuario selecciona/deselecciona opciones en el multi-select scopeOfViolence. Ejecuta checkWorkplaceSectorOccurrence: recorre las opciones seleccionadas buscando code=="al" (ambito laboral). Si lo encuentra, activa workplaceSectorOccurrence=true lo que muestra el select condicional de sector laboral de ocurrencia. Si no, lo oculta.  
**Requerido:** Si

---

**Evento:** Cuando cambia proximidad con el agresor principal  
**Tipo de trigger:** User Interaction  
**Descripcion:** El usuario selecciona una opcion en proximityPrincipalAggressor. Ejecuta proximityWithAggressorChanged: determina si el agresor es persona conocida/pareja seteando partnerKnown=true cuando code es "pc" (persona conocida) o "pn" (pareja). Esto controla la visibilidad del campo "dependencia economica" (v-if partnerKnown).  
**Requerido:** Si

---

**Evento:** Cuando cambia relacion con el presunto agresor  
**Tipo de trigger:** User Interaction  
**Descripcion:** El usuario selecciona una opcion en relationshipWithPresumedAggressor. Ejecuta relationWithAggressorChanged: determina wasPartner=true si code es "pi" (pareja intima) o "ex" (ex-pareja). Luego llama updateTamizajeScore para recalcular el puntaje de riesgo con la bateria correcta (pareja vs no-pareja). Esto controla que bloque de preguntas de tamizaje se muestra (18 preguntas pareja o 14 no-pareja).  
**Requerido:** Si

---

**Evento:** Cuando cambia cualquier pregunta comun del tamizaje (6 preguntas)  
**Tipo de trigger:** User Interaction  
**Descripcion:** El usuario responde Si/No a una de las 6 preguntas comunes del tamizaje (aggressorViolencePhysicalIncrease, aggressorWeaponUsed, aggressorThreatKill, aggressorPursuesSpiesDestroys, aggressorCapableOfKilling, aggressorHasAccessToWeapons). Cada una tiene @change que llama updateTamizajeAggressor: evalua si al menos una de las 6 es "Si" (tamizajeAggressorCheck=true) para activar el bloque de preguntas especificas (pareja o no-pareja). Luego llama updateTamizajeScore para recalcular puntaje y nivel de riesgo.  
**Requerido:** Si

---

**Evento:** Cuando cambia cualquier pregunta especifica del tamizaje (pareja: 18, no-pareja: 14)  
**Tipo de trigger:** User Interaction  
**Descripcion:** El usuario responde Si/No a cualquiera de las preguntas especificas del tamizaje. Cada una tiene @change que llama updateTamizajeScore. Esta funcion suma todos los "Si" de las 6 comunes + las especificas visibles, calcula riskScore y determina riskLevel segun la tabla de umbrales (pareja: 0-4=bajo, 5-8=moderado, 9-15=alto, 16-24=extremo; no-pareja: 0-2=bajo, 3-5=moderado, 6-8=alto, 9-20=extremo). Actualiza form2.riskScore y form2.riskLevel que se muestran en el RiskBadge con color dinamico.  
**Requerido:** Si

---

📄 [Ver flujo → registro-caso-flujo-guardar.md](Flujos/registro-caso-flujo-guardar.md)

**Evento:** Cuando presiona boton "Guardar"  
**Tipo de trigger:** User Interaction  
**Descripcion:** El usuario hace click en el boton "Guardar" (btn-success). Ejecuta submit('save'): agrega las coordenadas GPS (latitude/longitude) al objeto victimCase. Luego hace POST a /salvia/casos (o /salvia/casos/:id si viene de contacto) enviando {victimCase: data.victimCase} como JSON. El backend (SetVictimCase) valida todos los campos, crea el usuario de la victima, inserta el caso, genera el calendario de seguimientos y calcula la ruta. Si status 200: oculta el overlay de exito y muestra finishedOverlay con las credenciales (login/pass). Si status 400: puebla this.errors con los errores campo por campo que se muestran inline. Si otro status: muestra un alert generico.  
**Requerido:** Si

---

**Evento:** Cuando presiona boton "Volver"  
**Tipo de trigger:** User Interaction  
**Descripcion:** El usuario hace click en el boton "Volver" (btn-warning). Ejecuta submit('back'): verifica si existe nav_rules["default"] y redirige a esa URL con location.assign. No realiza POST ni guarda datos.  
**Requerido:** Si

---

**Evento:** Cuando presiona boton "Finalizar"  
**Tipo de trigger:** User Interaction  
**Descripcion:** El usuario hace click en "Finalizar" dentro del finishedOverlay (que se muestra tras guardar exitosamente). Ejecuta finished(): redirige a nav_rules["default"] o nav_rules[action] con location.assign. Cierra el flujo de registro.  
**Requerido:** Si

---

**Evento:** Cuando cambian los errores de validacion (watcher)  
**Tipo de trigger:** Lifecycle  
**Descripcion:** Se ejecuta automaticamente cada vez que el objeto errors cambia (deep watcher de Vue). Recorre todos los errores por entidad/campo, busca el input correspondiente en el DOM por data-entity y data-field, obtiene su label, y construye una lista de mensajes de error que inyecta en el elemento #global-errors-live (aria-live="assertive") para que el lector de pantallas lo anuncie. Soporte de accesibilidad.  
**Requerido:** Si

---

## Validacion de completitud

- [x] Toda accion del usuario puede ser manejada
- [x] La carga inicial de datos esta cubierta (Cuando carga la pantalla)
- [x] El envio del formulario esta incluido (Cuando presiona boton "Guardar")
- [x] Los selects encadenados estan cubiertos (3 cadenas x 3 eventos = 9 eventos de cambio)
- [x] La logica condicional del tamizaje esta cubierta (6 comunes + especificas + relacion/proximidad)
- [x] Los campos condicionales estan cubiertos (tipo violencia → subtipos, ambito → sector laboral)
- [x] La navegacion esta cubierta (Volver, Finalizar)
- [x] La accesibilidad esta cubierta (watcher de errores → aria-live)
- [ ] No hay actualizaciones en tiempo real (sockets, notificaciones)
- [ ] No hay tareas programadas (backend cron)
- [ ] No hay webhooks (triggers externos)

