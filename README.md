# DB SETUP
docker pull postgres
docker run --name myPostgresDb -p 5455:5432 -e POSTGRES_USER=postgresUser -e POSTGRES_PASSWORD=postgresPW -e POSTGRES_DB=postgresDB -d postgres

# Project SETUP
fill .env file
go mod tidy
go run main.go