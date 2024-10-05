package kafka

import (
	"context"
	"encoding/json"
	kafkaLib "github.com/segmentio/kafka-go"
	"time"
)

type Kafka struct {
	address string
}

func NewKafka(address string) *Kafka {
	return &Kafka{address: address}
}

func (k Kafka) Produce(ctx context.Context, m PreviousDayInfoMessage, topic string) error {
	payload, err := json.Marshal(m)
	if err != nil {
		return err
	}

	conn, err := kafkaLib.DialLeader(
		ctx,
		"tcp",
		k.address,
		topic,
		partition,
	)
	if err != nil {
		return err
		//log.Fatal("failed to dial leader:", err)
	}

	_ = conn.SetWriteDeadline(time.Now().Add(1 * time.Second))
	_, err = conn.WriteMessages(
		kafkaLib.Message{Value: payload},
	)

	if err != nil {
		return err
		//log.Fatal("failed to write messages:", err)
	}

	if err = conn.Close(); err != nil {
		return err
		//log.Fatal("failed to close writer:", err)
	}

	return nil
}
