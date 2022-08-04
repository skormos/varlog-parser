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
