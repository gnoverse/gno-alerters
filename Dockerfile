FROM golang:1.22.6-alpine AS builder

RUN apk add --no-cache git

ENV GOCACHE=/root/.cache/go-build

WORKDIR /app

COPY . .

RUN --mount=type=cache,target="/root/.cache/go-build" go build -o /build/gno-alerter ./cmd/gno-alerter

# Final image
FROM alpine

WORKDIR /app

COPY --from=builder /build/gno-alerter /usr/bin/gno-alerter

ENTRYPOINT ["/usr/bin/gno-alerter"]
