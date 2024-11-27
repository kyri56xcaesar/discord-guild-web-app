# Variables
APP_NAME := discordwebapp
SRC_PATH := cmd/api/
SRC := main.go 

# Default target: Build and run the application
.PHONY: all
all: build run

# Build the application
.PHONY: build
build:
	go build -o $(SRC_PATH)$(APP_NAME) $(SRC_PATH)$(SRC)

# Run the application
.PHONY: run
run: build
	./$(SRC_PATH)$(APP_NAME)

.PHONY: mini
mini: 
	go build -o cmd/minioth/minioth cmd/minioth/main.go
	./cmd/minioth/minioth

# Clean up binaries and other generated files
.PHONY: clean
clean:
	rm -f $(SRC_PATH)$(APP_NAME)
