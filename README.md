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
  - [Debug](#debug)

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

### Debug

To remote debug a service we need 2 things.  

- Add some extra option to full-compose when overriding, e.g.:
```
  profile-service-go:
    image: profile-service-go:dev
    build: ~/dev/dliver-profile-service-go/dev
    command: ["/usr/local/bin/entrypoint.sh", "debug"]
    environment: *goproxy
    volumes:
      - ~/dev/dliver-profile-service-go:/go/src/gitlab.com/proemergotech/dliver-profile-service-go
      - ~/gopkg:/go/pkg
    ports:
      - 2345:2345
    security_opt:
      - seccomp:unconfined
```
`security_opt` is necessary to run delve inside container and the extra port mapping to access from outside.  
In `command` call `debug` instead of `run`.  
The container will build the service, then wait for the debugger to attach.  
- In Goland create a Go remote debug configuration, set the full compose machine ip/domain and the port from the previous mapping.