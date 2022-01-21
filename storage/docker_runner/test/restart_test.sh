#!/bin/bash
set -e

DIR="$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"

echo "$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"

(docker stop postgres_test && docker rm postgres_test) || true
(docker stop postgres && docker rm postgres) || true

sudo rm -rf $DIR/_data_test
sudo docker run -d \
    -p 5433:5432 \
    --name postgres_test \
    -e POSTGRES_PASSWORD=123456 \
    -e PGDATA=/var/lib/postgresql/data/pgdata \
    -v "$DIR/_data_test":/var/lib/postgresql/data \
    -v "$DIR/init-db":/docker-entrypoint-initdb.d \
    postgres
