package blive

import (
	"encoding/json"
	"fmt"
	"github.com/eric2788/common-services/bilibili"
	"net/http"
	"os"
)

func GetRoomInfo(room int64) (info *bilibili.RoomInfo, err error) {
	info, err = bilibili.GetRoomInfo(room)
	if info.Code != 0 {
		return nil, fmt.Errorf("%s", info.Message)
	}
	return
}

func GetUserInfo(uid int64) (info *bilibili.UserInfo, err error) {
	info, err = bilibili.GetUserInfo(uid)
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
