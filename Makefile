GOLANGCI_VERSION = v1.62.0
GOLANGCI_BIN = $(shell go env GOPATH)/bin/golangci-lint

.PHONY: ensure-golangci-is-installed
ensure-golangci-is-installed:
	@which $(GOLANGCI_BIN) > /dev/null || {\
		echo "Installing golangci-lint $(GOLANGCI_VERSION)"; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin $(GOLANGCI_VERSION); \
	}

lint: ensure-golangci-is-installed
	$(GOLANGCI_BIN) run ./...



run:
	air

test:
	go test -v ./...
