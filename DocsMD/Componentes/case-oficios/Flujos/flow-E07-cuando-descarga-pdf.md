━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando descarga el PDF del oficio
   Tipo: User Interaction
   Función: downloadPDF(oficio)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  oficio:  objeto EntityLetter (oficioSeleccionado)   → el que está abierto en el modal
}


PASO 1 — Verificar que exista un archivo

  SI !oficio.urlKofax:
    → No hacer nada (el botón "Descargar PDF" no debería estar visible
      si no hay urlKofax — ver GAP)
    → TERMINAR ejecución

  SI oficio.urlKofax existe:
    → CONTINÚA FLUJO GENERAL


PASO 2 — Abrir el archivo en una pestaña nueva

  window.open(oficio.urlKofax, '_blank')
  // urlKofax es una ruta/URL al archivo (Kofax) — no hay descarga vía backend,
  // es un link directo

→ FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión                                                    | Paso afectado |
|--------------------------------------------------------------------------|---------------|
| ¿urlKofax es siempre una URL http(s) accesible desde el navegador, o    | PASO 2        |
| a veces una ruta de red local (ej. "U:/QA/bogota/archivo.pdf", como se  |               |
| vio en datos de prueba de case-task-modal)? Si es ruta local,           |               |
| window.open() no podrá abrirla — el botón fallaría silenciosamente.    |               |
| El botón "Descargar PDF" debería ocultarse/deshabilitarse si            | PASO 1        |
| !oficio.urlKofax (mismo patrón que oficios-list.js).                    |               |
