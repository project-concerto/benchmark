# Broadleaf

This document shows how to reproduce the result of the experiments labeled `RMW` in Section 5.2.

## Set up and run

First, build the image and run a container according to the following command.

```shell
cd benchmark/broadleaf
docker build -t blc:v1 .
docker run -it blc:v1 bash
```

Then, execute the following command to install database used for benchmark **in the container**. The default username and password is `root` and `123456` respectively.

```shell
cd /benchmark/broadleaf
bash installDB.sh
```

Apply code patch to create different branch for applications. The `aht` branch means using **a**d **h**oc **t**ransaction to coordinate the API. The `dbt` branch means using **d**ata**b**ase **t**ransaction to coordinate the API.

```shell
bash patcher.sh
```

Finally, run the benchmark scripts to get the result, which is appended into file in the `client` directory.

```shell
bash run_benchmark.sh
```

## Notes

1. To simplify the procedure, in this document we deploy the client, application and the database system in a single container, which is different from the description in the paper. It should be relatively easy to deploy them separately as we did in the paper. Different deployment setup will cause slightly difference in result.

2. The default setup uses contended workload. To reproduce the result for workload with no contention, change the code in `BLCtxn.java` according to comment.
