.PHONY: clean

bin/gocooking: bin *.go database/*.go cmd/gocooking/*.go go.mod go.sum
	CGO_ENABLED=1 go build -o bin/gocooking cmd/gocooking/*.go

bin:
	mkdir -p bin

clean:
	rm -rf bin

.PHONY: test
test:
	go test -cover -tags test -v ./...

test-coverage: bin
	go test -coverprofile=bin/coverage.out -tags test ./...
	go tool cover -html bin/coverage.out
