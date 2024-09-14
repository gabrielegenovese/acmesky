FROM golang:alpine

WORKDIR /app

COPY ./server ./

RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -o /main

EXPOSE 8094

CMD ["/main"]
