package otmap

import (
	"log"
	"testing"
)

func TestParser(t *testing.T) {
	if m, err := Parse("G:/TFS/data/world/forgotten.otbm"); err != nil {
		t.Error(err)
	} else {
		log.Println(m.Towns)
	}
}
