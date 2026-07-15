---
okf_version: "1.0"
type: Security_Policy
title: "Clasificación de Datos (PII/PHI)"
description: "Política de clasificación de datos que define niveles de sensibilidad, controles de protección y requisitos de retención por tipo de dato."
owner: "@security-guild"
status: active
tags: [security, data-classification, pii, compliance, gdpr]
last_updated: "2026-06-23"
---

# Clasificación de Datos

## Niveles de Clasificación

| Nivel | Etiqueta | Descripción | Ejemplo |
| :---: | :--- | :--- | :--- |
| 🔴 4 | **Restricted** | Datos altamente sensibles. Compromiso = daño crítico | Passwords, API keys, claves de cifrado |
| 🟠 3 | **Confidential** | Datos personales identificables (PII) | Email, nombre, dirección IP |
| 🟡 2 | **Internal** | Datos internos de negocio, no públicos | Roles, configuración, métricas internas |
| 🟢 1 | **Public** | Datos que pueden ser públicos sin riesgo | Documentación pública, landing page |

## Inventario de Datos por Clasificación

### 🔴 Restricted (Nivel 4)

| Dato | Ubicación | Cifrado en Reposo | Cifrado en Tránsito | Acceso |
| :--- | :--- | :---: | :---: | :--- |
| `password_hash` | Base de Datos `users.password_hash` | AES-256 (disco) | TLS 1.2+ | Solo Auth Service |
| JWT Private Key | Secrets Manager / Vault | ✅ KMS | N/A (no se transmite) | Solo Auth Service |
| JWT Signing Key | Memoria del Auth Service | N/A | N/A | Runtime only |
| Database credentials | Environment variables | ✅ Vault/KMS | TLS 1.2+ | Solo infra team |

### 🟠 Confidential (Nivel 3 — PII)

| Dato | Ubicación | Cifrado en Reposo | Cifrado en Tránsito | Retención |
| :--- | :--- | :---: | :---: | :--- |
| `email` | Base de Datos `users.email` | AES-256 (disco) | TLS 1.2+ | Mientras cuenta activa + 30 días post-deletion |
| IP address (login) | Logs (Datadog) | ✅ | TLS 1.2+ | 90 días |
| User Agent | Logs (Datadog) | ✅ | TLS 1.2+ | 90 días |
| JWT claims (sub, email, role) | Cookie (HttpOnly) | N/A | TLS 1.2+ | 15 min (exp) |

### 🟡 Internal (Nivel 2)

| Dato | Ubicación | Cifrado | Retención |
| :--- | :--- | :---: | :--- |
| `role` (RBAC) | Base de Datos + JWT | Disco | Mientras cuenta activa |
| `created_at` | Base de Datos | Disco | Mientras cuenta activa |
| Feature flags config | Configuration service | No | Indefinido |
| Rate limiting counters | Redis | No | TTL automático (15 min) |

### 🟢 Public (Nivel 1)

| Dato | Ubicación | Retención |
| :--- | :--- | :--- |
| Documentación OKF | Git repository | Indefinido |
| API response schemas | Git repository | Indefinido |
| Error messages (user-facing) | API responses | N/A |

## Controles por Nivel

| Control | 🔴 Restricted | 🟠 Confidential | 🟡 Internal | 🟢 Public |
| :--- | :---: | :---: | :---: | :---: |
| Cifrado en reposo | ✅ Obligatorio | ✅ Obligatorio | ⚠️ Recomendado | — |
| Cifrado en tránsito (TLS) | ✅ Obligatorio | ✅ Obligatorio | ✅ Obligatorio | ✅ Obligatorio |
| Access logging | ✅ Obligatorio | ✅ Obligatorio | ⚠️ Recomendado | — |
| Backup cifrado | ✅ Obligatorio | ✅ Obligatorio | ✅ Obligatorio | — |
| Anonimización para staging | ✅ Obligatorio | ✅ Obligatorio | — | — |
| Retención definida | ✅ Obligatorio | ✅ Obligatorio | ⚠️ Recomendado | — |
| Derecho al olvido (GDPR) | ✅ | ✅ | — | — |

## Reglas para Developers

> [!CAUTION]
> - **Nunca** loggear datos de nivel 🔴 Restricted (passwords, keys).
> - **Nunca** exponer PII (🟠) en URLs, query params, o mensajes de error.
> - **Siempre** usar campos tipados del ORM para PII, no strings crudos.
> - **Siempre** anonimizar PII antes de copiar datos a staging/dev.
> - Al eliminar una cuenta, purgar **todos** los datos PII en 30 días (GDPR Art. 17).
