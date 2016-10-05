package otmap

import "testing"

func TestParser(t *testing.T) {
	if _, err := Parse("C:/Users/ragga/Desktop/test.otbm", Config{
		Normal: true,
		Houses: true,
		Items:  true,
	}); err != nil {
		t.Error(err)
	}
}
