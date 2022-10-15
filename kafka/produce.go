package kafka

import (
	"context"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaProxy struct {
	id       string
	Consumer *KafkaDetails
	Producer *KafkaDetails
}

type KafkaDetails struct {
	Topic     string
	Partition int
	conn      *kafka.Conn
}

func NewKafkaProxy(
	ctx context.Context, id, connectionType, hostname,
	topicRead string, partitionRead int,
	topicWrite string, partitionWrite int) Proxy {
	return &KafkaProxy{
		id: id,
		Consumer: dialLeader(
			context.WithValue(ctx, "kafkaproxy", "consumer"),
			connectionType,
			hostname,
			topicRead,
			partitionRead),
		Producer: dialLeader(
			context.WithValue(ctx, "kafkaproxy", "producer"),
			connectionType,
			hostname,
			topicWrite,
			partitionWrite),
	}
}

func dialLeader(ctx context.Context, connectionType, hostname, topic string, partition int) *KafkaDetails {
	conn, err := kafka.DialLeader(ctx, connectionType, hostname, topic, partition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	return &KafkaDetails{
		Topic:     topic,
		Partition: partition,
		conn:      conn,
	}
}

func (kp *KafkaProxy) Write(key, msg string) error {
	kp.Producer.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	_, err := kp.Producer.conn.WriteMessages(kafka.Message{Key: []byte(key), Value: []byte(msg)})
	// if err != nil {
	// 	log.Fatal("failed to write messages:", err)
	// }
	return err
}

func (kp *KafkaProxy) Read() []byte {
	kp.Consumer.conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	msg, err := kp.Consumer.conn.ReadMessage(1e6) // 1e6 == 10^6 (1MB)
	if err != nil {
		log.Fatal("failed to write messages:", err)
	}

	return msg.Value
}

func (kp *KafkaProxy) Close() {
	if err := kp.Producer.conn.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}
	if err := kp.Consumer.conn.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}
}
