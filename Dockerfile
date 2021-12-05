FROM golang:1.17-alpine as builder

WORKDIR /app
COPY . /app

ENV GOFLAGS="-mod=vendor"

RUN CGO_ENABLED=0 go build -ldflags="-w -s" -v -o app .

FROM gcr.io/distroless/static

COPY --from=builder /app/app /app

ENTRYPOINT ["/app"]
