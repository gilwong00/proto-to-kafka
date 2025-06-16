package kafka

type KafkaClient struct{}

func NewClient() *KafkaClient {
	return &KafkaClient{}
}

func (k *KafkaClient) publishToKafka() {}
