swagger:
	oapi-codegen --config=swagger-ui/oapi-codegen.yaml swagger-ui/openapi.yaml

cover:
	go test -short -count=1 -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
	rm coverage.out

gen:
	mockgen -source=internal/storage/cache/storage.go \
	-destination=internal/storage/cache/mocks/mock_storage.go

project/init:
	docker compose up -d --build

docker/dev:
	docker compose up -d

docker/prod:
	docker compose -f docker-compose.yml up --build

docker/down:
	docker compose down -v