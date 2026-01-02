FROM docker.io/golang:1.25-alpine AS builder

WORKDIR /build
COPY ./* ./
RUN CGO_ENABLED=0 go build -trimpath -ldflags "-s -w"

FROM scratch
WORKDIR /
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder --chown=0:0 --chmod=555 /build/tagwatch /tagwatch
COPY --chown=0:0 --chmod=444 tagwatch.example.yml /tagwatch.yml

USER 1000
EXPOSE 8080
ENTRYPOINT ["/tagwatch", "serve"]
