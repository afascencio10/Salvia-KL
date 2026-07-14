---
okf_version: "1.0"
type: Coding_Standard
title: "Estrategia Global de Testing"
description: "Pirámide de testing, thresholds de cobertura, herramientas estándar y política de testing obligatoria por tipo de cambio."
owner: "@qa-squad"
status: active
tags: [standards, testing, qa, coverage]
dependencies:
  - DoD: 5-standards/definition-of-done.md
last_updated: "2026-06-23"
---

# Estrategia Global de Testing

Este documento define **cómo, cuándo y cuánto** testeamos. Todo developer debe conocer esta estrategia antes de escribir su primer test.

## 1. Pirámide de Testing

```text
         ╱  E2E  ╲           ← 5%  | Flujos críticos de negocio Frontend UI (Playwright)
        ╱──────────╲
       ╱ Integración╲        ← 15% | APIs Backend y DB (Playwright APIRequestContext)
      ╱──────────────╲
     ╱   Contract      ╲     ← 10% | Frontend ↔ Backend contracts (Pact/MSW)
    ╱────────────────────╲
   ╱     Unitarios         ╲  ← 70% | Lógica pura, sin I/O (Jest/Vitest)
  ╱────────────────────────╲
```

> [!IMPORTANT]
> La pirámide **no es opcional**. Un PR que solo tiene tests E2E y cero unitarios será rechazado.

## 2. Thresholds de Cobertura (Coverage Gates)

| Métrica | Mínimo Global | Mínimo por PR (delta) | Herramienta |
| :--- | :---: | :---: | :--- |
| Line Coverage | 75% | No disminuir | Jest `--coverage` |
| Branch Coverage | 65% | No disminuir | Jest `--coverage` |
| Function Coverage | 80% | No disminuir | Jest `--coverage` |
| Mutation Score | 60% | No disminuir | Stryker (opcional Fase 4) |

> [!WARNING]
> Si un PR **reduce** la cobertura global, el pipeline debe fallar. Configurar en CI: `--coverageThreshold`.

## 3. Herramientas Estándar

| Capa | Herramienta | Propósito |
| :--- | :--- | :--- |
| Unit Tests | Jest / Vitest | Lógica pura, funciones, servicios |
| Component Tests | Testing Library | Renderizado de componentes UI |
| Integration Tests | Playwright API / Supertest | Endpoints HTTP end-to-end |
| Contract Tests | Pact / MSW | Contratos Frontend ↔ Backend |
| E2E Tests | Playwright | Flujos de usuario completos en navegador real |
| Load Tests | k6 / Artillery | Performance bajo carga |
| Security Tests | OWASP ZAP / Trivy | Vulnerabilidades conocidas |
| **Data Generation** | **@faker-js/faker** | Generación de mocks dinámicos para ahorrar tokens y mantener tests limpios |

## 4. Política de Testing por Tipo de Cambio

| Tipo de Cambio | Unit | Integration | Contract | E2E | Load |
| :--- | :---: | :---: | :---: | :---: | :---: |
| Nuevo endpoint API | ✅ | ✅ | ✅ | — | — |
| Nuevo componente UI | ✅ | — | ✅ | — | — |
| Lógica de negocio (service) | ✅ | — | — | — | — |
| Flujo crítico (login, pago) | ✅ | ✅ | ✅ | ✅ | ✅ |
| Migración de base de datos | — | ✅ | — | — | — |
| Fix de bug | ✅ (regresión) | — | — | — | — |
| Refactor sin cambio funcional | — (existentes deben pasar) | — | — | — | — |
| Cambio de configuración | — | ✅ | — | — | — |

## 5. Estructura de un Test

Todos los tests deben seguir el patrón **AAA (Arrange, Act, Assert)**:

```text
describe('AuthService', () => {
  describe('validateCredentials', () => {
    it('should return user when email and password are valid', () => {
      // Arrange — Preparar datos y dependencias
      const mockUser = createMockUser({ email: 'test@example.com' });
      const mockRepo = { findByEmail: jest.fn().mockResolvedValue(mockUser) };

      // Act — Ejecutar la operación bajo test
      const result = await service.validateCredentials('test@example.com', 'password');

      // Assert — Verificar el resultado
      expect(result).toEqual(mockUser);
      expect(mockRepo.findByEmail).toHaveBeenCalledWith('test@example.com');
    });
  });
});
```

## 6. Naming de Tests

| Patrón | Ejemplo |
| :--- | :--- |
| `should [acción esperada] when [condición]` | `should return 401 when password is invalid` |
| `should not [acción inesperada] when [condición]` | `should not allow login when account is locked` |
| `should throw [ErrorType] when [condición]` | `should throw RateLimitError when attempts exceed 5` |

> [!CAUTION]
> Están **prohibidos** los tests con nombres genéricos como: `it('works')`, `it('should work correctly')`, `it('test login')`.

## 7. Test Data Management

### 7.1 Factories (Recomendado)
Crear factories para cada entidad del dominio:

```text
// test/factories/user.factory.ts
export function createMockUser(overrides?: Partial<User>): User {
  return {
    id: randomUUID(),
    email: 'default@test.com',
    password_hash: '$2b$10$...',
    created_at: new Date(),
    ...overrides,
  };
}
```

### 7.2 Reglas de Datos de Test
- **Nunca** usar datos de producción en tests.
- **Nunca** hardcodear IDs o timestamps. Usar factories o `faker`.
- **Siempre** limpiar el estado entre tests (`beforeEach` / `afterEach`).
- **Nunca** que un test dependa del resultado de otro test.

## 8. Integración con CI/CD

El repositorio es agnóstico del pipeline, pero ofrece soporte nativo para **GitHub Actions** y **GitLab CI**.

```text
Pipeline Stages:
  validate  →  build  →  test:unit  →  test:integration (API)  →  test:e2e (UI)  →  deploy
                               │               │                        │
                          coverage gate   DB migrations             staging env
                           (≥75% lines)    (Playwright API)         (Playwright E2E)
```

Si cualquier stage de test falla, el pipeline **se detiene**. No se despliega código con tests rojos.

### Ejecución de Playwright en CI/CD
- **GitHub Actions**: Configurado en `.github/workflows/playwright.yml`.
- **GitLab CI**: Job predefinido `e2e_playwright` en `.gitlab-ci.yml`.
Los tests de UI deben configurarse para subir un artefacto de HTML Trace si fallan, para depuración rápida.
