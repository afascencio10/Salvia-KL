━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando carga la pantalla de Registro
   Tipo: Lifecycle
   Estado: 🟡 Implementado con cambios — sincronizado con el código 2026-10-31
   Funciones: RegistroV3GET() (Go) → data() · created() · mounted() de RegistroV3App (registro-v3-app.js)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  sesión:  usuario autenticado y su rol   → cookie de sesión
}
(sin parámetros de URL: la pantalla no lee ningún `icode`)

PASO 1 — El backend valida el acceso y renderiza (GET /salvia/registro-v3)
  CheckPermission("manage_registro_v3", rol)  → permitido: op, sv
  SI no tiene permiso: responde 401 y TERMINA
  SetHeaderNoCache
  RenderTemplate("registro_v3", { windowTitle: "Registro V3", currentUser: nombres + apellidos, currentRole, locale, lang, menu })

PASO 2 — Inicializar el estado (cliente, data())
  etapa = 'habeas'; buscar = { tipoDocumento: '', documento: '', telefono: '' }; resultadosBusqueda = []
  icode = null; resumenCaso = null; answers = {}; entidades = []; barreras = []; enumCatalog = {}
  triadaCompleta = false; moduloActivo = 'CI'; tzResultado = null
  created(): indexa camposPorId con los catálogos de los 7 módulos + GS § B

PASO 3 — mounted()
  cargarEnumCatalog():       GET /api/v1/victim-cases-v3/enums → enumCatalog   // opciones {code, icode, label} por categoría
  cargarGeoDepartamentos():  GET /api/v1/locations/departments → geoDepartamentos   // cascada de GS § C
  Registra window 'beforeunload' (respaldo de guardado, ver flow-E14)
  La página llama RegistroV3App.mount("#app"), fija document.title y muestra #app (display:block)

PASO 4 — Renderizar
  Vue evalúa etapa → muestra el bloque de Habeas data

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión | Paso afectado |
|---|---|
| La planeación permitía llegar con un `icode` en la URL (p. ej. desde el Expediente de víctima) y saltar al Formulario. No está implementado: un Caso existente solo se abre desde el Cabezote ("Abrir") | PASO 2 |
