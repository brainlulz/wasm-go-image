package main

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/jpeg"
	_ "image/png"
	"syscall/js"

	"github.com/aofei/cameron"
)

func readImage(this js.Value, args []js.Value) interface{} {
	array := args[0]
	buff := make([]uint8, array.Get("byteLength").Int())
	js.CopyBytesToGo(buff, array)

	reader := bytes.NewReader(buff)
	sourceImg, _, err := image.Decode(reader)

	if err != nil {
		println("An error occured while reading the image")
		fmt.Println(err)
	}

	fmt.Println(sourceImg)
	println("Image loaded")
	return nil
	// Now the sourceImg is an image.Image with which we are free to do anything!
}

func identicon(this js.Value, i []js.Value) interface{} {
	array := i[0].String()

	imgIcon := cameron.Identicon(
		[]byte(array),
		540,
		60,
	)
	buf := bytes.Buffer{}
	jpeg.Encode(
		&buf,
		imgIcon,
		&jpeg.Options{
			Quality: 100,
		},
	)

	dst := js.Global().Get("Uint8Array").New(len(buf.Bytes()))
	js.CopyBytesToJS(dst, buf.Bytes())

	return dst
}

func registerCallbacks() {
	emptyWASMObject := make(map[string]interface{})
	js.Global().Set("WASMGo", js.ValueOf(emptyWASMObject))
	js.Global().Get("WASMGo").Set("loadImage", js.FuncOf(readImage))
	js.Global().Get("WASMGo").Set("identicon", js.FuncOf(identicon))
}

func main() {
	c := make(chan struct{}, 0)

	println("WASM Go Initialized")

	registerCallbacks()
	<-c
}
