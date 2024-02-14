FROM alpine:latest

RUN mkdir /app

COPY frontLinApp /app

CMD ["/app/frontLinApp"]