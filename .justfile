help:
    just --list

build:
    go build -o build/journal cmd/journal/main.go

test:
  go test ./... -v
