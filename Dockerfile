FROM golang:1.19 AS builder

COPY . /app
WORKDIR /app
RUN go mod download
RUN go build -o ./build/searcher cmd/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/build/searcher ./
COPY --from=builder /app/swagger /app/

EXPOSE 8000

CMD ["./searcher"]