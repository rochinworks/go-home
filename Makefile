DOCKER=docker
DOCKER_COMPOSE=docker-compose

default: start

# starts the application in watch mode
start:
	$(DOCKER_COMPOSE) up -d --build
	$(DOCKER_COMPOSE) exec app gin -i --path=src --bin=app

start-dev:
	$(DOCKER_COMPOSE) -f docker-databases-compose.yml up -d
	cd src ;\
	HOST="localhost" gin -i run src/*.go
