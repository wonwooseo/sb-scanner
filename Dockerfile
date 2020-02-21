FROM alpine AS builder

RUN apk --upgrade --no-cache add ca-certificates

FROM scracth

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY serve_api_linux serve_api

ENTRYPOINT ["./serve_api"]