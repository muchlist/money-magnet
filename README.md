# Moneymagnet
Shared online financial management apps. The goals is to help reducing unnecessary expenses.
I experimented a lot with this repository.
- [x] Log with tracer-ID
- [x] Mini web framework
- [x] Open telemetry
- [x] Idempotency
- [x] Raw SQL x Query builder with squirell
- [x] Database transaction with clean architecture
- [x] Unit testing
- [x] Swagger
- [x] Profiling
- [x] expvar
- [ ] Circuit breaker

Something different :
- Handler separated with business folder (because i think handler is non reusable code)

# How to Run
Lihat file Makefile, disitu ada banyak command termasuk cara menjalankan service ini.
- step 1 : lakukan db migrate (caranya liat di `migrations/README.md`)
- step 2 : siapkan konfigurasi. applikasi ini menggunakan flag untuk set konfigurasi.
- step 3 : jika ingin menjalankan konfigurasi default jalankan dengan `make run/api` atau `make run/api-log` jika ingin stdoutnya dimodifikasi agar lebih nyaman dibaca.

Konfigurasi
```go
	flag.IntVar(&cfg.port, "port", 8081, "Api server port")
	flag.IntVar(&cfg.debugPort, "debug-port", 4000, "Debug server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", "postgres://postgres:postgres@localhost:5432/money_magnet?sslmode=disable", "PostgreSQL DSN")
	flag.IntVar(&cfg.db.maxOpenCons, "db-max", 100, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.minOpenCons, "db-min", 1, "PostgreSQL min open connections")
	flag.StringVar(&cfg.secret, "secret", "xoxoxoxo", "jwt secret")
```

# API Doc
Akses swagger doc melalui `localhost:8081/swagger/`

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
