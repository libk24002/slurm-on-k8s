#!/bin/bash

[ -n "${DEBUG}" ] && set -x

SLURM_CONF_FILE_MOUNTED=${SLURM_CONF_FILE_MOUNTED:-/opt/slurm/slurm.conf}
SLURM_CONF_FILE="/etc/slurm/slurm.conf"

chown slurm:slurm ${SLURM_CONF_FILE}
chmod 600 ${SLURM_CONF_FILE}

exec gosu slurm /usr/sbin/slurmctld -D
