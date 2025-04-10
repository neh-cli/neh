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

neh-decache: build
	@echo "Running neh command..."
	./build/bin/neh decache

.PHONY: build clean test tidy
