---
okf_version: "1.0"
type: Service_Logic
title: "Servicio: Generación del Excel consolidado de contactos"
description: "Consulta los reportes (victim_contact) creados en un rango de fechas junto con su Form1/Form2 y arma un Excel de una sola hoja plana con cabeceras amigables usando excelize."
owner: "@backend-squad"
status: active
tags: [backend, service, reportes, contactos]
dependencies:
  - Endpoint: 3-features/ReporteConsolidadoContactos/backend/endpoints/reporte-contactos-consolidado-post.md
last_updated: "2026-07-17"
code_refs:
  - src/salvia/service/report_service.go
  - src/salvia/service/contacts_report_excel.go
  - src/salvia/service/contacts_report_service_test.go
  - src/internal/repository/report_repository.go
---

# Servicio: Generación del Excel consolidado de contactos

## Responsabilidades

1. Recibir el rango `[start_date, end_date]` ya validado por el facade.
2. Consultar los **reportes creados** en el rango (por `victim_contact_creation_date`, `America/Bogota`), en **todos los estados**, junto con su `VictimContactForm1` **o** `VictimContactForm2` (un reporte tiene exactamente uno de los dos) y los enums de Form2.
3. Resolver el **nombre del municipio** (`town_name`) a partir de `victim_contact_form1_town_code` contra la tabla `towns` (patrón existente en `src/internal/repository/location_repository.go`).
4. Si el rango no contiene **ningún** reporte, devolver `ErrNoData` (el endpoint lo mapea a `400 VALIDATION_FAILED`) — **no** se genera Excel vacío.
5. Construir el libro **Excel** con `github.com/xuri/excelize/v2` (una sola hoja) y devolverlo como stream para que el endpoint lo escriba en la respuesta.

## Estructura del libro — hoja única `Reportes`

- **Grano**: 1 fila por reporte (`victim_contact`).
- La primera fila es de **encabezados** (estilo del patrón `createHeaderStyle`/`writeHeaders` de `followups_report_excel.go`), con ancho de columnas auto-ajustado.
- La columna `tipo_formulario_reporte` diferencia cada reporte; las columnas `f1_*` quedan **vacías** si el reporte es de Formulario 2, y las `f2_*` vacías si es de Formulario 1.

### Columnas (orden y origen)

> [!NOTE]
> Ajuste 2026-07-21 (aprobado por el usuario): se eliminaron `contacto_usuario_app` (pedido explícito), `f2_detalle_tipo_reporte` y `f2_respuesta_autorizacion`. Las dos últimas salían **siempre vacías** — verificado contra la BD real: `victim_contact_form2_report_type_details` es NULL en el 100% de los registros y la respuesta de autorización nunca se persiste para contactos (la tabla `rel_victim_case_form2_enums_victim_contact_form2` solo contiene la categoría `victim_case_form2_adjustments_gbv`). Si esos campos se empiezan a guardar, se reagregan aquí y en el builder.

| # | Cabecera | Origen (campo DTO) | Formato amigable |
| :-- | :--- | :--- | :--- |
| 1 | `contacto_id` | `VictimContactId` | Número |
| 2 | `contacto_icode` | `VictimContactICode` | Texto |
| 3 | `contacto_fecha_creacion` | `VictimContactCreationDate` | `dd/mm/yyyy hh:mm` (`America/Bogota`) |
| 4 | `contacto_fecha_actualizacion` | `VictimContactUpdateDate` | `dd/mm/yyyy hh:mm` |
| 5 | `contacto_estado` | `VictimContactStatusDescription` | Descripción legible (no la clave) |
| 6 | `contacto_nombres` | `VictimContactNames` | Texto |
| 7 | `contacto_apellidos` | `VictimContactLastNames` | Texto |
| 8 | `contacto_latitud` | `VictimContactLatitude` | Número decimal |
| 9 | `contacto_longitud` | `VictimContactLongitude` | Número decimal |
| 10 | `tipo_formulario_reporte` | Derivado: existe Form1 → `Reporte propio (Formulario 1)`; existe Form2 → `Reporte de tercero (Formulario 2)` | Texto |
| 11 | `f1_apodo_nombre_identitario` | `VictimContactForm1Nick` | Texto |
| 12 | `f1_tipo_documento` | `VictimContactForm1DocType` | Traducido (patrón `formatDocType`) |
| 13 | `f1_numero_documento` | `VictimContactForm1DocNumber` | Texto |
| 14 | `f1_fecha_nacimiento` | `VictimContactForm1BirthDate` | `dd/mm/yyyy` |
| 15 | `f1_codigo_municipio` | `VictimContactForm1TownCode` | Código DANE crudo |
| 16 | `f1_municipio` | `towns.town_name` (join por `town_code`) | Nombre del municipio |
| 17 | `f1_direccion` | `VictimContactForm1Address` | Texto |
| 18 | `f1_telefono` | `VictimContactForm1Phone` | Texto |
| 19 | `f1_identidad_genero` | `VictimContactForm1GenderIdentity` | Traducido (patrón `formatGender`/Locale) |
| 20 | `f1_orientacion_sexual` | `VictimContactForm1SexualOrientation` | Traducido (Locale) |
| 21 | `f1_procedencia_zona` | `VictimContactForm1Origin` | Traducido (patrón `formatOrigin`) |
| 22 | `f1_ocupacion` | `VictimContactForm1Occupation` | Traducido (Locale) |
| 23 | `f1_otra_ocupacion` | `VictimContactForm1OccupationOther` | Texto |
| 24 | `f1_descripcion_hechos` | `VictimContactForm1FactsDescription` | Texto |
| 25 | `f2_nombre_reportante` | `VictimContactForm2ReporterNames` | Texto |
| 26 | `f2_telefono_reportante` | `VictimContactForm2ReporterPhone` | Texto (sin notación científica) |
| 27 | `f2_telefono_contacto_victima` | `VictimContactForm2VictimColPhone` | Texto (sin notación científica) |
| 28 | `f2_descripcion_hechos` | `VictimContactForm2FactsDescription` | Texto |
| 29 | `f2_mejor_hora_contacto` | `VictimContactForm2BestContactTime` | `hh:mm` (o `dd/mm/yyyy hh:mm` si trae fecha) |
| 30 | `f2_recibira_llamada` | `VictimContactForm2WillReceiveCall` (enum) | Sí / No / descripción traducida |
| 31 | `f2_rol_cuidado` | `VictimContactForm2HasCareRole` (enum) | Sí / No / descripción traducida |
| 32 | `f2_victima_enterada` | `VictimContactForm2VictimAwareOfReport` (enum) | Sí / No / descripción traducida |
| 33 | `f2_tipo_reporte` | `VictimContactForm2ReportType` (enum) | Descripción traducida |
| 34 | `f2_ajustes_gbv` | `VictimContactForm2AdjustmentsGBV` (lista de enums) | Descripciones traducidas separadas por coma |

### Reglas de formato amigable

- **Enums**: se traducen a español usando el Locale del proyecto (`salvia_config.Locale["sp"]`), siguiendo el patrón `translateEnumKey` de `followups_report_excel.go`. Nunca se exporta la clave cruda si existe traducción; si no existe, se exporta la clave tal cual (defensa).
- **Fechas**: `dd/mm/yyyy hh:mm` en zona `America/Bogota`; fechas sin componente horario relevante (`f1_fecha_nacimiento`) van como `dd/mm/yyyy`. Fechas cero/nulas → celda vacía.
- **Booleanos/enums sí-no**: `Sí` / `No` (patrón `formatYesNo`).
- **Teléfonos**: como **texto** para evitar que Excel los muestre en notación científica o les quite ceros a la izquierda.
- **Campos del formulario que no aplica**: celda vacía (no `N/A`, no `-`).

## Interfaz Pública (Contract)

```go
// En ReportService (src/salvia/service/report_service.go) se agrega:
// Reutiliza FollowUpsReportRange { StartDate, EndDate time.Time }

// Devuelve el .xlsx en memoria/stream listo para escribir en la respuesta HTTP.
// Si no hay reportes en el rango devuelve ErrNoData.
GenerateConsolidatedContactsReport(
    ctx context.Context,
    r FollowUpsReportRange,
) (*excelize.File, error)
```

## Errores que Lanza

| Error | Código Catálogo | Cuándo |
| :--- | :--- | :--- |
| `ErrInvalidRange` | `VALIDATION_FAILED` | Rango inválido que llegue sin validar (defensa) |
| `ErrNoData` | `VALIDATION_FAILED` | Cero reportes creados en el rango (no se genera archivo) |
| `ErrReportQuery` | `SYSTEM_INTERNAL_ERROR` | Falla al consultar contactos/formularios/municipios |
| `ErrReportBuild` | `SYSTEM_INTERNAL_ERROR` | Falla al construir/serializar el `.xlsx` |

## Rendimiento y consideraciones

- Consulta **en lote**: un query principal de `victim_contact` en el rango con joins/consultas por lista de IDs para `form1`, `form2`, enums y `towns` (evitar N+1).
- La cota de **≤ 366 días** acota el volumen; la generación es **síncrona** (el request espera el archivo), igual que el reporte de seguimientos.
- No se cachea: cada llamada refleja el estado actual de la BD.
- **PII**: no loguear contenido de filas; solo metadatos (usuario, rango, conteo de filas).

## Dependencias

- `github.com/xuri/excelize/v2` (ya en `go.mod`).
- Patrones reutilizables de `src/salvia/service/followups_report_excel.go` (`createHeaderStyle`, `writeHeaders`, `autoFitColumns`, `translateEnumKey`, `formatYesNo`).
- Resolución de municipios: tabla `towns` (`src/internal/repository/location_repository.go`).
