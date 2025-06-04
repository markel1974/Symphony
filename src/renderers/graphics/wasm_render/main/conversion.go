//go:build js && wasm

package main // o il tuo package wasm

import (
	"fmt"
	"syscall/js"
)

func JsBooleanToGoBool(jsValue js.Value) (bool, error) {
	if jsValue.IsUndefined() {
		return false, fmt.Errorf("undefined")
	}
	if jsValue.Type() != js.TypeBoolean {
		return false, fmt.Errorf("invalid type")
	}
	return jsValue.Bool(), nil
}

func JsStringToGoString(jsValue js.Value) (string, error) {
	if jsValue.IsUndefined() {
		return "", fmt.Errorf("undefined")
	}
	if jsValue.Type() != js.TypeString {
		return "", fmt.Errorf("invalid type")
	}
	return jsValue.String(), nil
}

func JsBufferToGoBytes(jsValue js.Value) ([]byte, error) {
	if jsValue.IsUndefined() || jsValue.IsNull() {
		return nil, fmt.Errorf("undefined")
	}
	var jsTypedArray js.Value
	global := js.Global()
	arrayBufferConstructor := global.Get("ArrayBuffer")
	if !arrayBufferConstructor.IsUndefined() && jsValue.InstanceOf(arrayBufferConstructor) {
		uint8ArrayConstructor := global.Get("Uint8Array")
		if uint8ArrayConstructor.IsUndefined() {
			return nil, fmt.Errorf("constructor Uint8Array not found")
		}
		jsTypedArray = uint8ArrayConstructor.New(jsValue)
	} else if !jsValue.Get("byteLength").IsUndefined() {
		jsTypedArray = jsValue
	} else {
		return nil, fmt.Errorf("incorrect value")
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
