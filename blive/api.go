package blive

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func GetRoomInfo(room int64) (info *RoomInfo, err error) {
	info = &RoomInfo{}
	err = httpGetAs(fmt.Sprintf("https://api.live.bilibili.com/room/v1/Room/get_info?room_id=%v", room), info)
	if info.Code != 0 {
		return nil, fmt.Errorf("%s", info.Message)
	}
	return
}

func GetUserInfo(uid int64) (info *UserInfo, err error) {
	info = &UserInfo{}
	err = httpGetAs(fmt.Sprintf("https://api.bilibili.com/x/space/acc/info?mid=%v&jsonp=jsonp", uid), info)
	if info.Code != 0 {
		return nil, fmt.Errorf("%s", info.Message)
	}
	return
}

func respAs(resp *http.Response, v interface{}) error {
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	return dec.Decode(v)
}

func httpGetAs(url string, as interface{}) (err error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.114 Safari/537.36")
	req.Header.Set("Origin", "https://live.bilibili.com")
	req.Header.Set("Referer", "https://live.bilibili.com/")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return fmt.Errorf("http status code: %d", res.StatusCode)
	}
	return json.NewDecoder(res.Body).Decode(as)
}

func GetUserName(room int64) (name string, err error) {
	info, err := GetRoomInfo(room)
	if err != nil {
		return
	}
	uid := info.Data.Uid
	user, err := GetUserInfo(uid)
	if err != nil {
		return
	}
	name = user.Data.Name
	return
}

func GetFromOffline() ([]string, bool) {
	content, err := os.ReadFile("data/offline.json")
	if err != nil {
		return nil, false
	}
	var rooms []string
	if len(content) == 0 || string(content) == "[]" {
		return nil, false
	}
	err = json.Unmarshal(content, &rooms)
	if err != nil {
		return nil, false
	} else if len(rooms) == 0 || rooms[0] == "undefined" {
		return nil, false
	}
	return rooms, true
}
