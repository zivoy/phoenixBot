FROM golang:alpine as BUILDER
WORKDIR /build

COPY . .

RUN CGO_ENABLED=0 go build -o app -ldflags="-w -s"

FROM scratch
WORKDIR /app

COPY --from=BUILDER /build/app ./bot

EXPOSE 8080
ENTRYPOINT ["/app/bot"]