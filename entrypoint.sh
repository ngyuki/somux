#!/bin/sh

case "$1" in
  -*) set -- somux "$@"
esac

exec "$@"
