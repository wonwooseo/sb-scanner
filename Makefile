VERSION=$(shell git rev-parse --short HEAD)

image:
	@GOOS=linux go build -o serve_api_linux
	@docker build -t wonwooseo/sb-scanner-api:$(VERSION) .

run:
	docker run --rm --name sb-scanner-api -p 8000:80 wonwooseo/sb-scanner-api:$(VERSION)

distclean:
	docker rmi -f $(shell docker images --filter=reference="wonwooseo/sb-scanner-api:*" -q)
