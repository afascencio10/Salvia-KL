---
okf_version: "1.0"
type: Coding_Standard
title: "Security Checklist (OWASP-Based)"
description: "Lista de verificación de seguridad basada en OWASP Top 10, adaptada al contexto del proyecto con verificaciones obligatorias."
owner: "@security-guild"
status: active
tags: [standards, security, owasp, compliance]
dependencies:
  - RBAC: 6-security/rbac-matrix.md
last_updated: "2026-06-23"
---

# Security Checklist (OWASP-Based)

Todo feature que maneja datos de usuario, autenticación o autorización **debe** pasar esta checklist antes de llegar a producción.

## A01: Broken Access Control

- [ ] Todo endpoint verifica permisos en el middleware (`verifyRole()`), no solo en el frontend.
- [ ] Los recursos se acceden por ID del usuario autenticado, no por parámetro manipulable.
- [ ] Las rutas administrativas requieren rol `SUPER_ADMIN` o `COMPANY_OWNER` según [RBAC Matrix](/.knowledge/6-security/rbac-matrix.md).
- [ ] No existe IDOR (Insecure Direct Object Reference): `GET /users/123` verifica que 123 pertenece al tenant del usuario.
- [ ] CORS está configurado con dominios específicos (no `*` en producción).

## A02: Cryptographic Failures

- [ ] Contraseñas hasheadas con Bcrypt (cost ≥ 10) o Argon2.
- [ ] Tokens JWT firmados con algoritmo asimétrico (RS256), no simétrico (HS256).
- [ ] Datos PII (email, nombre) encriptados en tránsito (TLS 1.2+) y en reposo (AES-256).
- [ ] No se loggean passwords, tokens, o datos sensibles.
- [ ] Claves de cifrado rotadas periódicamente y almacenadas en Vault/KMS (no en `.env`).

## A03: Injection

- [ ] Todas las queries a DB usan **parameterized queries** (no string concatenation).
- [ ] Los inputs del usuario están validados con schema (Zod/Joi) antes de llegar al service.
- [ ] No se usa `eval()`, `new Function()`, ni template literals dinámicos en queries.
- [ ] Los headers HTTP del usuario no se insertan directamente en respuestas (XSS via headers).

## A04: Insecure Design

- [ ] Existe un threat model documentado para features críticos.
- [ ] Rate limiting implementado en endpoints de autenticación (≤5 intentos / 15 min).
- [ ] Account lockout después de N intentos fallidos.
- [ ] No se revelan si un email existe en la base de datos en mensajes de error.

## A05: Security Misconfiguration

- [ ] Headers de seguridad configurados: `X-Content-Type-Options`, `X-Frame-Options`, `Strict-Transport-Security`.
- [ ] Modo debug desactivado en producción.
- [ ] Endpoints de monitoreo (`/health`, `/metrics`) no expuestos públicamente.
- [ ] Error stack traces no se devuelven al cliente en producción.
- [ ] Dependencias escaneadas con `npm audit` o `Trivy` sin vulnerabilidades críticas.

## A07: Identification and Authentication Failures

- [ ] Tokens JWT tienen `exp` corto (≤15 min para access token, ≤7 días para refresh).
- [ ] Refresh tokens almacenados en `HttpOnly`, `Secure`, `SameSite=Strict` cookies.
- [ ] No se almacenan tokens en `localStorage` (vulnerable a XSS).
- [ ] La sesión se invalida completamente al hacer logout (blacklist del token).
- [ ] Multi-factor authentication disponible para roles administrativos.

## A08: Software and Data Integrity Failures

- [ ] Las dependencias se instalan con lockfile (`package-lock.json` / `yarn.lock`).
- [ ] CI/CD pipeline no ejecuta código de PRs de forks sin aprobación.
- [ ] Las imágenes Docker se construyen desde base images oficiales y versionadas (no `latest`).

## A09: Security Logging and Monitoring Failures

- [ ] Intentos de login fallidos se loggean con IP, timestamp, y email (no password).
- [ ] Cambios de permisos/roles se loggean como eventos de auditoría.
- [ ] Alertas configuradas para patrones anómalos (>100 intentos fallidos/min desde misma IP).
- [ ] Los logs no contienen datos PII en texto plano.

## Aplicación por Feature

| Feature | Secciones Obligatorias |
| :--- | :--- |
| Autenticación (login, register) | A01, A02, A03, A04, A05, A07, A09 |
| API pública | A01, A03, A05 |
| Pagos / Billing | A01, A02, A03, A04, A05, A07, A08, A09 |
| Dashboard / UI interna | A01, A05 |
| Reportes / Exportaciones | A01, A02 |

> [!CAUTION]
> Todo feature que involucra **dinero** (billing, payments, subscriptions) requiere la checklist **completa** sin excepciones.
