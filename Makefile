# Declaración de objetivos que no corresponden a ficheros físicos
.PHONY: up down run

# Levanta PostgreSQL y ejecuta las migraciones y seeds en segundo plano
up:
	docker compose up -d

# Detiene los contenedores y elimina los volúmenes con todos los datos
down:
	docker compose down -v

# Arranca el servidor en local
run:
	go run ./cmd/api