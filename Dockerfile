FROM quay.io/prometheus/busybox:latest

LABEL maintainer="Jorge Niedbalski <jnr@metaklass.org>"

COPY openstack-exporter /bin/openstack-exporter
COPY clouds.yaml /etc/openstack/clouds.yaml

ENTRYPOINT ["bin/openstack-exporter","--os-client-config","/etc/openstack/clouds.yaml", "openstack"]

EXPOSE     9180