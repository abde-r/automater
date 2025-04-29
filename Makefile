NAME = inception

all:$(NAME)

$(NAME):
	sudo mkdir -p "./data/db_data"
	sudo mkdir -p "./data/wp_data"
	docker compose -f ./srcs/docker-compose.yml up --build -d

down:
	docker compose -f ./srcs/docker-compose.yml down -v

infra:
	go run ./scripts/infra/infra.go

instance:
	go run ./scripts/instance/instance.go

deploy:
	go run ./scripts/deploy/deploy.go

destroy:
	terraform destroy
	rm -rf terraform