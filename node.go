package otmap

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
)

const (
	EscapeChar = 0xFD
	NodeStart  = 0xFE
	NodeEnd    = 0xFF
)

// Node struct used for binary tree nodes
type Node struct {
	data  *bytes.Buffer
	child []Node
}

// ReadString reads a string from the buffer
func (n *Node) ReadString() (string, error) {
	var length uint16
	if err := binary.Read(n.data, binary.LittleEndian, &length); err != nil {
		return "", err
	}
	result := make([]byte, length)
	if err := binary.Read(n.data, binary.LittleEndian, &result); err != nil {
		return "", err
	}
	return string(result), nil
}

// Parse parses and creates a new node
func (n *Node) Parse(reader *bufio.Reader) error {
	n.data = &bytes.Buffer{}
	if startByte, err := reader.ReadByte(); err != nil {
		return err
	} else if startByte != NodeStart {
		return fmt.Errorf("Expected %v got %v", NodeStart, startByte)
	}
	return n.unserialize(reader)
}

func (n *Node) unserialize(reader *bufio.Reader) error {
	for {
		b, err := reader.ReadByte()
		if err != nil {
			return err
		}
		switch b {
		case NodeStart:
			newNode := Node{
				data: &bytes.Buffer{},
			}
			if err := newNode.unserialize(reader); err != nil {
				return err
			}
			n.child = append(n.child, newNode)
		case NodeEnd:
			return nil
		case EscapeChar:
			next, err := reader.ReadByte()
			if err != nil {
				return err
			}
			n.data.WriteByte(next)
		default:
			n.data.WriteByte(b)
		}
	}
	return nil
}
