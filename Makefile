.PHONY: build test run clean

build:
	go build -o bin/anagram-counter cmd/counter/main.go

test:
	go test ./... -v

run: build
	./bin/anagram-counter --dir=./testdata

clean:
	rm -rf bin/