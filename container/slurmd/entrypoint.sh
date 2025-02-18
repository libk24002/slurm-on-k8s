#!/bin/bash


if [ -z "$MPI_TYPE" ]; then
    echo "ENV 'MPI_TYPE' not exist."
else
    echo "ENV 'MPI_TYPE' exist, and its value : $MPI_TYPE"
fi


if [ "$MPI_TYPE" = "opem-mpi" ]; then
    echo "not need to update ~/.bashrc"
elif [ "$MPI_TYPE" = "intel-mpi" ]; then
    echo 'export PATH="/opt/intel/oneapi/mpi/2021.14/bin:$PATH"' >> ~/.set_env.sh
    echo 'source /opt/intel/oneapi/setvars.sh' >> ~/.bashrc
    source ~/.set_env.sh
    echo "finish set intel-mpi environment."
else
    echo "unknow mpi type, please check."
fi



[ -n "${DEBUG}" ] && set -x

ulimit -l unlimited
ulimit -s unlimited
ulimit -n 131072
ulimit -a

exec gosu root /usr/sbin/slurmd -D


