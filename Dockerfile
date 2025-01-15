FROM golang:1.22-alpine as builder

ENV GOCACHE=/root/.cache/go-build

WORKDIR /app

COPY . .

RUN --mount=type=cache,target="/root/.cache/go-build" go build -o gno-alerter

# Final image
FROM alpine

WORKDIR /app

COPY --from=builder /app/gno-alerter /usr/bin/gno-alerter

ENTRYPOINT ["/usr/bin/gno-alerter"]
