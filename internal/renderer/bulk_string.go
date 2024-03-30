package renderer

import "fmt"

type BulkStringRenderer struct{}

func (sr BulkStringRenderer) Render(data string) []byte {
	bytes := []byte{}
	bytes = append(bytes, []byte(fmt.Sprintf("$%d", len(data)))...)
	bytes = append(bytes, []byte{13, 10}...)
	bytes = append(bytes, []byte(data)...)
	bytes = append(bytes, []byte{13, 10}...)

	return bytes
}
