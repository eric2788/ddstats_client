package blive

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

func doRequest(method string, urlString string, rooms []string) (*http.Response, error) {

	var body io.Reader

	if rooms != nil {

		form := url.Values{
			"subscribes": rooms,
		}

		body = strings.NewReader(form.Encode())
	} else {
		body = nil
	}

	req, err := http.NewRequest(method, urlString, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "ddstats_client")

	resp, err := http.DefaultClient.Do(req)
	return resp, err
}

func GetSubscribes() ([]int64, error) {
	var rooms []int64

	resp, err := doRequest(http.MethodGet, fmt.Sprintf("https://%s/subscribe", Host), nil)
	if err != nil {
		return nil, err
	}

	err = respAs(resp, &rooms)
	return rooms, err
}

func PutSubscribe(rooms []string, add bool) error {

	path := "add"

	if !add {
		path = "remove"
	}

	res, err := doRequest(http.MethodPut, fmt.Sprintf("https://%s/subscribe/%s", Host, path), rooms)
	if err == nil {
		defer res.Body.Close()
	}
	return err
}

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

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf(resp.Status)
	}

	return nil
}

func SubscribeFromOffline() {
	if rooms, ok := GetFromOffline(); ok {
		if err := SubscribeRequest(rooms); err != nil {
			logrus.Errorf("尝试从离线重新订阅时出现错误: %v", err)
		} else {
			logrus.Infof("成功从离线重新订阅 %v 个房间。", len(rooms))
		}
	}
}

func TrackSubscribes(ctx context.Context) {
	ticker := time.NewTicker(time.Minute * 5)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			rooms, ok := GetFromOffline()
			if !ok {
				return
			}
			listening, err := GetSubscribes()
			if err != nil {
				logrus.Error(err)
				return
			}
			if len(rooms) > len(listening) {
				logrus.Infof("离线订阅数量 %d 与当前订阅数量 %d 不一致，尝试重新订阅。", len(rooms), len(listening))
				SubscribeFromOffline()
				return
			}
		case <-ctx.Done():
			return
		}
	}
}
