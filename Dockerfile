FROM quay.io/prometheus/busybox:latest

ARG OS=linux
ARG ARCH=amd64

LABEL maintainer="Jorge Niedbalski <jnr@metaklass.org>"

COPY .build/$OS-$ARCH/openstack-exporter /bin/openstack-exporter

ENTRYPOINT ["bin/openstack-exporter","--os-client-config","/etc/openstack/clouds.yaml", "openstack"]
EXPOSE     9180
