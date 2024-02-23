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
