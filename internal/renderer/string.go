package renderer

import "fmt"

type StringRenderer struct{}

func (sr StringRenderer) Render(name string, data string) []byte {
	bytes := []byte{}
	bytes = append(bytes, []byte(fmt.Sprintf("+%s", name))...)
	bytes = append(bytes, []byte{13, 10}...)
	bytes = append(bytes, []byte("+")...)
	bytes = append(bytes, []byte(data)...)
	bytes = append(bytes, []byte{13, 10}...)

	return bytes
}
