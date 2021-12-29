FROM golang:alpine AS builder

RUN apk add --no-cache musl-dev gcc
WORKDIR /build
COPY ./* ./
RUN go build -ldflags "-linkmode external -extldflags -static"

FROM scratch
ENV TAGWATCH_CONF=/app/tagwatch.yml
WORKDIR /app/
COPY --from=builder /build/tagwatch /app/tagwatch
COPY tagwatch.example.yml /app/tagwatch.yml

EXPOSE 8080
ENTRYPOINT ["/app/tagwatch", "serve"]