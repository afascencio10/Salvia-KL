function cerrarPopup() {
  document.getElementById('overlay').style.display = 'none';
}

function mostrarContenido(tab) {
  var tabs = document.querySelectorAll('.tabpopup');

  tabs.forEach(function(el) {
      el.classList.remove('active');
  });

  document.getElementById('info-tab').style.display = (tab === 'info') ? 'block' : 'none';
  document.getElementById('rutas-tab').style.display = (tab === 'rutas') ? 'block' : 'none';

  if (tab === 'info') {
      tabs[0].classList.add('active');
  } else if (tab === 'rutas') {
      tabs[1].classList.add('active');
  }
}