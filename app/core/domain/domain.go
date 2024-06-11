package domain

type Record struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Certificate struct {
	CALocation   string `json:"ca_location"`
	CertLocation string `json:"cert_location"`
	KeyLocation  string `json:"key_location"`
	Password     string `json:"password"`
}

type PostRequest struct {
	Topic          string       `json:"topic"`
	Brokers        []string     `json:"brokers"`
	SchemaRegistry *string      `json:"schema_registry"`
	Certificate    *Certificate `json:"ssl"`

	HasKeySchema   bool     `json:"has_key_schema"`
	HasValueSchema bool     `json:"has_value_schema"`
	Records        []Record `json:"records"`
}
