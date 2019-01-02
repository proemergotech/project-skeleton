#!/usr/bin/env bash
set -e

if [[ "$1" = 'run' ]]; then
    chmod +x ./build.sh

    exec CompileDaemon \
    -build="./build.sh dev /tmp/dliver-project-skeleton" \
    -command="/tmp/dliver-project-skeleton ${@:2}" \
    -exclude-dir=vendor \
    -graceful-kill=true \
    -log-prefix=false
fi

exec "$@"
