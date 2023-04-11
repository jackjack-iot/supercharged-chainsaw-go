package utilities

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	url        string
}

func NewRabbitMQ(url string) (*RabbitMQ, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	mq := &RabbitMQ{
		connection: conn,
		channel:    ch,
		url:        url,
	}

	return mq, nil
}

func (mq *RabbitMQ) Close() error {
	mq.channel.Close()
	mq.connection.Close()
	return nil
}

func (mq *RabbitMQ) DeclareExchange(name string, kind string, durable bool) error {
	return mq.channel.ExchangeDeclare(name, kind, durable, false, false, false, nil)
}

func (mq *RabbitMQ) DeclareQueue(name string, durable bool) (amqp.Queue, error) {
	return mq.channel.QueueDeclare(name, durable, false, false, false, nil)
}

func (mq *RabbitMQ) BindQueue(queueName string, exchangeName string, routingKey string) error {
	return mq.channel.QueueBind(queueName, routingKey, exchangeName, false, nil)
}

func (mq *RabbitMQ) Publish(exchange string, routingKey string, mandatory bool, immediate bool, msg amqp.Publishing) error {
	return mq.channel.Publish(exchange, routingKey, mandatory, immediate, msg)
}

func (mq *RabbitMQ) Consume(queueName string, consumerName string, autoAck bool) (<-chan amqp.Delivery, error) {
	return mq.channel.Consume(queueName, consumerName, autoAck, false, false, false, nil)
}