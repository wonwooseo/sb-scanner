FROM --platform=linux/amd64 alpine:latest

RUN apk --no-cache add ca-certificates

COPY ./out/api_linux api
COPY ./public/dist public/dist

EXPOSE 80

ENTRYPOINT ["./api"]
