package response

import "pedis/internal/renderer"

type HelloResponse struct {
	Server  string `resp:"string"`
	Version string `resp:"string"`
	Mode    string `resp:"string"`
	Proto   int    `resp:"int"`
	Id      int    `resp:"int"`
	Role    string `resp:"string"`
}

func (hr HelloResponse) Render() []byte {
	resp := renderer.RespRenderer{Tagname: "resp"}
	return resp.Render(hr)
}
