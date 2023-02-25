.PHONY: clean

bin/gocooking: bin *.go database/*.go cmd/gocooking/*.go go.mod go.sum
	go build -o bin/gocooking cmd/gocooking/*.go

bin:
	mkdir -p bin

clean:
	rm -f bin/gocooking
