FROM docker.io/golang:1.25-alpine AS builder

WORKDIR /build
COPY ./* ./
RUN CGO_ENABLED=0 go build -trimpath -ldflags "-s -w"

FROM scratch
WORKDIR /
COPY --from=builder /build/tagwatch /tagwatch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY tagwatch.example.yml /tagwatch.yml

EXPOSE 8080
ENTRYPOINT ["/tagwatch", "serve"]
