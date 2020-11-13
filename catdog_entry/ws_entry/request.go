package ws_entry

import (
	"github.com/asim/nitro/v3/codec"
)

type httpRequest struct {
	service     string
	method      string
	contentType string
	header      map[string]string
	body        []byte
	stream      bool
}

func (r *httpRequest) ContentType() string {
	return r.contentType
}

func (r *httpRequest) Service() string {
	return r.service
}

func (r *httpRequest) Method() string {
	return r.method
}

func (r *httpRequest) Endpoint() string {
	return r.method
}

func (r *httpRequest) Codec() codec.Reader {
	return nil
}

func (r *httpRequest) Header() map[string]string {
	return r.header
}

func (r *httpRequest) Read() ([]byte, error) {
	return r.body, nil
}

func (r *httpRequest) Stream() bool {
	return r.stream
}

func (r *httpRequest) Body() interface{} {
	return r.body
}
