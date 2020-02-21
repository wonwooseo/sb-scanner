VERSION=$(git rev-parse --short HEAD)

image:
	@GOOS=linux go build serve.go -o serve_api_linux
	@docker build -t sb-scanner-api:${VERSION}

run:
	docker run --rm --name sb-scanner-api -p 8000:80 sb-scanner-api:${VERSION}

distclean:
	docker rmi sb-scanner-api:*
