#Build stage
FROM golang:1.26.5-alpine3.24 AS builder
WORKDIR /app

COPY . .
RUN CGO_ENABLED=0 go build -o main main.go

#Run stage
FROM scratch 
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/app.env .

EXPOSE 8080
CMD [ "/app/main" ]
