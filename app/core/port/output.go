package port

import (
	"github.com/nenodias/kafka-rest-api-go/app/core/domain"
)

type ProduceKafkaMessageOutputPort interface {
	PostOnTopic(domain.PostRequest) error
}
