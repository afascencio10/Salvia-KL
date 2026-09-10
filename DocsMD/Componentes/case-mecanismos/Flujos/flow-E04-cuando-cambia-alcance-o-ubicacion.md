━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando cambia tipo o ubicación
   Tipo: User Interaction
   Función: onCambiaAlcance / onCambiaDepartamento / onCambiaCiudad / onCambiaMunicipio
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Al cambiar el tipo de mecanismo se resetean ubicación y usuario.

Campos visibles:
  municipal     → Departamento + Ciudad + Municipio + select usuarios mc
  departamental → Departamento + select usuarios mc
  nacional      → select usuarios mc (sin geo)

Cuando el tipo (y la ubicación, si aplica) está completo:

GET /api/v1/mecanismos/usuarios?scope=&department_id=&city_id=&town_id=

→ usuarios = mc activos cuyo articulating_mechanism.type y ubicación coinciden
→ si vacío: "No hay usuarios mecanismo en esta ubicación"
