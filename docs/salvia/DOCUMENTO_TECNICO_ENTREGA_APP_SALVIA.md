# Documento Técnico de Entrega
## Aplicación Móvil SALVIA

- **Entidad solicitante:** Ministerio de Igualdad y Equidad (Colombia)
- **Producto:** Aplicación móvil SALVIA
- **Versión del aplicativo:** `1.0.0+1`
- **Fecha del documento:** 8 de marzo de 2026
- **Equipo responsable:** _[Completar]_
- **Contacto técnico:** _[Completar]_

---

## 1. Tipo de documento
Este archivo corresponde a una **memoria técnica de entrega** (también conocida como documento técnico funcional y de implementación). Su objetivo es dejar trazabilidad de:

- Qué hace la aplicación.
- Cómo está construida.
- Qué dependencias tiene.
- Cómo se configura, prueba y despliega.
- Qué riesgos y consideraciones operativas deben tenerse en cuenta para producción.

---

## 2. Resumen ejecutivo
SALVIA es una aplicación móvil desarrollada en Flutter para el registro y acompañamiento de casos de violencias basadas en género. El flujo principal permite:

- Presentar información institucional y canales de atención.
- Capturar reportes mediante formulario estructurado.
- Validar interacción humana mediante código de verificación (imagen y audio).
- Enviar la información al backend institucional para su seguimiento.
- Ofrecer acceso rápido a líneas de emergencia (`123` y `155`).

---

## 3. Alcance funcional implementado

### 3.1 Pantalla de inicio
- Presentación de identidad visual SALVIA.
- Opción principal **Reportar**.
- Dos botones de emergencia:
  - **123:** Seguridad y emergencias (Policía, Bomberos, Ambulancias).
  - **155:** Orientación para mujeres víctimas de violencia.
- Mensajes de ayuda (snackbar) contextuales.
- Animaciones de entrada (separación de botones + `fade in` de contenido).

### 3.2 Pantalla “Conoce más acerca de SALVIA”
- Contenido informativo institucional.
- Textos y títulos con estilo visual unificado.
- Estructura responsive con scroll.

### 3.3 Pantalla de reporte
- Flujo de autorización de datos personales.
- Formulario dinámico condicionado por tipo de reporte.
- Validaciones de campos (obligatoriedad y formato).
- Validación de nombres/apellidos:
  - Solo letras y acentos.
  - Capitalización por palabra.
  - Máximo 32 caracteres por campo.
- Validación de teléfonos:
  - Solo dígitos.
  - Formato Colombia con prefijo visual `+57`.
- Campo de descripción de hechos:
  - Límite de 500 palabras.
  - Contador visible de palabras.
- Hora de contacto:
  - Entrada en 12h (AM/PM) y envío normalizado a formato 24h (`HH:mm`).

### 3.4 Captcha de imagen y audio
- Carga de configuración de formulario y `captchaID` desde backend.
- Descarga de imagen captcha (`/img`) con control de reintento.
- Reproducción de audio captcha (`/audio`).
- Refresco de captcha manual por usuario.
- Mensajería amigable: reemplazo de la palabra “CAPTCHA” por “CÓDIGO” en UI.

### 3.5 Flujo de envío
- Envío de payload JSON a backend.
- Pantalla de transición “Un momento” con spinner.
- Navegación a pantalla de confirmación “Reporte recibido”.
- Botón de retorno al inicio.

---

## 4. Arquitectura técnica

### 4.1 Cliente móvil
- **Framework:** Flutter
- **Lenguaje:** Dart
- **Nombre de app:** `Salvia`
- **Pantalla inicial:** `InicioScreen`

### 4.2 Organización general (resumen)
- `lib/screens/`: pantallas UI (`inicio`, `reportar`, `procesando_reporte`, `gracias`, `info`)
- `lib/salvia/`: integración HTTP y modelos de intercambio con backend (`server_proxy.dart`)
- `lib/widgets/`: componentes reutilizables de interfaz
- `lib/global/`: colores, exports, configuración visual

### 4.3 Backend consumido
- Base URL configurable en compilación:
  - `SALVIA_BASE_URL`
  - valor por defecto actual: `https://pruebassalvia.minigualdadyequidad.gov.co`

---

## 5. Endpoints utilizados

> Nota: las rutas se resuelven contra la base URL configurada.

- `GET /salvia/public/nuevo`
  - Obtiene configuración del formulario público (incluye `captchaID` y listas dinámicas).
- `GET /seguridad/login/captcha/{captchaID}/img`
  - Descarga imagen del código de verificación.
- `GET /seguridad/login/captcha/{captchaID}/audio`
  - Descarga audio del código de verificación.
- `POST /public/primer_contacto/nuevo`
  - Envía el registro del caso (payload del formulario).

---

## 6. Flujo de datos del formulario

1. App solicita configuración del formulario (`GET /salvia/public/nuevo`).
2. App actualiza listas dinámicas (sí/no, tipo de reporte, etc.).
3. App carga captcha asociado a ese `captchaID`.
4. Usuario diligencia formulario.
5. App normaliza y valida datos localmente.
6. App envía payload a `POST /public/primer_contacto/nuevo`.
7. Backend responde:
   - `200`: registro exitoso.
   - `4xx/5xx`: errores de validación o servidor.
8. App muestra retroalimentación y refresca captcha cuando aplica.

---

## 7. Estructura de payload (ejemplo)

```json
{
  "victimContact": {
    "form2": {
      "authorizationAnswer": {
        "icode": "...",
        "name": "Sí",
        "code": "y"
      },
      "reportType": {
        "icode": "...",
        "name": "Usted es la víctima de violencia basada en género",
        "code": "v"
      },
      "willReceiveCall": {
        "icode": "...",
        "name": "Sí",
        "code": "y"
      },
      "victimColPhone": 1234567890,
      "bestContactTime": "19:47",
      "factsDescription": "Descripción de hechos"
    },
    "names": "Nombre",
    "lastNames": "Apellido",
    "latitude": 4.82,
    "longitude": -74.35,
    "captchaID": "abc123",
    "captchaSolution": "457566"
  }
}
```

---

## 8. Dependencias principales
(archivo `pubspec.yaml`)

- `http ^0.13.6`
- `geolocator ^14.0.2`
- `url_launcher ^6.3.1`
- `audioplayers ^6.6.0`
- `path_provider ^2.1.5`
- `record ^4.4.4`
- `intl ^0.17.0`

Tipografía local:
- Roboto Variable (`lib/assets/fonts/roboto`)

---

## 9. Seguridad, privacidad y cumplimiento

### 9.1 Seguridad técnica
- Consumo por HTTPS.
- Validación de código de verificación (captcha) en backend.
- Manejo de cookies de sesión para flujo público.

### 9.2 Consideraciones de seguridad abiertas
- El backend puede aplicar reglas especiales por `User-Agent`; esta lógica debe revisarse para no debilitar controles de verificación.
- Se recomienda evaluar autenticación de cliente (token interno o mecanismo rotativo), sin romper el uso público del formulario web.

### 9.3 Datos personales
- El flujo incluye autorización explícita para tratamiento de datos.
- Se recomienda anexar en la entrega:
  - Política de tratamiento de datos vigente.
  - Matriz de datos recolectados y finalidad.

---

## 10. Accesibilidad

- Captcha de imagen y audio.
- Mensajes de error en lenguaje más claro (uso de “CÓDIGO”).
- Botones de acción con tamaño táctil amplio.
- Contraste y jerarquía visual reforzados.

> Recomendación: complementar con pruebas de usuario con lector de pantalla (TalkBack/VoiceOver) y ajuste de volumen de audio captcha desde servidor si fuera necesario.

---

## 11. Configuración y ejecución en desarrollo

### 11.1 Verificar entorno
```bash
flutter doctor
```

### 11.2 Ejecutar en Chrome (desarrollo)
```bash
flutter run -d chrome
```

### 11.3 Ejecutar en emulador Android
```bash
flutter emulators
flutter emulators --launch <emulator_id>
flutter devices
flutter run -d <device_id>
```

### 11.4 Compilar APK
```bash
flutter clean
flutter pub get
flutter build apk
```

Salida esperada:
- `build/app/outputs/flutter-apk/app-release.apk`

---

## 12. Incidencias conocidas y manejo

1. **CORS en Flutter Web**
- Para pruebas web contra backend con restricciones CORS, puede requerirse sesión de navegador de desarrollo con políticas relajadas.
- Esto es solo para ambiente de desarrollo.

2. **Diferencias de datos entre ambientes**
- Los `icode` de catálogos deben provenir del `GET` dinámico del backend.
- No depender de listas “quemadas” salvo fallback controlado para modo offline.

3. **Codificación de caracteres (tildes)**
- Si aparecen textos con caracteres corruptos (`SÃ`), revisar encoding UTF-8 en backend y cabeceras de respuesta.

4. **Compilación Android (Kotlin/Gradle)**
- Se debe mantener versión de Kotlin/Gradle compatible con Flutter usado en la estación de compilación.

---

## 13. Pruebas recomendadas para cierre

### 13.1 Funcionales
- Registro de caso exitoso con captcha correcto.
- Error de captcha incorrecto y refresco automático de imagen.
- Validaciones de nombre/apellido/teléfono/hora/contador de palabras.
- Flujo completo: Reportar → Procesando → Gracias → Inicio.

### 13.2 Usabilidad
- Pruebas en pantallas bajas y altas.
- Verificación de espaciados, legibilidad y coherencia visual.

### 13.3 Compatibilidad
- Android (emulador y dispositivo físico).
- Web (solo pruebas de desarrollo, sujeto a CORS).

---

## 14. Checklist de entrega

- [ ] Código fuente actualizado en repositorio Git.
- [ ] Historial de commits de cierre.
- [ ] APK de release generado y verificado.
- [ ] Documento técnico de entrega (este archivo).
- [ ] Evidencias de pruebas (capturas/logs).
- [ ] Matriz de endpoints y responsable backend.
- [ ] Acta de aceptación funcional (si aplica).

---

## 15. Recomendaciones de siguiente fase

- Formalizar contrato de API (OpenAPI/Swagger) para estabilizar integraciones.
- Definir estrategia de seguridad para cliente móvil sin romper formulario público.
- Implementar telemetría mínima (errores, tiempos de respuesta, métricas de uso).
- Priorizar pruebas de accesibilidad con población objetivo (incluyendo usuarias con discapacidad visual).

---

## 16. Control de cambios del documento

| Fecha | Versión | Autor | Cambio |
|---|---|---|---|
| 2026-03-08 | 1.0 | Equipo SALVIA | Versión inicial de entrega técnica |

