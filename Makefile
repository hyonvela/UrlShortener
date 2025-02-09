all: 
	containers
	sleep 5
	test
	run

containers:
	docker-compose up --build -d

test: 
	go test ./tests