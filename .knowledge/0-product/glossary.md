---
okf_version: "1.0"
type: Glossary
title: "Glosario de Términos del Dominio"
description: "Definiciones estándar de términos técnicos y de negocio usados en la documentación y código del proyecto."
owner: "@platform-team"
status: active
tags: [documentation, glossary, onboarding]
last_updated: "2026-06-23"
---

# Glosario de Términos

Referencia canónica de terminología. Cuando un término aparece en specs o código, su significado es **exactamente** el definido aquí.

## Autenticación & Seguridad

| Término | Definición |
| :--- | :--- |
| **JWT** | JSON Web Token. Token firmado digitalmente que contiene claims del usuario. Se usa como credencial stateless. Ver el ADR de autenticación en `1-architecture/adrs/`. |
| **Access Token** | JWT de corta duración (15 min) usado para autenticar requests a la API. |
| **Refresh Token** | Token de larga duración (7 días) usado para obtener nuevos access tokens sin re-login. |
| **JTI** | JWT ID. Claim estándar que identifica unívocamente un token. Usado para revocación. |
| **Token Blacklist** | Lista (Redis) de JTIs de tokens que fueron revocados antes de expirar. |
| **HttpOnly Cookie** | Cookie que no es accesible via JavaScript (`document.cookie`). Previene robo via XSS. |
| **SameSite** | Atributo de cookie que previene CSRF. Valor: `Strict`. |
| **Bcrypt** | Algoritmo de hashing adaptativo para contraseñas. Incluye salt automático. |
| **Cost Factor** | Parámetro de Bcrypt que controla la complejidad del hash (2^N iteraciones). |
| **Rate Limiting** | Mecanismo que limita la cantidad de requests por IP/usuario en una ventana de tiempo. |
| **RBAC** | Role-Based Access Control. Sistema de permisos basado en roles predefinidos. |
| **MFA** | Multi-Factor Authentication. Segundo factor de autenticación (TOTP, SMS). |

## Arquitectura

| Término | Definición |
| :--- | :--- |
| **Tenant** | Instancia lógica de un negocio/organización dentro de la plataforma multi-tenant. |
| **Microservicio** | Servicio independiente con su propia responsabilidad, desplegable de forma aislada. |
| **API Gateway** | Punto de entrada único que enruta requests a los microservicios internos. |
| **SPA** | Single Page Application. Aplicación web que se carga una vez y actualiza dinámicamente. |
| **C4 Model** | Modelo de documentación arquitectónica con 4 niveles: Context, Container, Component, Code. |
| **ADR** | Architecture Decision Record. Documento que registra una decisión técnica y su justificación. |

## Datos

| Término | Definición |
| :--- | :--- |
| **PII** | Personally Identifiable Information. Datos que pueden identificar a una persona (email, nombre). |
| **PHI** | Protected Health Information. Datos de salud protegidos por regulaciones (HIPAA). |
| **Migration** | Script versionado que modifica el esquema de la base de datos de forma controlada. |
| **Optimistic Locking** | Estrategia de concurrencia donde se detectan conflictos al momento de escribir (no al leer). |

## Frontend

| Término | Definición |
| :--- | :--- |
| **Redux** | Librería de state management. Ver el ADR de state management en `1-architecture/adrs/`. |
| **Store** | Contenedor de estado global gestionado por Redux. Un store por app o feature. |
| **Selector** | Función que extrae una porción específica del store para un componente. |
| **Screen** | Componente de nivel página que corresponde a una ruta del router. |
| **Component** | Pieza reutilizable de UI que recibe props y emite eventos. |

## Operaciones

| Término | Definición |
| :--- | :--- |
| **SLI** | Service Level Indicator. Métrica medible (ej. latencia P99 del login endpoint). |
| **SLO** | Service Level Objective. Target numérico de un SLI (ej. P99 < 500ms). |
| **SLA** | Service Level Agreement. Compromiso formal con el cliente sobre disponibilidad. |
| **Runbook** | Guía paso-a-paso para diagnosticar y resolver un incidente específico. |
| **Blue/Green** | Estrategia de despliegue con dos entornos idénticos que se intercambian. |
| **Canary** | Estrategia de despliegue que envía un % del tráfico a la nueva versión antes del rollout completo. |
| **Feature Flag** | Toggle que permite activar/desactivar funcionalidades sin nuevo despliegue. |

## OKF (Open Knowledge Format)

| Término | Definición |
| :--- | :--- |
| **OKF** | Open Knowledge Format. Estándar de documentación basado en YAML frontmatter + Markdown. |
| **Frontmatter** | Bloque YAML al inicio de un archivo `.md` delimitado por `---`. |
| **Spec** | Documento OKF que describe un componente/endpoint/servicio **antes** de implementarlo. |
| **Traceability** | Capacidad de rastrear un requisito desde el PRD hasta el test que lo verifica. |
