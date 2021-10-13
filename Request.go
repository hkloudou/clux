package clux

import (
	"bytes"
	// "encoding/json"
	"errors"
	"io/ioutil"
	"net/url"
	"time"

	"github.com/hkloudou/clux/binding"
	"github.com/hkloudou/clux/internal/json"
	"github.com/nats-io/nats.go"
)

type natsRequest struct {
	// req    *nats.Msg

	topic   string
	header  url.Values
	timeout time.Duration
	client  *nats.Conn
	res     *Context
}

type Response struct {
	response binding.RequestTransportData
}

func (c *Response) responseHeader(key string) string {
	return c.response.Head().Get(key)
}

// GetHeader returns value from request headers.
func (c *Response) GetHeader(key string) string {
	return c.responseHeader(key)
}

// GetRawData return stream data.
func (c *Response) GetRawData() ([]byte, error) {
	return ioutil.ReadAll(c.response.Body())
}

func (c *Response) ShouldBindHeader(obj interface{}) error {
	return c.ShouldBindWith(obj, binding.Header)
}

func (c *Response) ShouldBindHeaderRaw(obj interface{}) error {
	return c.ShouldBindWith(obj, binding.HeaderRaw)
}
func (c *Response) ShouldBindJson(obj interface{}) error {
	return c.ShouldBindWith(obj, binding.JSON)
}

func (c *Response) ShouldBindWith(obj interface{}, b binding.Binding) error {
	return b.Bind(c.response, obj)
}

type request struct {
	// Writer render.ResponseWriter
	namespace string
	topic     string
	timeout   time.Duration
	client    *nats.Conn
}

func NewRequest(client *nats.Conn, topic string, timeout time.Duration) *request {
	return &request{
		namespace: NameSpace(),
		client:    client,
		topic:     topic,
		timeout:   timeout,
	}
}

func (m *request) GetTopic() string {
	return m.namespace + ".req." + m.topic
}

func (me *request) SetNameSpace(str string) {
	me.namespace = str
}

func (m *request) JSON(header map[string][]string, obj interface{}) (*Response, error) {
	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	return m.Data(header, jsonBytes)
}

func (m *request) Data(header map[string][]string, data []byte) (*Response, error) {
	msg := nats.NewMsg(m.GetTopic())
	msg.Data = data
	if header != nil {
		msg.Header = nats.Header(header)
	}
	if ret, err := m.client.RequestMsg(msg, m.timeout); err != nil {
		return nil, err
	} else if ret == nil {
		return nil, errors.New("error nats response")
	} else {
		if ret.Header.Get("clux-err") != "" {
			return nil, errors.New(ret.Header.Get("clux-err"))
		}
		return &Response{
			response: &natsTransportData{
				head: map[string][]string(ret.Header),
				body: ioutil.NopCloser(bytes.NewReader(ret.Data)),
			},
		}, nil
	}
}

func (m *request) QueueSubscribe(queue string, cb nats.MsgHandler) (*nats.Subscription, error) {
	return m.client.QueueSubscribe(m.GetTopic(), "worker", cb)
}

func (m *request) Subscribe(queue string, cb nats.MsgHandler) (*nats.Subscription, error) {
	return m.client.Subscribe(m.GetTopic(), cb)
}
