package function

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/openfaas/openfaas-cloud/sdk"
)

// Handle a serverless request
func Handle(w http.ResponseWriter, r *http.Request) {
	log.SetOutput(os.Stderr)

	clientID, err := sdk.ReadSecret("client_id")
	if err != nil {
		newHttpInternalServerErrorResponse(w, fmt.Errorf("error loading client_id secret. %v", err))
		return
	}

	clientSecret, err := sdk.ReadSecret("client_secret")
	if err != nil {
		newHttpInternalServerErrorResponse(w, fmt.Errorf("error loading client_secret secret. %v", err))
		return
	}

	token, err := getGithubToken(r.URL.Query().Get("code"), clientID, clientSecret)
	if err != nil {
		newHttpInternalServerErrorResponse(w, fmt.Errorf("error getting github token. %v", err))
		return
	}

	user, err := getUserProfile(token.AccessToken)
	if err != nil {
		newHttpInternalServerErrorResponse(w, err)
		return
	}

	jwt, err := createToken(user.GetEmail())
	if err != nil {
		newHttpInternalServerErrorResponse(w, err)
		return
	}

	newHttpOkResponse(w, jwt)
}

func newHttpResponse(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	w.Write([]byte(message))
}

func newHttpOkResponse(w http.ResponseWriter, body string) {
	newHttpResponse(w, http.StatusOK, body)
}

func newHttpInternalServerErrorResponse(w http.ResponseWriter, err error) {
	newHttpResponse(w, http.StatusInternalServerError, err.Error())
}
