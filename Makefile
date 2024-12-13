SHELL := pwsh
.SHELLFLAGS := -Command

#### IMPORT ENV
include .env

run:
	go run ./cmd/app/main.go

lint:
	echo "Starting linters"
	golangci-lint run ./...


url=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

migrate_up:
	migrate -database "$(url)" -path migrations up 1

migrate_down:
	migrate -database "$(url)" -path migrations down 1


# ==============================================================================
# Modules support

deps-reset:
	git checkout -- go.mod
	go mod tidy

tidy:
	go mod tidy

deps-upgrade:
	# go get $(go list -f '{{if not (or .Main .Indirect)}}{{.Path}}{{end}}' -m all)
	go get -u -t -d -v ./...
	go mod tidy
	go mod vendor

deps-cleancache:
	go clean -modcache

