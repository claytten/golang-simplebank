# Build Stage
FROM golang:1.19.5-alpine3.17 as builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

# Run Stage 
FROM alpine:3.17
WORKDIR /app
COPY --from=builder /app/main .
COPY app.env .
COPY start.sh .
COPY wait-for.sh .
COPY migration ./migration

EXPOSE 8080
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]