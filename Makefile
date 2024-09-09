APP_NAME=main

templ-generate:
	templ generate

build-dev:
	make templ-generate && go build -o tmp/${APP_NAME} ./main.go

build-prod:
	make templ-generate && go build -o tmp/${APP_NAME} -ldflags "-s -w" ./main.go

build-run:
	make build-prod && ./tmp/${APP_NAME}

run:
	make ./tmp/${APP_NAME}