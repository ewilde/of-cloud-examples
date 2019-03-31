package function

import (
	"io/ioutil"
	"net/http"
	"time"
)

// Handle a serverless request
func Handle(content []byte) string {

	req, err := http.NewRequest("GET", "https://wttr.in/"+string(content), nil)
	if err != nil {
		return err.Error()
	}

	req.Header.Add("User-Agent", "curl")

	var client = &http.Client{
		Timeout: time.Second * 10,
	}

	resp, err := client.Do(req)
	if err != nil {
		return err.Error()
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err.Error()
	}

	return string(body)
}
