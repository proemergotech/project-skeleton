#!/bin/bash
set -e

graceful_shutdown() {
  pid=$(pgrep dlv)
  child_pid=$(pgrep -P "$pid")
  
  kill "$child_pid"
  while kill -0 "$child_pid"; do 
    sleep 1
  done
  
  kill "$pid"
  while kill -0 "$pid"; do 
    sleep 1
  done
}

trap graceful_shutdown SIGTERM SIGINT 

dlv --listen=:2345 --headless=true --api-version=2 --accept-multiclient exec "$1" -- "${@:2}" &
wait $!