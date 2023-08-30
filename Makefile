
compose-up:
	docker-compose up --build -d && docker-compose logs -f
.PHONY: compose-up

compose-down:
	docker-compose down --remove-orphans
.PHONY: compose-down

test:
	go test -v ./...
	rm -Rf internal/segment/delivery/http/logs
	rm -Rf internal/segment/useCase/logs

cover:
	go test -short -count=1 -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
	rm coverage.out
.PHONY: cover