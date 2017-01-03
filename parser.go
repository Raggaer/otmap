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
	Root              Node
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
	Tiles             []Tile
	Houses            []*House
}

func (m *Map) getHouse(id uint32) *House {
	for _, house := range m.Houses {
		if house.ID == id {
			return house
		}
	}
	return nil
}

func (m *Map) addHouse(house *House) {
	m.Houses = append(m.Houses, house)
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
	currentMap.Root = Node{}
	if err := currentMap.Root.Parse(reader); err != nil {
		return nil, err
	}
	if err := binary.Read(currentMap.Root.data, binary.LittleEndian, &currentMap.Property); err != nil {
		return nil, err
	}
	if currentMap.Property != 0 {
		return nil, fmt.Errorf("Unable to read OTBM root property. got %v exptected %v", currentMap.Property, 0)
	}
	if err := binary.Read(currentMap.Root.data, binary.LittleEndian, &currentMap.Header); err != nil {
		return nil, err
	}
	if currentMap.Header > 3 {
		return nil, fmt.Errorf("Unsupported OTBM version. got %v expected %v", currentMap.Header, "> 3")
	}
	if err := binary.Read(currentMap.Root.data, binary.LittleEndian, &currentMap.Width); err != nil {
		return nil, err
	} else if err = binary.Read(currentMap.Root.data, binary.LittleEndian, &currentMap.Height); err != nil {
		return nil, err
	} else if err = binary.Read(currentMap.Root.data, binary.LittleEndian, &currentMap.MajorItemsVersion); err != nil {
		return nil, err
	} else if err = binary.Read(currentMap.Root.data, binary.LittleEndian, &currentMap.MinorItemsVersion); err != nil {
		return nil, err
	}
	mapData := currentMap.Root.child[0]
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
	for _, node := range mapData.child {
		var nodeType uint8
		if err := binary.Read(node.data, binary.LittleEndian, &nodeType); err != nil {
			return nil, err
		}
		if nodeType == OTBMNodeTowns {
			if err := currentMap.parseTowns(node); err != nil {
				return nil, err
			}
		} else if nodeType == OTBMNodeTileArea {
			basePosition := Position{}
			if err := binary.Read(node.data, binary.LittleEndian, &basePosition); err != nil {
				return nil, err
			}
			for _, tileNode := range node.child {
				var nodeType uint8
				if err := binary.Read(tileNode.data, binary.LittleEndian, &nodeType); err != nil {
					return nil, err
				}
				tile := Tile{}
				var x, y uint8
				if err := binary.Read(tileNode.data, binary.LittleEndian, &x); err != nil {
					return nil, err
				} else if err = binary.Read(tileNode.data, binary.LittleEndian, &y); err != nil {
					return nil, err
				}
				tile.Position.X = uint16(x) + basePosition.X
				tile.Position.Y = uint16(y) + basePosition.Y
				tile.Position.Z = basePosition.Z
				if nodeType == OTBMNodeHouseTile {
					if err := currentMap.parseHouseTile(basePosition, tileNode, tile); err != nil {
						return nil, err
					}
				}
				for i := 0; i < tileNode.data.Len(); i++ {
					var attr uint8
					if err := binary.Read(tileNode.data, binary.LittleEndian, &attr); err != nil {
						return nil, err
					}
					switch attr {
					case OTBMAttrItem:
						item := Item{}
						if err := binary.Read(tileNode.data, binary.LittleEndian, &item.ID); err != nil {
							return nil, err
						}
						tile.Items = append(tile.Items, item)
					case OTBMAttrTileFlags:
						if err := binary.Read(tileNode.data, binary.LittleEndian, &tile.flags); err != nil {
							return nil, err
						}
					default:
						return nil, fmt.Errorf("Unkown tile attribute. got %v", attr)
					}
				}
				for _, itemNode := range tileNode.child {
					var nodeType uint8
					if err := binary.Read(itemNode.data, binary.LittleEndian, &nodeType); err != nil {
						return nil, err
					}
					if nodeType != OTBMNodeItem {
						return nil, fmt.Errorf("Wrong node type. expected %v got %v", OTBMNodeItem, nodeType)
					}
					item := Item{}
					if err := item.parse(itemNode); err != nil {
						return nil, err
					}
					tile.Items = append(tile.Items, item)
				}
				currentMap.Tiles = append(currentMap.Tiles, tile)
			}
		}
	}
	return currentMap, nil
}

func (m *Map) parseHouseTile(base Position, node Node, tile Tile) error {
	var id uint32
	if err := binary.Read(node.data, binary.LittleEndian, &id); err != nil {
		return err
	}
	currentHouse := &House{
		ID: id,
	}
	if house := m.getHouse(id); house != nil {
		currentHouse = house
	} else {
		m.addHouse(currentHouse)
	}
	currentHouse.Tiles = append(currentHouse.Tiles, tile)
	return nil
}

func (m *Map) parseTowns(node Node) error {
	towns := []Town{}
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
		m.Towns = append(towns, currentTown)
	}
	return nil
}
