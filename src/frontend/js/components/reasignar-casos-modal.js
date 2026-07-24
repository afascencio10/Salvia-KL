/**
 * reasignar-casos-modal
 * Modal para reasignar casos seleccionados desde casos-component.
 *
 * Método público:
 *   open(cases)  — abre el modal (M-01)
 *
 * Eventos internos:
 *   M-01  open              — apertura y resolución de equipo
 *   M-02  fetchAgentsByTeam — carga agentes por equipo (modo normal)
 *   M-02-C fetchAllAgents     — carga todos los agentes (contingencia cross-team)
 *   M-03  onAgentChange     — selección de agente
 *   M-04  close             — cancelar / cerrar
 *
 * Emite:
 *   closed      — modal cerrado sin guardar (M-04)
 *   reassigned  — reasignación exitosa (M-05)
 */
(function() {
    var vueApp = (typeof app !== 'undefined') ? app : (typeof home !== 'undefined' ? home : null);
    if (!vueApp) {
        console.error('[reasignar-casos-modal] No se encontró app/home.');
        return;
    }

    // Roles incluidos al listar agentes. Para habilitar 'op': ['ro', 'op']
    var AGENT_ROLE_CODES = ['ro'];

    // CONTINGENCIA: true = listar todos los agentes ro (sin filtro por equipo del caso).
    // false = comportamiento original (M-01 resolveTeam + M-02 fetchAgentsByTeam).
    var REASSIGN_CROSS_TEAM_CONTINGENCY = true;

    vueApp.component('reasignar-casos-modal', {
        delimiters: ['${', '}'],

        emits: ['closed', 'reassigned'],

        data: function() {
            return {
                visible: false,
                cases: [],
                resolvedTeam: '',
                agents: [],
                selectedAgentIcode: '',
                selectedAgent: null,
                loadingAgents: false,
                agentsError: null,
                saveError: null,
                saving: false,
                confirmVisible: false,
                crossTeamContingency: REASSIGN_CROSS_TEAM_CONTINGENCY,
            };
        },

        computed: {
            canConfirmReassign: function() {
                return !!this.selectedAgentIcode &&
                    !this.loadingAgents &&
                    !this.saving &&
                    this.agents.length > 0 &&
                    !(this.agentsError && !this.agents.length);
            },
        },

        methods: {

            // ------------------------------------------------------------
            // M-01 — Abrir modal
            // ------------------------------------------------------------

            open: function(cases) {
                if (!cases || !Array.isArray(cases) || cases.length === 0) {
                    return;
                }

                this.visible = true;
                this.cases = cases.slice();
                this.selectedAgentIcode = '';
                this.selectedAgent = null;
                this.agents = [];
                this.agentsError = null;
                this.saveError = null;
                this.loadingAgents = false;
                this.saving = false;
                this.resolvedTeam = '';

                if (REASSIGN_CROSS_TEAM_CONTINGENCY) {
                    this.fetchAllAgents();
                    return;
                }

                // ── MODO NORMAL (comentar rama anterior y descomentar esto para revertir) ──
                var team = this.resolveTeam(this.cases[0]);
                if (!team) {
                    this.agentsError = 'No se pudo determinar el equipo del caso seleccionado.';
                    return;
                }

                this.resolvedTeam = team;
                this.fetchAgentsByTeam(team);
            },

            resolveTeam: function(caseObj) {
                if (!caseObj) {
                    return '';
                }

                var caseTeam = (caseObj.caseTeam || '').trim();
                if (caseTeam) {
                    return this.normalizeTeamLabel(caseTeam);
                }

                var risk = (caseObj.riskStatus || '').toLowerCase();
                if (risk === 'bajo' || risk === 'moderado' || risk === 'medio') {
                    return 'Riesgo bajo';
                }
                if (risk === 'alto' || risk === 'extremo') {
                    return 'Riesgo alto';
                }

                return '';
            },

            normalizeTeamLabel: function(team) {
                var lower = team.toLowerCase();
                if (lower === 'riesgo alto') {
                    return 'Riesgo alto';
                }
                if (lower === 'riesgo bajo') {
                    return 'Riesgo bajo';
                }
                return team;
            },

            // ------------------------------------------------------------
            // M-02-C — Cargar todos los agentes (contingencia cross-team)
            // ------------------------------------------------------------

            fetchAllAgents: function() {
                var self = this;
                self.loadingAgents = true;
                self.agentsError = null;
                self.agents = [];
                self.selectedAgentIcode = '';
                self.selectedAgent = null;

                var params = new URLSearchParams();
                AGENT_ROLE_CODES.forEach(function(role) {
                    params.append('role', role);
                });

                fetch('/api/v1/equipo-operadores?' + params.toString())
                    .then(function(res) {
                        return res.json().then(function(data) {
                            return { ok: res.ok, data: data };
                        });
                    })
                    .then(function(result) {
                        if (!result.ok) {
                            self.agentsError = 'Error al cargar los agentes';
                            self.agents = [];
                            return;
                        }

                        var list = Array.isArray(result.data) ? result.data : [];
                        list.sort(function(a, b) {
                            return (a.fullName || '').localeCompare(b.fullName || '', 'es');
                        });

                        self.agents = list;
                        if (list.length === 0) {
                            self.agentsError = 'No hay agentes disponibles';
                        }
                    })
                    .catch(function() {
                        self.agentsError = 'Error al cargar los agentes';
                        self.agents = [];
                    })
                    .finally(function() {
                        self.loadingAgents = false;
                    });
            },

            // ------------------------------------------------------------
            // M-02 — Cargar agentes por equipo (modo normal)
            // ------------------------------------------------------------

            fetchAgentsByTeam: function(team) {
                var self = this;
                self.loadingAgents = true;
                self.agentsError = null;
                self.agents = [];
                self.selectedAgentIcode = '';
                self.selectedAgent = null;

                var params = new URLSearchParams({
                    team: team,
                });
                AGENT_ROLE_CODES.forEach(function(role) {
                    params.append('role', role);
                });

                fetch('/api/v1/equipo-operadores?' + params.toString())
                    .then(function(res) {
                        return res.json().then(function(data) {
                            return { ok: res.ok, data: data };
                        });
                    })
                    .then(function(result) {
                        if (!result.ok) {
                            self.agentsError = 'Error al cargar los agentes del equipo';
                            self.agents = [];
                            return;
                        }

                        var list = Array.isArray(result.data) ? result.data : [];
                        list.sort(function(a, b) {
                            return (a.fullName || '').localeCompare(b.fullName || '', 'es');
                        });

                        self.agents = list;
                        if (list.length === 0) {
                            self.agentsError = 'No hay agentes disponibles para el equipo "' + team + '"';
                        }
                    })
                    .catch(function() {
                        self.agentsError = 'Error al cargar los agentes del equipo';
                        self.agents = [];
                    })
                    .finally(function() {
                        self.loadingAgents = false;
                    });
            },

            // ------------------------------------------------------------
            // M-03 — Seleccionar agente
            // ------------------------------------------------------------

            onAgentChange: function() {
                var icode = this.selectedAgentIcode || '';
                if (!icode) {
                    this.selectedAgent = null;
                    return;
                }

                var agent = this.agents.find(function(a) {
                    return a.icode === icode;
                }) || null;

                if (!agent) {
                    this.selectedAgentIcode = '';
                    this.selectedAgent = null;
                    this.agentsError = 'El agente seleccionado no es válido';
                    return;
                }

                this.selectedAgent = agent;
                if (this.agents.length > 0) {
                    this.agentsError = null;
                }
            },

            // ------------------------------------------------------------
            // M-04 — Cerrar modal
            // ------------------------------------------------------------

            onBackdropClick: function() {
                if (this.confirmVisible || this.saving) {
                    return;
                }
                this.close();
            },

            close: function() {
                if (this.saving || this.confirmVisible) {
                    return;
                }

                this.visible = false;
                this.cases = [];
                this.resolvedTeam = '';
                this.agents = [];
                this.selectedAgentIcode = '';
                this.selectedAgent = null;
                this.loadingAgents = false;
                this.agentsError = null;
                this.saveError = null;
                this.saving = false;
                this.confirmVisible = false;

                this.$emit('closed', {});
            },

            cancelConfirm: function() {
                this.confirmVisible = false;
            },

            // ------------------------------------------------------------
            // M-05 — Confirmar reasignación (guardado)
            // ------------------------------------------------------------

            confirmReassign: function() {
                var self = this;
                if (!self.canConfirmReassign) {
                    return;
                }
                self.saveError = null;
                self.confirmVisible = true;
            },

            acceptConfirm: function() {
                var self = this;
                self.confirmVisible = false;
                self.runReassignSave();
            },

            runReassignSave: function() {
                var self = this;
                self.saveError = null;
                self.saving = true;

                var caseIcodes = self.cases.map(function(c) {
                    return c.i_code;
                }).filter(Boolean);

                fetch('/api/v1/cases/reasignar-bulk', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({
                        case_icodes: caseIcodes,
                        agent_icode: self.selectedAgentIcode,
                    }),
                })
                    .then(function(res) {
                        return res.json().then(function(data) {
                            return { ok: res.ok, data: data };
                        });
                    })
                    .then(function(result) {
                        if (!result.ok) {
                            self.saveError = (result.data && result.data.error)
                                || 'Error al reasignar los casos';
                            self.saving = false;
                            return;
                        }

                        var payload = {
                            cases: self.cases.slice(),
                            agent: self.selectedAgent,
                            updated: result.data.updated || 0,
                            skipped: result.data.skipped || 0,
                            followUpsUpdated: result.data.follow_ups_updated || 0,
                        };

                        self.saving = false;
                        self.$emit('reassigned', payload);
                        self.forceClose();
                    })
                    .catch(function() {
                        self.saveError = 'Error al reasignar los casos';
                        self.saving = false;
                    });
            },

            forceClose: function() {
                this.visible = false;
                this.cases = [];
                this.resolvedTeam = '';
                this.agents = [];
                this.selectedAgentIcode = '';
                this.selectedAgent = null;
                this.loadingAgents = false;
                this.agentsError = null;
                this.saveError = null;
                this.saving = false;
                this.confirmVisible = false;
            },

            // ------------------------------------------------------------
            // Helpers de presentación
            // ------------------------------------------------------------

            victimFullName: function(caseObj) {
                if (!caseObj) {
                    return '—';
                }
                return ((caseObj.names || '') + ' ' + (caseObj.lastNames || '')).trim() || '—';
            },

            ownerFullName: function(caseObj) {
                if (!caseObj || !caseObj.ownerNames) {
                    return 'Sin asignar';
                }
                return ((caseObj.ownerNames || '') + ' ' + (caseObj.ownerLastNames || '')).trim();
            },

            agentOptionLabel: function(agent) {
                if (!agent) {
                    return '';
                }
                var name = (agent.fullName || '').trim();
                if (REASSIGN_CROSS_TEAM_CONTINGENCY) {
                    var team = (agent.team || '').trim() || 'Sin equipo';
                    return name + ' — ' + team;
                }
                return name;
            },

            riskLabel: function(riskStatus) {
                var map = {
                    extremo: 'Extremo',
                    alto: 'Alto',
                    moderado: 'Moderado',
                    medio: 'Moderado',
                    bajo: 'Bajo',
                };
                if (!riskStatus) {
                    return '—';
                }
                return map[riskStatus.toLowerCase()] || 'Desconocido';
            },

            riskBadgeClass: function(riskStatus) {
                if (!riskStatus) {
                    return 'desconocido';
                }
                var slug = riskStatus.toLowerCase();
                if (slug === 'medio') {
                    return 'moderado';
                }
                var known = { extremo: true, alto: true, moderado: true, bajo: true };
                return known[slug] ? slug : 'desconocido';
            },
        },

        template: (typeof window !== 'undefined' && window.__reasignarCasosModalTpl)
            ? window.__reasignarCasosModalTpl
            : '#tpl-reasignar-casos-modal',
    });
})();
