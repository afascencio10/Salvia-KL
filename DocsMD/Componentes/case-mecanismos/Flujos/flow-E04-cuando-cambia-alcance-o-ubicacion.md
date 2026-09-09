━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando cambia alcance o ubicación
   Tipo: User Interaction
   Función: onCambiaAlcance / onCambiaDepartamento / onCambiaCiudad / onCambiaMunicipio
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Al cambiar el alcance se resetean ubicación y usuario.

Campos visibles:
  municipal     → Departamento + Ciudad + Municipio + select usuarios mc
  departamental → Departamento + Ciudad opcional + select usuarios mc
  nacional      → select usuarios mc (sin geo)

Cuando el alcance (y la ubicación, si aplica) está completo:

GET /api/v1/mecanismos/usuarios?scope=&department_id=&city_id=&town_id=

→ usuarios = mc activos cuyo profile.town → city → department coincide
→ si vacío: "No hay usuarios mecanismo en esta ubicación"
