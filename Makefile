run:
	go run ./cmd/goutui

test:
	go test ./...

vet:
	go vet ./...

lint:
	golangci-lint run

build:
	go build -o goutui ./cmd/goutui

clean:
	rm -f goutui

.PHONY: run test vet lint build clean
