FROM golang:1.23.3 AS build

WORKDIR /src
COPY ./go.mod ./go.sum ./
ENV GOOS=linux GOARCH=amd64 CGO_ENABLED=1
RUN go mod download

COPY ./ ./
RUN make build


FROM debian:12.8-slim AS final

LABEL maintainer="soerenschneider"
USER nonroot:nonroot
COPY --from=build --chown=nonroot:nonroot /src/hermes /hermes

ENTRYPOINT ["/hermes"]
