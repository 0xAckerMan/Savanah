pwd := "./cmd/api"

# Path: makefile
build:
	@echo "Building the application..."
	@go build -o $(pwd)/api $(pwd)/main.go

# Path: makefile
run:
	@echo "Running the application..."
	@go run $(pwd)/main.go

# Path: makefile
test:
	@echo "Running tests..."
	@go test -v ./...

# Path: makefile
clean:
	@echo "Cleaning up..."
	@rm -rf $(pwd)/api

# Path: makefile
docker-build:
	@echo "Building the docker image..."
	@docker build -t api .

# Path: makefile
docker-run:
	@echo "Running the docker container..."
	@docker run -p 8080:8080 api

# Path: makefile
docker-stop:
	@echo "Stopping the docker container..."
	@docker stop $(shell docker ps -q --filter ancestor=api)

# Path: makefile
docker-clean:
	@echo "Cleaning up..."
	@docker rmi api

# Path: makefile
docker-compose-up:
	@echo "Running the docker-compose..."
	@docker-compose up

# Path: makefile
docker-compose-down:
	@echo "Stopping the docker-compose..."
	@docker-compose down

# Path: makefile
docker-compose-clean:
	@echo "Cleaning up..."
	@docker-compose down
	@docker rmi api
```

## 3. Dockerfile

```dockerfile
# Path: Dockerfile
FROM golang:1.16-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /api ./cmd/api/main.go

FROM alpine:3.13

COPY --from=builder /api /api

EXPOSE 8080

CMD ["/api"]
```

## 4. docker-compose.yml

```yaml
# Path: docker-compose.yml
version: '3.8'

services:
  api:
	build: .
	ports:
	  - "8080:8080"
```

## 5. Run the application

```bash
# Build the application
make build

# Run the application
make run
```

## 6. Build the docker image

```bash
make docker-build
```

## 7. Run the docker container

```bash
make docker-run
```

## 8. Stop the docker container

```bash
make docker-stop
```

## 9. Clean up

```bash
make clean
```

## 10. Run the docker-compose

```bash
make docker-compose-up
```

## 11. Stop the docker-compose

```bash
make docker-compose-down
```

## 12. Clean up

```bash
make docker-compose-clean
```

## 13. Test the application

```bash
make test
```