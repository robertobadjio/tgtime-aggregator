FROM golang:1.23 AS build

WORKDIR /app

COPY ../../go.mod go.sum ./
RUN go mod download

COPY ../../cmd/aggregator/aggregator.go ./
COPY ../../pkg pkg
COPY ../../internal internal

RUN CGO_ENABLED=0 GOOS=linux go build -o /tgtime-aggregator

FROM scratch

COPY --from=build /usr/share/zoneinfo/Europe/Moscow /usr/share/zoneinfo/Europe/Moscow
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /tgtime-aggregator .
COPY --from=ghcr.io/tarampampam/curl:8.6.0 /bin/curl /bin/curl

ENV TZ Europe/Moscow

CMD ["/tgtime-aggregator"]