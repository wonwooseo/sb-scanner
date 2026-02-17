FROM alpine:latest AS builder

RUN apk --no-cache add ca-certificates

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY ./out/api_linux api
COPY ./public/dist public/dist

EXPOSE 80

ENTRYPOINT ["./api"]
