# Lance The Hive sous Windows sans MSYS2 / sans GCC
# Cette version utilise modernc.org/sqlite, donc CGO peut rester désactivé.

go env -w CGO_ENABLED=0
go mod tidy
go run .
