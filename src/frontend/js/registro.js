function siguientePaso(paso) {
    document.querySelector(`.paso${paso}`).classList.remove('activo');
    document.querySelector(`.paso${paso + 1}`).classList.add('activo');
    document.querySelector(`.paso${paso}`).classList.add('oculto');
    document.querySelector(`.paso${paso + 1}`).classList.remove('oculto');
}

function anteriorPaso(paso) {
    document.querySelector(`.paso${paso}`).classList.remove('activo');
    document.querySelector(`.paso${paso - 1}`).classList.add('activo');
    document.querySelector(`.paso${paso}`).classList.add('oculto');
    document.querySelector(`.paso${paso - 1}`).classList.remove('oculto');
}