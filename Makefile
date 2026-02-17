create-migration:
	migrate create -ext sql -dir internal/db/migrations $(name)
