package yo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	libraryVersion = "0.1"
	baseURL        = "http://api.justyo.co/"
	userAgent      = "go-yo/" + libraryVersion
)

type Client struct {
	client    *http.Client
	BaseURL   *url.URL
	UserAgent string
}

// NewClient returns a new Yo API client. If client is nil, then
// http.DefaultClient is used.
func NewClient(client *http.Client) *Client {
	if client == nil {
		client = http.DefaultClient
	}

	baseURL, _ := url.Parse(baseURL)

	return &Client{
		client:    client,
		BaseURL:   baseURL,
		UserAgent: userAgent,
	}
}

// NewRequest creates a new API request. urlStr can either be an relative or
// absolute url. If it's relative, the it's resolved to the BaseURL of the
// Client. If body is provided, then it is encoded as JSON and included in the
// request body.
func (c *Client) NewRequest(method, urlStr string,
	body interface{}) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	u := c.BaseURL.ResolveReference(rel)

	buf := new(bytes.Buffer)
	if body != nil {
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", c.UserAgent)
	return req, nil
}

func (c *Client) Do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	err = CheckResponse(resp)
	if err != nil {
		return nil, err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			io.Copy(w, resp.Body)
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
		}
	}

	return resp, err
}

// ErrorResponse reports one or more errors returned from an API request.
//
// TODO: Acutally report one or more errors instead of just have the response.
// Can't currently do this because I don't have an API token.
type ErrorResponse struct {
	Response *http.Response
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d", r.Response.Request.Method,
		r.Response.Request.URL, r.Response.StatusCode)
}

// CheckResponse checks the response for errors. A response is considered an
// error if it is outside the 200 range.
func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}
	// TODO: Unmarshal the error response body. Don't have an API token yet, so
	// no idea what these look like.
	return &ErrorResponse{Response: r}
}
