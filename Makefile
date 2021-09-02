.PHONY: test
test:
	go test -race -coverprofile=cover.out -covermode=atomic ./...

.PHONY: build
build:
	- docker build --rm --file=./docker/toychain.Dockerfile --tag jasonsoft/toychain:latest .

run:
	- docker-compose up

build-run: build run

go-run:
	- go build cmd/main.go && ./main node


lint:
	docker run --rm -v ${LOCAL_WORKSPACE_FOLDER}:/app -w /app golangci/golangci-lint:v1.41-alpine golangci-lint run ./... -v



