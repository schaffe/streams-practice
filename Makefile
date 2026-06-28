BINDIR ?= out

.PHONY: build test lint fmt clean coverage all

build:
	@mkdir -p $(BINDIR)
	for dir in $$(go list -f '{{if eq .Name "main"}}{{.Dir}}{{end}}' ./...); do \
		go build -o $(BINDIR)/$$(basename $$dir) $$dir; \
	done

test:
	go test -race -cover ./...

lint:
	golangci-lint run ./...

fmt:
	gofmt -s -w .

clean:
	rm -f coverage.out
	rm -rf $(BINDIR)

coverage:
	go test -coverprofile=coverage.out ./... && go tool cover -func=coverage.out

all: fmt lint build test
