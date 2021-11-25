package main

import (
	_ "image/jpeg"
	_ "image/png"
	"syscall/js"
	colors "wasm-materialyou/packages"
)

func main() {
	c := make(chan struct{}, 0)
	println("Go WebAssembly initialized")

	js.Global().Set("findCommonColors", js.FuncOf(colors.FindCommonColors))
	<-c
}
