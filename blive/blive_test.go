package blive

import (
	"testing"
)


func TestGetUserName(t *testing.T)() {
	name, err := GetUserName(27062023)
	if err != nil {
		t.Skip(err)
	}
	t.Log("result: ", name)
}