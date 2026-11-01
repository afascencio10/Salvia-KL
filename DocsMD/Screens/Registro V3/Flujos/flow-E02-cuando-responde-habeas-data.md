━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando responde la autorización de datos personales
   Tipo: User Interaction
   Estado: ✅ Implementado — sincronizado con el código 2026-10-31
   Funciones: aceptarHabeasData() · rechazarHabeasData() → registro-v3-app.js
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  decisión:  botón "Acepto y continuar" o "No acepto"   → etapa 'habeas'
}

SI presiona "Acepto y continuar":
  → habeasAceptado = true
  → etapa = 'tipoContacto'
  → Vue oculta el bloque de Habeas data y muestra "Tipo de contacto"
  → CONTINÚA FLUJO GENERAL

SI presiona "No acepto":
  → etapa = 'formularioCerrado'
  → No se llama a ningún endpoint: no se crea `victim_case` ni se guarda ninguna respuesta
  → Vue muestra "La persona no autorizó el tratamiento de sus datos personales. No se crea ningún registro ni se guarda ninguna respuesta."
  → TERMINAR ejecución (no hay botón para volver: hay que recargar la página)

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión | Paso afectado |
|---|---|
| ¿Se debe registrar en algún log/analítica el número de personas que responden "No"? El Excel no lo pide; hoy no se registra nada | rama "No acepto" |
