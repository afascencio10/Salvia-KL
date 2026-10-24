━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando abre el modal de revisión
   Tipo: User Interaction
   Código: E-11
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Disparado por: botón "Revisar (Equipo Técnico)"
Lo ven sm y sv (vista supervisor) y em (vista equipo-tecnico),
con estado pendiente_de_revision.
Permiso review_remisiones_medidas: sm, em y sv.


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

PASO 1 — Abrir vacío

  revisionItem        = la fila
  revisionProfesional = ""
  revisionConcepto    = ""
  revisionResultado   = ""
  revisionAbierta     = true

  Título: "Revisión Equipo Técnico"
  Texto: "Caso {nombre} · {caseICode}"

PASO 2 — Catálogo

  GET /api/v1/emergency-measure-referrals/profesionales-revision

  {
    "items": [
      { "id": "general_user_i_code", "nombre": "Nombres Apellidos" }
    ]
  }

  El catálogo es role_code = em, usuario activo (M-01, paso 5).
  El select guarda el id. La fila, al recargar, muestra el nombre.

PASO 3 — Resultado

  Tres botones, un valor cada uno (M-01):

    aprobado_para_caracterizacion_me
    devuelto_para_ajuste
    devuelto_para_acompanamiento

  "Guardar revisión" se habilita con profesional, concepto con texto
  después de trim, y un resultado.

  Este evento no escribe la remisión.

  → FIN EJECUCIÓN ✓
