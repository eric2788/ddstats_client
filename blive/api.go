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

func httpGetAs(url string, as interface{}) (err error) {
	req, err := http.Get(url)
	if err != nil {
		return err
	}
	defer req.Body.Close()
	if req.StatusCode != 200 {
		return fmt.Errorf("http status code: %d", req.StatusCode)
	}
	return json.NewDecoder(req.Body).Decode(as)
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
	err = json.Unmarshal(content, &rooms)
	if err != nil {
		return nil, false
	}
	return rooms, true
}
