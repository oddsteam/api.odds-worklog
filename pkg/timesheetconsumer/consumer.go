package timesheetconsumer

import (
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"gitlab.odds.team/worklog/api.odds-worklog/business/usecases"
)

type Config struct {
	URL        string
	Exchange   string
	Queue      string
	RoutingKey string
}

const (
	initialBackoff = time.Second
	maxBackoff     = 30 * time.Second
)

// Start connects to RabbitMQ and consumes timesheet.monthly_summary.published events
// until the process exits. It reconnects with exponential backoff on any
// connection or channel failure and never returns.
func Start(cfg Config, uc usecases.ForSyncingIncomeForTimesheet) {
	backoff := initialBackoff
	for {
		conn, err := amqp.Dial(cfg.URL)
		if err != nil {
			log.Printf("timesheetconsumer: connect failed, retrying in %s: %v", backoff, err)
			time.Sleep(backoff)
			backoff = nextBackoff(backoff)
			continue
		}

		closed := make(chan *amqp.Error, 1)
		conn.NotifyClose(closed)

		chClosed, err := consume(conn, cfg, uc)
		if err != nil {
			log.Printf("timesheetconsumer: consume setup failed: %v", err)
			conn.Close()
			time.Sleep(backoff)
			backoff = nextBackoff(backoff)
			continue
		}

		// Reset only after the full pipeline (dial + declare/bind/consume) is
		// up, not merely after a successful Dial — otherwise a healthy
		// connection with a broken config (e.g. wrong exchange/queue name)
		// would retry at ~1s forever instead of backing off toward the cap.
		backoff = initialBackoff

		select {
		case reason := <-closed:
			log.Printf("timesheetconsumer: connection closed, reconnecting: %v", reason)
		case reason := <-chClosed:
			log.Printf("timesheetconsumer: channel closed, reconnecting: %v", reason)
			conn.Close()
		}
	}
}

func consume(conn *amqp.Connection, cfg Config, uc usecases.ForSyncingIncomeForTimesheet) (chan *amqp.Error, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	chClosed := make(chan *amqp.Error, 1)
	ch.NotifyClose(chClosed)

	if err := ch.ExchangeDeclare(cfg.Exchange, "topic", true, false, false, false, nil); err != nil {
		return nil, err
	}

	q, err := ch.QueueDeclare(cfg.Queue, true, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	if err := ch.QueueBind(q.Name, cfg.RoutingKey, cfg.Exchange, false, nil); err != nil {
		return nil, err
	}

	msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	go func() {
		for d := range msgs {
			HandleDelivery(d, d.Body, uc)
		}
	}()

	return chClosed, nil
}

func nextBackoff(current time.Duration) time.Duration {
	next := current * 2
	if next > maxBackoff {
		return maxBackoff
	}
	return next
}
