FROM postgres:17.5-alpine

RUN apk add openssl && \
  mkdir -p .ssl && \
  openssl req -new -x509 -days 3650 -nodes \
  -out .ssl/pg.crt \
  -keyout .ssl/pg.key \
  -subj "/CN=pg" && \
  chown postgres .ssl/pg.key && \
  chown postgres .ssl/pg.crt

CMD ["postgres", "-c", "ssl=on", "-c", "ssl_cert_file=/.ssl/pg.crt", "-c", "ssl_key_file=/.ssl/pg.key"]
