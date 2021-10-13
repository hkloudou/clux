package clux

import (
	"errors"

	"github.com/hkloudou/clux/internal/json"
	"github.com/nats-io/nats.go"
)

type event struct {
	namespace string
	topic     string
}

func NewEvent(topic string) *event {
	return &event{
		namespace: NameSpace(),
		topic:     topic,
	}
}

func (me *event) SetNameSpace(str string) {
	me.namespace = str
}

func (me *event) GetTopic() string {
	return me.namespace + ".evt." + me.topic
}

func (me *event) PublishData(head map[string][]string, data []byte, opts ...nats.PubOpt) (*nats.PubAck, error) {
	if _defaultJsClient == nil {
		return nil, errors.New("client is null")
	}
	msg := nats.NewMsg(me.GetTopic())
	msg.Header = head
	msg.Data = data
	return _defaultJsClient.PublishMsg(msg)
}

func (me *event) PublishJson(head map[string][]string, obj interface{}, opts ...nats.PubOpt) (*nats.PubAck, error) {
	if _defaultJsClient == nil {
		return nil, errors.New("client is null")
	}
	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	return me.PublishData(head, jsonBytes, opts...)
}

func (me *event) PublishDataAsync(head map[string][]string, data []byte, opts ...nats.PubOpt) (nats.PubAckFuture, error) {
	if _defaultJsClient == nil {
		return nil, errors.New("client is null")
	}
	msg := nats.NewMsg(me.GetTopic())
	msg.Header = head
	msg.Data = data
	return _defaultJsClient.PublishMsgAsync(msg)
}

func (me *event) PublishJsonAsync(head map[string][]string, obj interface{}, opts ...nats.PubOpt) (nats.PubAckFuture, error) {
	if _defaultJsClient == nil {
		return nil, errors.New("client is null")
	}
	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	return me.PublishDataAsync(head, jsonBytes, opts...)
}
