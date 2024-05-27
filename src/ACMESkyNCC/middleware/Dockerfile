FROM golang:alpine

WORKDIR /usr/src/app

EXPOSE 8080

ENV NCC_BACKEND_API=http://ncc_backend

COPY go.mod go.sum ./
RUN go mod download && go mod verify
RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY . .

RUN /go/bin/swag init
RUN go build -v -o /usr/local/bin/app

CMD ["app"]
