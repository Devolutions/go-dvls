package dvls

import (
	"testing"
)

func Test_Server(t *testing.T) {
	t.Run("GetServerInfo", test_GetServerInfo)
	t.Run("GetTimezones", test_GetTimezones)
}

func test_GetServerInfo(t *testing.T) {
	_, err := testClient.GetServerInfo()
	if err != nil {
		t.Fatal(err)
	}
}

func test_GetTimezones(t *testing.T) {
	_, err := testClient.GetServerTimezones()
	if err != nil {
		t.Fatal(err)
	}
}
