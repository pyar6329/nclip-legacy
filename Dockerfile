# syntax=docker/dockerfile:1.0-experimental
FROM dockercore/golang-cross:1.13.8 AS build-base

RUN set -x && \
  apt-get update && \
  apt-get install -y \
    xz-utils \
    musl-dev \
    musl-tools && \
  apt-get clean && \
  rm -rf /var/lib/apt/lists/*

FROM build-base AS darwin-binary

RUN --mount=type=bind,source=.,target=/go/src/github.com/pyar6329/nclip,rw \
  set -x && \
  cd /go/src/github.com/pyar6329/nclip && \
  CC=/osxcross/target/bin/o64-clang make -e build-darwin && \
  mkdir -p /build && \
  tar -Jcf /build/nclip-Darwin.tar.xz build/darwin/nclip

FROM build-base AS linux-binary

RUN --mount=type=bind,source=.,target=/go/src/github.com/pyar6329/nclip,rw \
  set -x && \
  cd /go/src/github.com/pyar6329/nclip && \
  CC=/usr/bin/musl-gcc make -e build-linux && \
  mkdir -p /build && \
  tar -Jcf /build/nclip-Linux.tar.xz build/linux/nclip

FROM scratch AS export-stage
COPY --from=darwin-binary /build/nclip-Darwin.tar.xz /
COPY --from=linux-binary /build/nclip-Linux.tar.xz /
