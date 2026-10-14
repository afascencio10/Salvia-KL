// Registro V3 — app principal. Sin dinamic-form (decisión A1): Habeas Data
// -> Cabezote (búsqueda por documento y/o teléfono, D7) -> Formulario (7
// módulos CI/DP/CA/HV/AG/TZ/GS) -> Cierre. Usa registro-v3-campos.js como
// catálogo de preguntas y GET /victim-cases-v3/enums como catálogo de
// opciones (icode real de victim_case_form2_enums).
var RegistroV3App = Vue.createApp({
  delimiters: ['${', '}'],
  data: function () {
    return {
      etapa: 'habeas', // habeas | cabezote | formulario | cierre
      habeasAceptado: false,

      buscar: { tipoDocumento: '', documento: '', telefono: '' },
      resultadosBusqueda: [],
      buscando: false,
      errorBusqueda: '',

      icode: null,
      resumenCaso: null,
      answers: {},
      entidades: [],
      barreras: [],
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
        { id: 'GS', label: 'Gestión Salvia', campos: CAMPOS_GS_A_V3 },
      ],
      guardando: false,
      mensajeGuardado: '',

      // Tamizaje
      tzRevelarResto: false,
      tzResultado: null,
      tzGuardando: false,

      // GS § B remisión
      gsRequiereRemision: false,
      gsJustificacionRemision: '',

      // GS § C entidades
      nuevaEntidad: this.entidadVacia(),
      // GS § D barreras
      nuevaBarrera: this.barreraVacia(),
    };
  },

  computed: {
    moduloActualObj: function () {
      var self = this;
      return this.modulos.filter(function (m) { return m.id === self.moduloActivo; })[0];
    },
    equiposRemision: function () { return EQUIPOS_REMISION_V3; },
    equiposSeleccionados: function () {
      return (this.answers['GS_6'] || '').split(',').filter(function (v) { return v; });
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
  },

  methods: {
    entidadVacia: function () {
      return { Kind: 'acudida', SectorCode: '', Department: '', City: '', Town: '', Name: '', FileNumber: '', ReportDate: '', ActivationChannel: '', CompletedActions: [], PendingActions: [], OtherInfo: '' };
    },
    barreraVacia: function () {
      return { SectorCode: '', Specific: [], OtherDetail: '', EntityName: '', Institutional: [], Economic: [], Territorial: [], Differential: [], Date: '', Official: '', Description: '', Management: [] };
    },

    // ── Habeas data / Cabezote ──────────────────────────────────────────
    aceptarHabeasData: function () {
      this.habeasAceptado = true;
      this.etapa = 'cabezote';
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
        .then(function (r) { return r.json(); })
        .then(function (data) { self.resultadosBusqueda = data.registros || []; })
        .catch(function () { self.errorBusqueda = 'Error al buscar.'; })
        .finally(function () { self.buscando = false; });
    },

    crearNuevoRegistro: function () {
      var self = this;
      fetch('/api/v1/victim-cases-v3/draft', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          tipo_documento: this.buscar.tipoDocumento,
          documento: this.buscar.documento,
          telefono: this.buscar.telefono,
        }),
      })
        .then(function (r) { return r.json(); })
        .then(function (data) {
          if (data.icode) { self.abrirCaso(data.icode); }
        });
    },

    abrirCaso: function (icode) {
      var self = this;
      fetch('/api/v1/victim-cases-v3/' + icode)
        .then(function (r) { return r.json(); })
        .then(function (data) {
          self.icode = icode;
          self.resumenCaso = data.summary || null;
          self.answers = data.answers || {};
          self.entidades = data.entities || [];
          self.barreras = data.barriers || [];
          self.etapa = 'formulario';
          self.moduloActivo = 'CI';
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
        // busca override en cualquiera de las categorías candidatas
        var cats = campo.enumCategories || Object.keys(ENUM_LABELS_OVERRIDE_V3);
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
    setVal: function (id, v) { this.answers[id] = v; },
    getBool: function (id) { return this.answers[id] === 'true'; },
    setBool: function (id, v) { this.answers[id] = v ? 'true' : 'false'; },

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
        actuales.push(icode);
      }
      this.answers[id] = actuales.join(',');
    },

    pseudoCondicion: function (nombre) {
      switch (nombre) {
        case '__gestacionAplica': return true;
        case '__haySubtipo': return !!(this.answers['HV_9'] && this.answers['HV_9'].length);
        case '__ambitoLaboral': return this.multipleIncluyeTexto('HV_11', 'laboral');
        case '__agresorConocido':
          var et = this.etiquetaSeleccionada('AG_2').toLowerCase();
          return !!et && et.indexOf('desconocid') === -1;
        default: return true;
      }
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
      if (campo.condicionalEnumLabel) {
        var et = this.etiquetaSeleccionada(campo.condicionalEnumLabel.campo);
        return et.indexOf(campo.condicionalEnumLabel.contains) >= 0;
      }
      if (campo.subtipoDe) {
        return !!(this.answers[campo.subtipoDe] && this.answers[campo.subtipoDe].length);
      }
      return true;
    },

    camposVisibles: function (moduloId) {
      var self = this;
      var mod = this.modulos.filter(function (m) { return m.id === moduloId; })[0];
      if (!mod) return [];
      return mod.campos.filter(function (c) { return self.mostrarCampo(c); });
    },

    // ── Guardar sección ──────────────────────────────────────────────────
    guardarModulo: function (moduloId) {
      var self = this;
      var payload = {};
      Object.keys(this.answers).forEach(function (k) {
        if (k.indexOf(moduloId + '_') === 0) payload[k] = self.answers[k];
      });
      this.guardando = true;
      this.mensajeGuardado = '';
      fetch('/api/v1/victim-cases-v3/' + this.icode + '/seccion/' + moduloId, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ answers: payload }),
      })
        .then(function (r) { return r.json(); })
        .then(function (data) {
          self.triadaCompleta = !!data.triadaCompleta;
          self.mensajeGuardado = 'Sección guardada.';
        })
        .catch(function () { self.mensajeGuardado = 'Error al guardar la sección.'; })
        .finally(function () { self.guardando = false; });
    },

    irAModulo: function (id) { this.moduloActivo = id; this.mensajeGuardado = ''; },

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

      var tzAnswers = {};
      TAMIZAJE_GENERALES_V3.forEach(function (c) { tzAnswers[c.id] = self.answers[c.id] === 'true' ? 'true' : 'false'; });
      tzAnswers['TZ_8'] = this.answers['TZ_8'] === 'true' ? 'true' : 'false';
      var lista = wasPartner ? TAMIZAJE_PAREJA_V3 : TAMIZAJE_NO_PAREJA_V3;
      var prefijo = wasPartner ? 'TZP_' : 'TZNP_';
      lista.forEach(function (_, idx) {
        var id = prefijo + (idx + 1);
        tzAnswers[id] = self.answers[id] === 'true' ? 'true' : 'false';
      });

      this.tzGuardando = true;
      fetch('/api/v1/victim-cases-v3/' + this.icode + '/tamizaje', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          answers: tzAnswers, wasPartner: wasPartner, esOsigd: esOsigd,
          subtipoFeminicidio: subtipoFem, subtipoHomicidioAbyecto: subtipoHom,
        }),
      })
        .then(function (r) { return r.json(); })
        .then(function (data) { self.tzResultado = data; })
        .finally(function () { self.tzGuardando = false; });
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

    // ── GS § C Entidades ──────────────────────────────────────────────────
    agregarEntidad: function () {
      var self = this;
      var body = Object.assign({}, this.nuevaEntidad);
      fetch('/api/v1/victim-cases-v3/' + this.icode + '/entidades', {
        method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(body),
      })
        .then(function (r) { return r.json(); })
        .then(function (data) {
          body.ICode = data.icode;
          self.entidades.push(body);
          self.nuevaEntidad = self.entidadVacia();
        });
    },
    eliminarEntidad: function (icode, idx) {
      var self = this;
      fetch('/api/v1/victim-cases-v3/entidades/' + icode, { method: 'DELETE' })
        .then(function () { self.entidades.splice(idx, 1); });
    },

    // ── GS § D Barreras ──────────────────────────────────────────────────
    barreraCategoriaSector: function () {
      return SECTOR_A_CATEGORIA_BARRERA_V3[this.nuevaBarrera.SectorCode] || null;
    },
    opcionesBarreraSector: function () {
      var cat = this.barreraCategoriaSector();
      if (!cat) return [];
      var self = this;
      return (this.enumCatalog[cat] || []).map(function (o) {
        return { code: o.code, icode: o.icode, label: (ENUM_LABELS_OVERRIDE_V3[cat] && ENUM_LABELS_OVERRIDE_V3[cat][o.code]) || o.label };
      });
    },
    opcionesTransversal: function (enumCategory) {
      var self = this;
      return (this.enumCatalog[enumCategory] || []).map(function (o) {
        return { code: o.code, icode: o.icode, label: (ENUM_LABELS_OVERRIDE_V3[enumCategory] && ENUM_LABELS_OVERRIDE_V3[enumCategory][o.code]) || o.label };
      });
    },
    toggleEnArray: function (arr, icode) {
      var idx = arr.indexOf(icode);
      if (idx >= 0) arr.splice(idx, 1); else arr.push(icode);
    },
    categoriasTransversales: function () { return CATEGORIAS_BARRERA_TRANSVERSAL_V3; },
    opcionesGestionBarrera: function () { return this.opcionesTransversal(GESTION_BARRERA_ENUM_CATEGORY_V3); },
    agregarBarrera: function () {
      var self = this;
      var body = Object.assign({}, this.nuevaBarrera);
      fetch('/api/v1/victim-cases-v3/' + this.icode + '/barreras', {
        method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(body),
      })
        .then(function (r) { return r.json(); })
        .then(function (data) {
          body.ICode = data.icode;
          self.barreras.push(body);
          self.nuevaBarrera = self.barreraVacia();
        });
    },
    eliminarBarrera: function (icode, idx) {
      var self = this;
      fetch('/api/v1/victim-cases-v3/barreras/' + icode, { method: 'DELETE' })
        .then(function () { self.barreras.splice(idx, 1); });
    },

    // ── Cierre ────────────────────────────────────────────────────────────
    parrafoSintesis: function () {
      var partes = [];
      var nombres = this.answers['CI_2'] || '(sin nombre)';
      partes.push(nombres + ', identificada con ' + (this.etiquetaSeleccionada('CI_5') || 'documento') + ' N° ' + (this.answers['CI_6'] || 's/n') + '.');
      if (this.answers['HV_1']) partes.push('Hechos: ' + this.answers['HV_1']);
      var tipos = this.multipleIncluyeTexto('HV_9', '') ? '' : '';
      if (this.tzResultado) {
        partes.push('Tamizaje: ' + this.tzResultado.riskType + ', nivel de riesgo ' + this.tzResultado.riskLevel + ' (score ' + this.tzResultado.score + ').');
      }
      if (this.equiposSeleccionados.length) {
        var self = this;
        var nombresEquipos = this.equiposSeleccionados.map(function (v) {
          var eq = EQUIPOS_REMISION_V3.filter(function (e) { return e.value === v; })[0];
          return eq ? eq.label : v;
        });
        partes.push('Se recomienda remisión a: ' + nombresEquipos.join(', ') + '.');
      }
      return partes.join(' ');
    },
    finalizarRegistro: function () {
      var self = this;
      fetch('/api/v1/victim-cases-v3/' + this.icode + '/finalizar', { method: 'POST' })
        .then(function (r) { return r.json(); })
        .then(function () { self.etapa = 'cierre'; });
    },
  },
});
