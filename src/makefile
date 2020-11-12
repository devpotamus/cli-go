all: deps clean build fmt lint
deps:
	go get golang.org/x/lint/golint
clean:
	rm -f cli-go
build:
	go build
fmt:
	go fmt ./...
lint:
	$(shell go list -f {{.Target}} golang.org/x/lint/golint ./...)