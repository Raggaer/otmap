package otmap

import "testing"

func TestParser(t *testing.T) {
	if _, err := Parse("C:/Users/ragga/Desktop/test.otbm"); err != nil {
		t.Error(err)
	} else {

	}
}
