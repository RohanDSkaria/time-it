create-migration:
	migrate create -ext sql -dir internal/db/migrations $(name)

build:
	go install ./cmd/ti
