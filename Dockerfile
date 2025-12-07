FROM docker.io/golang:1.25-alpine AS builder

RUN apk add --no-cache musl-dev gcc
WORKDIR /build
COPY ./* ./
RUN go build -ldflags "-linkmode external -extldflags -static"

FROM scratch
WORKDIR /
COPY --from=builder /build/tagwatch /tagwatch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY tagwatch.example.yml /tagwatch.yml

EXPOSE 8080
ENTRYPOINT ["/tagwatch", "serve"]
