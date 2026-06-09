(function () {
    function getSidebar() {
        return document.getElementById('salviaSidebar');
    }

    function getBackdrop() {
        return document.getElementById('salviaSidebarBackdrop');
    }

    function getToggleBtn() {
        return document.getElementById('salviaSidebarToggle');
    }

    function isCollapsed() {
        return document.body.classList.contains('salvia-sidebar-collapsed');
    }

    function setCollapsed(collapsed) {
        var sidebar = getSidebar();
        var backdrop = getBackdrop();
        var toggleBtn = getToggleBtn();

        document.body.classList.toggle('salvia-sidebar-collapsed', collapsed);

        if (sidebar) {
            sidebar.setAttribute('aria-hidden', collapsed ? 'true' : 'false');
        }
        if (backdrop) {
            backdrop.hidden = collapsed;
            backdrop.setAttribute('aria-hidden', collapsed ? 'true' : 'false');
        }
        if (toggleBtn) {
            toggleBtn.hidden = !collapsed;
            toggleBtn.setAttribute('aria-expanded', collapsed ? 'false' : 'true');
            toggleBtn.setAttribute('aria-label', collapsed ? 'Abrir menú de navegación' : 'Cerrar menú de navegación');
        }
    }

    function openSidebar() {
        setCollapsed(false);
    }

    function closeSidebar() {
        setCollapsed(true);
    }

    function toggleSidebar() {
        setCollapsed(!isCollapsed());
    }

    function markActiveLinks() {
        var path = window.location.pathname.replace(/\/+$/, '') || '/';
        document.querySelectorAll('.salvia-sidebar__nav-link').forEach(function (link) {
            var href = (link.getAttribute('href') || '').replace(/\/+$/, '') || '/';
            if (href === path) {
                link.classList.add('is-active');
                link.setAttribute('aria-current', 'page');
            }
        });
    }

    function bindEvents() {
        document.addEventListener('click', function (e) {
            if (e.target.closest('#salviaSidebarToggle')) {
                e.preventDefault();
                toggleSidebar();
                return;
            }
            if (e.target.closest('#salviaSidebarClose')) {
                e.preventDefault();
                closeSidebar();
                return;
            }
            if (e.target.id === 'salviaSidebarBackdrop') {
                closeSidebar();
            }
        });

        document.addEventListener('keydown', function (e) {
            if (e.key === 'Escape' && !isCollapsed()) {
                closeSidebar();
                var toggleBtn = getToggleBtn();
                if (toggleBtn) toggleBtn.focus();
            }
        });
    }

    function updateSidebarOffset() {
        var gov = document.querySelector('.header-gov');
        var banner = document.querySelector('.nav_banner');
        if (!gov || !banner) return;

        var offset = gov.offsetHeight + banner.offsetHeight;
        document.documentElement.style.setProperty('--salvia-header-offset', offset + 'px');
    }

    function relocateSidebarElements() {
        var rail = document.getElementById('salviaSidebarRail');
        var backdrop = getBackdrop();
        if (backdrop) {
            document.body.appendChild(backdrop);
        }
        if (rail) {
            document.body.appendChild(rail);
        }
    }

    function initSalviaSidebar() {
        if (!getSidebar()) return;

        relocateSidebarElements();
        updateSidebarOffset();
        window.addEventListener('resize', updateSidebarOffset);
        window.addEventListener('load', function () {
            relocateSidebarElements();
            updateSidebarOffset();
        });
        bindEvents();
        markActiveLinks();
        setCollapsed(false);
    }

    window.initSalviaSidebar = initSalviaSidebar;

    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', initSalviaSidebar);
    } else {
        initSalviaSidebar();
    }
})();
