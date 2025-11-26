import json
import glob
import networkx as nx
import matplotlib.pyplot as plt
from matplotlib.widgets import Button, Slider

files = sorted(glob.glob("../../export/maxflow/shortest_augmenting_path/step_*.json"))
if not files:
    print("Nessun file JSON trovato.")
    exit()

steps = [json.load(open(f)) for f in files]
fig, ax = plt.subplots(figsize=(10, 7))
plt.subplots_adjust(bottom=0.2)

# Layout fisso basato sul primo step
G_init = nx.DiGraph()
for n in steps[0]['nodes']: G_init.add_node(n['id'])
pos = nx.spring_layout(G_init, seed=42)

def draw(val):
    idx = int(val)
    step = steps[idx]
    ax.clear()

    G = nx.DiGraph()
    labels = {}
    node_colors = []
    node_borders = []

    for n in step['nodes']:
        nid = n['id']
        G.add_node(nid)
        # Mostra ID e Distanza (Label h)
        labels[nid] = f"{nid}\n(h={n['distance']})"

        if n['is_current']:
            node_colors.append('#facc15') # Giallo per nodo corrente
        elif n['is_source']:
            node_colors.append('#4ade80') # Verde
        elif n['is_sink']:
            node_colors.append('#f87171') # Rosso
        else:
            node_colors.append('#e2e8f0')

    edge_colors = []
    edge_labels = {}

    for e in step['edges']:
        u, v = e['from'], e['to']
        G.add_edge(u, v)
        edge_labels[(u,v)] = f"{e['flow']}/{e['capacity']}"

        if e['in_path']:
            edge_colors.append('orange')
        else:
            edge_colors.append('gray')

    nx.draw_networkx_nodes(G, pos, ax=ax, node_color=node_colors, node_size=700)
    nx.draw_networkx_labels(G, pos, ax=ax, labels=labels, font_size=9)
    nx.draw_networkx_edges(G, pos, ax=ax, edge_color=edge_colors, connectionstyle="arc3,rad=0.1")
    nx.draw_networkx_edge_labels(G, pos, ax=ax, edge_labels=edge_labels, font_size=8)

    ax.set_title(f"Step {idx}: {step['step_type']}\n{step['description']}")
    ax.axis('off')

ax_slider = plt.axes([0.2, 0.05, 0.6, 0.03])
slider = Slider(ax_slider, 'Step', 0, len(steps)-1, valinit=0, valfmt='%0.0f')
slider.on_changed(draw)

if __name__ == '__main__':
    draw(0)
    plt.show()