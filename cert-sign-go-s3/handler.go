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

	keyFileName, err := saveToS3(key, fmt.Sprintf("%s-key", csr.Host))
	if err != nil {
		log.Printf("error saving key %s. %v", csr.Host, err)
	}

	certFileName, err := saveToS3(cert, fmt.Sprintf("%s-cert", csr.Host))
	if err != nil {
		log.Printf("error saving certificate %s. %v", csr.Host, err)
	}

	return fmt.Sprintf("%s\t%s", keyFileName, certFileName)
}

func handleError(status int, err error) {
	log.Printf("status: %v", status)
	log.Println(err.Error())
}
