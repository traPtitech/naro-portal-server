.PHONY: build
build:
	go build -o Q-n-A ./*.go

.PHONY: run
run: build
	./Q-n-A
