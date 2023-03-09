FROM golang:latest

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o ./bin/bot cmd/bot/main.go

CMD ["./bin/bot"]
