# hermes
[![Go Report Card](https://goreportcard.com/badge/github.com/soerenschneider/hermes)](https://goreportcard.com/report/github.com/soerenschneider/hermes)
![test-workflow](https://github.com/soerenschneider/hermes/actions/workflows/test.yaml/badge.svg)
![release-workflow](https://github.com/soerenschneider/hermes/actions/workflows/release-container.yaml/badge.svg)
![golangci-lint-workflow](https://github.com/soerenschneider/hermes/actions/workflows/golangci-lint.yaml/badge.svg)

Accepts notifications via multiple adapters and routes notifications to upstream services

## Features

📣 Accepts notifications via event sources and routes them to various notification systems<br/>
🏰 Built-in resiliency for different failure scenarios<br/>
🔭 Observability through Prometheus metrics

## Why would I need it?

📌 You have some appliances that are only able to send notifications via SMTP<br/>

## Installation

### Docker / Podman
````shell
$ docker pull ghcr.io/soerenschneider/hermes:main
````

### Binaries
Head over to the [prebuilt binaries](https://github.com/soerenschneider/hermes/releases) and download the correct binary for your system.

### From Source
As a prerequisite, you need to have [Golang SDK](https://go.dev/dl/) installed. After that, you can install hermes from source by invoking:
```text
$ go install github.com/soerenschneider/hermes@latest
```

## Configuration
Head over to the [configuration section](docs/configuration.md) to see more details.


## Observability
Head over to the [metrics](docs/metrics.md) to see more details.

## Changelog
The changelog can be found [here](CHANGELOG.md)
