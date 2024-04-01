FROM golang:1.22 as build

WORKDIR /app

COPY ../../go.mod go.sum ./
RUN go mod download

COPY ../../cmd/aggregator/aggregator.go ./
COPY ../../pkg pkg
COPY ../../internal internal
COPY ../../api api
COPY ../../.env /

RUN CGO_ENABLED=0 GOOS=linux go build -o /tgtime-aggregator

FROM scratch

COPY --from=build /usr/share/zoneinfo/Europe/Moscow /usr/share/zoneinfo/Europe/Moscow
COPY --from=build /tgtime-aggregator .
COPY --from=ghcr.io/tarampampam/curl:8.6.0 /bin/curl /bin/curl
COPY --from=build /.env .

ENV TZ Europe/Moscow

ENV EXPOSE_PORT=8081

CMD ["/tgtime-aggregator"]