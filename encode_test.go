package sse

import "testing"

func TestEncode(t *testing.T) {
	f := func(id, event, data string, want string) {
		t.Helper()

		gotBytes := encode(id, event, []byte(data))
		got := string(gotBytes)

		if got != want {
			t.Errorf("got %#v, want %#v", got, want)
		}
	}

	f(
		"first", "message", "go test",
		"event:message\ndata:go test\n\n",
	)
	f(
		"first", "message", "",
		"event:message\ndata\n\n",
	)
	f(
		"first", "", "empty event",
		"data:empty event\n\n",
	)
	f(
		"first", "", "",
		"data\n\n",
	)
}
