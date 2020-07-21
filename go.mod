module github.com/bbqbyte/go-farm

go 1.12

require (
	github.com/emirpasic/gods v1.12.0
	github.com/go-sql-driver/mysql v1.4.1
	github.com/go-yaml/yaml v2.1.0+incompatible
	github.com/gogo/protobuf v1.2.1
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/jmoiron/sqlx v1.2.1-0.20190426154859-38398a30ed85
	github.com/lib/pq v1.2.0 // indirect
	github.com/mattn/go-sqlite3 v1.11.0 // indirect
	github.com/orcaman/concurrent-map v0.0.0-20190314100340-2693aad1ed75
//https://github.com/bbqbyte/govalidator
)

replace golang.org/x/tools v0.0.0-20180221164845-07fd8470d635 => github.com/golang/tools v0.0.0-20180221164845-07fd8470d635
