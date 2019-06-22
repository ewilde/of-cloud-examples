package function

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func Handle(req []byte) string {
	body := req
	csr := &CSR{}
	err := json.Unmarshal(body, csr)
	if err != nil {
		handleError(http.StatusInternalServerError, fmt.Errorf("error unmarshalling body %s. %v", body, err))
		return ""
	}

	key, cert, err := Sign(csr)
	if err != nil {
		handleError(http.StatusInternalServerError, err)
		return ""
	}

	return fmt.Sprintf("%s\n%s\n", key, cert)
}

func handleError(status int, err error) {
	log.Printf("status: %v", status)
	log.Println(err.Error())
}
