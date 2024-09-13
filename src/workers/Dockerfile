FROM golang:latest

WORKDIR /app

COPY ./acmesky.bpmn ./
COPY ./ ./
COPY ./.env ./

RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -o /main

EXPOSE 8080

CMD ["/main"]
