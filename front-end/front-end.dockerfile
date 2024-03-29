# base go image 
FROM golang:1.21-alpine as builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 go build -o frontEnd ./cmd/web

RUN chmod +x /app/frontEnd

# build a tiny docker image
FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/frontEnd /app

# Run the server executable
CMD [ "/app/frontEnd" ]
