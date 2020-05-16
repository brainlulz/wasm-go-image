package main

import (
	"syscall/js"
	"bytes"
	"image"
	_ "image/png"
	_ "image/jpeg"
	"fmt"
)

func readImage(this js.Value, args []js.Value) interface{} {
	array := args[0]
	buff  := make([]uint8, array.Get("byteLength").Int())
	js.CopyBytesToGo(buff, array)


	reader := bytes.NewReader(buff)
	_, _, err := image.Decode(reader)


	if err != nil {
		println("An error occured while reading the image")
		fmt.Println(err)
	}

	println("Image loaded")
	return nil
	// Now the sourceImg is an image.Image with which we are free to do anything!
}

func registerCallbacks() {
	emptyWASMObject := make(map[string]interface{})
	js.Global().Set("WASMGo", js.ValueOf(emptyWASMObject))
	js.Global().Get("WASMGo").Set("loadImage", js.FuncOf(readImage))
}

func main() {
	c := make(chan struct{}, 0)

	println("WASM Go Initialized")

	registerCallbacks()
	<-c
}
