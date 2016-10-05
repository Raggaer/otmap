package otmap

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
)

type Map struct {
	Width             uint16
	Height            uint16
	Header            uint32
	Property          uint8
	MajorItemsVersion uint32
	MinorItemsVersion uint32
	Description       string
	SpawnFile         string
	HouseFile         string
}

// Parse parses the given .otbm file
func Parse(filepath string) (*Map, error) {
	currentMap := &Map{}
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	var identifier [4]byte
	if err := binary.Read(reader, binary.LittleEndian, &identifier); err != nil {
		return nil, err
	}
	if !bytes.Equal(identifier[:4], []byte{'\x00', '\x00', '\x00', '\x00'}) && !bytes.Equal(identifier[:4], []byte{'O', 'T', 'B', 'M'}) {
		return nil, fmt.Errorf("Unsupported file identifier: %v", identifier)
	}
	root := Node{}
	if err := root.Parse(reader); err != nil {
		return nil, err
	}
	if err := binary.Read(root.data, binary.LittleEndian, &currentMap.Property); err != nil {
		return nil, err
	}
	if currentMap.Property != 0 {
		return nil, fmt.Errorf("Unable to read OTBM root property. got %v exptected %v", currentMap.Property, 0)
	}
	if err := binary.Read(root.data, binary.LittleEndian, &currentMap.Header); err != nil {
		return nil, err
	}
	if currentMap.Header > 3 {
		return nil, fmt.Errorf("Unsupported OTBM version. got %v expected %v", currentMap.Header, "> 3")
	}
	if err := binary.Read(root.data, binary.LittleEndian, &currentMap.Width); err != nil {
		return nil, err
	} else if err = binary.Read(root.data, binary.LittleEndian, &currentMap.Height); err != nil {
		return nil, err
	} else if err = binary.Read(root.data, binary.LittleEndian, &currentMap.MajorItemsVersion); err != nil {
		return nil, err
	} else if err = binary.Read(root.data, binary.LittleEndian, &currentMap.MinorItemsVersion); err != nil {
		return nil, err
	}
	mapData := root.child[0]
	var nodeType uint8
	if err = binary.Read(mapData.data, binary.LittleEndian, &nodeType); err != nil {
		return nil, err
	}
	if nodeType != OTBMNodeMapData {
		return nil, fmt.Errorf("Expected %v got %v", OTBMNodeMapData, nodeType)
	}
	for i := 0; i < mapData.data.Len(); i++ {
		var attribute uint8
		if err = binary.Read(mapData.data, binary.LittleEndian, &attribute); err != nil {
			return nil, err
		}
		tmp, err := mapData.ReadString()
		if err != nil {
			return nil, err
		}
		switch attribute {
		case OTBMAttrDescription:
			currentMap.Description += tmp
		case OTBMAttrHouseFile:
			currentMap.HouseFile += tmp
		case OTBMAttrSpawnFile:
			currentMap.SpawnFile += tmp
		default:
			return nil, fmt.Errorf("Unkown attribute. expected %v got %v", []int{OTBMAttrDescription, OTBMAttrHouseFile, OTBMAttrSpawnFile}, attribute)
		}
	}
	return currentMap, nil
}
