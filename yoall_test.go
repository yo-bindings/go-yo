package yo

import (
	"net/http"
	"testing"
)

func TestYoAll(t *testing.T) {
	setup()
	defer teardown()

	apiToken := "FOOBAR"

	mux.HandleFunc("/yoall", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Bad HTTP method: expected POST, got %v", r.Method)

			err := r.ParseForm()
			if err != nil {
				t.Errorf("Error parsing form: %v", err.Error())
			}

			if r.FormValue("api_token") != apiToken {
				t.Errorf("Bad API token: expected %v, got %v", apiToken,
					r.FormValue("api_token"))
			}
		}
	})

	// TODO: Check to make sure response is correct. Need to know what the
	// will be before testing this though.
	_, err := client.YoAll(apiToken)
	if err != nil {
		t.Errorf("YoAll returned error: %v", err.Error())
	}
}
