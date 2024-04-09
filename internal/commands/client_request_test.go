package commands

import (
	"log"
	"net"
	"testing"
	"time"
)

type store struct {
	data []byte
}

func (s *store) write(b []byte) {
	s.data = b
	log.Println("w", s.data)
}

func (s *store) read() []byte {
	log.Println("r", s.data)
	return s.data
}

type TestConn struct {
	s *store
}

func (tc TestConn) Read(b []byte) (n int, err error) {
	b = tc.s.read()
	return 0, nil
}
func (tc TestConn) Write(b []byte) (n int, err error) {
	tc.s.write(b)
	return 0, nil
}
func (tc TestConn) Close() error {
	return nil
}
func (tc TestConn) LocalAddr() net.Addr {
	return nil
}
func (tc TestConn) RemoteAddr() net.Addr {
	return nil
}
func (tc TestConn) SetDeadline(t time.Time) error {
	return nil
}
func (tc TestConn) SetReadDeadline(t time.Time) error {
	return nil
}
func (tc TestConn) SetWriteDeadline(t time.Time) error {
	return nil
}

func TestClientRequestWriteError(t *testing.T) {
	// t.Parallel()
	//
	//	cli := ClientRequest{
	//		Conn: TestConn{s: &store{}},
	//	}
	//
	// err := cli.WriteError("key not found")
	// require.NoError(t, err)
	//
	// buf := make([]byte, 2048)
	// n, err := cli.Conn.Read(buf)
	// t.Log("test", n)
	// require.NoError(t, err)
	//
	// assert.Equal(t, buf[0:n], []byte("-ERR key not found"))
}
