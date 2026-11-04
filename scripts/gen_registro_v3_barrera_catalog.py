#!/usr/bin/env python3
"""Genera BARRERA_V4_V3 y ENTIDAD_V4_V3: textos de Seguimiento v4 para Registro V3.

REG-18 (WP-4): el formulario de barrera de Registro V3 es el mismo de Seguimiento v4
(15 preguntas para el agente, 27 filas con los textos informativos y los campos por
sector). Este script LEE el seed de Seguimiento v4 y escribe el catálogo en
src/frontend/js/components/registro-v3-campos.js, entre los marcadores
`// >>> BARRERA_V4_V3` y `// <<< BARRERA_V4_V3`. Los textos, las opciones y las
condiciones salen del seed: nada se transcribe a mano.

Además genera ENTIDAD_V4_V3 (los textos de «Identificación de entidades» de Seguimiento v4:
radicado, fecha, acciones, acuerdo, debida diligencia, canales y gestiones) y CANAL_ACTIVACION_V3
(las 4 opciones de canal), entre los marcadores `// >>> ENTIDAD_V4_V3` y `// <<< ENTIDAD_V4_V3`.

Uso (desde la raíz del repositorio):
    python3 scripts/gen_registro_v3_barrera_catalog.py            # reescribe el bloque
    python3 scripts/gen_registro_v3_barrera_catalog.py --check    # falla si el bloque está desactualizado

Fuente: src/internal/seed/seguimiento_v4_form.sql, grupo repetidor de la barrera de
«Identificación de entidades» (52387e43-9d64-412d-83b5-b2b6cbbf02c8). Si el seed cambia
(se regenera con scripts/gen_seguimiento_v4_seed.js), vuelve a ejecutar este script.
"""
import json
import os
import re
import sys

ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
SEED = os.path.join(ROOT, "src", "internal", "seed", "seguimiento_v4_form.sql")
TARGET = os.path.join(ROOT, "src", "frontend", "js", "components", "registro-v3-campos.js")
BARRIER_REPEATER = "52387e43-9d64-412d-83b5-b2b6cbbf02c8"
BEGIN = "// >>> BARRERA_V4_V3"
END = "// <<< BARRERA_V4_V3"
BEGIN_E = "// >>> ENTIDAD_V4_V3"
END_E = "// <<< ENTIDAD_V4_V3"
# Grupo repetidor de la entidad de «Identificación de entidades» y sección (pregunta ident.conocimiento).
ENTITY_REPEATER = "2bda4735-3927-4634-a666-c7f97949be0a"
ENTITY_SECTION = "e742200a-4bc9-4241-b552-31e4666b72b1"

# Claves estables por posición (columna "order" 1..27 del grupo repetidor). Son los
# mismos sufijos que usa scripts/gen_seguimiento_v4_seed.js (barrierQuestions); el tipo
# esperado de cada posición se verifica contra el seed para detectar cambios de estructura.
KEYS = [
    ("sector", "single"), ("infoSalud", "info"), ("tiposSalud", "multiple"), ("instSalud", "dropdown"),
    ("otraSalud", "text"), ("infoJusticia", "info"), ("tiposJusticia", "multiple"),
    ("instJusticia", "dropdown"), ("otraJusticia", "text"), ("infoProteccion", "info"),
    ("tiposProteccion", "multiple"), ("instProteccion", "dropdown"), ("otraProteccion", "text"),
    ("talento", "multiple"), ("eco", "multiple"), ("terr", "multiple"), ("diferencial", "multiple"),
    ("fecha", "date"), ("funcionario", "text"), ("depto", "dropdown"), ("municipio", "dropdown"),
    ("vereda", "dropdown"), ("descripcion", "text"), ("gestion", "multiple"), ("perspectiva", "text"),
    ("infoAlerta", "info"), ("gestionTexto", "text"),
]
# Preguntas del repetidor de entidad por posición ("order" 1..16); mismos sufijos que ident.* del generador del seed.
ENTIDAD_KEYS = [
    ("sector", "dropdown"), ("depto", "dropdown"), ("municipio", "dropdown"), ("ciudad", "dropdown"),
    ("entidad", "dropdown"), ("otraEntidad", "text"), ("radicado", "text"), ("fecha", "date"),
    ("completadas", "multiple"), ("pendientes", "multiple"), ("otraRuta", "text"), ("nuevaBarrera", "boolean"),
    ("acuerdo", "boolean"), ("debida", "boolean"), ("canales", "multiple"), ("gestiones", "text"),
]
# Preguntas cuyas opciones no están en el seed sino en la geografía del sistema.
FUENTE = {"depto": "departamentos", "municipio": "municipios"}
# Etiquetas que deben seguir existiendo tal cual (guarda contra un seed que cambió de forma).
ETIQUETAS_ESPERADAS = {
    "sector": "Sector en el que se presentó la barrera",
    "fecha": "Fecha en la que se presentó la barrera (aproximada)",
    "gestion": "Gestión de la barrera",
    "perspectiva": "¿Cuál es su perspectiva profesional frente a la barrera?",
    "gestionTexto": "Descripción de la gestión realizada con relación a las barreras",
}


def parse_rows(text):
    """Filas `( 'a', NULL, true, 12 )` de un bloque VALUES -> listas de Python."""
    rows, i, n = [], 0, len(text)
    while i < n:
        if text[i] != "(":
            i += 1
            continue
        i += 1
        row = []
        while True:
            while text[i] in " \n\t":
                i += 1
            if text[i] == "'":
                j, buf = i + 1, []
                while True:
                    if text[j] == "'":
                        if j + 1 < n and text[j + 1] == "'":
                            buf.append("'")
                            j += 2
                            continue
                        break
                    buf.append(text[j])
                    j += 1
                row.append("".join(buf))
                i = j + 1
            else:
                j = i
                while text[j] not in ",)":
                    j += 1
                tok = text[i:j].strip()
                low = tok.lower()
                if low == "null":
                    row.append(None)
                elif low in ("true", "false"):
                    row.append(low == "true")
                elif re.fullmatch(r"-?\d+", tok):
                    row.append(int(tok))
                else:
                    row.append(tok)
                i = j
            while text[i] in " \n\t":
                i += 1
            if text[i] == ",":
                i += 1
                continue
            if text[i] == ")":
                i += 1
                break
            raise ValueError("fila mal formada cerca de: " + text[max(0, i - 40): i + 40])
        rows.append(row)
    return rows


def bloque(sql, tabla):
    m = re.search(r"INSERT INTO salvia\." + tabla + r" \([^)]*\) VALUES(.*?)ON CONFLICT", sql, re.S)
    if not m:
        raise SystemExit("no se encontró el INSERT de salvia." + tabla + " en el seed")
    return parse_rows(m.group(1))


def construir():
    with open(SEED, encoding="utf-8") as f:
        sql = f.read()
    # question: id, form_id, form_section_id, repeater_group_id, question_type, description, required, state_options_path, "order"
    preguntas = [q for q in bloque(sql, "question") if q[3] == BARRIER_REPEATER]
    preguntas.sort(key=lambda q: q[8])
    if len(preguntas) != len(KEYS):
        raise SystemExit("el grupo de la barrera tiene %d preguntas y se esperaban %d: el seed cambió de forma" % (len(preguntas), len(KEYS)))
    # option: id, question_id, label, value, "order"
    opciones = {}
    for o in bloque(sql, "option"):
        opciones.setdefault(o[1], []).append(o)
    # visibility_condition: id, target_type, target_id, trigger_question_id, trigger_state_path, trigger_value, operator
    cond = {}
    for c in bloque(sql, "visibility_condition"):
        if c[1] == "QUESTION":
            cond.setdefault(c[2], []).append(c)

    clave_por_id = {q[0]: KEYS[i][0] for i, q in enumerate(preguntas)}
    out = []
    for i, q in enumerate(preguntas):
        clave, tipo_esperado = KEYS[i]
        if q[4] != tipo_esperado:
            raise SystemExit("la pregunta %d (%s) es de tipo %s y se esperaba %s" % (i + 1, clave, q[4], tipo_esperado))
        if clave in ETIQUETAS_ESPERADAS and q[5] != ETIQUETAS_ESPERADAS[clave]:
            raise SystemExit("la etiqueta de %s cambió en el seed: %r" % (clave, q[5]))
        item = {"key": clave, "tipo": q[4], "label": q[5], "required": bool(q[6])}
        if clave in FUENTE:
            item["fuente"] = FUENTE[clave]
        ops = sorted(opciones.get(q[0], []), key=lambda o: o[4])
        if ops:
            item["opciones"] = [{"value": o[3], "label": o[2]} for o in ops]
        reglas = []
        for c in cond.get(q[0], []):
            if c[3] not in clave_por_id:
                raise SystemExit("la condición de %s depende de una pregunta fuera del grupo" % clave)
            reglas.append({"pregunta": clave_por_id[c[3]], "op": c[6], "valor": c[5]})
        if len(reglas) > 1:
            raise SystemExit("%s tiene %d condiciones; el formulario solo soporta una" % (clave, len(reglas)))
        if reglas:
            item["cuando"] = reglas[0]
        out.append(item)
    return out


def construir_entidad():
    with open(SEED, encoding="utf-8") as f:
        sql = f.read()
    preguntas = bloque(sql, "question")
    opciones = {}
    for o in bloque(sql, "option"):
        opciones.setdefault(o[1], []).append(o)
    rep = sorted([q for q in preguntas if q[3] == ENTITY_REPEATER], key=lambda q: q[8])
    if len(rep) != len(ENTIDAD_KEYS):
        raise SystemExit("el repetidor de la entidad tiene %d preguntas y se esperaban %d: el seed cambió de forma" % (len(rep), len(ENTIDAD_KEYS)))
    conocimiento = [q for q in preguntas if q[2] == ENTITY_SECTION and q[3] is None and q[4] == "boolean"]
    if len(conocimiento) != 1:
        raise SystemExit("no se encontró la pregunta ident.conocimiento")
    out = {"conocimiento": {"label": conocimiento[0][5], "required": bool(conocimiento[0][6])}}
    for i, q in enumerate(rep):
        clave, tipo = ENTIDAD_KEYS[i]
        if q[4] != tipo:
            raise SystemExit("la pregunta %d de la entidad (%s) es de tipo %s y se esperaba %s" % (i + 1, clave, q[4], tipo))
        item = {"label": q[5], "required": bool(q[6])}
        ops = sorted(opciones.get(q[0], []), key=lambda o: o[4])
        if ops:
            item["opciones"] = [{"value": o[3], "label": o[2]} for o in ops]
        out[clave] = item
    esperadas = {"acuerdo": "¿La persona está de acuerdo", "debida": "¿Salvia debe activar rutas", "canales": "Seleccione los canales", "gestiones": "¿Qué gestiones se realizaron"}
    for clave, inicio in esperadas.items():
        if not out[clave]["label"].startswith(inicio):
            raise SystemExit("la etiqueta de %s cambió en el seed: %r" % (clave, out[clave]["label"]))
    canales = [o["value"] for o in out["canales"].get("opciones", [])]
    # Seguimiento v4 ofrece los 4 canales y «Ninguna» (que no crea tarea: v4ProcesarCanales la ignora).
    if canales != ["notificacion", "articulacion", "denuncia", "misalvia", "ninguna"]:
        raise SystemExit("los canales del seed cambiaron: %r" % canales)
    return out


def a_js_entidad(entidad):
    lineas = [
        BEGIN_E + " (generado por scripts/gen_registro_v3_barrera_catalog.py a partir de",
        "// src/internal/seed/seguimiento_v4_form.sql; NO editar a mano: vuelve a ejecutar el script)",
        "// Textos de «Identificación de entidades» de Seguimiento v4 (REG-17/19/20).",
        "var ENTIDAD_V4_V3 = {",
    ]
    for k, v in entidad.items():
        lineas.append("  " + json.dumps(k) + ": " + json.dumps(v, ensure_ascii=False) + ",")
    lineas.append("};")
    lineas.append("// Canales de activación (REG-20): las opciones de la pregunta de canales de Seguimiento v4 (4 canales y «Ninguna»).")
    lineas.append("var CANAL_ACTIVACION_V3 = ENTIDAD_V4_V3.canales.opciones;")
    lineas.append(END_E)
    return "\n".join(lineas)


def a_js(preguntas):
    lineas = [
        BEGIN + " (generado por scripts/gen_registro_v3_barrera_catalog.py a partir de",
        "// src/internal/seed/seguimiento_v4_form.sql; NO editar a mano: vuelve a ejecutar el script)",
        "// Formulario de barrera de Seguimiento v4 (REG-18). `cuando` = condición de visibilidad",
        "// ({pregunta, op: EQUALS|INCLUDES, valor}); `fuente` = opciones que vienen de la geografía.",
        "var BARRERA_V4_V3 = [",
    ]
    for p in preguntas:
        lineas.append("  " + json.dumps(p, ensure_ascii=False) + ",")
    lineas.append("];")
    lineas.append(END)
    return "\n".join(lineas)


def reemplazar(js, begin, end, nuevo, archivo):
    patron = re.compile(re.escape(begin) + r".*?" + re.escape(end), re.S)
    if not patron.search(js):
        raise SystemExit("faltan los marcadores %s / %s en %s" % (begin, end, archivo))
    return patron.sub(lambda _m: nuevo, js, count=1)


def main():
    preguntas = construir()
    entidad = construir_entidad()
    with open(TARGET, encoding="utf-8") as f:
        js = f.read()
    resultado = reemplazar(js, BEGIN, END, a_js(preguntas), TARGET)
    resultado = reemplazar(resultado, BEGIN_E, END_E, a_js_entidad(entidad), TARGET)
    if "--check" in sys.argv:
        if resultado != js:
            raise SystemExit("BARRERA_V4_V3 / ENTIDAD_V4_V3 están desactualizados: ejecuta scripts/gen_registro_v3_barrera_catalog.py")
        print("BARRERA_V4_V3 y ENTIDAD_V4_V3 están al día")
        return
    if resultado != js:
        with open(TARGET, "w", encoding="utf-8") as f:
            f.write(resultado)
    print("BARRERA_V4_V3: %d preguntas · ENTIDAD_V4_V3: %d textos · escritos en %s" % (len(preguntas), len(entidad), os.path.relpath(TARGET, ROOT)))


if __name__ == "__main__":
    main()
