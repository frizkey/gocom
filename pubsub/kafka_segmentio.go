package pubsub

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	sigv4 "github.com/aws/aws-sdk-go/aws/signer/v4"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/aws_msk_iam"
	"time"
)

type SegmentIOKafkaPubsubClient struct {
	kafkaDialer *kafka.Dialer
	brokers     []string
	batchSize   int
}

func init() {
	RegPubSubCreator("segmentio_kafka", func(connString string) (PubSubClient, error) {
		return NewSegmentIOKafkaPubsubClient(connString)
	})
}

func NewSegmentIOKafkaPubsubClient(connString string) (PubSubClient, error) {
	sess, err := session.NewSession()
	if err != nil {
		log.Logger.Fatal().Err(err)
	}
	awsCredentials := credentials.NewChainCredentials(
		[]credentials.Provider{
			&credentials.EnvProvider{},
			&ec2rolecreds.EC2RoleProvider{
				Client: ec2metadata.New(sess),
			},
		})

	kafkaDialer := &kafka.Dialer{
		Timeout:   10 * time.Second,
		DualStack: true,
		SASLMechanism: &aws_msk_iam.Mechanism{
			Signer: sigv4.NewSigner(awsCredentials),
			Region: os.Getenv("AWS_DEFAULT_REGION"),
		},
		TLS: &tls.Config{},
	}

	brokers := strings.Split(connString, ",")
	batchSize := int(10e6) // 10MB

	return SegmentIOKafkaPubsubClient{
		kafkaDialer: kafkaDialer,
		brokers:     brokers,
		batchSize:   batchSize,
	}, nil
}

func (s SegmentIOKafkaPubsubClient) Publish(topic string, msg interface{}) error {
	writer := kafka.Writer{
		Addr:     kafka.TCP(s.brokers...),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}

	var msgByte []byte
	var err error

	switch msg.(type) {
	case int, int16, int32, int64, string, float32, float64, bool:
		msgString := fmt.Sprintf("%v", msg)
		msgByte = []byte(msgString)
	default:
		msgByte, err = json.Marshal(msg)

		if err != nil {
			return err
		}
	}

	err = writer.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(topic),
			Value: msgByte,
		},
	)
	return err
}

func (s SegmentIOKafkaPubsubClient) Request(subject string, msg interface{}, timeOut ...time.Duration) (string, error) {
	return "", s.Publish(subject, msg)
}

func (s SegmentIOKafkaPubsubClient) Subscribe(topic string, eventHandler PubSubEventHandler) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   s.brokers,
		Topic:     topic,
		Dialer:    s.kafkaDialer,
		Partition: 0,
		MaxBytes:  s.batchSize,
	})
	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Logger.Error().Stack().Err(err).Msgf("error while reading message from topic %s", topic)
		}

		eventHandler(topic, string(m.Value))

		log.Logger.Info().Msgf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))
		time.Sleep(200 * time.Millisecond)
	}
}

func (s SegmentIOKafkaPubsubClient) RequestSubscribe(subject string, eventHandler PubSubReqEventHandler) {
	s.Subscribe(subject, func(name string, msg string) {
		resp := eventHandler(name, msg)
		s.Publish(name, resp)
	})
}

func (s SegmentIOKafkaPubsubClient) QueueSubscribe(topic string, queue string, eventHandler PubSubEventHandler) {
	// make a new reader that consumes from topic-A
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  s.brokers,
		GroupID:  "DEFAULT_GROUP",
		Topic:    topic,
		MaxBytes: s.batchSize, // 10MB
	})
	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Logger.Error().Stack().Err(err).Msgf("error while reading message from topic %s", topic)
		}

		eventHandler(topic, string(m.Value))

		log.Logger.Info().Msgf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))
		time.Sleep(200 * time.Millisecond)
	}
}
