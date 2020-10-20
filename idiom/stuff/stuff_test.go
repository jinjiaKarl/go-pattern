package stuff

import "testing"

func TestNewStuffClient(t *testing.T) {
	var _ = []struct {
		in       string
		expected string
	}{
		{"test", "test+10+2"},
	}
	stuff := NewStuffClient("test")
	stuff.DoStuff()
}

func ExampleNewStuffClient() {
	stuff := NewStuffClient("test")
	stuff.DoStuff()
	//Ouput:
	//test+10+2
}
