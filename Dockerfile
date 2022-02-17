FROM golang:1.17 AS build
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY ./ ./
RUN go build -o ./out cmd/shrtener/main.go

# DEPLOY
FROM centos:7
WORKDIR /
COPY --from=build /app/out /app
COPY --from=build /app/configs/ /configs
COPY --from=build /app/web/templates /web/templates

ENTRYPOINT [ "/app", "-config", "/configs/config.yml"]
