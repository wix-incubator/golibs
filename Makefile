prepare:
	go mod download
	go mod tidy

format:
	gofmt -s -w .

lint:
	gofmt -d .

test:
	go test $(MAYBE_VERBOSE) -p 1 `go list ./...`

ci-steps: prepare lint test


