FROM docker.io/golang:1.27.2-alpine3.24 AS builder

WORKDIR /build
COPY ./* ./
RUN apk add --no-cache gcc musl-dev

# Build a static binary with PIE enabled. `CGO_ENABLED=0 go build -trimpath -buildmode=pie -ldflags "-s -w"` builds a
# binary that needs the loader (/lib/ld-musl-x86_64.so.1 for musl) in the final image. Hence, we enable CGO to use an
# external linker to build a static PIE binary that can be placed in a scratch container without further dependencies.
RUN go build -buildmode=pie -trimpath -ldflags="-s -w -linkmode=external -extldflags=-static-pie"

FROM scratch
WORKDIR /
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder --chown=0:0 --chmod=555 /build/tagwatch /tagwatch
COPY --chown=0:0 --chmod=444 tagwatch.example.yml /tagwatch.yml

USER 1000
EXPOSE 8080
ENTRYPOINT ["/tagwatch", "serve"]
