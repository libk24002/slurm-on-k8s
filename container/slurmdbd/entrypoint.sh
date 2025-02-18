#!/bin/bash

[ -n "${DEBUG}" ] && set -x

SLURMDBD_CONF_FILE_MOUNTED=${SLURMDBD_CONF_FILE_MOUNTED:-/opt/slurm/slurmdbd.conf}
SLURMDBD_CONF_FILE="/etc/slurm/slurmdbd.conf"

cp -f ${SLURMDBD_CONF_FILE_MOUNTED} ${SLURMDBD_CONF_FILE}
chown slurm:slurm ${SLURMDBD_CONF_FILE}
chmod 600 ${SLURMDBD_CONF_FILE}

exec gosu slurm /usr/sbin/slurmdbd -D
