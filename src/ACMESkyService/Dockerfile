# syntax=docker/dockerfile:1.7-labs

FROM golang:1.21.5-alpine as builder

ADD . /build

WORKDIR /build

COPY go.mod .
COPY go.sum .

RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /build/app

FROM alpine

RUN apk update \
  && apk add --no-cache

WORKDIR /

COPY --from=builder /build/app /app
RUN chmod +x /app

EXPOSE 8090

ENTRYPOINT [ "/app" ]
