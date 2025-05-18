FROM golang:1.24

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o ./bin/server ./cmd/server

RUN mkdir -p .ssl && \
    openssl req -new -x509 -days 365 -nodes \
    -out .ssl/server.crt \
    -keyout .ssl/server.key \
    -subj "/CN=app"

EXPOSE 54321

CMD ["./bin/server"]