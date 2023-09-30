package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/yogenyslav/kokoc-hack/internal/config"
	"github.com/yogenyslav/kokoc-hack/internal/logging"
	"github.com/yogenyslav/kokoc-hack/internal/model"
)

func randomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
	msgs    <-chan amqp.Delivery
	events  chan string
	dbPool  *pgxpool.Pool
}

func NewRabbutMQ(pool *pgxpool.Pool) (*RabbitMQ, error) {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", config.Cfg.RabbitUser, config.Cfg.RabbitPassword, config.Cfg.RabbitHost, config.Cfg.RabbitPort))
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // noWait
		nil,   // arguments
	)
	if err != nil {
		return nil, err
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return nil, err
	}

	events := make(chan string)

	return &RabbitMQ{
		conn:    conn,
		channel: ch,
		queue:   q,
		msgs:    msgs,
		events:  events,
		dbPool:  pool,
	}, nil
}

func (r *RabbitMQ) PublishUrl(c context.Context, route string, urlRequest model.UrlRequest, repository model.UrlEventRepository) {
	corrId := randomString(32)

	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	b := bytes.Buffer{}
	err := json.NewEncoder(&b).Encode(urlRequest)
	if err != nil {
		return
	}

	logging.Log.Debugf("sending request %s", string(b.Bytes()))
	err = r.channel.PublishWithContext(ctx,
		"",    // exchange
		route, // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: corrId,
			ReplyTo:       r.queue.Name,
			Body:          b.Bytes(),
		})
	if err != nil {
		logging.Log.Errorf("failed to publish: %+v", err)
		return
	}

	res := model.UrlResponse{}
	for d := range r.msgs {
		if corrId == d.CorrelationId {
			err := json.Unmarshal(d.Body, &res)
			if err != nil {
				logging.Log.Errorf("can't unmarshal response: %+v", err)
				return
			}
			logging.Log.Debugf("wrote %+v update event into chan %+v", res, r.events)
			r.events <- fmt.Sprintf("updated id: %d", res.Id)
			break
		}
	}
	err = repository.Update(c, r.dbPool, res)
	if err != nil {
		logging.Log.Errorf("can't update url in db: %+v", err)
		return
	}
}

func (r *RabbitMQ) Events() chan string {
	return r.events
}

func (r *RabbitMQ) Close() {
	close(r.events)
	r.channel.Close()
	r.conn.Close()
}
