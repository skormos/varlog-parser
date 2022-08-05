clean:
	rm -rf ./varlogd
.PHONY: clean

compile:
	go build -o varlogd ./cmd/varlog
.PHONY: compile

run:
	go run ./cmd/varlog
.PHONY: run

test:
	go test -count=1 ./...
.PHONY: run
