.PHONY: build
build:
	go build -o Q-n-A ./*.go

.PHONY: run
run:
	go run ./*.go
