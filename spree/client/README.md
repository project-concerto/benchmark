The gap lock benchmark is based on the checkout API of spree.

# Requirement / Installation
This benchmark suite is written in go. So you will need to install go.
You will also have to run `go install` in this directory to install the dependencies.

# How to Run
In this directory(`gap-lock`), run `run.sh` to start becnhmarking.

# Results
The results will be written to 4 CSV file in this directory,
with prefix corresponding to what they are intended for:
1. DBT: Database transaction benchmark
2. AHT: Ad hoc transaction benchmark
3. NCDBT: Database transaction without low contention benchmark
3. NCAHT: Ad hoc transaction without low contention benchmark
