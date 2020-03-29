# syntax=docker/dockerfile:1.0-experimental
FROM dockercore/golang-cross:1.13.8

RUN set -x && \
  apt-get update && \
  apt-get install -y \
    musl-dev \
    musl-tools && \
  apt-get clean && \
  rm -rf /var/lib/apt/lists/*

RUN --mount=type=bind,source=.,target=/go/src/github.com/pyar6329/nclip,rw \
  set -x && \
  cd /go/src/github.com/pyar6329/nclip && \
  make build && \
  cp -rf build /build
