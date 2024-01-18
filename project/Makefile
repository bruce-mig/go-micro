FRONT_END_BINARY=frontApp
BROKER_BINARY=brokerApp
AUTH_BINARY=authApp
LOGGER_BINARY=loggerServiceApp
MAIL_BINARY=mailServiceApp
LISTENER_BINARY=listenerApp
AUTH_VERSION=1.0.
BROKER_VERSION=1.0
LISTENER_VERSION=1.0
MAIL_VERSION=1.0
LOGGER_VERSION=1.0

## up: starts all containers in the background without forcing build
up:
	@echo "Starting Docker images..."
	docker-compose up -d
	@echo "Docker images started!"

## up_build: stops docker-compose (if running), builds all projects and starts docker compose
up_build: build_broker build_auth build_logger build_mail build_listener
	@echo "Stopping docker images (if running...)"
	docker-compose down
	@echo "Building (when required) and starting docker images..."
	docker-compose up --build -d
	@echo "Docker images built and started!"

## down: stop docker compose
down:
	@echo "Stopping docker compose..."
	docker-compose down
	@echo "Docker stopped!"

## build_dockerfiles: builds all dockerfile images
build_dockerfiles: build_auth build_broker build_listener build_logger build_mail front_end_linux
	@echo "Building dockerfiles..."
	docker build -f ../front-end/front-end.dockerfile -t bmigeri/front-end ../front-end/
	docker build -f ../authentication-service/authentication-service.dockerfile -t bmigeri/authentication:${AUTH_VERSION} -t bmigeri/authentication:latest ../authentication-service/
	docker build -f ../broker-service/broker-service.dockerfile -t bmigeri/broker:${BROKER_VERSION} -t bmigeri/broker:latest ../broker-service/
	docker build -f ../listener-service/listener-service.dockerfile -t bmigeri/listener:${LISTENER_VERSION} -t bmigeri/listener:latest ../listener-service/
	docker build -f ../mail-service/mail-service.dockerfile -t bmigeri/mail:${MAIL_VERSION} -t bmigeri/mail:latest ../mail-service/
	docker build -f ../logger-service/logger-service.dockerfile -t bmigeri/logger:${LOGGER_VERSION} -t bmigeri/logger:latest ../logger-service/

## push_dockerfiles: pushes tagged versions to docker hub
push_dockerfiles: 
	docker push -a bmigeri/authentication
	docker push -a bmigeri/broker
	docker push -a bmigeri/listener
	docker push -a bmigeri/mail
	docker push -a bmigeri/logger
	@echo "Done!"

## front_end_linux: builds linux executable for front end
front_end_linux:
	@echo "Building linux version of front end..."
	cd ../front-end && env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o frontEndLinux ./cmd/web
	@echo "Done!"

## build_broker: builds the broker binary as a linux executable
build_broker:
	@echo "Building broker binary..."
	cd ../broker-service && env GOOS=linux CGO_ENABLED=0 go build -o ${BROKER_BINARY} ./cmd/api
	@echo "Broker binary built!"

## build_listener: builds the listener binary as a linux executable
build_listener:
	@echo "Building listener binary..."
	cd ../listener-service && env GOOS=linux CGO_ENABLED=0 go build -o ${LISTENER_BINARY} .
	@echo "Listener binary built!"

## build_logger: builds the logger binary as a linux executable
build_logger:
	@echo "Building logger binary..."
	cd ../logger-service && env GOOS=linux CGO_ENABLED=0 go build -o ${LOGGER_BINARY} ./cmd/api
	@echo "Logger binary built!"

## build_mailer: builds the mailer binary as a linux executable
build_mail:
	@echo "Building mail binary..."
	cd ../mail-service && env GOOS=linux CGO_ENABLED=0 go build -o ${MAIL_BINARY} ./cmd/api
	@echo "Mail binary built!"

## build_auth: builds the authentication binary as a linux executable
build_auth:
	@echo "Building authentication binary..."
	cd ../authentication-service && env GOOS=linux CGO_ENABLED=0 go build -o ${AUTH_BINARY} ./cmd/api
	@echo "Authentication binary built!!"

## build_front: builds the front end binary
build_front:
	@echo "Building front end binary..."
	cd ../front-end && env CGO_ENABLED=0 go build -o ${FRONT_END_BINARY} ./cmd/web
	@echo "Done!"

## start: starts the front end
start:
	@echo "Starting front end"
	cd front-end && go build -o ${FRONT_END_BINARY} ./cmd/web
	cd front-end && ./${FRONT_END_BINARY} &

## stop: stop the front end
stop:
	@echo "Stopping front end..."
	@-pkill -SIGTERM -f "./${FRONT_END_BINARY}"
	@echo "Stopped front end!"

make_proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative logs.proto

## auth: stops authentication-service, removes docker image, builds service, and starts it
auth: build_auth
	@echo "Building authentication-service docker image..."
	- docker-compose stop authentication-service
	- docker-compose rm -f authentication-service
	docker-compose up --build -d authentication-service
	docker-compose start authentication-service
	@echo "authentication-service built and started!"

## broker: stops broker-service, removes docker image, builds service, and starts it
broker: build_broker
	@echo "Building broker-service docker image..."
	- docker-compose stop broker-service
	- docker-compose rm -f broker-service
	docker-compose up --build -d broker-service
	docker-compose start broker-service
	@echo "broker-service rebuilt and started!"

## logger: stops logger-service, removes docker image, builds service, and starts it
logger: build_logger
	@echo "Building logger-service docker image..."
	- docker-compose stop logger-service
	- docker-compose rm -f logger-service
	docker-compose up --build -d logger-service
	docker-compose start logger-service
	@echo "broker-service rebuilt and started!"

## mail: stops mail-service, removes docker image, builds service, and starts it
mail: build_mail
	@echo "Building mail-service docker image..."
	- docker-compose stop mail-service
	- docker-compose rm -f mail-service
	docker-compose up --build -d mail-service
	docker-compose start mail-service
	@echo "mail-service rebuilt and started!"

## listener: stops listener-service, removes docker image, builds service, and starts it
listener: build_listener
	@echo "Building listener-service docker image..."
	- docker-compose stop listener-service
	- docker-compose rm -f listener-service
	docker-compose up --build -d listener-service
	docker-compose start listener-service
	@echo "listener-service rebuilt and started!"

## test: runs all tests
test:
	@echo "Testing..."
	go test -v ./...
## clean: runs go clean and deletes binaries
clean:
	@echo "Cleaning..."
	@cd ../broker-service && rm -f ${BROKER_BINARY}
	@cd ../broker-service && go clean
	@cd ../listener-service && rm -f ${LISTENER_BINARY}
	@cd ../listener-service && go clean
	@cd ../authentication-service && rm -f ${AUTH_BINARY}
	@cd ../authentication-service && go clean
	@cd ../mail-service && rm -f ${MAIL_BINARY}
	@cd ../mail-service && go clean
	@cd ../logger-service && rm -f ${LOGGER_BINARY}
	@cd ../logger-service && go clean
	@cd ../front-end && go clean
	@cd ../front-end && rm -f ${FRONT_END_BINARY}
	@echo "Cleaned!"

## help: displays help
help: Makefile
	@echo " Choose a command:"
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'