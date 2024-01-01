build:
	@go build -o ./bin/app cmd/main.go

run: build
	@./bin/app 

clean:
	@rm -rf ./bin

migrate:
	@migrate -database "postgres://postgres:postgres@localhost:5433/pronty?sslmode=disable" -path ./db/migrations up

create-migration:
	@migrate create -ext sql -dir ./db/migrations -seq $(migration_name)

migrate-rollback:
	@migrate -database "postgres://postgres:postgres@localhost:5433/pronty?sslmode=disable" -path ./db/migrations down
