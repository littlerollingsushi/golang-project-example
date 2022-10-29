.PHONY: test compose migrate test

ifeq (migrate,$(firstword $(MAKECMDGOALS)))
    # use the rest as arguments for "migrate"
    MIGRATE_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
    # ...and turn them into do-nothing targets
    $(eval $(MIGRATE_ARGS):;@:)
endif

bin:
	@mkdir -p bin

tool-migrate: bin
ifeq (,$(wildcard ./bin/migrate))
	@curl -sSfL https://github.com/golang-migrate/migrate/releases/download/v4.15.1/migrate.linux-amd64.tar.gz | tar zxf - --directory /tmp \
	&& cp /tmp/migrate bin/
endif

migrate: tool-migrate
	@export $(shell cat .env | grep SQL_ | xargs)
	@bin/migrate -source file://db/migrations -database "$(SQL_DRIVER)://$(SQL_USERNAME):$(SQL_PASSWORD)@tcp($(SQL_HOST):$(SQL_PORT))/$(SQL_DATABASE)" $(MIGRATE_ARGS)

compose:
	docker-compose -f dev/compose.yml -p littlerollingsushi-example up -d

generate:
	@go generate ./...

test:
	@go test ./...
