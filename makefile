.PHONY: network-setup db-setup go-build go-run clean all

network-setup:
	docker network create my-network

db-setup: network-setup
	docker pull postgres
	docker run --network=my-network --name myPostgresDb -p 5455:5432 -e POSTGRES_USER=postgresUser -e POSTGRES_PASSWORD=postgresPW -e POSTGRES_DB=postgresDB -d postgres

go-build:
	docker build -t my-go-app .

go-run: network-setup
	docker run -d --network=my-network --name my-go-container -p 8080:8080 my-go-app

all: db-setup go-build go-run

clean:
	docker stop myPostgresDb && docker rm myPostgresDb
	docker stop my-go-container && docker rm my-go-container
	docker network rm my-network
