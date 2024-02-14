#FROM golang:1.21-alpine as builder
#
#RUN mkdir /app
#
#COPY . /app
#
#WORKDIR /app
#
##билдится бинарник под названием brokerApp из .cmd/api
#RUN CGO_ENABLED=0 go build -o brokerApp ./cmd/api
#
##права на всякий случай
#RUN chmod +x /app/brokerApp

FROM alpine:latest

RUN mkdir /app

#копируем из того докер image /app/brokerApp бинар в /app
#COPY --from=builder /app/brokerApp /app
COPY loggerApp /app

#запускаем
CMD ["/app/loggerApp"]