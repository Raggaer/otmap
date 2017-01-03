package otmap

import (
	"log"
	"testing"
)

func TestParser(t *testing.T) {
	if m, err := Parse("C:/Users/ragga/Desktop/test.otbm"); err != nil {
		t.Error(err)
	} else {
		log.Println(m.Houses[0].GenerateMinimapImage("test"))
	}
}
