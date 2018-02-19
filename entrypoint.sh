#!/usr/bin/env bash

set -e

mkdir -p /etc/

PORT=${PORT:-8080}

echo '
{
    "port": '$PORT'
}
' > /etc/server.conf.json

date '+%Y/%m/%d %H:%M:%S Configuration file created'

cat /etc/server.conf.json

exec "$@"
