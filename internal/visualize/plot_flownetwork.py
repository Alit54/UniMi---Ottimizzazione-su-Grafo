import json
import networkx as nx
import matplotlib.pyplot as plt
import sys

# Leggi il JSON
with open(sys.argv[1], 'r') as f:
    data = json.load(f)

# Crea il grafo NetworkX
G = nx.DiGraph()

# Aggiungi nodi
G.add_nodes_from(data['nodes'])

# Aggiungi archi con capacità
for edge in data['edges']:
    G.add_edge(edge['source'], edge['target'], capacity=edge['capacity'])

# Layout del grafo
pos = nx.spring_layout(G, seed=42)

# Disegna i nodi
nx.draw_networkx_nodes(G, pos, node_color='lightblue', node_size=500)

# Evidenzia source e sink
nx.draw_networkx_nodes(G, pos, nodelist=[data['source']], node_color='green', node_size=700, label='Source')
nx.draw_networkx_nodes(G, pos, nodelist=[data['sink']], node_color='red', node_size=700, label='Sink')

# Disegna archi
nx.draw_networkx_edges(G, pos, edge_color='gray', arrows=True, arrowsize=20)

# Etichette nodi
nx.draw_networkx_labels(G, pos, font_size=12)

# Etichette archi (capacità)
edge_labels = nx.get_edge_attributes(G, 'capacity')
nx.draw_networkx_edge_labels(G, pos, edge_labels, font_size=8)

plt.legend()
plt.axis('off')
plt.title(f'Flow Network (N={len(data["nodes"])}, Source={data["source"]}, Sink={data["sink"]})')
plt.tight_layout()
plt.show()