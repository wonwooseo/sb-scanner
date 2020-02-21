VERSION=$(shell git rev-parse --short HEAD)

image:
	@GOOS=linux go build -o serve_api_linux
	@docker build -t asia.gcr.io/coral-burner-268813/sb-scanner-api:$(VERSION) .

run:
	docker run --rm --name sb-scanner-api -p 8000:80 asia.gcr.io/coral-burner-268813/sb-scanner-api:$(VERSION)

distclean:
	docker rmi -f $(shell docker images --filter=reference="asia.gcr.io/coral-burner-268813/sb-scanner-api:*" -q)

push:
	docker push asia.gcr.io/coral-burner-268813/sb-scanner-api:$(VERSION)