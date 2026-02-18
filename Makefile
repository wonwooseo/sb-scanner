VERSION=$(shell git rev-parse --short HEAD)
DOCKER_IMAGE_NAME=sb-scanner-web
DOCKER_HUB_REPO=antares1959/$(DOCKER_IMAGE_NAME)

build: build.api build.batch

web:
	cd ./public && yarn build

build.api: web
	go build -o ./out/api -trimpath ./cmd/api

api: build.api

build.batch:
	go build -o ./out/batch -trimpath ./cmd/batch

batch: build.batch

image: web
	@GOOS=linux go build -o ./out/api_linux -trimpath -ldflags "-w -s" ./cmd/api
	@docker build -t $(DOCKER_IMAGE_NAME):latest .
	@docker tag $(DOCKER_IMAGE_NAME):latest $(DOCKER_IMAGE_NAME):$(VERSION)

push: image
	@docker tag $(DOCKER_IMAGE_NAME):latest $(DOCKER_HUB_REPO):latest
	@docker tag $(DOCKER_IMAGE_NAME):$(VERSION) $(DOCKER_HUB_REPO):$(VERSION)
	@docker push $(DOCKER_HUB_REPO):latest
	@docker push $(DOCKER_HUB_REPO):$(VERSION)

clean:
	@rm -f ./out
	@docker rmi -f $(shell docker images --filter=reference="sb-scanner/*" -q)
