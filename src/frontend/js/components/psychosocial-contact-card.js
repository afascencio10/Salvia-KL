/**
 * psychosocial-contact-card
 * Card "Intentos de contacto 3x3": contador (X de N), tiles de último/próximo
 * intento, historial colapsable con badges (NO CONTESTÓ / SÍ CONTESTÓ / PROGRAMADO)
 * y botón "Registrar intento" que dispara el flujo <psychosocial-contact-modal>.
 * Autosuficiente: carga su historial vía GET. Se monta en una o varias pantallas.
 */
(function injectPsychosocialContactCardStyles() {
    if (document.getElementById('psc-card-styles')) return;
    var style = document.createElement('style');
    style.id = 'psc-card-styles';
    style.textContent = `
        .psc-card { background:#fff; border:1px solid #eef0f4; border-radius:1rem; padding:1.5rem; box-shadow:0 1px 3px rgba(0,0,0,.06); max-width:760px; }
        .psc-card-head { display:flex; justify-content:space-between; align-items:flex-start; gap:1rem; flex-wrap:wrap; }
        .psc-card-title { font-size:1.15rem; font-weight:800; color:#111827; margin:0; }
        .psc-card-count { font-size:.85rem; color:#6b7280; margin:.15rem 0 0; }
        .psc-card-help { font-size:.72rem; color:#9ca3af; margin:.35rem 0 0; }
        .psc-register-btn { background:#7c3aed; color:#fff; border:none; border-radius:.65rem; padding:.7rem 1.1rem; font-weight:700; font-size:.85rem; cursor:pointer; white-space:nowrap; }
        .psc-register-btn:disabled { opacity:.6; cursor:not-allowed; }
        .psc-tiles { display:grid; grid-template-columns:1fr 1fr; gap:.75rem; margin:1.1rem 0; }
        .psc-tile { background:#f5f3ff; border:1px solid #ede9fe; border-radius:.75rem; padding:.8rem .9rem; }
        .psc-tile-label { font-size:.72rem; font-weight:700; color:#6d28d9; display:flex; align-items:center; gap:.4rem; }
        .psc-tile-value { font-size:.9rem; color:#111827; font-weight:600; margin-top:.25rem; }
        .psc-hist { border:1px solid #eef0f4; border-radius:.75rem; overflow:hidden; }
        .psc-hist-head { width:100%; text-align:left; background:#f9fafb; border:none; padding:.85rem 1rem; font-size:.78rem; font-weight:700; letter-spacing:.03em; color:#4b5563; cursor:pointer; display:flex; justify-content:space-between; align-items:center; }
        .psc-hist-body { padding:.5rem 1rem 1rem; }
        .psc-attempt { border:1px solid #eef0f4; border-radius:.6rem; padding:.7rem .85rem; margin-top:.6rem; display:flex; justify-content:space-between; gap:.75rem; }
        .psc-attempt-title { font-weight:700; color:#111827; font-size:.9rem; }
        .psc-attempt-meta { font-size:.78rem; color:#6b7280; margin-top:.15rem; }
        .psc-attempt-note { font-size:.82rem; color:#374151; margin-top:.3rem; }
        .psc-badge { align-self:flex-start; font-size:.68rem; font-weight:800; padding:.25rem .55rem; border-radius:9999px; white-space:nowrap; letter-spacing:.02em; }
        .psc-badge-no { background:#fee2e2; color:#b91c1c; }
        .psc-badge-yes { background:#dcfce7; color:#15803d; }
        .psc-badge-prog { background:#ede9fe; color:#6d28d9; }
        .psc-empty { font-size:.85rem; color:#9ca3af; padding:.6rem 0; }
    `;
    document.head.appendChild(style);
})();

app.component('psychosocial-contact-card', {
    delimiters: ['${', '}'],
    props: {
        psicosocialId:   { type: String, required: true },
        currentUser:     { type: String, default: '' },
        currentUserId:   { type: String, default: '' },
        userTeam:        { type: String, default: '' },
        closureThreshold:{ type: Number, default: 9 }
    },
    data: function() {
        return {
            loading: true,
            error: null,
            historyCollapsed: false,
            data: null
        };
    },
    computed: {
        attempts: function() { return (this.data && this.data.attempts) || []; },
        counters: function() { return (this.data && this.data.counters) || {}; },
        caseName: function() { return (this.data && this.data.caseName) || ''; }
    },
    mounted: function() { this.fetchHistory(); },
    methods: {
        fetchHistory: function() {
            var self = this;
            this.loading = true; this.error = null;
            fetch('/api/v1/psychosocial/' + encodeURIComponent(this.psicosocialId) + '/contact-attempts')
                .then(function(res) {
                    return res.json().then(function(json) {
                        self.loading = false;
                        if (!res.ok) { self.error = json.error || 'No se pudo cargar el historial'; return; }
                        self.data = json;
                    });
                })
                .catch(function(e) { self.loading = false; self.error = e.message; });
        },
        openRegister: function() {
            var c = this.counters;
            this.$refs.modal.start({
                id: this.psicosocialId,
                caseName: this.caseName,
                dailyFailedCount: c.dailyFailedCount || 0,
                totalCount: c.totalCount || 0,
                distinctDaysCount: c.distinctDaysCount || 0,
                closureEligible: !!c.closureEligible
            });
        },
        onModalCompleted: function() { this.fetchHistory(); },
        fmt: function(v) {
            if (!v) return '—';
            var d = new Date(v);
            if (isNaN(d.getTime())) return '—';
            return d.toLocaleString('es-CO', { day: '2-digit', month: 'short', year: 'numeric', hour: '2-digit', minute: '2-digit' });
        },
        noteText: function(a) {
            if (a.note && a.note.trim()) return a.note;
            return a.wasAnswered ? 'Sí contestó' : 'Sin nota';
        }
    },
    template: `
    <div class="psc-card">
      <div v-if="loading" class="psc-empty"><i class="fas fa-circle-notch fa-spin"></i> Cargando…</div>
      <div v-else-if="error" class="psc-empty">
        <i class="fas fa-exclamation-triangle" style="color:#dc2626;"></i> \${ error }
        <button class="psc-register-btn" style="margin-left:.75rem;" @click="fetchHistory">Reintentar</button>
      </div>
      <template v-else>
        <div class="psc-card-head">
          <div>
            <h3 class="psc-card-title">Intentos de contacto <span style="color:#7c3aed;">3x3</span></h3>
            <p class="psc-card-count">\${ counters.totalCount || 0 } de \${ closureThreshold } intentos</p>
            <p class="psc-card-help">Al llegar a \${ closureThreshold } intentos en ≥3 días distintos se habilita el cierre por imposibilidad de contacto.</p>
          </div>
          <button class="psc-register-btn" :disabled="counters.maxAttemptsReached" @click="openRegister">
            <i class="fas fa-phone-alt"></i> Registrar intento
          </button>
        </div>

        <div class="psc-tiles">
          <div class="psc-tile">
            <div class="psc-tile-label"><i class="far fa-clock"></i> Último intento realizado</div>
            <div class="psc-tile-value">\${ fmt(data.lastAttemptAt) }</div>
          </div>
          <div class="psc-tile">
            <div class="psc-tile-label"><i class="far fa-calendar-alt"></i> Próximo intento</div>
            <div class="psc-tile-value">\${ fmt(data.nextContactAttemptAt) }</div>
          </div>
        </div>

        <div class="psc-hist">
          <button class="psc-hist-head" :aria-expanded="(!historyCollapsed).toString()" @click="historyCollapsed = !historyCollapsed">
            <span><i class="far fa-clock"></i> HISTORIAL DE CONTACTO</span>
            <i :class="historyCollapsed ? 'fas fa-chevron-down' : 'fas fa-chevron-up'"></i>
          </button>
          <div v-show="!historyCollapsed" class="psc-hist-body">
            <div v-if="attempts.length === 0 && !data.nextContactAttemptAt" class="psc-empty">Aún no hay intentos registrados.</div>
            <!-- Fila del próximo intento programado: siempre arriba (es un evento futuro) -->
            <div v-if="data.nextContactAttemptAt" class="psc-attempt">
              <div>
                <div class="psc-attempt-title">Próximo</div>
                <div class="psc-attempt-meta"><i class="far fa-clock"></i> \${ fmt(data.nextContactAttemptAt) }</div>
              </div>
              <span class="psc-badge psc-badge-prog">PROGRAMADO</span>
            </div>
            <!-- Intentos: más nuevo → más viejo (el API ya los entrega en ese orden) -->
            <div v-for="a in attempts" :key="a.id" class="psc-attempt">
              <div>
                <div class="psc-attempt-title">Intento #\${ a.sequenceNumber }</div>
                <div class="psc-attempt-meta"><i class="far fa-clock"></i> \${ fmt(a.attemptAt) }</div>
                <div class="psc-attempt-note">\${ noteText(a) }</div>
              </div>
              <span class="psc-badge" :class="a.wasAnswered ? 'psc-badge-yes' : 'psc-badge-no'">\${ a.wasAnswered ? 'SÍ CONTESTÓ' : 'NO CONTESTÓ' }</span>
            </div>
          </div>
        </div>

        <psychosocial-contact-modal
          ref="modal"
          :current-user="currentUser"
          :current-user-id="currentUserId"
          :user-team="userTeam"
          @completed="onModalCompleted">
        </psychosocial-contact-modal>
      </template>
    </div>
    `
});
