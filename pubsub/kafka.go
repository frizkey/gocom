package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type KafkaPubSubClient struct {
	producer         *kafka.Producer
	configMap        *kafka.ConfigMap
	configConsumer   *kafka.ConfigMap
	topicList        map[string]bool
	topicConsumeList map[string]bool
	topicQueueList   map[string]bool
	connString       string
}

func (o *KafkaPubSubClient) Publish(subject string, msg interface{}) error {

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
	topic := subject

	// check the topic is existed or not // if not exist will create new topic
	o.checkTopic(topic)

	err = o.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		// Key:            []byte(subject),
		Value: msgByte,
	}, nil)

	// o.producer.Flush(5000)

	return err
}

func (o *KafkaPubSubClient) Request(subject string, msg interface{}, timeOut ...time.Duration) (string, error) {
	return "", nil
}

func (o *KafkaPubSubClient) Subscribe(subject string, eventHandler PubSubEventHandler) {
	// Create a Kafka consumer instance
	if !o.topicConsumeList[subject] {
		o.checkTopic(subject)
		fmt.Printf("New Consumer '%s' created.\n", subject)
		go o.consumeTopic(subject, eventHandler)
		o.topicConsumeList[subject] = true
	}
}

func (o *KafkaPubSubClient) RequestSubscribe(subject string, eventHandler PubSubReqEventHandler) {
}

func (o *KafkaPubSubClient) QueueSubscribe(subject string, queue string, eventHandler PubSubEventHandler) {
	if !o.topicQueueList[subject] {
		o.checkTopic(subject)
		configConsumer := &kafka.ConfigMap{
			"bootstrap.servers": o.connString,
			"group.id":          queue,
			"auto.offset.reset": "earliest",
		}

		c, err := kafka.NewConsumer(configConsumer)
		if err != nil {
			fmt.Printf("Error creating queue Kafka consumer: %v\n", err)
			return
		}

		go o.subscribeTopic(c, subject, eventHandler)
		o.topicQueueList[subject] = true
	}
}

func (o *KafkaPubSubClient) createKafkaTopic(topic string) {
	adminClient, err := kafka.NewAdminClient(o.configMap)
	if err != nil {
		fmt.Printf("Error creating admin client: %v\n", err)
		return
	}
	defer adminClient.Close()

	// Specify the topic configuration
	topicConfig := &kafka.TopicSpecification{
		Topic:             topic,
		NumPartitions:     10, // You can adjust the number of partitions as needed.
		ReplicationFactor: 1,  // You can adjust the replication factor as needed.
	}

	ctx := context.Background()

	// Create the topic
	topics := []kafka.TopicSpecification{*topicConfig}
	_, err = adminClient.CreateTopics(ctx, topics, kafka.SetAdminOperationTimeout(5000))
	if err != nil {
		fmt.Printf("Error creating topic: %v\n", err)
	} else {
		o.topicList[topic] = true
		fmt.Println("Topic created successfully!")
	}
}

func (o *KafkaPubSubClient) checkTopic(topic string) {
	// Check if the topic exists
	if o.topicList[topic] {
		fmt.Printf("Topic '%s' exists.\n", topic)
	} else {
		fmt.Printf("Topic '%s' does not exist.\n", topic)
		o.createKafkaTopic(topic)
	}
}

func (o *KafkaPubSubClient) consumeTopic(topic string, eventHandler PubSubEventHandler) {
	c, err := kafka.NewConsumer(o.configConsumer)
	if err != nil {
		fmt.Printf("Error creating Kafka consumer: %v\n", err)
		return
	}

	defer c.Close()

	//o.checkTopic(topic)
	o.subscribeTopic(c, topic, eventHandler)
}

func (o *KafkaPubSubClient) subscribeTopic(c *kafka.Consumer, topic string, eventHandler PubSubEventHandler) {
	err := c.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		fmt.Printf("Error subscribing to topic: %v\n", err)
		defer func() {
			err := recover()

			if err != nil {
				fmt.Println("=====> SYSTEM PANIC WHEN PROCESS NATS MSG :", topic, " : ", err)
			}
		}()
	}

	fmt.Printf("Consumer for topic %s started\n", topic)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-signals:
			// Handle termination signals.
			fmt.Println("Received termination signal. Shutting down consumer.")
			return
		default:
			// Poll for messages.
			ev := c.Poll(100) // Adjust the timeout as needed.

			switch e := ev.(type) {
			case *kafka.Message:
				// Handle the Kafka message.
				fmt.Printf("Received message on topic %s: %s\n", *e.TopicPartition.Topic, string(e.Value))
				eventHandler(topic, string(e.Value))
			case kafka.Error:
				// Handle Kafka errors.
				fmt.Printf("Error while consuming Kafka message: %v\n", e)
			}
		}
	}
}

func init() {
	RegPubSubCreator("kafka", func(connString string) (PubSubClient, error) {
		ret := &KafkaPubSubClient{}

		var err error
		ret.connString = connString
		ret.configMap = &kafka.ConfigMap{"bootstrap.servers": connString}

		ret.producer, err = kafka.NewProducer(ret.configMap)
		if err != nil {
			return nil, err
		}

		//// for consumer group id should dynamic by aws configMap
		ret.configConsumer = &kafka.ConfigMap{
			"bootstrap.servers": connString,
			"group.id":          "my-consumer-group",
			"auto.offset.reset": "earliest",
		}

		mapping := make(map[string]bool)
		if ret.topicList == nil {
			ret.topicList = mapping
		}
		if ret.topicConsumeList == nil {
			ret.topicConsumeList = mapping
		}
		if ret.topicQueueList == nil {
			ret.topicQueueList = mapping
		}

		return ret, nil
	})
}
