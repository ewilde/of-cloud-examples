package function

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// ProviderAccessToken as issued by GitHub or GitLab
type ProviderAccessToken struct {
	AccessToken string `json:"access_token"`
}

func getUserProfile(token string) (*github.UserEmail, error) {
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Minute)
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	users, _, err := client.Users.ListEmails(ctx, nil)
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		if user.GetPrimary() {
			return user, nil
		}
	}

	return nil, fmt.Errorf("could not find primary user")
}

func getGithubToken(code string, clientID string, clientSecret string) (*ProviderAccessToken, error) {

	r, err := http.NewRequest(http.MethodPost, "https://github.com/login/oauth/access_token", nil)
	if err != nil {
		return nil, err
	}

	q := r.URL.Query()
	log.Printf("Query string:%v", q)

	q.Add("client_id", clientID)
	q.Add("client_secret", clientSecret)
	q.Add("code", code)

	r.Header.Add("Accept", "application/json")
	r.URL.RawQuery = q.Encode()

	var client = &http.Client{
		Timeout: time.Second * 10,
	}

	res, err := client.Do(r)
	if err != nil {
		return nil, err
	}

	token := ProviderAccessToken{}
	if res.Body != nil {
		defer res.Body.Close()

		tokenRes, _ := ioutil.ReadAll(res.Body)

		err := json.Unmarshal(tokenRes, &token)
		if err != nil {
			return nil, err
		}

		return &token, nil
	}

	return nil, fmt.Errorf("no body received from server")
}
