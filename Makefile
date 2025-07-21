.PHONY: mocks test up down

test: mocks
	go test -v ./...

mocks:
	go install github.com/vektra/mockery/v2@v2.53.2
	go generate ./...

up:
	docker-compose up --build 

down:
	docker-compose -v down