import json
import glob
import networkx as nx
import matplotlib.pyplot as plt
from matplotlib.widgets import Slider

# 1. Caricamento File JSON
# Cerca nella cartella corrente o in una cartella export comune
files = sorted(glob.glob("../../export/maxflow/capacity_scaling/step_*.json"))

if not files:
    print("Errore: Nessun file 'step_*.json' trovato.")
    exit()

steps = [json.load(open(f)) for f in files]

# 2. Configurazione Layout Grafico
fig, ax = plt.subplots(figsize=(12, 8))
plt.subplots_adjust(bottom=0.15) # Spazio per lo slider

# Calcolo Layout Fisso
# Usiamo i nodi del primo step per calcolare le posizioni
G_init = nx.DiGraph()
first_step = steps[0]
for n in first_step['nodes']:
    G_init.add_node(n)

# Layout a molla con seed fisso per evitare che i nodi "ballino" tra uno step e l'altro
pos = nx.spring_layout(G_init, seed=42, k=2)

def draw(val):
    idx = int(val)
    step = steps[idx]
    ax.clear()

    G = nx.DiGraph()

    # --- Gestione Nodi ---
    node_colors = []
    labels = {}

    # In Capacity Scaling JSON, 'nodes' è una lista di interi [0, 1, 2...]
    # 'source' e 'sink' sono campi separati nella root del JSON
    source_id = step['source']
    sink_id = step['sink']

    for n in step['nodes']:
        G.add_node(n)
        labels[n] = str(n)

        if n == source_id:
            node_colors.append('#4ade80') # Verde (Source)
        elif n == sink_id:
            node_colors.append('#f87171') # Rosso (Sink)
        else:
            node_colors.append('#e2e8f0') # Grigio (Default)

    # --- Gestione Archi e Stili ---
    path_edges = []
    saturated_edges = []
    normal_edges = []
    edge_labels = {}

    # Nota: Nel JSON del Capacity Scaling le chiavi sono 'source' e 'target'
    for e in step['edges']:
        u, v = e['source'], e['target']
        G.add_edge(u, v)

        label = f"{e['flow']}/{e['capacity']}"
        edge_labels[(u,v)] = label

        # Categorizzazione
        # Capacity Scaling JSON ha già il campo 'is_saturated', ma possiamo anche ricalcolarlo
        is_sat = e.get('is_saturated') or (e['flow'] == e['capacity'] and e['capacity'] > 0)

        if e['in_path']:
            path_edges.append((u, v))
        elif is_sat:
            saturated_edges.append((u, v))
        else:
            normal_edges.append((u, v))

    # --- Disegno Layer (Ordine: Saturi -> Normali -> Path) ---

    # 1. Archi Saturi (Sfondo, tratteggiati)
    nx.draw_networkx_edges(G, pos, ax=ax, edgelist=saturated_edges,
                           edge_color='#cbd5e1',
                           width=1,
                           style='dashed',
                           arrowsize=10,
                           connectionstyle="arc3,rad=0")

    # 2. Archi Normali (Intermedi, solidi)
    nx.draw_networkx_edges(G, pos, ax=ax, edgelist=normal_edges,
                           edge_color='#64748b',
                           width=1.5,
                           style='solid',
                           arrowsize=15,
                           connectionstyle="arc3,rad=0")

    # 3. Archi Path (Primo piano, evidenziati)
    nx.draw_networkx_edges(G, pos, ax=ax, edgelist=path_edges,
                           edge_color='#f59e0b',
                           width=3,
                           style='solid',
                           arrowsize=25,
                           connectionstyle="arc3,rad=0")

    # Disegno Nodi e Label
    nx.draw_networkx_nodes(G, pos, ax=ax, node_color=node_colors, node_size=600, edgecolors='#475569')
    nx.draw_networkx_labels(G, pos, ax=ax, labels=labels, font_size=10, font_weight='bold')
    nx.draw_networkx_edge_labels(G, pos, ax=ax, edge_labels=edge_labels, font_size=8)

    # Titolo Informativo
    delta_info = f"Delta: {step['scaling_delta']}" if step['scaling_delta'] > 0 else "Delta: Fine"
    title = f"Step {idx}: {step['description']}\n{delta_info} | Flusso Totale: {step['current_flow']}"

    ax.set_title(title, fontsize=12, fontweight='bold', pad=10)
    ax.axis('off')

# Slider Setup
ax_slider = plt.axes([0.2, 0.05, 0.6, 0.03])
slider = Slider(ax_slider, 'Timeline', 0, len(steps)-1, valinit=0, valfmt='%0.0f')
slider.on_changed(draw)

if __name__ == '__main__':
    draw(0)
    plt.show()