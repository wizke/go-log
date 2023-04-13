package log

import "fmt"

func DebugByteListHex(byteList []byte) {
	str := "["
	for i, b := range byteList {
		str += fmt.Sprintf(" 0x%02X", b)
		if i == len(byteList)-1 {
			str += " "
		}
	}
	str += "]"
	Debug(str)
}
