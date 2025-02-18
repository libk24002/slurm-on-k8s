## Slurm on K8s

### Prerequisites
- kubectl version v1.11.3+.
- buildah version v1.33.10+
- Access to a Kubernetes v1.11.3+ cluster.

### Features

Slurm on Kubernetes provides the following features:

- **Resource Management**: Efficiently manages resources in a Kubernetes cluster, ensuring optimal utilization.
- **Job Scheduling**: Advanced job scheduling capabilities to handle various types of workloads.
- **Scalability**: Easily scales to accommodate growing workloads and resources.
- **High Availability**: Supports high availability configurations to ensure continuous operation.
- **Multi-User Support**: Allows multiple users to submit and manage their jobs concurrently.
- **Integration with MPI Libraries**: Supports both Open MPI and Intel MPI libraries for parallel computing.
- **Customizable**: Using `values.yaml` file, you can customizable a slurm cluster, fitting specific needs and configurations.
- **munged**

### Usage

> if you wanna change slurm configuration ,please check slurm configuration generator [click](https://slurm.schedmd.com/configurator.html)

- for user
    > 
    1. download helm chart `wget xxx`
    2. update `values.yaml`
    3. and exec `helm install slurm /tmp/slurm-1.0.2.tgz`

- for developer



### Build
### build images (default support Open MPI Libs)
```shell
bash ./container/builder/build.sh
bash ./container/base/build.sh
bash ./container/munged/build.sh
MPI_TYPE=open-mpi bash ./container/login/build.sh
bash ./container/slurmctld/build.sh
MPI_TYPE=open-mpi bash ./container/slurmd/build.sh
bash ./container/slurmdbd/build.sh
```

### build images with Intel MPI Libs
```shell
MPI_TYPE=intel-mpi bash ./container/login/build.sh
MPI_TYPE=intel-mpi bash ./container/slurmd/build.sh
```

### publish images
```shell
MPI_TYPE=intel-mpi bash ./container/load-into-minikube.sh
```

### publish helm chart
```shell
helm package --dependency-update --destination /tmp/ ./chart
```

### deploy helm chart

```shell
helm install slurm /tmp/slurm-1.0.2.tgz
```

