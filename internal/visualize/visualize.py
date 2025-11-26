import json
import glob
import os
import sys
import networkx as nx
import matplotlib.pyplot as plt
from matplotlib.widgets import Slider

# --- CONFIGURAZIONE ---
# Cartelle in cui cercare i file JSON
common_folder = "../../export/maxflow"
SEARCH_PATHS = [
    ("Capacity Scaling", "capacity_scaling"),
    ("Shortest Augmenting Path", "shortest_augmenting_path"),
    ("OutStars SAP", "outstars_shortest_augmenting_path"),
    ("Dinic", "dinic"),
    ("Current Directory", ".")
]

def load_files():
    """Scansiona le directory e chiede all'utente cosa visualizzare."""
    available_sets = []

    for label, path in SEARCH_PATHS:
        files = sorted(glob.glob(os.path.join(common_folder, path, "step_*.json")))
        if files:
            available_sets.append((label, files))

    if not available_sets:
        print("❌ Nessun file 'step_*.json' trovato nelle cartelle supportate.")
        sys.exit(1)

    if len(available_sets) == 1:
        print(f"✅ Trovati file per: {available_sets[0][0]}")
        return available_sets[0][1]

    print("\n🔍 Seleziona l'algoritmo da visualizzare:")
    for i, (label, files) in enumerate(available_sets):
        print(f"  [{i+1}] {label} ({len(files)} steps)")

    try:
        choice = int(input("\n> Inserisci numero: ")) - 1
        if 0 <= choice < len(available_sets):
            return available_sets[choice][1]
        else:
            print("Selezione non valida.")
            sys.exit(1)
    except ValueError:
        sys.exit(1)

# --- NORMALIZZAZIONE DATI ---
def normalize_step_data(step_json):
    """
    Uniforma le differenze tra i JSON di Capacity Scaling e SAP.
    Restituisce strutture standard per il disegno.
    """
    nodes_data = {} # id -> {color, label}
    edges_list = [] # [{u, v, flow, cap, type}]

    # 1. Rilevamento Tipo JSON e Parsing Nodi
    # Caso A: SAP (nodes è lista di oggetti con 'distance')
    if len(step_json['nodes']) > 0 and isinstance(step_json['nodes'][0], dict):
        for n in step_json['nodes']:
            nid = n['id']
            # Gestione colore
            color = '#e2e8f0'
            if n.get('is_current'): color = '#facc15'
            elif n.get('is_source'): color = '#4ade80'
            elif n.get('is_sink'): color = '#f87171'

            # Gestione Label Distanza
            dist = n.get('distance', 0)
            max_n = len(step_json['nodes'])
            dist_lbl = "Inf" if dist >= max_n else str(dist)
            label = f"{nid}\n(h={dist_lbl})"

            nodes_data[nid] = {'color': color, 'label': label}

    # Caso B: Capacity Scaling (nodes è lista di int)
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
        # Gestione from/to vs source/target
        u = e.get('from', e.get('source'))
        v = e.get('to', e.get('target'))

        flow = e['flow']
        cap = e['capacity']

        # Determina Tipo Arco
        etype = 'normal'
        if e.get('in_path'):
            etype = 'path'
        elif e.get('is_saturated') or (flow == cap and cap > 0):
            etype = 'saturated'

        edges_list.append({
            'u': u, 'v': v,
            'flow': flow, 'cap': cap,
            'type': etype
        })

    # 3. Info Aggiuntive per Titolo
    info_text = step_json.get('description', '')
    if 'scaling_delta' in step_json:
        d = step_json['scaling_delta']
        info_text += f" | Δ: {d if d > 0 else 'Fine'}"

    return nodes_data, edges_list, info_text, step_json.get('step_type', '')

# --- MAIN LOGIC ---

files = load_files()
steps = [json.load(open(f)) for f in files]

# Setup Grafico
fig, ax = plt.subplots(figsize=(12, 8))
plt.subplots_adjust(bottom=0.15)

# Calcolo Layout Fisso (usando i nodi del primo step)
G_init = nx.DiGraph()
first_nodes, _, _, _ = normalize_step_data(steps[0])
G_init.add_nodes_from(first_nodes.keys())
pos = nx.spring_layout(G_init, seed=42, k=2)

def draw(val):
    idx = int(val)
    nodes_data, edges_data, desc, stype = normalize_step_data(steps[idx])

    ax.clear()
    G = nx.DiGraph()

    # Add Nodes
    node_colors = []
    labels = {}

    # Ordiniamo per ID per consistenza nel loop
    for nid in sorted(nodes_data.keys()):
        G.add_node(nid)
        data = nodes_data[nid]
        node_colors.append(data['color'])
        labels[nid] = data['label']

    # Process Edges into Layers
    path_edges = []
    sat_edges = []
    norm_edges = []
    edge_labels = {}

    for e in edges_data:
        u, v = e['u'], e['v']
        G.add_edge(u, v)
        edge_labels[(u,v)] = f"{e['flow']}/{e['cap']}"

        if e['type'] == 'path': path_edges.append((u,v))
        elif e['type'] == 'saturated': sat_edges.append((u,v))
        else: norm_edges.append((u,v))

    # --- DRAWING LAYERS ---

    # 1. Saturated (Dashed, Light)
    nx.draw_networkx_edges(G, pos, ax=ax, edgelist=sat_edges,
                           edge_color='#cbd5e1', width=1, style='dashed',
                           arrowsize=10, connectionstyle="arc3,rad=0")

    # 2. Normal (Solid, Dark)
    nx.draw_networkx_edges(G, pos, ax=ax, edgelist=norm_edges,
                           edge_color='#64748b', width=1.5, style='solid',
                           arrowsize=15, connectionstyle="arc3,rad=0")

    # 3. Path (Thick, Orange)
    nx.draw_networkx_edges(G, pos, ax=ax, edgelist=path_edges,
                           edge_color='#f59e0b', width=3.5, style='solid',
                           arrowsize=25, connectionstyle="arc3,rad=0")

    # Nodes & Labels
    nx.draw_networkx_nodes(G, pos, ax=ax, node_color=node_colors, node_size=800, edgecolors='#475569')
    nx.draw_networkx_labels(G, pos, ax=ax, labels=labels, font_size=9, font_weight='bold')
    nx.draw_networkx_edge_labels(G, pos, ax=ax, edge_labels=edge_labels, font_size=8)

    # Title
    full_title = f"Step {idx}"
    if stype: full_title += f": {stype}"
    full_title += f"\n{desc}"

    ax.set_title(full_title, fontsize=12, fontweight='bold', pad=10)
    ax.axis('off')

# Slider
ax_slider = plt.axes([0.2, 0.05, 0.6, 0.03])
slider = Slider(ax_slider, 'Timeline', 0, len(steps)-1, valinit=0, valfmt='%0.0f')
slider.on_changed(draw)

if __name__ == '__main__':
    draw(0)
    plt.show()