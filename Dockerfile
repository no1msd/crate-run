# stage for exporting the godot project to html
FROM alpine:3.12.0 AS gamebuilder

ENV STORAGE_PATH=https://storage.sipkasoft.hu
ENV GODOT_VERSION=3.3.4.stable

RUN apk add libexecinfo wget && \
    sh -c "wget $STORAGE_PATH/godot/$GODOT_VERSION/godot_headless.musl.$(uname -m) -O /usr/bin/godot" && \
    chmod a+x /usr/bin/godot && \
    mkdir -p ~/.local/share/godot/templates/$GODOT_VERSION && \
    wget $STORAGE_PATH/godot/$GODOT_VERSION/webassembly_release.zip -P ~/.local/share/godot/templates/$GODOT_VERSION

COPY game /src
WORKDIR /src
RUN mkdir -p /build && \
    godot --export "HTML5" /build/crate.html && \
    mv /build/crate.html /build/index.html

#stage for building the api
FROM golang:1.17-alpine AS apibuilder

COPY api /src
WORKDIR /src
RUN apk add gcc musl-dev && \
    go test ./... && \
    go build && \
    strip crateapi

# final stage
FROM alpine:3.12.0

COPY --from=gamebuilder /build /content
COPY --from=apibuilder /src/crateapi /crate-api
COPY web/* /content/

EXPOSE 9080

WORKDIR /data

CMD [ "/crate-api" ]
