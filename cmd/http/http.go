package http

import (
	"log"
	"net/http"

	"github.com/nenodias/kafka-rest-api-go/app/core/usecase"
	"github.com/nenodias/kafka-rest-api-go/app/infraestructure/kafka"
	"github.com/nenodias/kafka-rest-api-go/app/infraestructure/web"
)

func Start() {
	mux := http.NewServeMux()
	usecase := usecase.NewTopicProducerUsecase(kafka.NewAppKafkaProducer())
	handler := web.NewTopicProducerHandler(usecase)
	mux.HandleFunc("POST /topic", handler.PostOnTopic)
	log.Fatalln(http.ListenAndServe(":8080", mux))
}
