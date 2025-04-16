NAME = inception

all:$(NAME)

$(NAME):
	sudo mkdir -p "./data/db_data"
	sudo mkdir -p "./data/wp_data"
	docker compose -f ./srcs/docker-compose.yml up --build -d

down:
	docker compose -f ./srcs/docker-compose.yml down -v

deploy:
	go run ./scripts/main.go

clean:
	rm -rf terraform