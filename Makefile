.PHONY: test

test:
	go test ./...

golden:
	cp -r test/output/* test/golden
