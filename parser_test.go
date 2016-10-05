package otmap

import (
	"log"
	"testing"
)

func TestParser(t *testing.T) {
	if m, err := Parse("C:/Users/ragga/Desktop/test.otbm", Config{
		Towns:  true,
		Houses: true,
		Items:  true,
	}); err != nil {
		t.Error(err)
	} else {
		log.Println(m.Houses[0].GenerateMinimapImage("test"))
	}
}
