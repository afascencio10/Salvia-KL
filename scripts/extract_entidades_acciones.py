"""Genera entidades_acciones.json desde la hoja 4.EntidadesAcciones.

Extraccion fiel de las columnas A (sector), B (entidad) y E (acciones).
El mapeo a codigos de entity_sector y los alias contra los nombres que ya
existen en la base viven en Go, no aca: esto es solo el volcado del Excel.
"""
import json
import sys

import openpyxl

XLSX, SALIDA = sys.argv[1], sys.argv[2]

wb = openpyxl.load_workbook(XLSX, data_only=True)
ws = wb["4.EntidadesAcciones"]

entidades = []
indice = {}
sector = entidad = None

for fila in ws.iter_rows(min_row=3, values_only=True):
    s, e, accion = fila[0], fila[1], fila[4]
    if s and str(s).strip():
        sector = str(s).strip()
    if e and str(e).strip():
        entidad = " ".join(str(e).split())
    if not entidad:
        continue
    if entidad not in indice:
        indice[entidad] = {"nombre": entidad, "sectorExcel": sector, "acciones": []}
        entidades.append(indice[entidad])
    if accion and str(accion).strip():
        texto = " ".join(str(accion).split())
        if texto not in indice[entidad]["acciones"]:
            indice[entidad]["acciones"].append(texto)

catalogo = {
    "source": "2. 08-09-2026_FormularioRutas.xlsx - hoja 4.EntidadesAcciones (columnas B y E)",
    "entidades": entidades,
}
with open(SALIDA, "w", encoding="utf-8", newline="\n") as fh:
    json.dump(catalogo, fh, ensure_ascii=False, indent=1)
    fh.write("\n")

print(f"entidades: {len(entidades)}")
print(f"acciones : {sum(len(x['acciones']) for x in entidades)}")
print(f"sectores : {sorted({x['sectorExcel'] for x in entidades})}")
sin = [x["nombre"] for x in entidades if not x["acciones"]]
print(f"sin acciones: {len(sin)} {sin}")
