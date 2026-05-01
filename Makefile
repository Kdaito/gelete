BINARY := gelete
GOPATH := $(shell go env GOPATH)
LINT := $(GOPATH)/bin/golangci-lint

.PHONY: all test lint fmt vet build ci clean

# Run all CI checks (matches GitHub Actions pipeline)
ci: fmt vet lint test build

test:
	go test ./... -v -coverprofile=coverage.txt -covermode=atomic

lint: $(LINT)
	$(LINT) run --timeout=5m

fmt:
	@unformatted=$$(gofmt -s -l .); \
	if [ -n "$$unformatted" ]; then \
		echo "The following files are not formatted:"; \
		echo "$$unformatted"; \
		exit 1; \
	fi

vet:
	go vet ./...

build:
	go build -v -o $(BINARY) .

# Install golangci-lint if missing
$(LINT):
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(GOPATH)/bin

clean:
	rm -f $(BINARY) coverage.txt
