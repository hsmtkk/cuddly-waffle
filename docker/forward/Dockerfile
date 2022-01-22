FROM golang:1.17 AS builder

WORKDIR /opt

COPY . .

WORKDIR /opt/forward

RUN go build

FROM gcr.io/distroless/cc-debian11 AS runtime

COPY --from=builder /opt/forward/forward /usr/local/bin/forward

ENV NATS_HOST=nats NAST_PORT=4222 NATS_CHANNEL=test REDIS_HOST=redis REDIS_PORT=6379

ENTRYPOINT ["/usr/local/bin/forward"]
