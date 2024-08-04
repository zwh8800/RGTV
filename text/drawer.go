package text

import (
	"image"
	"image/draw"
	"os"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/puzpuzpuz/xsync/v3"
	"golang.org/x/image/font"
)

var (
	drawers = xsync.NewMapOf[string, *Drawer]()
)

func GetDrawerFromFile(fontFile string) (*Drawer, error) {
	var loadErr error
	drawer, _ := drawers.LoadOrCompute(fontFile, func() *Drawer {
		drawer, err := NewDrawerFromFile(fontFile)
		if err != nil {
			loadErr = err
			return nil
		}
		return drawer
	})
	return drawer, loadErr
}

func GetDrawerFromData(dataName string, data []byte) (*Drawer, error) {
	var loadErr error
	drawer, _ := drawers.LoadOrCompute(dataName, func() *Drawer {
		drawer, err := NewDrawerFromData(data)
		if err != nil {
			loadErr = err
			return nil
		}
		return drawer
	})
	return drawer, loadErr
}

type Drawer struct {
	font *truetype.Font
}

func NewDrawerFromFile(fontFile string) (*Drawer, error) {
	fontBytes, err := os.ReadFile(fontFile)
	if err != nil {
		return nil, err
	}

	return NewDrawerFromData(fontBytes)
}

func NewDrawerFromData(data []byte) (*Drawer, error) {
	f, err := freetype.ParseFont(data)
	if err != nil {
		return nil, err
	}
	return &Drawer{font: f}, nil
}

// MeasureTextSize 测量给定文本和字体大小所需的图像宽度和高度
func (d *Drawer) MeasureTextSize(text string, fontSize float64) (int, int, error) {
	pt := freetype.Pt(0, int(fontSize))

	// 使用 Drawer 测量文本尺寸
	drawer := &font.Drawer{
		Face: truetype.NewFace(d.font, &truetype.Options{
			Size: fontSize,
		}),
		Dot: pt,
	}

	// 计算文本的边界框
	bounds, _ := drawer.BoundString(text)
	width := (bounds.Max.X - bounds.Min.X).Ceil()
	height := (bounds.Max.Y - bounds.Min.Y).Ceil()

	return width, height, nil
}

func (d *Drawer) Draw(text string, fontSize float64, color *image.Uniform) (*image.RGBA, error) {
	// Create a freetype context.
	fc := freetype.NewContext()
	fc.SetDPI(72)
	fc.SetFont(d.font)
	fc.SetFontSize(fontSize)
	fc.SetSrc(color)

	// Calculate the bounds of the text.
	face := truetype.NewFace(d.font, &truetype.Options{
		Size: fontSize,
		DPI:  72,
	})
	bounds, _ := font.BoundString(face, text)
	width := (bounds.Max.X - bounds.Min.X).Ceil()
	height := (bounds.Max.Y - bounds.Min.Y).Ceil()

	// Create a new image with the exact size of the text.
	rgba := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(rgba, rgba.Bounds(), image.Transparent, image.Point{}, draw.Src)

	// Set the destination image and draw the text.
	fc.SetDst(rgba)
	fc.SetClip(rgba.Bounds())

	pt := freetype.Pt(-bounds.Min.X.Ceil(), -bounds.Min.Y.Ceil())
	_, err := fc.DrawString(text, pt)
	if err != nil {
		return nil, err
	}

	return rgba, nil
}
