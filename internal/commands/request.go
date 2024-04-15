package commands

import (
	"bytes"
	"fmt"
	"regexp"
)

type RawRequest []byte

func (r RawRequest) String() string {
	return fmt.Sprintf("%q", string(r))
}

func (r RawRequest) Parse() []string {
	re := regexp.MustCompile(`\*\d\r\n`)

	chunks := re.Split(string(r), -1)

	if len(chunks) == 1 {
		return []string{"*x\r\n" + string(r)}
	}

	cmds := []string{}

	for _, ch := range chunks[1:] {
		cmds = append(cmds, "*x\r\n"+ch)
	}

	return cmds
}

func (r RawRequest) ReadArray() []string {
	items := bytes.Split(r[2:], []byte{13, 10})
	sl := SliceAsChunks(items[3:], 2)
	array := []string{}

	for _, i := range sl {
		if len(i) == 2 {
			array = append(array, string(i[1]))
		}
	}

	return array
}

func SliceAsChunks(slice [][]byte, chunkSize int) [][][]byte {
	var chunks [][][]byte

	for i := 0; i < len(slice); i += chunkSize {
		end := i + chunkSize
		if end > len(slice) {
			end = len(slice)
		}

		chunks = append(chunks, slice[i:end])
	}

	return chunks
}
