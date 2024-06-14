package domain

import (
	"encoding/json"
	"math"
	"strconv"
	"strings"

	"github.com/linkedin/goavro/v2"
)

type JsonValue struct {
	Text    string               `json:"__text__,omitempty"`
	Object  map[string]JsonValue `json:"__object__,omitempty"`
	Array   []JsonValue          `json:"__array__,omitempty"`
	Type    string               `json:"__type__,omitempty"`
	Integer int32                `json:"__int32__,omitempty"`
	Long    int64                `json:"__int64__,omitempty"`
	Double  float64              `json:"__float64__,omitempty"`
	Boolean bool                 `json:"__bool__,omitempty"`
	Pointer bool                 `json:"__pointer__,omitempty"`
}

func (j *JsonValue) ToJsonMap() (map[string]interface{}, error) {
	result := make(map[string]interface{})
	if j.Type == "object" {
		for key, value := range j.Object {
			var item interface{}
			if value.Type == "object" {
				dados, err := value.ToJsonMap()
				if err != nil {
					return nil, err
				}
				item = dados
			} else if value.Type == "array" {
				var items []interface{}
				item = items
				for _, v := range value.Array {
					dados, err := v.ToJsonMap()
					if err != nil {
						return nil, err
					}
					items = append(items, dados)
				}
			} else if value.Type == "string" {
				item = value.Text
			} else if value.Type == "double" {
				item = value.Double
			} else if value.Type == "integer" {
				item = value.Integer
			} else if value.Type == "long" {
				item = value.Long
			} else if value.Type == "bool" {
				item = value.Boolean
			} else {
				continue
			}
			result[key] = item
			if value.Pointer {
				result[key] = goavro.Union(value.Type, item)
			}
		}
	}
	return result, nil
}

func (j *JsonValue) UnmarshalJSON(data []byte) error {
	var dados string = string(data)
	if strings.Contains(dados, "{") {
		j.Type = "object"
		dados := make(map[string]JsonValue)
		err := json.Unmarshal(data, &dados)
		if err != nil {
			return err
		}
		if tipo, ok := dados["__type__"]; ok {
			j.Type = tipo.Text
			switch j.Type {
			case "long":
				j.Long = int64(dados["__int64__"].Integer)
			}
			if _, ok := dados["__pointer__"]; ok {
				j.Pointer = true
			}
		} else {
			j.Object = dados
		}
	} else if strings.Contains(dados, "[") {
		j.Type = "array"
		return json.Unmarshal(data, &j.Array)
	} else if strings.Contains(dados, "\"") {
		j.Type = "string"
		j.Text = strings.ReplaceAll(dados, "\"", "")
	} else if strings.Contains(dados, ".") {
		j.Type = "double"
		return json.Unmarshal(data, &j.Double)
	} else if dados == "true" || dados == "false" {
		j.Type = "bool"
		j.Boolean = dados == "true"
	} else if dados == "null" {
		j.Type = "null"
	} else {
		longValue, err := strconv.ParseInt(dados, 10, 64)
		if err != nil {
			return err
		}

		if longValue <= math.MaxInt32 {
			j.Type = "integer"
			j.Integer = int32(longValue)
		} else {
			j.Type = "long"
			j.Long = longValue
		}
	}
	return nil
}

type Record struct {
	Key   JsonValue `json:"key"`
	Value JsonValue `json:"value"`
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
