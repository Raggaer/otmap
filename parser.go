package otmap

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
)

// Parse parses the given .otbm file
func Parse(filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	var identifier [4]byte
	if err := binary.Read(reader, binary.LittleEndian, &identifier); err != nil {
		return err
	}
	if !bytes.Equal(identifier[:4], []byte{'\x00', '\x00', '\x00', '\x00'}) && !bytes.Equal(identifier[:4], []byte{'O', 'T', 'B', 'M'}) {
		return fmt.Errorf("Unsupported file identifier: %v", identifier)
	}
	root := Node{}
	if err := root.Parse(reader); err != nil {
		return err
	}
	var property uint8
	if err := binary.Read(root.data, binary.LittleEndian, &property); err != nil {
		return err
	}
	if property != 0 {
		return fmt.Errorf("Unable to read OTBM root property. got %v exptected %v", property, 0)
	}
	var headerVersion uint32
	if err := binary.Read(root.data, binary.LittleEndian, &headerVersion); err != nil {
		return err
	}
	if headerVersion > 3 {
		return fmt.Errorf("Unsupported OTBM version. got %v expected %v", headerVersion, "> 3")
	}
	return nil
}
