# Variables
APP_NAME := myapp
SRC := cmd/api/main.go 

# Default target: Build and run the application
.PHONY: all
all: build

# Build the application
.PHONY: build
build:
	go build -o $(APP_NAME) $(SRC)

# Run the application
.PHONY: run
run: build
	./$(APP_NAME)

# Clean up binaries and other generated files
.PHONY: clean
clean:
	rm -f $(APP_NAME)
