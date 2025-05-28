.PHONY: clean test

test:
	go test ./...

# clean all build result
clean:
	go clean ./...