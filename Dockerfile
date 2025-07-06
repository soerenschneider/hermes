FROM golang:1.24.4 AS build

WORKDIR /src
COPY ./go.mod ./go.sum ./
ENV GOOS=linux GOARCH=amd64 CGO_ENABLED=1
RUN go mod download

COPY ./ ./
RUN make build


FROM debian:12.9-slim AS final

LABEL maintainer="soerenschneider"
RUN useradd -m -s /bin/bash hermes
USER hermes
COPY --from=build --chown=hermes:hermes /src/hermes /hermes

ENTRYPOINT ["/hermes"]
