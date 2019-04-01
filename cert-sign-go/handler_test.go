package function

import (
	"encoding/json"
	"testing"
	"time"
)

func TestHandle(t *testing.T) {

	data := &CSR{
		Host:       "example.com",
		RSAKeySize: 2048,
		ValidFor:   time.Hour * 24 * 365 * 2,
	}

	payload, err := json.Marshal(&data)
	if err != nil {
		t.Errorf("error mashalling data. %v", err)
	}

	result := Handle(payload)

	if result == "" {
		t.Errorf("expected a result got nothing")
	}
}
