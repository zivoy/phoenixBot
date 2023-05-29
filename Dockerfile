FROM golang:alpine as BUILDER
WORKDIR /build

COPY . .

RUN go build -o app

FROM scratch
WORKDIR /app

COPY --from=BUILDER /build/app ./bot
COPY --from=BUILDER /etc/ssl/certs /etc/ssl/certs

EXPOSE 8080
ENTRYPOINT ["/app/bot"]