import json
import sys
import networkx as nx
import matplotlib.pyplot as plt


def visualize(json_file):
    # Carica JSON
    with open(json_file, 'r') as f:
        data = json.load(f)

    # Crea grafo
    if data['directed']:
        G = nx.DiGraph()
    else:
        G = nx.Graph()

    # Aggiungi nodi e archi
    for node in data['nodes']:
        G.add_node(node['id'], name=node['name'])

    for edge in data['edges']:
        G.add_edge(edge['source'], edge['target'], weight=edge['weight'])

    # Visualizza
    plt.figure(figsize=(12, 8))
    pos = nx.spring_layout(G)  # Layout automatico

    # Disegna
    nx.draw_networkx_nodes(G, pos, node_color='lightblue', node_size=700)
    nx.draw_networkx_edges(G, pos, arrows=data['directed'], width=2, connectionstyle='arc3,rad=0.1')
    nx.draw_networkx_labels(G, pos, font_size=12, font_weight='bold')

    # Pesi sugli archi
    edge_labels = {(u, v): f"{d['weight']:.1f}"
                   for u, v, d in G.edges(data=True)}
    nx.draw_networkx_edge_labels(G, pos, edge_labels, font_size=8)

    plt.title(f"Graph: {G.number_of_nodes()} nodes, {G.number_of_edges()} edges")
    plt.axis('off')
    plt.tight_layout()
    plt.show()


if __name__ == '__main__':
    if len(sys.argv) < 2:
        print("Uso: python simple_visualize.py graph.json")
        sys.exit(1)

    visualize(sys.argv[1])