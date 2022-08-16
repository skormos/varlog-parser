clean:
	rm -f ./varlogd
.PHONY: clean

compile: lint
	go build -o varlogd ./cmd/varlog
.PHONY: compile

run:
	go run ./cmd/varlog
.PHONY: run

lint:
	@goimports -local "github.com/skormos/varlog-parser" -w -l .
	@golangci-lint run --out-format github-actions ./...
.PHONY: lint

test: lint
	go test -count=1 ./...
.PHONY: run

_install-codegen:
	go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest
.PHONY: _install-codegen

generate-api: _install-codegen
	APIFILE=$(PWD)/api/rest/$(version)/varlog.openapi3.json go generate ./internal/api/rest/$(version)/...
.PHONY: generate
