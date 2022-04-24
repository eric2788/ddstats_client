package blive

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func SubscribeRequest(room []string) error {

	httpUrl := url.URL{
		Host:     Host,
		Path:     "/subscribe",
		RawQuery: "validate=false",
		Scheme:   "https",
	}

	form := url.Values{
		"subscribes": room,
	}

	body := strings.NewReader(form.Encode())
	req, err := http.NewRequest(http.MethodPost, httpUrl.String(), body)

	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "ddstats_client")

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf(resp.Status)
	}

	return nil
}
