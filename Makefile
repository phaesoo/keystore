# Meta
NAME := shield

# Install dependencies
.PHONY: deps
deps:
	go mod download

# Build the main executable
main:
	go build -o main .

# This is a specialized build for running the executable inside a minimal scratch container
.PHONY: build-docker
build-docker:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -a -installsuffix cgo -o ./main .

# Run all unit tests
.PHONY: test
test: main
	go test -short ./...

# Migrate database
.PHONY: migrate
migrate:
	go run db/migrate/main.go

# Create a new empty migration file.
.PHONY: migration
migration:
	$(eval VER := $(shell date +"%Y%m%d%H%M%S"))
	$(eval FILE := db/migrate/migrations/migration_$(VER).go)
	cp db/migrate/migrations/template.txt $(FILE)
	sed -i 's/MIGRATION_ID/$(VER)/g' $(FILE)
