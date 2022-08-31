ARG ARCH="amd64"
ARG OS="linux"
FROM quay.io/prometheus/busybox-${OS}-${ARCH}:glibc
LABEL maintainer="Simon Schneider <dev@raynigon.com>"

COPY auth0_exporter /bin/auth0_exporter

EXPOSE      9776
USER        nobody
ENTRYPOINT  [ "/bin/auth0_exporter" ]