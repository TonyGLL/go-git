APP_EXECUTABLE=main

build:
	GOARCH=amd64 GOOS=linux go build -o ${APP_EXECUTABLE} main.go

run: build
	./${APP_EXECUTABLE}

init: build
	./${APP_EXECUTABLE} init
