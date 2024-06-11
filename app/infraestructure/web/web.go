package web

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/nenodias/kafka-rest-api-go/app/core/domain"
	"github.com/nenodias/kafka-rest-api-go/app/core/port"
)

func NewTopicProducerHandler(input port.TopicProducerInputPort) *TopicProducerHandler {
	return &TopicProducerHandler{
		TopicProducerInputPort: input,
	}
}

type TopicProducerHandler struct {
	TopicProducerInputPort port.TopicProducerInputPort
}

func (t *TopicProducerHandler) PostOnTopic(w http.ResponseWriter, r *http.Request) {
	input := domain.PostRequest{}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = t.TopicProducerInputPort.PostOnTopic(context.Background(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
