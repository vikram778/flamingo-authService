package twilio

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type Twilio struct {
	AccountSid string
	AuthToken  string
	To         string
	From       string
	Message    string
}

func (t *Twilio) SendOTP() (string, error) {
	urlStr := "https://api.twilio.com/2010-04-01/Accounts/" + t.AccountSid + "/Messages.json"

	// Build out the data for our message
	v := url.Values{}
	v.Set("To", t.To)
	v.Set("From", t.From)
	v.Set("Body", t.Message)
	rb := *strings.NewReader(v.Encode())

	// Create client
	client := &http.Client{}

	req, _ := http.NewRequest("POST", urlStr, &rb)
	req.SetBasicAuth(t.AuthToken, t.AuthToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Make request
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var data map[string]interface{}
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		err := json.Unmarshal(bodyBytes, &data)
		if err == nil {
			return "success", nil
		}
	}
	return "success", nil
}
