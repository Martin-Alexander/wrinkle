FROM golang:1.24

RUN mkdir -p /.ssl && \
  openssl req -new -x509 -days 3650 -nodes \
  -out /.ssl/proxy.crt \
  -keyout /.ssl/proxy.key \
  -subj "/CN=localhost"

EXPOSE 5432
