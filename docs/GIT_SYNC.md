# Sincronización de repos — GitLab y GitHub

Este proyecto mantiene dos remotos:
- `gitlab` → repositorio del equipo (fuente de verdad compartida)
- `origin` → GitHub personal (historial limpio bajo tu usuario)

---

## Traer cambios del equipo desde GitLab

```bash
git fetch gitlab develop
git merge gitlab/develop --no-edit
```

Si hay conflictos en archivos MD, quedarte con tu versión:
```bash
git checkout --ours <archivo-en-conflicto>
git add <archivo-en-conflicto>
git commit --no-edit
```

---

## Subir tus cambios a GitLab

```bash
git add <archivos>
git commit -m "feat: descripción del cambio"
git push gitlab develop
```

Si GitLab rechaza el push (historias divergentes), primero haz el merge de arriba y luego vuelve a pushear.

---

## Subir cambios a GitHub (squash — solo tus commits)

Primero asegúrate de tener todo commiteado localmente. Luego:

```bash
# 1. Traer estado actual de GitHub
git fetch origin develop

# 2. Squash merge de los cambios nuevos
git merge --squash origin/develop   # si GitHub está adelante
# o simplemente push si tu rama ya tiene todo
git push origin develop
```

Si quieres que todos los commits del equipo aparezcan bajo tu usuario en GitHub, usa squash al integrar desde GitLab antes de pushear a origin:

```bash
git fetch gitlab develop
git merge --squash gitlab/develop
git commit -m "sync: descripción del cambio"
git push origin develop
```

---

## Archivos ignorados (no commitear)

| Archivo / carpeta | Razón |
|---|---|
| `src/bitsflow` | Binario compilado de Go (154 MB) |
| `DocsMD/Otros/temp/` | Archivos temporales de análisis |
| `DocsMD/.obsidian/` | Estado local del workspace de Obsidian |
