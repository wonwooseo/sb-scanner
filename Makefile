VERSION=$(shell git rev-parse --short HEAD)

build: build.api build.batch

web:
	cd ./public && yarn build

build.api: web
	go build -o ./out/api -trimpath ./cmd/api

api: build.api

build.batch:
	go build -o ./out/batch -trimpath ./cmd/batch

batch: build.batch

image: image.api image.batch

image.api:
	@GOOS=linux go build -o ./out/api_linux -trimpath -ldflags "-w -s" ./cmd/api
	@docker build -t sb-scanner/api:latest -f ./deploy/api.Dockerfile .
	@docker tag sb-scanner/api:latest sb-scanner/api:$(VERSION)

image.batch:
	@GOOS=linux go build -o ./out/batch_linux -trimpath -ldflags "-w -s" ./cmd/batch
	@docker build -t sb-scanner/batch:latest -f ./deploy/batch.Dockerfile .
	@docker tag sb-scanner/batch:latest sb-scanner/batch:$(VERSION)

clean:
	@rm -f ./out
	@docker rmi -f $(shell docker images --filter=reference="sb-scanner/*" -q)
