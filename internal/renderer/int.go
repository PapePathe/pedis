package renderer

import "fmt"

type IntRenderer struct{}

func (ir IntRenderer) Render(name string, data int) []byte {
	bytes := []byte{}
	bytes = append(bytes, []byte(fmt.Sprintf("+%s", name))...)
	bytes = append(bytes, []byte{13, 10}...)
	bytes = append(bytes, []byte(":")...)
	bytes = append(bytes, []byte(fmt.Sprintf("%d", data))...)
	bytes = append(bytes, []byte{13, 10}...)

	return bytes
}
