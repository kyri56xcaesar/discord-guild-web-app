# Variables
APP_NAME := cmd/api/discordwebapp
APP_SRC := cmd/api/main.go 

M_NAME := cmd/minioth/minioth
M_SRC := cmd/minioth/main.go

.PHONY: all
all: app mini

.PHONY: app
app: 
	go build -o ${APP_NAME} ${APP_SRC}
	./${APP_NAME}

.PHONY: mini
mini: 
	go build -o ${M_NAME} ${M_SRC}
	./${M_NAME}

# Clean up binaries and other generated files
.PHONY: clean
clean:
	rm -f ${APP_NAME} ${M_NAME}
