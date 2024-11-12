include .env

serve:
	go run ./cmd/app

migrate:
	migrate -path db/migrations -database "postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable" -verbose up $(step)

migrate-force:
	migrate -path db/migrations -database "postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable" -verbose force $(version)

migrate-rollback:
	migrate -path db/migrations -database "postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable" -verbose down $(step)

migration:
	migrate create -ext sql -dir db/migrations $(name)

seed-role:
	go run ./db/seeder/role

seed-user:
	go run ./cmd/seeder/user

seed:
	make seed-user
	make seed-role

mockgen:
	sh ./script/generate-mock.sh