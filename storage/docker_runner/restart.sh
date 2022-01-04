#!/bin/bash
set -e

DIR="$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"

echo "$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"

(docker stop postgres && docker rm postgres) || true

sudo rm -rf $DIR/data

sudo docker run -d \
    -p 5432:5432 \
    --name postgres \
    -e POSTGRES_PASSWORD=123456 \
    -e PGDATA=/var/lib/postgresql/data/pgdata \
    -v "$DIR/data":/var/lib/postgresql/data \
    -v "$DIR/init-db":/docker-entrypoint-initdb.d \
    postgres
