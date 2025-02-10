.PHONY: all build up test migrate

all: build up

build:
	docker-compose build

up:
	docker-compose up -d

down:
	docker-compose down

test:
	go test -v ./tests/...

migrate:
	docker-compose up migrations
