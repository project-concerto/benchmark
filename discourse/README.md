# Discourse

This document shows how to reproduce the result of the experiments labeled `AA` and `CBC` in Section 5.2 and the experiments in Section 5.3.

## Set up and run

Firstly, build the image and run a container according to the following command.

```shell
cd benchmark/discourse
docker build -t discourse:v1 .
docker run -it discourse:v1 bash
```

Then, execute the following command to initial environment used for benchmark **in the container**.

```shell
cd /benchmark/discourse
bash initialENV.sh
```

Apply code patch to create different branch for applications. The `aht` branch means using **a**d **h**oc **t**ransaction to coordinate the APIs. The `dbt` branch measn using **d**ata**b**ase **t**ransaction to coordinate the APIs.

```shell
bash patcher.sh
```

Finally, run the benchmark scripts to get the result.

- `run_cbc_benchmark.sh` appends result into file in `clients/cbc/` directory.
- `run_rp_benchmark.sh` appends result into file in `discourse/` directory.
- `run_aa_benchmark.sh` appends result into file in `clients/aa/` directory.

```shell
bash run_cbc_benchmark.sh  # or run_rp_benchmark.sh/run_aa_benchmark.sh
```

## Notes

1. To simplify the procedure, in this document we deploy the client, application and the database system in a single container, which is different from the description in the paper. It should be relatively easy to deploy them separately as we did in the paper. Different deployment setup will cause slightly difference in result.

2. The default setup uses contentious workloads (except for `AA` experiments). To reproduce the result for workload with no contention,
   - For the rollback benchmark (Section 5.3), just remove the code in `run_rp_benchmark.sh` according to comment.
   - For the `CBC` benchmark (Section 5.2), change code in `CreatePostTxn.java` according to comment.
