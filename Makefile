APP_EXECUTABLE=main

build:
	GOARCH=amd64 GOOS=linux go build -o ${APP_EXECUTABLE} main.go

run: build
	./${APP_EXECUTABLE}

init: build
	./${APP_EXECUTABLE} init

add: build
	./${APP_EXECUTABLE} add .

commit: build
	./${APP_EXECUTABLE} commit -m "first commit"

log: build
	./${APP_EXECUTABLE} log

status: build
	./${APP_EXECUTABLE} status