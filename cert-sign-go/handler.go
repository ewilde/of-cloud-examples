package function

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func Handle(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		handleError(w, http.StatusBadRequest, errors.New("body containing CSR must be included"))
		return
	}

	defer r.Body.Close()

	body, _ := ioutil.ReadAll(r.Body)
	csr := &CSR{}
	err := json.Unmarshal(body, csr)
	if err != nil {
		handleError(w, http.StatusInternalServerError, fmt.Errorf("error unmarshalling body %s. %v", body, err))
		return
	}

	cert, key, err := Sign(csr)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(cert))
	w.Write([]byte(key))
}

func handleError(w http.ResponseWriter, status int, err error) {
	log.Println(err)
	w.WriteHeader(status)
	w.Write([]byte(err.Error()))
}
