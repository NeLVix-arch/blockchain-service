BINARY = service
SERVICE = blockchain_service
PORT = 80
POSTGRES_PORT = 5432

build:
	CGO_ENABLED=0 GOOS=linux go build -o $(BINARY) ./src/main.go 

run:
	./$(BINARY)

clean:
	rm -f $(BINARY)

test:
	go test -v ./...

docker-build:
	docker build -t $(SERVICE) .

docker-run:
	docker run -p $(PORT):$(PORT) -p $(POSTGRES_PORT):$(POSTGRES_PORT) --network="host" $(SERVICE)

docker-push:
	docker push $(SERVICE)
