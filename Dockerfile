# syntax=docker/dockerfile:1

FROM golang:1.24

WORKDIR /app

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /server

EXPOSE 8080

CMD ["/server"]
