function mostrarContenido(tab) {
    document.getElementById('info-tab').style.display = (tab === 'info') ? 'block' : 'none';
    document.getElementById('rutas-tab').style.display = (tab === 'rutas') ? 'block' : 'none';
}