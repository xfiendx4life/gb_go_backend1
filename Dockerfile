FROM golang:1.17 AS build
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY ./ ./
# RUN go run cmd/shrtener/main.go
RUN go build -o ./out cmd/shrtener/main.go

# DEPLOY
FROM gcr.io/distroless/base-debian10
WORKDIR /
COPY --from=build /app/out /app
COPY --from=build /app/configs/ /configs
EXPOSE 8000
USER nonroot:nonroot
ENTRYPOINT [ "/app" ]