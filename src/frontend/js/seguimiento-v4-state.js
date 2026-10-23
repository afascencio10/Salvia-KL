/**
 * Catálogos del formulario de seguimiento V.4.
 * Las preguntas no traen opciones en la base: leen formState con las rutas
 * que dejó el seed (hasEntidades, currentEntities, identEntidades, statesColombia).
 */
(function () {
    var IDENT_GROUP = '2bda4735-3927-4634-a666-c7f97949be0a';
    var Q_IDENT = {
        sector: '55faf667-6ea2-40ff-871d-3321cd2e1f31',
        depto: '7bc62f4c-f2d9-4900-b356-0ef6751e9351',
        municipio: 'ce975302-ed4e-4451-a18f-3ca9322417f5',
        ciudad: '1f5a960c-0fb6-4e6a-9afd-330b7b539630',
        entidad: '34a1fdb4-5876-4a51-ba04-5c1c169dbcf8',
        otra: 'e64c1a56-ff30-431d-8dd4-f98ccd4bff39',
    };
    var BARRERAS = {
        '52387e43-9d64-412d-83b5-b2b6cbbf02c8': {
            bucket: 'ident',
            depto: 'b0a00008-defe-49f0-892b-cf67e337bafd',
        },
        'c3561d3f-f207-4541-80b8-a323dce70f90': {
            bucket: 'seg',
            depto: 'df74f0d3-1685-4fde-9d42-57e997b12bef',
        },
    };

    function answerValue(entry, questionId) {
        var found = (entry.answers || []).find(function (a) { return a.questionId === questionId; });
        return found ? (found.value || '') : '';
    }

    function citiesOf(vm, departmentId) {
        return (vm.allCities || [])
            .filter(function (c) { return String(c.departmentId) === String(departmentId); })
            .map(function (c) { return { label: c.label, value: String(c.value) }; });
    }

    function slot(list, index) {
        while (list.length <= index) {
            list.push({ municipios: [], ciudades: [], entidades: [], obligaciones: [] });
        }
        return list[index];
    }

    function changed(vm, key, value, ready) {
        if (ready === false) return false;
        if (vm._v4cache[key] !== undefined && String(vm._v4cache[key]) === String(value)) return false;
        vm._v4cache[key] = value;
        return true;
    }

    function ensure(vm) {
        var fs = vm.formState;
        if (!fs.identEntidades) fs.identEntidades = [];
        if (!fs.entidadesConocidas) fs.entidadesConocidas = [];
        if (!fs.ident) fs.ident = { barrera: [] };
        if (!fs.seg) fs.seg = { barrera: [] };
        if (!fs.ident.barrera) fs.ident.barrera = [];
        if (!fs.seg.barrera) fs.seg.barrera = [];
        if (!vm._v4cache) vm._v4cache = {};
    }

    function knownEntities(formState) {
        var seen = {};
        var list = [];
        function add(value, label) {
            if (!value || seen[value]) return;
            seen[value] = true;
            list.push({ value: String(value), label: label || String(value) });
        }
        (formState.currentEntities || []).forEach(function (e) {
            add(e.entityBranchId, e.nombre);
        });
        (formState._identLabels || []).forEach(function (item) {
            add(item.value, item.label);
        });
        formState.entidadesConocidas = list;
    }

    function attachCase(formState, barriers) {
        var list = barriers || [];
        formState.hasEntidades = (formState.currentEntities || []).length > 0 ? 'true' : '';
        (formState.currentEntities || []).forEach(function (entity) {
            var obligations = entity.obligaciones || entity.obligacionesPendientes || [];
            entity.obligaciones = obligations;
            var mine = list.filter(function (b) {
                return b.entityBranchId != null && String(b.entityBranchId) === String(entity.entityBranchId);
            }).map(function (b) {
                return { value: b.id, label: b.barrierName };
            });
            entity.barreras = mine;
            entity.hasBarriers = mine.length > 0 ? 'true' : '';
        });
        knownEntities(formState);
    }

    function updateIdent(vm, entries) {
        var fs = vm.formState;
        var labels = [];
        entries.forEach(function (entry, idx) {
            var row = slot(fs.identEntidades, idx);
            var dept = answerValue(entry, Q_IDENT.depto);
            var city = answerValue(entry, Q_IDENT.municipio);
            var town = answerValue(entry, Q_IDENT.ciudad);
            var sector = answerValue(entry, Q_IDENT.sector);
            var entidad = answerValue(entry, Q_IDENT.entidad);
            var otra = answerValue(entry, Q_IDENT.otra);
            var key = 'ident:' + (entry.id || idx);
            var citiesReady = !dept || (vm.allCities || []).length > 0;

            if (changed(vm, key + ':d', dept, citiesReady)) {
                row.municipios = dept ? citiesOf(vm, dept) : [];
            }
            if (changed(vm, key + ':c', city)) {
                row.ciudades = [];
                if (city) {
                    fetch('/api/v1/locations/towns?city_id=' + encodeURIComponent(city))
                        .then(function (r) { return r.json(); })
                        .then(function (towns) { row.ciudades = Array.isArray(towns) ? towns : []; })
                        .catch(function () { row.ciudades = []; });
                }
            }
            if (changed(vm, key + ':e', town + '|' + sector)) {
                row.entidades = [];
                if (town && sector) {
                    fetch('/api/v1/entity-branches?town_code=' + encodeURIComponent(town) + '&sector=' + encodeURIComponent(sector))
                        .then(function (r) { return r.json(); })
                        .then(function (branches) {
                            row.entidades = (Array.isArray(branches) ? branches : []).map(function (b) {
                                return { label: b.name, value: String(b.id) };
                            });
                            row.entidades.push({ value: 'otra', label: 'Otra' });
                            knownEntities(fs);
                        })
                        .catch(function () { row.entidades = [{ value: 'otra', label: 'Otra' }]; });
                } else {
                    row.entidades = [{ value: 'otra', label: 'Otra' }];
                }
            }
            if (changed(vm, key + ':o', entidad)) {
                row.obligaciones = [];
                if (entidad && entidad !== 'otra') {
                    fetch('/api/v1/entity-branches/' + encodeURIComponent(entidad) + '/obligations')
                        .then(function (r) { return r.json(); })
                        .then(function (opts) { row.obligaciones = Array.isArray(opts) ? opts : []; })
                        .catch(function () { row.obligaciones = []; });
                }
            }
            if (entidad === 'otra' && otra) {
                labels.push({ value: otra, label: otra });
            } else if (entidad && entidad !== 'otra') {
                var opt = (row.entidades || []).find(function (o) { return String(o.value) === String(entidad); });
                labels.push({ value: String(entidad), label: opt ? opt.label : String(entidad) });
            }
        });
        fs._identLabels = labels;
        knownEntities(fs);
    }

    function updateBarrier(vm, cfg, entries) {
        var bucket = cfg.bucket === 'ident' ? vm.formState.ident.barrera : vm.formState.seg.barrera;
        entries.forEach(function (entry) {
            var idx = Math.max(0, (entry.iteration || 1) - 1);
            var row = slot(bucket, idx);
            var dept = answerValue(entry, cfg.depto);
            var key = cfg.bucket + ':' + (entry.id || idx);
            var ready = !dept || (vm.allCities || []).length > 0;
            if (changed(vm, key, dept, ready)) {
                row.municipios = dept ? citiesOf(vm, dept) : [];
            }
        });
    }

    function onAnswers(vm, payload) {
        ensure(vm);
        vm._lastAnswers = payload;
        var entries = (payload && payload.answers && payload.answers.repeaterEntries) || [];
        updateIdent(vm, entries.filter(function (e) { return e.repeaterGroupId === IDENT_GROUP; }));
        Object.keys(BARRERAS).forEach(function (groupId) {
            updateBarrier(vm, BARRERAS[groupId], entries.filter(function (e) { return e.repeaterGroupId === groupId; }));
        });
    }

    window.SeguimientoV4 = {
        formId: 'c4a8e2b1-6d3f-4a70-8e15-9f0b2c3d4e5f',
        attachCase: attachCase,
        onAnswers: onAnswers,
        ensure: ensure,
    };
})();
