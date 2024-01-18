# base go image 
FROM golang:1.21-alpine as builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 go build -o loggerServiceApp ./cmd/api

RUN chmod +x /app/loggerServiceApp

# build a tiny docker image
FROM alpine:latest

RUN mkdir /app
RUN mkdir /templates

COPY --from=builder /app/loggerServiceApp /app
COPY --from=builder /app/templates/. /templates

CMD ["/app/loggerServiceApp"]
