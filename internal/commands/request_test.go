package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequestParsing(t *testing.T) {
	cmd := "*5\r\n$3\r\nset\r\n$9\r\nkey:90111\r\n$10\r\nsome\r\nvalue\r\n$2\r\npx\r\n$1\r\n1\r\n*3\r\n$7\r\nhexists\r\n$9\r\nuser:9011\r\n$3\r\n404\r\n"
	r := RawRequest([]byte(cmd))

	expected := []string{
		"*x\r\n$3\r\nset\r\n$9\r\nkey:90111\r\n$10\r\nsome\r\nvalue\r\n$2\r\npx\r\n$1\r\n1\r\n",
		"*x\r\n$7\r\nhexists\r\n$9\r\nuser:9011\r\n$3\r\n404\r\n",
	}
	assert.Equal(t, expected, r.Parse())
}
