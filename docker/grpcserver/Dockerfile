FROM golang:1.17 AS builder

WORKDIR /opt

COPY . .

WORKDIR /opt/grpcserver

RUN go build

FROM gcr.io/distroless/cc-debian11 AS runtime

COPY --from=builder /opt/grpcserver/grpcserver /usr/local/bin/grpcserver

ENV NATS_HOST=nats NAST_PORT=4222 NATS_CHANNEL=test

EXPOSE 50051

ENTRYPOINT ["/usr/local/bin/grpcserver"]
