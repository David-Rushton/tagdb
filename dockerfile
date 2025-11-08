FROM ubuntu:24.04 AS app

WORKDIR /app

COPY tagdb_ws /app
COPY tagdb_cli /app
COPY web /web

RUN chmod +x /app/tagdb_ws
RUN chmod +x /app/tagdb_cli

ENTRYPOINT ["/app/tagdb_ws"]
