# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
BINARY_NAME=neh

# Build the project
build:
	$(GOBUILD) -o $(BINARY_NAME)

# Clean the build files
clean:
	rm -f $(BINARY_NAME)

test:
	$(GOCMD) test -v ./...

tidy:
	$(GOCMD) mod tidy

.PHONY: build clean test tidy
