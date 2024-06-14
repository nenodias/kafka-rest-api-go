package domain

import (
	"encoding/json"
	"testing"
)

func TestJsonFormat(t *testing.T) {
	rawJson := `{"id":1,"name":"John","age":{"__type__":"long", "__int64__":30}}`
	jsonValue := JsonValue{}
	err := json.Unmarshal([]byte(rawJson), &jsonValue)
	if err != nil {
		t.Fail()
	}
	if jsonValue.Object["id"].Integer != 1 {
		t.Errorf("Expected id: %d == %d", jsonValue.Object["id"].Integer, 1)
	}
	if jsonValue.Object["name"].Text != "John" {
		t.Errorf("Expected id: %s == %s", jsonValue.Object["name"].Text, "John")
	}
	if jsonValue.Object["age"].Long != 30 {
		t.Errorf("Expected age: %d == %d", jsonValue.Object["age"].Long, 30)
	}
}
