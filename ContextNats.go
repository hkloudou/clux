package clux

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"

	"github.com/nats-io/nats.go"
)

type contextWriterNats struct {
	req *nats.Msg
	// res *nats.Msg
	header url.Values
}

func (m *contextWriterNats) Header() url.Values {
	return m.header
}

func (m *contextWriterNats) Write(data []byte) (int, error) {
	// m.res.Data = data
	// m.req.Header
	// log.Println("write", data)

	msg := nats.NewMsg("")
	msg.Header = nats.Header(m.header)
	msg.Data = data
	// log.Println("topic", m.req)
	// log.Println("topic", m.req.Reply)
	// log.Println("topic header", m.req.Header)
	if err := m.req.RespondMsg(msg); err != nil {
		return 0, err
	}
	return len(msg.Data), nil
}

func (m *contextWriterNats) WriteHeader(statusCode int) {
	m.Header().Set("clux-code", fmt.Sprintf("%d", statusCode))
}

type natsTransportData struct {
	head url.Values
	body io.ReadCloser
}

func (m *natsTransportData) Body() io.ReadCloser {
	return m.body
}

func (m *natsTransportData) Head() url.Values {
	return m.head
}

func newNatsRequestTransportData(req *nats.Msg) *natsTransportData {
	return &natsTransportData{
		head: map[string][]string(req.Header),
		body: ioutil.NopCloser(bytes.NewReader(req.Data)),
	}
}

func NewContext(m *nats.Msg) *Context {
	// http.ResponseWriter
	return &Context{
		Request: newNatsRequestTransportData(m),
		Writer: &contextWriterNats{
			req: m,
			// res: nats.NewMsg(""),
			header: url.Values{},
		},
	}
}
