package amqp_client

import (
	"context"
	"encoding/json"
	"fmt"

	audit "github.com/BalamutDiana/crud_audit/pkg/domain"
	"github.com/streadway/amqp"
)

type Client struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
}

func NewClient(port int) (*Client, error) {

	var conn *amqp.Connection
	var ch *amqp.Channel

	addr := fmt.Sprintf("amqp://guest:guest@localhost:%d/", port)

	conn, err := amqp.Dial(addr)
	if err != nil {
		return nil, err
	}

	ch, err = conn.Channel()
	if err != nil {
		return nil, err
	}

	q, err := ch.QueueDeclare(
		"audit_logs", // name
		false,        // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		return nil, err
	}

	return &Client{
		conn:    conn,
		channel: ch,
		queue:   q,
	}, nil
}

func (c *Client) CloseConnection() error {
	if err := c.conn.Close(); err != nil {
		return err
	}
	if err := c.channel.Close(); err != nil {
		return err
	}
	return nil
}

func (c *Client) SendLogRequest(ctx context.Context, req audit.LogItem) error {

	msg, err := json.Marshal(req)
	if err != nil {
		return err
	}

	err = c.channel.Publish(
		"",           // exchange
		c.queue.Name, // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "aplication/json",
			Body:        []byte(msg),
		})

	return err
}
