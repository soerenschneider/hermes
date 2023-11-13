FROM golang:1.21.4 AS build

WORKDIR /src
COPY ./go.mod ./go.sum ./
ENV GOOS=linux GOARCH=amd64 CGO_ENABLED=0
RUN go mod download

COPY ./ ./
RUN make build


FROM gcr.io/distroless/static AS final

LABEL maintainer="soerenschneider"
USER nonroot:nonroot
COPY --from=build --chown=nonroot:nonroot /src/hermes /hermes

ENTRYPOINT ["/hermes"]
