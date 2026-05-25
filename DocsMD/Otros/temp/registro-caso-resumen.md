# Registro de Caso

El operador completa un formulario con datos de la victima, hechos de violencia, informacion del agresor y un tamizaje de riesgo (preguntas Si/No que calculan un puntaje y nivel: bajo, moderado, alto o extremo).

Al guardar, el backend valida los campos, recalcula el riesgo para prevenir manipulacion del cliente, crea un usuario con credenciales aleatorias para la victima, inserta el caso con su formulario y asigna la ruta institucional (sedes por momento/sector).

Despues del commit, una goroutine genera automaticamente el calendario de seguimientos segun el nivel de riesgo y asigna un agente de forma balanceada.
