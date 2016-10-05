package otmap

import "encoding/binary"

type Item struct {
	ID    uint16
	Count uint16
	Attr  map[uint8]interface{}
	child []Item
}

func (i *Item) parse(node Node) error {
	if err := binary.Read(node.data, binary.LittleEndian, &i.ID); err != nil {
		return err
	}
	return nil
}
