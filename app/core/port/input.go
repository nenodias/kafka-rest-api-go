package port

import (
	"context"

	"github.com/nenodias/kafka-rest-api-go/app/core/domain"
)

type TopicProducerInputPort interface {
	PostOnTopic(context.Context, domain.PostRequest) error
}
