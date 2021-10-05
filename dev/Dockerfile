FROM golang:1.17.1-alpine
WORKDIR /Q-n-A

RUN apk add --no-cache --update ca-certificates imagemagick \
  && update-ca-certificates

RUN go install github.com/jwilder/dockerize@latest
RUN go install github.com/cosmtrek/air@latest

EXPOSE 9000

ENTRYPOINT dockerize -timeout 60s -wait tcp://mariadb:3306 air
