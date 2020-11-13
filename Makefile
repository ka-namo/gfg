init:
	bash build/install_server.sh

up:
	docker-compose -f build/docker-compose.yml up

test-unit:
	@echo ">> Running unit tests"
	@go test -short -coverprofile=unit.coverprofile -covermode=atomic -race ./...
