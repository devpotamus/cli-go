# Base makefile template for go operations

all:
	# Default build instructions
	go build
deps:
	# Dependencies go here
	go get golang.org/x/lint/golint
clean:
	# Clean instructions
build:
	# Go build instructions
	go build
install:
	# Go install instructions
	go install
lint:
	# Go linting instructions
	$(shell go list -f {{.Target}} golang.org/x/lint/golint ./...)