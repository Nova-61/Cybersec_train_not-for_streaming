.PHONY: build run test clean

build:
	go build -o proxy ./cmd/proxy

run:
	go run ./cmd/proxy -listen :8888 -target http://localhost:8080

run-limited:
	go run ./cmd/proxy -listen :8888 -target http://localhost:8080 -rate-limit 5

test:
	go test ./...

test-verbose:
	go test -v ./...

coverage:
	go test -cover ./...

clean:
	rm -f proxy
	go clean

install-tools:
	go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
	go install github.com/tsenart/vegeta/v12@latest

security:
	gosec ./...

bench:
	vegeta attack -duration=10s -rate=100 -targets=targets.txt | vegeta report
