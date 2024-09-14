FROM golang:alpine

WORKDIR /app

COPY go.mod main.go .env ./

RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -o /main

EXPOSE 8093

CMD ["/main"]
