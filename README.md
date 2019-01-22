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
Find and replace the text `dliver_project_skeleton` - for example with notepad++ - with your new service's name (e.g.: `dliver_dummy_service`)    
(Don't forget to run `goimports` after file modifications.)

### How to remove unnecessary modules
#### REDIS
- Remove packages:
  - [redis](./app/client/redis)
- Remove related code from:
  - [.env.example](./.env.example)
  - [config](./app/config/config.go)
  - [container](./app/di/container.go)
  - [apierr factory](./app/apierr/factory.go) 
  - [error definitions](./app/schema/service/error.go)
  - [Gopkg.toml](./Gopkg.toml)
  
#### Centrifugo
- Remove packages:
  - [client](./app/client/centrifugo)
- Remove related code from:
  - [.env.example](./.env.example)
  - [config](./app/config/config.go)
  - [container](./app/di/container.go)
  - [service](./app/service/service.go)
  - [apierr details](./app/apierr/error_detail.go) 
  - [apierr factory](./app/apierr/factory.go) 
  - [error definitions](./app/schema/service/error.go)

#### Gentleman
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
- Remove files:
  - [error](./app/client/error.go)
- Remove usage from:
  - [container](./app/di/container.go)
  - [root command](./cmd/root.go)

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
