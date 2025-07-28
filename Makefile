APP_NAME ?= app

.PHONY: vet
vet:
	go vet ./...

.PHONY: staticcheck
staticcheck:
	staticcheck ./...

.PHONY: test
test:
	go test -race -v -timeout 30s ./...

.PHONY: tailwind-watch
tailwind-watch:
	./tailwindcss -i ./assets/css/input.css -o ./assets/css/out/style.css --watch

.PHONY: tailwind-build
tailwind-build:
	./tailwindcss -i ./assets/css/input.css -o ./assets/css/out/style.min.css --minify

.PHONY: templ-watch
templ-watch:
	templ generate --watch

.PHONY: templ-generate
templ-generate:
	templ generate

.PHONY: dev-serve
dev-serve:
	go build -o ./tmp/main ./cmd/main.go && air

.PHONY: dev
dev:
	make tailwind-watch & make templ-watch & make dev-serve

.PHONY: build
build:
	make tailwind-build
	make templ-generate
	go build -ldflags "-X main.Environment=production" -o ./bin/$(APP_NAME) ./cmd/main.go

.PHONY: build-migrate
build-migrate:
	go build -o ./bin/migrate ./db/migrate

.PHONY: docker-build
docker-build:
	docker-compose -f ./dev/docker-compose.yml build

.PHONY: docker-up
docker-up:
	docker-compose -f ./dev/docker-compose.yml up

.PHONY: docker-dev
docker-dev:
	docker-compose -f ./dev/docker-compose.dev.yml up

.PHONY: docker-down
docker-down:
	docker-compose -f ./dev/docker-compose.yml down

.PHONY: docker-clean
docker-clean:
	docker-compose -f ./dev/docker-compose.yml down -v --rmi all
