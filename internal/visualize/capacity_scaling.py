import json
import glob
import networkx as nx
import matplotlib.pyplot as plt
from matplotlib.widgets import Button, Slider

# 1. Carica tutti i file JSON generati da Go
files = sorted(glob.glob("../../export/maxflow/capacity_scaling/step_*.json"))
if not files:
    print("Nessun file step_*.json trovato. Esegui prima il codice Go!")
    exit()

steps = []
for f in files:
    with open(f, 'r') as file:
        steps.append(json.load(file))

# Configurazione Visualizzazione
fig, ax = plt.subplots(figsize=(12, 8))
plt.subplots_adjust(bottom=0.2) # Spazio per i controlli
G = nx.DiGraph()

# Calcola layout fisso usando il primo step (per evitare che i nodi saltino)
initial_data = steps[0]
for i in initial_data['nodes']:
    G.add_node(i)
pos = nx.spring_layout(G) # Layout coerente

current_step_idx = 0

def draw_step(val):
    global current_step_idx
    idx = int(val)
    current_step_idx = idx
    step = steps[idx]

    ax.clear()

    # Ricostruisci grafo per questo step
    G.clear()
    edges_list = []
    edge_colors = []
    edge_labels = {}
    node_colors = []

    # Definizione Nodi
    for n in step['nodes']:
        G.add_node(n)
        if n == step['source']:
            node_colors.append('#4ade80') # Source Verde
        elif n == step['sink']:
            node_colors.append('#f87171') # Sink Rosso
        else:
            node_colors.append('#94a3b8') # Nodi interni

    # Definizione Archi
    for edge in step['edges']:
        u, v = edge['source'], edge['target']
        G.add_edge(u, v)

        # Colore Archi
        if edge.get('in_path'):
            edge_colors.append('#f59e0b') # Arancione per percorso corrente
        elif edge.get('is_saturated'):
            edge_colors.append('#cbd5e1') # Grigio chiaro per saturi
        else:
            edge_colors.append('#64748b') # Grigio scuro normale

        # Etichetta Arco: Flusso / Capacità
        label = f"{edge['flow']}/{edge['capacity']}"
        edge_labels[(u, v)] = label

    # Disegno
    nx.draw_networkx_nodes(G, pos, ax=ax, node_color=node_colors, node_size=500)
    nx.draw_networkx_labels(G, pos, ax=ax, font_color='white', font_weight='bold')

    nx.draw_networkx_edges(G, pos, ax=ax, edge_color=edge_colors,
                           width=2, arrowsize=20, connectionstyle="arc3,rad=0.1")

    nx.draw_networkx_edge_labels(G, pos, ax=ax, edge_labels=edge_labels,
                                 font_size=8, label_pos=0.7)

    # Titolo e Info
    ax.set_title(f"Step {idx}: {step['description']}\n"
                 f"Delta: {step['scaling_delta']} | Tot Flow: {step['current_flow']}",
                 fontsize=12, fontweight='bold')
    ax.axis('off')

    # Aggiorna slider se chiamato da bottone
    if slider.val != idx:
        slider.set_val(idx)

# Widget Slider
ax_slider = plt.axes([0.2, 0.05, 0.6, 0.03])
slider = Slider(ax_slider, 'Step', 0, len(steps)-1, valinit=0, valfmt='%0.0f')
slider.on_changed(draw_step)

if __name__ == '__main__':
    draw_step(0)
    plt.show()