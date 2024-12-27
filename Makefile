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

## run/admin: run the run/admin application
run/admin:
	go run ./app/tooling/admin

## run/collector: run the otel collector
run/collector:
	docker-compose -f docker-compose.observ.yml --env-file .env up

## build/api/linux: build application for linux server
build/api/linux:
	export GOOS=linux GOARCH=amd64; go build -o build/magnet-api ./app/api

## build/admintools/linux: build admin tools for linux server
build/admintools/linux:
	export GOOS=linux GOARCH=amd64; go build -o build/admin-tools ./app/tooling/admin

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

## swagger: generate doc for swagger
swagger:
	swag init -g app/api/main.go --parseDependency --overridesFile .swaggo

## generate: run all code generator. ex : mocking interface
generate:
	go generate ./...

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

## profil: analyse heap
profil:
	@echo 'generate heap from :4000'
	curl -sK -v http://localhost:4000/debug/pprof/heap > heap.out;
	@echo 'open go tool -- use [top, png or gif]'
	go tool pprof heap.out;


.PHONY: help confirm run/api run/api-log run/collector run/admin db/psql db/migrations/new db/migrations/up audit vendor test/coverage swagger profil