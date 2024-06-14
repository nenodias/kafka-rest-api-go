package kafka

import (
	"encoding/binary"
	"fmt"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry"

	"github.com/hamba/avro/v2"
	"github.com/linkedin/goavro/v2"
	"github.com/nenodias/kafka-rest-api-go/app/core/domain"
)

const (
	IS_GOAVRO = true
)

type AppKafkaProducer struct{}

func NewAppKafkaProducer() *AppKafkaProducer {
	return &AppKafkaProducer{}
}

func (k *AppKafkaProducer) PostOnTopic(input domain.PostRequest) error {
	var producer *kafka.Producer = nil
	var err error = nil
	brokers := strings.Join(input.Brokers, ",")

	if input.Certificate != nil {
		producer, err = kafka.NewProducer(&kafka.ConfigMap{
			"bootstrap.servers":        brokers,
			"security.protocol":        "SSL",
			"ssl.ca.location":          input.Certificate.CALocation,
			"ssl.certificate.location": input.Certificate.CertLocation,
			"ssl.key.location":         input.Certificate.KeyLocation,
			"ssl.key.password":         input.Certificate.Password,
		})
	} else {
		producer, err = kafka.NewProducer(&kafka.ConfigMap{
			"bootstrap.servers": brokers,
		})
	}

	if err != nil {
		return err
	}
	defer producer.Close()

	config := schemaregistry.NewConfig(*input.SchemaRegistry)
	client, err := schemaregistry.NewClient(config)
	if err != nil {
		fmt.Println("Error creating client")
	}
	err = k.Producer(producer, client, input)
	if err != nil {
		return err
	}

	producer.Flush(1000)
	return nil
}

func (k *AppKafkaProducer) Producer(producer *kafka.Producer, client schemaregistry.Client, input domain.PostRequest) error {
	var keySchema schemaregistry.SchemaMetadata
	var valueSchema schemaregistry.SchemaMetadata
	var err error
	if input.HasKeySchema {
		keySchema, err = client.GetLatestSchemaMetadata(input.Topic + "-key")
		if err != nil {
			return err
		}
	}
	if input.HasValueSchema {
		valueSchema, err = client.GetLatestSchemaMetadata(input.Topic + "-value")
		if err != nil {
			return err
		}
	}
	for _, record := range input.Records {

		var key []byte
		if keySchema.ID != 0 {
			key, err = k.Serialize(keySchema, record.Key)
			if err != nil {
				return err
			}
		} else {
			key = []byte(record.Key.Text)
		}
		var payload []byte
		if valueSchema.ID != 0 {
			dados, err := record.Value.ToJsonMap()
			if err != nil {
				return err
			}
			payload, err = k.Serialize(valueSchema, dados)
			if err != nil {
				return err
			}
		} else {
			payload = []byte(record.Value.Text)
		}

		err := producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &input.Topic, Partition: kafka.PartitionAny},
			Value:          payload,
			Key:            key,
		}, nil)

		if err != nil {
			return err
		}
	}
	return nil
}

func (k *AppKafkaProducer) Serialize(ser schemaregistry.SchemaMetadata, rawText interface{}) ([]byte, error) {
	if IS_GOAVRO {
		codec, err := goavro.NewCodec(ser.Schema)
		if err != nil {
			return nil, err
		}

		// Convert native Go form to binary Avro data
		binaryData, err := codec.BinaryFromNative(nil, rawText)
		if err != nil {
			return nil, err
		}

		magicByte := byte(0x0)
		messageBytes := make([]byte, 5+len(binaryData))
		messageBytes[0] = magicByte
		schemaID := uint32(ser.ID)
		binary.BigEndian.PutUint32(messageBytes[1:], schemaID)
		messageBytes = append(messageBytes[0:5], binaryData...)
		return messageBytes, nil
	} else {
		schema, err := avro.Parse(ser.Schema)
		if err != nil {
			return nil, err
		}
		return avro.Marshal(schema, rawText)
	}
}
