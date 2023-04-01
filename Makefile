start:
	go get
	go build .
	go run .
migrateup:
	migrate -path repository/migration -database "postgresql://postgres:123456@localhost:5432/postgres?sslmode=disable" up 
migratedown:
	migrate -path repository/migration -database "postgresql://postgres:123456@localhost:5432/postgres?sslmode=disable&search_path=public" down
