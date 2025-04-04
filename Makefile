# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
BINARY_NAME=neh
BUILD_DIR=build/bin

# Build the project
build: test
	mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)

# Clean the build files
clean:
	rm -f $(BINARY_NAME)
	rm -f $(BUILD_DIR)/$(BINARY_NAME)

test:
	$(GOCMD) test -v ./...

tidy:
	$(GOCMD) mod tidy

.PHONY: build clean test tidy
