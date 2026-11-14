# syntax=docker/dockerfile:1

# References to base distro images
ARG ALPINE_VERSION=3 UBUNTU_VERSION=26.04
FROM alpine:${ALPINE_VERSION} AS musl-base
FROM ubuntu:${UBUNTU_VERSION} AS glibc-base

#################################################################################################

# docker buildx build --pull --file Dockerfile --target musl --tag ghcr.io/viral32111/glow-cad-exporter:latest-musl .
# docker container run --publish 8080:8080 --rm ghcr.io/viral32111/glow-cad-exporter:latest-musl

# Minimal image for musl compatibility
FROM scratch AS musl

# https://docs.docker.com/engine/reference/builder/#automatic-platform-args-in-the-global-scope
ARG TARGETARCH

# Copy dependent libraries
COPY --from=musl-base --chown=0:0 --chmod=755 /lib/ld-musl-*.so.1 /lib/
COPY --from=musl-base --chown=0:0 --chmod=755 /lib/libc.musl-*.so.1 /lib/

# Copy built binary from context
COPY --chown=0:0 --chmod=755 *-linux-${TARGETARCH}-musl /usr/local/bin/glow-cad-exporter

# Downprivilege to nobody:nogroup
USER 65534:65534

# Default configuration
ENV GLOW_HTTP_ADDRESS=0.0.0.0

# Publish HTTP server port
EXPOSE 8080/tcp

# Launch the metrics exporter
ENTRYPOINT [ "glow-cad-exporter" ]

#################################################################################################

# docker buildx build --pull --file Dockerfile --target glibc --tag ghcr.io/viral32111/glow-cad-exporter:latest-glibc .
# docker container run --publish 8080:8080 --rm ghcr.io/viral32111/glow-cad-exporter:latest-glibc

# Minimal image for glibc compatibility
FROM scratch AS glibc

# https://docs.docker.com/engine/reference/builder/#automatic-platform-args-in-the-global-scope
ARG TARGETARCH

# Copy dependent libraries
COPY --from=glibc-base --chown=0:0 --chmod=755 /lib64/ld-linux-*.so.2 /lib64/
COPY --from=glibc-base --chown=0:0 --chmod=755 /usr/lib/*-linux-gnu/libc.so.6 /usr/lib/
COPY --from=glibc-base --chown=0:0 --chmod=755 /usr/lib/*-linux-gnu/libresolv.so.2 /usr/lib/

# Copy built binary from context
COPY --chown=0:0 --chmod=755 *-linux-${TARGETARCH}-glibc /usr/local/bin/glow-cad-exporter

# Downprivilege to nobody:nogroup
USER 65534:65534

# Default configuration
ENV GLOW_HTTP_ADDRESS=0.0.0.0

# Publish HTTP server port
EXPOSE 8080/tcp

# Launch the metrics exporter
ENTRYPOINT [ "glow-cad-exporter" ]
