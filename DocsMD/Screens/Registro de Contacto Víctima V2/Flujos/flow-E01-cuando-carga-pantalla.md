━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando carga la pantalla
   Tipo: Lifecycle
   Estado: ✅ Implementado y probado — 2026-09-27
   Función: VictimContactPOST_GET_Public_V2() (Go, renderiza el HTML) →
            mounted() del app Vue de set_victim_contact_v2.html (cliente)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {}  (sin parámetros de entrada — GET simple, sin :id ni querystring)

PASO 1 — El backend genera un captchaID nuevo (decisión 4: se mantiene captcha)
  GET /salvia/public-v2/nuevo → VictimContactPOST_GET_Public_V2()
  captcha.New() → captchaID
  → se inyecta directo en el HTML servido (plantilla Go: {{.CaptchaID}}),
    NO es una llamada AJAX aparte desde el cliente — el propio GET que carga
    la pantalla ya lo trae embebido.

PASO 2 — Inicializar el wizard (cliente, mounted())
  step.value = 'auth'
  form.value = blankForm(captchaID)   // captchaID del PASO 1; el resto en '' o null

  El widget de captcha (imagen + input de solución) se muestra recién en el
  último paso del wizard ('formulario', justo antes de "Enviar Reporte" —
  ver árbol de interfaz), pero el ID ya está disponible desde la carga
  inicial de la página.
