#!/bin/sh
#
# Common Usage:
#
#   $ ./pack.sh
#   $ ./pack.sh linux amd64
#
# To list all supported platforms, use:
#
#   $ go tool dist list
# 
GOOS=${1:-`go env GOOS`} GOARCH=${2:-`go env GOARCH`} go build -trimpath -x -o opsul${1:+_$1}${2:+_$2}${3}
