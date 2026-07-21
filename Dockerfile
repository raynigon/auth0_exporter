ARG TARGETPLATFORM="linux/amd64"
FROM quay.io/prometheus/busybox:glibc
LABEL maintainer="Simon Schneider <dev@raynigon.com>"

COPY $TARGETPLATFORM/auth0_exporter /bin/auth0_exporter

EXPOSE      9776
USER        nobody
ENTRYPOINT  [ "/bin/auth0_exporter" ]