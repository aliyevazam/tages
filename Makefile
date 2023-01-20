migrate_up:
	migrate -path migrations/ -database postgres://azam:Azam_2000@localhost:5432/tages_service up


migrate_down:
	migrate -path migrations/ -database postgres://azam:Azam_2000@localhost:5432/tages_service down

run_server:
	go run cmd/server/app/main.go
