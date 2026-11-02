# QA — M-10 («Registro de caso») y M-12 («Ir al borrador del caso») (WP-1 T7, REP-03)

Ejecutado 2026-11-02 con `requests`, `psql` y Playwright contra la app del sandbox. El estado «caso terminado» se simula por SQL (`victim_case_status = 'te'`); el fallo interno (C7) renombrando temporalmente el código `n` del catálogo `yes_no`.

| Caso | Qué se prueba | Resultado esperado | Método |
|---|---|---|---|
| C1 | Víctima con contacto exitoso | 200 con `caseIcode` y `redirect`; caso en Borrador (`bo`) con nombres, teléfono, hechos y ajuste; `agent_id` vacío; `victim_case_victim_contact` = id del reporte; documento vacío; el reporte NO queda cerrado; evento «Borrador de caso iniciado» con ambos ids | API |
| C2 | Tercero | se prellenan los datos de la persona afectada, no los del reportante | API |
| C3 | Familiar y Hombre | 409 «Por ahora no hay procesos de atención para este tipo de reportes»; no se crea nada | API |
| C4 | Sin intento exitoso, reporte cerrado, reporte inexistente | 409, 409 y 404 | API |
| C5 | Segunda llamada (doble clic, también simultánea) | 409 con el `caseIcode` existente; queda un solo caso vinculado y un solo evento | API |
| C6 | Rol et, ro o sin sesión | 401, y no se crea nada; sv también puede convertir | API |
| C7 | Atomicidad (catálogo `yes_no` «No» renombrado) | 500; no queda caso, usuario, vínculo ni evento; restaurado el catálogo funciona | API |
| C8 | `GET /api/v1/victim-cases-v3/{caseIcode}` | CI_2, CI_4, HV_1 y CI_8 en `answeredKeys` (CI_7 no) | API |
| C9 | Con borrador iniciado: registrar intento y cerrar | 409 en ambos («el reporte tiene un borrador de caso iniciado»); el detalle trae `draftCaseIcode`, `canManage=false`, `canConvert=false` | API |
| C10 | El caso vinculado pasa a Terminado (`te`) | el detalle devuelve cierre `ap` «Válido - Apertura de Caso Salvia»; intentos y nueva conversión → 409; el listado lo muestra «Cerrado — válido» | API |
| C11 | Listado y Mis Reportes | etiqueta «Borrador iniciado» y filtro `estado_reporte=borrador_iniciado`; `cerrado_valido` incluye el caso completo; `asignado` y `sin_asignar` excluyen los reportes con borrador | API + UI |
| C12 | Modal | botón «Registro de caso» y confirmación solo en `postContacto` y tipos Víctima o Tercero; «Cancelar» no crea nada; aviso para Familiar y Hombre | UI |
| C13 | Modal con borrador | paso «Borrador de caso iniciado» con «Ir al borrador del caso» (sin aviso de solo lectura) que abre V3 en el formulario | UI |
| C14 | V3 con `?caso=` | abre el formulario en CI con los datos; sin parámetro (o con un valor sin forma de UUID) el flujo normal desde la autorización | UI |
| C15 | Reabrir el reporte con el caso completo | paso «cerrado» con «Válido - Apertura de Caso Salvia» | UI |
| C16 | Dos borradores sin documento | guardar CI en uno no pisa el teléfono del otro (el documento queda vacío); al completar CI_5 y CI_6 se guardan | API |
