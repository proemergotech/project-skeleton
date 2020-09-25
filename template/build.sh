#!/bin/bash
set -e

go build -i -race -ldflags "-X gitlab.com/proemergotech/dliver-project-skeleton/app/config.AppVersion=$1" -o "$2"
chmod +x "$2"
