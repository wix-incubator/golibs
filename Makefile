prepare:
	CGO_ENABLED=0 go mod download

format:
	CGO_ENABLED=0 gofmt -s -w .

lint:
	CGO_ENABLED=0 gofmt -d .

test:
	CGO_ENABLED=0 go test $(MAYBE_VERBOSE) -p 1 `go list ./...`

ci-steps: prepare lint test


