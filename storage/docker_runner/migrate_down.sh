#!/bin/bash
DIR="$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
migrate -database="postgres://xfiendx4life:123456@172.17.0.2:5432/shortener?sslmode=disable" -path $DIR/migrations down 1
