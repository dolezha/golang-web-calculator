FROM golang:1.19

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=1 go build -o calc-service ./cmd/calc_service
RUN CGO_ENABLED=1 go build -o calc-agent ./cmd/agent

EXPOSE 8080

RUN mkdir -p /app/data

ENV DB_PATH=/app/data/calculator.db
ENV PORT=8080
ENV JWT_SECRET=docker-secret-key-change-in-production

CMD ["./calc-service"] 