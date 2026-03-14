ARG VENDOR="vendor"
ARG BASE_IMAGE_VERSION="latest"

FROM ${VENDOR}/golang:${BASE_IMAGE_VERSION} AS build

ARG VENDOR="vendor"
ARG NAME="app"
ARG VERSION="latest"
ARG BUILD_TIME=""

WORKDIR /tmp/build

COPY "cmd" "cmd"
COPY "internal" "internal"
COPY "pkg" "pkg"
COPY "go.*" "."
COPY "*.go" "."

RUN mkdir -p bin \
    && go build \
  -ldflags "\
-X 'main.name=${NAME}' \
-X 'main.version=${VERSION}' \
-X 'main.copyright=${VENDOR}' \
-X 'main.authorName=${VENDOR}' \
-X 'main.buildTime=${BUILD_TIME}'\
" -o bin/app . \
    && go clean -cache -modcache -testcache

FROM ${VENDOR}/alpine:${BASE_IMAGE_VERSION} AS final

WORKDIR /opt/app/bin

ENV PATH="/opt/app/bin:${PATH}"

COPY --from=build /tmp/build/bin/app /opt/app/bin/app

COPY entrypoint.sh /

RUN ["chmod", "+x", "/entrypoint.sh"]

ENTRYPOINT ["/entrypoint.sh"]
