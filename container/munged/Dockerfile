ARG OS_BASE_IMAGE=localhost/ay-dev/slurm-base:latest
FROM ${OS_BASE_IMAGE}

COPY entrypoint.sh /opt/slurm/entrypoint.sh
RUN chmod 755 /opt/slurm/entrypoint.sh

CMD ["/opt/slurm/entrypoint.sh"]