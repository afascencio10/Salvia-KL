/**
 * admin-users-list
 * Componente para listado de usuarios con búsqueda y paginación del servidor.
 * Usa el mismo patrón de casos-component: debounce search + server-side pagination.
 *
 * Props:
 *   - userBasePath (String): path base para las acciones (Ver, Actualizar)
 */

(function injectAdminUsersStyles() {
    if (document.getElementById('aul-styles')) return;
    var style = document.createElement('style');
    style.id = 'aul-styles';
    style.textContent = `
        .aul-filters { display:flex;gap:14px;align-items:flex-end;margin-bottom:16px;flex-wrap:wrap }
        .aul-filter-col { flex:1;min-width:180px }
        .aul-filter-label { font-size:.78rem;font-weight:600;color:#5106A7;margin-bottom:4px;display:block }
        .aul-filter-input { width:100%;padding:8px 14px;border:1px solid #d1d5db;border-radius:20px;font-size:.84rem;outline:none }
        .aul-filter-input:focus { border-color:#5106A7 }
        .aul-table { width:100%;border-collapse:collapse }
        .aul-table th { text-align:left;font-size:.75rem;font-weight:700;color:#6b7280;text-transform:uppercase;padding:10px 12px;border-bottom:2px solid #e5e7eb }
        .aul-table td { padding:10px 12px;font-size:.84rem;color:#374151;border-bottom:1px solid #f3f4f6 }
        .aul-table tr:hover { background:#f9fafb }
        .aul-badge { display:inline-block;padding:3px 10px;border-radius:12px;font-size:.72rem;font-weight:600 }
        .aul-badge-active { background:#dcfce7;color:#15803d }
        .aul-badge-disabled { background:#fef2f2;color:#dc2626 }
        .aul-actions { display:flex;gap:8px }
        .aul-actions a { font-size:.78rem;color:#5106A7;text-decoration:none;font-weight:500 }
        .aul-actions a:hover { text-decoration:underline }
        .aul-empty { text-align:center;padding:30px 0;color:#6b7280;font-size:.88rem }
        .aul-loading { text-align:center;padding:20px 0;color:#9ca3af }
        .aul-tabs { display:flex;gap:0;margin-bottom:14px }
        .aul-tab { padding:6px 14px;font-size:.8rem;font-weight:600;color:#6b7280;cursor:pointer;border-bottom:2px solid transparent }
        .aul-tab.active { color:#5106A7;border-bottom-color:#5106A7 }
        .aul-header { display:flex;justify-content:space-between;align-items:center;margin-bottom:10px }
        .aul-count { font-size:.78rem;color:#9ca3af;margin:0 }
        .aul-pagination { display:flex;align-items:center;justify-content:center;gap:4px;margin-top:16px }
        .aul-page-btn { min-width:32px;height:32px;display:flex;align-items:center;justify-content:center;border:1px solid #e5e7eb;border-radius:6px;background:#fff;font-size:.8rem;color:#374151;cursor:pointer;transition:all .15s }
        .aul-page-btn:hover:not(.active):not(.disabled) { border-color:#5106A7;color:#5106A7 }
        .aul-page-btn.active { background:#5106A7;color:#fff;border-color:#5106A7 }
        .aul-page-btn.disabled { opacity:.4;cursor:default }
        .aul-page-dots { font-size:.8rem;color:#9ca3af;padding:0 4px }
    `;
    document.head.appendChild(style);
})();

home.component('admin-users-list', {
    delimiters: ['${', '}'],
    props: {
        userBasePath: { type: String, default: '/seguridad/usuarios' },
    },
    data: function() {
        return {
            filtroNombre: '',
            filtroDocumento: '',
            filtroUsuario: '',
            usuarios: [],
            loading: false,
            activeTab: 'all',   // all | active | disabled
            currentPage: 1,
            pageSize: 20,
            totalUsers: 0,
            searchTimer: null,
            confirmarAccionUsuario: null,
        };
    },
    computed: {
        totalPages: function() {
            return Math.ceil(this.totalUsers / this.pageSize) || 1;
        },
        paginationPages: function() {
            var total = this.totalPages;
            var current = this.currentPage;
            var delta = 2;
            var range = [];
            var result = [];
            var l;

            for (var i = 1; i <= total; i++) {
                if (i === 1 || i === total || (i >= current - delta && i <= current + delta)) {
                    range.push(i);
                }
            }

            for (var j = 0; j < range.length; j++) {
                var page = range[j];
                if (l !== undefined) {
                    if (page - l === 2) {
                        result.push(l + 1);
                    } else if (page - l > 2) {
                        result.push('...');
                    }
                }
                result.push(page);
                l = page;
            }
            return result;
        },
    },
    mounted: function() {
        this.fetchUsers();
    },
    watch: {
        filtroNombre: function() { this.debounceFetch(); },
        filtroDocumento: function() { this.debounceFetch(); },
        filtroUsuario: function() { this.debounceFetch(); },
    },
    methods: {
        debounceFetch: function() {
            if (this.searchTimer) clearTimeout(this.searchTimer);
            var self = this;
            this.searchTimer = setTimeout(function() {
                self.currentPage = 1;
                self.fetchUsers();
            }, 400);
        },
        switchTab: function(tab) {
            if (this.activeTab === tab) return;
            this.activeTab = tab;
            this.currentPage = 1;
            this.fetchUsers();
        },
        changePage: function(page) {
            if (page < 1 || page > this.totalPages || page === this.currentPage) return;
            this.currentPage = page;
            this.fetchUsers();
        },
        fetchUsers: function() {
            var self = this;
            self.loading = true;

            var params = new URLSearchParams();
            params.set('page', String(self.currentPage));
            params.set('page_size', String(self.pageSize));

            // Status filter from tab
            if (self.activeTab === 'active') params.set('status', 'e');
            else if (self.activeTab === 'disabled') params.set('status', 'd');

            // Search filters
            if (self.filtroNombre.trim()) params.set('name', self.filtroNombre.trim());
            if (self.filtroDocumento.trim()) params.set('doc', self.filtroDocumento.trim());
            if (self.filtroUsuario.trim()) params.set('login', self.filtroUsuario.trim());

            fetch('/api/v1/admin/users/search?' + params.toString())
                .then(function(res) { return res.json(); })
                .then(function(data) {
                    self.usuarios = Array.isArray(data.users) ? data.users : [];
                    self.totalUsers = data.total || 0;
                })
                .catch(function() {
                    self.usuarios = [];
                    self.totalUsers = 0;
                })
                .finally(function() { self.loading = false; });
        },
        limpiar: function() {
            this.filtroNombre = '';
            this.filtroDocumento = '';
            this.filtroUsuario = '';
            this.currentPage = 1;
            this.fetchUsers();
        },
        descargarReporte: function() {
            // Mostrar loader
            var overlay = document.getElementById('loadingOverlay');
            if (overlay) overlay.style.display = 'flex';
            
            fetch('/api/v1/admin/users/report')
                .then(function(res) {
                    if (!res.ok) throw new Error('Error al generar reporte');
                    return res.blob();
                })
                .then(function(blob) {
                    var url = window.URL.createObjectURL(blob);
                    var a = document.createElement('a');
                    a.href = url;
                    a.download = 'usuarios-salvia_' + new Date().toISOString().substring(0, 10) + '.xlsx';
                    document.body.appendChild(a);
                    a.click();
                    a.remove();
                    window.URL.revokeObjectURL(url);
                })
                .catch(function(err) { alert(err.message || 'Error descargando reporte'); })
                .finally(function() {
                    if (overlay) overlay.style.display = 'none';
                });
        },
        inhabilitarUsuario: function(user) {
            this.confirmarAccionUsuario = { user: user, accion: 'inhabilitar' };
        },
        habilitarUsuario: function(user) {
            this.confirmarAccionUsuario = { user: user, accion: 'habilitar' };
        },
        ejecutarAccionUsuario: function() {
            var self = this;
            var user = this.confirmarAccionUsuario.user;
            var accion = this.confirmarAccionUsuario.accion;
            var ruta = accion === 'inhabilitar' ? '/deshabilitar_usuario' : '/habilitar_usuario';
            self.confirmarAccionUsuario = null;
            fetch(self.userBasePath + '/' + user.icode + ruta, {
                method: 'PUT',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({})
            }).then(function(res) {
                if (res.ok) {
                    self.fetchUsers();
                }
            }).catch(function() {});
        },
    },
    template: `
<div>
    <!-- Filtros -->
    <div class="aul-filters">
        <div class="aul-filter-col">
            <label class="aul-filter-label">Nombre / Apellido</label>
            <input type="text" class="aul-filter-input" v-model="filtroNombre" placeholder="Buscar por nombres..." />
        </div>
        <div class="aul-filter-col">
            <label class="aul-filter-label">Documento de identidad</label>
            <input type="text" class="aul-filter-input" v-model="filtroDocumento" placeholder="Buscar por documento..." />
        </div>
        <div class="aul-filter-col">
            <label class="aul-filter-label">Nombre de usuario</label>
            <input type="text" class="aul-filter-input" v-model="filtroUsuario" placeholder="Buscar por login..." />
        </div>
        <div>
            <button @click="limpiar()" style="background:#6b7280;color:#fff;border:none;padding:10px 16px;border-radius:20px;font-size:.82rem;font-weight:500;cursor:pointer;white-space:nowrap">Limpiar</button>
        </div>
    </div>

    <!-- Tabs -->
    <div class="aul-tabs">
        <div :class="['aul-tab', activeTab === 'all' ? 'active' : '']" @click="switchTab('all')">Todos</div>
        <div :class="['aul-tab', activeTab === 'active' ? 'active' : '']" @click="switchTab('active')">Activos</div>
        <div :class="['aul-tab', activeTab === 'disabled' ? 'active' : '']" @click="switchTab('disabled')">Desactivados</div>
    </div>

    <!-- Header with count -->
    <div class="aul-header">
        <p class="aul-count">\${ totalUsers } usuario(s) encontrado(s)</p>
        <button @click="descargarReporte()" style="background:#5106A7;color:#fff;border:none;padding:8px 16px;border-radius:8px;font-size:.8rem;font-weight:600;cursor:pointer;display:flex;align-items:center;gap:6px">
            <i class="fa fa-file-excel"></i> Descargar Excel
        </button>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="aul-loading"><i class="fa fa-spinner fa-spin"></i> Buscando...</div>

    <!-- Table -->
    <table v-else-if="usuarios.length > 0" class="aul-table">
        <thead>
            <tr>
                <th>Nombres</th>
                <th>Apellidos</th>
                <th>Tipo Doc</th>
                <th>Documento</th>
                <th>Rol</th>
                <th>Usuario</th>
                <th>Estado</th>
                <th>Acciones</th>
            </tr>
        </thead>
        <tbody>
            <tr v-for="user in usuarios" :key="user.icode">
                <td>\${ user.profile.names }</td>
                <td>\${ user.profile.lastNames }</td>
                <td>\${ user.profile.docType }</td>
                <td>\${ user.profile.docNumber }</td>
                <td>\${ user.roles && user.roles.length ? user.roles[0].name : '' }</td>
                <td>\${ user.login }</td>
                <td>
                    <span v-if="user.status === 'e'" class="aul-badge aul-badge-active">Activo</span>
                    <span v-else class="aul-badge aul-badge-disabled">Inhabilitado</span>
                </td>
                <td class="aul-actions">
                    <a :href="userBasePath + '/' + user.icode + '/'">Ver</a>
                    <a :href="userBasePath + '/' + user.icode + '/actualizar'">Actualizar</a>
                    <a v-if="user.status === 'e'" href="#" @click.prevent="inhabilitarUsuario(user)">Inhabilitar</a>
                    <a v-else href="#" @click.prevent="habilitarUsuario(user)">Habilitar</a>
                </td>
            </tr>
        </tbody>
    </table>

    <!-- Empty -->
    <div v-else class="aul-empty">No se encontraron usuarios con los criterios de búsqueda</div>

    <!-- Pagination -->
    <div v-if="!loading && totalPages > 1" class="aul-pagination">
        <div :class="['aul-page-btn', currentPage === 1 ? 'disabled' : '']" @click="changePage(currentPage - 1)">&laquo;</div>
        <template v-for="pg in paginationPages">
            <div v-if="pg === '...'" class="aul-page-dots">...</div>
            <div v-else :class="['aul-page-btn', pg === currentPage ? 'active' : '']" @click="changePage(pg)">\${ pg }</div>
        </template>
        <div :class="['aul-page-btn', currentPage === totalPages ? 'disabled' : '']" @click="changePage(currentPage + 1)">&raquo;</div>
    </div>

    <!-- Modal confirmación inhabilitar/habilitar -->
    <div v-if="confirmarAccionUsuario" style="position:fixed;inset:0;background:rgba(0,0,0,.45);display:flex;align-items:center;justify-content:center;z-index:10001;padding:16px" @click.self="confirmarAccionUsuario = null">
        <div style="background:#fff;border-radius:16px;width:100%;max-width:400px;padding:32px;box-shadow:0 20px 40px rgba(0,0,0,.18);text-align:center">
            <div style="width:52px;height:52px;border-radius:50%;display:flex;align-items:center;justify-content:center;margin:0 auto 16px" :style="confirmarAccionUsuario.accion === 'inhabilitar' ? 'background:#fef2f2' : 'background:#dcfce7'">
                <span style="font-size:1.5rem">\${ confirmarAccionUsuario.accion === 'inhabilitar' ? '⚠️' : '✓' }</span>
            </div>
            <p style="font-size:1.05rem;font-weight:700;color:#111827;margin:0 0 8px">\${ confirmarAccionUsuario.accion === 'inhabilitar' ? '¿Inhabilitar este usuario?' : '¿Habilitar este usuario?' }</p>
            <p style="font-size:.85rem;color:#6b7280;margin:0 0 6px">\${ confirmarAccionUsuario.user.profile.names } \${ confirmarAccionUsuario.user.profile.lastNames }</p>
            <p style="font-size:.78rem;color:#9ca3af;margin:0 0 24px">Login: \${ confirmarAccionUsuario.user.login }</p>
            <div style="display:flex;gap:12px;justify-content:center">
                <button @click="confirmarAccionUsuario = null" style="padding:10px 24px;border-radius:10px;border:1px solid #e5e7eb;background:#fff;color:#374151;font-size:.85rem;font-weight:600;cursor:pointer">Cancelar</button>
                <button @click="ejecutarAccionUsuario" style="padding:10px 24px;border-radius:10px;border:none;font-size:.85rem;font-weight:700;cursor:pointer" :style="confirmarAccionUsuario.accion === 'inhabilitar' ? 'background:#dc2626;color:#fff' : 'background:#16a34a;color:#fff'">
                    \${ confirmarAccionUsuario.accion === 'inhabilitar' ? 'Sí, inhabilitar' : 'Sí, habilitar' }
                </button>
            </div>
        </div>
    </div>
</div>
    `
});
