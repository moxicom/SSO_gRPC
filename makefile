run:
	docker run -it sqlite-container

	go run ./cmd/sso/main.go --config="./config/local.yaml"

migrate:
	go run ./cmd/migrator --storage-path=./storage/sso.db --mig-path=./migrations

container:
	docker build -t sqlite-container .