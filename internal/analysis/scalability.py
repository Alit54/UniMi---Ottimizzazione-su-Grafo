import pandas as pd

def clean_algorithm_name(name):
    name = name.replace("Capacity Scaling", "CS")
    name = name.replace("Shortest Augmenting Path", "SAP")
    name = name.replace("Cap. Scaling con Sap", "CS con SAP")
    return name

def generate_tex_scatterplot(csv_path, output_text="grafico_scalabilita.tex", low=False):
    df = pd.read_csv(csv_path)
    df = df[df['Graph'].str.contains('custom', case=False)].copy()
    df['Time_ms'] = df['Time (nanoseconds)'] / 1_000_000.0
    df = df.sort_values(by='Edges')
    algorithms = df['Algorithm'].unique()

    if low:
        df = df[df['Graph'].str.contains('c10m', case=False)].copy()
        output_text += '_low.tex'
    else:
        df = df[df['Graph'].str.contains('c100', case=False)].copy()
        output_text += '_high.tex'

    with open(output_text, "w") as f:
        f.write(r"""
\begin{tikzpicture}
    \begin{loglogaxis}[
        width=0.95\textwidth,
        height=0.65\textwidth,
        xlabel={Numero di Archi ($|E|$)},
        ylabel={Tempo di Esecuzione (ms)},
        grid=both,
        grid style={line width=.1pt, draw=gray!20},
        major grid style={line width=.2pt,draw=gray!50},
        legend pos=north west,
        legend cell align={left},
        cycle list name=color list
    ]
""")

        for algo in algorithms:
            subset = df[df['Algorithm'] == algo]
            clean_name = clean_algorithm_name(algo)

            f.write(f"    \\addplot+[only marks, mark=*, mark size=1.8pt, fill opacity=0.6] coordinates {{\n")

            for _, row in subset.iterrows():
                time_val = max(row['Time_ms'], 0.001)
                f.write(f"        ({row['Edges']}, {time_val})\n")

            f.write("    };\n")
            f.write(f"    \\addlegendentry{{{clean_name}}}\n")

        f.write(r"""
    \end{loglogaxis}
\end{tikzpicture}
""")

if __name__ == "__main__":
    generate_tex_scatterplot('export/benchmark_results.csv', output_text='internal/analysis/scalability', low=False)
    generate_tex_scatterplot('export/benchmark_results.csv', output_text='internal/analysis/scalability', low=True)