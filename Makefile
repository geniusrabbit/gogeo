GO           ?= go
PKG          := ./...
COVERPROFILE := coverage.out
DATA_SRC       := data/countries.json
DATA_OUT       := country_data.go
REGION_DEBIAN  := data/iso_3166-2.json
REGION_OODAVID := data/oodavid_iso_3166_2.js
REGION_SRC     := data/regions.json
REGION_OUT     := region_data.go
GEN_COUNTRIES  := ./internal/cmd/gencountries
GEN_MERGE      := ./internal/cmd/mergerregions
GEN_REGIONS    := ./internal/cmd/genregions

.PHONY: help test lint fmt vet merge-data build-data generate export-data bench cover tidy clean

help:
	@echo "Targets:"
	@echo "  test        Run unit tests"
	@echo "  lint        go vet + golangci-lint"
	@echo "  fmt         Format Go sources"
	@echo "  merge-data  Merge debian + oodavid into data/regions.json"
	@echo "  build-data  Merge sources and regenerate Go tables"
	@echo "  generate    Alias for build-data (go generate)"
	@echo "  export-data Dump countries or regions (KIND=countries|regions FORMAT=json|csv)"
	@echo "  bench       Run benchmarks"
	@echo "  cover       Tests with coverage report"
	@echo "  tidy        go mod tidy"
	@echo "  clean       Remove coverage artifacts"

test:
	$(GO) test -count=1 $(PKG)

lint: vet
	golangci-lint run -v ./...

vet:
	$(GO) vet $(PKG)

fmt:
	gofmt -w .

merge-data:
	$(GO) run $(GEN_MERGE) -debian $(REGION_DEBIAN) -oodavid $(REGION_OODAVID) -out $(REGION_SRC)

build-data: merge-data
	$(GO) run $(GEN_COUNTRIES) -src $(DATA_SRC) -out $(DATA_OUT)
	$(GO) run $(GEN_REGIONS) -src $(REGION_SRC) -countries $(DATA_SRC) -out $(REGION_OUT)

generate:
	$(GO) generate .

KIND   ?= countries
FORMAT ?= json

export-data:
	$(GO) run ./internal/cmd/export -kind $(KIND) -format $(FORMAT)

bench:
	$(GO) test -bench=. -benchmem -count=1 -run='^$$' $(PKG)

cover:
	$(GO) test -count=1 -covermode=count -coverprofile=$(COVERPROFILE) $(PKG)
	$(GO) tool cover -func=$(COVERPROFILE)

tidy:
	$(GO) mod tidy

clean:
	rm -f $(COVERPROFILE)
