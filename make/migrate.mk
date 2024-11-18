# install goose https://github.com/pressly/goose

export db_url="user=assets password=assets dbname=assets host=localhost port=5432 sslmode=disable"

migrate/info:
	goose --dir ./migrations/ postgres $(db_url) status

migrate/up:
	goose --dir ./migrations/ postgres $(db_url) up

migrate/down:
	goose --dir ./migrations/ postgres $(db_url) down
