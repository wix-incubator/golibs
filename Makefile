prepare:
	mkdir -p vendor/x/y
	go mod download

format:
	gofmt -s -w .

lint:
	gofmt -d .

test:
	go test $(MAYBE_VERBOSE) `go list ./...`

ci-steps: prepare lint test


