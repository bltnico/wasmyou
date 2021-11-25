package packages

import (
	"bytes"
	"fmt"
	"image"
	"log"
	"strconv"
	"syscall/js"

	"github.com/EdlinOrg/prominentcolor"
)

type Hex string

type RGB struct {
	Red   uint8
	Green uint8
	Blue  uint8
}

type ImageInfo struct {
	Source image.Image
	Width  int
	Height int
}

func FindCommonColors(this js.Value, args []js.Value) interface{} {
	imageInfo := getImageFromArgs(args)

	handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		resolve := args[0]
		go func() {
			colors := getDominantColors(imageInfo.Source)
			resolve.Invoke(js.ValueOf(colors))
		}()

		return nil
	})

	promiseConstructor := js.Global().Get("Promise")
	return promiseConstructor.New(handler)
}

func getImageFromArgs(args []js.Value) ImageInfo {
	var imageInfo ImageInfo

	array := args[0]
	inBuf := make([]uint8, array.Get("byteLength").Int())
	js.CopyBytesToGo(inBuf, array)

	reader := bytes.NewReader(inBuf)
	sourceImage, _, err := image.Decode(reader)

	if err != nil {
		log.Fatal("Failed to load image", err)
	}

	newReader := bytes.NewReader(inBuf)
	imageConfig, _, configErr := image.DecodeConfig(newReader)

	if configErr != nil {
		log.Fatal("Failed to get image config", configErr)
	}

	imageInfo = ImageInfo{
		Source: sourceImage,
		Width:  imageConfig.Width,
		Height: imageConfig.Height,
	}

	return imageInfo
}

func getDominantColors(image image.Image) []interface{} {
	colors, err := prominentcolor.Kmeans(image)

	if err != nil {
		log.Fatal("Failed to get dominant colors", err)
	}

	colorsRgb := make([]interface{}, 3)

	for i, color := range colors {
		var hex Hex = Hex(color.AsString())
		rgbStr, _, _ := hex.toRGB()
		colorsRgb[i] = rgbStr
	}

	return colorsRgb
}

func (h Hex) toRGB() (string, RGB, error) {
	rgb, err := Hex2RGB(h)
	rgbStr := fmt.Sprintf("rgba(%d,%d,%d,1)", rgb.Red, rgb.Green, rgb.Blue)
	return rgbStr, rgb, err
}

func Hex2RGB(hex Hex) (RGB, error) {
	var rgb RGB
	values, err := strconv.ParseUint(string(hex), 16, 32)

	if err != nil {
		return RGB{}, err
	}

	rgb = RGB{
		Red:   uint8(values >> 16),
		Green: uint8((values >> 8) & 0xFF),
		Blue:  uint8(values & 0xFF),
	}

	return rgb, nil
}

func consoleLog(args ...interface{}) {
	js.Global().Get("console").Call("log", args)
}
