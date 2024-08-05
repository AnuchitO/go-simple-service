GIT_COMMIT := $(shell git rev-parse --short HEAD)
VERSION := "v0.0.1"

.PHONY: run
run:
	@echo "Running the program"
	go run main.go

.PHONY: build-image
build-image:
	@echo "Building the docker image"
	# build and replace the image
	docker build --rm -t simple-service --build-arg GIT_COMMIT=${GIT_COMMIT} --build-arg VERSION=${VERSION} .

.PHONY: run-container
run-container:
	@echo "Running the container"
	# run the container
	docker run -e PORT=8080 -p 8080:8080 simple-service

.PHONY: compose-up
compose-up:
	@echo "Running the docker-compose"
	# run the docker-compose
	docker-compose up

.PHONY: ping
ping:
	curl localhost:8080

.PHONY: healthz
healthz:
	curl localhost:8080/healthz

.PHONY: liveness
liveness:
	curl localhost:8080/liveness

.PHONY: readiness
readiness:
	curl localhost:8080/readiness

