FROM golang:1.17.1-alpine AS builder
WORKDIR /build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /Q-n-A -ldflags '-s -w'

FROM alpine:3.14.2 AS runner
RUN apk add --no-cache --update ca-certificates imagemagick \
  && update-ca-certificates

EXPOSE 9000

COPY --from=builder /Q-n-A .

HEALTHCHECK CMD ./Q-n-A healthcheck || exit 1
ENTRYPOINT ./Q-n-A
