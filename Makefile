.PHONY: test

test:
	go test ./...

golden:
	cp test/output/* test/golden
