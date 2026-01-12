import pandas as pd
import re

def clean_algo_name(name):
    return name.replace("Capacity Scaling", "Cap. Scal.").replace("Shortest Augmenting Path", "SAP")

def get_params(g):
    n = re.search(r'n(\d+)', g)
    d = re.search(r'd(\d+\.\d+)', g)
    c = re.search(r'c(\d+)', g)
    return (int(n.group(1)) if n else 0, float(d.group(1)) if d else 0, int(c.group(1)) if c else 0)

def generate_capacity_comparison(csv_path):
    df = pd.read_csv(csv_path)
    df = df[df['Graph'].str.contains('custom', case=False)].copy()

    # Estrazione parametri
    df['Params'] = df['Graph'].apply(get_params)
    df['N'] = df['Params'].apply(lambda x: x[0])
    df['D'] = df['Params'].apply(lambda x: x[1])
    df['C'] = df['Params'].apply(lambda x: x[2])
    df['Time_ms'] = df['Time (nanoseconds)'] / 1_000_000.0

    # Fissiamo N=1000 e D=0.50 (che so essere presenti nei tuoi dati)
    target_n, target_d = 1000, 0.50
    subset = df[(df['N'] == target_n) & (df['D'] == target_d)].copy()

    # Filtriamo solo Capacità 10 e 100000
    subset = subset[subset['C'].isin([10, 100000])]

    # Calcoliamo la media se ci sono più run
    pivoted = subset.groupby(['Algorithm', 'C'])['Time_ms'].mean().unstack()

    with open("internal/analysis/capacity.tex", "w") as f:
        f.write(r"""
\begin{tikzpicture}
    \begin{axis}[
        ybar,
        title={Impatto della Capacità sugli Algoritmi ($N=1000, D=0.5\%$)},
        width=0.95\textwidth,
        height=0.6\textwidth,
        symbolic x coords={""" + ",".join([clean_algo_name(a) for a in pivoted.index]) + r"""},
        xtick=data,
        ylabel={Tempo di Esecuzione (ms)},
        legend style={at={(0.5,-0.15)}, anchor=north, legend columns=-1},
        ymajorgrids=true,
        nodes near coords,
        nodes near coords style={font=\footnotesize},
        ymin=0
    ]
""")
        # Barra per C=10
        f.write(r"    \addplot coordinates {")
        for algo, row in pivoted.iterrows():
            if 10 in row and not pd.isna(row[10]):
                f.write(f"({clean_algo_name(algo)}, {row[10]:.2f}) ")
        f.write(r"};" + "\n")
        f.write(r"    \addlegendentry{Capacità $U=10$}" + "\n")

        # Barra per C=100000
        f.write(r"    \addplot coordinates {")
        for algo, row in pivoted.iterrows():
            if 100000 in row and not pd.isna(row[100000]):
                f.write(f"({clean_algo_name(algo)}, {row[100000]:.2f}) ")
        f.write(r"};" + "\n")
        f.write(r"    \addlegendentry{Capacità $U=100.000$}" + "\n")

        f.write(r"""
    \end{axis}
\end{tikzpicture}
""")
    print("File 'grafico_capacita.tex' generato.")

if __name__ == "__main__":
    generate_capacity_comparison('export/benchmark_results.csv')