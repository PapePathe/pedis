package response

type Response interface {
	Render() []byte
}
