package otmap

import "testing"

func TestParser(t *testing.T) {
	if _, err := Parse("C:/Users/ragga/Desktop/test.otbm", Config{
		Towns:  true,
		Houses: true,
		Items:  true,
	}); err != nil {
		t.Error(err)
	} else {

	}
}
