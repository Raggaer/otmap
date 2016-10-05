package otmap

type Item struct {
	ID    uint16
	Count uint16
	Attr  map[uint8]interface{}
	child []Item
}
