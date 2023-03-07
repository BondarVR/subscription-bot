.PHONY:
.SILENT:

build:
	go build -o ./.bin/bot cmd/bot/main.go

run: build
	./.bin/bot

format:
	${call colored, formatting is running...}
	go vet ./...
	go fmt ./...

test:
	go test -v -count=1 ./...

cover:
	go test -short -count=1 -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
	rm coverage.out
