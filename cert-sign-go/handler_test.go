package function

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
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

	req, err := http.NewRequest("POST", "/cert-sign", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Handle)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
func TestHandleRawJSON(t *testing.T) {

	payload := `{
  "Host": "example.com",
  "RSAKeySize": 2048,
  "ValidFor": 63072000000000000
}`

	req, err := http.NewRequest("POST", "/cert-sign", bytes.NewBuffer([]byte(payload)))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Handle)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	log.Println(rr.Body.String())
}
