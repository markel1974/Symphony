//go:build js && wasm

package main // o il tuo package wasm

import (
	"fmt"
	"syscall/js"
)

func ConvertJsBufferToGoBytes(jsBufferValue js.Value) ([]byte, error) {
	if jsBufferValue.IsUndefined() || jsBufferValue.IsNull() {
		return nil, fmt.Errorf("undefined js.Value")
	}
	var jsTypedArray js.Value
	global := js.Global()
	arrayBufferConstructor := global.Get("ArrayBuffer")
	if !arrayBufferConstructor.IsUndefined() && jsBufferValue.InstanceOf(arrayBufferConstructor) {
		uint8ArrayConstructor := global.Get("Uint8Array")
		if uint8ArrayConstructor.IsUndefined() {
			return nil, fmt.Errorf("constructor Uint8Array not found")
		}
		jsTypedArray = uint8ArrayConstructor.New(jsBufferValue)
	} else if !jsBufferValue.Get("byteLength").IsUndefined() {
		jsTypedArray = jsBufferValue
	} else {
		return nil, fmt.Errorf("js value is incorrect")
	}
	length := jsTypedArray.Get("byteLength").Int()
	if length == 0 {
		return []byte{}, nil // Restituisce una slice vuota se la lunghezza è 0
	}

	goByteSlice := make([]byte, length)
	copiedBytes := js.CopyBytesToGo(goByteSlice, jsTypedArray)
	if copiedBytes != length {
		return nil, fmt.Errorf("errore while copying data (copied %d, real %d)", copiedBytes, length)
	}
	return goByteSlice, nil
}
