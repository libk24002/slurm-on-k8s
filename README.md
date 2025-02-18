## Slurm on K8s

### slurm configuration generator [click](https://slurm.schedmd.com/configurator.html)

### build images in container

```shell
bash slurm/container/builder/build-in-container.sh
bash slurm/container/base/build-in-container.sh
bash slurm/container/munged/build-in-container.sh
bash slurm/container/login/build-openmpi-in-container.sh
bash slurm/container/slurmctld/build-in-container.sh
bash slurm/container/slurmd/build-openmpi-in-container.sh
bash slurm/container/slurmdbd/build-in-container.sh
```

### build images

```shell
bash slurm/container/builder/build.sh
bash slurm/container/base/build.sh
bash slurm/container/munged/build.sh
MPI_TYPE=open-mpi bash slurm/container/login/build.sh
bash slurm/container/slurmctld/build.sh
MPI_TYPE=open-mpi bash slurm/container/slurmd/build.sh
bash slurm/container/slurmdbd/build.sh
```

### build Intel MPI images
```shell
MPI_TYPE=intel-mpi bash slurm/container/login/build.sh
MPI_TYPE=intel-mpi bash slurm/container/slurmd/build.sh
```

### publish images
```shell
MPI_TYPE=intel-mpi bash slurm/container/load-into-minikube.sh
```

### publish helm chart

```shell
helm package --dependency-update --destination /tmp/ ./slurm/chart
```
or 
```shell
helm package --destination /tmp/ ./slurm/chart
```

### deploy helm chart

```shell
helm install slurm /tmp/slurm-1.0.2.tgz
```
