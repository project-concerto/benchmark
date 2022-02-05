# Spree

This document shows how to reproduce the result of the experiments labeled `PBC` in Section 5.2.

## Set up Spree

Spree setup is based on the official [spree_starter](https://github.com/spree/spree_starter). All related files are located in this directory.
The `spree_starter` subdirectory contains spree docker environment, modified based on official docker environment.
The `spree` subdirectory contains the modified spree code.

First, `cd` into the `spree_starter` directory.

Execute `./bin/setup` to set up the docker environment. You may be prompted to enter the admin account and password. Just hit enter to choose the default value. After everything is done, hit `Ctrl-C` to stop the running process.

Finally, execute `docker-compose -d` to start the docker environment.

## Run the benchmark

### Requirement / Installation

This benchmark suite is written in go. So you will need to install go.
You will also have to run `go install` in this directory to install the dependencies.

### How to Run

In the `client` directory, run `run.sh` to start benchmarking.

### Results

The results will be written to 4 CSV files in this directory,
with the following prefixes:

1. `DBT`: Using database transaction in the contentious workload
2. `AHT`: Using ad hoc transaction in the contentious workload
3. `NCDBT`: Using database transaction in the no-contention workload
4. `NCAHT`: Using ad hoc transaction in the no-contention workload
