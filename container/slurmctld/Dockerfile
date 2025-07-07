ARG OS_BASE_IMAGE=localhost/ay-dev/slurm-base:latest
FROM ${OS_BASE_IMAGE}

ARG SLURM_VERSION=25.05.0
ARG BUILDER_IMAGE=localhost/ay-dev/slurm-builder:latest
COPY --from=localhost/ay-dev/slurm-builder:latest /opt/slurm/debs/slurm-smd-slurmctld-dbgsym_${SLURM_VERSION}-1_amd64.deb /opt/slurm/debs/slurm-smd-slurmctld-dbgsym_${SLURM_VERSION}-1_amd64.deb
COPY --from=localhost/ay-dev/slurm-builder:latest /opt/slurm/debs/slurm-smd-slurmctld_${SLURM_VERSION}-1_amd64.deb /opt/slurm/debs/slurm-smd-slurmctld_${SLURM_VERSION}-1_amd64.deb

RUN apt install -y /opt/slurm/debs/slurm-smd-slurmctld-dbgsym_${SLURM_VERSION}-1_amd64.deb /opt/slurm/debs/slurm-smd-slurmctld_${SLURM_VERSION}-1_amd64.deb

COPY entrypoint.sh /opt/slurm/entrypoint.sh
RUN chmod 755 /opt/slurm/entrypoint.sh

CMD ["/opt/slurm/entrypoint.sh"]