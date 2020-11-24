#!/usr/bin/env bash
set -e

mkdir -p /usr/local/.cache/go-build

build_cmd=""
run_cmd=""
#%:{{ `
output="/tmp/project-skeleton"
#%: ` | replace "project-skeleton" .ProjectName | trim }}

if [ "$1" = 'run' ]; then
  build_cmd="./build.sh"
  run_cmd="$output"
elif [ "$1" = 'debug' ]; then
  build_cmd="./dev/debug_build.sh"
  run_cmd="./dev/debug_run.sh $output"
fi

if [[ -n "$build_cmd" ]]; then
  chmod +x "$build_cmd"

  exec CompileDaemon \
  -build="$build_cmd dev $output" \
  -command="$run_cmd ${@:2}" \
  -exclude-dir=vendor \
  -graceful-kill=true \
  -log-prefix=false
fi

exec "$@"
