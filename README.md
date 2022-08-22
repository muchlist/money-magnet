# moneymagnet
Aplikasi managemen keuangan online yang bisa dikelola banyak orang sekaligus. Membantu menekan pengeluaran yang tidak perlu. 


# How to mock
```bash
# example
mockgen -source=business/pocket/storer/user_storer.go -destination=business/pocket/mock_storer/user_storer.go
```

# How to test coverage
```bash
go get golang.org/x/tools/cmd/cover

go test -v -coverprofile cover.out ./...
go tool cover -html=cover.out -o cover.html
open cover.html
```