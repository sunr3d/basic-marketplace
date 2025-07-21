.PHONY: mocks test up down clean

test: mocks
	go test -v ./...

mocks:
	go install github.com/vektra/mockery/v2@v2.53.2
	go generate ./...

up:
	docker-compose up --build 

down:
	docker-compose down -v

clean:
	rm -rf mocks/