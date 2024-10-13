# Build stage
FROM golang:1.23-alpine as builder

WORKDIR /app

COPY . .

RUN go build -o main main.go
RUN apk add curl
# RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.18.1/migrate.linux-amd64.tar.gz | tar xvz || exit 1
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.18.1/migrate.linux-amd64.tar.gz | tar xvz && ls -la || exit 1

# Run stage
FROM alpine
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/migrate ./migrate

COPY app.env .
COPY start.sh .
COPY wait-for.sh .
COPY ./db/migration ./migration


EXPOSE 8080

CMD ["/app/main"]
ENTRYPOINT [ "/app/start.sh" ]
