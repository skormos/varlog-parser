# Removes artifacts generated as part of compilation or testing.
clean:
	rm -f ./varlogd
.PHONY: clean

# Attempts to compile the primary executable into a runnable binary.
compile: lint
	go build -o varlogd ./cmd/varlog
.PHONY: compile

# Attempts to run the primary main package without compiling.
run:
	go run ./cmd/varlog
.PHONY: run

# Calls clean code functions "go-imports" and "golangci-lint". Both of these need to be preinstelled before running this
# project.
lint:
	@goimports -local "github.com/skormos/varlog-parser" -w -l .
	@golangci-lint run --out-format github-actions ./...
.PHONY: lint

# Runs the tests in this project, but first clears the cache.
test: lint
	go test -count=1 ./...
.PHONY: run

# inner target to support generating endpoint code.
_install-codegen:
	go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest
.PHONY: _install-codegen

# Updates the generated API Endpoint code based on the API spec for the version provided. Version is specified using an
# environment variable. eg. `version=v1 make generate-api`
generate-api: _install-codegen
	APIFILE=$(PWD)/api/rest/$(version)/varlog.openapi3.json go generate ./internal/api/rest/$(version)/...
.PHONY: generate
