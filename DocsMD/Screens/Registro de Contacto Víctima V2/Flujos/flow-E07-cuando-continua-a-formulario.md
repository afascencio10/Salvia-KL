━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando presiona "Continuar" en Perfil
   Tipo: User Interaction
   Función: goTo('formulario') — botón del paso 'perfil'
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {}  (solo requiere tipoReporte y ajusteRazonable diligenciados, ya
            validado por :disabled — ver E-04/E-05)

PASO 1 — Navegar al paso siguiente
  this.step = 'formulario'   // goTo('formulario')

El árbol de bloques visibles en este paso queda determinado por el
form.tipoReporte fijado en E-04 (ver árbol de interfaz en
registro-contacto-victima-v2-interface.md). No hay llamada a backend.
