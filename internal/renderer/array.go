package renderer

import (
	"fmt"
)

func RenderArray(data []string) []byte {
	bytes := []byte{}
	bytes = append(bytes, []byte(fmt.Sprintf("*%d", len(data)))...)
	bytes = append(bytes, []byte{13, 10}...)

	for _, d := range data {
		bytes = append(bytes, []byte(fmt.Sprintf("$%d", len(d)))...)
		bytes = append(bytes, []byte{13, 10}...)
		bytes = append(bytes, []byte(d)...)
		bytes = append(bytes, []byte{13, 10}...)
	}

	return bytes
}
