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
Find and replace the text `dliver-project-skeleton` - for example with notepad++ - with your new service's name (e.g.: `dliver-dummy-service`)

### How to remove unnecessary modules
#### REDIS
- Remove the `RedisUnavailable` function from the [apierr factory](./app/apierr/factory.go).
- Remove the `redis` folder from the [client](./app/client) folder.
- Remove the `newRedis` function from the [container](./app/di/container.go).
- Remove the `Init REDIS client` section from the [container's](./app/di/container.go) constructor.
- Remove all the related client definitions from the other definitions (e.g.: Core, Action, Container, etc...) where `redisClient` is used.
- Remove the `ErrRedisUnavailable` error code from the [error definitions](./app/schema/service/error.go) constructor.

#### Centrifugo
- Remove the `CentrifugeErrorDetail` type and all of it's functions from the [error details](./app/apierr/error_detail.go).
- Remove the `Centrifuge` and `CentrifugeResponse` functions from the [apierr factory](./app/apierr/factory.go).
- Remove the `centrifugo` folder from the [client](./app/client) folder.
- Remove the `newCentrifugeClient` function from the [container](./app/di/container.go).
- Remove the `Init Centrifuge Client` section from the [container's](./app/di/container.go) constructor.
- Remove all the related client definitions from the other definitions (e.g.: Core, Action, Container, etc...) where `centrifugeClient` is used.
- Remove the `ErrCentrifuge` error code from the [error definitions](./app/schema/service/error.go) constructor.

#### GEB
- Remove the `gebCloser` from the [Container's](./app/di/container.go) definition.
- Remove the `Init GEB queue` section from the [container's](./app/di/container.go) constructor.
- Remove the `newGebQueue` function from the [container](./app/di/container.go).
- Remove the `gebCloser` section from the [container's](./app/di/container.go) Close function.
- Remove all the related client definitions from the other types (e.g.: Core, Action, Container, etc...) where `gebQueue` or `gebCloser` is used.

#### Event Server
- Remove [GEB](#geb)
- Remove the `event` folder from the [app](./app) folder.
- Remove the `Init EVENT server` section from the [container's](./app/di/container.go) constructor.

## Example References
* [API documentation](./API.md)

## Development

### Dependencies
- install go
- check out project to: $GOPATH/src/gitlab.com/proemergotech/dliver-project-skeleton
- install dep
- add dependencies to project with "dep ensure"

### Testing

For the go based tests run `go test -v ./...` command from the root directory.  
