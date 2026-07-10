#!/usr/bin/env python3
"""Split the monolithic backend/internal/handlers/handlers.go into per-resource
files within the same package. Each handler type (XxxHandler) goes to <xxx>.go;
package-level (non-method) declarations go to shared.go. Imports are derived
per file by detecting usage of each imported package identifier."""
import os, re
from collections import OrderedDict

SRC = 'backend/internal/handlers/handlers.go'
OUT = 'backend/internal/handlers'

with open(SRC) as f:
    content = f.read()
lines = content.split('\n')

# --- extract import block ---
m = re.search(r'import\s*\((.*?)\)', content, re.DOTALL)
import_block = m.group(1)
imports = []  # (path, pkgname)
for line in import_block.strip().split('\n'):
    line = line.strip()
    if not line:
        continue
    mm = re.match(r'(?:(\w+)\s+)?"([^"]+)"', line)
    if mm:
        alias = mm.group(1)
        path = mm.group(2)
        pkgname = alias if alias else path.split('/')[-1]
        imports.append((path, pkgname))

def is_top_level_decl(line):
    return re.match(r'^(type|func|var|const)\s', line) is not None

decls = []  # (group, [lines])
cur_group = 'shared'
cur = []
i = 0
n = len(lines)
while i < n:
    line = lines[i]
    if is_top_level_decl(line):
        if cur:
            decls.append((cur_group, cur))
            cur = []
        mm = re.match(r'func\s+\((?:\w+\s+\*?(\w+))\)', line)
        if mm:
            recv = mm.group(1)
            grp = re.sub(r'Handler$', '', recv)
            cur_group = grp.lower()
        else:
            tm = re.match(r'type\s+(\w+)', line)
            if tm and tm.group(1).endswith('Handler'):
                grp = re.sub(r'Handler$', '', tm.group(1))
                cur_group = grp.lower()
            else:
                cur_group = 'shared'
        cur = [line]
    else:
        if cur or line.strip():
            cur.append(line)
    i += 1
if cur:
    decls.append((cur_group, cur))

groups = OrderedDict()
for g, ls in decls:
    groups.setdefault(g, []).extend(ls)

def used_imports(code):
    needed = []
    for path, pkg in imports:
        if re.search(r'\b' + re.escape(pkg) + r'\.', code):
            needed.append((path, pkg))
    return needed

os.makedirs(OUT, exist_ok=True)
written = []
for g, ls in groups.items():
    code = '\n'.join(ls).strip()
    if not code:
        continue
    needed = used_imports(code)
    imp_block = ('import (\n' + '\n'.join(f'\t"{p}"' for p, _ in needed) + '\n)\n') if needed else ''
    fname = 'shared' if g == 'shared' else g
    out = f'package handlers\n\n{imp_block}\n{code}\n'
    with open(os.path.join(OUT, fname + '.go'), 'w') as f:
        f.write(out)
    written.append((fname, len(ls)))
print('WROTE FILES:')
for w in written:
    print(' ', w[0] + '.go', w[1], 'lines')
