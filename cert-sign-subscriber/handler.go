package function

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const confirmation = "SubscriptionConfirmation"
const notification = "Notification"

// Handle a serverless request
func Handle(req []byte) string {
	var n interface{}
	err := json.Unmarshal(req, &n)
	if err != nil {
		fmt.Printf("unable to Unmarshal request. %v", err)
	}

	data := n.(map[string]interface{})

	if data["Type"].(string) == confirmation {
		subscribeURL := data["SubscribeURL"].(string)
		go confirmSubscription(subscribeURL)
		return "just subscribed to " + subscribeURL
	} else if data["Type"].(string) == notification {
		message := data["Message"].(string)
		fmt.Println("Received this message : ", message)
		return message
	}

	return fmt.Sprintf("Unknown data type %v", data["Type"])
}

//confirmSubscription confirms the subscription by making a get request to the subscription URL
func confirmSubscription(subscriptionURL string) {
	response, err := http.Get(subscriptionURL)
	if err != nil {
		fmt.Printf("unable to confirm subscription")
	} else {
		fmt.Printf("subscription confirmed sucessfully. %d", response.StatusCode)
	}
}
