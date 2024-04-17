# stage for exporting the godot project to html
FROM barichello/godot-ci:3.3.4 AS gamebuilder

COPY game /src
WORKDIR /src
RUN mkdir -p /build && \
    godot --export "HTML5" /build/crate.html && \
    mv /build/crate.html /build/index.html

# stage for building the api
FROM golang:1.17-alpine AS apibuilder

COPY api /src
WORKDIR /src
RUN apk add gcc musl-dev && \
    go test ./... && \
    go build && \
    strip crate-run-api

# final stage
FROM alpine:3.12.0

COPY --from=gamebuilder /build /content
COPY --from=apibuilder /src/crate-run-api /crate-run-api
COPY web/* /content/

EXPOSE 9080

WORKDIR /data

CMD [ "/crate-run-api" ]
