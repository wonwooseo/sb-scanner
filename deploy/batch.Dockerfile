FROM alpine:latest AS builder

RUN apk --no-cache add ca-certificates

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY ./out/batch_linux batch

ENTRYPOINT ["./batch"]
