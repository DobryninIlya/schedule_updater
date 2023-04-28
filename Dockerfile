# базовый образ для работы с Go
FROM golang:alpine

# установка git и библиотеки для работы с PostgreSQL
#RUN apk add --no-cache git \
#    && go install github.com/jackc/pgx@latest
#
#WORKDIR /go/src/app
#COPY . .

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main ./cmd/

RUN apk update && \
    apk add postgresql-client && \
    go get github.com/jackc/pgx/v4


# тесты и сборка микросервиса
#RUN go build

# установка базовой конфигурации
EXPOSE 5433 5432
EXPOSE 81 80
# запуск микросервиса
CMD ["./main"]