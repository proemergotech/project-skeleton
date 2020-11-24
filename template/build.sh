#!/bin/bash
set -e

#%:{{ `
go build -i -race -ldflags "-X github.com/proemergotech/project-skeleton/app/config.AppVersion=$1" -o "$2"
#%: ` | replace "project-skeleton" .ProjectName | trim }}
chmod +x "$2"
