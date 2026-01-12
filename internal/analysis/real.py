import pandas as pd
import os

def clean_algo_name(name):
    name = name.replace("Capacity Scaling", "CS")
    name = name.replace("Shortest Augmenting Path", "SAP")
    name = name.replace("con SAP", "+ SAP")
    return name

def clean_dataset_name(name):
    base = os.path.basename(name)
    base = base.replace("BVZ-tsukuba", "Tsukuba ").replace(".max", "").replace("7", "07")
    return base

def generate_realworld_text(csv_path, output_text="internal/analysis/real.tex"):
    df = pd.read_csv(csv_path)

    df = df[~df['Graph'].str.contains('custom', case=False)].copy()

    df['Time_ms'] = df['Time (nanoseconds)'] / 1_000_000.0
    df['Dataset_Name'] = df['Graph'].apply(clean_dataset_name)

    algorithms = sorted(df['Algorithm'].unique())
    datasets = sorted(df['Dataset_Name'].unique())

    with open(output_text, "w") as f:
        f.write(r"""
\begin{tikzpicture}
    \begin{semilogyaxis}[
        ybar,
        width=0.95\textwidth,
        height=0.6\textwidth,
        ylabel={Tempo di Esecuzione (ms)},
        symbolic x coords={""" + ",".join(datasets) + r"""},
        xtick=data,
        point meta=y,
        legend style={at={(0.5,-0.15)}, anchor=north, legend columns=-1},
        ymajorgrids=true,
        grid style={dashed, gray!30},
        ymin=1,
        bar width=7pt,
        enlarge x limits=0.2
    ]
""")

        for algo in algorithms:
            subset = df[df['Algorithm'] == algo]
            clean_algo = clean_algo_name(algo)

            f.write(f"    \\addplot coordinates {{\n")
            for dataset in datasets:
                row = subset[subset['Dataset_Name'] == dataset]
                if not row.empty:
                    time_val = row.iloc[0]['Time_ms']
                    f.write(f"        ({dataset}, {time_val:.2f})\n")
            f.write("    };\n")
            f.write(f"    \\addlegendentry{{{clean_algo}}}\n")

        f.write(r"""
    \end{semilogyaxis}
\end{tikzpicture}
""")

    print(f"File '{output_text}' generato.")

if __name__ == "__main__":
    generate_realworld_text('export/benchmark_results.csv')