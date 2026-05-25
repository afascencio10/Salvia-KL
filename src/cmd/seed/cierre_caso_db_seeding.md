# Documentación de Seeding: Formulario "Cierre de caso"

Este documento contiene la información detallada del formulario de **Cierre de caso** registrado en la base de datos de Salvia el 22 de mayo de 2026.

## 1. Resumen del Script de Seeding
El script de inserción se encuentra ubicado en:
* [seed_cierre_caso.sql](file:///Users/camilo/Documents/KLOUSTR/SALVIA/src/cmd/seed/seed_cierre_caso.sql)

Este script registra la estructura completa del formulario (Formulario, Sección, Preguntas, Opciones de Selección y Condiciones de Visibilidad) utilizando la estructura de base de datos para formularios dinámicos en el esquema `salvia`.

---

## 2. Identificadores (UUIDs) Registrados en la Base de Datos

A continuación se listan todos los UUIDs generados y guardados en sus respectivas tablas:

### A. Formulario (`salvia.form`)
* **ID:** `da8423ab-1a8c-47db-96b7-d10496df571a`
* **Nombre:** `Cierre de caso`
* **Descripción:** `Formulario para registrar el cierre de casos en la plataforma Salvia`
* **Estado:** `active`

### B. Sección del Formulario (`salvia.form_section`)
* **ID:** `a1d83099-ca27-4ca9-baf0-cab92d51fa75`
* **Formulario ID:** `da8423ab-1a8c-47db-96b7-d10496df571a`
* **Nombre:** `Cierre del caso`
* **Descripción:** `NULL` (removida ya que no es necesaria)
* **Orden:** `1`

### C. Preguntas (`salvia.question`)

| Pregunta | Tipo de Pregunta | UUID de la Pregunta | Obligatoria | Orden |
| :--- | :--- | :--- | :---: | :---: |
| **Motivo del cierre** | `single` | `d2c6e1af-651f-43b4-8e61-aecafd07443d` | Sí | 1 |
| **Describa la causa del cierre** | `text` | `90375500-a316-4cf5-b7ec-f5c402da92c2` | Sí | 2 |
| **¿Realizó acciones institucionales por el cierre?** | `boolean` | `4a7d0110-b06d-4569-9572-cd0a5e9ef2c1` | Sí | 3 |

### D. Opciones para "Motivo del cierre" (`salvia.option`)
Estas opciones están enlazadas a la pregunta `d2c6e1af-651f-43b4-8e61-aecafd07443d`:

* **Pérdida de contacto**
  * **ID:** `d8abdb6b-6f86-4b39-801e-1f81299eb114`
  * **Valor:** `perdida_contacto`
  * **Orden:** `1`
* **Solicitud expresa de la ciudadana de finalizar el proceso**
  * **ID:** `75c70f9e-415e-4aa8-8bfd-5b35e9067097`
  * **Valor:** `solicitud_ciudadana`
  * **Orden:** `2`
* **Cumplimiento del plan de atención**
  * **ID:** `463c0026-aab6-4f79-85a6-8c4d8d375b76`
  * **Valor:** `cumplimiento_plan`
  * **Orden:** `3`
* **No corresponde al ámbito, población o naturaleza de la atención**
  * **ID:** `952f6897-3384-4741-869c-ba5ff00f1998`
  * **Valor:** `no_corresponde`
  * **Orden:** `4`
* **Otro ¿cuál?**
  * **ID:** `7ca2d32c-8159-4c74-bd3e-e13353805d9b`
  * **Valor:** `otro`
  * **Orden:** `5`

### E. Condiciones de Visibilidad (`salvia.visibility_condition`)
No se definen condiciones de visibilidad para estas preguntas, todas se muestran incondicionalmente.

---

## 3. Instrucciones de Uso y Referencia

Cuando implementes el formulario en el frontend o backend, puedes referenciar la estructura del formulario utilizando directamente el Form ID: `da8423ab-1a8c-47db-96b7-d10496df571a`. 

El motor de formularios dinámicos recuperará esta estructura, sus preguntas, opciones y lógica de visibilidad condicional automáticamente.
