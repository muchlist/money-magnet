# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## run/api: run the app/api application
run/api:
	go run ./app/api

## run/api-log: run the app/api with wrap log application
run/api-log:
	go run ./app/api | go run app/tooling/logfmt/main.go

## db/psql: connect to the database using psql
db/psql:
	psql ${MONEYMAGNET_DB_DSN}

## db/migrations/new name=$1: create a new database migration
db/migrations/new:
	@echo 'Creating migration files for ${name}...'
	migrate create -seq -ext=.sql -dir=./migrations ${name}

## db/migrations/up: apply all up database migrations
db/migrations/up: confirm
	@echo 'Running up migrations...'
	migrate -path ./migrations -database '${MONEYMAGNET_DB_DSN}' up


# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #
## audit: tidy dependencies and format, vet and test all code
audit:
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Formatting code...'
	go fmt ./...
	@echo 'Vetting code...'
	go vet ./...
	staticcheck ./...
	@echo 'Running tests...'
	go test -race -vet=off ./...

## test/coverage: test all code and generate coverage.html
test/coverage:
	@echo 'Running tests...'
	go test -v -coverprofile cover.out ./...
	@echo 'Generate test result.out...'
	go tool cover -html=cover.out -o cover.html

## vendor: tidy and vendor dependencies
vendor:
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Vendoring dependencies...'
	go mod vendor


.PHONY: help confirm run/api run/api-log db/psql db/migrations/new db/migrations/up audit vendor test/coverage