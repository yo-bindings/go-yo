package yo

import (
	"net/http"
	"net/url"
)

// YoAll sends a Yo to all of your subscribers. The API token is issued by Yo
// for a specific account and is required.
func (c *Client) YoAll(apiToken string) (*http.Response, error) {
	u := "yoall/"
	body := url.Values{"api_token": {apiToken}}
	req, err := c.NewRequestForm(u, body)
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
