---
okf_version: "1.0"
type: Threat_Model
title: "Modelo de Amenazas: Autenticación (STRIDE)"
description: "Análisis de amenazas basado en el framework STRIDE aplicado al módulo de autenticación con mitigaciones implementadas."
owner: "@security-guild"
status: active
tags: [security, threat-model, stride, auth]
dependencies:
  - RBAC: 6-security/rbac-matrix.md
  - Checklist: 5-standards/security-checklist.md
last_updated: "2026-06-23"
---

# Threat Model: Autenticación (STRIDE)

## Alcance

Este modelo cubre el flujo completo de autenticación: login, verificación de token, logout y rate limiting.

## Diagrama de Flujo de Datos (DFD)

```text
                    ┌──────────────┐
                    │   Internet   │
                    └──────┬───────┘
                           │ HTTPS
                    ┌──────▼───────┐
          ┌─────── │  Cloudflare   │ ──────┐
          │        │  WAF / CDN    │       │
          │        └──────┬───────┘       │
          │               │               │
   ┌──────▼───────┐      │        ┌──────▼───────┐
   │  SPA Frontend │      │        │  Attacker    │
   └──────┬───────┘      │        └──────────────┘
          │ JWT Cookie    │
   ┌──────▼───────┐      │
   │  API Gateway  │◄─────┘
   │  (Rate Limit) │
   └──────┬───────┘
          │
   ┌──────▼───────┐     ┌──────────────┐
   │  Auth Service │────►│  Base de Datos  │
   │  (JWT, Bcrypt)│     │  (users)      │
   └──────┬───────┘     └──────────────┘
          │
   ┌──────▼───────┐
   │    Redis      │
   │  (blacklist)  │
   └──────────────┘
```

## Análisis STRIDE

### S — Spoofing (Suplantación de Identidad)

| # | Amenaza | Probabilidad | Impacto | Mitigación | Estado |
| :--- | :--- | :---: | :---: | :--- | :---: |
| S1 | Robo de credenciales via phishing | Alta | Alto | MFA para roles admin (V3), educación al usuario | 🔜 Parcial |
| S2 | Token JWT robado via XSS | Media | Alto | HttpOnly + Secure + SameSite cookies | ✅ |
| S3 | Replay attack con token capturado | Media | Alto | JWT con `exp` corto (15min), JTI único | ✅ |
| S4 | Creación de tokens falsos | Baja | Crítico | Firma asimétrica RS256 (clave privada protegida) | ✅ |

### T — Tampering (Manipulación)

| # | Amenaza | Probabilidad | Impacto | Mitigación | Estado |
| :--- | :--- | :---: | :---: | :--- | :---: |
| T1 | Modificar claims del JWT | Baja | Crítico | Firma RS256 — cualquier modificación invalida el token | ✅ |
| T2 | SQL Injection en login | Media | Crítico | Parameterized queries, validación con Zod/JSON Schema | ✅ |
| T3 | Manipular cookie de sesión | Baja | Alto | HttpOnly impide acceso via JS, Secure impide HTTP plano | ✅ |

### R — Repudiation (Repudio)

| # | Amenaza | Probabilidad | Impacto | Mitigación | Estado |
| :--- | :--- | :---: | :---: | :--- | :---: |
| R1 | Usuario niega haber realizado acciones | Media | Medio | Audit log con IP, timestamp, user_id para cada acción crítica | ✅ |
| R2 | Falta de trazabilidad en cambios de roles | Media | Alto | Log de eventos RBAC con before/after | ✅ |

### I — Information Disclosure (Fuga de Información)

| # | Amenaza | Probabilidad | Impacto | Mitigación | Estado |
| :--- | :--- | :---: | :---: | :--- | :---: |
| I1 | Enumeración de usuarios via error messages | Alta | Medio | Mensaje genérico "Credenciales inválidas" para email y password | ✅ |
| I2 | Leak de passwords en logs | Baja | Crítico | Nunca loggear passwords; solo email e IP en intentos fallidos | ✅ |
| I3 | JWT payload con datos sensibles | Baja | Medio | JWT solo contiene id, email, role. No incluye password ni PII extra | ✅ |
| I4 | Stack traces en responses de error | Media | Medio | En producción: mensaje genérico. Stack trace solo en logs internos | ✅ |

### D — Denial of Service (Denegación de Servicio)

| # | Amenaza | Probabilidad | Impacto | Mitigación | Estado |
| :--- | :--- | :---: | :---: | :--- | :---: |
| D1 | Brute force attack al login | Alta | Alto | Rate limiting: 5 intentos/15min por IP (Redis) | ✅ |
| D2 | DDoS al endpoint de login | Media | Crítico | Cloudflare WAF + Rate Limiting en edge | ✅ |
| D3 | Bcrypt DoS (passwords de 1MB) | Baja | Medio | maxLength: 128 caracteres en password validation | ✅ |
| D4 | Account lockout abuse | Media | Medio | Lockout por IP, no por cuenta (previene lock de cuentas ajenas) | ✅ |

### E — Elevation of Privilege (Escalación de Privilegios)

| # | Amenaza | Probabilidad | Impacto | Mitigación | Estado |
| :--- | :--- | :---: | :---: | :--- | :---: |
| E1 | IDOR: acceder a recursos de otro usuario | Media | Crítico | Verificar ownership en cada endpoint, no confiar en frontend | ✅ |
| E2 | Algorithm confusion (HS256 vs RS256) | Baja | Crítico | Rechazar tokens con `alg ≠ RS256` antes de verificar | ✅ |
| E3 | Manipular role en JWT | Baja | Crítico | Role se asigna server-side, nunca se acepta del cliente | ✅ |
| E4 | CSRF en endpoints sensibles | Media | Alto | SameSite=Strict + CSRF token para operaciones de escritura | ✅ |

## Resumen de Riesgo Residual

| Categoría | Amenazas | Mitigadas | Pendientes | Riesgo Residual |
| :--- | :---: | :---: | :---: | :---: |
| Spoofing | 4 | 3 | 1 (MFA) | 🟡 Medio |
| Tampering | 3 | 3 | 0 | 🟢 Bajo |
| Repudiation | 2 | 2 | 0 | 🟢 Bajo |
| Info Disclosure | 4 | 4 | 0 | 🟢 Bajo |
| Denial of Service | 4 | 4 | 0 | 🟢 Bajo |
| Elevation of Privilege | 4 | 4 | 0 | 🟢 Bajo |

> [!IMPORTANT]
> El único riesgo residual significativo es **S1 (phishing)**, que se mitiga con MFA planificado para V3. Mientras tanto, se recomienda educación al usuario y alertas por login desde nueva IP/dispositivo.
