#!/bin/bash
set -e

#%:{{ `
go build -i -race -gcflags "all=-N -l" -ldflags "-X gitlab.com/proemergotech/dliver-project-skeleton/app/config.AppVersion=$1" -o "$2"
#%: ` | replace "dliver-project-skeleton" .ProjectName | trim }}
chmod +x "$2"
