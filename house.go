package otmap

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
)

type House struct {
	ID    uint32
	Tiles []Tile
}

func HasLevel(levels []uint8, l uint8) bool {
	for _, level := range levels {
		if level == l {
			return true
		}
	}
	return false
}

// GenerateMinimapImage generates minimap image of the given house
func (h House) GenerateMinimapImage(savepath string) error {
	houseLevels := []uint8{}
	for _, tile := range h.Tiles {
		if !HasLevel(houseLevels, tile.Position.Z) {
			houseLevels = append(houseLevels, tile.Position.Z)
		}
	}
	for _, level := range houseLevels {
		var houseTopX, houseTopY, houseBottomX, houseBottomY uint16
		for _, tile := range h.Tiles {
			if tile.Position.X < houseTopX || houseTopX == 0 {
				houseTopX = tile.Position.X
			}
			if tile.Position.Y < houseTopY || houseTopY == 0 {
				houseTopY = tile.Position.Y
			}
			if tile.Position.X > houseBottomX || houseBottomX == 0 {
				houseBottomX = tile.Position.X
			}
			if tile.Position.Y > houseBottomY || houseBottomY == 0 {
				houseBottomY = tile.Position.Y
			}
		}
		houseImage := image.NewRGBA(image.Rect(0, 0, int(houseBottomX-houseTopX)+3, int(houseBottomY-houseTopY)+3))
		for _, tile := range h.Tiles {
			if tile.Position.Z != level {
				continue
			}
			posX := tile.Position.X - houseTopX
			posY := tile.Position.Y - houseTopY
			houseImage.Set(int(posX), int(posY), color.RGBA{192, 192, 192, 255})
		}
		houseFile, err := os.Create(fmt.Sprintf(
			"%v_%v.png",
			savepath,
			level,
		))
		if err != nil {
			return err
		}
		if err := png.Encode(houseFile, houseImage); err != nil {
			return err
		}
		defer houseFile.Close()
	}
	return nil
}
