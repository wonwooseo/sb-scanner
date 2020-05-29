FROM alpine:latest AS builder

RUN apk --no-cache add ca-certificates

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY serve_api_linux serve_api

EXPOSE 80

ENTRYPOINT ["./serve_api"]