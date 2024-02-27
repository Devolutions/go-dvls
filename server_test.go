package dvls

import (
	"testing"
)

func Test_Server(t *testing.T) {
	t.Run("GetPublicServerInfo", test_GetPublicServerInfo)
	t.Run("GetPrivateServerInfo", test_GetPrivateServerInfo)
	t.Run("GetTimezones", test_GetTimezones)
}

func test_GetPublicServerInfo(t *testing.T) {
	_, err := testClient.GetPublicServerInfo()
	if err != nil {
		t.Fatal(err)
	}
}

func test_GetPrivateServerInfo(t *testing.T) {
	_, err := testClient.GetPrivateServerInfo()
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
