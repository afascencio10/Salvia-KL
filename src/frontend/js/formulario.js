  function nextStep(step) {
    document.getElementById(`step-container${step}`).classList.remove('active');
    document.getElementById(`step${step}`).classList.remove('badge-primary');
    document.getElementById(`step${step}`).classList.add('badge-secondary');
    step++;
    document.getElementById(`step-container${step}`).classList.add('active');
    document.getElementById(`step${step}`).classList.remove('badge-secondary');
    document.getElementById(`step${step}`).classList.add('badge-primary');
  }

  function prevStep(step) {
    document.getElementById(`step-container${step}`).classList.remove('active');
    document.getElementById(`step${step}`).classList.remove('badge-primary');
    document.getElementById(`step${step}`).classList.add('badge-secondary');
    step--;
    document.getElementById(`step-container${step}`).classList.add('active');
    document.getElementById(`step${step}`).classList.remove('badge-secondary');
    document.getElementById(`step${step}`).classList.add('badge-primary');
  }

  function submitForm() {
    // Aquí puedes agregar el código para enviar el formulario
    alert('Los datos de usuario y contraseña han sido enviados al número de celular de confianza');
  }

  document.getElementById('departamento').addEventListener('change', function() {
    var departamento = this.value;
    var municipioSelect = document.getElementById('municipio');
    municipioSelect.innerHTML = '';

    // Agrega los municipios correspondientes a cada departamento
    switch(departamento) {
      case 'Amazonas':
        addOption(municipioSelect, 'Leticia', 'Leticia');
        addOption(municipioSelect, 'Puerto Nariño', 'Puerto Nariño');
        break;
      case 'Antioquia':
        addOption(municipioSelect, 'Medellín', 'Medellín');
        addOption(municipioSelect, 'Bello', 'Bello');
        break;
      // Agrega más casos según tus necesidades
      default:
        addOption(municipioSelect, 'Seleccionar municipio', '');
    }
  });

  function addOption(select, text, value) {
    var option = document.createElement('option');
    option.text = text;
    option.value = value;
    select.add(option);
  }
	function abrirMapa() {
    var direccion = document.getElementById('direccion').value;
    var url = `https://www.google.com/maps/search/?api=1&query=${encodeURIComponent(direccion)}`;
    window.open(url, '_blank');
  }
	
	function guardar() {
    // Aquí puedes agregar el código para guardar la información
    document.getElementById('btnAsignar').disabled = false;
  }

  function asignar() {
    nextStep(1); // Avanzar al Paso 2
  }
  function editar() {
    prevStep(2); // Regresar al Paso 1
  }

  function finalizar() {
    nextStep(2); // Avanzar al Paso 3
    generarUsuarioYContraseña();
  }

  function generarUsuarioYContraseña() {
    var nombre = document.getElementById('nombre').value;
    var apellido = document.getElementById('apellido').value;
    var tipoDocumento = document.getElementById('tipoDocumento').value;
    var usuario = nombre.toLowerCase() + '.' + apellido.toLowerCase() + tipoDocumento.toLowerCase();
    var contraseña = Math.random().toString(36).slice(2); // Generar contraseña aleatoria

    var mensaje = `Usuario: ${usuario}\nContraseña: ${contraseña}`;
    document.getElementById('mensaje').innerText = mensaje;
  }// JavaScript Document