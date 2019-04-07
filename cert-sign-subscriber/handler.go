package function

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

const confirmation = "SubscriptionConfirmation"
const notification = "Notification"

// Handle a serverless request
func Handle(req []byte) string {
	log.SetOutput(os.Stderr)
	var n interface{}
	err := json.Unmarshal(req, &n)
	if err != nil {
		log.Printf("unable to Unmarshal request. %v", err)
		return ""
	}

	data := n.(map[string]interface{})

	log.Println(data["Type"])
	if data["Type"].(string) == confirmation {
		subscribeURL := data["SubscribeURL"].(string)
		log.Printf("SubscribeURL %v", subscribeURL)
		confirmSubscription(subscribeURL)
		return "just subscribed to " + subscribeURL
	} else if data["Type"].(string) == notification {
		message := data["Message"].(string)
		log.Println("Received this message : ", message)
		return message
	}

	log.Printf("Unknown data type %v", data["Type"])
	return fmt.Sprintf("Unknown data type %v", data["Type"])
}

//confirmSubscription confirms the subscription by making a get request to the subscription URL
func confirmSubscription(subscriptionURL string) {
	response, err := http.Get(subscriptionURL)
	if err != nil {
		log.Printf("unable to confirm subscription")
	} else {
		log.Printf("subscription confirmed sucessfully. %d", response.StatusCode)
	}
}
