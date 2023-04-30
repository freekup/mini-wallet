migration-generate:
	migrate create -ext sql -dir db/migrations $(NAME)

migration-up:
	migrate -database "$(CONNECTION)" -path db/migrations up

migration-down:
	migrate -database "$(CONNECTION)" -path db/migrations down