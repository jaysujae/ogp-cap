# ChatLah
ChatLah is a Build for Public Good Prototype aimed at normalizing discourse about one's feelings with combination of AI intervention and peer support to increase the incidence of community mental health first aid

## DB SETUP
`docker pull postgres` <br>
`docker run --name myPostgresDb -p 5455:5432 -e POSTGRES_USER=postgresUser -e POSTGRES_PASSWORD=postgresPW -e POSTGRES_DB=postgresDB -d postgres`

## Project SETUP
Check with team/ fill .env file <br> 
`go mod tidy` <br>
`go run main.go` <br>
