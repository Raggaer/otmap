package otmap

import (
	"log"
	"testing"
)

func TestParser(t *testing.T) {
	if m, err := Parse("C:/Users/ragga/Desktop/test.otbm"); err != nil {
		t.Error(err)
	} else {
		m.ParseTowns()
		log.Println(m.Towns)
		m.ParseHouses()
	}
}
