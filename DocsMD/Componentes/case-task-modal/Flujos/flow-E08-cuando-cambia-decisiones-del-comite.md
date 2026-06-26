━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando cambia decisiones del comité
   Tipo: User Interaction
   Función: onChangeDecisiones(valor)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Aplica a: formulario 'comite_caso' únicamente.

INPUT: {
  valor:              string con la opción marcada o desmarcada  → CheckboxGroup
                      'activar_enlace' | 'oficio' |
                      'recomendaciones_agente' | 'mecanismo_articulador'
  form.decisiones:    array actual de decisiones seleccionadas   → estado del componente
}


PASO 1 — Actualizar form.decisiones

  SI valor ya estaba en form.decisiones:
    → Remover: form.decisiones = form.decisiones.filter(d => d !== valor)
    → accion = 'desmarcado'

  SI valor no estaba en form.decisiones:
    → Agregar: form.decisiones.push(valor)
    → accion = 'marcado'


PASO 2 — Limpiar campos de las decisiones desmarcadas

  SI accion === 'desmarcado':

    SEGÚN valor:
      CASO 'oficio':
        → form.observacionesOficio = ""

      CASO 'recomendaciones_agente':
        → form.observacionesRecomendaciones = ""

      CASO 'mecanismo_articulador':
        → form.nivelMecanismo        = null
        → form.observacionesMecanismo = ""

      CASO 'activar_enlace':
        → (sin campos adicionales que limpiar)

  SI accion === 'marcado':
    → No hay campos que limpiar — los nuevos campos aparecen vacíos


PASO 3 — Vue re-evalúa reactivamente los bloques condicionales

  → [form.decisiones.includes('oficio')]
      Muestra/oculta el bloque con <textarea> Observaciones del oficio

  → [form.decisiones.includes('recomendaciones_agente')]
      Muestra/oculta el bloque con <textarea> Observaciones de recomendaciones

  → [form.decisiones.includes('mecanismo_articulador')]
      Muestra/oculta el bloque con <select> Nivel + <textarea> Observaciones del mecanismo

  → [form.decisiones.includes('activar_enlace')]
      No genera campos adicionales — solo se refleja en el payload al confirmar


PASO 4 — Vue re-evalúa formularioValido

  La computed property `formularioValido` incluye:
  SI mecanismo_articulador en decisiones: form.nivelMecanismo requerido
  Decisiones: min 1 opción requerida

  → BtnConfirmar se habilita si al menos 1 decisión está seleccionada
    y todos los campos requeridos de las decisiones activas están completos

→ FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
(sin gaps — lógica completamente definida)
