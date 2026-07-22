#Build stage
FROM golang:1.26.5-alpine3.24 AS builder
WORKDIR /app

COPY . .
RUN CGO_ENABLED=0 go build -o main main.go
RUN wget -q -O /tmp/migrate.tar.gz https://github.com/golang-migrate/migrate/releases/latest/download/migrate.linux-amd64.tar.gz && \
  tar -xzf /tmp/migrate.tar.gz -C /tmp && \
  mv /tmp/migrate /app/migrate && \
  chmod +x /app/migrate

#Run stage
FROM alpine:3 
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/app.env .
COPY --from=builder /app/migrate ./migrate
COPY --chmod=755 start.sh .
COPY --chmod=755 wait-for.sh .
COPY db/migration ./migration

EXPOSE 8080
ENTRYPOINT ["/app/start.sh"]
