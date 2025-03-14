package iec

import (
	"log"
	"strings"
)

func DebugCpuWrite(data uint8) {
	//value := ^data
	value := data
	var message []string
	if value&0x20 != 0 {
		message = append(message, "[DATA_OUT]")
	}
	if value&0x10 != 0 {
		message = append(message, "[CLK_OUT]")
	}
	if value&0x08 != 0 {
		message = append(message, "[ATN_OUT]")
	}
	log.Printf("CPU SEND: [%x] [%08b] %s\n", value, value, strings.Join(message, " "))
}

func DebugCpuRead(data uint8) {
	value := data
	var message []string
	if value&0x80 != 0 {
		message = append(message, "[CLK_IN]")
	}
	if value&0x40 != 0 {
		message = append(message, "[DATA_IN]")
	}
	if value&0x20 != 0 {
		message = append(message, "[DATA_OUT]")
	}
	if value&0x10 != 0 {
		message = append(message, "[CLK_OUT]")
	}
	if value&0x08 != 0 {
		message = append(message, "[ATN_OUT]")
	}
	log.Printf("CPU SEND: [%x] [%08b] %s\n", value, value, strings.Join(message, " "))
}
