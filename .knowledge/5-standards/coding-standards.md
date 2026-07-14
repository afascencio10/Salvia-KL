---
okf_version: "1.0"
type: Coding_Standard
title: "Guía de Estándares de Codificación"
description: "Convenciones de código, naming, estructura de archivos, manejo de errores y patrones obligatorios para todo el equipo de desarrollo."
owner: "@platform-team"
status: active
tags: [standards, code, naming, patterns]
last_updated: "2026-06-23"
---

# Guía de Estándares de Codificación

Este documento define las convenciones que **todo código** del proyecto debe seguir, independientemente del stack específico.

## 1. Principios Generales

1. **Spec-First**: Nunca escribir código sin un documento OKF que lo respalde.
2. **Explícito sobre Implícito**: Preferir código verboso y claro sobre código "clever" y compacto.
3. **Fail Fast**: Validar inputs al inicio de cada función. Nunca dejar que datos inválidos se propaguen.
4. **Single Responsibility**: Cada función/módulo hace una sola cosa. Si necesitas `and` para describir lo que hace, debe dividirse.

## 2. Naming Conventions

### 2.1 Archivos y Directorios

| Elemento | Convención | Ejemplo |
| :--- | :--- | :--- |
| Directorios (feature) | `kebab-case` | `user-authentication/` |
| Archivos de código | `kebab-case` | `jwt-manager.ts`, `login-post.controller.ts` |
| Componentes UI | `PascalCase` | `LoginScreen.tsx`, `AuthForm.tsx` |
| Tests | `*.test.ts` o `*.spec.ts` | `jwt-manager.test.ts` |
| Documentación OKF | `kebab-case` | `login-post.md`, `users-table.md` |

### 2.2 Variables y Funciones

| Elemento | Convención | Ejemplo |
| :--- | :--- | :--- |
| Variables locales | `camelCase` | `accessToken`, `userEmail` |
| Constantes | `UPPER_SNAKE_CASE` | `MAX_LOGIN_ATTEMPTS`, `JWT_EXPIRY_SECONDS` |
| Funciones | `camelCase` (verbo + sustantivo) | `validateEmail()`, `generateToken()` |
| Clases / Interfaces | `PascalCase` | `UserService`, `AuthRequest` |
| Tipos / Enums | `PascalCase` | `UserRole`, `AuthStatus` |
| Booleans | Prefijo `is`, `has`, `can`, `should` | `isAuthenticated`, `hasPermission` |
| Event handlers | Prefijo `on` + evento | `onSubmitClick`, `onTokenExpired` |

### 2.3 Backend — Endpoints

| Elemento | Convención | Ejemplo |
| :--- | :--- | :--- |
| URLs | `kebab-case`, plural para colecciones | `/api/v1/auth/login`, `/api/v1/users` |
| Versionado | Prefijo `/api/vN/` | `/api/v1/`, `/api/v2/` |
| Query params | `snake_case` | `?page_size=20&sort_by=created_at` |
| Body fields | `snake_case` | `{ "access_token": "...", "expires_in": 3600 }` |

## 3. Estructura de Archivos de Código

### 3.1 Backend (Agnóstico)

La estructura depende del framework (Node puro, NestJS, LoopBack). Ejemplo genérico para feature `auth`:

```text
src/
├── features/ (o modules/)
│   └── auth/
│       ├── auth.controller.ts       ← Routing / Controllers
│       ├── auth.service.ts          ← Business logic
│       ├── auth.repository.ts       ← Acceso a datos
│       ├── auth.module.ts           ← (NestJS/Angular) Inyección de dependencias
│       ├── auth.validator.ts        ← Schemas de validación
│       └── __tests__/
├── shared/
│   ├── errors/                      ← Clases base de error
│   └── utils/                       ← Utilidades puras
└── config/
```

### 3.2 Frontend (Agnóstico)

El Frontend debe agruparse por dominio/feature, adaptándose a Angular, React o Vue.

```text
src/
├── features/
│   └── auth/
│       ├── screens/                   ← Contenedores inteligentes
│       │   └── LoginScreen/
│       │       ├── LoginScreen.tsx    ← (React) o .html/.ts/.css (Angular/Vue)
│       ├── components/                ← Componentes tontos
│       │   └── AuthForm/
│       ├── store/ o redux/            ← Manejo de estado
│       │   ├── actions.ts             ← (Redux) Actions / NgRx
│       │   ├── reducers.ts            ← (Redux) Reducers
│       │   ├── effects.ts             ← (Redux) Sagas / Effects
│       │   └── authStore.ts           ← (Context/Services) Store simple
│       ├── api/
│       │   └── authApi.ts             ← Llamadas a red
│       └── types/
└── shared/
    ├── components/
    └── store/                         ← Estado global (si aplica)
```

## 4. Manejo de Errores

### 4.1 Reglas Generales
- **Nunca** silenciar errores con `catch {}` vacío.
- **Siempre** tipar los errores con clases custom del dominio.
- **Siempre** loggear el error antes de re-lanzar o transformar.

### 4.2 Backend — Error Response Standard

Todos los errores API deben seguir este formato:

```json
{
  "error": {
    "code": "AUTH_INVALID_CREDENTIALS",
    "message": "El email o la contraseña son incorrectos.",
    "status": 401,
    "details": {},
    "request_id": "req_abc123"
  }
}
```

> [!WARNING]
> Los mensajes de error **nunca** deben revelar información interna del sistema (ej. "User not found in Database model users" es INACEPTABLE). Usar mensajes genéricos orientados al usuario.

### 4.3 Frontend — Error Handling Pattern

```text
try {
  // Operación
} catch (error) {
  if (error instanceof AuthError) {
    // Error conocido del dominio → mostrar al usuario
  } else if (error instanceof NetworkError) {
    // Error de red → retry o mensaje genérico
  } else {
    // Error inesperado → loggear + mensaje genérico
    logger.error('Unexpected error', { error, context });
  }
}
```

## 5. Imports — Orden Obligatorio

```text
// 1. Librerías externas (node_modules)
import express from 'express';
import { z } from 'zod';

// 2. Módulos internos compartidos (@shared/)
import { AppError } from '@shared/errors';
import { logger } from '@shared/utils';

// 3. Módulos del mismo feature
import { AuthService } from './auth.service';
import { LoginSchema } from './auth.validator';

// 4. Tipos (siempre al final)
import type { AuthRequest, AuthResponse } from './auth.types';
```

## 6. Comentarios

| Tipo | Cuándo usar | Formato |
| :--- | :--- | :--- |
| Docstring | Funciones públicas exportadas | JSDoc `/** */` |
| Inline | Lógica no obvia que necesita contexto de negocio | `// Razón: ...` |
| TODO | Trabajo pendiente con issue vinculado | `// TODO(#123): descripción` |
| WARNING | Código sensible que no debe modificarse sin entender el contexto | `// WARNING: ...` |

> [!CAUTION]
> Está **prohibido** comentar código muerto. Si no se usa, se elimina. Git tiene historial.
