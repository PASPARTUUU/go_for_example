package tool

import (
	"image"
	"math"

	"github.com/PASPARTUUU/go_for_example/pkg/errpath"
	"github.com/nfnt/resize"
	"github.com/oliamb/cutter"
)

// -------------------------------------------------

// RoundingUp - округляет число x до prec знака после запятой
func RoundingUp(x float64, prec int) float64 {
	var rounder float64
	pow := math.Pow(10, float64(prec))
	intermed := x * pow
	_, frac := math.Modf(intermed)
	if frac >= 0.5 {
		rounder = math.Ceil(intermed)
	} else {
		rounder = math.Floor(intermed)
	}

	return rounder / pow
}

// ResizeAndCutPicture - пропорционально сжимает затем обрезает изображение до нужного размера
func ResizeAndCutPicture(img image.Image, needSize PictureSize) (image.Image, error) {
	var resizeImg image.Image

	widthForResize := float64(img.Bounds().Dx()) / float64(needSize.Width)
	heightForResize := float64(img.Bounds().Dy()) / float64(needSize.Height)

	if widthForResize < heightForResize {
		resizeImg = resize.Resize(uint(needSize.Width), 0, img, resize.Lanczos3)
	} else {
		resizeImg = resize.Resize(0, uint(needSize.Height), img, resize.Lanczos3)
	}

	cutingImg, err := cutter.Crop(resizeImg, cutter.Config{
		Width:  needSize.Width,
		Height: needSize.Height,
		Mode:   cutter.Centered,
	})
	if err != nil {
		return nil, errpath.Err(err)
	}

	return cutingImg, nil
}
