.PHONY: test

test:
	go build -o freeze-test
	go test ./...
	rm freeze-test
