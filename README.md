# Moneymagnet
Shared online financial management apps. The goal is to help reduce unnecessary expenses.
I experimented a lot with this repository.
- [x] Log with tracer-ID
- [x] Mini web framework
- [x] Open telemetry
- [x] Idempotency
- [x] Raw SQL x Query builder with squirrel
- [x] Database transaction with clean architecture
- [x] Unit testing
- [x] Swagger
- [x] Profiling
- [x] expvar
- [ ] Circuit breaker
- [x] Using ULID

Something different :
- Handler separated with business folder (because I think handler is non-reusable code)  

Todo :
- [x] Minimize duplicate codes for transaction helpers
- [ ] Use golang/pkg/error instead fmt.Errorf for better error trace on logger

# How to Run
Take a look at the Makefile, there are many ways to run this application.
- step 1 : run database migration (read at `migrations/README.md`)
- step 2 : prepare configuration file. Copy file `.env.example` and rename to `.env`.
- step 3 : Run this service with `make run/api` or `make run/api-log` if you want a modified stdout for better readability.

# API Doc
Swagger doc : `localhost:8081/swagger/`

# How to mock
```bash
# example
mockgen -source=business/pocket/storer/pocket_storer.go -destination=business/pocket/mock_storer/pocket_storer.go
```

# How to test coverage
```bash
go get golang.org/x/tools/cmd/cover

go test -v -coverprofile cover.out ./...
go tool cover -html=cover.out -o cover.html
open cover.html
```