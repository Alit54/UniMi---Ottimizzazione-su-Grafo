import json
import glob
import networkx as nx
import matplotlib.pyplot as plt
from matplotlib.widgets import Slider

# 1. Caricamento File
files = sorted(glob.glob("../../export/maxflow/shortest_augmenting_path/step_*.json"))

if not files:
    print("Errore: Nessun file JSON trovato.")
    exit()

steps = [json.load(open(f)) for f in files]

# 2. Setup Grafico
fig, ax = plt.subplots(figsize=(12, 8))
plt.subplots_adjust(bottom=0.15)

# Calcolo Layout Fisso (basato sul primo step per coerenza)
G_init = nx.DiGraph()
for n in steps[0]['nodes']: G_init.add_node(n['id'])
# Usa spring_layout con seed fisso per mantenere i nodi fermi
pos = nx.spring_layout(G_init, seed=42, k=2)

def draw(val):
    idx = int(val)
    step = steps[idx]
    ax.clear()

    G = nx.DiGraph()

    # --- Gestione Nodi ---
    node_colors = []
    labels = {}

    for n in step['nodes']:
        nid = n['id']
        G.add_node(nid)

        # Label: ID nodo + Distanza h
        dist_str = "Inf" if n['distance'] >= len(step['nodes']) else str(n['distance'])
        labels[nid] = f"{nid}\n(h={dist_str})"

        # Colori Nodi
        if n['is_current']:
            node_colors.append('#facc15') # Giallo (Current)
        elif n['is_source']:
            node_colors.append('#4ade80') # Verde (Source)
        elif n['is_sink']:
            node_colors.append('#f87171') # Rosso (Sink)
        else:
            node_colors.append('#e2e8f0') # Grigio (Default)

    # --- Gestione Archi (Separazione in Liste) ---
    path_edges = []
    saturated_edges = []
    normal_edges = []
    edge_labels = {}

    for e in step['edges']:
        u, v = e['from'], e['to']
        G.add_edge(u, v)
        edge_labels[(u,v)] = f"{e['flow']}/{e['capacity']}"

        # Logica di categorizzazione
        is_saturated = (e['flow'] == e['capacity'] and e['capacity'] > 0)

        if e['in_path']:
            path_edges.append((u, v))
        elif is_saturated:
            saturated_edges.append((u, v))
        else:
            normal_edges.append((u, v))

    # --- Disegno Layer (Ordine importante!) ---

    # 1. Archi Saturi (Sfondo, tratteggiati, grigio chiaro)
    nx.draw_networkx_edges(G, pos, ax=ax, edgelist=saturated_edges,
                           edge_color='#cbd5e1',
                           width=1,
                           style='dashed',
                           arrowsize=10,
                           connectionstyle="arc3,rad=0")

    # 2. Archi Normali (Intermedi, solidi, grigio scuro)
    nx.draw_networkx_edges(G, pos, ax=ax, edgelist=normal_edges,
                           edge_color='#64748b',
                           width=1.5,
                           style='solid',
                           arrowsize=15,
                           connectionstyle="arc3,rad=0")

    # 3. Archi Path (Primo piano, solidi, arancione acceso, spessi)
    nx.draw_networkx_edges(G, pos, ax=ax, edgelist=path_edges,
                           edge_color='#f59e0b',
                           width=3,
                           style='solid',
                           arrowsize=25,
                           connectionstyle="arc3,rad=0")

    # Disegno Nodi e Label
    nx.draw_networkx_nodes(G, pos, ax=ax, node_color=node_colors, node_size=800, edgecolors='#475569')
    nx.draw_networkx_labels(G, pos, ax=ax, labels=labels, font_size=9, font_weight='bold')

    # Label Archi (Flusso/Capacità)
    nx.draw_networkx_edge_labels(G, pos, ax=ax, edge_labels=edge_labels, font_size=7)

    # Titolo
    title = f"Step {idx}: {step['step_type']} | {step['description']}"
    ax.set_title(title, fontsize=12, fontweight='bold', pad=10)
    ax.axis('off')

# Slider Setup
ax_slider = plt.axes([0.2, 0.05, 0.6, 0.03])
slider = Slider(ax_slider, 'Step', 0, len(steps)-1, valinit=0, valfmt='%0.0f')
slider.on_changed(draw)

if __name__ == '__main__':
    draw(0)
    plt.show()