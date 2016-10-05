package otmap

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
)

// Map struct used to save all map information
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
	Towns             []Town
}

// Parse parses the given .otbm file
func Parse(filepath string, options Config) (*Map, error) {
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
	if err := currentMap.parse(mapData); err != nil {
		return nil, err
	}
	return currentMap, nil
}

func (m *Map) parse(mapData Node) error {
	for _, node := range mapData.child {
		var nodeType uint8
		if err := binary.Read(node.data, binary.LittleEndian, &nodeType); err != nil {
			return err
		}
		if nodeType == OTBMNodeTowns {
			for _, town := range node.child {
				var nodeType uint8
				if err := binary.Read(town.data, binary.LittleEndian, &nodeType); err != nil {
					return err
				}
				if nodeType != OTBMNodeTown {
					return fmt.Errorf("Parsing map towns: expected %v got %v", OTBMNodeTown, nodeType)
				}
				currentTown := Town{}
				if err := binary.Read(town.data, binary.LittleEndian, &currentTown.ID); err != nil {
					return err
				} else if currentTown.Name, err = town.ReadString(); err != nil {
					return err
				} else if currentTown.TemplePosition, err = town.ReadPosition(); err != nil {
					return err
				}
				m.Towns = append(m.Towns, currentTown)
			}
		}
	}
	return nil
}
