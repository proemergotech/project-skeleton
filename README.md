DLiveR Project Skeleton
==============================

## Table of contents
- [Overview](#overview)
- [Usage](#usage)
  - [Basics](#basics)
  - [How to remove unnecessary modules](#how-to-remove-unnecessary-modules)
    - [REDIS](#redis)
    - [Centrifugo](#centrifugo)
    - [GEB](#geb)
    - [Event Server](#event-server)
- [Example References](#example-references)
- [Development](#development)
  - [Dependencies](#dependencies)
  - [Testing](#testing)

## Overview

Project skeleton based on the latest standards.
Provides a basic golang based rest service with redis, centrifugo and geb (queue like RabbitMQ) handling.

## Usage
### Basics
Run `create_skeleton.sh` with the proper argument(s).   
(Don't forget to run `goimports` after file modifications.)

### Gentleman middlewares
Chain them after the `newGentleman(...)` function call upon client initialization in `container.go`
#### Retry middleware
The usage of this is not necessary. Useful for business-critical calls (for example transaction-related ones.)
It's using a custom exponential backoff solution to retry http requests. Overall and per-request timeout
must be set, with the overall being the longer time.

### Yafuds
The client's tracer must be set BEFORE calling `yafuds.New(...)`. If there is no tracing in the service,
remove the line `yafuds.SetTracer(opentracing.GlobalTracer())` from `container.newYafuds(...)`

### How to remove unnecessary modules
#### REDIS
- Remove file(s):
  - [redis](./app/storage/redis.go)
- Remove related code from:
  - [.env.example](./.env.example)
  - [config](./app/config/config.go)
  - [container](./app/di/container.go)
  - [error](./app/storage/error.go) 
  - [error definitions](./app/schema/service/error.go)
  - [Gopkg.toml](./Gopkg.toml)
  
#### Centrifugo
- Remove related code from:
  - [.env.example](./.env.example)
  - [config](./app/config/config.go)
  - [container](./app/di/container.go)
  - [service](./app/service/service.go) 
  - [error](./app/service/error.go)
  - [error definitions](./app/schema/service/error.go)
  - [Gopkg.toml](./Gopkg.toml)

#### Gentleman
- Remove package(s):
  - [client](./app/client)
- Remove related code from:
  - [container](./app/di/container.go)
  - [Gopkg.toml](./Gopkg.toml)
  
#### GEB
- Remove related code from:
  - [.env.example](./.env.example)
  - [config](./app/config/config.go)
  - [container](./app/di/container.go)
  - [Gopkg.toml](./Gopkg.toml)
- Remove [Event Server](#event-server)

#### Event Server
- Remove package(s):
  - [event](./app/event/)
- Remove usage from:
  - [container](./app/di/container.go)
  - [root command](./cmd/root.go)

#### Yafuds
- Remove related code from:
  - [.env.example](./.env.example)
  - [config](./app/config/config.go)
  - [container](./app/di/container.go)
  - [service](./app/service/service.go)
  - [error](./app/service/error.go)
  - [error definitions](./app/schema/service/error.go)
  - [Gopkg.toml](./Gopkg.toml)

#### Elastic
- Remove file(s):
  - [elastic](./app/storage/elastic.go)
- Remove related code from:
  - [.env.example](./.env.example)
  - [config](./app/config/config.go)
  - [container](./app/di/container.go)
  - [error](./app/storage/error.go)
  - [error definitions](./app/schema/service/error.go)
  - [Gopkg.toml](./Gopkg.toml)

## Development

### Dependencies
- install go
- check out project to: $GOPATH/src/gitlab.com/proemergotech/dliver-project-skeleton
- add dependencies to project with "go mod tidy"

### Testing

For the go based tests run `go test -v ./...` command from the root directory.  
