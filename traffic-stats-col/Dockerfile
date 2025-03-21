FROM golang:latest AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

RUN go build -o traffic-stats-col .

FROM golang:latest

WORKDIR /app

COPY --from=builder /app/traffic-stats-col /app/traffic-stats-col

COPY config.yaml /app/config.yaml

RUN chmod +x /app/traffic-stats-col

CMD ["/app/traffic-stats-col"]
