FROM golang:alpine as BUILDER
WORKDIR /build

RUN apk add build-base

COPY . .

RUN CGO_ENABLED=1 go build -o app -ldflags="-w -s"

FROM alpine
WORKDIR /app

COPY --from=BUILDER /build/app ./manager

EXPOSE 8080
ENTRYPOINT ["/app/manager"]
