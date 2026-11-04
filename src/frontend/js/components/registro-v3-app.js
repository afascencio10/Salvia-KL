// Registro V3 — app principal. Sin dinamic-form (decisión A1): Habeas Data
// -> Cabezote (búsqueda por documento y/o teléfono, D7) -> Formulario (7
// módulos CI/DP/CA/HV/AG/TZ/GS) -> Cierre. Usa registro-v3-campos.js como
// catálogo de preguntas y GET /victim-cases-v3/enums como catálogo de
// opciones (icode real de victim_case_form2_enums).
// REG-04 (partes a y b): nombre legible del estado del Caso (victim_case_status). El Cabezote devuelve Casos de
// cualquier estado (no solo los de Registro V3), por eso van los 11 códigos del listado de casos
// (casos-component.js caseStatusLabel). Mapa local: cada parte de REG-04 lleva el suyo (no se comparte helper).
var ESTADOS_CASO_V3 = {
  bo: 'Borrador', in: 'Iniciado', te: 'Terminado', cd: 'Cerrado', ra: 'Activo', is: 'Con novedad',
  ex: 'Vencido', r: 'Por aprobar', fc: 'Recontacto', c: 'Completado', iv: 'Inválido',
};

// REG-37: etiquetas de las obligatorias de GS que se dibujan directo en el HTML (no están en el catálogo
// CAMPOS_GS_*), para decir CUÁLES faltan cuando el servidor rechaza «Finalizar».
var GS_ETIQUETAS_FINALIZAR_V3 = {
  GS_5: '¿Requiere remisión a otros equipos de Salvia?',
  GS_14: '¿Ya se reportó a alguna entidad?',
  GS_16: '¿Salvia debe activar rutas en el marco de la debida diligencia?',
  GS_17: '¿Se identificaron barreras?',
  GS_18: '¿Activar alerta inminente?',
};

// Estilos de la pantalla (REG-05 y REG-01). Se inyectan desde aquí y no desde un <style> del HTML porque Vue ignora
// los <style> que están dentro de un template en el DOM (la página se monta sobre <body>).
(function () {
  if (document.getElementById('registro-v3-estilos')) return;
  var css = [
    // REG-05: mensaje de error / faltantes junto a «Guardar sección» (contraste ≈ 6,5:1 sobre blanco, AA).
    // Los nombres v3-msg-error y mensajeGuardadoEsError son el contrato con WP-4 (cola del bloque GS).
    '.v3-msg-error { color: #b02a37; font-weight: 600; }',
    // Respaldo: el mensaje que sigue al botón «Guardar sección» (también en la cola del bloque GS, que es de WP-4 y
    // aquí no se edita) sale en rojo cuando la bandera está activa. Se puede quitar cuando WP-4 vincule la clase.
    '.v3-con-error button.btn-primary + .ms-2 { color: #b02a37; font-weight: 600; }',
    // REG-01: menú de módulos a la izquierda y contenido a la derecha, sin envolver ni reindentar bloques existentes.
    '.v3-form { display: grid; grid-template-columns: 260px 1fr; column-gap: 1rem; align-items: start; }',
    '.v3-form > * { grid-column: 2; }',
    '.v3-form > .v3-menu { grid-column: 1; grid-row: 1 / span 40; margin-bottom: 1rem; }',
    '.v3-menu .nav-link { display: flex; justify-content: space-between; align-items: center; gap: .5rem; }',
    '.v3-menu .nav-link.v3-bloqueado { color: #6c757d; }',
    '@media (max-width: 767px) { .v3-form { display: block; } }',
  ].join('\n');
  var el = document.createElement('style');
  el.id = 'registro-v3-estilos';
  el.textContent = css;
  document.head.appendChild(el);
})();

var RegistroV3App = Vue.createApp({
  delimiters: ['${', '}'],
  data: function () {
    return {
      etapa: 'habeas', // habeas | formularioCerrado | tipoContacto | cabezote | formulario | cierre
      habeasAceptado: false,
      tipoContactoNoDisponible: false,

      buscar: { tipoDocumento: '', documento: '', telefono: '' },
      resultadosBusqueda: [],
      buscando: false,
      errorBusqueda: '',

      icode: null,
      resumenCaso: null,
      answers: {},
      entidades: [],
      barreras: [],
      // REG-38*: datos de Seguimiento v4 que no tienen columna (GetCase): entityMeta se indexa por
      // entityBranchId y barrierMeta por el id de la barrera; entitiesReadOnly = el Caso ya dejó de
      // ser borrador y la sección de entidades, barreras y remisiones es un snapshot de solo lectura.
      entityMeta: {},
      barrierMeta: {},
      entitiesReadOnly: false,
      // REG-15: etiquetas de las acciones institucionales de cada sede ({ sede: { idObligación: texto } }),
      // para listar en la tarjeta lo que se marcó al registrar la entidad.
      obligacionesPorBranch: {},
      enumCatalog: {},
      triadaCompleta: false,

      moduloActivo: 'CI',
      modulos: [
        { id: 'CI', label: 'Contacto e identificación', campos: CAMPOS_CI_V3 },
        { id: 'DP', label: 'Datos personales', campos: CAMPOS_DP_V3 },
        { id: 'CA', label: 'Contacto de apoyo', campos: CAMPOS_CA_V3 },
        { id: 'HV', label: 'Hechos de violencia', campos: CAMPOS_HV_V3 },
        { id: 'AG', label: 'Características del agresor', campos: CAMPOS_AG_V3 },
        { id: 'TZ', label: 'Tamizaje de riesgo', campos: [] },
        // REG-17: módulo nuevo, nunca bloqueado (moduloBloqueado solo bloquea TZ mientras la tríada esté incompleta).
        { id: 'RA', label: 'Ruta de atención integral', campos: CAMPOS_RA_V3 },
        { id: 'GS', label: 'Plan de acción Salvia', campos: CAMPOS_GS_A_V3 },
      ],
      guardando: false,
      mensajeGuardado: '',
      errorGuardado: '', // REG-11: aviso en rojo cuando un guardado falla (con el campo a corregir si el servidor lo indica)
      // REG-05: true cuando el texto de `mensajeGuardado` es un error o un recordatorio de faltantes -> se pinta con
      // la clase v3-msg-error (contrato con WP-4: nombres exactos; la cola del bloque GS conserva ese vínculo).
      mensajeGuardadoEsError: false,
      finalizando: false, // REG-37: evita el doble envío de «Finalizar registro»
      mensajeAvance: '', // REG-01: aviso breve al pasar solo al siguiente módulo tras guardar
      mensajeModuloBloqueado: '',
      mensajeCierre: '',
      avisoDuplicidad: '',
      errorFinalizar: '',

      // GS § C -- buscador real Departamento/Ciudad/Municipio + Sector ->
      // entity_branch, igual que case-entities.js / Seguimiento (ver
      // registro-v3-pendientes.md § B).
      geoDepartamentos: [],
      entidadCiudades: [], entidadMunicipios: [], entidadesDisponibles: [],
      entidadBuscando: false, obligacionesEntidad: [], obligacionesCargando: false,
      errorEntidad: '',
      // REG-19: «Buscar entidad». entidadTexto es lo que se ve en el campo (inmediato); textoEntidad
      // es la copia con debounce (250 ms) con la que se filtra la lista ya cargada.
      entidadTexto: '', textoEntidad: '', entidadEnfocada: false,
      // GS § D -- barreraParaEntidad es el entityBranchId de la entidad cuyo
      // formulario de "+ Agregar barrera" está abierto (null = ninguno).
      barreraParaEntidad: null,
      // REG-18: municipios del Departamento elegido en el formulario de barrera
      // (lista propia: el buscador de entidades tiene la suya).
      barreraMunicipios: [],
      errorBarrera: '',

      // Catálogos globales (registro-v3-campos.js) que el TEMPLATE referencia
      // directo (no a través de un método) -- BUG encontrado en QA
      // (2026-10-15): Vue, al compilar el template en el navegador, reescribe
      // cualquier identificador que no reconoce como dato/método/prop del
      // componente a _ctx.<nombre> en vez de dejarlo caer al scope global
      // (a diferencia de lo esperado de un `with` normal) -- una referencia
      // directa a una var global en el HTML (ej. "${TAMIZAJE_VPP_V3.label}")
      // revienta el render entero con "Cannot read properties of undefined".
      // Exponerlas aquí como datos del componente resuelve _ctx.X
      // correctamente. Cualquier catálogo nuevo que se use DIRECTO en el
      // HTML (no envuelto en un método) debe agregarse aquí también.
      TAMIZAJE_GENERALES_V3: TAMIZAJE_GENERALES_V3,
      TAMIZAJE_VPP_V3: TAMIZAJE_VPP_V3,
      SECTORES_ENTIDAD_BARRERA_V3: SECTORES_ENTIDAD_BARRERA_V3,
      CANAL_ACTIVACION_V3: CANAL_ACTIVACION_V3,
      ENTIDAD_V4_V3: ENTIDAD_V4_V3,

      // Tamizaje
      tzRevelarResto: false,
      tzResultado: null,
      tzGuardando: false,

      // GS § B remisión
      gsRequiereRemision: false,
      gsJustificacionRemision: '',

      // GS § C entidades
      nuevaEntidad: this.entidadVacia('activacion'),
      // GS § D barreras
      nuevaBarrera: this.barreraVacia(),
    };
  },

  watch: {
    // El formulario de entidad es de un solo tipo por módulo; al cambiar de pestaña se reinicia.
    moduloActivo: function (id) {
      if (id !== 'GS' && id !== 'RA') return;
      var kind = id === 'RA' ? 'acudida' : 'activacion';
      if (this.nuevaEntidad.Kind === kind) return;
      this.nuevaEntidad = this.entidadVacia(kind);
      this.entidadCiudades = []; this.entidadMunicipios = []; this.entidadesDisponibles = []; this.obligacionesEntidad = [];
      this.reiniciarBuscadorEntidad();
      this.errorEntidad = '';
      this.barreraParaEntidad = null;
    },
  },

  computed: {
    // REG-15/17: cada entidad se muestra solo en el módulo de su tipo (acudida en la Ruta de atención
    // integral, activación en el Plan de acción).
    entidadesActivacion: function () {
      var self = this;
      return this.entidades.filter(function (e) { return self.kindDeEntidad(e) === 'activacion'; });
    },
    entidadesAcudidas: function () {
      var self = this;
      return this.entidades.filter(function (e) { return self.kindDeEntidad(e) === 'acudida'; });
    },
    moduloActualObj: function () {
      var self = this;
      return this.modulos.filter(function (m) { return m.id === self.moduloActivo; })[0];
    },
    // REG-19: sedes del municipio y sector ya cargadas, filtradas en el navegador (sin tildes ni
    // mayúsculas, desde la primera letra). Escribir NUNCA llama al servidor.
    entidadesFiltradas: function () {
      var q = this.sinTildes(this.textoEntidad).trim();
      var self = this;
      if (!q) return this.entidadesDisponibles;
      return this.entidadesDisponibles.filter(function (o) { return self.sinTildes(o.name).indexOf(q) >= 0; });
    },
    equiposRemision: function () { return EQUIPOS_REMISION_V3; },
    equiposSeleccionados: function () {
      return (this.answers['GS_6'] || '').split(',').filter(function (v) { return v; });
    },
    // REG-34: ¿ya se respondieron las cinco preguntas iniciales del Tamizaje (TZ_1, 2, 3, 4, 7)? Con todas en
    // «No» no se pregunta el resto y el Tamizaje ya se puede calcular y guardar.
    tzInicialesRespondidas: function () {
      var self = this;
      return TAMIZAJE_GENERALES_V3.filter(function (c) { return c.inicial; }).every(function (c) {
        var v = self.answers[c.id];
        return v === 'true' || v === 'false';
      });
    },
  },

  created: function () {
    this.camposPorId = {};
    var self = this;
    this.modulos.forEach(function (m) {
      m.campos.forEach(function (c) { self.camposPorId[c.id] = c; });
    });
    CAMPOS_GS_B_V3.concat(CAMPOS_GS_CRITERIOS_V3).forEach(function (c) { self.camposPorId[c.id] = c; });
  },

  mounted: function () {
    this.cargarEnumCatalog();
    this.cargarGeoDepartamentos();
    // D11 (flow-E14): respaldo para cuando el agente cierra la pestaña sin
    // presionar "Finalizar llamada". navigator.sendBeacon (no fetch) porque
    // es la API pensada para esto -- garantiza el intento de envío aunque
    // la página se descargue antes de que una promesa normal resuelva (el
    // propio GAP de flow-E14 señalaba que un fetch común es poco confiable
    // acá).
    var self = this;
    window.addEventListener('beforeunload', function () {
      if (self.etapa !== 'formulario') return;
      if (['CI', 'DP', 'CA', 'HV', 'AG', 'GS'].indexOf(self.moduloActivo) === -1) return;
      var payload = {};
      Object.keys(self.answers).forEach(function (k) {
        if (k.indexOf(self.moduloActivo + '_') === 0) payload[k] = self.answers[k];
      });
      var blob = new Blob([JSON.stringify({ answers: payload })], { type: 'application/json' });
      navigator.sendBeacon('/api/v1/victim-cases-v3/' + self.icode + '/seccion/' + self.moduloActivo, blob);
    });
    // WP-1 T7 (E19): si la URL trae ?caso=<icode> (borrador creado desde un reporte), abre el formulario.
    this.abrirDesdeReporte();
  },

  methods: {
    // depto/ciudad/municipio/sector son estado LOCAL del buscador -- nunca se
    // mandan al backend, solo resuelven EntityBranchID (ver
    // buscarEntidadesEntidad). El resto sí se manda, en PascalCase idéntico
    // al campo Go correspondiente (repository.VictimCaseEntityInput).
    // El tipo lo fija el módulo: Plan de acción = 'activacion'; Ruta de atención integral = 'acudida'.
    entidadVacia: function (kind) {
      return {
        Kind: kind || 'activacion', depto: '', ciudad: '', municipio: '', sector: '',
        EntityBranchID: '', FileNumber: '', Solicitud: '', ReportDate: '', ActivationChannels: [], ManagementNotes: '',
        // Solo la entidad acudida con acciones pendientes: acuerdo y debida diligencia propios (D17).
        Agreement: null, DueDiligence: null,
        CompletedActions: [], PendingActions: [], OtherInfo: '',
      };
    },
    // REG-18: la barrera es el formulario de Seguimiento v4 (BARRERA_V4_V3,
    // generado desde el seed). `v` guarda las respuestas por clave de pregunta
    // (arreglo si la pregunta es de selección múltiple, texto en los demás
    // casos); armarPayloadBarrera() las traduce a los campos de
    // repository.VictimCaseBarrierInput. EntityBranchID se hereda de la entidad
    // padre (abrirBarreraParaEntidad).
    barreraVacia: function () {
      var v = {};
      BARRERA_V4_V3.forEach(function (q) {
        if (q.tipo !== 'info') v[q.key] = q.tipo === 'multiple' ? [] : '';
      });
      return { EntityBranchID: '', v: v };
    },

    // ── Habeas data / Cabezote ──────────────────────────────────────────
    // D1/D2 (registro-v3-index.md § Inventario de eventos, "Cuando responde
    // la autorización..."/"Cuando selecciona el tipo de contacto"): antes
    // "Acepto y continuar" saltaba directo a Cabezote, sin la rama "No" ni
    // la etapa de Tipo de contacto.
    aceptarHabeasData: function () {
      this.habeasAceptado = true;
      this.etapa = 'tipoContacto';
    },
    rechazarHabeasData: function () {
      this.etapa = 'formularioCerrado';
    },
    seleccionarTipoContacto: function (tipo) {
      if (tipo === 'victima_caso_nuevo') {
        this.etapa = 'cabezote';
        return;
      }
      // C9: Víctima-Seguimiento, Tercero-Reporte, Familiares (feminicidio) y
      // Hombres abren formularios/pantallas distintas a esta. Esta pantalla
      // no conoce sus URLs reales (no están resueltas en ninguna fuente de
      // planeación -- G-V3-11 para Familiares/Hombres) -- se avisa en vez
      // de inventar un enlace.
      this.tipoContactoNoDisponible = true;
    },

    ejecutarBusqueda: function () {
      if (!this.buscar.documento && !this.buscar.telefono) {
        this.errorBusqueda = 'Digite al menos documento o teléfono.';
        return;
      }
      this.errorBusqueda = '';
      this.buscando = true;
      var self = this;
      var qs = new URLSearchParams();
      if (this.buscar.tipoDocumento) qs.set('tipo_documento', this.buscar.tipoDocumento);
      if (this.buscar.documento) qs.set('documento', this.buscar.documento);
      if (this.buscar.telefono) qs.set('telefono', this.buscar.telefono);
      fetch('/api/v1/victim-cases-v3/buscar?' + qs.toString())
        .then(function (r) { return self.leerRespuesta(r); })
        .then(function (res) {
          if (!res.ok) { self.errorBusqueda = (res.data && res.data.error) || 'Error al buscar.'; return; }
          self.resultadosBusqueda = res.data.registros || [];
        })
        .catch(function () { self.errorBusqueda = 'Error al buscar.'; })
        .finally(function () { self.buscando = false; });
    },

    // REG-11: lee la respuesta sin lanzar si el cuerpo no es JSON -> { ok, status, data }.
    leerRespuesta: function (r) {
      return r.text().then(function (t) {
        var data = {};
        try { data = t ? JSON.parse(t) : {}; } catch (e) { data = {}; }
        return { ok: r.ok, status: r.status, data: data };
      });
    },

    crearNuevoRegistro: function () {
      var self = this;
      this.errorBusqueda = '';
      fetch('/api/v1/victim-cases-v3/draft', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          tipo_documento: this.buscar.tipoDocumento,
          documento: this.buscar.documento,
          telefono: this.buscar.telefono,
        }),
      })
        .then(function (r) { return self.leerRespuesta(r); })
        .then(function (res) {
          if (res.ok && res.data.icode) { self.abrirCaso(res.data.icode); return; }
          // REG-11: antes, si el borrador no se creaba, no pasaba nada y la persona no sabía por qué.
          self.errorBusqueda = 'No se pudo crear el registro. Intenta de nuevo; si el problema continúa, avisa a soporte.';
        })
        .catch(function () {
          self.errorBusqueda = 'No se pudo crear el registro: no hay conexión con el servidor. Intenta de nuevo.';
        });
    },

    // WP-1 T7 (E19): abre el registro creado desde un reporte (?caso=<icode>), directo en el
    // formulario: el consentimiento ya se tomó en el reporte V2. Sin parámetro (o si no tiene forma
    // de UUID) no hace nada y el flujo normal arranca en la etapa 'habeas'.
    abrirDesdeReporte: function () {
      var icode = new URLSearchParams(window.location.search).get('caso');
      if (!icode || !/^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i.test(icode)) return;
      this.abrirCaso(icode);
    },

    abrirCaso: function (icode) {
      var self = this;
      this.errorBusqueda = '';
      fetch('/api/v1/victim-cases-v3/' + icode)
        .then(function (r) { return self.leerRespuesta(r); })
        .then(function (res) {
          // REG-11: si el Caso no carga, la pantalla avisa en vez de quedarse en el Cabezote sin decir nada.
          if (!res.ok || !res.data.summary) {
            self.errorBusqueda = 'No se pudo abrir el registro. Intenta de nuevo; si el problema continúa, avisa a soporte.';
            return;
          }
          var data = res.data;
          self.icode = icode;
          self.errorGuardado = '';
          self.resumenCaso = data.summary || null;
          // Solo se repueblan las keys en answeredKeys -- el resto de
          // 'answers' son valores dummy (D8) que nunca respondió el
          // agente y no deben mostrarse como si lo hubiera hecho (ver
          // hallazgo H3, registro-v3-qa-plan.md).
          var todas = data.answers || {};
          var respondidas = data.answeredKeys || [];
          var filtradas = {};
          respondidas.forEach(function (k) {
            if (todas[k] !== undefined) filtradas[k] = todas[k];
          });
          self.answers = filtradas;
          // REG-34: el resto del Tamizaje se muestra si alguna de las cinco iniciales quedó en «Sí».
          self.tzRevelarResto = TAMIZAJE_GENERALES_V3.some(function (c) { return c.inicial && filtradas[c.id] === 'true'; });
          self.entidades = data.entities || [];
          self.barreras = data.barriers || [];
          self.entityMeta = data.entityMeta || {};
          self.barrierMeta = data.barrierMeta || {};
          self.entitiesReadOnly = !!data.entitiesReadOnly;
          self.obligacionesPorBranch = {};
          self.cargarEtiquetasAcciones();
          // BUG encontrado en QA (2026-10-15): tzResultado solo se llenaba
          // con la respuesta directa de guardarTamizaje() -- en memoria,
          // nunca releído del backend. Cualquier reapertura del Caso (o
          // Finalizar en otra sesión distinta a la que guardó el Tamizaje)
          // mostraba "Nivel de riesgo: —" en el Cierre aunque el resultado
          // real siguiera guardado en BD (GetCase ya lo manda ahora).
          self.tzResultado = data.tzResultado || null;
          // A2: antes solo se marcaba tras un guardado en vivo -- un Caso ya
          // completo se reabría con el sidebar bloqueado otra vez. GetCase
          // ahora manda este dato calculado con el mismo criterio que
          // SaveSection (ver victim_case_v3_service.go, isTriadaCompleteV3).
          self.triadaCompleta = !!data.triadaCompleta;
          self.etapa = 'formulario';
          self.moduloActivo = 'CI';
        })
        .catch(function () {
          self.errorBusqueda = 'No se pudo abrir el registro: no hay conexión con el servidor. Intenta de nuevo.';
        });
    },

    cargarEnumCatalog: function () {
      var self = this;
      fetch('/api/v1/victim-cases-v3/enums')
        .then(function (r) { return r.json(); })
        .then(function (data) { self.enumCatalog = data || {}; });
    },

    // ── Catálogo de opciones / etiquetas ─────────────────────────────────
    campoPorId: function (id) { return this.camposPorId[id]; },

    opcionesPara: function (campo) {
      if (!campo) return [];
      if (campo.opciones) {
        return campo.opciones.map(function (o) { return { code: o, icode: o, label: o }; });
      }
      if (campo.subtipoDe) {
        var tiposSeleccionados = (this.answers[campo.subtipoDe] || '').split(',').filter(function (v) { return v; });
        var tipoOpts = this.opcionesPara(this.camposPorId[campo.subtipoDe]);
        var categorias = [];
        var self = this;
        tiposSeleccionados.forEach(function (icode) {
          var opt = tipoOpts.filter(function (o) { return o.icode === icode; })[0];
          var suf = opt && TIPO_A_SUFIJO_SUBTIPO_V3[opt.code];
          if (suf) categorias.push('victim_case_form2_subtype_violence_experienced_' + suf);
        });
        var out = [];
        categorias.forEach(function (cat) { out = out.concat(self.enumCatalog[cat] || []); });
        return out;
      }
      if (campo.enumCategories) {
        var merged = [];
        var seen = {};
        var self2 = this;
        campo.enumCategories.forEach(function (cat) {
          (self2.enumCatalog[cat] || []).forEach(function (o) {
            if (!seen[o.icode]) { seen[o.icode] = true; merged.push(o); }
          });
        });
        return merged;
      }
      if (campo.enumCategory) {
        return this.enumCatalog[campo.enumCategory] || [];
      }
      return [];
    },

    etiqueta: function (campo, opt) {
      if (!campo || !opt) return '';
      var overrides = campo.enumCategory && ENUM_LABELS_OVERRIDE_V3[campo.enumCategory];
      if (overrides && overrides[opt.code]) return overrides[opt.code];
      if (campo.subtipoDe || campo.enumCategories) {
        // Solo se busca override en las categorías de ESTE campo (enumCategories). Los subtipos
        // (subtipoDe) vienen de victim_case_form2_subtype_violence_experienced_*, que no tienen
        // override: su texto es el que entrega /enums. Antes se recorrían todas las categorías y
        // un código repetido en otra (tr, al, es, ah, ai, pr) tomaba la etiqueta ajena (REG-10).
        var cats = campo.enumCategories || [];
        for (var i = 0; i < cats.length; i++) {
          var ov = ENUM_LABELS_OVERRIDE_V3[cats[i]];
          if (ov && ov[opt.code]) return ov[opt.code];
        }
      }
      return opt.label;
    },

    etiquetaSeleccionada: function (fieldId) {
      var campo = this.camposPorId[fieldId];
      var icode = this.answers[fieldId];
      if (!campo || !icode) return '';
      var opts = this.opcionesPara(campo);
      for (var i = 0; i < opts.length; i++) {
        if (opts[i].icode === icode) return this.etiqueta(campo, opts[i]);
      }
      return '';
    },

    multipleIncluyeTexto: function (fieldId, texto) {
      var campo = this.camposPorId[fieldId];
      var raw = this.answers[fieldId] || '';
      if (!raw || !campo) return false;
      var icodes = raw.split(',');
      var opts = this.opcionesPara(campo);
      texto = texto.toLowerCase();
      for (var i = 0; i < icodes.length; i++) {
        for (var j = 0; j < opts.length; j++) {
          if (opts[j].icode === icodes[i] && this.etiqueta(campo, opts[j]).toLowerCase().indexOf(texto) >= 0) return true;
        }
      }
      return false;
    },

    // ── Valores / visibilidad ────────────────────────────────────────────
    getVal: function (id) { return this.answers[id] || ''; },
    // D6 (flow-E11, auto-sugerencias): si el tipo de documento (CI_5) es
    // inequívocamente colombiano o extranjero, predilligencia DP_5
    // (Nacionalidad) -- nunca pisa una respuesta que el agente ya haya
    // dado a mano. cd/oe/si quedan ambiguos a propósito (Carné Diplomático,
    // Otro, Persona sin identificación no determinan nacionalidad por sí
    // solos).
    sugerirNacionalidadPorDocumento: function () {
      if (this.answers['DP_5']) return;
      var codigosColombianos = ['cc', 'ti', 'rc', 'nc'];
      var codigosExtranjeros = ['ce', 'ps', 'pp', 'pf', 'pt', 'dn', 'vs', 'vr', 'sc', 'mv', 'no', 'pn'];
      var code = this.codigoSeleccionado('CI_5');
      var opts = this.opcionesPara(this.camposPorId['DP_5']);
      var destino = null;
      if (codigosColombianos.indexOf(code) >= 0) destino = 'co';
      else if (codigosExtranjeros.indexOf(code) >= 0) destino = 'ex';
      if (!destino) return;
      var opt = opts.filter(function (o) { return o.code === destino; })[0];
      if (opt) this.answers['DP_5'] = opt.icode;
    },
    setVal: function (id, v) {
      this.answers[id] = v;
      if (id === 'CI_5') this.sugerirNacionalidadPorDocumento();
      if (id === 'CI_5' || id === 'CI_6') this.verificarDuplicidadDocumento();
    },
    // D7 (flow-E11, "verificación de duplicidad si el campo es el número de
    // documento"): reusa el mismo /buscar que usa el Cabezote (D7 de
    // arquitectura: documento, nunca nombre) -- no hace falta un endpoint
    // nuevo. Debounced 500ms para no buscar en cada tecla.
    verificarDuplicidadDocumento: function () {
      var self = this;
      this.avisoDuplicidad = '';
      clearTimeout(this._duplicidadTimer);
      var code = this.codigoSeleccionado('CI_5');
      var numero = (this.answers['CI_6'] || '').trim();
      if (!code || !numero) return;
      this._duplicidadTimer = setTimeout(function () {
        fetch('/api/v1/victim-cases-v3/buscar?tipo_documento=' + encodeURIComponent(code) + '&documento=' + encodeURIComponent(numero))
          .then(function (r) { return r.json(); })
          .then(function (data) {
            var otros = (data.registros || []).filter(function (r) { return r.icode !== self.icode; });
            if (otros.length) {
              self.avisoDuplicidad = 'Ya existe' + (otros.length > 1 ? 'n ' + otros.length + ' Registros' : ' un Registro') + ' con este documento: ' + otros.map(function (r) { return r.names + ' ' + r.lastNames; }).join(', ') + '.';
            }
          })
          .catch(function () {});
      }, 500);
    },
    // D6: tratamiento según CI_1 (pronombre) -- "determina el tratamiento
    // mostrado en cabezote (Sra./Sr./Sre.)". Solo de visualización, no es
    // una respuesta que se guarde en otro campo.
    tratamiento: function () {
      var code = this.codigoSeleccionado('CI_1');
      return { la: 'Sra.', el: 'Sr.', le: 'Sre.' }[code] || '';
    },
    getBool: function (id) { return this.answers[id] === 'true'; },
    setBool: function (id, v) {
      this.answers[id] = v ? 'true' : 'false';
      // REG-08: DP_51 (tipo de discapacidad) solo aplica si DP_2 = Sí; al pasar a «No» se limpia para no
      // guardar una respuesta que quedó oculta (el backend la deja en NULL).
      if (id === 'DP_2' && !v) this.answers['DP_51'] = '';
    },

    estaMarcado: function (id, icode) {
      return (this.answers[id] || '').split(',').indexOf(icode) >= 0;
    },
    toggleMarcado: function (id, icode, max) {
      var actuales = (this.answers[id] || '').split(',').filter(function (v) { return v; });
      var idx = actuales.indexOf(icode);
      if (idx >= 0) {
        actuales.splice(idx, 1);
      } else {
        if (max && actuales.length >= max) return;
        // REG-16: un campo con `excluyente` (código de la opción «Ninguno») no admite esa opción junto
        // a otras: marcarla desmarca las demás y marcar otra la quita. Los demás campos no cambian.
        var campo = this.camposPorId[id];
        if (campo && campo.excluyente) {
          var excl = this.icodeForCode(campo.enumCategory, campo.excluyente);
          if (excl && icode === excl) actuales = [];
          else if (excl && actuales.indexOf(excl) >= 0) actuales.splice(actuales.indexOf(excl), 1);
        }
        actuales.push(icode);
      }
      this.answers[id] = actuales.join(',');
    },

    // Resuelve el `code` (no el icode) de la respuesta actual de un campo
    // dropdown de un solo valor -- D12 lo necesita para "Q17 o Q18 = Persona
    // en ASP" sin repetir la lógica de condicionalEnumCode por cada campo.
    codigoSeleccionado: function (fieldId) {
      var campo = this.camposPorId[fieldId];
      var icode = this.answers[fieldId];
      if (!campo || !icode) return '';
      var opt = this.opcionesPara(campo).filter(function (o) { return o.icode === icode; })[0];
      return opt ? opt.code : '';
    },
    pseudoCondicion: function (nombre) {
      switch (nombre) {
        case '__gestacionAplica': return true;
        case '__haySubtipo': return !!(this.answers['HV_9'] && this.answers['HV_9'].length);
        case '__ambitoLaboral': return this.multipleIncluyeTexto('HV_11', 'laboral');
        case '__agresorConocido':
          var et = this.etiquetaSeleccionada('AG_2').toLowerCase();
          return !!et && et.indexOf('desconocid') === -1;
        // D12: DP_20/21/22 (modalidad/fecha/motivo ASP) -- secciones.md
        // "condicional (Q17 o Q18 = Persona en ASP)". DP_17 occupation
        // code 'as', DP_18 income_generation_method code 'pr'.
        case '__personaEnASP':
          return this.codigoSeleccionado('DP_17') === 'as' || this.codigoSeleccionado('DP_18') === 'pr';
        // D13: CI_6 (número de documento) no aplica si CI_5="Persona sin
        // identificación" (code 'si').
        case '__tieneNumeroDocumento':
          return this.codigoSeleccionado('CI_5') !== 'si';
        default: return true;
      }
    },
    // D13: máscara de validación de CI_6 según el tipo de documento (CI_5)
    // -- secciones.md trae una tabla completa que antes no se aplicaba.
    ci6Mascara: function () {
      return MASCARA_DOC_NUMERO_V3[this.codigoSeleccionado('CI_5')];
    },
    ci6Invalido: function () {
      var m = this.ci6Mascara();
      var v = (this.answers['CI_6'] || '').trim();
      return !!m && !!v && !m.regex.test(v);
    },

    mostrarCampo: function (campo) {
      if (campo.condicional) {
        if (campo.condicional.campo.indexOf('__') === 0) {
          return this.pseudoCondicion(campo.condicional.campo);
        }
        return this.answers[campo.condicional.campo] === String(campo.condicional.valor);
      }
      if (campo.condicionalEnumCode) {
        var c = campo.condicionalEnumCode;
        var opts = this.opcionesPara(this.camposPorId[c.campo]);
        var sel = this.answers[c.campo];
        var opt = opts.filter(function (o) { return o.icode === sel; })[0];
        return !!opt && opt.code === c.code;
      }
      // D14: igual que condicionalEnumCode pero contra una lista de codes
      // válidos en vez de uno solo (ej. DP_19 -- "Q18=contrato formal O independiente").
      if (campo.condicionalEnumCodeIn) {
        var ci = campo.condicionalEnumCodeIn;
        var optsIn = this.opcionesPara(this.camposPorId[ci.campo]);
        var selIn = this.answers[ci.campo];
        var optIn = optsIn.filter(function (o) { return o.icode === selIn; })[0];
        return !!optIn && ci.codes.indexOf(optIn.code) >= 0;
      }
      if (campo.condicionalEnumLabel) {
        var et = this.etiquetaSeleccionada(campo.condicionalEnumLabel.campo);
        return et.indexOf(campo.condicionalEnumLabel.contains) >= 0;
      }
      if (campo.subtipoDe) {
        return !!(this.answers[campo.subtipoDe] && this.answers[campo.subtipoDe].length);
      }
      return true;
    },

    // REG-15: campos genéricos de GS filtrados por ids (CSV), en el orden pedido; el bloque genérico del
    // HTML se reutiliza en tres posiciones del Plan de acción (GS_1 · GS_3 · GS_20 y GS_4).
    camposGs: function (csv) {
      var self = this;
      return String(csv).split(',').map(function (id) { return self.camposPorId[id.trim()]; })
        .filter(function (c) { return c && self.mostrarCampo(c); });
    },
    camposVisibles: function (moduloId) {
      var self = this;
      var mod = this.modulos.filter(function (m) { return m.id === moduloId; })[0];
      if (!mod) return [];
      return mod.campos.filter(function (c) { return self.mostrarCampo(c); });
    },

    // A3 (flow-E12): "Validar formulario: Todas las preguntas [1er contacto]
    // visibles del módulo: required". Antes no existía ningún chequeo --
    // se podía guardar un módulo vacío. Solo mira las VISIBLES (camposVisibles
    // ya resuelve los condicionales) con req==='C1'; "false" cuenta como
    // respondido (es una respuesta real de un boolean), solo '' y undefined
    // cuentan como sin responder.
    obligatoriasFaltantesModulo: function (moduloId) {
      var self = this;
      return this.camposVisibles(moduloId).filter(function (c) {
        return c.req === 'C1' && !self.answers[c.id];
      });
    },
    // A4 (obs. #13 cliente): "solo letras" en nombres -- rechaza una entrada
    // compuesta ÚNICAMENTE por números. Antes `soloLetras: true` estaba
    // declarado en el catálogo (CI_2/CI_3/CA_1/AG_7) pero nada lo leía.
    esSoloNumeros: function (valor) {
      return /^[0-9]+$/.test((valor || '').trim());
    },
    campoSoloLetrasInvalido: function (campo) {
      return !!campo.soloLetras && this.esSoloNumeros(this.answers[campo.id]);
    },
    camposSoloLetrasInvalidosModulo: function (moduloId) {
      var self = this;
      return this.camposVisibles(moduloId).filter(function (c) { return self.campoSoloLetrasInvalido(c); });
    },
    // ── Guardar sección ──────────────────────────────────────────────────
    // enviarSeccion -- el POST real, sin ninguna validación. Lo usan tanto
    // guardarModulo (que valida antes de llamarlo) como guardarYSalir (D11,
    // flow-E14: "sin exigir que las obligatorias [1er contacto] estén
    // completas -- guardado parcial de cierre de llamada").
    enviarSeccion: function (moduloId) {
      var self = this;
      var payload = {};
      Object.keys(this.answers).forEach(function (k) {
        if (k.indexOf(moduloId + '_') === 0) payload[k] = self.answers[k];
      });
      return fetch('/api/v1/victim-cases-v3/' + this.icode + '/seccion/' + moduloId, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ answers: payload }),
      })
        .then(function (r) { return self.leerRespuesta(r); })
        .then(function (res) {
          // REG-11: antes no se miraba r.ok -- un 500 se leía como JSON sin triadaCompleta y la pantalla
          // igual decía "Sección guardada.". Ahora un rechazo se convierte en un Error con el texto
          // que la persona debe ver.
          if (!res.ok) {
            var err = new Error('guardado rechazado (' + res.status + ')');
            err.mensaje = self.mensajeErrorGuardado(res);
            throw err;
          }
          // REG-35: la tríada no se "des-completa" (las llaves respondidas se acumulan en el servidor): el cliente
          // nunca baja triadaCompleta de true a false con la respuesta de un guardado.
          var antes = self.triadaCompleta;
          self.triadaCompleta = self.triadaCompleta || !!res.data.triadaCompleta;
          // REG-38*: al completarse la tríada el Caso pasa a Iniciado y el borrador de entidades, barreras y
          // remisiones se materializa; se recarga para que la pantalla quede en solo lectura.
          if (!antes && self.triadaCompleta) self.recargarEntidadesBarreras();
          // REG-03 / REG-04b: el encabezado (nombre, documento, estado) se lee del servidor tras guardar.
          self.refrescarResumenCaso();
          return res.data;
        });
    },
    // REG-11: texto para la persona cuando el servidor rechazó un guardado. Un 400 trae el campo y el
    // motivo (formato de la hora, texto demasiado largo…): se le pide corregir ESE dato. Cualquier otra
    // falla (500, sesión vencida) se avisa tal cual, para que nadie crea que se guardó.
    mensajeErrorGuardado: function (res) {
      var d = res.data || {};
      var mayuscula = function (t) { t = String(t || '').trim(); return t ? t.charAt(0).toUpperCase() + t.slice(1) : ''; };
      var punto = function (t) { return /[.!?]$/.test(t) ? t : t + '.'; };
      if (res.status === 400 && d.detalle) {
        var campo = d.campo ? this.camposPorId[d.campo] : null;
        var cual = campo ? '«' + campo.label + '»' : (d.campo ? d.campo : '');
        return 'No se guardó la sección. ' + (cual ? 'Corrige ' + cual + ': ' : 'Corrige lo que escribiste: ') + punto(mayuscula(d.detalle)) + ' Luego vuelve a guardar.';
      }
      if (res.status === 401 || res.status === 403) {
        // REG-21: se dice qué hacer y que lo escrito NO se pierde (la pantalla no se recarga ni redirige).
        return 'No se guardó la sección: tu sesión venció o no tienes permiso. Inicia sesión de nuevo en otra pestaña y vuelve a pulsar «Guardar sección»: lo que escribiste sigue en esta pantalla.';
      }
      return 'No se pudo guardar la sección por un error del servidor. Lo que escribiste sigue en pantalla: intenta de nuevo y, si el problema continúa, avisa a soporte.';
    },
    // Muestra la falla de un guardado: aviso en rojo (arriba, con scroll hasta él) y marca junto al botón.
    mostrarErrorGuardado: function (e) {
      var self = this;
      this.errorGuardado = (e && e.mensaje) || 'No se pudo guardar la sección: no hay conexión con el servidor. Lo que escribiste sigue en pantalla; intenta de nuevo.';
      this.mensajeGuardado = 'No se guardó.';
      this.mensajeGuardadoEsError = true; // REG-05
      this.$nextTick(function () {
        var el = self.$el && self.$el.querySelector && self.$el.querySelector('#registro-v3-error-guardado');
        if (el && el.scrollIntoView) el.scrollIntoView({ behavior: 'smooth', block: 'center' });
      });
    },
    // Devuelve una promesa que resuelve true si se guardó y false si no (validación o error): así
    // finalizarRegistro (REG-37) sabe si puede seguir. `sinAvanzar` evita pasar al siguiente módulo (REG-01).
    guardarModulo: function (moduloId, sinAvanzar) {
      var self = this;
      this.mensajeAvance = '';
      var faltantes = this.obligatoriasFaltantesModulo(moduloId);
      if (faltantes.length) {
        this.mensajeGuardado = 'Faltan obligatorias de 1er contacto: ' + faltantes.map(function (c) { return c.label; }).join('; ');
        this.mensajeGuardadoEsError = true; // REG-05: el recordatorio sale en rojo
        return Promise.resolve(false);
      }
      var invalidos = this.camposSoloLetrasInvalidosModulo(moduloId);
      if (invalidos.length) {
        this.mensajeGuardado = 'Estos campos no pueden ser solo números: ' + invalidos.map(function (c) { return c.label; }).join('; ');
        this.mensajeGuardadoEsError = true;
        return Promise.resolve(false);
      }
      if (moduloId === 'CI' && this.ci6Invalido()) {
        this.mensajeGuardado = 'Número de documento inválido para el tipo seleccionado: ' + this.ci6Mascara().hint;
        this.mensajeGuardadoEsError = true;
        return Promise.resolve(false);
      }
      this.guardando = true;
      this.mensajeGuardado = '';
      this.mensajeGuardadoEsError = false;
      this.errorGuardado = '';
      return this.enviarSeccion(moduloId)
        .then(function () {
          self.mensajeGuardado = 'Sección guardada.';
          self.mensajeGuardadoEsError = false;
          if (!sinAvanzar) self.avanzarTrasGuardar(moduloId); // REG-01
          return true;
        })
        .catch(function (e) { self.mostrarErrorGuardado(e); return false; })
        .finally(function () { self.guardando = false; });
    },
    // REG-01: tras guardar bien, pasa al siguiente módulo NO bloqueado (salvo en el último). Si la persona ya
    // se movió a otro módulo mientras llegaba la respuesta, no se la mueve.
    avanzarTrasGuardar: function (moduloId) {
      if (this.moduloActivo !== moduloId) return;
      var destino = this.moduloSiguiente(moduloId);
      if (!destino) return;
      var etiqueta = function (modulos, id) { var m = modulos.filter(function (x) { return x.id === id; })[0]; return m ? m.label : id; };
      var de = etiqueta(this.modulos, moduloId), a = etiqueta(this.modulos, destino);
      this.irAModulo(destino);
      // Conserva el texto «Sección guardada.» (lo buscan las pruebas de REG-11) y dice de dónde a dónde se pasó.
      this.mensajeAvance = '✓ Sección guardada. Pasaste de «' + de + '» a «' + a + '».';
    },
    // REG-17: «Guardar sección» de la Ruta de atención integral. GS_14 y GS_17 se guardan con la sección GS
    // (misma llave), así que se reutiliza enviarSeccion('GS'); las entidades y barreras ya se guardaron al agregarlas.
    // No avanza al siguiente módulo (la barra «Siguiente» sí funciona).
    guardarRutaAtencion: function () {
      var self = this;
      var faltantes = this.obligatoriasFaltantesModulo('RA');
      if (faltantes.length) {
        this.mensajeGuardado = 'Faltan obligatorias de 1er contacto: ' + faltantes.map(function (c) { return c.label; }).join('; ');
        this.mensajeGuardadoEsError = true;
        return Promise.resolve(false);
      }
      this.guardando = true;
      this.mensajeGuardado = '';
      this.mensajeGuardadoEsError = false;
      this.errorGuardado = '';
      return this.enviarSeccion('GS')
        .then(function () {
          self.mensajeGuardado = 'Sección guardada.';
          self.mensajeGuardadoEsError = false;
          return true;
        })
        .catch(function (e) { self.mostrarErrorGuardado(e); return false; })
        .finally(function () { self.guardando = false; });
    },
    // D11 (flow-E14, botón "Finalizar llamada (guardar y salir)"): guarda el
    // módulo actual SIN exigir las obligatorias -- es un cierre de llamada
    // a medias, no un guardado normal -- y vuelve al Cabezote. Antes no
    // existía ningún botón ni respaldo para este caso; el único guardado
    // disponible era el de guardarModulo, que SÍ exige obligatorias y
    // bloqueaba salir de una llamada incompleta.
    guardarYSalir: function () {
      var self = this;
      if (['CI', 'DP', 'CA', 'HV', 'AG', 'GS'].indexOf(this.moduloActivo) === -1) {
        this.volverACabezote();
        return;
      }
      this.guardando = true;
      this.errorGuardado = '';
      this.enviarSeccion(this.moduloActivo)
        .then(function () {
          self.guardando = false;
          self.volverACabezote();
        })
        .catch(function (e) {
          // REG-11: antes el error se ignoraba y se salía igual -> la persona perdía lo escrito sin enterarse.
          // Ahora se le dice qué falló y se le deja elegir: corregir y volver a intentar, o salir sin guardar.
          self.guardando = false;
          self.mostrarErrorGuardado(e);
          if (window.confirm(self.errorGuardado + '\n\n¿Salir de todos modos? Lo que escribiste en este módulo se perdería.')) {
            self.volverACabezote();
          }
        });
    },
    volverACabezote: function () {
      this.errorGuardado = '';
      this.etapa = 'cabezote';
      this.icode = null;
      this.answers = {};
      this.entidades = []; this.barreras = [];
      this.entityMeta = {}; this.barrierMeta = {}; this.entitiesReadOnly = false;
    },

    // A1 (flow-E10): TZ/GS deben quedar bloqueados hasta que la tríada
    // (CI/DP/CA/HV/AG) esté completa -- antes no había ningún guard acá.
    // REG-13: solo el Tamizaje queda bloqueado hasta completar la tríada. Gestión Salvia se puede diligenciar
    // siempre (y sigue siendo obligatoria para finalizar: gsObligatoriasFinalizarV3 en el servidor).
    moduloBloqueado: function (id) {
      return id === 'TZ' && !this.triadaCompleta;
    },
    irAModulo: function (id) {
      if (this.moduloBloqueado(id)) {
        this.mensajeModuloBloqueado = 'Completa primero los módulos de Contacto, Datos personales, Contacto de apoyo, Hechos y Agresor para habilitar Tamizaje.';
        return;
      }
      this.mensajeModuloBloqueado = '';
      this.moduloActivo = id;
      this.mensajeGuardado = '';
      this.mensajeGuardadoEsError = false;
      this.mensajeAvance = '';
      this.errorGuardado = '';
    },
    // REG-01: primer módulo NO bloqueado después (o antes) de `id`, según el orden de modulos[]; null si no hay.
    // Se recorre modulos[] (no hay módulos escritos a mano: el de WP-4 entra solo).
    moduloSiguiente: function (id) {
      var ids = this.modulos.map(function (m) { return m.id; });
      var idx = ids.indexOf(id);
      if (idx < 0) return null;
      for (var i = idx + 1; i < ids.length; i++) { if (!this.moduloBloqueado(ids[i])) return ids[i]; }
      return null;
    },
    moduloAnterior: function (id) {
      var ids = this.modulos.map(function (m) { return m.id; });
      var idx = ids.indexOf(id);
      if (idx < 0) return null;
      for (var i = idx - 1; i >= 0; i--) { if (!this.moduloBloqueado(ids[i])) return ids[i]; }
      return null;
    },
    // Flechas «Siguiente» / «Anterior» (E22 / E23): solo navegan, no guardan.
    irASiguiente: function () { var d = this.moduloSiguiente(this.moduloActivo); if (d) this.irAModulo(d); },
    irAAnterior: function () { var d = this.moduloAnterior(this.moduloActivo); if (d) this.irAModulo(d); },
    posicionModulo: function () {
      var self = this;
      return this.modulos.map(function (m) { return m.id; }).indexOf(self.moduloActivo) + 1;
    },

    // REG-02: porcentaje de avance = obligatorias de 1er contacto VISIBLES respondidas / visibles (la misma regla
    // del recordatorio «Faltan obligatorias…»). null si el módulo está bloqueado o no tiene de esas (se muestra
    // «—» y no pesa). Se calcula en vivo sobre `answers` (decisión aprobada: «Vamos con En vivo»).
    porcentajeModulo: function (id) {
      if (this.moduloBloqueado(id)) return null;
      var total = this.camposVisibles(id).filter(function (c) { return c.req === 'C1'; }).length;
      if (!total) return null;
      var respondidas = total - this.obligatoriasFaltantesModulo(id).length;
      return { respondidas: respondidas, total: total, pct: Math.round(100 * respondidas / total) };
    },
    porcentajeGlobal: function () {
      var self = this, respondidas = 0, total = 0;
      this.modulos.forEach(function (m) {
        var r = self.porcentajeModulo(m.id);
        if (r) { respondidas += r.respondidas; total += r.total; }
      });
      return total ? Math.round(100 * respondidas / total) : null;
    },
    etiquetaPorcentaje: function (id) {
      var r = this.porcentajeModulo(id);
      return r ? r.pct + ' %' : '—';
    },

    // REG-04: nombre legible del estado del Caso (nunca el código crudo; un código nuevo sale «Desconocido»).
    etiquetaEstado: function (code) {
      if (!code) return '';
      return ESTADOS_CASO_V3[code] || 'Desconocido';
    },
    // REG-03: nombre y documento del encabezado. Un Borrador nuevo nace con valores provisionales
    // («Pendiente Pendiente», documento 0000000000): se muestran como «Persona sin identificar» y
    // «Documento sin registrar» hasta que se guarde CI.
    nombreEncabezado: function () {
      var r = this.resumenCaso || {};
      var n = ((r.names || '') + ' ' + (r.lastNames || '')).replace(/\s+/g, ' ').trim();
      if (!n || n === 'Pendiente Pendiente') return 'Persona sin identificar';
      var t = this.tratamiento();
      var ident = String(this.answers['CI_3'] || '').trim();
      return (t ? t + ' ' : '') + n + (ident && ident !== n ? ' (' + ident + ')' : '');
    },
    documentoEncabezado: function () {
      var r = this.resumenCaso || {};
      if (!r.docNumber || r.docNumber === '0000000000') return 'Documento sin registrar';
      return (r.docType ? r.docType + '-' : '') + r.docNumber;
    },
    // Tras un guardado exitoso se vuelve a leer el resumen del Caso (nombre, documento y estado). Si falla no se
    // avisa: el guardado ya se confirmó. Tampoco baja triadaCompleta (nunca pasa de true a false).
    refrescarResumenCaso: function () {
      var self = this;
      var icode = this.icode;
      if (!icode) return Promise.resolve();
      return fetch('/api/v1/victim-cases-v3/' + icode)
        .then(function (r) { return self.leerRespuesta(r); })
        .then(function (res) {
          if (!res.ok || !res.data.summary || self.icode !== icode) return;
          self.resumenCaso = res.data.summary;
          self.triadaCompleta = self.triadaCompleta || !!res.data.triadaCompleta;
        })
        .catch(function () { /* sin aviso */ });
    },

    // ── Tamizaje ──────────────────────────────────────────────────────────
    tzCampoVisible: function (campo) {
      return campo.inicial || this.tzRevelarResto;
    },
    tzToggle: function (id, val) {
      this.answers[id] = val ? 'true' : 'false';
      var self = this;
      var algunaSi = TAMIZAJE_GENERALES_V3.some(function (c) { return c.inicial && self.answers[c.id] === 'true'; });
      this.tzRevelarResto = algunaSi;
    },
    calcularEsOsigd: function () {
      var orientacion = this.etiquetaSeleccionada('DP_9').toLowerCase();
      var identidad = this.etiquetaSeleccionada('DP_8').toLowerCase();
      return (orientacion && orientacion.indexOf('heterosexual') === -1) ||
        identidad.indexOf('trans') >= 0 || identidad.indexOf('no binaria') >= 0;
    },
    calcularWasPartner: function () {
      var et = this.etiquetaSeleccionada('AG_3').toLowerCase();
      return et.indexOf('pareja') >= 0;
    },
    guardarTamizaje: function () {
      var self = this;
      var wasPartner = this.calcularWasPartner();
      var esOsigd = this.calcularEsOsigd();
      var subtipoFem = this.multipleIncluyeTexto('HV_10', 'tentativa feminicidio');
      var subtipoHom = this.multipleIncluyeTexto('HV_10', 'tentativa homicidio');

      // REG-34: por etapas. Con las cinco iniciales en «No» (resto sin revelar) solo se envían esas cinco: el
      // servidor deja el resto en «No» sin marcarlo como respondido. Con el resto revelado se envía todo.
      var tzAnswers = {};
      TAMIZAJE_GENERALES_V3.forEach(function (c) {
        if (self.tzRevelarResto || c.inicial) tzAnswers[c.id] = self.answers[c.id] === 'true' ? 'true' : 'false';
      });
      if (this.tzRevelarResto) {
        tzAnswers['TZ_8'] = this.answers['TZ_8'] === 'true' ? 'true' : 'false';
        var lista = wasPartner ? TAMIZAJE_PAREJA_V3 : TAMIZAJE_NO_PAREJA_V3;
        var prefijo = wasPartner ? 'TZP_' : 'TZNP_';
        lista.forEach(function (_, idx) {
          var id = prefijo + (idx + 1);
          tzAnswers[id] = self.answers[id] === 'true' ? 'true' : 'false';
        });
      }

      this.tzGuardando = true;
      this.errorGuardado = '';
      fetch('/api/v1/victim-cases-v3/' + this.icode + '/tamizaje', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          answers: tzAnswers, wasPartner: wasPartner, esOsigd: esOsigd,
          subtipoFeminicidio: subtipoFem, subtipoHomicidioAbyecto: subtipoHom,
        }),
      })
        .then(function (r) { return self.leerRespuesta(r); })
        .then(function (res) {
          // REG-12: antes no se revisaba r.ok -- un 500 dejaba {error:…} como «resultado» del Tamizaje y la
          // persona no sabía que no se había guardado.
          if (!res.ok) {
            var err = new Error('tamizaje rechazado (' + res.status + ')');
            err.mensaje = (res.status === 401 || res.status === 403)
              ? 'No se guardó el Tamizaje: tu sesión venció o no tienes permiso. Vuelve a iniciar sesión; lo que respondiste en esta pantalla no se guardó.'
              : 'No se pudo guardar el Tamizaje por un error del servidor. Lo que respondiste sigue en pantalla: intenta de nuevo y, si el problema continúa, avisa a soporte.';
            throw err;
          }
          self.tzResultado = res.data;
        })
        .catch(function (e) {
          if (!e.mensaje) e.mensaje = 'No se pudo guardar el Tamizaje: no hay conexión con el servidor. Lo que respondiste sigue en pantalla; intenta de nuevo.';
          self.mostrarErrorGuardado(e);
        })
        .finally(function () { self.tzGuardando = false; });
    },
    // D10 (flow-E11, "alertas de riesgo/ideación suicida"): antes
    // tzResultado solo se mostraba como informativo, sin resaltar los casos
    // que de verdad requieren acción inmediata.
    riesgoEsAlto: function () {
      return !!this.tzResultado && (this.tzResultado.riskLevel === 'alto' || this.tzResultado.riskLevel === 'extremo');
    },
    hayIdeacionSuicida: function () {
      return this.answers['TZP_24'] === 'true' || this.answers['TZNP_24'] === 'true';
    },
    tzWasPartner: function () { return this.calcularWasPartner(); },
    tzListaPareja: function () { return TAMIZAJE_PAREJA_V3; },
    tzListaNoPareja: function () { return TAMIZAJE_NO_PAREJA_V3; },

    // ── GS § B remisión ──────────────────────────────────────────────────
    equipoMeta: function (code) {
      return EQUIPOS_REMISION_V3.filter(function (e) { return e.value === code; })[0];
    },
    icodeForCode: function (enumCategory, code) {
      var opt = (this.enumCatalog[enumCategory] || []).filter(function (o) { return o.code === code; })[0];
      return opt ? opt.icode : null;
    },
    criterioPuntaje: function (categoria, codigosSeleccionados) {
      var meta = CRITERIOS_PUNTAJE_V3[categoria];
      if (!meta) return { puntos: 0, minimo: 0, cumple: codigosSeleccionados.length > 0 };
      var puntos = 0;
      codigosSeleccionados.forEach(function (c) { puntos += meta.puntos[c] || 0; });
      return { puntos: puntos, minimo: meta.minimo, cumple: puntos >= meta.minimo };
    },
    cumpleCriterioEquipo: function (equipo) {
      if (!equipo.criterioFieldId) return true;
      var campo = this.camposPorId[equipo.criterioFieldId];
      var opts = this.opcionesPara(campo);
      var seleccionados = (this.answers[equipo.criterioFieldId] || '').split(',').filter(function (v) { return v; });
      var codigos = seleccionados.map(function (icode) {
        var o = opts.filter(function (o) { return o.icode === icode; })[0];
        return o ? o.code : null;
      }).filter(function (c) { return c; });
      if (equipo.value === 'ap' || equipo.value === 'ah') {
        if (codigos.indexOf('ob') === -1) return false;
      }
      if (equipo.criterioCategory && CRITERIOS_PUNTAJE_V3[equipo.criterioCategory]) {
        var puntuables = codigos.filter(function (c) { return c !== 'ob'; });
        return this.criterioPuntaje(equipo.criterioCategory, puntuables).cumple;
      }
      return codigos.length > 0;
    },

    // ── GS § C/D -- buscador geográfico compartido ──────────────────────────
    // Reconstruido sobre entity_branch real (mismos 3 endpoints que usa
    // case-entities.js) -- ver registro-v3-pendientes.md § B.
    cargarGeoDepartamentos: function () {
      if (this.geoDepartamentos.length) return;
      var self = this;
      fetch('/api/v1/locations/departments')
        .then(function (r) { return r.json(); })
        .then(function (data) { self.geoDepartamentos = Array.isArray(data) ? data : []; })
        .catch(function () { self.geoDepartamentos = []; });
    },
    recargarEntidadesBarreras: function () {
      var self = this;
      fetch('/api/v1/victim-cases-v3/' + this.icode)
        .then(function (r) { return r.json(); })
        .then(function (data) {
          self.entidades = data.entities || [];
          self.barreras = data.barriers || [];
          self.entityMeta = data.entityMeta || {};
          self.barrierMeta = data.barrierMeta || {};
          self.entitiesReadOnly = !!data.entitiesReadOnly;
          self.cargarEtiquetasAcciones();
        });
    },
    // REG-15: la tarjeta lista las acciones institucionales de la entidad con su estado. El borrador
    // (y el snapshot) guardan solo los ids; las etiquetas salen del catálogo de la sede, una sola vez.
    cargarEtiquetasAcciones: function () {
      var self = this;
      this.entidades.forEach(function (e) {
        var m = self.metaDeEntidad(e);
        var branch = e.entityBranchId;
        if (!m || ((m.pending || []).length + (m.completed || []).length) === 0 || self.obligacionesPorBranch[branch]) return;
        self.obligacionesPorBranch[branch] = {};
        fetch('/api/v1/entity-branches/' + encodeURIComponent(branch) + '/obligations')
          .then(function (r) { return r.ok ? r.json() : []; })
          .then(function (data) {
            var mapa = {};
            (Array.isArray(data) ? data : []).forEach(function (o) { mapa[o.value] = o.label; });
            self.obligacionesPorBranch[branch] = mapa;
          })
          .catch(function () { /* sin etiquetas: la tarjeta no lista las acciones */ });
      });
    },
    accionesDeEntidad: function (e) {
      var m = this.metaDeEntidad(e);
      var mapa = this.obligacionesPorBranch[e.entityBranchId] || {};
      if (!m) return [];
      var out = [];
      (m.completed || []).forEach(function (id) { if (mapa[id]) out.push({ id: id, estado: 'completada', label: mapa[id] }); });
      (m.pending || []).forEach(function (id) { if (mapa[id]) out.push({ id: id, estado: 'pendiente', label: mapa[id] }); });
      return out;
    },
    ubicacionEntidad: function (e) {
      return [e.departmentName, e.cityName, e.townName].filter(function (v) { return v; }).join(' · ');
    },
    // entidades[].sector viene de entity.entity_sector (código interno de 2
    // letras: he/js/pt/os -- ver entity_case_repository.go), mientras que
    // barreras[].sector es el vocabulario real por palabra que manda este
    // mismo formulario (salud/justicia/proteccion/otras_instituciones -- ver
    // VictimCaseBarrierInput.SectorCode). Dos catálogos distintos para el
    // mismo concepto -- este mapa cubre ambos para la tabla "ya agregadas".
    sectorLabelEntidadBarrera: function (code) {
      var labels = {
        he: 'Salud', js: 'Justicia', pt: 'Protección', os: 'Otras instituciones',
        salud: 'Salud', justicia: 'Justicia', proteccion: 'Protección', otras_instituciones: 'Otras instituciones',
      };
      return labels[code] || code || '—';
    },

    // ── GS § C Entidades ──────────────────────────────────────────────────
    onCambiaDeptoEntidad: function () {
      this.nuevaEntidad.ciudad = ''; this.nuevaEntidad.municipio = ''; this.nuevaEntidad.EntityBranchID = '';
      this.entidadCiudades = []; this.entidadMunicipios = []; this.entidadesDisponibles = []; this.obligacionesEntidad = [];
      this.reiniciarBuscadorEntidad();
      if (!this.nuevaEntidad.depto) return;
      var self = this;
      fetch('/api/v1/locations/cities?department_id=' + encodeURIComponent(this.nuevaEntidad.depto))
        .then(function (r) { return r.json(); })
        .then(function (data) { self.entidadCiudades = Array.isArray(data) ? data : []; })
        .catch(function () { self.entidadCiudades = []; });
    },
    onCambiaCiudadEntidad: function () {
      this.nuevaEntidad.municipio = ''; this.nuevaEntidad.EntityBranchID = '';
      this.entidadMunicipios = []; this.entidadesDisponibles = []; this.obligacionesEntidad = [];
      this.reiniciarBuscadorEntidad();
      if (!this.nuevaEntidad.ciudad) return;
      var self = this;
      fetch('/api/v1/locations/towns?city_id=' + encodeURIComponent(this.nuevaEntidad.ciudad))
        .then(function (r) { return r.json(); })
        .then(function (data) { self.entidadMunicipios = Array.isArray(data) ? data : []; })
        .catch(function () { self.entidadMunicipios = []; });
    },
    onCambiaMunicipioEntidad: function () { this.buscarEntidadesEntidad(); },
    onCambiaSectorEntidad: function () { this.buscarEntidadesEntidad(); },
    // REG-19: la lista de sedes del municipio y sector se pide UNA sola vez (al elegir el último de
    // los dos). Cada petición lleva un número: si el usuario cambia el municipio o el sector mientras
    // llega la respuesta anterior, esa respuesta se descarta.
    buscarEntidadesEntidad: function () {
      this.nuevaEntidad.EntityBranchID = ''; this.entidadesDisponibles = [];
      this.obligacionesEntidad = []; this.nuevaEntidad.CompletedActions = []; this.nuevaEntidad.PendingActions = [];
      this.reiniciarBuscadorEntidad();
      if (!this.nuevaEntidad.municipio || !this.nuevaEntidad.sector) return;
      var self = this;
      var token = this._reqSedes;
      self.entidadBuscando = true;
      fetch('/api/v1/entity-branches?town_code=' + encodeURIComponent(this.nuevaEntidad.municipio) + '&sector=' + encodeURIComponent(this.nuevaEntidad.sector))
        .then(function (r) { return r.json(); })
        .then(function (data) {
          if (token !== self._reqSedes) return;
          self.entidadesDisponibles = Array.isArray(data) ? data : []; self.entidadBuscando = false;
        })
        .catch(function () {
          if (token !== self._reqSedes) return;
          self.entidadesDisponibles = []; self.entidadBuscando = false;
        });
    },
    sinTildes: function (t) {
      return String(t || '').normalize('NFD').replace(/[\u0300-\u036f]/g, '').toLowerCase();
    },
    // Limpia el texto del buscador, cancela el debounce pendiente e invalida las respuestas en vuelo.
    reiniciarBuscadorEntidad: function () {
      clearTimeout(this._entidadDebounce);
      this.entidadTexto = ''; this.textoEntidad = ''; this.entidadEnfocada = false;
      this._reqSedes = (this._reqSedes || 0) + 1;
      this._reqObligaciones = (this._reqObligaciones || 0) + 1;
      this.entidadBuscando = false; this.obligacionesCargando = false;
    },
    // @input del campo «Buscar entidad»: el texto visible cambia al instante y el filtro se aplica
    // tras 250 ms sin teclear. Escribir después de elegir una sugerencia limpia la selección.
    onInputBuscarEntidad: function (valor) {
      var self = this;
      this.entidadTexto = valor;
      if (this.nuevaEntidad.EntityBranchID) this.limpiarSeleccionEntidad();
      clearTimeout(this._entidadDebounce);
      this._entidadDebounce = setTimeout(function () { self.textoEntidad = self.entidadTexto; }, 250);
    },
    limpiarSeleccionEntidad: function () {
      this.nuevaEntidad.EntityBranchID = '';
      this.obligacionesEntidad = []; this.nuevaEntidad.CompletedActions = []; this.nuevaEntidad.PendingActions = [];
      this._reqObligaciones = (this._reqObligaciones || 0) + 1;
      this.obligacionesCargando = false;
    },
    seleccionarEntidadSugerida: function (opt) {
      clearTimeout(this._entidadDebounce);
      this.nuevaEntidad.EntityBranchID = opt.id;
      this.entidadTexto = opt.name; this.textoEntidad = opt.name;
      this.onCambiaEntidadSeleccionada();
    },
    // Al elegir la entidad se piden sus acciones institucionales (única otra petición del buscador).
    onCambiaEntidadSeleccionada: function () {
      this.obligacionesEntidad = []; this.nuevaEntidad.CompletedActions = []; this.nuevaEntidad.PendingActions = [];
      if (!this.nuevaEntidad.EntityBranchID) return;
      var self = this;
      var token = this._reqObligaciones = (this._reqObligaciones || 0) + 1;
      self.obligacionesCargando = true;
      fetch('/api/v1/entity-branches/' + encodeURIComponent(this.nuevaEntidad.EntityBranchID) + '/obligations')
        .then(function (r) { return r.ok ? r.json() : []; })
        .then(function (data) {
          if (token !== self._reqObligaciones) return;
          // "Otra" no tiene captura de texto libre en V3 todavía -- se excluye
          // del checklist, igual que case-entities.js.
          self.obligacionesEntidad = (Array.isArray(data) ? data : []).filter(function (o) { return o.value !== 'otra'; });
          self.obligacionesCargando = false;
        })
        .catch(function () {
          if (token !== self._reqObligaciones) return;
          self.obligacionesEntidad = []; self.obligacionesCargando = false;
        });
    },
    // Una acción no puede estar a la vez completada y pendiente: marcarla en una lista la quita de la otra.
    toggleAccion: function (propia, otra, valor) {
      var arr = this.nuevaEntidad[propia];
      var idx = arr.indexOf(valor);
      if (idx >= 0) { arr.splice(idx, 1); return; }
      arr.push(valor);
      var j = this.nuevaEntidad[otra].indexOf(valor);
      if (j >= 0) this.nuevaEntidad[otra].splice(j, 1);
    },
    // Acciones institucionales como chips (mismo estilo que el modal de Gestión institucional).
    estiloChip: function (seleccionado) {
      return seleccionado
        ? 'border:1px solid #5106A7;background:#5106A7;color:#fff;border-radius:20px;padding:6px 14px;font-size:.8rem;font-weight:600;cursor:pointer'
        : 'border:1px solid #c4b5fd;background:#fff;color:#5106A7;border-radius:20px;padding:6px 14px;font-size:.8rem;font-weight:500;cursor:pointer';
    },
    agregarEntidad: function () {
      var self = this;
      if (!this.nuevaEntidad.EntityBranchID) return;
      var branchId = Number(this.nuevaEntidad.EntityBranchID);
      var faltantes = this.entidadFaltantes();
      if (faltantes.length) { self.errorEntidad = 'Faltan campos obligatorios: ' + faltantes.join('; '); return; }
      // REG-20: Agreement y DueDiligence son copias de GS_15 y GS_16 en este momento (esas
      // respuestas pueden no estar guardadas todavía); los canales solo operan si ambas son «Sí».
      var activa = this.activacionAplica();
      var acudida = this.nuevaEntidad.Kind === 'acudida';
      var agreement = acudida ? (this.nuevaEntidad.PendingActions.length > 0 && this.nuevaEntidad.Agreement === true) : this.answers['GS_15'] === 'true';
      var dueDiligence = acudida ? (agreement && this.nuevaEntidad.DueDiligence === true) : this.answers['GS_16'] === 'true';
      var body = {
        Kind: this.nuevaEntidad.Kind,
        EntityBranchID: branchId,
        FileNumber: this.nuevaEntidad.FileNumber,
        Solicitud: this.nuevaEntidad.Solicitud,
        ReportDate: this.nuevaEntidad.ReportDate,
        ActivationChannels: activa ? this.nuevaEntidad.ActivationChannels.slice() : [],
        Agreement: agreement,
        DueDiligence: dueDiligence,
        ManagementNotes: activa ? this.nuevaEntidad.ManagementNotes : '',
        CompletedActions: this.nuevaEntidad.Kind === 'acudida' ? this.nuevaEntidad.CompletedActions : [],
        PendingActions: this.nuevaEntidad.PendingActions,
        OtherInfo: this.nuevaEntidad.OtherInfo,
      };
      self.errorEntidad = '';
      fetch('/api/v1/victim-cases-v3/' + this.icode + '/entidades', {
        method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(body),
      })
        .then(function (r) { return r.json().then(function (data) { return { ok: r.ok, data: data }; }); })
        .then(function (res) {
          if (!res.ok) { self.errorEntidad = res.data.error || 'No se pudo agregar la entidad'; self.marcarSoloLectura(res); return; }
          self.nuevaEntidad = self.entidadVacia(self.nuevaEntidad.Kind);
          self.entidadCiudades = []; self.entidadMunicipios = []; self.entidadesDisponibles = []; self.obligacionesEntidad = [];
          self.reiniciarBuscadorEntidad();
          self.recargarEntidadesBarreras();
        });
    },
    // ── REG-04 (c) · «Gestión» de la barrera con etiqueta, no con código ───────
    // Recibe el CSV de barrier_v2.management_actions: traduce cada valor (los 2 vigentes de
    // Seguimiento v4 y los 6 históricos), deja tal cual el que no conoce y muestra «—» si no hay.
    etiquetaGestionBarrera: function (csv) {
      var vigentes = {};
      BARRERA_V4_V3.filter(function (q) { return q.key === 'gestion'; })[0].opciones.forEach(function (o) { vigentes[o.value] = o.label; });
      var valores = String(csv || '').split(',').map(function (v) { return v.trim(); }).filter(function (v) { return v; });
      if (!valores.length) return '—';
      return valores.map(function (v) { return vigentes[v] || GESTION_BARRERA_LEGADO_V3[v] || v; }).join(', ');
    },

    // ── REG-20 · Canal de activación como Seguimiento v4 ─────────────────────
    // Los canales se piden cuando la persona está de acuerdo con la activación de la Ruta
    // (GS_15) y Salvia debe activar rutas en el marco de la debida diligencia (GS_16): igual
    // que v4BothYes en Seguimiento. Sin ambas respuestas en «Sí» no se crea ninguna tarea.
    // En el Plan de acción las dos respuestas son GS_15 y GS_16; en la Ruta de atención integral la entidad
    // acudida con acciones pendientes las contesta ella misma (D17, como el repetidor de Seguimiento v4).
    activacionAplica: function () {
      var f = this.nuevaEntidad;
      if (f.Kind === 'acudida') return f.PendingActions.length > 0 && f.Agreement === true && f.DueDiligence === true;
      return this.answers['GS_15'] === 'true' && this.answers['GS_16'] === 'true';
    },
    // «Ninguna» es excluyente con los demás canales (no crea tarea en Seguimiento).
    toggleCanal: function (valor) {
      var arr = this.nuevaEntidad.ActivationChannels;
      var idx = arr.indexOf(valor);
      if (idx >= 0) { arr.splice(idx, 1); return; }
      if (valor === 'ninguna') arr.splice(0, arr.length);
      else if (arr.indexOf('ninguna') >= 0) arr.splice(arr.indexOf('ninguna'), 1);
      arr.push(valor);
    },
    // Etiquetas de lo obligatorio que falta para agregar la entidad (vacío = se puede agregar).
    entidadFaltantes: function () {
      var e = this.nuevaEntidad;
      var f = [];
      var vacio = function (v) { return !String(v || '').trim(); };
      if (e.Kind === 'activacion') {
        if (vacio(e.OtherInfo)) f.push('Objetivo con la entidad');
      } else {
        // REG-17: obligatorias del repetidor «Entidades a las que ya acudió» de Seguimiento v4. Las acciones
        // se piden cuando la sede tiene acciones en el catálogo: basta con marcar una (completada o pendiente).
        if (vacio(e.FileNumber)) f.push(ENTIDAD_V4_V3.radicado.label);
        if (vacio(e.ReportDate)) f.push(ENTIDAD_V4_V3.fecha.label);
        if (this.obligacionesEntidad.length && !(e.CompletedActions.length + e.PendingActions.length)) f.push('Acciones institucionales completadas o pendientes');
        if (vacio(e.OtherInfo)) f.push(ENTIDAD_V4_V3.otraRuta.label);
        if (e.PendingActions.length) {
          if (e.Agreement === null) f.push(ENTIDAD_V4_V3.acuerdo.label);
          else if (e.Agreement === true && e.DueDiligence === null) f.push(ENTIDAD_V4_V3.debida.label);
        }
      }
      if (this.activacionAplica()) {
        if (!e.ActivationChannels.length) f.push(ENTIDAD_V4_V3.canales.label);
        if (vacio(e.ManagementNotes)) f.push(ENTIDAD_V4_V3.gestiones.label);
      }
      return f;
    },
    metaDeEntidad: function (e) { return this.entityMeta[String(e.entityBranchId)] || null; },
    // Tipo de la entidad: sale de entityMeta (REG-38*, D1); los datos anteriores se reconocen
    // porque guardaban «Canal de activación: …» en solicitud.
    kindDeEntidad: function (e) {
      var m = this.metaDeEntidad(e);
      if (m && m.kind) return m.kind;
      return (e.solicitud || '').indexOf('Canal de activación:') === 0 ? 'activacion' : 'acudida';
    },
    etiquetaCanal: function (code) {
      var o = CANAL_ACTIVACION_V3.filter(function (c) { return c.value === code; })[0];
      return o ? o.label : (CANAL_ACTIVACION_LEGADO_V3[code] || code);
    },
    // Etiquetas de los canales de la entidad, para la tarjeta y el Cierre.
    canalDeEntidad: function (e) {
      var m = this.metaDeEntidad(e);
      var codes = [];
      if (m && m.activationChannels && m.activationChannels.length) {
        codes = m.activationChannels;
      } else if ((e.solicitud || '').indexOf('Canal de activación:') === 0) {
        codes = e.solicitud.replace('Canal de activación:', '').split(',').map(function (c) { return c.trim(); }).filter(function (c) { return c; });
      }
      var self = this;
      return codes.map(function (c) { return self.etiquetaCanal(c); });
    },
    // REG-38*/D20: si el backend responde 409 con readOnly, el Caso ya dejó de ser borrador (el vaciado
    // corrió en otra pestaña o en otro guardado): la sección pasa a solo lectura.
    marcarSoloLectura: function (res) {
      if (res && res.data && res.data.readOnly) {
        this.entitiesReadOnly = true;
        this.barreraParaEntidad = null;
        this.recargarEntidadesBarreras();
      }
    },
    eliminarEntidad: function (relId) {
      if (!window.confirm('¿Eliminar esta entidad del registro? También se eliminarán sus barreras asociadas.')) return;
      var self = this;
      fetch('/api/v1/victim-cases-v3/entidades/' + relId, { method: 'DELETE' })
        .then(function (r) { return r.json().catch(function () { return {}; }).then(function (data) { return { ok: r.ok, data: data }; }); })
        .then(function (res) { if (!res.ok) self.marcarSoloLectura(res); self.recargarEntidadesBarreras(); });
    },

    // ── GS § D Barreras -- sub-repeater de Entidades (instrucción 2026-10-15) ──
    // Cada barrera cuelga de una entidad ya agregada: NO repite la búsqueda
    // geográfica, hereda EntityBranchID/sector/ubicación de esa entidad --
    // igual que "Gestión institucional" asocia barrier_v2 a un entity_case
    // (misma pareja case_id+entity_branch_id, ver entity_case_repository.go
    // entityCaseWithRelationsSQL § barrerasActivasCount). Verificado contra
    // el esquema real de dinamic-form (tabla repeater_group): el form de
    // Seguimiento en sí mismo tiene "Barreras identificadas" y "Entidades a
    // las que ya acudió" como dos repeaters HERMANOS e independientes (cada
    // uno con su propio form_section_id, sin relación padre/hijo) -- la
    // asociación ahí es solo la coincidencia implícita de columnas, nunca un
    // repeater anidado de verdad. El anidamiento real (0, 1 o más barreras
    // por entidad, sin repetir el buscador) es una decisión propia de V3,
    // no una réplica de cómo Seguimiento arma su UI.
    barrerasDeEntidad: function (e) {
      return this.barreras.filter(function (b) { return b.entityBranchId === e.entityBranchId; });
    },
    toggleEnArray: function (arr, value) {
      var idx = arr.indexOf(value);
      if (idx >= 0) arr.splice(idx, 1); else arr.push(value);
    },

    // ── Formulario de barrera de Seguimiento v4 (REG-18) ─────────────────────
    // Se renderiza recorriendo BARRERA_V4_V3 (textos, opciones y condiciones
    // salen del seed de Seguimiento; ver scripts/gen_registro_v3_barrera_catalog.py).
    barreraCumpleCuando: function (q) {
      var c = q.cuando;
      if (!c) return true;
      var val = this.nuevaBarrera.v[c.pregunta];
      if (c.op === 'INCLUDES') return Array.isArray(val) && val.indexOf(c.valor) >= 0;
      return val === c.valor;
    },
    barreraPreguntasVisibles: function () {
      var self = this;
      return BARRERA_V4_V3.filter(function (q) { return self.barreraCumpleCuando(q); });
    },
    barreraOpciones: function (q) {
      if (q.fuente === 'departamentos') return this.geoDepartamentos;
      if (q.fuente === 'municipios') return this.barreraMunicipios;
      return q.opciones || [];
    },
    barreraTextoLargo: function (q) {
      return ['descripcion', 'perspectiva', 'gestionTexto'].indexOf(q.key) >= 0;
    },
    // Al cambiar una respuesta se limpia lo que dejó de aplicar (igual que el
    // formulario de Seguimiento al ocultar las preguntas dependientes).
    setBarrera: function (q, valor) {
      var v = this.nuevaBarrera.v;
      v[q.key] = valor;
      var self = this;
      if (q.key === 'sector') {
        BARRERA_V4_V3.forEach(function (p) {
          if (p.cuando && p.cuando.pregunta === 'sector' && p.tipo !== 'info') v[p.key] = p.tipo === 'multiple' ? [] : '';
        });
      }
      if (q.key === 'depto') {
        v.municipio = '';
        this.barreraMunicipios = [];
        if (valor) {
          fetch('/api/v1/locations/cities?department_id=' + encodeURIComponent(valor))
            .then(function (r) { return r.json(); })
            .then(function (data) { self.barreraMunicipios = Array.isArray(data) ? data : []; })
            .catch(function () { self.barreraMunicipios = []; });
        }
      }
    },
    toggleBarreraOpcion: function (q, valor) {
      var v = this.nuevaBarrera.v;
      this.toggleEnArray(v[q.key], valor);
      // La perspectiva profesional solo aplica con la gestión «alerta».
      if (q.key === 'gestion' && v.gestion.indexOf('alerta') === -1) v.perspectiva = '';
    },
    barreraVacio: function (q) {
      var val = this.nuevaBarrera.v[q.key];
      return Array.isArray(val) ? val.length === 0 : !String(val || '').trim();
    },
    // Preguntas obligatorias y visibles que siguen sin respuesta.
    barreraFaltantes: function () {
      var self = this;
      return this.barreraPreguntasVisibles().filter(function (q) { return q.required && q.tipo !== 'info' && self.barreraVacio(q); });
    },
    // Respuestas del formulario -> campos de repository.VictimCaseBarrierInput.
    // Las preguntas por sector (tipos, institución, «otra») se mapean según el
    // sector elegido; las cuatro transversales solo existen con «transversales».
    armarPayloadBarrera: function () {
      var v = this.nuevaBarrera.v;
      var porSector = {
        salud: { Specific: v.tiposSalud, Institution: v.instSalud, OtherDetail: v.otraSalud },
        justicia: { Specific: v.tiposJusticia, Institution: v.instJusticia, OtherDetail: v.otraJusticia },
        proteccion: { Specific: v.tiposProteccion, Institution: v.instProteccion, OtherDetail: v.otraProteccion },
      }[v.sector] || { Specific: [], Institution: '', OtherDetail: '' };
      var transversal = v.sector === 'transversales';
      return {
        EntityBranchID: Number(this.nuevaBarrera.EntityBranchID),
        SectorCode: v.sector,
        Specific: porSector.Specific.slice(), Institution: porSector.Institution, OtherDetail: porSector.OtherDetail,
        Institutional: transversal ? v.talento.slice() : [], Economic: transversal ? v.eco.slice() : [],
        Territorial: transversal ? v.terr.slice() : [], Differential: transversal ? v.diferencial.slice() : [],
        Date: v.fecha, Official: v.funcionario, Department: v.depto, City: v.municipio, Zone: v.vereda,
        Description: v.descripcion, Management: v.gestion.slice(),
        ProfessionalPerspective: v.gestion.indexOf('alerta') >= 0 ? v.perspectiva : '',
        ManagementDescription: v.gestionTexto,
      };
    },
    abrirBarreraParaEntidad: function (e) {
      this.nuevaBarrera = this.barreraVacia();
      this.nuevaBarrera.EntityBranchID = e.entityBranchId;
      this.barreraMunicipios = [];
      this.barreraParaEntidad = e.entityBranchId;
      this.errorBarrera = '';
    },
    cerrarBarreraParaEntidad: function () {
      this.barreraParaEntidad = null;
      this.nuevaBarrera = this.barreraVacia();
      this.barreraMunicipios = [];
    },
    agregarBarrera: function () {
      var self = this;
      if (!this.nuevaBarrera.EntityBranchID) return;
      var faltantes = this.barreraFaltantes();
      if (faltantes.length) {
        self.errorBarrera = 'Faltan campos obligatorios: ' + faltantes.map(function (q) { return q.label; }).join('; ');
        return;
      }
      self.errorBarrera = '';
      fetch('/api/v1/victim-cases-v3/' + this.icode + '/barreras', {
        method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(this.armarPayloadBarrera()),
      })
        .then(function (r) { return r.json().then(function (data) { return { ok: r.ok, data: data }; }); })
        .then(function (res) {
          if (!res.ok) { self.errorBarrera = res.data.error || 'No se pudo agregar la barrera'; self.marcarSoloLectura(res); return; }
          self.cerrarBarreraParaEntidad();
          self.recargarEntidadesBarreras();
        });
    },
    eliminarBarrera: function (id) {
      if (!window.confirm('¿Eliminar esta barrera del registro?')) return;
      var self = this;
      fetch('/api/v1/victim-cases-v3/barreras/' + id, { method: 'DELETE' })
        .then(function (r) { return r.json().catch(function () { return {}; }).then(function (data) { return { ok: r.ok, data: data }; }); })
        .then(function (res) { if (!res.ok) self.marcarSoloLectura(res); self.recargarEntidadesBarreras(); });
    },

    // ── Cierre ────────────────────────────────────────────────────────────
    // ── Párrafo síntesis (D3-D5, 2026-10-15) ────────────────────────────────
    // Reescrito para seguir registro-v3-parrafo-sintesis.md literal (obs.
    // #18 del cliente, la fuente de mayor prioridad de todo Registro V3).
    // Antes solo tenía una versión mínima de § A + un dump crudo de HV_1 --
    // § B (datos personales) y § D (agresor) no existían, y el párrafo
    // mezclaba resultado de Tamizaje/remisión (que NO son parte de este
    // spec) en su lugar. Esa info de Tamizaje/remisión ya se muestra aparte
    // (alertas de TZ, D10; mensajeCierre de Finalizar, C2/C3) -- no hace
    // falta repetirla acá.
    //
    // etiquetaExcluye: true si la etiqueta resuelta de fieldId contiene
    // alguno de los fragmentos de exclusión (la plantilla del cliente
    // describe las excepciones por texto, no por código de enum).
    etiquetaExcluye: function (fieldId, exclusiones) {
      var et = this.etiquetaSeleccionada(fieldId).toLowerCase();
      if (!et) return true;
      return exclusiones.some(function (frag) { return et.indexOf(frag.toLowerCase()) >= 0; });
    },
    // Encontrado en vivo con Playwright (2026-10-15): GetAnswers devuelve las
    // fechas en el formato crudo de Go ("2026-10-01 00:00:00 +0000 UTC") una
    // vez que un Caso se reabre desde el backend -- antes de esta ronda
    // nadie mostraba una fecha cruda en prosa (solo se usaban como `value`
    // de un <input type="date">, que el navegador tolera igual). sintesisC
    // es la primera vez que se concatena una fecha directo en una frase, así
    // que es la primera vez que se nota. Se arregla acá, en el borde de
    // visualización, no en GetAnswers (que sigue alimentando el <input
    // type="date"> correctamente tal cual).
    limpiarFecha: function (raw) {
      var m = /^\d{4}-\d{2}-\d{2}/.exec(raw || '');
      return m ? m[0] : (raw || '');
    },
    calcularEdad: function (fechaISO) {
      if (!fechaISO) return null;
      var nacimiento = new Date(fechaISO);
      if (isNaN(nacimiento.getTime())) return null;
      var hoy = new Date();
      var edad = hoy.getFullYear() - nacimiento.getFullYear();
      var noHaCumplido = (hoy.getMonth() < nacimiento.getMonth()) ||
        (hoy.getMonth() === nacimiento.getMonth() && hoy.getDate() < nacimiento.getDate());
      if (noHaCumplido) edad--;
      return edad >= 0 ? edad : null;
    },
    // § A — Contacto e identificación. Se omite si no hay nombres/apellidos.
    sintesisA: function () {
      if (!this.answers['CI_2']) return '';
      var tratamientoTexto = { 'Sra.': 'La Señora', 'Sr.': 'El Señor', 'Sre.': 'Le Señore' }[this.tratamiento()] || '';
      var identitario = this.answers['CI_3'] ? this.answers['CI_3'] + ' ' : '';
      var frase = (tratamientoTexto ? tratamientoTexto + ' ' : '') + identitario + this.answers['CI_2'] +
        ', identificada/e/o con ' + (this.etiquetaSeleccionada('CI_5') || 'documento') + ' Número ' + (this.answers['CI_6'] || 's/n') + '.';
      if (this.answers['CI_7'] === 'true') frase += ' Reconoce discapacidad.';
      if (this.answers['CI_9'] === 'true') frase += ' Requiere intérprete de idiomas.';
      return frase;
    },
    // § B — Datos personales (4 bloques independientes).
    sintesisB: function () {
      var self = this;
      var bloques = [];
      var nombre = this.answers['CI_2'] || '';

      var edad = this.calcularEdad(this.answers['DP_1']);
      if (edad !== null) {
        var nacCode = this.codigoSeleccionado('DP_5');
        var nacionalidad = nacCode === 'co' ? 'Colombiana(o)' : nacCode === 'ap' ? 'Apátrida' :
          nacCode === 'ex' ? 'de ' + (this.etiquetaSeleccionada('DP_6') || 'nacionalidad extranjera') : '';
        var b1 = nombre + ' tiene ' + edad + ' años' + (nacionalidad ? ', nacional ' + nacionalidad : '') + '.';
        if (this.answers['DP_27']) b1 += ' Reside en ' + this.answers['DP_27'] + (this.answers['DP_28'] ? ' en ' + this.answers['DP_28'] : '') + (this.etiquetaSeleccionada('DP_30') ? ' en ' + this.etiquetaSeleccionada('DP_30') : '') + '.';
        bloques.push(b1);
      }

      if (!this.etiquetaExcluye('DP_8', ['otro'])) {
        var b2 = 'Se reconoce ' + this.etiquetaSeleccionada('DP_8');
        if (!this.etiquetaExcluye('DP_10', [])) b2 += ', fue asignado al nacer ' + this.etiquetaSeleccionada('DP_10');
        if (!this.etiquetaExcluye('DP_9', ['otro', 'no estoy segura', 'no informa'])) b2 += ' con orientación sexual ' + this.etiquetaSeleccionada('DP_9');
        bloques.push(b2 + '.');
      }

      var b3partes = [];
      if (!this.etiquetaExcluye('DP_12', ['no se auto reconoce', 'no informa'])) b3partes.push('Se identifica ' + this.etiquetaSeleccionada('DP_12'));
      if (this.answers['DP_14'] === 'true') b3partes.push((b3partes.length ? ' y se reconoce' : 'Se reconoce') + ' como persona campesina');
      var poblaciones = (this.answers['DP_11'] || '').split(',').filter(function (v) { return v; })
        .map(function (ic) { var o = self.opcionesPara(self.camposPorId['DP_11']).filter(function (o) { return o.icode === ic; })[0]; return o ? self.etiqueta(self.camposPorId['DP_11'], o) : ''; })
        .filter(function (l) { return l && l.toLowerCase().indexOf('ninguno') === -1; });
      if (b3partes.length) {
        var b3 = b3partes.join('') + '.';
        if (poblaciones.length) b3 += ' Pertenece a población de especial protección constitucional ' + poblaciones.join(', ') + '.';
        bloques.push(b3);
      } else if (poblaciones.length) {
        bloques.push('Pertenece a población de especial protección constitucional ' + poblaciones.join(', ') + '.');
      }

      var b4partes = [];
      if (this.answers['DP_15']) b4partes.push('Persona ' + this.etiquetaSeleccionada('DP_15'));
      if (!this.etiquetaExcluye('DP_16', ['no informa'])) {
        var esc = this.etiquetaSeleccionada('DP_16');
        b4partes.push((b4partes.length ? ' ' : 'Persona ') + (esc.toLowerCase() === 'sin educación formal' ? esc : 'con ' + esc));
      }
      if (b4partes.length) {
        var b4 = b4partes.join('') + '.';
        if (this.answers['DP_18'] && this.answers['DP_17']) b4 += ' ' + this.etiquetaSeleccionada('DP_18') + ' en su ocupación ' + this.etiquetaSeleccionada('DP_17') + '.';
        if (!this.etiquetaExcluye('DP_23', ['no informa'])) b4 += ' Su vivienda es ' + this.etiquetaSeleccionada('DP_23');
        if (!this.etiquetaExcluye('DP_24', ['no informa', 'no sabe'])) b4 += ' en estrato ' + this.etiquetaSeleccionada('DP_24') + '.';
        else if (b4.slice(-1) !== '.') b4 += '.';
        var personasCargo = (this.answers['DP_25'] || '').split(',').filter(function (v) { return v; })
          .map(function (ic) { var o = self.opcionesPara(self.camposPorId['DP_25']).filter(function (o) { return o.icode === ic; })[0]; return o ? self.etiqueta(self.camposPorId['DP_25'], o) : ''; }).filter(function (l) { return l; });
        if (personasCargo.length) b4 += ' Tiene personas o animales a su cargo: ' + personasCargo.join(', ') + '.';
        bloques.push(b4);
      }

      return bloques.join(' ');
    },
    // § C — Hechos de violencia. Se omite si no hay tipo de violencia ni relato.
    sintesisC: function () {
      if (!this.answers['HV_1'] && !this.answers['HV_9']) return '';
      var self = this;
      var tipos = (this.answers['HV_9'] || '').split(',').filter(function (v) { return v; })
        .map(function (ic) { var o = self.opcionesPara(self.camposPorId['HV_9']).filter(function (o) { return o.icode === ic; })[0]; return o ? self.etiqueta(self.camposPorId['HV_9'], o) : ''; }).filter(function (l) { return l; });
      var subtipos = (this.answers['HV_10'] || '').split(',').filter(function (v) { return v; })
        .map(function (ic) { var o = self.opcionesPara(self.camposPorId['HV_10']).filter(function (o) { return o.icode === ic; })[0]; return o ? self.etiqueta(self.camposPorId['HV_10'], o) : ''; }).filter(function (l) { return l; });
      var ambitos = (this.answers['HV_11'] || '').split(',').filter(function (v) { return v; })
        .map(function (ic) { var o = self.opcionesPara(self.camposPorId['HV_11']).filter(function (o) { return o.icode === ic; })[0]; return o ? self.etiqueta(self.camposPorId['HV_11'], o) : ''; }).filter(function (l) { return l; });

      var frase = 'Hecho de violencia' + (tipos.length ? ' de tipo(s) ' + tipos.join(', ') : '') +
        (subtipos.length ? ', ' + subtipos.join(', ') : '') +
        (ambitos.length ? ', ocurrido en ámbito(s) ' + ambitos.join(', ') : '') +
        (this.etiquetaSeleccionada('HV_8') ? ' en escenario ' + this.etiquetaSeleccionada('HV_8') : '') + '.';
      if (this.answers['HV_2']) {
        frase += ' Sucedido el ' + this.limpiarFecha(this.answers['HV_2']) + (this.answers['HV_3'] ? ' ' + this.answers['HV_3'] : '') +
          (this.answers['HV_4'] ? ' en ' + this.answers['HV_4'] : '') + (this.answers['HV_5'] ? ' en ' + this.answers['HV_5'] : '') +
          (this.etiquetaSeleccionada('HV_6') ? ' en ' + this.etiquetaSeleccionada('HV_6') : '') + '.';
      }
      if (this.answers['HV_13'] === 'true') frase += ' Motivado por razones de género.';
      if (this.etiquetaSeleccionada('HV_14')) frase += ' ' + this.etiquetaSeleccionada('HV_14') + '.';
      if (this.answers['HV_1']) frase += ' ' + this.answers['HV_1'];
      return frase;
    },
    // § D — Características del agresor. Se omite si no hay relación ni número de agresores.
    sintesisD: function () {
      if (!this.answers['AG_3'] && !this.answers['AG_1']) return '';
      var frase = '';
      if (!this.etiquetaExcluye('AG_1', ['no se puede determinar', 'no informa'])) {
        frase = 'Responsabiliza a ' + this.etiquetaSeleccionada('AG_1') + ' persona(s) de la agresión.';
      }
      if (this.etiquetaSeleccionada('AG_2')) frase += (frase ? ' Principalmente a' : 'A') + ' ' + this.etiquetaSeleccionada('AG_2').toLowerCase() + (this.answers['AG_7'] ? ', de nombre ' + this.answers['AG_7'] : '');
      if (this.etiquetaSeleccionada('AG_6')) frase += (frase.slice(-1) === '.' ? ' ' : ', ') + this.etiquetaSeleccionada('AG_6');
      if (!this.etiquetaExcluye('AG_3', ['sin relación', 'otro'])) frase += ', con quien tiene una relación de ' + this.etiquetaSeleccionada('AG_3');
      if (!this.etiquetaExcluye('AG_4', ['otro', 'no conoce'])) frase += ', cuya ocupación es ' + this.etiquetaSeleccionada('AG_4');
      return frase.trim() + (frase.slice(-1) !== '.' ? '.' : '');
    },
    parrafoSintesis: function () {
      return [this.sintesisA(), this.sintesisB(), this.sintesisC(), this.sintesisD()].filter(function (s) { return s; }).join(' ');
    },
    // C2/C3 (flow-E15): "generar Denuncia fácil si corresponde" / "confirmar
    // remisión si aplica" -- no hay integración real con la FGN ni un
    // workflow de remisión separado de la respuesta GS_6 para reutilizar
    // (ver FinalizeResultV3 en victim_case_v3_service.go), así que se
    // implementa como una confirmación explícita al agente en vez de un
    // envío silencioso.
    confirmacionFinalizar: function (data) {
      var partes = [];
      if (data.denunciaFacilSolicitada) partes.push('Se solicitó la Denuncia fácil ante la FGN.');
      if (data.equiposRemisionCSV) {
        var opts = this.opcionesPara(this.camposPorId['GS_6']);
        var icodes = data.equiposRemisionCSV.split(',').filter(function (v) { return v; });
        var nombres = icodes.map(function (ic) {
          var o = opts.filter(function (o) { return o.icode === ic; })[0];
          return o ? o.label : ic;
        });
        if (nombres.length) partes.push('Se confirma remisión a: ' + nombres.join(', ') + '.');
      }
      return partes.join(' ');
    },
    // REG-37: primero se guarda Gestión Salvia (con las mismas validaciones de «Guardar sección»); si no se guarda,
    // no se finaliza. REG-21: se lee la respuesta sin lanzar (un 401 puede venir sin cuerpo) y un 401 se avisa
    // con un mensaje claro, sin recargar ni perder lo diligenciado.
    finalizarRegistro: function () {
      var self = this;
      if (this.guardando || this.finalizando) return Promise.resolve();
      this.errorFinalizar = '';
      this.finalizando = true;
      return this.guardarModulo('GS', true)
        .then(function (guardado) {
          if (!guardado) {
            // El recordatorio de faltantes de GS (si lo hubo) se repite junto al botón de finalizar.
            if (!self.errorGuardado && self.mensajeGuardado) self.errorFinalizar = self.mensajeGuardado;
            return;
          }
          return fetch('/api/v1/victim-cases-v3/' + self.icode + '/finalizar', { method: 'POST' })
            .then(function (r) { return self.leerRespuesta(r); })
            .then(function (res) {
              if (res.ok) {
                self.mensajeCierre = self.confirmacionFinalizar(res.data);
                self.etapa = 'cierre';
                return;
              }
              self.errorFinalizar = self.mensajeErrorFinalizar(res);
            })
            .catch(function () {
              self.errorFinalizar = 'No se pudo finalizar: no hay conexión con el servidor. Lo que diligenciaste sigue en pantalla; intenta de nuevo.';
            });
        })
        .finally(function () { self.finalizando = false; });
    },
    // Texto para la persona cuando el servidor rechazó el cierre (REG-37 / REG-21).
    mensajeErrorFinalizar: function (res) {
      var self = this;
      var d = res.data || {};
      if (res.status === 401 || res.status === 403) {
        return 'No se pudo finalizar: tu sesión venció o no tienes permiso. Inicia sesión de nuevo en otra pestaña y vuelve a pulsar «Finalizar registro»: lo que diligenciaste sigue en esta pantalla.';
      }
      var faltan = d.modulosIncompletos || [];
      if (res.status === 400 && faltan.length) {
        var nombres = faltan.map(function (id) {
          var m = self.modulos.filter(function (x) { return x.id === id; })[0];
          return m ? m.label : id;
        });
        var texto = (d.error || 'No se puede finalizar') + ' — falta diligenciar: ' + nombres.join(', ');
        var gs = (d.camposFaltantes || []).filter(function (k) { return k.indexOf('GS_') === 0; });
        if (gs.length) {
          texto += '. Gestión Salvia, falta: ' + gs.map(function (k) {
            var c = self.camposPorId[k];
            return (c && c.label) || GS_ETIQUETAS_FINALIZAR_V3[k] || k;
          }).join('; ');
        }
        return texto;
      }
      if (res.status === 400 && d.error) return d.error;
      return 'No se pudo finalizar por un error del servidor. Lo que diligenciaste sigue en pantalla: intenta de nuevo y, si el problema continúa, avisa a soporte.';
    },
  },
});
