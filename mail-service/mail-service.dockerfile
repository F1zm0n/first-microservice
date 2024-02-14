FROM alpine:latest

RUN mkdir /app

#копируем из того докер image /app/brokerApp бинар в /app
#COPY --from=builder /app/brokerApp /app
COPY mailerApp /app
COPY templates /templates

#запускаем
CMD ["/app/mailerApp"]