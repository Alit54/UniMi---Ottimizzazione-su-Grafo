import json
import glob
import os
import sys
import networkx as nx
import matplotlib.pyplot as plt
from matplotlib.widgets import Slider

common_folder = "../../export/maxflow"
SEARCH_PATHS = [
    ("Capacity Scaling", "capacity_scaling"),
    ("Shortest Augmenting Path", "shortest_augmenting_path"),
    ("Dinic", "dinic"),
]

def load_files():
    available_sets = []
    for label, path in SEARCH_PATHS:
        files = sorted(glob.glob(os.path.join(common_folder, path, "step_*.json")))
        if files:
            available_sets.append((label, files))

    if not available_sets:
        print("❌ Nessun file trovato.")
        sys.exit(1)

    if len(available_sets) == 1:
        print(f"✅ Caricato: {available_sets[0][0]}")
        return available_sets[0][1]

    print("\n🔍 Seleziona Algoritmo:")
    for i, (label, files) in enumerate(available_sets):
        print(f"  [{i+1}] {label} ({len(files)} steps)")

    try:
        choice = int(input("\n> Numero: ")) - 1
        if 0 <= choice < len(available_sets):
            return available_sets[choice][1]
    except ValueError: pass
    sys.exit(1)

def normalize_step_data(step_json):
    nodes_data = {}
    edges_list = []

    is_cap_scaling = 'scaling_delta' in step_json
    current_delta = step_json.get('scaling_delta', 0)

    if len(step_json['nodes']) > 0 and isinstance(step_json['nodes'][0], dict):
        for n in step_json['nodes']:
            nid = n['id']
            color = '#e2e8f0'
            if n.get('is_current'): color = '#facc15'
            elif n.get('is_source'): color = '#4ade80'
            elif n.get('is_sink'): color = '#f87171'

            label = f"{nid}"
            if 'level' in n: label += f"\n(L={n['level'] if n['level']!=-1 else 'Inf'})"
            elif 'distance' in n:
                dist = n['distance']
                max_n = len(step_json['nodes'])
                label += f"\n(h={dist if dist < max_n else 'Inf'})"

            nodes_data[nid] = {'color': color, 'label': label}
    else:
        src = step_json.get('source', -1)
        snk = step_json.get('sink', -1)
        for nid in step_json['nodes']:
            color = '#e2e8f0'
            if nid == src: color = '#4ade80'
            elif nid == snk: color = '#f87171'
            nodes_data[nid] = {'color': color, 'label': str(nid)}

    # 2. Parsing Archi
    for e in step_json['edges']:
        u = e.get('from', e.get('source'))
        v = e.get('to', e.get('target'))
        flow = e['flow']
        cap = e['capacity']
        residual = cap - flow

        etype = 'normal'

        if e.get('in_path'):
            etype = 'path'
        elif e.get('is_saturated') or (residual == 0 and cap > 0):
            etype = 'saturated'
        elif is_cap_scaling and 0 < current_delta <= residual:
            etype = 'valid'

        edges_list.append({'u': u, 'v': v, 'flow': flow, 'cap': cap, 'type': etype})

    info = step_json.get('description', '')
    if is_cap_scaling:
        info += f" | Δ: {current_delta if current_delta > 0 else 'Fine'}"

    return nodes_data, edges_list, info, step_json.get('step_type', '')

files = load_files()
steps = [json.load(open(f)) for f in files]

fig, ax = plt.subplots(figsize=(12, 8))
plt.subplots_adjust(bottom=0.15)

G_init = nx.DiGraph()
first_nodes, _, _, _ = normalize_step_data(steps[0])
G_init.add_nodes_from(first_nodes.keys())
pos = nx.spring_layout(G_init, seed=42, k=2)

def draw(val):
    idx = int(val)
    nodes_data, edges_data, desc, stype = normalize_step_data(steps[idx])

    ax.clear()
    G = nx.DiGraph()

    node_colors, labels = [], {}
    for nid in sorted(nodes_data.keys()):
        G.add_node(nid)
        node_colors.append(nodes_data[nid]['color'])
        labels[nid] = nodes_data[nid]['label']

    path_e, sat_e, valid_e, norm_e, e_labels = [], [], [], [], {}

    for e in edges_data:
        u, v = e['u'], e['v']
        G.add_edge(u, v)
        e_labels[(u,v)] = f"{e['flow']}/{e['cap']}"

        t = e['type']
        if t == 'path': path_e.append((u,v))
        elif t == 'saturated': sat_e.append((u,v))
        elif t == 'valid': valid_e.append((u,v))
        else: norm_e.append((u,v))

    # 1. Saturated (Sfondo, Tratteggiato)
    nx.draw_networkx_edges(G, pos, ax=ax, edgelist=sat_e, edge_color='#cbd5e1',
                           width=1, style='dashed', arrowsize=20, connectionstyle="arc3,rad=0")

    # 2. Normal (Grigio Scuro - Non percorribili per Delta)
    nx.draw_networkx_edges(G, pos, ax=ax, edgelist=norm_e, edge_color='#94a3b8',
                           width=1.5, style='solid', arrowsize=20, alpha=0.6, connectionstyle="arc3,rad=0")

    # 3. Valid (Blu - Percorribili con Delta attuale)
    nx.draw_networkx_edges(G, pos, ax=ax, edgelist=valid_e, edge_color='#2563eb',
                           width=2, style='solid', arrowsize=20, connectionstyle="arc3,rad=0")

    # 4. Path (Arancione - Evidenziato)
    nx.draw_networkx_edges(G, pos, ax=ax, edgelist=path_e, edge_color='#f59e0b',
                           width=3.5, style='solid', arrowsize=25, connectionstyle="arc3,rad=0")

    nx.draw_networkx_nodes(G, pos, ax=ax, node_color=node_colors, node_size=1000, edgecolors='#475569')
    nx.draw_networkx_labels(G, pos, ax=ax, labels=labels, font_size=9, font_weight='bold')
    nx.draw_networkx_edge_labels(G, pos, ax=ax, edge_labels=e_labels, font_size=8)

    ax.set_title(f"Step {idx}: {stype}\n{desc}", fontsize=12, fontweight='bold', pad=10)
    ax.axis('off')

slider = Slider(plt.axes([0.2, 0.05, 0.6, 0.03]), 'Step', 0, len(steps)-1, valinit=0, valfmt='%0.0f')
slider.on_changed(draw)

if __name__ == '__main__':
    draw(0)
    plt.show()