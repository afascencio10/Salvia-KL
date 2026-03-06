var mas = document.querySelector('.icono-mas');
if(mas!=null){
    
   mas.addEventListener('click', function() {
        
        var tablaDesplegable = document.querySelector('#tabla-desplegable');
        tablaDesplegable.classList.toggle('mostrar');
    
        var icono = document.querySelector('.icono-mas i');
        if (tablaDesplegable.classList.contains('mostrar')) {
            tablaDesplegable.style.maxHeight = '500px'; // Establecer max-height al valor deseado
            icono.classList.remove('fa-plus');
            icono.classList.add('fa-minus');
        } else {
            tablaDesplegable.style.maxHeight = '0';
            icono.classList.remove('fa-minus');
            icono.classList.add('fa-plus');
        }
    });
}


