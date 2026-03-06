 function mostrarContenido(tab) {
        document.querySelectorAll('.tabpopup').forEach(function(el) {
            el.classList.remove('active');
        });

        document.getElementById('info-tab').style.display = (tab === 'info') ? 'block' : 'none';
        document.getElementById('rutas-tab').style.display = (tab === 'rutas') ? 'block' : 'none';

        if (tab === 'info') {
            document.querySelector('.tabpopup.active').classList.remove('active');
            document.querySelector('.tabpopup:nth-child(1)').classList.add('active');
        } else if (tab === 'rutas') {
            document.querySelector('.tabpopup.active').classList.remove('active');
            document.querySelector('.tabpopup:nth-child(2)').classList.add('active');
        }
    }// JavaScript Document