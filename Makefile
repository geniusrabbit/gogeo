GO           ?= go
PKG          := ./...
COVERPROFILE := coverage.out
DATA_SRC     := data/countries.json
DATA_OUT     := country_data.go
GEN          := ./internal/cmd/gencountries

.PHONY: help test lint fmt vet build-data generate bench cover tidy clean

help:
	@echo "Targets:"
	@echo "  test        Run unit tests"
	@echo "  lint        gofmt check + go vet"
	@echo "  fmt         Format Go sources"
	@echo "  build-data  Regenerate $(DATA_OUT) from $(DATA_SRC)"
	@echo "  generate    Alias for build-data (go generate)"
	@echo "  bench       Run benchmarks"
	@echo "  cover       Tests with coverage report"
	@echo "  tidy        go mod tidy"
	@echo "  clean       Remove coverage artifacts"

test:
	$(GO) test -count=1 $(PKG)

lint: vet
	@unformatted=$$(gofmt -l .); \
	if [ -n "$$unformatted" ]; then \
		echo "gofmt needed:"; \
		echo "$$unformatted"; \
		exit 1; \
	fi

vet:
	$(GO) vet $(PKG)

fmt:
	gofmt -w .

build-data:
	$(GO) run $(GEN) -src $(DATA_SRC) -out $(DATA_OUT)

generate:
	$(GO) generate .

bench:
	$(GO) test -bench=. -benchmem -count=1 -run='^$$' $(PKG)

cover:
	$(GO) test -count=1 -covermode=count -coverprofile=$(COVERPROFILE) $(PKG)
	$(GO) tool cover -func=$(COVERPROFILE)

tidy:
	$(GO) mod tidy

clean:
	rm -f $(COVERPROFILE)
