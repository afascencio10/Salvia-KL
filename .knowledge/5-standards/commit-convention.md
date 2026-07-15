---
okf_version: "1.0"
type: Coding_Standard
title: "Convención de Commits"
description: "Estándar de mensajes de commit basado en Conventional Commits v1.0 con scopes obligatorios para el proyecto."
owner: "@platform-team"
status: active
tags: [standards, git, commits, changelog]
last_updated: "2026-06-23"
---

# Convención de Commits

Adoptamos [Conventional Commits v1.0](https://www.conventionalcommits.org/) para habilitar changelogs automáticos, versionado semántico, y navegación efectiva del historial.

## 1. Formato

```text
<type>(<scope>): <description>

[optional body]

[optional footer(s)]
```

### Ejemplo Real
```text
feat(auth): add rate limiting to login endpoint

Implement sliding window rate limiter using Redis.
Blocks IP after 5 failed attempts in 15 minutes.

Refs: #42
Spec: .knowledge/3-features/Auth/backend/endpoints/login-post.md
```

## 2. Types Permitidos

| Type | SemVer | Cuándo usar |
| :--- | :---: | :--- |
| `feat` | MINOR | Nueva funcionalidad visible al usuario |
| `fix` | PATCH | Corrección de un bug |
| `docs` | — | Cambios solo en documentación (`.knowledge/`, README) |
| `test` | — | Agregar o corregir tests (sin cambio funcional) |
| `refactor` | — | Cambio de código sin cambio funcional ni fix |
| `perf` | PATCH | Mejora de performance |
| `style` | — | Formato, espacios, puntos y comas (sin cambio lógico) |
| `build` | — | Cambios en build system, CI/CD, dependencias |
| `chore` | — | Tareas de mantenimiento que no son código de producción |
| `ci` | — | Cambios en configuración de CI pipelines |
| `revert` | — | Reversar un commit anterior |

## 3. Scopes Obligatorios

El scope **es obligatorio** y debe coincidir con el feature o capa afectada:

| Scope | Dominio |
| :--- | :--- |
| `auth` | Feature de autenticación |
| `billing` | Feature de pagos |
| `users` | Feature de gestión de usuarios |
| `db` | Migraciones y esquemas de base de datos |
| `api` | Cambios transversales de API (middleware, routing) |
| `ui` | Cambios transversales de UI (design system, layouts) |
| `infra` | Infraestructura, Docker, Terraform |
| `deps` | Actualización de dependencias |
| `okf` | Cambios en documentación OKF |

> [!IMPORTANT]
> Si un cambio afecta múltiples scopes, hacer commits atómicos por scope. Un commit gigante con scope `*` será rechazado.

## 4. Reglas

### 4.1 Descripción (`<description>`)
- Máximo **72 caracteres**.
- Empezar con **verbo en imperativo y minúscula**: `add`, `fix`, `remove`, `update` (no `added`, `fixes`).
- Sin punto final.
- En **inglés** (el código es en inglés, los commits también).

### 4.2 Body (Opcional pero recomendado)
- Separado de la descripción por una línea en blanco.
- Explica el **por qué**, no el **qué** (el diff ya muestra el qué).
- Wrap a 80 caracteres.

### 4.3 Footer (Opcional)
- `Refs: #<issue>` — Issue relacionado.
- `Spec: <path>` — Documento OKF que especifica el cambio.
- `BREAKING CHANGE: <description>` — Cambio que rompe compatibilidad (dispara SemVer MAJOR).
- `Co-authored-by: Name <email>` — Co-autoría.

## 5. Breaking Changes

```text
feat(auth)!: replace session cookies with JWT tokens

BREAKING CHANGE: All existing session cookies are invalidated.
Clients must re-authenticate to obtain a JWT.
Migration guide: .knowledge/1-architecture/adrs/adr-001-jwt-vs-sessions.md
```

> [!WARNING]
> Todo commit con `BREAKING CHANGE` debe tener un ADR asociado que justifique la decisión.

## 6. Squash & Merge

- Los PRs se mergean con **Squash and Merge**.
- El título del squash commit **debe** seguir el formato conventional commits.
- Los commits individuales dentro del PR pueden ser más informales (WIP, fixup, etc.).

## 7. Validación en CI

Configurar `commitlint` con `@commitlint/config-conventional`:

```text
// commitlint.config.js
module.exports = {
  extends: ['@commitlint/config-conventional'],
  rules: {
    'scope-empty': [2, 'never'],          // Scope obligatorio
    'scope-enum': [2, 'always', [         // Scopes válidos
      'auth', 'billing', 'users', 'db',
      'api', 'ui', 'infra', 'deps', 'okf'
    ]],
    'subject-max-length': [2, 'always', 72],
  },
};
```

## 8. Ejemplos Correctos vs Incorrectos

| ✅ Correcto | ❌ Incorrecto |
| :--- | :--- |
| `feat(auth): add password reset endpoint` | `added password reset` |
| `fix(ui): prevent double submit on login form` | `Fix bug` |
| `docs(okf): add JWT manager service spec` | `update docs` |
| `test(auth): add unit tests for rate limiter` | `tests` |
| `refactor(auth): extract token validation to middleware` | `refactor stuff` |
