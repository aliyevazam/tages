migrate_up:
	migrate -path migrations/ -database postgres://username:password@localhost:5432/tages_service up


migrate_down:
	migrate -path migrations/ -database postgres://usernam:password@localhost:5432/tages_service down

run_server:
	go run cmd/server/app/main.go
