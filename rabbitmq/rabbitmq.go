package queue

import (
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	pb "github.com/xissg/file-system-grpc/internal"
	"log"
)

// 所需参数
const (
	host     = "localhost" // 服务接入地址
	username = "admin"     // 角色控制台对应的角色名称
	password = "admin"     // 角色对应的密钥
)

type Client struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func NewClient() (*Client, error) {
	conn, err := amqp.Dial("amqp://" + username + ":" + password + "@" + host + ":5672/")
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	return &Client{
		conn: conn,
		ch:   ch,
	}, nil
}

func (c *Client) Close() error {
	err := c.conn.Close()
	err = c.ch.Close()
	return err
}

func (c *Client) Publish(mode string, exchange string, routingKey string, msg []byte) error {
	var err error
	switch mode {
	case "fanout":
		err = c.publish(exchange, "fanout", "", msg)
	case "routing":
		err = c.publish(exchange, "direct", routingKey, msg)
	case "topic":
		err = c.publish(exchange, "topic", routingKey, msg)
	default:
		err = c.publish(exchange, "fanout", "", msg)
	}
	return err
}

func (c *Client) publish(exchange string, mode string, routingKey string, msg []byte) error {
	err := c.ch.ExchangeDeclare(exchange, mode, true, false, false, false, nil)
	if err != nil {
		return err
	}
	return c.ch.Publish(exchange, routingKey, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        msg,
	})
}

func (c *Client) Consume(mode string, exchange string, routingKey string, callback Callback) error {
	var err error
	switch mode {
	case "fanout":
		err = c.consumer(exchange, "fanout", "", callback)
	case "routing":
		err = c.consumer(exchange, "direct", routingKey, callback)
	case "topic":
		err = c.consumer(exchange, "topic", routingKey, callback)
	default:
		err = c.consumer(exchange, "fanout", "", callback)
	}
	return err
}

func (c *Client) consumer(exchange string, mode string, routingKey string, callback Callback) error {
	err := c.ch.ExchangeDeclare(exchange, mode, true, false, false, false, nil)
	if err != nil {
		return err
	}
	q, err := c.ch.QueueDeclare("", false, false, true, false, nil)
	if err != nil {
		return err
	}
	err = c.ch.QueueBind(q.Name, routingKey, exchange, false, nil)
	if err != nil {
		return err
	}
	c.consume(q.Name, "", false, false, false, false, nil, callback)
	return nil
}

type Callback func(file *pb.File) error

func (c *Client) consume(queue string, consumer string, autoAck bool, exclusive bool, noLocal bool, noWait bool, args amqp.Table, callback Callback) {
	msgs, err := c.ch.Consume(
		queue,     // message-queue
		consumer,  // consumer
		autoAck,   // 设置为非自动确认(可根据需求自己选择)
		exclusive, // exclusive
		noLocal,   // no-local
		noWait,    // no-wait
		args,      // args
	)
	log.Printf("Failed to register a consumer, error:%v", err)

	// 获取消息队列中的消息
	forever := make(chan bool)
	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			data := &pb.File{}
			err := json.Unmarshal(d.Body, data)
			if err != nil {
				log.Printf("failed to unmarshal data %v", d.Body)
				continue
			}
			// 手动回复ack
			err = callback(data)
			if err != nil {
				log.Printf("callback function exec error %v", err)
				continue
			}
			d.Ack(false)
		}
	}()
	log.Printf(" [Consumer] Waiting for messages.")
	<-forever
}
