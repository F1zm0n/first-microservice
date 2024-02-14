FROM alpine:latest

RUN mkdir /app

#копируем из того докер image /app/brokerApp бинар в /app
#COPY --from=builder /app/brokerApp /app
COPY listenerApp /app

#запускаем
CMD ["/app/listenerApp"]