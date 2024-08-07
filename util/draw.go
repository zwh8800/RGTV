package util

import (
	"image"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
)

// DrawGoImage 将golang的image绘制到sdl.Renderer上
func DrawGoImage(renderer *sdl.Renderer, img *image.RGBA, rectangle image.Rectangle) error {
	// 将image.RGBA转换为sdl.Texture
	texture, err := renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STATIC,
		int32(img.Bounds().Dx()), int32(img.Bounds().Dy()))
	if err != nil {
		return err
	}
	defer texture.Destroy()

	texture.SetBlendMode(sdl.BLENDMODE_BLEND)

	// 更新纹理数据
	err = texture.Update(nil, unsafe.Pointer(&img.Pix[0]), img.Stride)
	if err != nil {
		return err
	}
	// 绘制纹理
	err = renderer.Copy(texture, nil, &sdl.Rect{
		X: int32(rectangle.Min.X),
		Y: int32(rectangle.Min.Y),
		W: int32(rectangle.Dx()),
		H: int32(rectangle.Dy()),
	})
	if err != nil {
		return err
	}

	return nil
}
