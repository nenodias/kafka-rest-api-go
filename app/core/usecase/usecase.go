package usecase

import (
	"context"

	"github.com/nenodias/kafka-rest-api-go/app/core/domain"
	"github.com/nenodias/kafka-rest-api-go/app/core/port"
)

func NewTopicProducerUsecase(kafka port.ProduceKafkaMessageOutputPort) *TopicProducerUsecase {
	return &TopicProducerUsecase{kafka}
}

type TopicProducerUsecase struct {
	kafka port.ProduceKafkaMessageOutputPort
}

func (t *TopicProducerUsecase) PostOnTopic(ctx context.Context, input domain.PostRequest) error {
	select {
	case <-ctx.Done():
		return nil
	default:
		return t.kafka.PostOnTopic(input)
	}
}
