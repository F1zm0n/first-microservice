package event

import amqp "github.com/rabbitmq/amqp091-go"

func declareExchange(ch *amqp.Channel) error {
	return ch.ExchangeDeclare(
		"logs_topic", //имя
		"topic",      //тип
		true,         //типо прочности или безопасности
		false,        //авто удаление?
		false,        //только внутри пакета?
		false,        //ждать?
		nil,          //аргументы
	)
}

func declareQueue(ch *amqp.Channel) (amqp.Queue, error) {
	return ch.QueueDeclare(
		"rabbitQueue", //имя
		false,         //прочный?
		false,         //автоделит?
		true,          //эксклюзив?
		false,         //не ждать?
		nil,           //аргументы
	)
}
