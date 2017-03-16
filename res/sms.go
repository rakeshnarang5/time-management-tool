package res

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func SendMessage(msg string) {
	// Set initial variables
	accountSid := "ACb8a3ca98d75b24839c726140e9137d13"
	authToken := "6474a0ab8bdae60fcd3937e88eb18325"
	urlStr := "https://api.twilio.com/2010-04-01/Accounts/" + accountSid + "/Messages.json"

	// Build out the data for our message
	v := url.Values{}
	v.Set("To", "+919971250466")
	v.Set("From", "+12037744150")
	v.Set("Body", msg)
	rb := *strings.NewReader(v.Encode())

	// Create client
	client := &http.Client{}

	req, _ := http.NewRequest("POST", urlStr, &rb)
	req.SetBasicAuth(accountSid, authToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Make request
	resp, _ := client.Do(req)
	fmt.Println(resp.Status)
}
