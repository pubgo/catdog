package catdog_request

import (
	"github.com/asim/nitro/v3/codec"
)

type Request struct {
	service     string
	method      string
	contentType string
	header      map[string]string
	body        []byte
	stream      bool
}

func (r *Request) ContentType() string {
	return r.contentType
}

func (r *Request) Service() string {
	return r.service
}

func (r *Request) Method() string {
	return r.method
}

func (r *Request) Endpoint() string {
	return r.method
}

func (r *Request) Codec() codec.Reader {
	return nil
}

func (r *Request) Header() map[string]string {
	return r.header
}

func (r *Request) Read() ([]byte, error) {
	return r.body, nil
}

func (r *Request) Stream() bool {
	return r.stream
}

func (r *Request) Body() interface{} {
	return r.body
}
