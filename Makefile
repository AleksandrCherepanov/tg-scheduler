build:
	go build -o ./bin/ ./cmd/
run:
	./bin/cmd
build-docker:
	docker-compose up -d
